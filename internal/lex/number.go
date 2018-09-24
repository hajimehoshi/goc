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

package lex

import (
	"fmt"
	"io"

	"github.com/hajimehoshi/goc/internal/ctype"
	gio "github.com/hajimehoshi/goc/internal/io"
)

type IntegerSuffix int

const (
	IntegerSuffixNone IntegerSuffix = iota
	IntegerSuffixL
	IntegerSuffixLL
	IntegerSuffixU
	IntegerSuffixUL
	IntegerSuffixULL
)

func ReadIntegerSuffix(src gio.Source) (IntegerSuffix, error) {
	bs, err := src.Peek(3)
	if err != nil && err != io.EOF {
		return 0, err
	}
	s := ""
	for _, b := range bs {
		if ('A' <= b && b <= 'Z') || ('a' <= b && b <= 'z') {
			s += string(b)
			continue
		}
		break
	}

	switch s {
	case "":
		return IntegerSuffixNone, nil
	case "l", "L":
		src.Discard(1)
		return IntegerSuffixL, nil
	case "ll", "LL":
		src.Discard(2)
		return IntegerSuffixLL, nil
	case "u", "U":
		src.Discard(1)
		return IntegerSuffixU, nil
	case "ul", "UL":
		src.Discard(2)
		return IntegerSuffixUL, nil
	case "ull", "ULL":
		src.Discard(3)
		return IntegerSuffixULL, nil
	}

	return 0, fmt.Errorf("lex: unexpected suffix %q", s)
}

func ReadNumber(src gio.Source) (ctype.IntegerValue, error) {
	b, err := gio.ShouldReadByte(src)
	if err != nil {
		return ctype.IntegerValue{}, err
	}

	// TODO: Float number

	if !IsDigit(b) {
		return ctype.IntegerValue{}, fmt.Errorf("lex: non-digit character")
	}

	v := int64(0)

	if b == '0' {
		bs, err := src.Peek(1)
		if err != nil && err != io.EOF {
			return ctype.IntegerValue{}, err
		}
		if len(bs) < 1 {
			return ctype.IntegerValue{
				Type:  ctype.Int,
				Value: 0,
			}, nil
		}
		if bs[0] == 'x' || bs[0] == 'X' {
			src.Discard(1)
			for {
				bs, err := src.Peek(1)
				if err != nil && err != io.EOF {
					return ctype.IntegerValue{}, err
				}
				if len(bs) < 1 {
					break
				}
				if !isHexDigit(bs[0]) {
					break
				}
				src.Discard(1)
				v *= 16
				v += int64(hex(bs[0]))
			}
		}
		if IsDigit(bs[0]) {
			for {
				bs, err := src.Peek(1)
				if err != nil && err != io.EOF {
					return ctype.IntegerValue{}, err
				}
				if len(bs) < 1 {
					break
				}
				if !IsDigit(bs[0]) {
					break
				}
				if !isOctDigit(bs[0]) {
					return ctype.IntegerValue{}, fmt.Errorf("lex: malformed octal constant")
				}
				src.Discard(1)
				v *= 8
				v += int64(bs[0] - '0')
			}
		}
	} else {
		v = int64(b - '0')
		for {
			bs, err := src.Peek(1)
			if err != nil && err != io.EOF {
				return ctype.IntegerValue{}, err
			}
			if len(bs) < 1 {
				return ctype.IntegerValue{
					Type:  ctype.Int,
					Value: v,
				}, nil
			}
			if !IsDigit(bs[0]) {
				break
			}
			src.Discard(1)
			v *= 10
			v += int64(bs[0] - '0')
		}
	}

	s, err := ReadIntegerSuffix(src)
	if err != nil {
		return ctype.IntegerValue{}, err
	}
	switch s {
	case IntegerSuffixNone:
		if v >= 0x80000000 {
			return ctype.IntegerValue{
				Type:  ctype.LongLong,
				Value: v,
			}, nil
		}
		return ctype.IntegerValue{
			Type:  ctype.Int,
			Value: v,
		}, nil
	case IntegerSuffixL:
		if v >= 0x80000000 {
			return ctype.IntegerValue{
				Type:  ctype.LongLong,
				Value: v,
			}, nil
		}
		return ctype.IntegerValue{
			Type:  ctype.Long,
			Value: v,
		}, nil
	case IntegerSuffixLL:
		return ctype.IntegerValue{
			Type:  ctype.LongLong,
			Value: v,
		}, nil
	case IntegerSuffixU:
		if v >= 0x100000000 {
			return ctype.IntegerValue{
				Type:  ctype.ULongLong,
				Value: v,
			}, nil
		}
		return ctype.IntegerValue{
			Type:  ctype.UInt,
			Value: v,
		}, nil
	case IntegerSuffixUL:
		if v >= 0x100000000 {
			return ctype.IntegerValue{
				Type:  ctype.ULongLong,
				Value: v,
			}, nil
		}
		return ctype.IntegerValue{
			Type:  ctype.ULong,
			Value: v,
		}, nil
	case IntegerSuffixULL:
		return ctype.IntegerValue{
			Type:  ctype.ULongLong,
			Value: v,
		}, nil
	default:
		panic("not reached")
	}
}
