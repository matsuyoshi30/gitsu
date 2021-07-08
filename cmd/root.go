package cmd

import (
	"os"

	"github.com/matsuyoshi30/gitsu/cmd/prompts"

	"github.com/urfave/cli/v2"
)

func Execute() error {
	app := &cli.App{
		Name:  "gitsu",
		Usage: "Easily switch between multiple git users",
		Commands: []*cli.Command{
			DeleteCommand(),
			ModifyCommand(),
			SelectCommand(),
			ResetCommand(),
			InitCommand(),
			AddCommand(),
		},
		Action: func(c *cli.Context) error {
			action, _, err := prompts.SelectionCustom(
				"Select action",
				[]string{
					"Select user",
					"Add new user",
					"Delete user",
					"Modify user",
				},
			)
			if err != nil {
				return err
			}

			switch action {
			case 0:
				return c.App.Command("select").Run(c)
			case 1:
				return c.App.Command("add").Run(c)
			case 2:
				return c.App.Command("delete").Run(c)
			case 3:
				return c.App.Command("modify").Run(c)
			}

			return nil
		},
	}
	return app.Run(os.Args)
}
