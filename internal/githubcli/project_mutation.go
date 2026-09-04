package githubcli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/MiguelRodo/projects/internal/contract"
)

// ProjectSchema holds the discovered GraphQL schema for a ProjectV2.
type ProjectSchema struct {
	ID     string
	Number int
	Title  string
	Fields map[string]ProjectField // lowercased field name -> field
}

// ProjectField holds one field's definition and its single-select options.
type ProjectField struct {
	ID       string
	Name     string
	DataType string
	Options  map[string]string // lowercased option name -> option ID
}

// FindField looks up a field by name (case-insensitive).
func (s ProjectSchema) FindField(name string) (ProjectField, bool) {
	field, ok := s.Fields[strings.ToLower(strings.TrimSpace(name))]
	return field, ok
}

// FindOptionID looks up a single-select option ID by name (case-insensitive).
func (f ProjectField) FindOptionID(name string) (string, bool) {
	if f.Options == nil {
		return "", false
	}
	id, ok := f.Options[strings.ToLower(strings.TrimSpace(name))]
	return id, ok
}

// ProjectItemState contains the inspected fields for an item in a Project.
type ProjectItemState struct {
	ItemID      string            `json:"itemId"`
	IssueNumber int               `json:"issueNumber,omitempty"`
	URL         string            `json:"url"`
	Fields      map[string]string `json:"fields"`
}

