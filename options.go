package vio

type TypeOfBackend string

const (
	Posix TypeOfBackend = "Posix"

// Git                = "Git"
)

type Options struct {
	// Type of object database to use. Default: Memory
	BackendType TypeOfBackend

	// Path to snapshots
	SnapshotsPath string

	// Path to the repository
	RepoPath string

	// Config file
	ConfigFile string
}

func NewOptions() (o Options) {
	o.BackendType = Posix
	o.SnapshotsPath = "./.snapshots"
	o.RepoPath = "."
	o.ConfigFile = ".vioconfig"
	return
}
