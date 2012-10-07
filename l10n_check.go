package main

import (
	"fmt"
	"os"

	"github.com/PhiCode/l10n_check/config/properties"
	"github.com/PhiCode/l10n_check/validate"
)

type result struct {
	file  string
	props *properties.Properties
	valid *validate.Results
}

func main() {
	args := os.Args()
	l := len(args)
	if l < 1 {
		usage()
	}

	results := make([]result, 0, l)

	for _, file := range args {
		props, valid := properties.ReadAndParse(file)
		results = append(file, results, result{props, valid})
	}

	fmt.Println(results)
}

func usage() {
	fmt.Printf("usage: %s <file-name> [<file-name>]...\n", os.Args()[0])
	os.Exit(1)
}
