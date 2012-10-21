package properties

import (
	//"bytes"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	//"unicode"

	"github.com/PhiCode/l10n_check/validate"
)

type Property struct {
	Key   string
	Value string
	Line  int
}

type Properties struct {
	props []*Property
	byKey map[string]*Property
}

type context struct {
	key      []byte
	val      []byte
	props    *Properties
	validate *validate.Results
	lineNr   int
}

func ReadAndParse(filename string) (*Properties, *validate.Results) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		msg := fmt.Sprintf("could not open/read file '%s': %s", filename, err.Error())
		v := new(validate.Results)
		v.AddError(msg)
		return nil, v
	}
	var v validate.Results
	var p Properties
	parse(data, &p, &v)
	return &p, &v
}

func parse(data []byte, props *Properties, validate *validate.Results) {
	lines := splitLines(data)
	props.props = make([]*Property, 0, len(lines)/2)
	props.byKey = make(map[string]*Property)
	partialLine := false

	ctx := context{
		key:      make([]byte, 0, 4096),
		val:      make([]byte, 0, 4096),
		props:    props,
		validate: validate,
	}

	for x, line := range lines {
		if !partialLine {
			if isEmptyOrComment(line) {
				continue
			}
			ctx.lineNr = x + 1
			partialLine = ctx.readStart(line)
		} else {
			partialLine = ctx.readContinue(line)
		}
		if !partialLine {
			ctx.finishKeyValue()
		}
	}
	ctx.finishKeyValue()
}

func (ctx *context) appendKey(b byte) { ctx.key = append(ctx.key, b) }
func (ctx *context) appendVal(b byte) { ctx.val = append(ctx.val, b) }
func (ctx *context) unreadVal() {
	if l := len(ctx.val); l > 0 {
		ctx.val = ctx.val[:l-1]
	}
}

func (ctx *context) readStart(line []byte) bool {
	// 1. consume whitespace
	// 2. consume key
	// 3. consume whitespace and : and =
	// 4. consume value
	// return true if last char is \ => partial line
	state := 1
	var prev byte
	for _, v := range line {
		switch state {
		case 1:
			if !isWhiteSpace(v) {
				ctx.appendKey(v)
				state = 2
			}
		case 2:
			if isWhiteSpace(v) {
				state = 3
			} else {
				if (v == ':' || v == '=') && prev != '\\' {
					state = 3
				} else {
					ctx.appendKey(v)
				}
			}
		case 3:
			if !isWhiteSpace(v) && v != ':' && v != '=' {
				ctx.appendVal(v)
				state = 4
			}
		case 4:
			ctx.appendVal(v)
		}
		prev = v
	}
	return ctx.finishLine(prev)
}

func (ctx *context) readContinue(line []byte) bool {
	// 1. consume whitespace
	// 2. consume value
	// return true if last char is \ => partial line
	state := 1
	var prev byte
	for _, v := range line {
		switch state {
		case 1:
			if !isWhiteSpace(v) {
				ctx.appendVal(v)
				state = 2
			}
		case 2:
			ctx.val = append(ctx.val, v)
		}
		prev = v
	}
	return ctx.finishLine(prev)
}

func (ctx *context) finishLine(prev byte) bool {
	// TODO: handle empty value
	if prev == '\\' {
		ctx.unreadVal()
		return true
	}
	return false
}

func (ctx *context) finishKeyValue() {
	fmt.Printf("line=%d, key='%s', value='%s'\n", ctx.lineNr, ctx.key, ctx.val)

	line := ctx.lineNr
	key := string(ctx.key)
	val := string(ctx.val)
	p := &Property{key, val, line}
	ctx.props.props = append(ctx.props.props, p)
	old, contains := ctx.props.byKey[key]
	if contains {
		msg := fmt.Sprintf("duplicate key '%s' from line %d overwrite previous key-value pair from line %d", key, line, old.Line)
		ctx.validate.AddWarningN(msg, line)
	}
	ctx.props.byKey[key] = p

	// reset read-buffers
	ctx.key = ctx.key[:0]
	ctx.val = ctx.val[:0]
}

// TODO: make "lines" a container.List
func splitLines(data []byte) [][]byte {
	var lines [][]byte = make([][]byte, 0, 256)
	var line []byte = make([]byte, 0, 4096)
	var prev byte
	for _, v := range data {
		if v == '\r' || v == '\n' {
			if prev == '\r' && v == '\n' {
				prev = v
				continue
			}
			l := make([]byte, len(line))
			copy(l, line)
			lines = append(lines, l)
			line = line[:0] // empty
		} else {
			line = append(line, v)
		}
		prev = v
	}
	return lines
}

// sorted byte slice
// 0x09 = tab
// 0x0A = LF
// 0x0C = form feed
// 0x0D = CR
// 0x20 = space
var whitespaces = []byte{0x09, 0x0A, 0x0C, 0x0D, 0x20}

func isWhiteSpace(b byte) bool {
	n := len(whitespaces)
	i := sort.Search(n, func(i int) bool { return whitespaces[i] >= b })
	//fmt.Printf("isWhitespace(%s): %b\n", b, (i < n && whitespaces[i] == b))
	return i < n && whitespaces[i] == b
}

// empty / comment lines 
// are those whos first non-whitespace character is # or !
func isEmptyOrComment(line []byte) bool {
	if len(line) == 0 {
		return true
	}
	for _, b := range line {
		if !isWhiteSpace(b) {
			return b == '#' || b == '!'
		}
	}
	return true
}

func (props *Properties) String() string {
	parts := make([]string, len(props.props))
	for i, prop := range props.props {
		parts[i] = fmt.Sprintf("line %d '%s' = '%s'", prop.Line, prop.Key, prop.Value)
	}
	return strings.Join(parts, "\n")
}
