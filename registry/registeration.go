package registry

type Registeration struct {
	ServiceName ServiceName
	ServiceURL string
	RequiredService []ServiceName
	ServiceUpdateURL string
}

type ServiceName string

const (
	LogService = ServiceName("LogService")
	GradeService = ServiceName("GradeService")
)

type patchEntry struct {
	Name ServiceName
	URL string
}

type patch struct {
	Added []patchEntry
	Removed []patchEntry
}