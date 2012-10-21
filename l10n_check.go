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
	if len(os.Args) < 2 {
		usage()
	}
	args := os.Args[1:]
	results := make([]result, 0, len(args))
	for _, file := range args {
		props, valid := properties.ReadAndParse(file)
		results = append(results, result{file, props, valid})
		fmt.Println("props:", props, "valid:", valid)
	}

	fmt.Println(results)
}

func usage() {
	fmt.Printf("usage: %s <file-name> [<file-name>]...\n", os.Args[0])
	os.Exit(1)
}
