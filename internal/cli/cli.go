// Package cli implements the command-line interface for projectctl / projects.
package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/MiguelRodo/projects/internal/config"
	"github.com/MiguelRodo/projects/internal/git"
	"github.com/MiguelRodo/projects/internal/runner"
	"github.com/MiguelRodo/projects/pkg/project"
	"github.com/MiguelRodo/projects/pkg/version"
)

// App encapsulates dependencies and CLI execution context.
type App struct {
	Stdout     io.Writer
	Stderr     io.Writer
	GitClient  git.Client
	WorkingDir string
}

// NewApp creates a new CLI App with standard system defaults.
func NewApp(stdout, stderr io.Writer) *App {
	cwd, _ := os.Getwd()
	gitCli, _ := git.NewExecClient()
	return &App{
		Stdout:     stdout,
		Stderr:     stderr,
		GitClient:  gitCli,
		WorkingDir: cwd,
	}
}

// Run parses arguments and executes the requested command.
func Run(args []string, stdout, stderr io.Writer) int {
	app := NewApp(stdout, stderr)
	return app.Execute(context.Background(), args)
}

func (a *App) print(s string) {
	_, _ = fmt.Fprint(a.Stdout, s)
}

func (a *App) printf(format string, args ...any) {
	_, _ = fmt.Fprintf(a.Stdout, format, args...)
}

func (a *App) println(args ...any) {
	_, _ = fmt.Fprintln(a.Stdout, args...)
}

func (a *App) errPrintf(format string, args ...any) {
	_, _ = fmt.Fprintf(a.Stderr, format, args...)
}

func (a *App) errPrintln(args ...any) {
	_, _ = fmt.Fprintln(a.Stderr, args...)
}

// parseFlagsAndArgs separates flags from positional arguments before parsing with flag.FlagSet.
func parseFlagsAndArgs(fs *flag.FlagSet, args []string) ([]string, error) {
	var flags []string
	var pos []string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			pos = append(pos, args[i+1:]...)
			break
		}
		if strings.HasPrefix(arg, "-") {
			cleanName := strings.TrimLeft(arg, "-")
			// Handle --flag=value
			if idx := strings.Index(cleanName, "="); idx != -1 {
				flagName := cleanName[:idx]
				if fs.Lookup(flagName) != nil {
					flags = append(flags, arg)
				} else {
					pos = append(pos, arg)
				}
				continue
			}

			// Handle --flag value or boolean --flag
			if f := fs.Lookup(cleanName); f != nil {
				flags = append(flags, arg)
				if bf, ok := f.Value.(interface{ IsBoolFlag() bool }); ok && bf.IsBoolFlag() {
					continue
				}
				if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
					i++
					flags = append(flags, args[i])
				}
			} else {
				// Flag not registered in this FlagSet; treat as positional/passthrough
				pos = append(pos, arg)
			}
		} else {
			pos = append(pos, arg)
		}
	}

	if err := fs.Parse(flags); err != nil {
		return nil, err
	}
	return pos, nil
}

// Execute runs the CLI given arguments and a context.
func (a *App) Execute(ctx context.Context, args []string) int {
	if len(args) == 0 {
		a.printHelp()
		return 0
	}

	// Global flags
	var rootDir string
	var configFile string
	var showVersion bool
	var showHelp bool

	flags := flag.NewFlagSet("projectctl", flag.ContinueOnError)
	flags.SetOutput(a.Stderr)
	flags.StringVar(&rootDir, "C", "", "Run as if started in <path>")
	flags.StringVar(&rootDir, "workspace", "", "Workspace root directory")
	flags.StringVar(&configFile, "config", "", "Path to projects manifest file")
	flags.BoolVar(&showVersion, "v", false, "Show version")
	flags.BoolVar(&showVersion, "version", false, "Show version")
	flags.BoolVar(&showHelp, "h", false, "Show help")
	flags.BoolVar(&showHelp, "help", false, "Show help")

	parsedPos, err := parseFlagsAndArgs(flags, args)
	if err != nil {
		return 2
	}

	if showHelp {
		a.printHelp()
		return 0
	}

	if showVersion {
		a.println(version.String())
		return 0
	}

	if len(parsedPos) == 0 {
		a.printHelp()
		return 0
	}

	if rootDir != "" {
		a.WorkingDir = rootDir
	}

	subcmd := parsedPos[0]
	cmdArgs := parsedPos[1:]

	switch subcmd {
	case "init":
		return a.cmdInit(cmdArgs, configFile)
	case "list", "ls":
		return a.cmdList(cmdArgs, configFile)
	case "add":
		return a.cmdAdd(cmdArgs, configFile)
	case "remove", "rm":
		return a.cmdRemove(cmdArgs, configFile)
	case "sync", "clone":
		return a.cmdSync(ctx, cmdArgs, configFile)
	case "status", "st":
		return a.cmdStatus(ctx, cmdArgs, configFile)
	case "exec":
		return a.cmdExec(ctx, cmdArgs, configFile)
	case "version":
		return a.cmdVersion(cmdArgs)
	case "help":
		a.printHelp()
		return 0
	default:
		a.errPrintf("Unknown command %q. Run 'projectctl help' for usage.\n", subcmd)
		return 1
	}
}

