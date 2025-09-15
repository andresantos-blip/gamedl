package main

import (
	"gamedl/cmd"
	"gamedl/lib/app/build"
)

// Build info variables.
// These should be set using -ldflags during the build process.
// Currently only goreleaser is setting up these variables when a new release/build is created.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.Execute(build.Info{
		Version: version,
		Commit:  commit,
		Date:    date,
	})
}
