/*
 * Copyright (c) 2006-2011 Philipp Meinen <philipp@bind.ch>
 * 
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software
 * is furnished to do so, subject to the following conditions:
 * 
 * The above copyright notice and this permission notice shall be included
 * in all copies or substantial portions of the Software.
 * 
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
 * IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
 * TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH
 * THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

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

var (
	verbose *bool = flag.Bool("v", false, "enable verbose mode")
	nowarn  *bool = flag.Bool("nowarn", false, "do not print warnings")
)

func main() {
	flag.Usage = func() {
		fmt.Printf("usage: %s [-v] <file-name> [<file-name> ...]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
	}

	anyFault := false
	args := flag.Args()
	results := make([]result, 0, len(args))
	for _, file := range args {
		props, valid, err := properties.ReadAndParse(file)
		if err != nil {
			fmt.Println(err)
			anyFault = true
			continue
		}
		r := result{file, props, valid}
		results = append(results, r)
	}

	for _, result := range results {
		fmt.Printf("%s: %d keys\n", result.file, len(result.props.ByKey))
		if *verbose {
			result.props.PrintAll("\t")
		}
	}

	fmt.Println()

	for _, result := range results {
		n := result.valid.Print(*nowarn)
		anyFault = anyFault || n > 0
	}

	if len(results) > 1 {
		master := results[0]
		for _, other := range results[1:] {
			n := properties.AnalyzeAll(master.props, other.props)
			anyFault = anyFault || n > 0
		}
	}

	if anyFault {
		os.Exit(2)
	}
}

func (r *result) String() string {
	return fmt.Sprintf("file: %s\nprops:\n%s\nvalid:\n%s\n", r.file, r.props, r.valid)
}
