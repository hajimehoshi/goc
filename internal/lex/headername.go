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
)

func ReadHeaderName(src Source) (string, error) {
	b, err := shouldReadByte(src)
	if err != nil {
		return "", err
	}
	if b != '<' && b != '"' {
		return "", fmt.Errorf("lex: '<' or '\"' expected but %q", string(b))
	}

	bs := []byte{}
	end := byte(0)
	switch b {
	case '<':
		end = '>'
	case '"':
		end = '"'
	default:
		panic("not reached")
	}
	for {
		b, err := shouldReadByte(src)
		if err != nil {
			return "", err
		}
		switch b {
		case end:
			return string(bs), nil
		case '\n':
			return "", fmt.Errorf("lex: unterminated header-name")
		default:
			bs = append(bs, b)
		}
	}
}
