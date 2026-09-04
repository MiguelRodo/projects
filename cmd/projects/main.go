package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/MiguelRodo/projects/internal/cli"
	"github.com/MiguelRodo/projects/internal/githubcli"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	os.Exit(cli.Run(ctx, os.Args[1:], os.Stdout, os.Stderr, githubcli.ExecRunner{}))
}
