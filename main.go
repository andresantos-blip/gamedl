package main

import (
	"gamedl/cmd"
	"gamedl/lib/app/build"
)

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
