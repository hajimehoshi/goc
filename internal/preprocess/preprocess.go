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
)

type ppTokenBufReader struct {
	tokens  PPTokenReader
	buf     []*Token
	current *Token
}

func (t *ppTokenBufReader) NextPPToken() (*Token, error) {
	if len(t.buf) > 0 {
		tk := t.buf[0]
		t.buf = t.buf[1:]
		t.current = tk
		return tk, nil
	}
	tk, err := t.tokens.NextPPToken()
	if err != nil {
		return nil, err
	}
	t.current = tk
	return tk, nil
}

func (t *ppTokenBufReader) PeekPPToken() (*Token, error) {
	if len(t.buf) > 0 {
		return t.buf[0], nil
	}
	tk, err := t.tokens.NextPPToken()
	if err != nil {
		return nil, err
	}
	t.buf = append(t.buf, tk)
	return tk, nil
}

func (t *ppTokenBufReader) AtLineHead() bool {
	if t.current == nil {
		return true
	}
	if t.current.Type == '\n' {
		return true
	}
	return false
}

type ppTokenSliceReader struct {
	tokens []*Token
	pos    int
}

func (t *ppTokenSliceReader) NextPPToken() (*Token, error) {
	if len(t.tokens) <= t.pos {
		return &Token{
			Type: EOF,
		}, nil
	}
	tk := t.tokens[t.pos]
	t.pos++
	return tk, nil
}

func (t *ppTokenSliceReader) PeekPPToken() (*Token, error) {
	if len(t.tokens) <= t.pos {
		return &Token{
			Type: EOF,
		}, nil
	}
	return t.tokens[t.pos], nil
}

type preprocessor struct {
	src     *ppTokenBufReader
	tokens  map[string]PPTokenReader
	sub     []*Token
	visited map[string]struct{}
	macros  map[string]macro
}

func (p *preprocessor) NextPPToken() (*Token, error) {
	for {
		t, err := p.next()
		if t == nil && err == nil {
			continue
		}
		return t, err
	}
}

