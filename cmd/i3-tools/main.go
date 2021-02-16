package main

import (
	"github.com/geshido/i3-tools/pkg/selector"
	"github.com/geshido/i3-tools/pkg/usecases"
	"github.com/geshido/i3-tools/pkg/usecases/scratchpad_print"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "i3-tools"
	app.Usage = "tools for i3wm"
	app.Authors = []*cli.Author{{Name: "Yury Ignatev", Email: "geshido@gmail.com"}}

	registry := usecases.Build()
	app.Commands = []*cli.Command{
		{
			Name:  "scratchpad",
			Usage: "scratchpad actions",
			Subcommands: []*cli.Command{
				{
					Name:  "toggle",
					Usage: "toggle scratchpad visibility by its name",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "name",
							Required: true,
						},
					},
					Action: func(c *cli.Context) error {
						return registry.ScratchpadToggle.Run(c.String("name"))
					},
				},
				{
					Name:  "detect",
					Usage: "continuously print scratchpad window title",
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "name", Required: true},
						&cli.IntFlag{Name: "truncate", Value: 20},
						&cli.StringFlag{Name: "highlight", Value: "fff"},
						&cli.StringFlag{Name: "remove-suffix"},
					},
					Action: func(c *cli.Context) error {
						return registry.ScratchpadPrint.Run(scratchpad_print.Input{
							ScratchPadName: c.String("name"),
							Truncate:       c.Int("truncate"),
							ActiveColor:    c.String("highlight"),
							RemoveSuffix:   c.String("remove-suffix"),
						})
					},
				},
			},
		},
		{
			Name:  "workspace",
			Usage: "workspace actions",
			Subcommands: []*cli.Command{
				{
					Name:  "switch",
					Usage: "switch to workspace",
					Flags: []cli.Flag{
						&cli.IntFlag{Name: "idx"},
					},
					Action: func(c *cli.Context) error {
						idx := c.Int("idx")
						if idx != 0 {
							return registry.WorkspaceSwitch.Run(selector.NewByIndex(idx))
						}

						rofiSelector := selector.NewRofi("Switch to")
						return registry.WorkspaceSwitch.Run(rofiSelector)
					},
				},
				{
					Name:  "rename",
					Usage: "rename workspace",
					Action: func(_ *cli.Context) error {
						return registry.WorkspaceRename.Run(selector.NewRofi("Rename to", selector.RofiOptionNoChoices()))
					},
				},
			},
		},
		{
			Name:  "container",
			Usage: "container actions",
			Subcommands: []*cli.Command{
				{
					Name:  "wsmove",
					Usage: "move container to workspace",
					Flags: []cli.Flag{
						&cli.IntFlag{Name: "idx", Usage: "workspace index"},
					},
					Action: func(c *cli.Context) error {
						idx := c.Int("idx")
						if idx != 0 {
							return registry.ContainerMove.Run(selector.NewByIndex(idx))
						}
						return registry.ContainerMove.Run(selector.NewRofi("Move to"))
					},
				},
				{
					Name:  "alttab_listen",
					Usage: "listen to window event to track previous container id",
					Action: func(c *cli.Context) error {
						return registry.ContainerAltTab.Subscribe()
					},
				},
				{
					Name:  "alttab_switch",
					Usage: "switch to previous container",
					Action: func(c *cli.Context) error {
						return registry.ContainerAltTab.SwitchToPrevious()
					},
				},
				{
					Name:  "promote",
					Usage: "promote container to biggest one",
					Action: func(c *cli.Context) error {
						return registry.PromoteWindow.Run()
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
