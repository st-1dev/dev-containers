package commands

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"

	"runner/pkg/ssh"

	fp "runner/pkg/filepath"
)

type AttachCmd struct {
	imageTag        string
	hostWorkDirPath string
	hostHomeDir     string
	host            string
	port            int
	user            string
	password        string
}

func (*AttachCmd) Name() string {
	return "attach"
}

func (*AttachCmd) Synopsis() string {
	return "attach."
}

func (*AttachCmd) Usage() string {
	return `
`
}

func (p *AttachCmd) SetFlags(f *flag.FlagSet) {
	workDir, _ := os.Getwd()
	homeDir, _ := os.UserHomeDir()

	f.StringVar(&p.imageTag, "image", "", "Dev image tag.")
	f.StringVar(&p.hostWorkDirPath, "workDir", workDir, "Work dir on host to mount inside container.")
	f.StringVar(&p.hostHomeDir, "homeDir", homeDir, "Home dir to use SSH.")
	f.StringVar(&p.host, "host", "localhost", "The host to bind containers ports to.")
	f.IntVar(&p.port, "port", 2221, "The SSH port to bind from container.")
	f.StringVar(&p.user, "user", "user", "The container user username.")
	f.StringVar(&p.password, "password", "user", "The container user password.")
}

func (p *AttachCmd) validateCliArguments() (err error) {
	if p.imageTag == "" {
		return fmt.Errorf("'image' must be set with image tag")
	}
	if !fp.IsDir(p.hostWorkDirPath) {
		return fmt.Errorf("'workDir' must be exists and be directory")
	}
	if !fp.IsDir(p.hostHomeDir) {
		return fmt.Errorf("'homeDir' must be exists and be directory")
	}
	return nil
}

func (p *AttachCmd) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	var err error
	err = p.validateCliArguments()
	if err != nil {
		return subcommands.ExitFailure
	}

	err = ssh.Shell(p.host, p.port, p.user, p.password)
	if err != nil {
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
