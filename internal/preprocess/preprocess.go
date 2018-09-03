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

package preprocess

import (
	"fmt"

	"github.com/hajimehoshi/goc/internal/token"
)

type tokenReader struct {
	tokens   []*token.Token
	pos      int
	linehead bool
}

func (t *tokenReader) Next() *token.Token {
	if t.pos >= len(t.tokens) {
		return nil
	}
	tk := t.tokens[t.pos]
	t.pos++
	return tk
}

func (t *tokenReader) Peek() *token.Token {
	if t.pos >= len(t.tokens) {
		return nil
	}
	tk := t.tokens[t.pos]
	t.pos++
	return tk
}

func (t *tokenReader) AtLineHead() bool {
	if t.pos == 0 {
		return true
	}
	if len(t.tokens) == 0 {
		return true
	}
	if t.tokens[t.pos-1].Type == '\n' {
		return true
	}
	return false
}

func Preprocess(tokens []*token.Token) ([]*token.Token, error) {
	src := &tokenReader{
		tokens: tokens,
	}
	ts := []*token.Token{}

	for {
		t := src.Next()
		if t == nil {
			break
		}

		switch t.Type {
		case token.Ident:
			// TODO: Apply macros
			ts = append(ts, t)
		case '#':
			if !src.AtLineHead() {
				ts = append(ts, t)
				continue
			}
			t = src.Next()
			if t.Type == '\n' {
				// Empty
				ts = append(ts, t)
				continue
			}
			if t.Type != token.Ident {
				return nil, fmt.Errorf("preprocess: expected %s but %s", token.Ident, t.Type)
			}
			switch t.Name {
			case "define":
				return nil, fmt.Errorf("preprocess: #define is not implemented")
			case "undef":
				return nil, fmt.Errorf("preprocess: #undef is not implemented")
			case "include":
				return nil, fmt.Errorf("preprocess: #include is not implemented")
			case "if":
				return nil, fmt.Errorf("preprocess: #if is not implemented")
			case "ifdef":
				return nil, fmt.Errorf("preprocess: #ifdef is not implemented")
			case "ifndef":
				return nil, fmt.Errorf("preprocess: #ifndef is not implemented")
			case "else":
				return nil, fmt.Errorf("preprocess: #else is not implemented")
			case "endif":
				return nil, fmt.Errorf("preprocess: #line is not implemented")
			case "line":
				return nil, fmt.Errorf("preprocess: #line is not implemented")
			case "elif":
				return nil, fmt.Errorf("preprocess: #elif is not implemented")
			case "pragma":
				return nil, fmt.Errorf("preprocess: #pragma is not implemented")
			case "error":
				msg := ""
				for {
					t := src.Next()
					if t == nil || t.Type == '\n' {
						break
					}
					// TODO: Define RawString() and use it?
					msg += " " + t.String()
				}
				return nil, fmt.Errorf("preprocess: #error" + msg)
			default:
				return nil, fmt.Errorf("preprocess: invalid preprocessing directive %s", t.Name)
			}
		default:
			ts = append(ts, t)
		}
	}
	return ts, nil
}
