package main

import (
	"context"
	"flag"
	"os"

	"dev-runner/cmd/dev-runner/commands"

	"github.com/google/subcommands"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&commands.AttachCmd{}, "")
	subcommands.Register(&commands.LoadCmd{}, "")
	subcommands.Register(&commands.LogsCmd{}, "")
	subcommands.Register(&commands.RunCmd{}, "")
	subcommands.Register(&commands.StopCmd{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
