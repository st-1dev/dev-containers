package management

type MountPoint struct {
	HostPath      string
	ContainerPath string
	ReadOnly      bool
}

type EnvironmentVariable struct {
	Name  string
	Value string
}

type PortBinding struct {
	ContainerPort int
	HostPort      int
}

type Label struct {
	Name  string
	Value string
}

type NetworkMode string

const (
	NetworkBridge NetworkMode = "bridge"
	NetworkHost   NetworkMode = "host"
	NetworkNat    NetworkMode = "nat"
)

func GetNetworkModes() []NetworkMode {
	return []NetworkMode{NetworkBridge, NetworkHost, NetworkNat}
}
