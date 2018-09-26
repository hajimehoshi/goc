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
	"bufio"
	"bytes"
	"testing"

	"github.com/hajimehoshi/goc/internal/ctype"
	. "github.com/hajimehoshi/goc/internal/lex"
)

func TestReadIntegerSuffix(t *testing.T) {
	cases := []struct {
		In  string
		Out IntegerSuffix
		Err bool
	}{
		{``, IntegerSuffixNone, false},
		{`+`, IntegerSuffixNone, false},
		{`l`, IntegerSuffixL, false},
		{`ll`, IntegerSuffixLL, false},
		{`u`, IntegerSuffixU, false},
		{`ul`, IntegerSuffixUL, false},
		{`ull`, IntegerSuffixULL, false},
		{`L`, IntegerSuffixL, false},
		{`LL`, IntegerSuffixLL, false},
		{`U`, IntegerSuffixU, false},
		{`UL`, IntegerSuffixUL, false},
		{`ULL`, IntegerSuffixULL, false},
		{`l+`, IntegerSuffixL, false},
		{`ll-`, IntegerSuffixLL, false},
		{`u*`, IntegerSuffixU, false},
		{`ul/`, IntegerSuffixUL, false},
		{`ull `, IntegerSuffixULL, false},
		{`lL`, 0, true},
		{`Ll`, 0, true},
		{`ulL`, 0, true},
		{`uLl`, 0, true},
		{`Ull`, 0, true},
		{`la`, 0, true},
		{`ZZ`, 0, true},
	}
	for _, c := range cases {
		got, err := ReadIntegerSuffix(bufio.NewReader(bytes.NewReader([]byte(c.In))))
		if err != nil && !c.Err {
			t.Errorf("ReadIntegerSuffix(%q) should not return error but did: %v", c.In, err)
		}
		if err == nil && c.Err {
			t.Errorf("ReadIntegerSuffix(%q) should return error but not", c.In)
		}
		if got != c.Out {
			t.Errorf("ReadIntegerSuffix(%q): got: %d, want: %d", c.In, got, c.Out)
		}
	}
}

