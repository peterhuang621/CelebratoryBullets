package registry

type Registration struct {
	ServiceName      ServiceName
	ServiceURL       string
	RequiredServices []ServiceName
	ServiceUpdateURL string
	HeartBeatURL     string
}

type ServiceName string

const (
	LogService    = ServiceName("Log Service")
	GradingServie = ServiceName("Grading Service")
)

type patchEntry struct {
	Name ServiceName
	URL  string
}

type patch struct {
	Added   []patchEntry
	Removed []patchEntry
}
