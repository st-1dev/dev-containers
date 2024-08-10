package podman

import (
	"context"
	"fmt"
	"io"
	"os"

	"dev-runner/pkg/conainer/management"

	nettypes "github.com/containers/common/libnetwork/types"
	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/bindings/images"
	"github.com/containers/podman/v5/pkg/domain/entities/types"
	"github.com/containers/podman/v5/pkg/specgen"
	"github.com/opencontainers/runtime-spec/specs-go"
)

type podmanManager struct {
	conCtx context.Context
}

func NewPodmanManager() management.ContainerManager {
	return &podmanManager{}
}

func (m *podmanManager) Init(ctx context.Context) (err error) {
	var conCtx context.Context
	conCtx, err = bindings.NewConnection(ctx, "")
	if err != nil {
		return fmt.Errorf("cannot create posman connection: %w", err)
	}

	m.conCtx = conCtx
	return err
}

func (m *podmanManager) LoadImage(
	_ context.Context,
	r io.Reader,
) (err error) {
	_, err = images.Load(m.conCtx, r)
	if err != nil {
		return fmt.Errorf("cannot load image: %w", err)
	}
	return nil
}

func (m *podmanManager) GetImageLabels(
	_ context.Context,
	imageName string,
) (labels []management.Label, err error) {
	var inspect *types.ImageInspectReport
	inspect, err = images.GetImage(m.conCtx, imageName, nil)
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

func (m *podmanManager) RunContainer(
	_ context.Context,
	imageName string,
	containerName string,
	mountPoints []management.MountPoint,
	environmentVariables []management.EnvironmentVariable,
	portBindings []management.PortBinding,
	networkMode management.NetworkMode,
) (containerId string, err error) {
	var environmentVariables_ map[string]string
	for _, item := range environmentVariables {
		environmentVariables_[os.ExpandEnv(item.Name)] = os.ExpandEnv(item.Value)
	}

	var mountPoints_ []specs.Mount
	for _, item := range mountPoints {
		mountPoints_ = append(
			mountPoints_,
			specs.Mount{
				Destination: item.ContainerPath,
				Source:      item.HostPath,
			},
		)
	}

	var portBindings_ []nettypes.PortMapping
	if len(portBindings) > 0 {
		for _, item := range portBindings {
			portBindings_ = append(
				portBindings_,
				nettypes.PortMapping{
					ContainerPort: uint16(item.ContainerPort),
					HostPort:      uint16(item.HostPort),
				},
			)
		}
	}

	var networkMode_ specgen.NamespaceMode
	networkMode_, err = getNetworkMode(networkMode)
	if err != nil {
		return "", err
	}

	s := specgen.NewSpecGenerator(imageName, false)
	s.Name = containerName
	s.CapAdd = []string{
		"CAP_AUDIT_WRITE",
		"SYS_PTRACE",
		"NET_RAW",
		"NET_ADMIN",
	}
	s.Env = environmentVariables_
	s.Mounts = mountPoints_
	s.PortMappings = portBindings_
	s.SeccompPolicy = "unconfined"
	s.NetNS.NSMode = networkMode_

	var containerResp_ types.ContainerCreateResponse
	containerResp_, err = containers.CreateWithSpec(m.conCtx, s, nil)
	if err != nil {
		return "", fmt.Errorf("cannot create container from image '%s': %w", imageName, err)
	}

	containerId = containerResp_.ID

	err = containers.Start(m.conCtx, containerName, nil)
	if err != nil {
		return "", fmt.Errorf("cannot start container '%s' from image '%s': %w", containerId, imageName, err)
	}

	return containerId, nil
}

func (m *podmanManager) StopContainer(
	_ context.Context,
	containerName string,
) (err error) {
	err = containers.Stop(m.conCtx, containerName, nil)
	if err != nil {
		return fmt.Errorf("cannot stop container '%s': %w", containerName, err)
	}

	removeOptions := new(containers.RemoveOptions)
	removeOptions.WithVolumes(true)
	_, err = containers.Remove(m.conCtx, containerName, removeOptions)
	if err != nil {
		return fmt.Errorf("cannot remove container '%s': %w", containerName, err)
	}

	return nil
}

func (m *podmanManager) PrintContainerLogs(
	_ context.Context,
	containerName string,
	w io.Writer,
) (err error) {
	done := make(chan bool)
	stdOut := make(chan string, 20)
	stdErr := make(chan string, 20)

	go func() {
		for {
			select {
			case msg := <-stdOut:
				_, _ = w.Write([]byte(msg))
			case msg := <-stdErr:
				_, _ = w.Write([]byte(msg))
			case <-done:
				return
			}
		}
	}()

	options := new(containers.LogOptions)
	options.
		WithStderr(true).
		WithStdout(true).
		WithTimestamps(true)

	err = containers.Logs(
		m.conCtx,
		containerName,
		options,
		stdOut,
		stdErr,
	)
	if err != nil {
		return fmt.Errorf("cannot get container logs from '%s': %w", containerName, err)
	}
	done <- true

	return nil
}

func getNetworkMode(mode management.NetworkMode) (network specgen.NamespaceMode, err error) {
	switch mode {
	case management.NetworkBridge:
		return specgen.Bridge, nil
	case management.NetworkHost:
		return specgen.Host, nil
	default:
		return specgen.Default, fmt.Errorf("network mode '%v' is not supported", mode)
	}
}
