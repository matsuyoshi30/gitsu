package prompts

import (
	"github.com/matsuyoshi30/gitsu/internal/fixes"

	"github.com/manifoldco/promptui"
)

var CustomTemplate = &promptui.SelectTemplates{
	Label:    "{{ . }}",
	Active:   "▶ {{ . | cyan }}",
	Inactive: "  {{ . }}",
	Selected: "▶ {{ . | green }}",
}

// Selection runs a selection prompt and returns the index and value of the selected item
func Selection(label string, items []string) (int, string, error) {
	s := promptui.Select{
		Label:  label,
		Items:  items,
		Stdout: &fixes.BellSkipper{},
	}
	return s.Run()
}

// SelectionCustom runs a selection prompt with custom template and returns the index and value of the selected item
func SelectionCustom(label string, items []string) (int, string, error) {
	s := promptui.Select{
		Label:     label,
		Items:     items,
		Stdout:    &fixes.BellSkipper{},
		Templates: CustomTemplate,
	}
	return s.Run()
}
