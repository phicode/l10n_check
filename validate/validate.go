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

package validate

import (
	"bytes"
	"fmt"
)

type Results struct {
	Resource string
	warnings []Result
	errors   []Result
}

type Result struct {
	// optional line number, 0 if no line number
	// is available for this validation result
	Line int
	// the validation message
	Msg string
}

func (r *Results) AddWarning(msg string) {
	r.warnings = append(r.warnings, Result{Msg: msg})
}

func (r *Results) AddWarningN(msg string, line int) {
	r.warnings = append(r.warnings, Result{Msg: msg, Line: line})
}

func (r *Results) AddError(msg string) {
	r.errors = append(r.errors, Result{Msg: msg})
}

func (r *Results) AddErrorN(msg string, line int) {
	r.errors = append(r.errors, Result{Msg: msg, Line: line})
}

func (r *Results) Any() bool {
	return len(r.warnings) > 0 || len(r.errors) > 0
}

func (r *Results) Print(nowarn bool) int {
	var buf bytes.Buffer
	buf.WriteString("file: ")
	buf.WriteString(r.Resource)
	buf.WriteByte('\n')
	n := 0
	if l := len(r.errors); l > 0 {
		add("errors:\n", &buf, r.errors)
		n += l
	}
	if l := len(r.warnings); !nowarn && l > 0 {
		add("warnings:\n", &buf, r.warnings)
		n += l
	}
	if n == 0 {
		if nowarn {
			buf.WriteString("\tno errors\n")
		} else {
			buf.WriteString("\tno warnings or errors\n")
		}
	}
	fmt.Println(buf.String())
	return n
}

func (r *Result) String() string {
	if r.Line == 0 {
		return r.Msg
	}
	return fmt.Sprintf("Line %d: %s", r.Line, r.Msg)
}

func (r *Results) String() string {
	var buf bytes.Buffer
	buf.WriteString("file: ")
	buf.WriteString(r.Resource)
	if r.Any() {
		buf.WriteByte('\n')
		add("warnings:\n", &buf, r.warnings)
		add("errors:\n", &buf, r.errors)
	} else {
		buf.WriteString(" => no warnings or errors\n")
	}
	return buf.String()
}

func add(header string, buf *bytes.Buffer, results []Result) {
	if l := len(results); l > 0 {
		buf.WriteString(header)
		for _, res := range results {
			buf.WriteByte('\t')
			buf.WriteString(res.String())
			buf.WriteByte('\n')
		}
	}
}
