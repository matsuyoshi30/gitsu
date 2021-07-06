package cmd

import (
	"fmt"

	"github.com/matsuyoshi30/gitsu/cmd/prompts"
	"github.com/matsuyoshi30/gitsu/internal/config"
	"github.com/matsuyoshi30/gitsu/internal/models"

	"github.com/urfave/cli/v2"
)

func ModifyCommand() *cli.Command {
	return &cli.Command{
		Name:    "modify",
		Aliases: []string{"m"},
		Usage:   "Modify existing user",
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

			index, _, err := prompts.SelectionCustom("Select git user", list)
			if err != nil {
				return err
			}

			name, err := prompts.Input("New git user name, leave empty for no change")
			if err != nil {
				return err
			}

			email, err := prompts.Input("New git email address, leave empty for no change")
			if err != nil {
				return err
			}

			var keyID string
			if c.Bool("gpg") {
				keyID, err = prompts.Input("GPG key ID")
				if err != nil {
					return err
				}
			}

			alias, err := prompts.Input("User alias, leave empty for no alias")
			if err != nil {
				return err
			}

			user := models.NewUser(name, email, alias, keyID)
			err = cfg.ModifyUser(index, user)
			if err != nil {
				return err
			}

			return config.Write(cfg)
		},
	}
}
