package commands

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"runner/pkg/conainer/management"
	"runner/pkg/dev"

	"github.com/google/subcommands"

	fp "runner/pkg/filepath"
)

type LogsCmd struct {
	imageTag        string
	hostWorkDirPath string
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

func (p *LogsCmd) Execute(ctx context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	var err error
	err = p.validateCliArguments()
	if err != nil {
		log.Fatalf("command line validation failed: %s\n", err.Error())
		return subcommands.ExitFailure
	}

	manager := management.NewDockerManager()
	err = manager.Init(ctx)
	if err != nil {
		log.Fatalf("docker manager initialization failed: %s\n", err.Error())
		return subcommands.ExitFailure
	}

	containerName := dev.GenContainerName(p.imageTag, p.hostWorkDirPath)

	err = manager.PrintContainerLogs(ctx, containerName, os.Stdout)
	if err != nil {
		log.Fatalf("print container logs failed: %s\n", err.Error())
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
