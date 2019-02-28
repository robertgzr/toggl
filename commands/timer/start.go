package timer

import (
	"fmt"

	"github.com/dougEfresh/gtoggl-api/gtproject"
	gttimeentry "github.com/dougEfresh/gtoggl-api/gttimentry"
	"github.com/dougEfresh/gtoggl-api/gtworkspace"
	"github.com/urfave/cli"

	"github.com/robertgzr/toggl/commands"
)

var startCommand = cli.Command{
	Name:  "start",
	Usage: "start a timer",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name: "description, d",
		},
		cli.StringFlag{
			Name: "project, p",
		},
		cli.StringFlag{
			Name: "workspace, w",
		},
		cli.StringSliceFlag{
			Name: "tags, t",
		},
	},
	Action: func(ctx *cli.Context) error {
		var (
			description = ctx.String("description")
			project     = ctx.String("project")
			workspace   = ctx.String("workspace")
			tags        = ctx.StringSlice("tags")
		)

		tc, err := commands.NewClient(ctx)
		if err != nil {
			return err
		}

		if project == "" && workspace == "" {
			return cli.NewExitError("workspace or project is required", 1)
		}

		timeEntry := gttimeentry.TimeEntry{
			Description: description,
			Tags:        tags,
			CreatedWith: commands.CreatedWith(ctx),
		}

		var ws *gtworkspace.Workspace
		if workspace != "" {
			var err error
			ws, err = commands.FindWorkspace(tc, workspace)
			if err != nil {
				return err
			}
			timeEntry.Wid = ws.Id
		}

		var pj *gtproject.Project
		if project != "" {
			var err error
			pj, err = commands.FindProject(tc, project)
			if err != nil {
				return err
			}
			timeEntry.Wid = pj.WId
			timeEntry.Pid = pj.Id
		}

		timeentries := gttimeentry.NewClient(tc)
		te, err := timeentries.Start(&timeEntry)
		if err != nil {
			return err
		}
		fmt.Printf("Started timer (id %d)\n", te.Id)
		return nil
	},
}
