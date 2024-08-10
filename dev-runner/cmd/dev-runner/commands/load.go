package commands

import (
	"context"
	"flag"
	"log"

	"github.com/google/subcommands"
)

type LoadCmd struct {
	containerManagerName string
}

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

func (p *LoadCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.containerManagerName, "cm", "docker", "Containers manager. Values: docker or podman.")
}

func (p *LoadCmd) execute(ctx context.Context, _ *flag.FlagSet) (err error) {
	// TODO: wire code here
	//var manager management.ContainerManager
	//manager, err = creator.CreateContainerManager(p.containerManagerName)
	//if err != nil {
	//	return fmt.Errorf("cannot create container manager: %w", err)
	//}

	return nil
}

func (p *LoadCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	err := p.execute(ctx, f)
	if err != nil {
		log.Fatalf("got error: %s\n", err.Error())
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
