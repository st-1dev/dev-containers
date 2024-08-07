package management

import (
	"context"
	"io"
)

type ContainerManager interface {
	Init(ctx context.Context) (err error)

	LoadImage(ctx context.Context, r io.Reader) (err error)

	RunContainer(
		ctx context.Context,
		imageTag string,
		containerName string,
		mountPoints map[string]string,
		environmentVariables map[string]string,
		portBindings map[int]int,
		networkMode string,
	) (containerId string, err error)

	StopContainer(ctx context.Context, containerId string) (err error)

	PrintContainerLogs(ctx context.Context, containerId string, w io.Writer) (err error)
}
