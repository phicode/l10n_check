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

package properties

import (
	"fmt"
	"strings"
)

func Analyze(a, b *Properties, sameValues bool) int {
	n := findEmptyValues(a, b)
	n += findMissingKeys(a, b)
	n += findMissingKeys(b, a)
	if sameValues {
		n += findSameValues(a, b)
	}
	return n
}

func findMissingKeys(a, b *Properties) int {
	numFaults := 0
	for key := range a.ByKey {
		if _, ok := b.ByKey[key]; !ok {
			if numFaults == 0 {
				fmt.Printf("Key(s) in '%s' but not in '%s'\n", a.file, b.file)
			}
			fmt.Printf("\t%s\n", key)
			numFaults++
		}
	}
	if numFaults > 0 {
		fmt.Println()
	}
	return numFaults
}

func findEmptyValues(a, b *Properties) int {
	numFaults := 0
	for key, vala := range a.ByKey {
		la := len(vala.Value)
		if valb, ok := b.ByKey[key]; ok {
			lb := len(valb.Value)
			if (la == 0 && lb > 0) || (la > 0 && lb == 0) {
				if numFaults == 0 {
					fmt.Printf("Key(s) empty/non-empty in '%s' but not in '%s'\n", a.file, b.file)
				}
				fmt.Printf("\t%s\n", key)
				numFaults++
			}
		}
	}
	if numFaults > 0 {
		fmt.Println()
	}
	return numFaults
}

func findSameValues(a, b *Properties) int {
	numFaults := 0
	for key, vala := range a.ByKey {
		la := len(vala.Value)
		if valb, ok := b.ByKey[key]; ok {
			lb := len(valb.Value)
			if (la > 0 && lb > 0) && strings.EqualFold(vala.Value, valb.Value) {
				if numFaults == 0 {
					fmt.Printf("Keys(s) with same values in '%s' and '%s'\n", a.file, b.file)
				}
				fmt.Printf("\t%s = %q\n", key, vala.Value)
				numFaults++
			}
		}
	}
	if numFaults > 0 {
		fmt.Println()
	}
	return numFaults
}
