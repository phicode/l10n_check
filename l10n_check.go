package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()

	args := flag.Args()
	l := len(args)
	if l < 2 {
		usage()
	}
	encoding := args[0]
	files := args[1:]
	result := check.RunCheck(encoding, files)

	fmt.Println(result)
}

func usage() {
	fmt.Println("usage: <encoding> <file-name> [<file-name>]...")
	os.Exit(1)
}
