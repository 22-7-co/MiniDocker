package container

type ContainerInfo struct {
	Pid          string   `json:"pid"`
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Command      string   `json:"command"`
	CreateTime   string   `json:"create_time"`
	Status       string   `json:"status"`
	PortMappings []string `json:"port_mappings"`
}

var (
	RUNNING             = "running"
	STOP                = "stop"
	EXIT                = "exited"
	DefaultInfoLocation = "/var/run/mydocker/%s/"
	ConfigName          = "config.json"
	ContainerLogFile    = "container.json"
	RootUrl             = "var/lib/mydocker/container"
	BusyboxPath         = "/opt/busybox"
)
