package githubcli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var issueURLPattern = regexp.MustCompile(`/issues/(\d+)\s*$`)

// IssueView represents the inspected state of an issue.
type IssueView struct {
	Number    int             `json:"number"`
	Title     string          `json:"title"`
	Body      string          `json:"body"`
	State     string          `json:"state"`
	URL       string          `json:"url"`
	Labels    []IssueLabel    `json:"labels"`
	Assignees []IssueAssignee `json:"assignees"`
	Milestone *IssueMilestone `json:"milestone,omitempty"`
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
	State           *string // "open" or "closed"
	CloseReason     string  // "completed" or "not_planned"
}

// ViewIssue inspects an issue using gh issue view --json.
func ViewIssue(ctx context.Context, runner Runner, repo string, number int) (IssueView, error) {
	output, err := runner.Run(
		ctx,
		"issue", "view", strconv.Itoa(number),
		"--repo", repo,
		"--json", "number,title,body,state,labels,assignees,milestone,url",
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

	matches := issueURLPattern.FindStringSubmatch(strings.TrimSpace(string(output)))
	if len(matches) < 2 {
		return IssueView{}, fmt.Errorf("could not parse issue number from gh output: %q", string(output))
	}
	number, err := strconv.Atoi(matches[1])
	if err != nil {
		return IssueView{}, fmt.Errorf("invalid parsed issue number %q: %w", matches[1], err)
	}

	view, err := ViewIssue(ctx, runner, input.Repo, number)
	if err != nil {
		return IssueView{}, fmt.Errorf("read back created issue %s#%d: %w", input.Repo, number, err)
	}

	if view.Title != input.Title {
		return IssueView{}, fmt.Errorf("read back title disagrees: got %q, want %q", view.Title, input.Title)
	}
	if !strings.EqualFold(view.State, "OPEN") {
		return IssueView{}, fmt.Errorf("created issue state is %s, want OPEN", view.State)
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
	for _, l := range input.AddLabels {
		if strings.TrimSpace(l) != "" {
			editArgs = append(editArgs, "--add-label", strings.TrimSpace(l))
			hasEdits = true
		}
	}
	for _, l := range input.RemoveLabels {
		if strings.TrimSpace(l) != "" {
			editArgs = append(editArgs, "--remove-label", strings.TrimSpace(l))
			hasEdits = true
		}
	}
	for _, a := range input.AddAssignees {
		if strings.TrimSpace(a) != "" {
			editArgs = append(editArgs, "--add-assignee", strings.TrimSpace(a))
			hasEdits = true
		}
	}
	for _, a := range input.RemoveAssignees {
		if strings.TrimSpace(a) != "" {
			editArgs = append(editArgs, "--remove-assignee", strings.TrimSpace(a))
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

	if input.Title != nil && after.Title != *input.Title {
		return IssueView{}, fmt.Errorf("title readback disagrees: got %q, want %q", after.Title, *input.Title)
	}
	if input.State != nil && !strings.EqualFold(after.State, *input.State) {
		return IssueView{}, fmt.Errorf("state readback disagrees: got %q, want %q", after.State, *input.State)
	}

	return after, nil
}
