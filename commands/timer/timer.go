package timer

import (
	"fmt"
	"strconv"

	gttimeentry "github.com/dougEfresh/gtoggl-api/gttimentry"
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
