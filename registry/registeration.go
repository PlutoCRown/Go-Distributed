package registry

type Registeration struct {
	ServiceName     ServiceName
	ServiceURL      string
	RequiredService []ServiceName
	// 如果依赖服务有变化，将会通知此接口
	ServiceUpdateURL string
	HeartbeatURL     string
}

type ServiceName string

const (
	LogService   = ServiceName("Log Service")
	GradeService = ServiceName("Grade Service")
)

type patchEntry struct {
	Name ServiceName
	URL  string
}

type patch struct {
	Added   []patchEntry
	Removed []patchEntry
}
