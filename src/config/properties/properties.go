package properties

import (
	"io"
	"os"
)

type Property struct {
	Key   string
	Value string
	Line  int
}

type Properties struct {
	props  []Property
	byKey  map[string]*Property
	errors []string
}

func Parse(r *io.Reader) (p *Properties, err error) {

}
