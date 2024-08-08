package management

import (
	"context"
	"io"
)

type ContainerManager interface {
	Init(ctx context.Context) (err error)

	LoadImage(
		ctx context.Context,
		r io.Reader,
	) (err error)

	GetImageLabels(
		ctx context.Context,
		imageName string,
	) (labels []Label, err error)

	RunContainer(
		ctx context.Context,
		imageName string,
		containerName string,
		mountPoints []MountPoint,
		environmentVariables []EnvironmentVariable,
		portBindings []PortBinding,
		networkMode NetworkMode,
	) (containerId string, err error)

	StopContainer(
		ctx context.Context,
		containerName string,
	) (err error)

	PrintContainerLogs(
		ctx context.Context,
		containerName string,
		w io.Writer,
	) (err error)
}
