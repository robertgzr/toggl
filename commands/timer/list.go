package timer

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	gttimeentry "github.com/dougEfresh/gtoggl-api/gttimentry"
	"github.com/urfave/cli"

	"github.com/robertgzr/toggl/commands"
)

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

		if !all {
			filterRunningTimers(&tes)
		}

		if ctx.GlobalBool("json") {
			return json.NewEncoder(os.Stdout).Encode(tes)
		}

		tw := tabwriter.NewWriter(os.Stdout, 1, 8, 2, ' ', 0)
		fmt.Fprintln(tw, "ID\tDESC\tSTART\tSTOP\tDURATION\tPROJECT\tWORKSPACE")
		for _, te := range tes {
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

func filterRunningTimers(tes *gttimeentry.TimeEntries) {
	var filtered gttimeentry.TimeEntries
	for _, te := range *tes {
		if te.Duration < 0 {
			filtered = append(filtered, te)
		}
	}
	*tes = filtered
}
