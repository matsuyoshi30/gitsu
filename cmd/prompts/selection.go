package prompts

import (
	"github.com/matsuyoshi30/gitsu/internal/fixes"

	"github.com/manifoldco/promptui"
)

// Selection runs a selection prompt and returns the index and value of the selected item
func Selection(label string, items []string) (int, string, error) {
	s := promptui.Select{
		Label:  label,
		Items:  items,
		Stdout: &fixes.BellSkipper{},
	}
	return s.Run()
}
