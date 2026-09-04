package githubcli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// IssueView represents the inspected state of an issue.
type IssueView struct {
	Number       int               `json:"number"`
	Title        string            `json:"title"`
	Body         string            `json:"body"`
	State        string            `json:"state"`
	StateReason  string            `json:"stateReason,omitempty"`
	URL          string            `json:"url"`
	Labels       []IssueLabel      `json:"labels"`
	Assignees    []IssueAssignee   `json:"assignees"`
	Milestone    *IssueMilestone   `json:"milestone,omitempty"`
	IssueType    *IssueType        `json:"issueType,omitempty"`
	ProjectItems []json.RawMessage `json:"projectItems,omitempty"`
}

// IssueLabel is an issue label name.
type IssueLabel struct {
	Name string `json:"name"`
}

// IssueAssignee is an assigned user's login.
type IssueAssignee struct {
	Login string `json:"login"`
}

// IssueMilestone is an assigned milestone title.
type IssueMilestone struct {
	Title string `json:"title"`
}

// IssueType is an organisation-native GitHub issue type.
type IssueType struct {
	Name string `json:"name"`
}

// IssueSummary is the stable subset used for exact-title duplicate checks.
type IssueSummary struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	State  string `json:"state"`
	URL    string `json:"url"`
}

// CreateIssueInput holds parameters for creating an issue.
type CreateIssueInput struct {
	Repo      string
	Title     string
	Body      string
	Labels    []string
	Assignees []string
	Milestone string
}

// EditIssueInput holds parameters for editing an existing issue.
type EditIssueInput struct {
	Repo            string
	Number          int
	Title           *string
	Body            *string
	AddLabels       []string
	RemoveLabels    []string
	AddAssignees    []string
	RemoveAssignees []string
	Milestone       *string
	IssueType       *string // empty removes the current issue type
	State           *string // "open" or "closed"
	CloseReason     string  // "completed" or "not_planned"
}

// ViewIssue inspects an issue using gh issue view --json.
func ViewIssue(ctx context.Context, runner Runner, repo string, number int) (IssueView, error) {
	output, err := runner.Run(
		ctx,
		"issue", "view", strconv.Itoa(number),
		"--repo", repo,
		"--json", "number,title,body,state,stateReason,labels,assignees,milestone,issueType,projectItems,url",
	)
	if err != nil {
		return IssueView{}, fmt.Errorf("view issue %s#%d: %w", repo, number, err)
	}
	var view IssueView
	if err := json.Unmarshal(output, &view); err != nil {
		return IssueView{}, fmt.Errorf("decode issue view: %w", err)
	}
	return view, nil
}

// FindIssuesByExactTitle scans the complete issue repository and returns
// issues whose title exactly matches title. REST's issue collection also
// contains pull requests, which are deliberately excluded.
func FindIssuesByExactTitle(ctx context.Context, runner Runner, repo, title string) ([]IssueSummary, error) {
	filter := fmt.Sprintf(
		`.[] | select((.pull_request == null) and (.title == %s)) | {number, title, state, url: .html_url}`,
		strconv.Quote(title),
	)
	args := []string{"api", "--paginate"}
	args = append(args, apiHeaders()...)
	args = append(args, fmt.Sprintf("repos/%s/issues?state=all&per_page=100", repo), "--jq", filter)
	output, err := runner.Run(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("scan issues in %s for an exact-title match: %w", repo, err)
	}
	matches := make([]IssueSummary, 0)
	decoder := json.NewDecoder(strings.NewReader(string(output)))
	for {
		var issue IssueSummary
		if err := decoder.Decode(&issue); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("decode exact-title issue results: %w", err)
		}
		matches = append(matches, issue)
	}
	sort.Slice(matches, func(i, j int) bool { return matches[i].Number < matches[j].Number })
	return matches, nil
}

