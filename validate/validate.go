package validate

type Results struct {
	Resource string
	Warnings []Result
	Errors   []Result
}

type Result struct {
	// optional line number, 0 if no line number
	// is available for this validation result
	Line int
	// the validation message
	Msg string
}

func (r *Results) AddWarning(msg string) {
	r.Warnings = resAppend(r.Warnings, Result{Msg: msg})
}

func (r *Results) AddWarningN(msg string, line int) {
	r.Warnings = resAppend(r.Warnings, Result{Msg: msg, Line: line})
}

func (r *Results) AddError(msg string) {
	r.Errors = resAppend(r.Errors, Result{Msg: msg})
}

func (r *Results) AddErrorN(msg string, line int) {
	r.Errors = resAppend(r.Errors, Result{Msg: msg, Line: line})
}

func (r *Results) Any() bool {
	return len(r.Warnings) > 0 || len(r.Errors) > 0
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
