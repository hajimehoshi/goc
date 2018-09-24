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

	gio "github.com/hajimehoshi/goc/internal/io"
)

// IsWhitespace returns true if c is a whitespace char, otherwise false.
func IsWhitespace(c byte) bool {
	switch c {
	case ' ', '\t', '\v', '\f', '\r', '\n':
		return true
	default:
		return false
	}
}

// IsNondigit returns true if c is nondigit character, otherwise false.
// "6.4.2.1 General" [spec]
func IsNondigit(c byte) bool {
	if 'A' <= c && c <= 'Z' {
		return true
	}
	if 'a' <= c && c <= 'z' {
		return true
	}
	if c == '_' {
		return true
	}
	return false
}

// IsDigit returns true if c is digit character, otherwise false.
// "6.4.2.1 General" [spec]
func IsDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// IsSingleCharPunctuator returns true if c is a single character punctuator, otherwise false.
// "6.4.6 Punctuators" [spec]
func IsSingleCharPunctuator(c byte) bool {
	switch c {
	case '[', ']', '(', ')', '{', '}', '.', '&', '*', '+', '-', '~', '!', '/', '%', '<', '>', '^', '|', '?', ':', ';', '=', ',', '#':
		return true
	default:
		return false
	}
}

func ReadIdentifier(src gio.Source) (string, error) {
	b, err := gio.ShouldReadByte(src)
	if err != nil {
		return "", err
	}

	if !IsNondigit(b) {
		return "", fmt.Errorf("lex: expected nondigit but %q", string(rune(b)))
	}

	r := []byte{b}
	for {
		bs, err := src.Peek(1)
		if err != nil && err != io.EOF {
			return "", err
		}
		if len(bs) < 1 {
			break
		}
		if !IsDigit(bs[0]) && !IsNondigit(bs[0]) {
			break
		}
		src.Discard(1)
		r = append(r, bs[0])
	}

	return string(r), nil
}