// CreateIssue creates an issue on GitHub and independently verifies it via readback.
func CreateIssue(ctx context.Context, runner Runner, input CreateIssueInput) (IssueView, error) {
	if input.Repo == "" {
		return IssueView{}, errors.New("issue creation requires a repository")
	}
	if strings.TrimSpace(input.Title) == "" {
		return IssueView{}, errors.New("issue creation requires a non-empty title")
	}

	args := []string{"issue", "create", "--repo", input.Repo, "--title", input.Title}
	if input.Body != "" {
		args = append(args, "--body", input.Body)
	}
	for _, l := range input.Labels {
		if strings.TrimSpace(l) != "" {
			args = append(args, "--label", strings.TrimSpace(l))
		}
	}
	for _, a := range input.Assignees {
		if strings.TrimSpace(a) != "" {
			args = append(args, "--assignee", strings.TrimSpace(a))
		}
	}
	if input.Milestone != "" {
		args = append(args, "--milestone", input.Milestone)
	}

	output, err := runner.Run(ctx, args...)
	if err != nil {
		return IssueView{}, fmt.Errorf("create issue in %s: %w", input.Repo, err)
	}

	target, err := ResolveGitHubItemTarget(input.Repo, 0, strings.TrimSpace(string(output)))
	if err != nil {
		return IssueView{}, fmt.Errorf("could not resolve created issue identity from gh output %q: %w", string(output), err)
	}
	if target.Kind != "issues" {
		return IssueView{}, fmt.Errorf("created issue command returned a non-issue URL: %s", target.URL)
	}

	view, err := ViewIssue(ctx, runner, input.Repo, target.Number)
	if err != nil {
		return IssueView{}, fmt.Errorf("read back created issue %s#%d: %w", input.Repo, target.Number, err)
	}

	if err := verifyCreatedIssue(view, input); err != nil {
		return IssueView{}, err
	}

	return view, nil
}

// EditIssue modifies an issue and independently verifies applied changes and preservation.
func EditIssue(ctx context.Context, runner Runner, input EditIssueInput) (IssueView, error) {
	if input.Repo == "" {
		return IssueView{}, errors.New("issue edit requires a repository")
	}
	if input.Number <= 0 {
		return IssueView{}, errors.New("issue edit requires a positive issue number")
	}
	if err := validateIssueEditInput(input); err != nil {
		return IssueView{}, err
	}

	before, err := ViewIssue(ctx, runner, input.Repo, input.Number)
	if err != nil {
		return IssueView{}, fmt.Errorf("inspect issue before edit: %w", err)
	}

	editArgs := []string{"issue", "edit", strconv.Itoa(input.Number), "--repo", input.Repo}
	hasEdits := false

	if input.Title != nil && *input.Title != before.Title {
		editArgs = append(editArgs, "--title", *input.Title)
		hasEdits = true
	}
	if input.Body != nil && *input.Body != before.Body {
		editArgs = append(editArgs, "--body", *input.Body)
		hasEdits = true
	}
	beforeLabels := issueLabelNames(before.Labels)
	for _, l := range cleanNames(input.AddLabels) {
		if !containsName(beforeLabels, l) {
			editArgs = append(editArgs, "--add-label", l)
			hasEdits = true
		}
	}
	for _, l := range cleanNames(input.RemoveLabels) {
		if containsName(beforeLabels, l) {
			editArgs = append(editArgs, "--remove-label", l)
			hasEdits = true
		}
	}
	beforeAssignees := issueAssigneeNames(before.Assignees)
	for _, a := range cleanNames(input.AddAssignees) {
		if !containsName(beforeAssignees, a) {
			editArgs = append(editArgs, "--add-assignee", a)
			hasEdits = true
		}
	}
	for _, a := range cleanNames(input.RemoveAssignees) {
		if containsName(beforeAssignees, a) {
			editArgs = append(editArgs, "--remove-assignee", a)
			hasEdits = true
		}
	}
	if input.Milestone != nil {
		if *input.Milestone == "" {
			if before.Milestone != nil {
				editArgs = append(editArgs, "--remove-milestone")
				hasEdits = true
			}
		} else if before.Milestone == nil || before.Milestone.Title != *input.Milestone {
			editArgs = append(editArgs, "--milestone", *input.Milestone)
			hasEdits = true
		}
	}
	if input.IssueType != nil {
		requestedType := strings.TrimSpace(*input.IssueType)
		if requestedType == "" {
			if before.IssueType != nil {
				editArgs = append(editArgs, "--remove-type")
				hasEdits = true
			}
		} else if before.IssueType == nil || !strings.EqualFold(before.IssueType.Name, requestedType) {
			editArgs = append(editArgs, "--type", requestedType)
			hasEdits = true
		}
	}

	if hasEdits {
		if _, err := runner.Run(ctx, editArgs...); err != nil {
			return IssueView{}, fmt.Errorf("apply issue edits: %w", err)
		}
	}

	if input.State != nil {
		targetState := strings.ToUpper(strings.TrimSpace(*input.State))
		if targetState == "CLOSED" && strings.EqualFold(before.State, "OPEN") {
			closeArgs := []string{"issue", "close", strconv.Itoa(input.Number), "--repo", input.Repo}
			if input.CloseReason != "" {
				closeArgs = append(closeArgs, "--reason", input.CloseReason)
			}
			if _, err := runner.Run(ctx, closeArgs...); err != nil {
				return IssueView{}, fmt.Errorf("close issue %s#%d: %w", input.Repo, input.Number, err)
			}
		} else if targetState == "OPEN" && strings.EqualFold(before.State, "CLOSED") {
			reopenArgs := []string{"issue", "reopen", strconv.Itoa(input.Number), "--repo", input.Repo}
			if _, err := runner.Run(ctx, reopenArgs...); err != nil {
				return IssueView{}, fmt.Errorf("reopen issue %s#%d: %w", input.Repo, input.Number, err)
			}
		}
	}

	after, err := ViewIssue(ctx, runner, input.Repo, input.Number)
	if err != nil {
		return IssueView{}, fmt.Errorf("read back edited issue: %w", err)
	}

	if err := verifyEditedIssue(before, after, input); err != nil {
		return IssueView{}, err
	}

	return after, nil
}

