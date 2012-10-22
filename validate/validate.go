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
	r.warnings = resAppend(r.warnings, Result{Msg: msg})
}

func (r *Results) AddWarningN(msg string, line int) {
	r.warnings = resAppend(r.warnings, Result{Msg: msg, Line: line})
}

func (r *Results) AddError(msg string) {
	r.errors = resAppend(r.errors, Result{Msg: msg})
}

func (r *Results) AddErrorN(msg string, line int) {
	r.errors = resAppend(r.errors, Result{Msg: msg, Line: line})
}

func (r *Results) Any() bool {
	return len(r.warnings) > 0 || len(r.errors) > 0
}

func resAppend(orig []Result, res Result) []Result {
	l, c := len(orig), cap(orig)
	if l >= c {
		xs := make([]Result, (c+1)*2)
		copy(xs, orig)
		orig = xs[:l]
	}
	return append(orig, res)
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
