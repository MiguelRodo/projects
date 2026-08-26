// Command projects is an alias executable for projectctl.
package main

import (
	"os"

	"github.com/MiguelRodo/projects/internal/cli"
)

func main() {
	exitCode := cli.Run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(exitCode)
}