func verifyCreatedIssue(view IssueView, input CreateIssueInput) error {
	if view.Number <= 0 || strings.TrimSpace(view.URL) == "" {
		return errors.New("created issue readback is missing its stable number or URL")
	}
	if view.Title != input.Title {
		return fmt.Errorf("created issue title readback disagrees: got %q, want %q", view.Title, input.Title)
	}
	if view.Body != input.Body {
		return errors.New("created issue body readback disagrees with the requested body")
	}
	if !strings.EqualFold(view.State, "OPEN") {
		return fmt.Errorf("created issue state is %s, want OPEN", view.State)
	}
	if err := compareNameSets("created issue labels", issueLabelNames(view.Labels), input.Labels); err != nil {
		return err
	}
	if err := compareNameSets("created issue assignees", issueAssigneeNames(view.Assignees), input.Assignees); err != nil {
		return err
	}
	if input.Milestone == "" {
		if view.Milestone != nil {
			return fmt.Errorf("created issue unexpectedly has milestone %q", view.Milestone.Title)
		}
	} else if view.Milestone == nil || view.Milestone.Title != input.Milestone {
		observed := ""
		if view.Milestone != nil {
			observed = view.Milestone.Title
		}
		return fmt.Errorf("created issue milestone readback disagrees: got %q, want %q", observed, input.Milestone)
	}
	return nil
}

func validateIssueEditInput(input EditIssueInput) error {
	if overlap := overlappingNames(input.AddLabels, input.RemoveLabels); len(overlap) > 0 {
		return fmt.Errorf("labels cannot be both added and removed: %s", strings.Join(overlap, ", "))
	}
	if overlap := overlappingNames(input.AddAssignees, input.RemoveAssignees); len(overlap) > 0 {
		return fmt.Errorf("assignees cannot be both added and removed: %s", strings.Join(overlap, ", "))
	}
	if input.State != nil {
		switch strings.ToLower(strings.TrimSpace(*input.State)) {
		case "open", "closed":
		default:
			return fmt.Errorf("invalid issue state %q; expected open or closed", *input.State)
		}
	}
	if input.CloseReason != "" {
		switch strings.ToLower(strings.TrimSpace(input.CloseReason)) {
		case "completed", "not_planned":
		default:
			return fmt.Errorf("invalid close reason %q; expected completed or not_planned", input.CloseReason)
		}
	}
	return nil
}

