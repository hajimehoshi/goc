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

type macro struct {
	name      string
	tokens    []*Token
	paramsLen int
}

func (m *macro) apply(src PPTokenReadPeeker, expandedFrom map[string]struct{}) ([]*Token, error) {
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
	t, err := src.PeekPPToken()
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
		r = append(r, args[t.ParamIndex]...)
	}
	return r, nil
}
