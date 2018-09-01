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
		Out interface{}
		Err bool
	}{
		{`0`, ctype.Int(0), false},
		{`00`, ctype.Int(0), false},
		{`000u`, ctype.UInt(0), false},

		{`1`, ctype.Int(1), false},
		{`123`, ctype.Int(123), false},
		{`987654321`, ctype.Int(987654321), false},
		{`1l`, ctype.Long(1), false},
		{`16777216ULL`, ctype.ULongLong(16777216), false},
		{`42+`, ctype.Int(42), false},
		{`141421356ul-`, ctype.ULong(141421356), false},

		// Oct
		{`0377`, ctype.Int(255), false},
		{`04444ll`, ctype.LongLong(2340), false},
		{`0?`, ctype.Int(0), false},

		// Hex
		{`0x0`, ctype.Int(0), false},
		{`0x0000`, ctype.Int(0), false},
		{`0xdeadbeefUL`, ctype.ULong(0xdeadbeef), false},
		{`0xffff`, ctype.Int(65535), false},
		{`0Xffff`, ctype.Int(65535), false},
		{`0xFFFF`, ctype.Int(65535), false},
		{`0XFFFF`, ctype.Int(65535), false},

		// Limit
		{`0x7fffffff`, ctype.Int(0x7fffffff), false},
		{`0x7fffffffl`, ctype.Long(0x7fffffff), false},
		{`0x7fffffffll`, ctype.LongLong(0x7fffffff), false},
		{`0x7fffffffu`, ctype.UInt(0x7fffffff), false},
		{`0x7ffffffful`, ctype.ULong(0x7fffffff), false},
		{`0x7fffffffull`, ctype.ULongLong(0x7fffffff), false},
		{`0x80000000`, ctype.LongLong(0x80000000), false},
		{`0x80000000l`, ctype.LongLong(0x80000000), false},
		{`0x80000000ll`, ctype.LongLong(0x80000000), false},
		{`0x80000000u`, ctype.UInt(0x80000000), false},
		{`0x80000000ul`, ctype.ULong(0x80000000), false},
		{`0x80000000ull`, ctype.ULongLong(0x80000000), false},
		{`0xffffffff`, ctype.LongLong(0xffffffff), false},
		{`0xffffffffl`, ctype.LongLong(0xffffffff), false},
		{`0xffffffffll`, ctype.LongLong(0xffffffff), false},
		{`0xffffffffu`, ctype.UInt(0xffffffff), false},
		{`0xfffffffful`, ctype.ULong(0xffffffff), false},
		{`0xffffffffull`, ctype.ULongLong(0xffffffff), false},
		{`0x100000000`, ctype.LongLong(0x100000000), false},
		{`0x100000000l`, ctype.LongLong(0x100000000), false},
		{`0x100000000ll`, ctype.LongLong(0x100000000), false},
		{`0x100000000u`, ctype.ULongLong(0x100000000), false},
		{`0x100000000ul`, ctype.ULongLong(0x100000000), false},
		{`0x100000000ull`, ctype.ULongLong(0x100000000), false},
		{`0x7fffffffffffffff`, ctype.LongLong(0x7fffffffffffffff), false},
		{`0x7fffffffffffffffl`, ctype.LongLong(0x7fffffffffffffff), false},
		{`0x7fffffffffffffffll`, ctype.LongLong(0x7fffffffffffffff), false},
		{`0x7fffffffffffffffu`, ctype.ULongLong(0x7fffffffffffffff), false},
		{`0x7ffffffffffffffful`, ctype.ULongLong(0x7fffffffffffffff), false},
		{`0x7fffffffffffffffull`, ctype.ULongLong(0x7fffffffffffffff), false},

		{`08`, nil, true},
		{`x`, nil, true},
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
