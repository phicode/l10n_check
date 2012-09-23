package properties

import (
	"io"
	"validate"

//	"os"
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

func Parse(r *io.Reader) (p *Properties, validate *validate.Results) {
	p = new(Properties)
	validate = new(validate.Results)
	var buf [1024]byte
	left := 0
	line := 0
	for {
		n, err := r.Read(buf[left:])
		if err != nil {
			left = left + n
			consumed, line := parse(buf[:left], line, p, validate)
			left = left - consumed
			if left > 0 {
				copy(buf, buf[consumed:consumed+left])
			}
		}
	}
	return
}

func ReadAndParse(file string) (p *Properties, validate *validate.Results) {

}

func parse(buf []byte, p *Properteis, validate *validate.Results) (consumed int, line int) {

}
