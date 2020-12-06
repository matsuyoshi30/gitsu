package main

import (
	"flag"
	"fmt"
	"os"
)

func usage() {
	format := `Usage:
  gitsu [flags]

Flags:
  --global              Set user as global.

Author:
  matsuyoshi30 <sfbgwm30@gmail.com>
`
	fmt.Fprintln(os.Stderr, format)
}

var (
	isGlobal = flag.Bool("global", false, "Set user as global")

	// these are set in build step
	version = "unversioned"
	commit  = "?"
	date    = "?"
)
