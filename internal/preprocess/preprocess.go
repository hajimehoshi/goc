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
	"io"
	"strings"

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

func (t *tokenReader) NextExpected(expected ...token.Type) (*token.Token, error) {
	tk := t.Next()
	if tk == nil {
		return nil, fmt.Errorf("preprocess: unexpected EOF")
	}
	for _, e := range expected {
		if tk.Type == e {
			return tk, nil
		}
	}

	s := []string{}
	for _, e := range expected {
		s = append(s, fmt.Sprintf("%s", e))
	}
	return nil, fmt.Errorf("preprocess: expected %s but %s", strings.Join(s, ","), tk.Type)
}

func (t *tokenReader) Peek() *token.Token {
	if t.pos >= len(t.tokens) {
		return nil
	}
	return t.tokens[t.pos]
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

type macro struct {
	tokens    []*token.Token
	paramsLen int
}

type preprocessor struct {
	src     *tokenReader
	tokens  map[string][]*token.Token
	sub     []*token.Token
	visited map[string]struct{}
	macros  map[string]macro
}

func (p *preprocessor) Next() (*token.Token, error) {
	for {
		t, err := p.next()
		if t == nil && err == nil {
			continue
		}
		return t, err
	}
}

func (p *preprocessor) next() (*token.Token, error) {
	if len(p.sub) > 0 {
		t := p.sub[0]
		p.sub = p.sub[1:]
		return t, nil
	}

	wasLineHead := p.src.AtLineHead()

	t := p.src.Next()
	if t == nil {
		return nil, io.EOF
	}

	switch t.Type {
	case token.Ident:
		m, ok := p.macros[t.Name]
		if !ok {
			return t, nil
		}

		// Apply object-like macro.
		if m.paramsLen == -1 {
			p.sub = m.tokens
			return nil, nil
		}

		// Apply function-like macro.
		// Parse arguments
		if _, err := p.src.NextExpected('('); err != nil {
			return nil, err
		}
		args := [][]*token.Token{}
		if p.src.Peek().Type == ')' {
			if _, err := p.src.NextExpected(')'); err != nil {
				panic("not reached")
			}
		} else {
		args:
			for {
				arg := []*token.Token{}
				level := 0
				for {
					t := p.src.Next()
					if t.Type == ')' && level == 0 {
						args = append(args, arg)
						break args
					}
					if t.Type == ',' && level == 0 {
						args = append(args, arg)
						break
					}
					arg = append(arg, t)
					if t.Type == '(' {
						level++
					}
					if t.Type == ')' {
						level--
					}
				}
			}
		}
		if len(args) != m.paramsLen {
			return nil, fmt.Errorf("preprocess: expected %d args but %d", m.paramsLen, len(args))
		}
		p.sub = []*token.Token{}
		for _, t := range m.tokens {
			if t.Type != token.Param {
				p.sub = append(p.sub, t)
				continue
			}
			p.sub = append(p.sub, args[t.ParamIndex]...)
		}
	case '#':
		if !wasLineHead {
			return t, nil
		}
		// The tokens must end with '\n', so nil check is not needed.
		t = p.src.Next()
		if t.Type == '\n' {
			// Empty directive
			return t, nil
		}
		if t.Type != token.Ident {
			return nil, fmt.Errorf("preprocess: expected %s but %s", token.Ident, t.Type)
		}
		switch t.Name {
		case "define":
			t, err := p.src.NextExpected(token.Ident)
			if err != nil {
				return nil, err
			}
			// TODO: What if the same macro is redefined?

			paramsLen := -1
			var params []string
			if t := p.src.Peek(); t.Type == '(' && t.Adjacent {
				if _, err := p.src.NextExpected('('); err != nil {
					panic("not reached")
				}
				params = []string{}
				if p.src.Peek().Type == ')' {
					if _, err := p.src.NextExpected(')'); err != nil {
						panic("not reached")
					}
				} else {
					for {
						t, err := p.src.NextExpected(token.Ident)
						if err != nil {
							return nil, err
						}
						params = append(params, t.Name)
						t, err = p.src.NextExpected(')', ',')
						if err != nil {
							return nil, err
						}
						if t.Type == ')' {
							break
						}
					}
				}
				paramsLen = len(params)
			}

			ts := []*token.Token{}
			for {
				t := p.src.Next()
				if t.Type == '\n' {
					break
				}
				ts = append(ts, t)
			}

			// Replace parameter identifier tokens with Param tokens.
			if len(params) > 0 {
				for i, t := range ts {
					if t.Type != token.Ident {
						continue
					}
					idx := -1
					for i, p := range params {
						if t.Name == p {
							idx = i
							break
						}
					}
					if idx == -1 {
						continue
					}
					ts[i] = &token.Token{
						Type:       token.Param,
						ParamIndex: idx,
					}
				}
			}
			p.macros[t.Name] = macro{
				tokens:    ts,
				paramsLen: paramsLen,
			}
		case "undef":
			return nil, fmt.Errorf("preprocess: #undef is not implemented")
		case "include":
			t, err := p.src.NextExpected(token.HeaderName)
			if err != nil {
				return nil, err
			}
			path := t.StringValue
			if _, ok := p.visited[path]; ok {
				return nil, fmt.Errorf("preprocess: recursive #include: %s", path)
			}
			p.visited[path] = struct{}{}
			ts, err := preprocessImpl(path, p.tokens, p.visited)
			if err != nil {
				return nil, err
			}
			p.sub = ts
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
				if t.Type == '\n' {
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

	// Preprocessing derective is processed correctly.
	// There is no token to return.
	return nil, nil
}

func Preprocess(path string, tokens map[string][]*token.Token) ([]*token.Token, error) {
	return preprocessImpl(path, tokens, map[string]struct{}{})
}

func preprocessImpl(path string, tokens map[string][]*token.Token, visited map[string]struct{}) ([]*token.Token, error) {
	ts, ok := tokens[path]
	if !ok {
		return nil, fmt.Errorf("preprocess: file not found: %s", path)
	}
	p := &preprocessor{
		src: &tokenReader{
			tokens: ts,
		},
		tokens:  tokens,
		visited: visited,
		macros:  map[string]macro{},
	}
	r := []*token.Token{}
	for {
		t, err := p.Next()
		if t != nil {
			r = append(r, t)
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}
	return r, nil
}
