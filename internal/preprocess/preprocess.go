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

type bufPPTokenReader struct {
	tokens  []*Token
	pos     int
	current *Token
}

func (t *bufPPTokenReader) NextPPToken() (*Token, error) {
	if t.pos >= len(t.tokens) {
		tk := &Token{
			Type: EOF,
		}
		t.current = tk
		return tk, nil
	}
	tk := t.tokens[t.pos]
	t.pos++
	t.current = tk
	return tk, nil
}

func (t *bufPPTokenReader) Peek() (*Token, error) {
	if t.pos >= len(t.tokens) {
		return &Token{
			Type: EOF,
		}, nil
	}
	return t.tokens[t.pos], nil
}

func (t *bufPPTokenReader) AtLineHead() bool {
	if t.current == nil {
		return true
	}
	if t.current.Type == '\n' {
		return true
	}
	return false
}

type macro struct {
	name      string
	tokens    []*Token
	paramsLen int
}

type preprocessor struct {
	src          *bufPPTokenReader
	tokens       map[string][]*Token
	sub          []*Token
	visited      map[string]struct{}
	macros       map[string]macro
	expandedFrom map[*Token][]string
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

func (p *preprocessor) applyMacro(src *bufPPTokenReader, m *macro) ([]*Token, map[int]struct{}, error) {
	// Apply object-like macro.
	if m.paramsLen == -1 {
		return m.tokens, nil, nil
	}

	// Apply function-like macro.
	// Parse arguments
	if _, err := nextExpected(src, '('); err != nil {
		return nil, nil, err
	}

	args := [][]*Token{}
	t, err := src.Peek()
	if err != nil {
		return nil, nil, err
	}
	if t.Type == ')' {
		if _, err := nextExpected(src, ')'); err != nil {
			panic("not reached")
		}
	} else {
	args:
		for {
			arg := []*Token{}
			level := 0
			for {
				t, err := src.NextPPToken()
				if err != nil {
					return nil, nil, err
				}
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
		return nil, nil, fmt.Errorf("preprocess: expected %d args but %d", m.paramsLen, len(args))
	}

	wasParam := map[int]struct{}{}
	r := []*Token{}
	for _, t := range m.tokens {
		if t.Type != Param {
			r = append(r, t)
			continue
		}
		for i := range args[t.ParamIndex] {
			wasParam[len(r)+i] = struct{}{}
		}
		r = append(r, args[t.ParamIndex]...)
	}
	return r, wasParam, nil
}

func (p *preprocessor) next() (*Token, error) {
	if len(p.sub) > 0 {
		t := p.sub[0]
		p.sub = p.sub[1:]

		var e []string
		if p.expandedFrom != nil {
			e = p.expandedFrom[t]
		}

		if t.Type != Identifier {
			return t, nil
		}

		m, ok := p.macros[t.Val]
		if !ok {
			return t, nil
		}

		// The token came from the same macro.
		if p.expandedFrom != nil {
			for _, name := range p.expandedFrom[t] {
				if m.name == name {
					return t, nil
				}
			}
		}

		// "6.10.3.4 Rescanning and further replacement" [spec]
		src := &bufPPTokenReader{
			tokens: p.sub,
			pos:    0,
		}
		tks, wasParam, err := p.applyMacro(src, &m)
		if err != nil {
			return nil, err
		}
		p.sub = append(tks, p.sub[src.pos:]...)

		for i, t := range tks {
			if _, ok := wasParam[i]; ok {
				continue
			}
			// TODO: duplicated?
			p.expandedFrom[t] = append(p.expandedFrom[t], e...)
			p.expandedFrom[t] = append(p.expandedFrom[t], m.name)
		}

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
		tks, wasParam, err := p.applyMacro(p.src, &m)
		if err != nil {
			return nil, err
		}
		p.sub = tks
		p.expandedFrom = map[*Token][]string{}
		for i, t := range tks {
			if _, ok := wasParam[i]; ok {
				continue
			}
			p.expandedFrom[t] = []string{m.name}
		}
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
			t, err := p.src.NextPPToken()
			if err != nil {
				return nil, err
			}
			if t.Type != Identifier {
				return nil, fmt.Errorf("preprocess: expected ident or keyword but %s", t.Type)
			}
			name := t.Val
			// TODO: What if the same macro is redefined?

			paramsLen := -1
			var params []string
			t, err = p.src.Peek()
			if err != nil {
				return nil, err
			}
			if t.Type == '(' && t.Adjacent {
				if _, err := nextExpected(p.src, '('); err != nil {
					panic("not reached")
				}
				params = []string{}
				t, err := p.src.Peek()
				if err != nil {
					return nil, err
				}
				if t.Type == ')' {
					if _, err := nextExpected(p.src, ')'); err != nil {
						panic("not reached")
					}
				} else {
					for {
						t, err := p.src.NextPPToken()
						if err != nil {
							return nil, err
						}
						if t.Type != Identifier {
							return nil, fmt.Errorf("preprocess: expected ident or keyword but %s", t.Type)
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
			if len(params) > 0 {
				for i, t := range ts {
					if t.Type != Identifier {
						continue
					}
					idx := -1
					for i, p := range params {
						if t.Val == p {
							idx = i
							break
						}
					}
					if idx == -1 {
						continue
					}
					ts[i] = &Token{
						Type:       Param,
						ParamIndex: idx,
					}
				}
			}
			p.macros[name] = macro{
				name:      name,
				tokens:    ts,
				paramsLen: paramsLen,
			}
		case "undef":
			t, err := p.src.NextPPToken()
			if err != nil {
				return nil, err
			}
			if t.Type != Identifier {
				return nil, fmt.Errorf("preprocess: expected ident or keyword but %s", t.Type)
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

func Preprocess(path string, tokens map[string][]*Token) ([]*Token, error) {
	return preprocessImpl(path, tokens, map[string]struct{}{})
}

func preprocessImpl(path string, tokens map[string][]*Token, visited map[string]struct{}) ([]*Token, error) {
	ts, ok := tokens[path]
	if !ok {
		return nil, fmt.Errorf("preprocess: file not found: %s", path)
	}
	p := &preprocessor{
		src: &bufPPTokenReader{
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
