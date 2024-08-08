package commands

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"dev-runner/pkg/conainer/management/docker"

	"dev-runner/pkg/dev"

	"github.com/google/subcommands"

	fp "dev-runner/pkg/filepath"
)

type StopCmd struct {
	imageTag        string
	hostWorkDirPath string
}

func (*StopCmd) Name() string {
	return "stop"
}

func (*StopCmd) Synopsis() string {
	return "stop."
}

func (*StopCmd) Usage() string {
	return `
`
}

func (p *StopCmd) SetFlags(f *flag.FlagSet) {
	workDir, _ := os.Getwd()

	f.StringVar(&p.imageTag, "image", "", "Dev image tag.")
	f.StringVar(&p.hostWorkDirPath, "workDir", workDir, "Work dir on host to mount inside container.")
}

func (p *StopCmd) validateCliArguments() (err error) {
	if p.imageTag == "" {
		return fmt.Errorf("'image' must be set with image tag")
	}
	if !fp.IsDir(p.hostWorkDirPath) {
		return fmt.Errorf("'workDir' must be exists and be directory")
	}
	return nil
}

func (p *StopCmd) Execute(ctx context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	var err error
	err = p.validateCliArguments()
	if err != nil {
		log.Fatalf("command line validation failed: %s\n", err.Error())
		return subcommands.ExitFailure
	}

	manager := docker.NewDockerManager()
	err = manager.Init(ctx)
	if err != nil {
		log.Fatalf("docker manager initialization failed: %s\n", err.Error())
		return subcommands.ExitFailure
	}

	containerName := dev.GenContainerName(p.imageTag, p.hostWorkDirPath)

	err = manager.StopContainer(ctx, containerName)
	if err != nil {
		log.Fatalf("stop container failed: %s\n", err.Error())
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
