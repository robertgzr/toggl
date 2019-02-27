package app

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/robertgzr/toggl/commands/project"
	"github.com/robertgzr/toggl/commands/timer"
)

func init() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s version %s (%s)\n", c.App.Name, c.App.Version, BUILDTIME)
	}
}

func New() *cli.App {
	app := cli.NewApp()
	app.Name = "toggl"
	app.Version = VERSION
	app.Usage = `toggl.com CLI`
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "token, t",
			Usage: "toggl.com API token",
		},
		cli.BoolFlag{
			Name:  "json, j",
			Usage: "output json instead",
		},
	}
	app.Commands = []cli.Command{
		timer.Command,
		project.Command,
	}
	return app
}
