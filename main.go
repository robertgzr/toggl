package main // import "github.com/robertgzr/toggl"

import (
	"fmt"
	"os"

	"github.com/robertgzr/toggl/app"
)

func main() {
	app := app.New()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "toggl: %s\n", err)
		os.Exit(1)
	}
}