func TestReadNumber(t *testing.T) {
	cases := []struct {
		In  string
		Out ctype.IntegerValue
		Err bool
	}{
		{`0`, ctype.IntegerValue{Type: ctype.Int, Value: 0}, false},
		{`00`, ctype.IntegerValue{Type: ctype.Int, Value: 0}, false},
		{`000u`, ctype.IntegerValue{Type: ctype.UInt, Value: 0}, false},

		{`1`, ctype.IntegerValue{Type: ctype.Int, Value: 1}, false},
		{`123`, ctype.IntegerValue{Type: ctype.Int, Value: 123}, false},
		{`987654321`, ctype.IntegerValue{Type: ctype.Int, Value: 987654321}, false},
		{`1l`, ctype.IntegerValue{Type: ctype.Long, Value: 1}, false},
		{`16777216ULL`, ctype.IntegerValue{Type: ctype.ULongLong, Value: 16777216}, false},
		{`42+`, ctype.IntegerValue{Type: ctype.Int, Value: 42}, false},
		{`141421356ul-`, ctype.IntegerValue{Type: ctype.ULong, Value: 141421356}, false},

		// Oct
		{`0377`, ctype.IntegerValue{Type: ctype.Int, Value: 255}, false},
		{`04444ll`, ctype.IntegerValue{Type: ctype.LongLong, Value: 2340}, false},
		{`0?`, ctype.IntegerValue{Type: ctype.Int, Value: 0}, false},

		// Hex
		{`0x0`, ctype.IntegerValue{Type: ctype.Int, Value: 0}, false},
		{`0x0000`, ctype.IntegerValue{Type: ctype.Int, Value: 0}, false},
		{`0xdeadbeefUL`, ctype.IntegerValue{Type: ctype.ULong, Value: 0xdeadbeef}, false},
		{`0xffff`, ctype.IntegerValue{Type: ctype.Int, Value: 65535}, false},
		{`0Xffff`, ctype.IntegerValue{Type: ctype.Int, Value: 65535}, false},
		{`0xFFFF`, ctype.IntegerValue{Type: ctype.Int, Value: 65535}, false},
		{`0XFFFF`, ctype.IntegerValue{Type: ctype.Int, Value: 65535}, false},

		// Limit
		{`0x7fffffff`, ctype.IntegerValue{Type: ctype.Int, Value: 0x7fffffff}, false},
		{`0x7fffffffl`, ctype.IntegerValue{Type: ctype.Long, Value: 0x7fffffff}, false},
		{`0x7fffffffll`, ctype.IntegerValue{Type: ctype.LongLong, Value: 0x7fffffff}, false},
		{`0x7fffffffu`, ctype.IntegerValue{Type: ctype.UInt, Value: 0x7fffffff}, false},
		{`0x7ffffffful`, ctype.IntegerValue{Type: ctype.ULong, Value: 0x7fffffff}, false},
		{`0x7fffffffull`, ctype.IntegerValue{Type: ctype.ULongLong, Value: 0x7fffffff}, false},
		{`0x80000000`, ctype.IntegerValue{Type: ctype.LongLong, Value: 0x80000000}, false},
		{`0x80000000l`, ctype.IntegerValue{Type: ctype.LongLong, Value: 0x80000000}, false},
		{`0x80000000ll`, ctype.IntegerValue{Type: ctype.LongLong, Value: 0x80000000}, false},
		{`0x80000000u`, ctype.IntegerValue{Type: ctype.UInt, Value: 0x80000000}, false},
		{`0x80000000ul`, ctype.IntegerValue{Type: ctype.ULong, Value: 0x80000000}, false},
		{`0x80000000ull`, ctype.IntegerValue{Type: ctype.ULongLong, Value: 0x80000000}, false},
		{`0xffffffff`, ctype.IntegerValue{Type: ctype.LongLong, Value: 0xffffffff}, false},
		{`0xffffffffl`, ctype.IntegerValue{Type: ctype.LongLong, Value: 0xffffffff}, false},
		{`0xffffffffll`, ctype.IntegerValue{Type: ctype.LongLong, Value: 0xffffffff}, false},
		{`0xffffffffu`, ctype.IntegerValue{Type: ctype.UInt, Value: 0xffffffff}, false},
		{`0xfffffffful`, ctype.IntegerValue{Type: ctype.ULong, Value: 0xffffffff}, false},
		{`0xffffffffull`, ctype.IntegerValue{Type: ctype.ULongLong, Value: 0xffffffff}, false},
		{`0x100000000`, ctype.IntegerValue{Type: ctype.LongLong, Value: 0x100000000}, false},
		{`0x100000000l`, ctype.IntegerValue{Type: ctype.LongLong, Value: 0x100000000}, false},
		{`0x100000000ll`, ctype.IntegerValue{Type: ctype.LongLong, Value: 0x100000000}, false},
		{`0x100000000u`, ctype.IntegerValue{Type: ctype.ULongLong, Value: 0x100000000}, false},
		{`0x100000000ul`, ctype.IntegerValue{Type: ctype.ULongLong, Value: 0x100000000}, false},
		{`0x100000000ull`, ctype.IntegerValue{Type: ctype.ULongLong, Value: 0x100000000}, false},
		{`0x7fffffffffffffff`, ctype.IntegerValue{Type: ctype.LongLong, Value: 0x7fffffffffffffff}, false},
		{`0x7fffffffffffffffl`, ctype.IntegerValue{Type: ctype.LongLong, Value: 0x7fffffffffffffff}, false},
		{`0x7fffffffffffffffll`, ctype.IntegerValue{Type: ctype.LongLong, Value: 0x7fffffffffffffff}, false},
		{`0x7fffffffffffffffu`, ctype.IntegerValue{Type: ctype.ULongLong, Value: 0x7fffffffffffffff}, false},
		{`0x7ffffffffffffffful`, ctype.IntegerValue{Type: ctype.ULongLong, Value: 0x7fffffffffffffff}, false},
		{`0x7fffffffffffffffull`, ctype.IntegerValue{Type: ctype.ULongLong, Value: 0x7fffffffffffffff}, false},

		{`08`, ctype.IntegerValue{}, true},
		{`x`, ctype.IntegerValue{}, true},
	}
	for _, c := range cases {
		got, err := ReadNumber(bufio.NewReader(bytes.NewReader([]byte(c.In))))
		if err != nil && !c.Err {
			t.Errorf("ReadNumber(%q) should not return error but did: %v", c.In, err)
		}
		if err == nil && c.Err {
			t.Errorf("ReadNumber(%q) should return error but not", c.In)
		}
		if got != c.Out {
			t.Errorf("ReadNumber(%q): got: %[2]d (%[2]T), want: %[3]d (%[3]T)", c.In, got, c.Out)
		}
	}
}