func verifyEditedIssue(before, after IssueView, input EditIssueInput) error {
	if after.Number != before.Number || after.URL != before.URL {
		return fmt.Errorf("issue identity changed during edit: got #%d %s, want #%d %s", after.Number, after.URL, before.Number, before.URL)
	}

	wantTitle := before.Title
	if input.Title != nil {
		wantTitle = *input.Title
	}
	if after.Title != wantTitle {
		return fmt.Errorf("title readback disagrees: got %q, want %q", after.Title, wantTitle)
	}
	wantBody := before.Body
	if input.Body != nil {
		wantBody = *input.Body
	}
	if after.Body != wantBody {
		return errors.New("body readback disagrees with the requested or preserved body")
	}

	wantLabels := applyNameDelta(issueLabelNames(before.Labels), input.AddLabels, input.RemoveLabels)
	if err := compareNameSets("labels", issueLabelNames(after.Labels), wantLabels); err != nil {
		return err
	}
	wantAssignees := applyNameDelta(issueAssigneeNames(before.Assignees), input.AddAssignees, input.RemoveAssignees)
	if err := compareNameSets("assignees", issueAssigneeNames(after.Assignees), wantAssignees); err != nil {
		return err
	}

	wantMilestone := ""
	if before.Milestone != nil {
		wantMilestone = before.Milestone.Title
	}
	if input.Milestone != nil {
		wantMilestone = *input.Milestone
	}
	gotMilestone := ""
	if after.Milestone != nil {
		gotMilestone = after.Milestone.Title
	}
	if gotMilestone != wantMilestone {
		return fmt.Errorf("milestone readback disagrees: got %q, want %q", gotMilestone, wantMilestone)
	}

	wantType := ""
	if before.IssueType != nil {
		wantType = before.IssueType.Name
	}
	if input.IssueType != nil {
		wantType = strings.TrimSpace(*input.IssueType)
	}
	gotType := ""
	if after.IssueType != nil {
		gotType = after.IssueType.Name
	}
	if !strings.EqualFold(gotType, wantType) {
		return fmt.Errorf("issue type readback disagrees: got %q, want %q", gotType, wantType)
	}

	wantState := before.State
	if input.State != nil {
		wantState = *input.State
	}
	if !strings.EqualFold(after.State, wantState) {
		return fmt.Errorf("state readback disagrees: got %q, want %q", after.State, wantState)
	}
	if input.State == nil && after.StateReason != before.StateReason {
		return fmt.Errorf("unrequested state reason changed: got %q, want preserved %q", after.StateReason, before.StateReason)
	}
	if input.State != nil && strings.EqualFold(*input.State, "closed") && strings.EqualFold(before.State, "open") {
		wantReason := strings.ToUpper(strings.TrimSpace(input.CloseReason))
		if wantReason == "" {
			wantReason = "COMPLETED"
		}
		if !strings.EqualFold(after.StateReason, wantReason) {
			return fmt.Errorf("close reason readback disagrees: got %q, want %q", after.StateReason, wantReason)
		}
	}

	if !reflect.DeepEqual(canonicalRawSet(before.ProjectItems), canonicalRawSet(after.ProjectItems)) {
		return errors.New("unrequested Project membership or Project item summary changed during issue edit")
	}
	return nil
}

func issueLabelNames(labels []IssueLabel) []string {
	result := make([]string, 0, len(labels))
	for _, label := range labels {
		result = append(result, label.Name)
	}
	return result
}

func issueAssigneeNames(assignees []IssueAssignee) []string {
	result := make([]string, 0, len(assignees))
	for _, assignee := range assignees {
		result = append(result, assignee.Login)
	}
	return result
}

func cleanNames(values []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		key := strings.ToLower(trimmed)
		if trimmed == "" || seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, trimmed)
	}
	return result
}

func containsName(values []string, candidate string) bool {
	for _, value := range values {
		if strings.EqualFold(strings.TrimSpace(value), strings.TrimSpace(candidate)) {
			return true
		}
	}
	return false
}

func overlappingNames(left, right []string) []string {
	var overlap []string
	for _, value := range cleanNames(left) {
		if containsName(right, value) {
			overlap = append(overlap, value)
		}
	}
	sort.Strings(overlap)
	return overlap
}

func applyNameDelta(before, additions, removals []string) []string {
	result := cleanNames(before)
	for _, removal := range cleanNames(removals) {
		filtered := result[:0]
		for _, value := range result {
			if !strings.EqualFold(value, removal) {
				filtered = append(filtered, value)
			}
		}
		result = filtered
	}
	for _, addition := range cleanNames(additions) {
		if !containsName(result, addition) {
			result = append(result, addition)
		}
	}
	return result
}

func compareNameSets(subject string, got, want []string) error {
	canonical := func(values []string) []string {
		result := make([]string, 0, len(values))
		for _, value := range cleanNames(values) {
			result = append(result, strings.ToLower(value))
		}
		sort.Strings(result)
		return result
	}
	gotCanonical := canonical(got)
	wantCanonical := canonical(want)
	if !reflect.DeepEqual(gotCanonical, wantCanonical) {
		return fmt.Errorf("%s readback disagrees: got %v, want %v", subject, got, want)
	}
	return nil
}

func canonicalRawSet(values []json.RawMessage) []string {
	result := make([]string, 0, len(values))
	for _, raw := range values {
		var decoded any
		if err := json.Unmarshal(raw, &decoded); err != nil {
			result = append(result, string(raw))
			continue
		}
		encoded, err := json.Marshal(decoded)
		if err != nil {
			result = append(result, string(raw))
			continue
		}
		result = append(result, string(encoded))
	}
	sort.Strings(result)
	return result
}
