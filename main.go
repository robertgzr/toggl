package main // import "github.com/robertgzr/toggl"

import (
	"fmt"
	"os"

	"github.com/robertgzr/toggl/app"
)

var (
	version = "undefined"
	commit  = "undefined"
	date    = "undefined"
)

func main() {
	app := app.New(version, commit, date)
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "toggl: %s\n", err)
		os.Exit(1)
	}
}