func (p *preprocessor) next() (*Token, error) {
	if len(p.sub) > 0 {
		t := p.sub[0]
		p.sub = p.sub[1:]

		if t.Type != Identifier {
			return t, nil
		}

		m, ok := p.macros[t.Val]
		if !ok {
			return t, nil
		}

		// The token came from the same macro.
		if t.ExpandedFrom != nil {
			if _, ok := t.ExpandedFrom[m.name]; ok {
				return t, nil
			}
		}

		// "6.10.3.4 Rescanning and further replacement" [spec]
		src := &ppTokenSliceReader{
			tokens: p.sub,
			pos:    0,
		}
		tks, err := m.apply(src, t.ExpandedFrom)
		if err != nil {
			return nil, err
		}
		p.sub = append(tks, p.sub[src.pos:]...)
		return nil, nil
	}

	wasLineHead := p.src.AtLineHead()

	t, err := p.src.NextPPToken()
	if err != nil {
		return nil, err
	}

	switch t.Type {
	case Identifier:
		m, ok := p.macros[t.Val]
		if !ok {
			return t, nil
		}
		tks, err := m.apply(p.src, nil)
		if err != nil {
			return nil, err
		}
		p.sub = tks
	case '#':
		if !wasLineHead {
			return t, nil
		}
		// The tokens must end with '\n', so nil check is not needed.
		t, err := p.src.NextPPToken()
		if err != nil {
			return nil, err
		}
		if t.Type == '\n' {
			// Empty directive
			return t, nil
		}
		if t.Type != Identifier {
			return nil, fmt.Errorf("preprocess: expected %s but %s", Identifier, t.Type)
		}
		switch t.Val {
		case "define":
			t, err := nextExpected(p.src, Identifier)
			if err != nil {
				return nil, err
			}
			name := t.Val
			// TODO: What if the same macro is redefined?

			paramsLen := -1
			var params []string
			t, err = p.src.PeekPPToken()
			if err != nil {
				return nil, err
			}
			if t.Type == '(' && t.Adjacent {
				if _, err := nextExpected(p.src, '('); err != nil {
					panic("not reached")
				}
				params = []string{}
				t, err := p.src.PeekPPToken()
				if err != nil {
					return nil, err
				}
				if t.Type == ')' {
					if _, err := nextExpected(p.src, ')'); err != nil {
						panic("not reached")
					}
				} else {
					for {
						t, err := nextExpected(p.src, Identifier)
						if err != nil {
							return nil, err
						}
						params = append(params, t.Val)
						t, err = nextExpected(p.src, ')', ',')
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

			ts := []*Token{}
			for {
				t, err := p.src.NextPPToken()
				if err != nil {
					return nil, err
				}
				if t.Type == '\n' {
					break
				}
				ts = append(ts, t)
			}

			// Replace parameter identifier-like tokens with Param tokens.
			if paramsLen >= 0 {
				ts2 := []*Token{}
				for i := 0; i < len(ts); i++ {
					t := ts[i]
					switch t.Type {
					case '#':
						i++
						if i >= len(ts) {
							return nil, fmt.Errorf("preprocess: '#' is not followed by a macro parameter")
						}
						t := ts[i]
						if t.Type != Identifier {
							return nil, fmt.Errorf("preprocess: '#' is not followed by a macro parameter")
						}
						idx := -1
						for i, p := range params {
							if t.Val == p {
								idx = i
								break
							}
						}
						if idx == -1 {
							return nil, fmt.Errorf("preprocess: '#' is not followed by a macro parameter")
						}
						ts2 = append(ts2, &Token{
							Type:       Param,
							ParamIndex: idx,
							ParamHash:  true,
						})
					case Identifier:
						idx := -1
						for i, p := range params {
							if t.Val == p {
								idx = i
								break
							}
						}
						if idx != -1 {
							ts2 = append(ts2, &Token{
								Type:       Param,
								ParamIndex: idx,
							})
						} else {
							ts2 = append(ts2, t)
						}
					default:
						ts2 = append(ts2, t)
					}
				}
				ts = ts2
			}

			p.macros[name] = macro{
				name:      name,
				tokens:    ts,
				paramsLen: paramsLen,
			}
		case "undef":
			t, err := nextExpected(p.src, Identifier)
			if err != nil {
				return nil, err
			}
			delete(p.macros, t.Val)
			if _, err := nextExpected(p.src, '\n'); err != nil {
				return nil, err
			}
		case "include":
			t, err := nextExpected(p.src, HeaderName)
			if err != nil {
				return nil, err
			}
			path := t.Val
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
				t, err := p.src.NextPPToken()
				if err != nil {
					return nil, err
				}
				if t.Type == '\n' {
					break
				}
				// TODO: Define RawString() and use it?
				msg += " " + t.String()
			}
			return nil, fmt.Errorf("preprocess: #error" + msg)
		default:
			return nil, fmt.Errorf("preprocess: invalid preprocessing directive %s", t.Val)
		}
	default:
		return t, nil
	}

	// Preprocessing derective is processed correctly.
	// There is no token to return.
	return nil, nil
}

func Preprocess(path string, tokens map[string]PPTokenReader) ([]*Token, error) {
	return preprocessImpl(path, tokens, map[string]struct{}{
		path: {},
	})
}

func preprocessImpl(path string, tokens map[string]PPTokenReader, visited map[string]struct{}) ([]*Token, error) {
	ts, ok := tokens[path]
	if !ok {
		return nil, fmt.Errorf("preprocess: file not found: %s", path)
	}
	p := &preprocessor{
		src: &ppTokenBufReader{
			tokens: ts,
		},
		tokens:  tokens,
		visited: visited,
		macros:  map[string]macro{},
	}
	r := []*Token{}
	for {
		t, err := p.NextPPToken()
		if err != nil {
			return nil, err
		}
		if t == nil {
			continue
		}
		if t.Type == '\n' {
			continue
		}
		if t.Type == EOF {
			break
		}
		r = append(r, t)
	}
	return r, nil
}
