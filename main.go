package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

func main() {
	flag.Usage = usage
	flag.Parse()
	os.Exit(run())
}

func run() int {
	action := promptui.Select{
		Label: "Select action",
		Items: listCommands(),
	}

	_, actionType, err := action.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to select action: %v\n", err)
		return 1
	}

	cmd, ok := commands[actionType]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unexpected action type\n")
		return 1
	}

	if err := cmd(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to %s: %v\n",
			strings.ToLower(actionType), err)
	}

	return 0
}
