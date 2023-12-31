package main

import (
	"fmt"
	"os"

	"github.com/jodydadescott/home-dns-server/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
