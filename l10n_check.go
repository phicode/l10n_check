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

	anyFault := true

	for _, result := range results {
		fmt.Println(result.valid)
		anyFault = anyFault || result.valid.Any()
	}

	if len(results) > 1 {
		master := results[0]
		for _, other := range results[1:] {
			n := analyzeKeys(master, other)
			n += analyzeKeys(other, master)
			anyFault = anyFault || n > 0
		}
	}

	if anyFault {
		os.Exit(1)
	}
	os.Exit(0)
}

func usage() {
	fmt.Printf("usage: %s [-v] <file-name> [<file-name>]...\n", os.Args[0])
	os.Exit(1)
}

func analyzeKeys(a, b result) int {
	numFaults := 0
	diff := false
	for key := range a.props.ByKey {
		if _, ok := b.props.ByKey[key]; !ok {
			if !diff {
				fmt.Printf("Key(s) in '%s' but not in '%s'\n", a.file, b.file)
				diff = true
			}
			fmt.Println("\t", key)
			numFaults++
		}
	}
	diff = false
	for key, vala := range a.props.ByKey {
		la := len(vala.Value)
		if valb, ok := b.props.ByKey[key]; ok {
			if lb := len(valb.Value); (la == 0 && lb > 0) || (la > 0 && lb == 0) {
				if !diff {
					fmt.Printf("Key(s) empty/non-empty in '%s' but not in '%s'\n", a.file, b.file)
					diff = true
				}
				fmt.Println("\t", key)
				numFaults++
			}
		}
	}
	return numFaults
}

func (r *result) String() string {
	return fmt.Sprintf("file: %s\nprops:\n%s\nvalid:\n%s\n", r.file, r.props, r.valid)
}
