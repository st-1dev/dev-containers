package management

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"

	"github.com/docker/docker/client"
)

type dockerManager struct {
	dockerClient *client.Client
}

func NewDockerManager() ContainerManager {
	return &dockerManager{}
}

func (m *dockerManager) Init(ctx context.Context) (err error) {
	var client_ *client.Client
	client_, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("cannot create docker client: %w", err)
	}

	m.dockerClient = client_
	return err
}

func (m *dockerManager) LoadImage(ctx context.Context, r io.Reader) (err error) {
	_, err = m.dockerClient.ImageLoad(ctx, r, true)
	if err != nil {
		return fmt.Errorf("cannot load docker image: %w", err)
	}
	return nil
}

func (m *dockerManager) RunContainer(
	ctx context.Context,
	imageTag string,
	containerName string,
	mountPoints map[string]string,
	environmentVariables map[string]string,
	portBindings map[int]int,
	networkMode string,
) (containerId string, err error) {
	var environmentVariables_ []string
	for name, value := range environmentVariables {
		expandedName := os.ExpandEnv(name)
		expandedValue := os.ExpandEnv(value)
		environmentVariables_ = append(environmentVariables_, fmt.Sprintf("%s=%s", expandedName, expandedValue))
	}

	var mountPoints_ []mount.Mount
	for pathInContainer, pathOnHost := range mountPoints {
		source := os.ExpandEnv(pathOnHost)
		target := os.ExpandEnv(pathInContainer)
		mountPoints_ = append(mountPoints_, mount.Mount{Type: mount.TypeBind, Source: source, Target: target})
	}

	var portBindings_ nat.PortMap
	if len(portBindings) > 0 {
		portBindings_ = make(nat.PortMap, len(portBindings))
		for portOnHost, portInContainer := range portBindings {
			portOnHostStr := fmt.Sprintf("%d", portOnHost)
			portInContainerStr := nat.Port(fmt.Sprintf("%d", portInContainer))
			portBindings_[portInContainerStr] = []nat.PortBinding{{HostPort: portOnHostStr}}
		}
	}

	var networkMode_ container.NetworkMode
	networkMode_, err = getNetworkMode(networkMode)
	if err != nil {
		return "", err
	}

	var containerResp_ container.CreateResponse
	containerResp_, err = m.dockerClient.ContainerCreate(
		ctx,
		&container.Config{
			Image: imageTag,
			Env:   environmentVariables_,
		},
		&container.HostConfig{
			Mounts:       mountPoints_,
			NetworkMode:  networkMode_,
			PortBindings: portBindings_,
			CapAdd: strslice.StrSlice{
				"CAP_AUDIT_WRITE",
				"SYS_PTRACE",
				"NET_RAW",
				"NET_ADMIN",
			},
			SecurityOpt: []string{"seccomp=unconfined"},
		},
		nil,
		nil,
		containerName,
	)
	if err != nil {
		return "", fmt.Errorf("cannot create docker container for image '%s': %w", imageTag, err)
	}

	containerId = containerResp_.ID

	err = m.dockerClient.ContainerStart(ctx, containerId, container.StartOptions{})
	if err != nil {
		return "", fmt.Errorf("cannot start container '%s' for image '%s': %w", containerId, imageTag, err)
	}

	return containerId, nil
}

func (m *dockerManager) StopContainer(ctx context.Context, containerId string) (err error) {
	err = m.dockerClient.ContainerStop(ctx, containerId, container.StopOptions{})
	if err != nil {
		return fmt.Errorf("cannot stop docker container with id '%s': %w", containerId, err)
	}

	err = m.dockerClient.ContainerRemove(ctx, containerId, container.RemoveOptions{RemoveVolumes: true})
	if err != nil {
		return fmt.Errorf("cannot remove docker container with id '%s': %w", containerId, err)
	}

	return nil
}

func (m *dockerManager) PrintContainerLogs(ctx context.Context, containerId string, w io.Writer) (err error) {
	var reader_ io.ReadCloser
	reader_, err = m.dockerClient.ContainerLogs(
		ctx,
		containerId,
		container.LogsOptions{
			ShowStderr: true,
			ShowStdout: true,
			Timestamps: true,
		})
	if err != nil {
		return fmt.Errorf("cannot get docker container logs with id '%s': %w", containerId, err)
	}
	defer func() {
		_ = reader_.Close()
	}()

	_, err = io.Copy(w, reader_)
	if err != nil {
		return fmt.Errorf("cannot print docker container logs with id '%s': %w", containerId, err)
	}

	return nil
}

func getNetworkMode(mode string) (networkMode container.NetworkMode, err error) {
	switch mode {
	case "default":
		return network.NetworkDefault, nil
	case "host":
		return network.NetworkHost, nil
	case "none":
		return network.NetworkNone, nil
	case "bridge":
		return network.NetworkBridge, nil
	case "nat":
		return network.NetworkNat, nil
	default:
		return network.NetworkDefault, fmt.Errorf("network mode '%s' is not supported")
	}
}
