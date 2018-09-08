// Copyright 2018 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package literal

import (
	"fmt"
	"io"

	"github.com/hajimehoshi/goc/internal/ctype"
	"github.com/hajimehoshi/goc/internal/ioutil"
)

var escapedChars = map[byte]ctype.Int{
	'a':  '\a',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
	'v':  '\v',
	'\\': '\\',
	'\'': '\'',
	'"':  '"',
	'?':  '?',
}

func isOctDigit(c byte) bool {
	return '0' <= c && c <= '7'
}

func isHexDigit(c byte) bool {
	if '0' <= c && c <= '9' {
		return true
	}
	if 'A' <= c && c <= 'F' {
		return true
	}
	if 'a' <= c && c <= 'f' {
		return true
	}
	return false
}

func hex(c byte) byte {
	if '0' <= c && c <= '9' {
		return c - '0'
	}
	if 'A' <= c && c <= 'F' {
		return c - 'A' + 10
	}
	if 'a' <= c && c <= 'f' {
		return c - 'a' + 10
	}
	panic("not reached")
}

func ReadEscapedChar(src Source) (ctype.Int, error) {
	if err := ioutil.ShouldRead(src, '\\'); err != nil {
		return 0, err
	}

	b, err := ioutil.ShouldReadByte(src)
	if err != nil {
		return 0, err
	}

	if ch, ok := escapedChars[b]; ok {
		return ch, nil
	}

	// Hex
	if b == 'x' {
		bs, err := ioutil.ShouldPeek(src, 2)
		if err != nil {
			return 0, err
		}
		if !isHexDigit(bs[0]) {
			return 0, fmt.Errorf("literal: non-hex character in escape sequence: %q", bs[0])
		}
		if !isHexDigit(bs[1]) {
			return 0, fmt.Errorf("literal: non-hex character in escape sequence: %q", bs[1])
		}
		src.Discard(2)
		return ctype.Int((hex(bs[0]) << 4) | hex(bs[1])), nil
	}

	// Oct
	if isOctDigit(b) {
		x := ctype.Int(b - '0')

		bs, err := src.Peek(1)
		if err != nil && err != io.EOF {
			return 0, err
		}
		if len(bs) < 1 {
			return x, nil
		}
		if !isOctDigit(bs[0]) {
			return x, nil
		}
		src.Discard(1)
		x *= 8
		x += ctype.Int(bs[0] - '0')

		bs, err = src.Peek(1)
		if err != nil && err != io.EOF {
			return 0, err
		}
		if len(bs) < 1 {
			return x, nil
		}
		if !isOctDigit(bs[0]) {
			return x, nil
		}
		src.Discard(1)
		x *= 8
		x += ctype.Int(bs[0] - '0')
		if x >= 256 {
			return 0, fmt.Errorf("literal: octal escape value > 255: %d", x)
		}
		return x, nil
	}

	if b == 'u' {
		// TODO
		return 0, fmt.Errorf("literal: \\uxxxx is not implemented yet")
	}

	if b == 'U' {
		// TODO
		return 0, fmt.Errorf("literal: \\Uxxxxxxxx is not implemented yet")
	}

	return 0, fmt.Errorf("literal: unknown escape sequence: %q", b)
}

func ReadChar(src Source) (ctype.Int, error) {
	if err := ioutil.ShouldRead(src, '\''); err != nil {
		return 0, err
	}

	b, err := ioutil.ShouldPeekByte(src)
	if err != nil {
		return 0, err
	}

	v := ctype.Int(0)
	if b != '\\' {
		if b == '\r' || b == '\n' {
			return 0, fmt.Errorf("literal: newline in character literal")
		}
		if b == '\'' {
			return 0, fmt.Errorf("literal: empty character literal or unescaped ' in character literal")
		}
		src.Discard(1)
		v = ctype.Int(b)
	} else {
		b, err := ReadEscapedChar(src)
		if err != nil {
			return 0, err
		}
		v = ctype.Int(b)
	}

	if err := ioutil.ShouldRead(src, '\''); err != nil {
		return 0, err
	}

	return v, nil
}
