package githubcli

import (
	"context"
	"strings"
	"testing"
)

func TestCreateIssue(t *testing.T) {
	fake := &fakeRunner{t: t, responses: []fakeResponse{
		{
			args:   []string{"issue", "create", "--repo", "owner/repo", "--title", "Test Title", "--body", "Test Body", "--label", "bug", "--assignee", "monalisa"},
			output: []byte("https://github.com/owner/repo/issues/42\n"),
		},
		{
			args:   []string{"issue", "view", "42", "--repo", "owner/repo", "--json", "number,title,body,state,labels,assignees,milestone,url"},
			output: []byte(`{"number":42,"title":"Test Title","body":"Test Body","state":"OPEN","url":"https://github.com/owner/repo/issues/42","labels":[{"name":"bug"}],"assignees":[{"login":"monalisa"}]}`),
		},
	}}

	view, err := CreateIssue(context.Background(), fake, CreateIssueInput{
		Repo:      "owner/repo",
		Title:     "Test Title",
		Body:      "Test Body",
		Labels:    []string{"bug"},
		Assignees: []string{"monalisa"},
	})
	if err != nil {
		t.Fatalf("CreateIssue error = %v", err)
	}
	if view.Number != 42 || view.Title != "Test Title" || view.State != "OPEN" {
		t.Fatalf("view = %+v", view)
	}
	if len(view.Labels) != 1 || view.Labels[0].Name != "bug" {
		t.Fatalf("labels = %+v", view.Labels)
	}
}

func TestEditIssue(t *testing.T) {
	newTitle := "Updated Title"
	state := "closed"
	fake := &fakeRunner{t: t, responses: []fakeResponse{
		{
			// Initial view before edit
			args:   []string{"issue", "view", "42", "--repo", "owner/repo", "--json", "number,title,body,state,labels,assignees,milestone,url"},
			output: []byte(`{"number":42,"title":"Old Title","body":"Old Body","state":"OPEN","url":"https://github.com/owner/repo/issues/42"}`),
		},
		{
			// gh issue edit
			args:   []string{"issue", "edit", "42", "--repo", "owner/repo", "--title", "Updated Title", "--add-label", "enhancement"},
			output: []byte("https://github.com/owner/repo/issues/42\n"),
		},
		{
			// gh issue close
			args:   []string{"issue", "close", "42", "--repo", "owner/repo", "--reason", "completed"},
			output: []byte("Closed issue #42\n"),
		},
		{
			// View after edit for readback
			args:   []string{"issue", "view", "42", "--repo", "owner/repo", "--json", "number,title,body,state,labels,assignees,milestone,url"},
			output: []byte(`{"number":42,"title":"Updated Title","body":"Old Body","state":"CLOSED","url":"https://github.com/owner/repo/issues/42","labels":[{"name":"enhancement"}]}`),
		},
	}}

	view, err := EditIssue(context.Background(), fake, EditIssueInput{
		Repo:        "owner/repo",
		Number:      42,
		Title:       &newTitle,
		AddLabels:   []string{"enhancement"},
		State:       &state,
		CloseReason: "completed",
	})
	if err != nil {
		t.Fatalf("EditIssue error = %v", err)
	}
	if view.Title != "Updated Title" || view.State != "CLOSED" {
		t.Fatalf("view = %+v", view)
	}
}

func TestCreateIssueRejectsEmptyTitle(t *testing.T) {
	fake := &fakeRunner{t: t}
	_, err := CreateIssue(context.Background(), fake, CreateIssueInput{
		Repo:  "owner/repo",
		Title: "   ",
	})
	if err == nil || !strings.Contains(err.Error(), "non-empty title") {
		t.Fatalf("error = %v, want non-empty title", err)
	}
}
