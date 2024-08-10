package creator

import (
	"fmt"

	"dev-runner/pkg/conainer/management"
	"dev-runner/pkg/conainer/management/docker"
	"dev-runner/pkg/conainer/management/podman"
)

func CreateContainerManager(name string) (manager management.ContainerManager, err error) {
	switch name {
	case "docker":
		return docker.NewDockerManager(), nil
	case "podman":
		return podman.NewPodmanManager(), nil
	}
	return nil, fmt.Errorf("conatainer manager '%s' is not supported")
}
