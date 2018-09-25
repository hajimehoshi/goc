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

func TestReadString(t *testing.T) {
	cases := []struct {
		In  string
		Out string
		Err bool
	}{
		{`""`, "", false},
		{`""+`, "", false},
		{`"a"`, "a", false},
		{`"xyz"`, "xyz", false},
		{`"The quick brown fox jumps over the lazy dog."`, "The quick brown fox jumps over the lazy dog.", false},
		{`"\\\\"`, "\\\\", false},
		{`"\n"`, "\n", false},
		{`"\"\""`, "\"\"", false},
		{`"'"`, "'", false},
		{`"\'"`, "'", false},
		{`"\00"`, string(0), false},
		{`"\08"`, string(0) + "8", false},
		{`"\xff"`, string([]byte{0xff}), false},

		{"\"\n\"", "", true},
		{`"`, "", true},
		{`"'`, "", true},
		{`"a`, "", true},
	}
	for _, c := range cases {
		got, err := ReadString(io.NewSource([]byte(c.In), ""))
		if err != nil && !c.Err {
			t.Errorf("ReadString(%q) should not return error but did: %v", c.In, err)
		}
		if err == nil && c.Err {
			t.Errorf("ReadString(%q) should return error but not", c.In)
		}
		if got != c.Out {
			t.Errorf("ReadString(%q): got: %q, want: %q", c.In, got, c.Out)
		}
	}
}
