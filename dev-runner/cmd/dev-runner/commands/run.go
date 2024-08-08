package commands

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"dev-runner/pkg/conainer/management"
	"dev-runner/pkg/conainer/management/docker"

	"dev-runner/pkg/dev"

	"github.com/google/subcommands"

	fp "dev-runner/pkg/filepath"
)

type RunCmd struct {
	imageTag         string
	hostWorkDirPath  string
	hostHomeDir      string
	user             string
	host             string
	containerSshPort int
	networkMode      string
	interactive      bool
	sort.StringSlice
}

func (*RunCmd) Name() string {
	return "run"
}

func (*RunCmd) Synopsis() string {
	return "run."
}

func (*RunCmd) Usage() string {
	return `
`
}

func (p *RunCmd) SetFlags(f *flag.FlagSet) {
	workDir, _ := os.Getwd()
	homeDir, _ := os.UserHomeDir()

	f.StringVar(&p.imageTag, "image", "", "Dev image tag.")
	f.StringVar(&p.hostWorkDirPath, "workDir", workDir, "Work dir on host to mount inside container.")
	f.StringVar(&p.hostHomeDir, "homeDir", homeDir, "Home dir on host to mount directories(.ssh, .docker and etc) inside container.")
	f.StringVar(&p.user, "user", "user", "The container user username.")
	f.StringVar(&p.host, "host", "localhost", "The host to bind containers ports to.")
	f.IntVar(&p.containerSshPort, "containerSshPort", 2221, "The SSH port to bind from container.")
	f.StringVar(&p.networkMode, "network", "host", "The network mode for container.")
	f.BoolVar(&p.interactive, "interactive", false, "Run container in interactive mode to debug.")
}

func (p *RunCmd) validateCliArguments() (err error) {
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

func (p *RunCmd) Execute(ctx context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
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

	var networkMode management.NetworkMode
	networkMode, err = getNetworkMode(p.networkMode)
	if err != nil {
		log.Fatalf("%s\n", err.Error())
		return subcommands.ExitFailure
	}

	var mountPoints []management.MountPoint
	mountPoints, err = getMountPoints(p.imageTag, p.hostWorkDirPath, p.hostHomeDir, p.user)
	if err != nil {
		log.Fatalf("cannot get mount points: %s\n", err.Error())
		return subcommands.ExitFailure
	}

	var environmentVariables []management.EnvironmentVariable
	var portBindings []management.PortBinding

	var containerId string
	containerId, err = manager.RunContainer(
		ctx,
		p.imageTag,
		containerName,
		mountPoints,
		environmentVariables,
		portBindings,
		networkMode,
	)
	if err != nil {
		log.Fatalf("start container failed: %s\n", err.Error())
		return subcommands.ExitFailure
	}

	log.Printf("container started '%s'\n", containerId)

	return subcommands.ExitSuccess
}

func getMountPoints(
	imageTag string,
	hostWorkDir string,
	hostHomeDir string,
	userInsideContainer string,
) (mountPoints []management.MountPoint, err error) {
	devHomeDirName := dev.GenDevHomeDirName(imageTag, hostWorkDir)
	devHomeDir := filepath.Join(hostWorkDir, "..", devHomeDirName)

	homeDirInsideContainer := fmt.Sprintf("/home/%s", userInsideContainer)

	dirsMustBeMounted := map[string]string{
		filepath.Join(homeDirInsideContainer, ".cache"):  filepath.Join(devHomeDir, ".cache"),
		filepath.Join(homeDirInsideContainer, ".config"): filepath.Join(devHomeDir, ".config"),
		filepath.Join(homeDirInsideContainer, ".java"):   filepath.Join(devHomeDir, ".java"),
		filepath.Join(homeDirInsideContainer, ".jdks"):   filepath.Join(devHomeDir, ".jdks"),
		filepath.Join(homeDirInsideContainer, ".local"):  filepath.Join(devHomeDir, ".local"),
		filepath.Join(homeDirInsideContainer, ".m2"):     filepath.Join(devHomeDir, ".m2"),
		filepath.Join("/", "work"):                       hostWorkDir,
	}
	for containerDirPath, hostDirPath := range dirsMustBeMounted {
		err = fp.MakePaths(hostDirPath)
		if err != nil {
			return nil, fmt.Errorf("cannot create required directory '%s': %w", hostDirPath, err)
		}
		mountPoints = append(
			mountPoints,
			management.MountPoint{
				HostPath:      hostDirPath,
				ContainerPath: containerDirPath,
			},
		)
	}

	filesMustBeMounted := map[string]string{
		filepath.Join(homeDirInsideContainer, ".bash_history"): filepath.Join(devHomeDir, ".bash_history"),
	}
	for containerFilePath, hostFilePath := range filesMustBeMounted {
		err = fp.MakeFiles(hostFilePath)
		if err != nil {
			return nil, fmt.Errorf("cannot create required file '%s': %w", hostFilePath, err)
		}
		mountPoints = append(
			mountPoints,
			management.MountPoint{
				HostPath:      hostFilePath,
				ContainerPath: containerFilePath,
			},
		)
	}

	dirsMaybeMounted := map[string]string{
		filepath.Join(homeDirInsideContainer, ".ssh"):    filepath.Join(hostHomeDir, ".ssh"),
		filepath.Join(homeDirInsideContainer, ".docker"): filepath.Join(hostHomeDir, ".docker"),
	}
	for containerDirPath, hostDirPath := range dirsMaybeMounted {
		if !fp.IsDir(hostDirPath) {
			continue
		}
		mountPoints = append(
			mountPoints,
			management.MountPoint{
				HostPath:      hostDirPath,
				ContainerPath: containerDirPath,
			},
		)
	}

	filesMaybeMounted := map[string]string{
		filepath.Join(homeDirInsideContainer, ".gitconfig"): filepath.Join(hostHomeDir, ".gitconfig"),
		filepath.Join("/", "var", "run", "docker.sock"):     "/var/run/docker.sock",
	}
	for containerFilePath, hostFilePath := range filesMaybeMounted {
		if !fp.IsFile(hostFilePath) {
			continue
		}
		mountPoints = append(
			mountPoints,
			management.MountPoint{
				HostPath:      hostFilePath,
				ContainerPath: containerFilePath,
			},
		)
	}

	return mountPoints, nil
}

func getNetworkMode(value string) (networkMode management.NetworkMode, err error) {
	for _, item := range management.GetNetworkModes() {
		if string(item) == value {
			return item, nil
		}
	}
	return management.NetworkHost, fmt.Errorf("incorrect network mode '%s'", value)
}
