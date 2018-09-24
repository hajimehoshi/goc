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
	"strings"

	"github.com/hajimehoshi/goc/internal/io"
	"github.com/hajimehoshi/goc/internal/lex"
)

type macro struct {
	name      string
	tokens    []*Token
	paramsLen int
}

func (m *macro) apply(src ppTokenReadPeeker, expandedFrom map[string]struct{}) ([]*Token, error) {
	// Apply object-like macro.
	if m.paramsLen == -1 {
		for _, t := range m.tokens {
			if t.ExpandedFrom == nil {
				t.ExpandedFrom = map[string]struct{}{}
			}
			for n := range expandedFrom {
				t.ExpandedFrom[n] = struct{}{}
			}
			t.ExpandedFrom[m.name] = struct{}{}
		}
		return m.tokens, nil
	}

	// Apply function-like macro.
	// Parse arguments
	if _, err := nextExpected(src, '('); err != nil {
		return nil, err
	}

	args := [][]*Token{}
	t, err := src.peekPPToken()
	if err != nil {
		return nil, err
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
					return nil, err
				}
				if t.Type == EOF {
					return nil, fmt.Errorf("preprocess: unexpected EOF")
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
		return nil, fmt.Errorf("preprocess: expected %d args but %d", m.paramsLen, len(args))
	}

	r := []*Token{}
	for _, t := range m.tokens {
		if t.Type != Param {
			if t.ExpandedFrom == nil {
				t.ExpandedFrom = map[string]struct{}{}
			}
			for n := range expandedFrom {
				t.ExpandedFrom[n] = struct{}{}
			}
			t.ExpandedFrom[m.name] = struct{}{}
			r = append(r, t)
			continue
		}
		if !t.ParamHash {
			r = append(r, args[t.ParamIndex]...)
			continue
		}

		// "6.10.3.2 The # operator" [spec]
		lit := ""
		for _, p := range args[t.ParamIndex] {
			raw := p.Raw
			if p.Type == StringLiteral {
				raw = strings.Replace(strings.Replace(p.Raw, `\`, `\\`, -1), `"`, `\"`, -1)
			}
			if p.Adjacent || lit == "" {
				lit += raw
			} else {
				lit += " " + raw
			}
		}
		raw := `"` + lit + `"`
		val, err := lex.ReadString(io.NewByteSource([]byte(raw)))
		if err != nil {
			return nil, err
		}
		r = append(r, &Token{
			Type: StringLiteral,
			Val:  val,
			Raw:  raw,
		})
	}
	return r, nil
}
