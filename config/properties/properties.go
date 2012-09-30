package properties

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"unicode"
	"validate"
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

func Parse(r *io.Reader) (*Properties, *validate.Results) {
	p := new(Properties)
	validate := new(validate.Results)
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		msg := fmt.Sprintf("error while reading file content: %s", err.Error())
		validate.AddError(msg)
		return nil, validate
	}
	parse(buf, p, validate)
	return p, validate
}

func ReadAndParse(fileName string) (p *Properties, validate *validate.Results) {
	file, err := os.Open(fileName)
	if err != nil {
		v := new(validate.Results)
		msg := fmt.Sprintf("could not open file '%s': %s", fileName, err.Error())
		v.AddError(msg)
		return nil, v
	}
}

func parse(buf []byte, props *Properties, validate *validate.Results) {
	readKey := true
	var key []rune = make([]rune, 512)
	var value []rune = make([]rune, 512)
	key = key[:0]
	value = key[:0]
	line := 0
	for _, v := range buf {
		switch {
		case v == '=':
			if readKey {
				readKey = false
			} else {
				value = append(value, '=')
			}
		case v == '\n':
			if readKey {

			} else {
				prop := Property{string(key), string(value), line}
				props
			}
			line++
		case readKey && isWhiteSpace(v):
		default:
		}
	}
}

func isWhiteSpace(b byte) bool {
	v := rune(b)
	return unicode.IsSpace(v)
}
