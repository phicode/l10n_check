package properties

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"unicode"

	"github.com/PhiCode/l10n_check/validate"
)

type Property struct {
	Key   string
	Value string
	Line  int
}

type Properties struct {
	props []Property
	byKey map[string]*Property
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
	props.props = make([]Property, 0, len(lines)/2)
	lines := splitLines(data)
	partialLine := false
	for x, line := range lines {
		if !partialLine {
			if isEmptyOrComment(line) {
				continue
			}
			key, value, partialLine := readKeyValue(line, false)
		} else {
			value, partialLine := 
		}

		line := x + 1
		if readEmptyOrComment(reader, n, validate) {
			line++
			continue
		}

		key, ok := readKey(line, n, validate)
		if !ok {
			continue
		}
		val, ok := readVal(line, n, validate)
		if ok {
			prop := Property{Key: key, Value: val, Line: n}
			props.props = append(props.props, prop)
			props.byKey[key] = prop
		}
	}
}

func readEmptyOrComment(reader *bytes.Reader, n int, validate *validate.Results) bool {
	num := reader.Len()
	if num <= 0 {
		return false
	}
	v1 := reader.ReadByte()
	if v1 == '\r' {
		if num > 1 {
			v2 := reader.ReadByte()
			if v2 == '\n' {
				return true
			} else {
				reader.UnreadByte()
				reader.UnreadByte()
				validate.AddErrorN("\\r is not followed by a \\n", n)
			}
		} else {
			reader.UnreadByte()
			validate.AddErrorN("\\r at end of file", n)
		}

	}
	reader.ReadAt(b, off)
	line = bytes.TrimSpace(line)
	return len(line) == 0 || line[0] == '#'
}

// TODO: make "lines" a container/List
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
}

const whitespaces = []byte{' ', '\r', '\n', '\t', '\f'}

func isWhiteSpace(b byte) bool {
	return bytes.Contains(whitespaces, b)
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
