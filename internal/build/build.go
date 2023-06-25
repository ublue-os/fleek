package build

// Variables in this file are set via ldflags.
var (
	IsDev = Version == "0.0.0-dev"

	Version    = "0.0.0-dev"
	Commit     = "none"
	CommitDate = "unknown"
)
