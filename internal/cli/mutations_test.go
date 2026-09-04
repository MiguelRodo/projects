package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/MiguelRodo/projects/internal/githubcli"
)

func TestIssueCreatePlanDefault(t *testing.T) {
	var stdout, stderr bytes.Buffer
	fake := &runner{t: t, responses: []response{
		{
			args: []string{
				"api", "--paginate", "--slurp",
				"-H", "Accept: application/vnd.github+json",
				"-H", "X-GitHub-Api-Version: 2026-03-10",
				"repos/octo-org/example/issues?state=all&per_page=100",
			},
			output: `[[]]`,
		},
	}}
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
	if !strings.Contains(stderr.String(), "[1/2] Inspecting") || !strings.Contains(stderr.String(), "[2/2] Planning issue creation") {
		t.Fatalf("stderr = %s", stderr.String())
	}
}

func TestIssueCreateApply(t *testing.T) {
	var stdout, stderr bytes.Buffer
	fake := &runner{t: t, responses: []response{
		{
			args: []string{
				"api", "--paginate", "--slurp",
				"-H", "Accept: application/vnd.github+json",
				"-H", "X-GitHub-Api-Version: 2026-03-10",
				"repos/octo-org/example/issues?state=all&per_page=100",
			},
			output: `[[]]`,
		},
		{
			args:   []string{"issue", "create", "--repo", "octo-org/example", "--title", "Created Issue", "--body", "Body text", "--label", "bug"},
			output: "https://github.com/octo-org/example/issues/55\n",
		},
		{
			args:   []string{"issue", "view", "55", "--repo", "octo-org/example", "--json", "number,title,body,state,stateReason,labels,assignees,milestone,issueType,projectItems,url"},
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
			args:   []string{"issue", "view", "55", "--repo", "octo-org/example", "--json", "number,title,body,state,stateReason,labels,assignees,milestone,issueType,projectItems,url"},
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
			args:   []string{"issue", "view", "55", "--repo", "octo-org/example", "--json", "number,title,body,state,stateReason,labels,assignees,milestone,issueType,projectItems,url"},
			output: `{"number":55,"title":"Original Title","body":"Original Body","state":"OPEN","url":"https://github.com/octo-org/example/issues/55"}`,
		},
		{
			args:   []string{"issue", "edit", "55", "--repo", "octo-org/example", "--title", "New Title"},
			output: "https://github.com/octo-org/example/issues/55\n",
		},
		{
			args:   []string{"issue", "view", "55", "--repo", "octo-org/example", "--json", "number,title,body,state,stateReason,labels,assignees,milestone,issueType,projectItems,url"},
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
	fake := &runner{t: t, responses: []response{
		{
			args:   []string{"project", "view", "12", "--owner", "octo-org", "--format", "json"},
			output: `{"number":12,"owner":{"login":"octo-org"},"title":"Example planning","url":"https://example.invalid/project"}`,
		},
		{
			args:   []string{"project", "item-list", "12", "--owner", "octo-org", "--limit", "10000", "--format", "json"},
			output: `{"items":[],"totalCount":0}`,
		},
	}}
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
			args:   []string{"project", "view", "12", "--owner", "octo-org", "--format", "json"},
			output: `{"number":12,"owner":{"login":"octo-org"},"title":"Example planning","url":"https://example.invalid/project"}`,
		},
		{
			args:   []string{"project", "item-list", "12", "--owner", "octo-org", "--limit", "10000", "--format", "json"},
			output: `{"items":[],"totalCount":0}`,
		},
		{
			args:   []string{"project", "item-add", "12", "--owner", "octo-org", "--url", "https://github.com/octo-org/example/issues/55", "--format", "json"},
			output: `{"id": "PVTI_ITEM_55"}`,
		},
		{
			args:   []string{"project", "view", "12", "--owner", "octo-org", "--format", "json"},
			output: `{"number":12,"owner":{"login":"octo-org"},"title":"Example planning","url":"https://example.invalid/project"}`,
		},
		{
			args:   []string{"project", "item-list", "12", "--owner", "octo-org", "--limit", "10000", "--format", "json"},
			output: `{"items":[{"id":"PVTI_ITEM_55","content":{"number":55,"repository":"octo-org/example","title":"Issue","type":"Issue","url":"https://github.com/octo-org/example/issues/55"}}],"totalCount":1}`,
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
	fake := &runner{t: t, responses: []response{
		{
			args: []string{
				"api", "graphql",
				"-f", "query=" + githubcli.ProjectSchemaQuery("organization"),
				"-f", "login=octo-org",
				"-F", "number=12",
			},
			output: `{"data":{"organization":{"projectV2":{"id":"PVT_12","number":12,"title":"Example planning","fields":{"nodes":[],"pageInfo":{"hasNextPage":false}}}}}}`,
		},
		{
			args: []string{
				"api", "--paginate", "--slurp",
				"-H", "Accept: application/vnd.github+json",
				"-H", "X-GitHub-Api-Version: 2026-03-10",
				"orgs/octo-org/issue-fields?per_page=100",
			},
			output: `[[{"id":1,"name":"Priority","data_type":"single_select","options":[{"id":2,"name":"High"}]}]]`,
		},
		{
			args:   []string{"project", "view", "12", "--owner", "octo-org", "--format", "json"},
			output: `{"number":12,"owner":{"login":"octo-org"},"title":"Example planning","url":"https://example.invalid/project"}`,
		},
		{
			args:   []string{"project", "item-list", "12", "--owner", "octo-org", "--limit", "10000", "--format", "json"},
			output: `{"items":[{"id":"PVTI_ITEM_55","priority":"Medium","content":{"number":55,"repository":"octo-org/example","title":"Issue","type":"Issue","url":"https://github.com/octo-org/example/issues/55"}}],"totalCount":1}`,
		},
	}}
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

func TestIssueCreateRejectsExactTitleDuplicate(t *testing.T) {
	var stdout, stderr bytes.Buffer
	fake := &runner{t: t, responses: []response{{
		args: []string{
			"api", "--paginate", "--slurp",
			"-H", "Accept: application/vnd.github+json",
			"-H", "X-GitHub-Api-Version: 2026-03-10",
			"repos/octo-org/example/issues?state=all&per_page=100",
		},
		output: `[[{"number":55,"title":"Existing","state":"open","html_url":"https://github.com/octo-org/example/issues/55"}]]`,
	}}}
	exitCode := Run(context.Background(), []string{
		"issue", "create", "--root", fixture(t, "single"), "--title", "Existing", "--apply",
	}, &stdout, &stderr, fake)
	if exitCode != 1 || !strings.Contains(stderr.String(), "exact-title issue already exists") {
		t.Fatalf("exit code = %d, stderr = %s", exitCode, stderr.String())
	}
}

func TestMutationCommandsRejectRepositoryDisagreementBeforeGitHub(t *testing.T) {
	var stdout, stderr bytes.Buffer
	exitCode := Run(context.Background(), []string{
		"issue", "edit", "--root", fixture(t, "single"), "--repo", "other/repo", "--issue", "55", "--title", "Title",
	}, &stdout, &stderr, &runner{t: t})
	if exitCode != 2 || !strings.Contains(stderr.String(), "disagrees with contract repository") {
		t.Fatalf("exit code = %d, stderr = %s", exitCode, stderr.String())
	}
}

func TestProjectMutationTargetsAreMutuallyExclusive(t *testing.T) {
	for _, command := range []string{"item-add", "item-edit"} {
		var stdout, stderr bytes.Buffer
		args := []string{"project", command, "--root", fixture(t, "single"), "--issue", "55", "--url", "https://github.com/octo-org/example/issues/55"}
		if command == "item-edit" {
			args = append(args, "--status", "Todo")
		}
		exitCode := Run(context.Background(), args, &stdout, &stderr, &runner{t: t})
		if exitCode != 2 || !strings.Contains(stderr.String(), "mutually exclusive") {
			t.Fatalf("%s exit code = %d, stderr = %s", command, exitCode, stderr.String())
		}
	}
}

func TestDispatcherIssueCreatePlanIncludesRoutingLabel(t *testing.T) {
	var stdout, stderr bytes.Buffer
	fake := &runner{t: t, responses: []response{{
		args: []string{
			"api", "--paginate", "--slurp",
			"-H", "Accept: application/vnd.github+json",
			"-H", "X-GitHub-Api-Version: 2026-03-10",
			"repos/octo-user/issues/issues?state=all&per_page=100",
		},
		output: `[[]]`,
	}}}
	exitCode := Run(context.Background(), []string{
		"issue", "create", "--root", fixture(t, "dispatcher"), "--project-key", "alpha", "--title", "Routed", "--json",
	}, &stdout, &stderr, fake)
	if exitCode != 0 {
		t.Fatalf("exit code = %d, stderr = %s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"project:alpha"`) {
		t.Fatalf("stdout = %s", stdout.String())
	}
}
