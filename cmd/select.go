package cmd

import (
	"fmt"

	"github.com/matsuyoshi30/gitsu/cmd/prompts"
	"github.com/matsuyoshi30/gitsu/internal/config"
	"github.com/matsuyoshi30/gitsu/internal/git"
	"github.com/matsuyoshi30/gitsu/internal/models"

	"github.com/urfave/cli/v2"
)

func SelectCommand() *cli.Command {
	return &cli.Command{
		Name:    "select",
		Aliases: []string{"s"},
		Usage:   "Select existing user",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "global",
				Value: false,
				Usage: "Set git user globally",
			},
		},
		Action: func(c *cli.Context) error {
			cfg, err := config.Read()
			if err != nil {
				return err
			}

			list := cfg.UserList()
			if len(list) == 0 {
				fmt.Println("No users")
				return nil
			}

			index, _, err := prompts.Selection("Select git user", list)
			if err != nil {
				return err
			}

			user, err := cfg.SelectUser(index)
			if err != nil {
				return err
			}

			var scope = models.Local
			if c.Bool("global") {
				scope = models.Global
			}

			return git.SetConfig(user, scope)
		},
	}
}
