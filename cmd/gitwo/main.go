package main

import "github.com/gitwohq/gitwo/cmd"

var (
	// Set by -ldflags: -X main.version=... etc.
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// hand build info to the CLI layer (see SetBuildInfo below)
	cmd.SetBuildInfo(version, commit, date)
	cmd.Execute()
}
