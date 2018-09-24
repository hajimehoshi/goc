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

package lex_test

import (
	"testing"

	"github.com/hajimehoshi/goc/internal/io"
	. "github.com/hajimehoshi/goc/internal/lex"
)

func TestReadPPNumber(t *testing.T) {
	cases := []struct {
		In  string
		Out string
		Err bool
	}{
		{`.123`, `.123`, false},
		{`.12.3`, `.12.3`, false},
		{`.12...3`, `.12...3`, false},

		{`0`, `0`, false},
		{`00`, `00`, false},
		{`000u`, `000u`, false},

		{`1`, `1`, false},
		{`123`, `123`, false},
		{`987654321`, `987654321`, false},
		{`1l`, `1l`, false},
		{`16777216ULL`, `16777216ULL`, false},
		{`42+`, `42`, false},
		{`141421356ul-`, `141421356ul`, false},

		{`1e1`, `1e1`, false},
		{`1e+1`, `1e+1`, false},
		{`1E-1`, `1E-1`, false},
		{`1x+1`, `1x`, false},
		{`1eee`, `1eee`, false},
		{`1+`, `1`, false},

		{`0377`, `0377`, false},
		{`04444ll`, `04444ll`, false},
		{`0?`, `0`, false},
		{`08`, `08`, false},

		{`0x0`, `0x0`, false},
		{`0x0000`, `0x0000`, false},
		{`0xdeadbeefUL`, `0xdeadbeefUL`, false},
		{`0xffff`, `0xffff`, false},
		{`0Xffff`, `0Xffff`, false},
		{`0xFFFF`, `0xFFFF`, false},
		{`0XFFFF`, `0XFFFF`, false},

		{`x`, ``, true},
		{`.`, ``, true},
		{`..`, ``, true},
		{`.+`, ``, true},
	}
	for _, c := range cases {
		got, err := ReadPPNumber(io.NewByteSource([]byte(c.In)))
		if err != nil && !c.Err {
			t.Errorf("ReadPPNumber(%q) should not return error but did: %v", c.In, err)
		}
		if err == nil && c.Err {
			t.Errorf("ReadPPNumber(%q) should return error but not", c.In)
		}
		if got != c.Out {
			t.Errorf("ReadPPNumber(%q): got: %q, want: %q", c.In, got, c.Out)
		}
	}
}
