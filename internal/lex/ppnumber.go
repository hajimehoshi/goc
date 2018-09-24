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

func ReadPPNumber(src gio.Source) (string, error) {
	b, err := gio.ShouldReadByte(src)
	if err != nil {
		return "", err
	}

	r := []byte{}

	switch {
	case IsDigit(b):
		r = append(r, b)
	case b == '.':
		r = append(r, b)
		b, err := gio.ShouldReadByte(src)
		if err != nil {
			return "", err
		}
		if !IsDigit(b) {
			return "", fmt.Errorf("lex: expected digit but %q", string(rune(b)))
		}
		r = append(r, b)
	default:
		return "", fmt.Errorf("lex: expected digit or . but %q", string(rune(b)))
	}

	for {
		bs, err := src.Peek(1)
		if err != nil && err != io.EOF {
			return "", err
		}
		if len(bs) == 0 {
			if err != io.EOF {
				panic("not reached")
			}
			break
		}

		b := bs[0]
		if !IsDigit(b) && b != '.' && !IsNondigit(b) {
			break
		}
		src.Discard(1)
		r = append(r, b)

		if b != 'e' && b != 'E' && b != 'p' && b != 'P' {
			continue
		}

		bs, err = src.Peek(1)
		if err != nil && err != io.EOF {
			return "", err
		}
		if len(bs) == 0 {
			if err != io.EOF {
				panic("not reached")
			}
			break
		}
		b = bs[0]

		if b == '+' || b == '-' {
			src.Discard(1)
			r = append(r, b)
			continue
		}
	}

	return string(r), nil
}