func (a *App) printHelp() {
	helpText := `projectctl - Multi-repository project and workspace management tool

Usage:
  projectctl [flags] <command> [command-flags] [arguments]

Commands:
  init      Initialize a new workspace manifest
  list      List repositories configured in the workspace
  add       Add a new repository to the workspace manifest
  remove    Remove a repository from the workspace manifest
  sync      Clone missing repositories and pull updates for existing ones
  status    Check git status and branch across all repositories
  exec      Execute a command across all repositories in the workspace
  version   Show version and build information
  help      Show help for projectctl

Flags:
  -C, --workspace <path>   Set workspace root directory (default: current directory)
      --config <path>      Path to workspace manifest file (default: auto-detect)
  -v, --version            Show version information
  -h, --help               Show help information

Examples:
  projectctl init --name my-project
  projectctl add https://github.com/example-org/repo1.git
  projectctl sync --pull
  projectctl status
  projectctl exec -- git fetch --all
`
	a.print(helpText)
}

func (a *App) cmdInit(args []string, configPath string) int {
	var name string
	var format string

	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(a.Stderr)
	fs.StringVar(&name, "name", "", "Workspace name")
	fs.StringVar(&format, "format", "json", "Manifest format: json or list")
	pos, err := parseFlagsAndArgs(fs, args)
	if err != nil {
		return 2
	}

	if name == "" && len(pos) > 0 {
		name = pos[0]
	}
	if name == "" {
		name = filepath.Base(a.WorkingDir)
		if name == "." || name == "/" {
			name = "projects"
		}
	}

	ws := project.NewWorkspace(name, a.WorkingDir)

	targetFile := configPath
	if targetFile == "" {
		if format == "list" {
			targetFile = filepath.Join(a.WorkingDir, config.DefaultReposListFile)
		} else {
			targetFile = filepath.Join(a.WorkingDir, config.DefaultProjectsFile)
		}
	}

	if _, err := os.Stat(targetFile); err == nil {
		a.errPrintf("Error: manifest file %q already exists\n", targetFile)
		return 1
	}

	if err := config.SaveWorkspace(ws, targetFile); err != nil {
		a.errPrintf("Error saving workspace: %v\n", err)
		return 1
	}

	a.printf("Initialized empty workspace %q in %s\n", name, targetFile)
	return 0
}

