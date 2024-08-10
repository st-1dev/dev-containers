package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/docker/docker/api/types"

	"dev-runner/pkg/conainer/management"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"

	"github.com/docker/docker/client"
)

type dockerManager struct {
	con *client.Client
}

func NewDockerManager() management.ContainerManager {
	return &dockerManager{}
}

func (m *dockerManager) Init(_ context.Context) (err error) {
	var con *client.Client
	con, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("cannot create client: %w", err)
	}

	m.con = con
	return err
}

func (m *dockerManager) LoadImage(
	ctx context.Context,
	r io.Reader,
) (err error) {
	_, err = m.con.ImageLoad(ctx, r, true)
	if err != nil {
		return fmt.Errorf("cannot load image: %w", err)
	}
	return nil
}

func (m *dockerManager) GetImageLabels(
	ctx context.Context,
	imageName string,
) (labels []management.Label, err error) {
	var inspect types.ImageInspect
	inspect, _, err = m.con.ImageInspectWithRaw(ctx, imageName)
	if err != nil {
		return nil, fmt.Errorf("cannot inspect image '%s': %w", err)
	}

	if len(inspect.Config.Labels) > 0 {
		labels = make([]management.Label, len(inspect.Config.Labels))
		for name, value := range inspect.Config.Labels {
			labels = append(labels, management.Label{Name: name, Value: value})
		}
	}
	return labels, nil
}

func (m *dockerManager) RunContainer(
	ctx context.Context,
	imageName string,
	containerName string,
	mountPoints []management.MountPoint,
	environmentVariables []management.EnvironmentVariable,
	portBindings []management.PortBinding,
	networkMode management.NetworkMode,
) (containerId string, err error) {
	var environmentVariables_ []string
	for _, item := range environmentVariables {
		environmentVariables_ = append(
			environmentVariables_,
			fmt.Sprintf("%s=%s", os.ExpandEnv(item.Name), os.ExpandEnv(item.Value)),
		)
	}

	var mountPoints_ []mount.Mount
	for _, item := range mountPoints {
		mountPoints_ = append(
			mountPoints_,
			mount.Mount{
				Type:     mount.TypeBind,
				Source:   os.ExpandEnv(item.HostPath),
				Target:   os.ExpandEnv(item.ContainerPath),
				ReadOnly: item.ReadOnly,
			},
		)
	}

	var portBindings_ nat.PortMap
	if len(portBindings) > 0 {
		portBindings_ = make(nat.PortMap, len(portBindings))
		for _, item := range portBindings {
			containerPort := nat.Port(strconv.Itoa(item.ContainerPort))
			hostPort := strconv.Itoa(item.HostPort)
			portBindings_[containerPort] = []nat.PortBinding{{
				HostPort: hostPort,
			}}
		}
	}

	var networkMode_ container.NetworkMode
	networkMode_, err = getNetworkMode(networkMode)
	if err != nil {
		return "", err
	}

	var containerResp_ container.CreateResponse
	containerResp_, err = m.con.ContainerCreate(
		ctx,
		&container.Config{
			Image: imageName,
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
		return "", fmt.Errorf("cannot create container from image '%s': %w", imageName, err)
	}

	containerId = containerResp_.ID

	err = m.con.ContainerStart(ctx, containerId, container.StartOptions{})
	if err != nil {
		return "", fmt.Errorf("cannot start container '%s' from image '%s': %w", containerId, imageName, err)
	}

	return containerId, nil
}

func (m *dockerManager) StopContainer(
	ctx context.Context,
	containerName string,
) (err error) {
	err = m.con.ContainerStop(ctx, containerName, container.StopOptions{})
	if err != nil {
		return fmt.Errorf("cannot stop container '%s': %w", containerName, err)
	}

	err = m.con.ContainerRemove(ctx, containerName, container.RemoveOptions{RemoveVolumes: true})
	if err != nil {
		return fmt.Errorf("cannot remove container '%s': %w", containerName, err)
	}

	return nil
}

func (m *dockerManager) PrintContainerLogs(
	ctx context.Context,
	containerName string,
	w io.Writer,
) (err error) {
	var reader_ io.ReadCloser
	reader_, err = m.con.ContainerLogs(
		ctx,
		containerName,
		container.LogsOptions{
			ShowStderr: true,
			ShowStdout: true,
			Timestamps: true,
		},
	)
	if err != nil {
		return fmt.Errorf("cannot get container logs from '%s': %w", containerName, err)
	}
	defer func() {
		_ = reader_.Close()
	}()

	_, err = io.Copy(w, reader_)
	if err != nil {
		return fmt.Errorf("cannot print container logs from '%s': %w", containerName, err)
	}

	return nil
}

func getNetworkMode(mode management.NetworkMode) (networkMode container.NetworkMode, err error) {
	switch mode {
	case management.NetworkBridge:
		return network.NetworkBridge, nil
	case management.NetworkHost:
		return network.NetworkHost, nil
	default:
		return network.NetworkDefault, fmt.Errorf("network mode '%v' is not supported", mode)
	}
}
