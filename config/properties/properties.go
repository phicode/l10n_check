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
	lines := bytes.Split(data, []byte("\n"))
	props.props = make([]Property, 0, len(lines)/2)

	reader := bytes.NewReader(data)
	//	line := 1
	for x, line := range lines {
		n := x + 1
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

func isWhiteSpace(b byte) bool {
	v := rune(b)
	return unicode.IsSpace(v)
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
	line = bytes.TrimSpace(line)
	return len(line) == 0 || line[0] == '#'
}
