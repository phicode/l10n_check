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

package properties

import (
	"bytes"
	"container/list"
	"fmt"
	"strings"
	"unicode"

	"github.com/phicode/l10n_check/validate"
)

type context struct {
	// current state
	lineNr int
	key    []byte
	val    []byte

	// overall state
	props    *Properties
	validate *validate.Results
}

func parse(data []byte, props *Properties, validate *validate.Results) {
	lines := splitLines(data)
	props.props = make([]*Property, 0, lines.Len()/2)
	props.ByKey = make(map[string]*Property)

	ctx := context{
		key:      make([]byte, 0, 4096),
		val:      make([]byte, 0, 4096),
		props:    props,
		validate: validate,
	}

	var res parseResult
	for nr, iter := 1, lines.Front(); iter != nil; nr, iter = nr+1, iter.Next() {
		line, ok := iter.Value.([]byte)
		if !ok {
			panic("internal error: not a byte-slice")
		}
		if res != PARTIAL_LINE {
			if isEmptyOrComment(line) {
				continue
			}
			ctx.lineNr = nr
			res = ctx.readStart(line)
		} else {
			res = ctx.readContinue(line)
		}
		if res == KEY_VALUE {
			ctx.finishKeyValue()
		} else if res == ONLY_KEY {
			msg := fmt.Sprintf("line contains only a key: '%s'", string(ctx.key))
			validate.AddErrorN(msg, nr)
			ctx.reset()
		}
	}
	if !ctx.isEmpty() {
		ctx.finishKeyValue()
	}
}

func (ctx *context) appendKey(b byte) { ctx.key = append(ctx.key, b) }
func (ctx *context) appendVal(b byte) { ctx.val = append(ctx.val, b) }
func (ctx *context) unreadVal() {
	if l := len(ctx.val); l > 0 {
		ctx.val = ctx.val[:l-1]
	}
}

type parseResult int

const (
	KEY_VALUE    parseResult = iota // finished key-value pair
	PARTIAL_LINE                    // key and partial value, which continues on the next line
	ONLY_KEY                        // line contained only a key
)

func (ctx *context) readStart(line []byte) parseResult {
	// 1. consume whitespace
	// 2. consume key
	// 3. consume whitespace, :, =
	// 4. consume value
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
	if state != 4 {
		return ONLY_KEY
	}
	return ctx.finishLine(prev)
}

