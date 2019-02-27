package timer

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/dougEfresh/gtoggl-api/gtproject"
	gttimeentry "github.com/dougEfresh/gtoggl-api/gttimentry"
	"github.com/dougEfresh/gtoggl-api/gtworkspace"
	"github.com/urfave/cli"

	"github.com/robertgzr/toggl/commands"
)

var Command = cli.Command{
	Name:    "timer",
	Aliases: []string{"timers"},
	Usage:   "control timers",
	Subcommands: cli.Commands{
		listCommand,
		startCommand,
		stopCommand,
		removeCommand,
	},
}

var listCommand = cli.Command{
	Name:        "list",
	Aliases:     []string{"ls"},
	Usage:       "list timers",
	Description: `list timers, running timers are marked with * on their duration fields`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "print all",
		},
	},
	Action: func(ctx *cli.Context) error {
		var (
			all = ctx.Bool("all")
		)

		tc, err := commands.NewClient(ctx)
		if err != nil {
			return err
		}
		timeentries := gttimeentry.NewClient(tc)
		tes, err := timeentries.List()
		if err != nil {
			return err
		}

		if ctx.GlobalBool("json") {
			return json.NewEncoder(os.Stdout).Encode(tes)
		}

		tw := tabwriter.NewWriter(os.Stdout, 1, 8, 2, ' ', 0)
		fmt.Fprintln(tw, "ID\tDESC\tSTART\tSTOP\tDURATION\tPROJECT\tWORKSPACE")
		for _, te := range tes {
			if !all && te.Duration > 0 {
				// skip non-running timers
				continue
			}
			var duration string
			if te.Duration < 0 {
				// calculate duration of a running timer as:
				// current_unix_time + elapsed_time (as negative offset from unix epoch)
				d := time.Duration(time.Now().UnixNano()) + time.Duration(te.Duration)*time.Second
				// mark a running timer
				d = d.Round(1 * time.Second)
				duration = "*" + d.String()
			} else {
				duration = (time.Duration(te.Duration) * time.Second).String()
			}
			fmt.Fprintf(tw, "%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
				te.Id,
				te.Description,
				commands.FormatTime(te.Start),
				commands.FormatTime(te.Stop),
				duration,
				te.Pid,
				te.Wid)
		}
		return tw.Flush()
	},
}

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

var stopCommand = cli.Command{
	Name:        "stop",
	Usage:       "stop a timer",
	ArgsUsage:   "<id>",
	Description: `stops the specified timer, use "-" for the currently running timer`,
	Action: func(ctx *cli.Context) error {
		idArg := ctx.Args().First()
		if idArg == "" {
			return cli.NewExitError("please specify a timer id", 1)
		}
		tc, err := commands.NewClient(ctx)
		if err != nil {
			return err
		}
		timeentries := gttimeentry.NewClient(tc)

		var t *gttimeentry.TimeEntry
		if idArg == "-" {
			// get timer from timers/current endpoint
			var err error
			t, err = timeentries.Current()
			if err != nil {
				return err
			}

		} else {
			// use a timer id, that's really all we need
			timerId, err := strconv.ParseInt(idArg, 10, 64)
			if err != nil {
				return err
			}
			t = &gttimeentry.TimeEntry{Id: uint64(timerId)}
		}

		te, err := timeentries.Stop(t)
		if err != nil {
			return err
		}
		fmt.Printf("Stopped timer (id %d)\n", te.Id)
		return nil
	},
}

var removeCommand = cli.Command{
	Name:      "rm",
	Usage:     "delete a timer",
	ArgsUsage: "<id>",
	Action: func(ctx *cli.Context) error {
		idArg := ctx.Args().First()
		if idArg == "" {
			return cli.NewExitError("please specify a timer id", 1)
		}
		timerId, err := strconv.ParseInt(idArg, 10, 64)
		if err != nil {
			return err
		}

		tc, err := commands.NewClient(ctx)
		if err != nil {
			return err
		}
		timeentries := gttimeentry.NewClient(tc)
		if err := timeentries.Delete(uint64(timerId)); err != nil {
			return err
		}
		fmt.Printf("Removed timer (id %d)\n", timerId)
		return nil
	},
}
