package main

import (
	"flag"
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

var verbose *bool = flag.Bool("v", false, "verbose")

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		usage()
	}
	args := flag.Args()
	results := make([]result, 0, len(args))
	for _, file := range args {
		props, valid, err := properties.ReadAndParse(file)
		if props == nil {
			fmt.Println(err)
			continue
		}
		r := result{file, props, valid}
		results = append(results, r)
		//fmt.Printf(r.String())
	}

	for _, result := range results {
		fmt.Printf("%s: %d keys\n", result.file, len(result.props.ByKey))
		if *verbose {
			result.props.PrintAll("\t")
		}
	}

	fmt.Println()

	for _, result := range results {
		fmt.Println(result.valid)
	}

	//fmt.Println(results)

	if len(results) > 1 {
		master := results[0]
		for _, other := range results[1:] {
			analyzeKeys(master, other)
			analyzeKeys(other, master)
		}
	}
}

func usage() {
	fmt.Printf("usage: %s [-v] <file-name> [<file-name>]...\n", os.Args[0])
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

func (r *result) String() string {
	return fmt.Sprintf("file: %s\nprops:\n%s\nvalid:\n%s\n", r.file, r.props, r.valid)
}
