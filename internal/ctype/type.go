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

type (
	Int       int32
	UInt      uint32
	Char      int8
	UChar     uint8
	Short     int16
	UShort    uint16
	Long      int32
	ULong     uint32
	LongLong  int64
	ULongLong uint64
	Float     float32
	Double    float64
)

func (Int) TypeString() string {
	return "int"
}

func (UInt) TypeString() string {
	return "unsigned int"
}

func (Char) TypeString() string {
	return "char"
}

func (UChar) TypeString() string {
	return "unsigned char"
}

func (Long) TypeString() string {
	return "long"
}

func (ULong) TypeString() string {
	return "unsigned long"
}

func (LongLong) TypeString() string {
	return "long long"
}

func (ULongLong) TypeString() string {
	return "unsigned long long"
}

func (Float) TypeString() string {
	return "float"
}

func (Double) TypeString() string {
	return "double"
}
