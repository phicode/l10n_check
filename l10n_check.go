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

	if len(results) > 1 {
		master := results[0]
		for _, other := range results[1:] {
			analyzeKeys(master, other)
			analyzeKeys(other, master)
		}
	}
}

func usage() {
	fmt.Printf("usage: %s <file-name> [<file-name>]...\n", os.Args[0])
	os.Exit(1)
}

func analyzeKeys(a, b result) {
	diff := false
	for key := range a.props.ByKey {
		if _, ok := b.props.ByKey[key]; !ok {
			if !diff {
				fmt.Printf("Key(s) in '%s' but not in '%s'\n", a.file, b.file)
				diff = true
			}
			fmt.Println("\t", key)
		}
	}
}
