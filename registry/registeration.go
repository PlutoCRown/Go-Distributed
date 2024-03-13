package registry

type Registeration struct {
	ServiceName ServiceName
	ServiceURL string
}

type ServiceName string

const (
	LogService = ServiceName("LogService")
)