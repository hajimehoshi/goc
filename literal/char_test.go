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

package literal_test

import (
	"bytes"
	"bufio"
	"testing"

	"github.com/hajimehoshi/goc/ctype"
	. "github.com/hajimehoshi/goc/literal"
)

func TestReadEscapedChar(t *testing.T) {
	cases := []struct {
		In  string
		Out ctype.Int
		Err bool
	}{
		{`\n`, '\n', false},
		{`\t`, '\t', false},
		{`\\`, '\\', false},
		{`\'`, '\'', false},
		{`\"`, '"', false},
		{`\x00`, 0, false},
		{`\x12`, 0x12, false},
		{`\x20`, ' ', false},
		{`\xab`, 0xab, false},
		{`\xff`, 0xff, false},
		{`\0`, 0, false},
		{`\0?`, 0, false},
		{`\07`, 7, false},
		{`\08`, 0, false},
		{`\377`, 255, false},
		{`\778`, 63, false},

		{``, 0, true},
		{`a`, 0, true},
		{`\u`, 0, true},
		{`\x0g`, 0, true},
		{`\xf`, 0, true},
		{`\8`, 0, true},
		{`\400`, 0, true},
		{`\777`, 0, true},
		{"\n", 0, true},
		{"\r", 0, true},
	}
	for _, c := range cases {
		got, err := ReadEscapedChar(bufio.NewReader(bytes.NewReader([]byte(c.In))))
		if err != nil && !c.Err {
			t.Errorf("ReadEscapedChar(%q) should not return error but did: %v", c.In, err)
		}
		if err == nil && c.Err {
			t.Errorf("ReadEscapedChar(%q) should return error but not", c.In)
		}
		if got != c.Out {
			t.Errorf("ReadEscapedChar(%q): got: %q, want: %q", c.In, got, c.Out)
		}
	}
}

func TestReadChar(t *testing.T) {
	cases := []struct {
		In  string
		Out ctype.Int
		Err bool
	}{
		{`'0'`, '0', false},
		{`'a'`, 'a', false},
		{`' '`, ' ', false},
		{`'\n'`, '\n', false},
		{`'\t'`, '\t', false},
		{`'\\'`, '\\', false},
		{`'\''`, '\'', false},
		{`'\"'`, '"', false},
		{`'"'`, '"', false},
		{`'\x00'`, 0, false},
		{`'\xff'`, 0xff, false},
		{`'\0'`, 0, false},
		{`'\377'`, 255, false},

		{`'\778'`, 0, true},
		{`''`, 0, true},
		{`'''`, 0, true},
		{`'\8'`, 0, true},
		{`'\400'`, 0, true},
		{`'\777'`, 0, true},
		{"'\n'", 0, true},
		{"'\r'", 0, true},
	}
	for _, c := range cases {
		got, err := ReadChar(bufio.NewReader(bytes.NewReader([]byte(c.In))))
		if err != nil && !c.Err {
			t.Errorf("ReadChar(%q) should not return error but did: %v", c.In, err)
		}
		if err == nil && c.Err {
			t.Errorf("ReadChar(%q) should return error but not", c.In)
		}
		if got != c.Out {
			t.Errorf("ReadChar(%q): got: %q, want: %q", c.In, got, c.Out)
		}
	}
}