func (ctx *context) readContinue(line []byte) parseResult {
	// 1. consume whitespace
	// 2. consume value
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

func (ctx *context) finishLine(prev byte) parseResult {
	if prev == '\\' {
		ctx.unreadVal()
		return PARTIAL_LINE
	}
	return KEY_VALUE
}

func (ctx *context) isEmpty() bool {
	return len(ctx.key) == 0 && len(ctx.val) == 0
}

func (ctx *context) finishKeyValue() {
	line := ctx.lineNr
	key, _ := ctx.sliceToStr(ctx.key)
	val, endEscaped := ctx.sliceToStr(ctx.val)
	// due to the layout of property files there can only be unescaped whitespace at the end of a value
	if !endEscaped && hasTrailingWhitespace(val) {
		msg := fmt.Sprintf("value of key '%s' contains unescaped trailing spaces", key)
		ctx.validate.AddWarningN(msg, line)
	}

	p := &Property{key, val, line}
	ctx.props.props = append(ctx.props.props, p)
	old, contains := ctx.props.ByKey[key]
	if contains {
		msg := fmt.Sprintf("key '%s' overwrites line %d", key, old.Line)
		ctx.validate.AddWarningN(msg, line)
	}
	ctx.props.ByKey[key] = p
	ctx.reset()
}

func (ctx *context) reset() {
	// reset read-buffers
	ctx.key = ctx.key[:0]
	ctx.val = ctx.val[:0]
}

// returns the string and a bool which tells us if the last character was escaped
func (ctx *context) sliceToStr(xs []byte) (string, bool) {
	l := len(xs)
	if l == 0 {
		return "", false
	}
	// states
	// 1. reading regular characters
	// 2. reading char after \
	// 3. reading unicode value (\uxxxx)
	// 4. skip n chars, switch to 1 afterwards
	state := 1
	var buf bytes.Buffer
	skip := 0
	var escaped bool
	for idx, x := range xs {
		switch state {
		case 1:
			if x == '\\' {
				escaped = true
				state = 2
			} else {
				escaped = false
				ctx.addRune(&buf, x, idx)
			}
		case 2:
			switch x {
			case 't':
				buf.WriteRune('\t')
				state = 1
			case 'n':
				buf.WriteRune('\n')
				state = 1
			case 'r':
				buf.WriteRune('\r')
				state = 1
			case 'f':
				buf.WriteRune('\f')
				state = 1
			case 'u': // unicode sequence
				state = 3
			default:
				ctx.addRune(&buf, x, idx)
				state = 1
			}
		case 3:
			// idx: 012345
			// val: \uffff
			// pos:   ^ => rem = len - idx
			// =>  rem = 6 - 2 = 4
			remaining := l - idx
			if remaining < 4 {
				msg := fmt.Sprintf("unicode sequence start found (\\u) but there are too few remaining bytes in the value")
				ctx.validate.AddErrorN(msg, ctx.lineNr)
			} else {
				unicodeSeq := xs[idx:(idx + 4)]
				ctx.parseUnicodeSeq(unicodeSeq, &buf)
			}
			// skip the next 3 chars since we already read them
			skip = 3
			state = 4
		case 4:
			skip--
			if skip == 0 {
				state = 1
			}
		}
	}
	return buf.String(), escaped
}

func (ctx *context) addRune(buf *bytes.Buffer, x byte, idx int) {
	r := rune(x)
	if unicode.IsSpace(r) || unicode.IsGraphic(r) {
		buf.WriteRune(r)
	} else {
		msg := fmt.Sprintf("non-graphic character found, code: %d, index in value: %d", int(x), idx)
		ctx.validate.AddErrorN(msg, ctx.lineNr)
	}
}

func (ctx *context) parseUnicodeSeq(xs []byte, buf *bytes.Buffer) {
	var symbol uint32
	for _, x := range xs {
		if v, ok := fromHexChar(x); ok {
			symbol = symbol*16 + v
		} else {
			msg := fmt.Sprintf("invalid unicode sequence: %s", string(xs))
			ctx.validate.AddErrorN(msg, ctx.lineNr)
			return
		}
	}
	// fmt.Printf("unicode char: %x\n", symbol)
	// TODO: validate symbol
	buf.WriteRune(rune(symbol))
}

func fromHexChar(x byte) (hex uint32, ok bool) {
	if x >= '0' && x <= '9' {
		return uint32(x - '0'), true
	}
	if x >= 'a' && x <= 'f' {
		return uint32(x-'a') + 10, true
	}
	if x >= 'A' && x <= 'F' {
		return uint32(x-'A') + 10, true
	}
	return 0, false
}

func splitLines(data []byte) *list.List {
	var lines *list.List = list.New()
	var line []byte = make([]byte, 0, 4096)
	var prev byte
	for _, v := range data {
		if v == '\r' || v == '\n' {
			if prev == '\r' && v == '\n' {
				prev = v
				continue
			}
			pushLine(lines, line)
			line = line[:0] // empty
		} else {
			line = append(line, v)
		}
		prev = v
	}
	if len(line) > 0 {
		pushLine(lines, line)
	}
	return lines
}

func pushLine(lines *list.List, line []byte) {
	l := make([]byte, len(line))
	copy(l, line)
	lines.PushBack(l)
}

const (
	WS_TAB byte = 0x09 // tab
	WS_LF       = 0x0A // line feed
	WS_FF       = 0x0C // form feed
	WS_CR       = 0x0D // carriage return
	WS_SP       = 0x20 // space
)

func isWhiteSpace(b byte) bool {
	return b == WS_TAB || b == WS_LF || b == WS_FF || b == WS_CR || b == WS_SP
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
	// all whitespace line
	return true
}

func hasTrailingWhitespace(s string) bool {
	return len(s) > len(strings.TrimRight(s, " \t\r\n"))
}
