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
	r.Warnings = resAppend(r.Warnings, Result{line: -1, Msg: msg})
}

func (r *Results) AddWarning(line int, msg string) {
	r.Warnings = resAppend(r.Warnings, Result{line: -1, Msg: msg})
}

func (r *Results) AddError(msg string) {
	r.Errors = resAppend(r.Errors, Result{line: line, Msg: msg})
}

func (r *Results) AddError(line int, msg string) {
	r.Errors = resAppend(r.Errors, Result{line: line, Msg: msg})
}

func (r *Results) Any() bool {
	return len(r.Warnings) > 0 || len(r.Errors) > 0
}

func resAppend(orig []Result, res Result) []Result {
	l := len(orig)
	if l < cap(orig) {
		// returns a new slice with the same backing array, len+1
		return append(orig, res)
	} else {
		if l == 0 {
			return []Result{res}
		} else {
			xs := make([]Result, l*2)
			copy(xs, orig)
			xs[l] = res
			return xs[0 : l+1]
		}
	}
}