type graphQLFieldNode struct {
	Typename string `json:"__typename"`
	ID       string `json:"id"`
	Name     string `json:"name"`
	DataType string `json:"dataType"`
	Options  []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"options"`
}

type graphQLProjectData struct {
	ID     string `json:"id"`
	Number int    `json:"number"`
	Title  string `json:"title"`
	Fields struct {
		Nodes []graphQLFieldNode `json:"nodes"`
	} `json:"fields"`
}

type graphQLSchemaResponse struct {
	Data struct {
		User struct {
			ProjectV2 *graphQLProjectData `json:"projectV2"`
		} `json:"user"`
		Organization struct {
			ProjectV2 *graphQLProjectData `json:"projectV2"`
		} `json:"organization"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// ProjectSchemaQuery returns the GraphQL query for fetching a ProjectV2's fields.
func ProjectSchemaQuery(ownerRoot string) string {
	return fmt.Sprintf(`query($login: String!, $number: Int!) {
  %s(login: $login) {
    projectV2(number: $number) {
      id
      number
      title
      fields(first: 100) {
        nodes {
          __typename
          ... on ProjectV2FieldCommon { id name dataType }
          ... on ProjectV2SingleSelectField { options { id name } }
        }
      }
    }
  }
}`, ownerRoot)
}

// QueryProjectSchema fetches the Project node ID, fields, and options via GraphQL.
func QueryProjectSchema(ctx context.Context, runner Runner, project contract.Project) (ProjectSchema, error) {
	ownerRoot := project.OwnerType
	if ownerRoot == "" {
		discovered, err := discoverOwnerType(ctx, runner, project.Owner)
		if err != nil {
			return ProjectSchema{}, fmt.Errorf("discover owner type for %s: %w", project.Owner, err)
		}
		ownerRoot = discovered
	}

	query := ProjectSchemaQuery(ownerRoot)

	output, err := runner.Run(
		ctx,
		"api", "graphql",
		"-f", "query="+query,
		"-F", "login="+project.Owner,
		"-F", "number="+strconv.Itoa(project.Number),
	)
	if err != nil {
		return ProjectSchema{}, fmt.Errorf("query Project schema for %s/%d: %w", project.Owner, project.Number, err)
	}

	var resp graphQLSchemaResponse
	if err := json.Unmarshal(output, &resp); err != nil {
		return ProjectSchema{}, fmt.Errorf("decode Project schema response: %w", err)
	}
	if len(resp.Errors) > 0 {
		return ProjectSchema{}, fmt.Errorf("GraphQL error querying Project schema: %s", resp.Errors[0].Message)
	}

	var data *graphQLProjectData
	if ownerRoot == "organization" {
		data = resp.Data.Organization.ProjectV2
	} else {
		data = resp.Data.User.ProjectV2
	}
	if data == nil {
		return ProjectSchema{}, fmt.Errorf("Project %s/%d not found at %s root", project.Owner, project.Number, ownerRoot)
	}

	schema := ProjectSchema{
		ID:     data.ID,
		Number: data.Number,
		Title:  data.Title,
		Fields: make(map[string]ProjectField),
	}
	for _, node := range data.Fields.Nodes {
		f := ProjectField{
			ID:       node.ID,
			Name:     node.Name,
			DataType: node.DataType,
			Options:  make(map[string]string),
		}
		for _, opt := range node.Options {
			f.Options[strings.ToLower(strings.TrimSpace(opt.Name))] = opt.ID
		}
		schema.Fields[strings.ToLower(strings.TrimSpace(node.Name))] = f
	}

	return schema, nil
}

func discoverOwnerType(ctx context.Context, runner Runner, owner string) (string, error) {
	output, err := runner.Run(ctx, "api", "users/"+owner, "--jq", ".type")
	if err != nil {
		return "", err
	}
	trimmed := strings.TrimSpace(string(output))
	switch trimmed {
	case "User":
		return "user", nil
	case "Organization":
		return "organization", nil
	default:
		return "", fmt.Errorf("unsupported owner type %q for %s", trimmed, owner)
	}
}

// AddProjectItem adds an issue or pull request to a Project via gh project item-add.
func AddProjectItem(ctx context.Context, runner Runner, projectNumber int, projectOwner, url string) (string, error) {
	output, err := runner.Run(
		ctx,
		"project", "item-add", strconv.Itoa(projectNumber),
		"--owner", projectOwner,
		"--url", url,
		"--format", "json",
	)
	if err != nil {
		return "", fmt.Errorf("add item to Project %s/%d: %w", projectOwner, projectNumber, err)
	}
	var res struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(output, &res); err != nil {
		return "", fmt.Errorf("decode item-add response: %w", err)
	}
	if res.ID == "" {
		return "", fmt.Errorf("item-add returned empty item ID for %s", url)
	}
	return res.ID, nil
}

type graphQLIssueProjectResponse struct {
	Data struct {
		Repository struct {
			Issue *struct {
				ID           string `json:"id"`
				Number       int    `json:"number"`
				URL          string `json:"url"`
				ProjectItems struct {
					Nodes []struct {
						ID      string `json:"id"`
						Project struct {
							ID     string `json:"id"`
							Number int    `json:"number"`
							Title  string `json:"title"`
						} `json:"project"`
						FieldValues struct {
							Nodes []struct {
								Typename string `json:"__typename"`
								Field    struct {
									ID   string `json:"id"`
									Name string `json:"name"`
								} `json:"field"`
								Name     string `json:"name,omitempty"`     // for single-select
								OptionID string `json:"optionId,omitempty"` // for single-select
								Date     string `json:"date,omitempty"`     // for date
								Text     string `json:"text,omitempty"`     // for text
							} `json:"nodes"`
						} `json:"fieldValues"`
					} `json:"nodes"`
				} `json:"projectItems"`
			} `json:"issue"`
		} `json:"repository"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// IssueProjectItemQuery is the GraphQL query for reading an issue's project items and field values.
const IssueProjectItemQuery = `query($owner: String!, $repo: String!, $number: Int!) {
  repository(owner: $owner, name: $repo) {
    issue(number: $number) {
      id
      number
      url
      projectItems(first: 20) {
        nodes {
          id
          project { id number title }
          fieldValues(first: 30) {
            nodes {
              __typename
              ... on ProjectV2ItemFieldSingleSelectValue {
                field { ... on ProjectV2SingleSelectField { id name } }
                name
                optionId
              }
              ... on ProjectV2ItemFieldDateValue {
                field { ... on ProjectV2FieldCommon { id name } }
                date
              }
              ... on ProjectV2ItemFieldTextValue {
                field { ... on ProjectV2FieldCommon { id name } }
                text
              }
            }
          }
        }
      }
    }
  }
}`

// QueryIssueProjectItem checks if an issue is in a Project and returns its item state.
func QueryIssueProjectItem(ctx context.Context, runner Runner, repoOwner, repoName string, issueNumber, projectNumber int) (*ProjectItemState, error) {
	query := IssueProjectItemQuery

	output, err := runner.Run(
		ctx,
		"api", "graphql",
		"-f", "query="+query,
		"-F", "owner="+repoOwner,
		"-F", "repo="+repoName,
		"-F", "number="+strconv.Itoa(issueNumber),
	)
	if err != nil {
		return nil, fmt.Errorf("query issue project items for %s/%s#%d: %w", repoOwner, repoName, issueNumber, err)
	}

	var resp graphQLIssueProjectResponse
	if err := json.Unmarshal(output, &resp); err != nil {
		return nil, fmt.Errorf("decode issue project response: %w", err)
	}
	if len(resp.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", resp.Errors[0].Message)
	}
	if resp.Data.Repository.Issue == nil {
		return nil, fmt.Errorf("issue %s/%s#%d not found", repoOwner, repoName, issueNumber)
	}

	issue := resp.Data.Repository.Issue
	for _, itemNode := range issue.ProjectItems.Nodes {
		if itemNode.Project.Number == projectNumber {
			fields := make(map[string]string)
			for _, fv := range itemNode.FieldValues.Nodes {
				fieldName := fv.Field.Name
				if fieldName == "" {
					continue
				}
				if fv.Name != "" {
					fields[fieldName] = fv.Name
				} else if fv.Date != "" {
					fields[fieldName] = fv.Date
				} else if fv.Text != "" {
					fields[fieldName] = fv.Text
				}
			}
			return &ProjectItemState{
				ItemID:      itemNode.ID,
				IssueNumber: issue.Number,
				URL:         issue.URL,
				Fields:      fields,
			}, nil
		}
	}

	return nil, nil // not in project
}

// EditItemFieldInput specifies parameters for editing a Project item field.
type EditItemFieldInput struct {
	ItemID               string
	ProjectNodeID        string
	FieldID              string
	SingleSelectOptionID string
	DateValue            string
	TextValue            string
	Clear                bool
}

// EditProjectItemField calls gh project item-edit to change a single field.
func EditProjectItemField(ctx context.Context, runner Runner, input EditItemFieldInput) error {
	args := []string{
		"project", "item-edit",
		"--id", input.ItemID,
		"--project-id", input.ProjectNodeID,
		"--field-id", input.FieldID,
	}
	if input.Clear {
		args = append(args, "--clear")
	} else if input.SingleSelectOptionID != "" {
		args = append(args, "--single-select-option-id", input.SingleSelectOptionID)
	} else if input.DateValue != "" {
		args = append(args, "--date", input.DateValue)
	} else if input.TextValue != "" {
		args = append(args, "--text", input.TextValue)
	} else {
		return errors.New("edit item field requires a value or --clear")
	}

	_, err := runner.Run(ctx, args...)
	if err != nil {
		return fmt.Errorf("edit Project item field: %w", err)
	}
	return nil
}

// MutateProjectItemInput defines the full high-level mutation request.
type MutateProjectItemInput struct {
	Project     contract.Project
	Repo        string // "owner/name"
	IssueNumber int
	URL         string
	Priority    string // common value: P0, P1, P2, P3
	Class       string // common or declared class
	Status      string // status name
	TargetDate  string // YYYY-MM-DD
	ClearFields []string
}

// MutateProjectItemResult is the verified result after editing Project item fields.
type MutateProjectItemResult struct {
	Project ProjectIdentity   `json:"project"`
	ItemID  string            `json:"itemId"`
	URL     string            `json:"url"`
	Fields  map[string]string `json:"fields"`
}

// MutateProjectItem performs verified Project item addition and field updates.
func MutateProjectItem(ctx context.Context, runner Runner, input MutateProjectItemInput) (MutateProjectItemResult, error) {
	if input.Priority != "" {
		if _, err := input.Project.ResolvePriority(input.Priority); err != nil {
			return MutateProjectItemResult{}, err
		}
	}
	if input.Class != "" {
		if _, err := input.Project.ValidateClass(input.Class); err != nil {
			return MutateProjectItemResult{}, err
		}
	}

	schema, err := QueryProjectSchema(ctx, runner, input.Project)
	if err != nil {
		return MutateProjectItemResult{}, err
	}

	repoParts := strings.Split(input.Repo, "/")
	if len(repoParts) != 2 && input.URL != "" {
		// Try parsing from URL: https://github.com/owner/repo/issues/123
		parts := strings.Split(strings.TrimPrefix(input.URL, "https://github.com/"), "/")
		if len(parts) >= 4 && parts[2] == "issues" {
			repoParts = []string{parts[0], parts[1]}
			if input.IssueNumber == 0 {
				input.IssueNumber, _ = strconv.Atoi(parts[3])
			}
		}
	}

	if input.URL == "" && input.IssueNumber > 0 && len(repoParts) == 2 {
		input.URL = fmt.Sprintf("https://github.com/%s/%s/issues/%d", repoParts[0], repoParts[1], input.IssueNumber)
	}

	if input.URL == "" {
		return MutateProjectItemResult{}, errors.New("cannot determine target issue URL")
	}

	var currentItem *ProjectItemState
	if len(repoParts) == 2 && input.IssueNumber > 0 {
		currentItem, err = QueryIssueProjectItem(ctx, runner, repoParts[0], repoParts[1], input.IssueNumber, input.Project.Number)
		if err != nil {
			return MutateProjectItemResult{}, fmt.Errorf("check current project item state: %w", err)
		}
	}

	itemID := ""
	if currentItem != nil {
		itemID = currentItem.ItemID
	} else {
		// Add to project
		addedID, err := AddProjectItem(ctx, runner, input.Project.Number, input.Project.Owner, input.URL)
		if err != nil {
			return MutateProjectItemResult{}, fmt.Errorf("add item to project: %w", err)
		}
		itemID = addedID
	}

	// Apply field edits
	if input.Priority != "" {
		providerVal, err := input.Project.ResolvePriority(input.Priority)
		if err != nil {
			return MutateProjectItemResult{}, err
		}
		pField, ok := schema.FindField("Priority")
		if !ok {
			return MutateProjectItemResult{}, fmt.Errorf("Priority field not found in Project %s/%d schema", input.Project.Owner, input.Project.Number)
		}
		optID, ok := pField.FindOptionID(providerVal)
		if !ok {
			return MutateProjectItemResult{}, fmt.Errorf("option %q for Priority not found in Project %s/%d", providerVal, input.Project.Owner, input.Project.Number)
		}
		if err := EditProjectItemField(ctx, runner, EditItemFieldInput{
			ItemID:               itemID,
			ProjectNodeID:        schema.ID,
			FieldID:              pField.ID,
			SingleSelectOptionID: optID,
		}); err != nil {
			return MutateProjectItemResult{}, fmt.Errorf("set Priority: %w", err)
		}
	}

	if input.Class != "" {
		validClass, err := input.Project.ValidateClass(input.Class)
		if err != nil {
			return MutateProjectItemResult{}, err
		}
		// Check where Class lives
		loc := input.Project.FieldLocations["Class"]
		if loc.Location == "organization issue type" {
			// Update issue type via gh issue edit
			if len(repoParts) == 2 && input.IssueNumber > 0 {
				_, err := runner.Run(ctx, "issue", "edit", strconv.Itoa(input.IssueNumber), "--repo", input.Repo, "--type", validClass)
				if err != nil {
					return MutateProjectItemResult{}, fmt.Errorf("set issue type %q: %w", validClass, err)
				}
			}
		} else {
			cField, ok := schema.FindField("Class")
			if !ok {
				return MutateProjectItemResult{}, fmt.Errorf("Class field not found in Project %s/%d schema", input.Project.Owner, input.Project.Number)
			}
			optID, ok := cField.FindOptionID(validClass)
			if !ok {
				return MutateProjectItemResult{}, fmt.Errorf("option %q for Class not found in Project %s/%d", validClass, input.Project.Owner, input.Project.Number)
			}
			if err := EditProjectItemField(ctx, runner, EditItemFieldInput{
				ItemID:               itemID,
				ProjectNodeID:        schema.ID,
				FieldID:              cField.ID,
				SingleSelectOptionID: optID,
			}); err != nil {
				return MutateProjectItemResult{}, fmt.Errorf("set Class: %w", err)
			}
		}
	}

	if input.Status != "" {
		statusVal := input.Project.ResolveStatus(input.Status)
		sField, ok := schema.FindField("Status")
		if !ok {
			return MutateProjectItemResult{}, fmt.Errorf("Status field not found in Project %s/%d schema", input.Project.Owner, input.Project.Number)
		}
		optID, ok := sField.FindOptionID(statusVal)
		if !ok {
			return MutateProjectItemResult{}, fmt.Errorf("option %q for Status not found in Project %s/%d", statusVal, input.Project.Owner, input.Project.Number)
		}
		if err := EditProjectItemField(ctx, runner, EditItemFieldInput{
			ItemID:               itemID,
			ProjectNodeID:        schema.ID,
			FieldID:              sField.ID,
			SingleSelectOptionID: optID,
		}); err != nil {
			return MutateProjectItemResult{}, fmt.Errorf("set Status: %w", err)
		}
	}

	if input.TargetDate != "" {
		dateFieldName := "Target date"
		if loc, ok := input.Project.FieldLocations["Due date"]; ok && loc.Field != "" {
			dateFieldName = loc.Field
		}
		dField, ok := schema.FindField(dateFieldName)
		if !ok {
			return MutateProjectItemResult{}, fmt.Errorf("date field %q not found in Project %s/%d", dateFieldName, input.Project.Owner, input.Project.Number)
		}
		if err := EditProjectItemField(ctx, runner, EditItemFieldInput{
			ItemID:        itemID,
			ProjectNodeID: schema.ID,
			FieldID:       dField.ID,
			DateValue:     input.TargetDate,
		}); err != nil {
			return MutateProjectItemResult{}, fmt.Errorf("set %s: %w", dateFieldName, err)
		}
	}

	for _, clearName := range input.ClearFields {
		f, ok := schema.FindField(clearName)
		if !ok {
			return MutateProjectItemResult{}, fmt.Errorf("cannot clear unknown field %q", clearName)
		}
		if err := EditProjectItemField(ctx, runner, EditItemFieldInput{
			ItemID:        itemID,
			ProjectNodeID: schema.ID,
			FieldID:       f.ID,
			Clear:         true,
		}); err != nil {
			return MutateProjectItemResult{}, fmt.Errorf("clear field %q: %w", clearName, err)
		}
	}

	// Readback verification
	afterItem, err := QueryIssueProjectItem(ctx, runner, repoParts[0], repoParts[1], input.IssueNumber, input.Project.Number)
	if err != nil {
		return MutateProjectItemResult{}, fmt.Errorf("read back project item: %w", err)
	}
	if afterItem == nil {
		return MutateProjectItemResult{}, fmt.Errorf("project item readback failed: item not found in Project %s/%d", input.Project.Owner, input.Project.Number)
	}

	return MutateProjectItemResult{
		Project: ProjectIdentity{
			Number: input.Project.Number,
			Owner:  input.Project.Owner,
			Title:  input.Project.Title,
		},
		ItemID: itemID,
		URL:    input.URL,
		Fields: afterItem.Fields,
	}, nil
}
