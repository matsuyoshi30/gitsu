package cmd

import (
	"github.com/matsuyoshi30/gitsu/cmd/prompts"
	"github.com/matsuyoshi30/gitsu/internal/config"
	"github.com/matsuyoshi30/gitsu/internal/models"

	"github.com/urfave/cli/v2"
)

// AddCommand returns the definition for the 'gitsu add' command
func AddCommand() *cli.Command {
	return &cli.Command{
		Name:    "add",
		Aliases: []string{"a"},
		Usage:   "Add new user",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "gpg",
				Value: false,
				Usage: "Add GPG key ID",
			},
		},
		Action: func(c *cli.Context) error {
			name, err := prompts.Input("Git user name")
			if err != nil {
				return err
			}

			email, err := prompts.InputWithValidation(
				"Git user email",
				func(s string) error {
					return models.ValidateEmail(s, false)
				},
			)
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
			cfg, err := config.Read()
			if err != nil {
				return err
			}

			err = cfg.AddUser(user)
			if err != nil {
				return err
			}

			return config.Write(cfg)
		},
	}
}
