package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestIssueCreatePlanDefault(t *testing.T) {
	var stdout, stderr bytes.Buffer
	fake := &runner{t: t} // Should not be called in plan mode
	exitCode := Run(
		context.Background(),
		[]string{"issue", "create", "--root", fixture(t, "single"), "--title", "Sample Plan Issue", "--json"},
		&stdout,
		&stderr,
		fake,
	)
	if exitCode != 0 {
		t.Fatalf("exit code = %d, stderr = %s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"action": "create_issue"`) ||
		!strings.Contains(stdout.String(), `"apply": false`) ||
		!strings.Contains(stdout.String(), `"title": "Sample Plan Issue"`) {
		t.Fatalf("stdout = %s", stdout.String())
	}
	if !strings.Contains(stderr.String(), "[1/1] Planning issue creation") {
		t.Fatalf("stderr = %s", stderr.String())
	}
}

func TestIssueCreateApply(t *testing.T) {
	var stdout, stderr bytes.Buffer
	fake := &runner{t: t, responses: []response{
		{
			args:   []string{"issue", "create", "--repo", "octo-org/example", "--title", "Created Issue", "--body", "Body text", "--label", "bug"},
			output: "https://github.com/octo-org/example/issues/55\n",
		},
		{
			args:   []string{"issue", "view", "55", "--repo", "octo-org/example", "--json", "number,title,body,state,labels,assignees,milestone,url"},
			output: `{"number":55,"title":"Created Issue","body":"Body text","state":"OPEN","url":"https://github.com/octo-org/example/issues/55","labels":[{"name":"bug"}]}`,
		},
	}}

	exitCode := Run(
		context.Background(),
		[]string{"issue", "create", "--root", fixture(t, "single"), "--title", "Created Issue", "--body", "Body text", "--label", "bug", "--apply", "--json"},
		&stdout,
		&stderr,
		fake,
	)
	if exitCode != 0 {
		t.Fatalf("exit code = %d, stderr = %s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"applied": true`) ||
		!strings.Contains(stdout.String(), `"number": 55`) {
		t.Fatalf("stdout = %s", stdout.String())
	}
	if !strings.Contains(stderr.String(), "[2/4] Creating issue") ||
		!strings.Contains(stderr.String(), "[3/4] Verified created issue") {
		t.Fatalf("stderr = %s", stderr.String())
	}
}

func TestIssueEditPlan(t *testing.T) {
	var stdout, stderr bytes.Buffer
	fake := &runner{t: t, responses: []response{
		{
			args:   []string{"issue", "view", "55", "--repo", "octo-org/example", "--json", "number,title,body,state,labels,assignees,milestone,url"},
			output: `{"number":55,"title":"Original Title","body":"Original Body","state":"OPEN","url":"https://github.com/octo-org/example/issues/55"}`,
		},
	}}

	exitCode := Run(
		context.Background(),
		[]string{"issue", "edit", "--root", fixture(t, "single"), "--issue", "55", "--title", "New Title", "--json"},
		&stdout,
		&stderr,
		fake,
	)
	if exitCode != 0 {
		t.Fatalf("exit code = %d, stderr = %s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"action": "edit_issue"`) ||
		!strings.Contains(stdout.String(), `"apply": false`) ||
		!strings.Contains(stdout.String(), `"New Title"`) {
		t.Fatalf("stdout = %s", stdout.String())
	}
}

func TestIssueEditApply(t *testing.T) {
	var stdout, stderr bytes.Buffer
	fake := &runner{t: t, responses: []response{
		{
			args:   []string{"issue", "view", "55", "--repo", "octo-org/example", "--json", "number,title,body,state,labels,assignees,milestone,url"},
			output: `{"number":55,"title":"Original Title","body":"Original Body","state":"OPEN","url":"https://github.com/octo-org/example/issues/55"}`,
		},
		{
			args:   []string{"issue", "edit", "55", "--repo", "octo-org/example", "--title", "New Title"},
			output: "https://github.com/octo-org/example/issues/55\n",
		},
		{
			args:   []string{"issue", "view", "55", "--repo", "octo-org/example", "--json", "number,title,body,state,labels,assignees,milestone,url"},
			output: `{"number":55,"title":"New Title","body":"Original Body","state":"OPEN","url":"https://github.com/octo-org/example/issues/55"}`,
		},
	}}

	exitCode := Run(
		context.Background(),
		[]string{"issue", "edit", "--root", fixture(t, "single"), "--issue", "55", "--title", "New Title", "--apply", "--json"},
		&stdout,
		&stderr,
		fake,
	)
	if exitCode != 0 {
		t.Fatalf("exit code = %d, stderr = %s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"applied": true`) ||
		!strings.Contains(stdout.String(), `"New Title"`) {
		t.Fatalf("stdout = %s", stdout.String())
	}
}

func TestProjectItemAddPlan(t *testing.T) {
	var stdout, stderr bytes.Buffer
	fake := &runner{t: t}
	exitCode := Run(
		context.Background(),
		[]string{"project", "item-add", "--root", fixture(t, "single"), "--issue", "55", "--json"},
		&stdout,
		&stderr,
		fake,
	)
	if exitCode != 0 {
		t.Fatalf("exit code = %d, stderr = %s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"action": "project_item_add"`) ||
		!strings.Contains(stdout.String(), `"apply": false`) ||
		!strings.Contains(stdout.String(), "issues/55") {
		t.Fatalf("stdout = %s", stdout.String())
	}
}

func TestProjectItemAddApply(t *testing.T) {
	var stdout, stderr bytes.Buffer
	fake := &runner{t: t, responses: []response{
		{
			args:   []string{"project", "item-add", "12", "--owner", "octo-org", "--url", "https://github.com/octo-org/example/issues/55", "--format", "json"},
			output: `{"id": "PVTI_ITEM_55"}`,
		},
	}}

	exitCode := Run(
		context.Background(),
		[]string{"project", "item-add", "--root", fixture(t, "single"), "--issue", "55", "--apply", "--json"},
		&stdout,
		&stderr,
		fake,
	)
	if exitCode != 0 {
		t.Fatalf("exit code = %d, stderr = %s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"applied": true`) ||
		!strings.Contains(stdout.String(), "PVTI_ITEM_55") {
		t.Fatalf("stdout = %s", stdout.String())
	}
}

func TestProjectItemEditPlan(t *testing.T) {
	var stdout, stderr bytes.Buffer
	fake := &runner{t: t}
	exitCode := Run(
		context.Background(),
		[]string{"project", "item-edit", "--root", fixture(t, "single"), "--issue", "55", "--priority", "P1", "--json"},
		&stdout,
		&stderr,
		fake,
	)
	if exitCode != 0 {
		t.Fatalf("exit code = %d, stderr = %s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"action": "project_item_edit"`) ||
		!strings.Contains(stdout.String(), `"apply": false`) ||
		!strings.Contains(stdout.String(), `"priority": "High"`) {
		t.Fatalf("stdout = %s", stdout.String())
	}
}
