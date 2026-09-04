package buildinfo

// These values are replaced by GoReleaser. Keeping useful development values
// makes local builds honest without requiring linker flags.
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

// Info describes the binary that is currently running.
type Info struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

// Current returns build information embedded in the binary.
func Current() Info {
	return Info{
		Version: version,
		Commit:  commit,
		Date:    date,
	}
}
