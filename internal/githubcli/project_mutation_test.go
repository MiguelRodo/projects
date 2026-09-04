package githubcli

import (
	"context"
	"strings"
	"testing"

	"github.com/MiguelRodo/projects/internal/contract"
)

func projectSchemaJSON() []byte {
	return []byte(`{
  "data": {
    "user": {
      "projectV2": {
        "id": "PVT_123",
        "number": 40,
        "title": "Planning",
        "fields": {
          "nodes": [
            {
              "__typename": "ProjectV2SingleSelectField",
              "id": "FIELD_STATUS",
              "name": "Status",
              "dataType": "SINGLE_SELECT",
              "options": [
                {"id": "OPT_TODO", "name": "Todo"},
                {"id": "OPT_IN_PROGRESS", "name": "In progress"},
                {"id": "OPT_DONE", "name": "Done"}
              ]
            },
            {
              "__typename": "ProjectV2SingleSelectField",
              "id": "FIELD_PRIORITY",
              "name": "Priority",
              "dataType": "SINGLE_SELECT",
              "options": [
                {"id": "OPT_P0", "name": "P0"},
                {"id": "OPT_P1", "name": "P1"},
                {"id": "OPT_P2", "name": "P2"},
                {"id": "OPT_P3", "name": "P3"}
              ]
            },
            {
              "__typename": "ProjectV2SingleSelectField",
              "id": "FIELD_CLASS",
              "name": "Class",
              "dataType": "SINGLE_SELECT",
              "options": [
                {"id": "OPT_TASK", "name": "Task"},
                {"id": "OPT_BUG", "name": "Bug"}
              ]
            },
            {
              "__typename": "ProjectV2Field",
              "id": "FIELD_TARGET_DATE",
              "name": "Target date",
              "dataType": "DATE"
            }
          ]
        }
      }
    }
  }
}`)
}

func issueProjectItemJSON(status, priority, class string) []byte {
	return []byte(`{
  "data": {
    "repository": {
      "issue": {
        "id": "ISSUE_ID",
        "number": 42,
        "url": "https://github.com/owner/repo/issues/42",
        "projectItems": {
          "nodes": [
            {
              "id": "PVTI_ITEM_42",
              "project": {
                "id": "PVT_123",
                "number": 40,
                "title": "Planning"
              },
              "fieldValues": {
                "nodes": [
                  {
                    "__typename": "ProjectV2ItemFieldSingleSelectValue",
                    "field": {"id": "FIELD_STATUS", "name": "Status"},
                    "name": "` + status + `"
                  },
                  {
                    "__typename": "ProjectV2ItemFieldSingleSelectValue",
                    "field": {"id": "FIELD_PRIORITY", "name": "Priority"},
                    "name": "` + priority + `"
                  },
                  {
                    "__typename": "ProjectV2ItemFieldSingleSelectValue",
                    "field": {"id": "FIELD_CLASS", "name": "Class"},
                    "name": "` + class + `"
                  }
                ]
              }
            }
          ]
        }
      }
    }
  }
}`)
}

func TestQueryProjectSchema(t *testing.T) {
	p := contract.Project{
		Owner:     "octo-user",
		OwnerType: "user",
		Number:    40,
	}
	fake := &fakeRunner{t: t, responses: []fakeResponse{
		{
			args: []string{
				"api", "graphql",
				"-f", "query=" + ProjectSchemaQuery("user"),
				"-F", "login=octo-user",
				"-F", "number=40",
			},
			output: projectSchemaJSON(),
		},
	}}

	schema, err := QueryProjectSchema(context.Background(), fake, p)
	if err != nil {
		t.Fatalf("QueryProjectSchema error = %v", err)
	}
	if schema.ID != "PVT_123" || schema.Number != 40 {
		t.Fatalf("schema = %+v", schema)
	}
	statusField, ok := schema.FindField("Status")
	if !ok {
		t.Fatal("Status field not found")
	}
	optID, ok := statusField.FindOptionID("In progress")
	if !ok || optID != "OPT_IN_PROGRESS" {
		t.Fatalf("option id = %q, want OPT_IN_PROGRESS", optID)
	}
}

func TestMutateProjectItem(t *testing.T) {
	p := contract.Project{
		Owner:     "octo-user",
		OwnerType: "user",
		Number:    40,
		Title:     "Planning",
		Priority: map[string]string{
			"P0": "P0",
			"P1": "P1",
			"P2": "P2",
			"P3": "P3",
		},
		ClassValues: []string{"Task", "Bug"},
	}

	fake := &fakeRunner{t: t, responses: []fakeResponse{
		// 1. QueryProjectSchema
		{
			args: []string{
				"api", "graphql",
				"-f", "query=" + ProjectSchemaQuery("user"),
				"-F", "login=octo-user",
				"-F", "number=40",
			},
			output: projectSchemaJSON(),
		},
		// 2. QueryIssueProjectItem (before)
		{
			args: []string{
				"api", "graphql",
				"-f", "query=" + IssueProjectItemQuery,
				"-F", "owner=owner",
				"-F", "repo=repo",
				"-F", "number=42",
			},
			output: issueProjectItemJSON("Todo", "P2", "Task"),
		},
		// 3. EditProjectItemField for Priority
		{
			args: []string{
				"project", "item-edit",
				"--id", "PVTI_ITEM_42",
				"--project-id", "PVT_123",
				"--field-id", "FIELD_PRIORITY",
				"--single-select-option-id", "OPT_P1",
			},
			output: []byte("{}"),
		},
		// 4. EditProjectItemField for Status
		{
			args: []string{
				"project", "item-edit",
				"--id", "PVTI_ITEM_42",
				"--project-id", "PVT_123",
				"--field-id", "FIELD_STATUS",
				"--single-select-option-id", "OPT_IN_PROGRESS",
			},
			output: []byte("{}"),
		},
		// 5. QueryIssueProjectItem (readback after)
		{
			args: []string{
				"api", "graphql",
				"-f", "query=" + IssueProjectItemQuery,
				"-F", "owner=owner",
				"-F", "repo=repo",
				"-F", "number=42",
			},
			output: issueProjectItemJSON("In progress", "P1", "Task"),
		},
	}}

	result, err := MutateProjectItem(context.Background(), fake, MutateProjectItemInput{
		Project:     p,
		Repo:        "owner/repo",
		IssueNumber: 42,
		Priority:    "P1",
		Status:      "In progress",
	})
	if err != nil {
		t.Fatalf("MutateProjectItem error = %v", err)
	}
	if result.ItemID != "PVTI_ITEM_42" {
		t.Fatalf("result.ItemID = %q, want PVTI_ITEM_42", result.ItemID)
	}
	if result.Fields["Priority"] != "P1" || result.Fields["Status"] != "In progress" {
		t.Fatalf("result.Fields = %+v", result.Fields)
	}
}

func TestMutateProjectItemRejectsPendingPriority(t *testing.T) {
	p := contract.Project{
		Owner:        "octo-user",
		OwnerType:    "user",
		Number:       40,
		Pending:      true,
		ContractPath: ".projects/project.md",
	}
	fake := &fakeRunner{t: t}

	_, err := MutateProjectItem(context.Background(), fake, MutateProjectItemInput{
		Project:     p,
		Repo:        "owner/repo",
		IssueNumber: 42,
		Priority:    "P1",
	})
	if err == nil || !strings.Contains(err.Error(), "pending") {
		t.Fatalf("error = %v, want pending priority error", err)
	}
}
