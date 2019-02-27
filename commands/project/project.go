package project

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/dougEfresh/gtoggl-api/gtuser"
	"github.com/urfave/cli"

	"github.com/robertgzr/toggl/commands"
)

var Command = cli.Command{
	Name:    "project",
	Aliases: []string{"projects"},
	Usage:   "manipulate projects",
	Subcommands: cli.Commands{
		listCommand,
	},
}

var listCommand = cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "list projects",
	Action: func(ctx *cli.Context) error {
		tc, err := commands.NewClient(ctx)
		if err != nil {
			return err
		}
		me, err := gtuser.NewClient(tc).Get(true)
		if err != nil {
			return err
		}

		if ctx.GlobalBool("json") {
			return json.NewEncoder(os.Stdout).Encode(me.Projects)
		}

		tw := tabwriter.NewWriter(os.Stdout, 1, 8, 2, ' ', 0)
		fmt.Fprintln(tw, "ID\tNAME\tCLIENT\tWORKSPACE\tPRIVATE")
		for _, pj := range me.Projects {
			fmt.Fprintf(tw, "%v\t%v\t%v\t%v\t%t\n", pj.Id, pj.Name, pj.CId, pj.WId, *pj.IsPrivate)
		}
		return tw.Flush()
	},
}
