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

package lex

import (
	"fmt"
	"io"
)

var escapedChars = map[byte]byte{
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

func ReadEscapedChar(src ByteReadPeeker) (byte, error) {
	if err := shouldRead(src, '\\'); err != nil {
		return 0, err
	}

	b, err := shouldReadByte(src)
	if err != nil {
		return 0, err
	}

	if ch, ok := escapedChars[b]; ok {
		return ch, nil
	}

	// Hex
	if b == 'x' {
		bs, err := shouldPeek(src, 2)
		if err != nil {
			return 0, err
		}
		if !isHexDigit(bs[0]) {
			return 0, fmt.Errorf("lex: non-hex character in escape sequence: %q", bs[0])
		}
		if !isHexDigit(bs[1]) {
			return 0, fmt.Errorf("lex: non-hex character in escape sequence: %q", bs[1])
		}
		mustDiscard(src, 2)
		return (hex(bs[0]) << 4) | hex(bs[1]), nil
	}

	// Oct
	if isOctDigit(b) {
		x := int(b - '0')

		bs, err := src.Peek(1)
		if err != nil && err != io.EOF {
			return 0, err
		}
		if len(bs) < 1 {
			return byte(x), nil
		}
		if !isOctDigit(bs[0]) {
			return byte(x), nil
		}
		mustDiscard(src, 1)
		x *= 8
		x += int(bs[0] - '0')

		bs, err = src.Peek(1)
		if err != nil && err != io.EOF {
			return 0, err
		}
		if len(bs) < 1 {
			return byte(x), nil
		}
		if !isOctDigit(bs[0]) {
			return byte(x), nil
		}
		mustDiscard(src, 1)
		x *= 8
		x += int(bs[0] - '0')
		if x >= 256 {
			return 0, fmt.Errorf("lex: octal escape value > 255: %d", x)
		}
		return byte(x), nil
	}

	if b == 'u' {
		// TODO
		return 0, fmt.Errorf("lex: \\uxxxx is not implemented yet")
	}

	if b == 'U' {
		// TODO
		return 0, fmt.Errorf("lex: \\Uxxxxxxxx is not implemented yet")
	}

	return 0, fmt.Errorf("lex: unknown escape sequence: %q", b)
}

func ReadChar(src ByteReadPeeker) (byte, error) {
	if err := shouldRead(src, '\''); err != nil {
		return 0, err
	}

	b, err := shouldPeekByte(src)
	if err != nil {
		return 0, err
	}

	v := byte(0)
	if b != '\\' {
		if b == '\r' || b == '\n' {
			return 0, fmt.Errorf("lex: newline in character literal")
		}
		if b == '\'' {
			return 0, fmt.Errorf("lex: empty character literal or unescaped ' in character literal")
		}
		mustDiscard(src, 1)
		v = b
	} else {
		b, err := ReadEscapedChar(src)
		if err != nil {
			return 0, err
		}
		v = b
	}

	if err := shouldRead(src, '\''); err != nil {
		return 0, err
	}

	return v, nil
}
