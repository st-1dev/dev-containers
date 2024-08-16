package commands

import (
	"context"
	"dev-runner/pkg/dev/naming"
	"flag"
	"fmt"
	"log"
	"os"

	"dev-runner/pkg/conainer/management"
	"dev-runner/pkg/conainer/management/creator"

	"github.com/google/subcommands"

	fp "dev-runner/pkg/filepath"
)

type LogsCmd struct {
	containerManagerName string
	imageTag             string
	hostWorkDirPath      string
}

func (*LogsCmd) Name() string {
	return "logs"
}

func (*LogsCmd) Synopsis() string {
	return "logs."
}

func (*LogsCmd) Usage() string {
	return `
`
}

func (p *LogsCmd) SetFlags(f *flag.FlagSet) {
	workDir, _ := os.Getwd()

	f.StringVar(&p.containerManagerName, "cm", "docker", "Containers manager. Values: docker or podman.")
	f.StringVar(&p.imageTag, "image", "", "Dev image tag.")
	f.StringVar(&p.hostWorkDirPath, "workDir", workDir, "Work dir on host to mount inside container.")
}

func (p *LogsCmd) validateCliArguments() (err error) {
	if p.imageTag == "" {
		return fmt.Errorf("'image' must be set with image tag")
	}
	if !fp.IsDir(p.hostWorkDirPath) {
		return fmt.Errorf("'workDir' must be exists and be directory")
	}
	return nil
}

func (p *LogsCmd) execute(ctx context.Context, _ *flag.FlagSet) (err error) {
	err = p.validateCliArguments()
	if err != nil {
		return fmt.Errorf("command line validation failed: %w", err)
	}

	var manager management.ContainerManager
	manager, err = creator.CreateContainerManager(p.containerManagerName)
	if err != nil {
		return fmt.Errorf("cannot create container manager: %w", err)
	}

	err = manager.Init(ctx)
	if err != nil {
		return fmt.Errorf("container manager initialization failed: %w", err)
	}

	containerName := naming.GenContainerName(p.imageTag, p.hostWorkDirPath)

	err = manager.PrintContainerLogs(ctx, containerName, os.Stdout)
	if err != nil {
		return fmt.Errorf("print container logs failed: %w", err)
	}

	return nil
}

func (p *LogsCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	err := p.execute(ctx, f)
	if err != nil {
		log.Fatalf("got error: %s\n", err.Error())
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
