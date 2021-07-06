package cmd

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/matsuyoshi30/gitsu/cmd/prompts"
	"github.com/matsuyoshi30/gitsu/internal/config"

	"github.com/urfave/cli/v2"
)

const resetOutputTemplate = `The following %d user profile(s) will be deleted 
{{ range . }}  {{ . }}
{{ end }}`

func ResetCommand() *cli.Command {
	return &cli.Command{
		Name:    "reset",
		Aliases: []string{"r"},
		Usage:   "Remove all saved user profiles",
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

			rawTemplate := fmt.Sprintf(resetOutputTemplate, len(list))
			t, err := template.New("reset").Parse(rawTemplate)
			if err != nil {
				return err
			}

			b := &bytes.Buffer{}
			err = t.Execute(b, list)
			if err != nil {
				return err
			}

			fmt.Println(b)

			selection, _, err := prompts.SelectionCustom(
				"Delete above profiles?",
				[]string{"Yes", "No"},
			)
			if err != nil {
				return err
			}

			if selection == 1 {
				return nil
			}

			cfg.Reset()
			return config.Write(cfg)
		},
	}
}
