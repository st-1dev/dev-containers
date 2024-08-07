package commands

import (
	"context"
	"flag"

	"github.com/google/subcommands"
)

type UnloadCmd struct{}

func (*UnloadCmd) Name() string {
	return "unload"
}

func (*UnloadCmd) Synopsis() string {
	return "unload."
}

func (*UnloadCmd) Usage() string {
	return `
`
}

func (p *UnloadCmd) SetFlags(_ *flag.FlagSet) {
}

func (p *UnloadCmd) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	return subcommands.ExitSuccess
}
