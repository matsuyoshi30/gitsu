package cmd

import (
	"fmt"

	"github.com/matsuyoshi30/gitsu/internal/config"
	"github.com/matsuyoshi30/gitsu/internal/git"
	"github.com/matsuyoshi30/gitsu/internal/models"

	"github.com/urfave/cli/v2"
)

func InitCommand() *cli.Command {
	return &cli.Command{
		Name:    "init",
		Aliases: []string{"i"},
		Usage:   "Initialize user config by providing an alias",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "global",
				Value: false,
				Usage: "Set git user globally",
			},
		},
		Action: func(c *cli.Context) error {
			alias := c.Args().First()

			cfg, err := config.Read()
			if err != nil {
				return err
			}

			list := cfg.UserList()
			if len(list) == 0 {
				fmt.Println("No users")
				return nil
			}

			var scope = models.Local
			if c.Bool("global") {
				scope = models.Global
			}

			err = git.IsInsideWorktree(scope)
			if err != nil {
				return err
			}

			if alias == "" {
				defaultUser, err := cfg.SelectDefaultUser()
				if err != nil {
					return err
				}

				err = git.SetConfig(defaultUser, scope)
				if err != nil {
					return err
				}

				fmt.Printf("Setting default profile %s", defaultUser.Format(0))
				return nil
			}

			user, err := cfg.SelectUserByAlias(alias)
			if err != nil {
				return err
			}

			err = git.SetConfig(user, scope)
			if err != nil {
				return err
			}

			fmt.Printf("Setting profile %s", user.Format(0))
			return nil
		},
	}
}
