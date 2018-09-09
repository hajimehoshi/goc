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

package ioutil_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	. "github.com/hajimehoshi/goc/internal/ioutil"
)

func TestBackslashNewLineStripper(t *testing.T) {
	cases := []struct {
		In  string
		Out string
	}{
		{"", ""},
		{"ABC", "ABC"},
		{"A\\B", "A\\B"},
		{"A\\\nB", "AB"},
		{"A\\\nB\\\nC", "ABC"},
		{"A\\\\\nB", "A\\B"},
		{"A\\\\\n\nB", "A\\\nB"},
		{"\\", "\\"},
		{"\\\n", ""},
		{"\n", "\n"},
	}
	for _, c := range cases {
		s := NewBackslashNewLineStripper(bytes.NewReader([]byte(c.In)))
		got, err := ioutil.ReadAll(s)
		if err != nil {
			t.Errorf("NewBackslashNewLineStripper(%q): err: %v", c.In, err)
			continue
		}
		if string(got) != c.Out {
			t.Errorf("NewBackslashNewLineStripper(%q): got %q, want: %q", c.In, string(got), c.Out)
		}
	}
}

func TestLastNewLineAdder(t *testing.T) {
	cases := []struct {
		In  string
		Out string
	}{
		{"", "\n"},
		{"ABC", "ABC\n"},
		{"ABC\n", "ABC\n"},
		{"ABC\n\n\n", "ABC\n\n\n"},
		{"ABC\n\n\n", "ABC\n\n\n"},
		{"A\nB\nC", "A\nB\nC\n"},
		{"\nA", "\nA\n"},
	}
	for _, c := range cases {
		s := NewLastNewLineAdder(bytes.NewReader([]byte(c.In)))
		got, err := ioutil.ReadAll(s)
		if err != nil {
			t.Errorf("NewLastNewLineAdder(%q): err: %v", c.In, err)
			continue
		}
		if string(got) != c.Out {
			t.Errorf("NewLastNewLineAdder(%q): got %q, want: %q", c.In, string(got), c.Out)
		}
	}
}