func (a *App) cmdList(args []string, configPath string) int {
	var jsonOutput bool
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	fs.SetOutput(a.Stderr)
	fs.BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	if _, err := parseFlagsAndArgs(fs, args); err != nil {
		return 2
	}

	ws, cfgFile, err := config.LoadWorkspace(a.WorkingDir, configPath)
	if err != nil {
		a.errPrintf("Error loading workspace: %v\n", err)
		return 1
	}

	if jsonOutput {
		enc := json.NewEncoder(a.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(ws.Repositories); err != nil {
			a.errPrintf("Error encoding JSON: %v\n", err)
			return 1
		}
		return 0
	}

	if len(ws.Repositories) == 0 {
		if cfgFile == "" {
			a.println("No workspace manifest found. Run 'projectctl init' or 'projectctl add <url>' to start.")
		} else {
			a.printf("Workspace %q has no repositories.\n", ws.Name)
		}
		return 0
	}

	w := tabwriter.NewWriter(a.Stdout, 0, 0, 3, ' ', 0)
	_, _ = fmt.Fprintln(w, "NAME\tPATH\tBRANCH\tURL")
	for _, r := range ws.Repositories {
		branch := r.Branch
		if branch == "" {
			branch = "-"
		}
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", r.Name, r.Path, branch, r.URL)
	}
	_ = w.Flush()
	return 0
}

func (a *App) cmdAdd(args []string, configPath string) int {
	var name, path, branch, desc string
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	fs.SetOutput(a.Stderr)
	fs.StringVar(&name, "name", "", "Repository name (defaults to repository basename)")
	fs.StringVar(&path, "path", "", "Local directory path relative to workspace root")
	fs.StringVar(&branch, "branch", "", "Default branch to clone")
	fs.StringVar(&desc, "description", "", "Repository description")

	pos, err := parseFlagsAndArgs(fs, args)
	if err != nil {
		return 2
	}

	if len(pos) == 0 {
		a.errPrintln("Error: repository URL is required. Usage: projectctl add <url> [flags]")
		return 2
	}

	rawURL := pos[0]

	ws, cfgFile, err := config.LoadWorkspace(a.WorkingDir, configPath)
	if err != nil {
		a.errPrintf("Error loading workspace: %v\n", err)
		return 1
	}

	if path == "" && len(pos) >= 2 {
		path = pos[1]
	}
	if branch == "" && len(pos) >= 3 {
		branch = pos[2]
	}

	repo := project.Repository{
		Name:        name,
		URL:         rawURL,
		Path:        path,
		Branch:      branch,
		Description: desc,
	}

	if repo.Name == "" {
		parsedList, err := project.ParseReposList(strings.NewReader(rawURL + " " + path + " " + branch))
		if err == nil && len(parsedList) > 0 {
			repo.Name = parsedList[0].Name
			if repo.Path == "" {
				repo.Path = parsedList[0].Path
			}
		}
	}

	if err := ws.AddRepository(repo); err != nil {
		a.errPrintf("Error adding repository: %v\n", err)
		return 1
	}

	if cfgFile == "" {
		cfgFile = filepath.Join(a.WorkingDir, config.DefaultProjectsFile)
	}

	if err := config.SaveWorkspace(ws, cfgFile); err != nil {
		a.errPrintf("Error saving workspace: %v\n", err)
		return 1
	}

	a.printf("Added repository %q (%s) to %s\n", repo.Name, repo.URL, cfgFile)
	return 0
}

func (a *App) cmdRemove(args []string, configPath string) int {
	if len(args) == 0 {
		a.errPrintln("Error: repository name is required. Usage: projectctl remove <name>")
		return 2
	}
	name := args[0]

	ws, cfgFile, err := config.LoadWorkspace(a.WorkingDir, configPath)
	if err != nil {
		a.errPrintf("Error loading workspace: %v\n", err)
		return 1
	}

	if err := ws.RemoveRepository(name); err != nil {
		a.errPrintf("Error removing repository: %v\n", err)
		return 1
	}

	if cfgFile == "" {
		cfgFile = filepath.Join(a.WorkingDir, config.DefaultProjectsFile)
	}

	if err := config.SaveWorkspace(ws, cfgFile); err != nil {
		a.errPrintf("Error saving workspace: %v\n", err)
		return 1
	}

	a.printf("Removed repository %q from %s\n", name, cfgFile)
	return 0
}

func (a *App) cmdSync(ctx context.Context, args []string, configPath string) int {
	var pull bool
	var concurrency int

	fs := flag.NewFlagSet("sync", flag.ContinueOnError)
	fs.SetOutput(a.Stderr)
	fs.BoolVar(&pull, "pull", true, "Pull updates for already cloned repositories")
	fs.IntVar(&concurrency, "concurrency", 4, "Number of concurrent operations")

	if _, err := parseFlagsAndArgs(fs, args); err != nil {
		return 2
	}

	if a.GitClient == nil {
		a.errPrintln("Error: git client is not available")
		return 1
	}

	ws, _, err := config.LoadWorkspace(a.WorkingDir, configPath)
	if err != nil {
		a.errPrintf("Error loading workspace: %v\n", err)
		return 1
	}

	if len(ws.Repositories) == 0 {
		a.println("No repositories configured to sync.")
		return 0
	}

	a.printf("Syncing %d repositories in workspace %q...\n", len(ws.Repositories), ws.Name)
	results := runner.SyncWorkspace(ctx, ws, a.GitClient, runner.SyncOptions{
		Concurrency: concurrency,
		Pull:        pull,
	})

	hasErrors := false
	for _, res := range results {
		if res.Err != nil {
			hasErrors = true
			a.errPrintf("[-] %s: %s (%v)\n", res.Repo.Name, res.Action, res.Err)
		} else {
			extra := ""
			if res.Message != "" {
				extra = fmt.Sprintf(" - %s", res.Message)
			}
			a.printf("[+] %s: %s%s\n", res.Repo.Name, res.Action, extra)
		}
	}

	if hasErrors {
		return 1
	}
	return 0
}

func (a *App) cmdStatus(ctx context.Context, args []string, configPath string) int {
	var jsonOutput bool
	fs := flag.NewFlagSet("status", flag.ContinueOnError)
	fs.SetOutput(a.Stderr)
	fs.BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	if _, err := parseFlagsAndArgs(fs, args); err != nil {
		return 2
	}

	if a.GitClient == nil {
		a.errPrintln("Error: git client is not available")
		return 1
	}

	ws, _, err := config.LoadWorkspace(a.WorkingDir, configPath)
	if err != nil {
		a.errPrintf("Error loading workspace: %v\n", err)
		return 1
	}

	if len(ws.Repositories) == 0 {
		a.println("No repositories in workspace.")
		return 0
	}

	results := runner.StatusWorkspace(ctx, ws, a.GitClient)

	if jsonOutput {
		enc := json.NewEncoder(a.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(results); err != nil {
			a.errPrintf("Error encoding JSON: %v\n", err)
			return 1
		}
		return 0
	}

	w := tabwriter.NewWriter(a.Stdout, 0, 0, 3, ' ', 0)
	_, _ = fmt.Fprintln(w, "REPOSITORY\tCLONED\tBRANCH\tCLEAN\tCHANGES")
	for _, res := range results {
		if !res.Exists {
			_, _ = fmt.Fprintf(w, "%s\tno\t-\t-\t-\n", res.Repo.Name)
			continue
		}
		cleanStr := "yes"
		if !res.IsClean {
			cleanStr = "dirty"
		}
		summary := "-"
		if res.Status != "" {
			lines := strings.Split(strings.TrimSpace(res.Status), "\n")
			summary = fmt.Sprintf("%d files modified", len(lines))
		}
		_, _ = fmt.Fprintf(w, "%s\tyes\t%s\t%s\t%s\n", res.Repo.Name, res.Branch, cleanStr, summary)
	}
	_ = w.Flush()
	return 0
}

func (a *App) cmdExec(ctx context.Context, args []string, configPath string) int {
	var concurrency int
	fs := flag.NewFlagSet("exec", flag.ContinueOnError)
	fs.SetOutput(a.Stderr)
	fs.IntVar(&concurrency, "concurrency", 4, "Concurrency limit")

	pos, err := parseFlagsAndArgs(fs, args)
	if err != nil {
		return 2
	}

	if len(pos) == 0 {
		a.errPrintln("Error: command is required. Usage: projectctl exec [flags] -- <command> [args...]")
		return 2
	}

	ws, _, err := config.LoadWorkspace(a.WorkingDir, configPath)
	if err != nil {
		a.errPrintf("Error loading workspace: %v\n", err)
		return 1
	}

	if len(ws.Repositories) == 0 {
		a.println("No repositories in workspace.")
		return 0
	}

	command := pos[0]
	execArgs := pos[1:]

	results := runner.ExecInWorkspace(ctx, ws, command, execArgs, concurrency)

	hasErrors := false
	for _, res := range results {
		a.printf("==> Repository: %s (exit code %d)\n", res.Repo.Name, res.ExitCode)
		if res.Stdout != "" {
			a.print(res.Stdout)
		}
		if res.Stderr != "" {
			a.print(res.Stderr)
		}
		if res.ExitCode != 0 {
			hasErrors = true
		}
	}

	if hasErrors {
		return 1
	}
	return 0
}

func (a *App) cmdVersion(args []string) int {
	var jsonOutput bool
	fs := flag.NewFlagSet("version", flag.ContinueOnError)
	fs.SetOutput(a.Stderr)
	fs.BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	if _, err := parseFlagsAndArgs(fs, args); err != nil {
		return 2
	}

	if jsonOutput {
		enc := json.NewEncoder(a.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(version.GetInfo()); err != nil {
			a.errPrintf("Error encoding JSON: %v\n", err)
			return 1
		}
		return 0
	}

	a.println(version.String())
	return 0
}
