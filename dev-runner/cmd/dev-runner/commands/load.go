package commands

import (
	"context"
	"flag"

	"github.com/google/subcommands"
)

type LoadCmd struct{}

func (*LoadCmd) Name() string {
	return "load"
}

func (*LoadCmd) Synopsis() string {
	return "load."
}

func (*LoadCmd) Usage() string {
	return `
`
}

func (p *LoadCmd) SetFlags(_ *flag.FlagSet) {
}

func (p *LoadCmd) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	return subcommands.ExitSuccess
}
