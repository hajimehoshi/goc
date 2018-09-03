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

type FileTokenizer interface {
	TokenizeFile(path string) []*token.Token
}

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

type preprocessor struct {
	src           *tokenReader
	fileTokenizer FileTokenizer
}

func (p *preprocessor) Next() (*token.Token, error) {
	t := p.src.Next()
	if t == nil {
		return nil, nil
	}

	switch t.Type {
	case token.Ident:
		// TODO: Apply macros
		return t, nil
	case '#':
		if !p.src.AtLineHead() {
			return t, nil
		}
		t = p.src.Next()
		if t == nil || t.Type == '\n' {
			// Empty directive
			return t, nil
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
			t = p.src.Next()
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
				t := p.src.Next()
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
		return t, nil
	}
}

func Preprocess(tokens []*token.Token, fileTokenizer FileTokenizer) ([]*token.Token, error) {
	p := &preprocessor{
		src: &tokenReader{
			tokens: tokens,
		},
		fileTokenizer: fileTokenizer,
	}
	ts := []*token.Token{}
	for {
		t, err := p.Next()
		if err != nil {
			return nil, err
		}
		if t == nil {
			break
		}
		ts = append(ts, t)
	}
	return ts, nil
}
