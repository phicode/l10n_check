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
	"io/ioutil"
	"strings"

	"github.com/PhiCode/l10n_check/validate"
)

type Property struct {
	Key   string
	Value string
	Line  int
}

type Properties struct {
	file  string
	props []*Property
	ByKey map[string]*Property
}

func ReadAndParse(filename string) (*Properties, *validate.Results, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open/read file '%s': %s", filename, err.Error())
	}
	var v validate.Results
	v.Resource = filename
	var p Properties
	p.file = filename
	parse(data, &p, &v)
	return &p, &v, nil
}

func (props *Properties) String() string {
	parts := make([]string, len(props.props))
	for i, prop := range props.props {
		parts[i] = fmt.Sprintf("line %d '%s' = '%s'", prop.Line, prop.Key, prop.Value)
	}
	return strings.Join(parts, "\n")
}

func (props *Properties) PrintAll(indent string) {
	for _, prop := range props.props {
		fmt.Printf("%sline %d '%s' = '%s'\n", indent, prop.Line, prop.Key, prop.Value)
	}
}
