package main

import (
	"fmt"
	"os"

	"github.com/matsuyoshi30/gitsu/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
