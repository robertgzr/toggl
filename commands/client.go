package commands

import (
	"fmt"

	"github.com/dougEfresh/gtoggl-api/gthttp"
	"github.com/dougEfresh/gtoggl-api/gtproject"
	"github.com/dougEfresh/gtoggl-api/gtuser"
	"github.com/dougEfresh/gtoggl-api/gtworkspace"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func NewClient(ctx *cli.Context) (*gthttp.TogglHttpClient, error) {
	client, err := gthttp.NewClient(ctx.GlobalString("token"))
	// gthttp.AlwaysUseAPIToken(ctx.GlobalString("token")),
	// gthttp.SetTraceLogger(log.New(os.Stderr, "gthttp", log.LstdFlags)))
	if err != nil {
		return nil, errors.Wrap(err, "creating client")
	}
	return client, nil
}

func CreatedWith(ctx *cli.Context) string {
	return fmt.Sprintf("%s version %s", ctx.App.Name, ctx.App.Version)
}

func FindWorkspace(client *gthttp.TogglHttpClient, name string) (*gtworkspace.Workspace, error) {
	user := gtuser.NewClient(client)
	me, err := user.Get(true)
	if err != nil {
		return nil, errors.Wrap(err, "listing workspaces")
	}
	for _, ws := range me.Workspaces {
		if ws.Name == name {
			return &ws, nil
		}
	}
	return nil, errors.Errorf("workspace %q not found", name)
}

func FindProject(client *gthttp.TogglHttpClient, name string) (*gtproject.Project, error) {
	user := gtuser.NewClient(client)
	me, err := user.Get(true)
	if err != nil {
		return nil, errors.Wrap(err, "listing projects")
	}
	for _, pj := range me.Projects {
		if pj.Name == name {
			return &pj, nil
		}
	}
	return nil, errors.Errorf("project %q not found", name)
}
