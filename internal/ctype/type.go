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

package ctype

type IntegerType int

type FloatType int

const (
	Int IntegerType = iota
	UInt
	Char
	UChar
	Short
	UShort
	Long
	ULong
	LongLong
	ULongLong
)

const (
	Float FloatType = iota
	Double
)

type IntegerValue struct {
	Type  IntegerType
	Value int64
}

type FloatValue struct {
	Type  FloatType
	Value float64
}

func (t IntegerType) String() string {
	switch t {
	case Int:
		return "int"
	case UInt:
		return "unsigned int"
	case Char:
		return "char"
	case UChar:
		return "unsigned char"
	case Short:
		return "short"
	case UShort:
		return "unsigned short"
	case Long:
		return "long"
	case ULong:
		return "unsigned long"
	case LongLong:
		return "long long"
	case ULongLong:
		return "unsigned long long"
	default:
		panic("not reached")
	}
}

func (t FloatType) String() string {
	switch t {
	case Float:
		return "float"
	case Double:
		return "double"
	default:
		panic("not reached")
	}
}
