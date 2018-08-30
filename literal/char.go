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
	"bufio"
	"fmt"
	"io"
)

var escapedChars = map[byte]rune{
	'a':  '\a',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
	'v':  '\v',
	'e':  '\033',
	'E':  '\033',
	'\\': '\\',
}

func isOctalDigit(c byte) bool {
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

func ReadChar(src io.Reader) (rune, error) {
	buf := bufio.NewReader(src)
	b, err := buf.ReadByte()
	if err != nil {
		if err == io.EOF {
			return 0, fmt.Errorf("literal: unexpected EOF")
		}
		return 0, err
	}
	if b != '\\' {
		return rune(b), nil
	}

	b2, err := buf.ReadByte()
	if err != nil {
		return 0, err
	}

	if ch, ok := escapedChars[b2]; ok {
		return ch, nil
	}

	// Hex
	if b2 == 'x' {
		bs, err := buf.ReadBytes(2)
		if err != nil && err != io.EOF {
			return 0, err
		}
		if len(bs) < 2 {
			return 0, fmt.Errorf("literal: unexpected EOF")
		}
		if !isHexDigit(bs[0]) {
			return 0, fmt.Errorf("literal: non-hex character in escape sequence: %q", bs[0])
		}
		if !isHexDigit(bs[1]) {
			return 0, fmt.Errorf("literal: non-hex character in escape sequence: %q", bs[1])
		}
		return rune((hex(bs[0]) << 4) | hex(bs[1])), nil
	}

	// Oct
	if isOctalDigit(b2) {
		// TODO: Implement this
	}

	// TODO: UCS-2 and UCS-4

	return 0, fmt.Errorf("literal: unknown escape sequence")
}
