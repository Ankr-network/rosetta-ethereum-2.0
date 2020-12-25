package main

import (
	"os"

	"rosetta-ethereum-2.0/cmd"

	"github.com/fatih/color"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}
}
