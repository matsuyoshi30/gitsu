package prompts

import (
	"github.com/matsuyoshi30/gitsu/internal/fixes"

	"github.com/manifoldco/promptui"
)

func Input(label string) (string, error) {
	i := &promptui.Prompt{
		Label:  label,
		Stdout: &fixes.BellSkipper{},
	}
	return i.Run()
}

func InputWithValidation(label string, v promptui.ValidateFunc) (string, error) {
	i := &promptui.Prompt{
		Label:    label,
		Validate: v,
		Stdout:   &fixes.BellSkipper{},
	}
	return i.Run()
}
