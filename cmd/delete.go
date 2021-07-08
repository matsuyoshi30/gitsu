package cmd

import (
	"fmt"

	"github.com/matsuyoshi30/gitsu/cmd/prompts"
	"github.com/matsuyoshi30/gitsu/internal/config"

	"github.com/urfave/cli/v2"
)

func DeleteCommand() *cli.Command {
	return &cli.Command{
		Name:    "delete",
		Aliases: []string{"d"},
		Usage:   "Delete existing user",
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

			err = cfg.DeleteUser(index)
			if err != nil {
				return err
			}

			return config.Write(cfg)
		},
	}
}
