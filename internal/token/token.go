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

package token

import (
	"fmt"

	"github.com/hajimehoshi/goc/internal/ctype"
)

type Token struct {
	Type Type

	IntegerValue ctype.IntegerValue
	FloatValue   ctype.FloatValue
	StringValue  string

	Name string
}

func (t *Token) String() string {
	switch t.Type {
	case IntegerLiteral:
		return fmt.Sprintf("integer: %v (%s)", t.IntegerValue, t.Type)
	case FloatLiteral:
		return fmt.Sprintf("float: %v (%s)", t.FloatValue, t.Type)
	case StringLiteral:
		return fmt.Sprintf("string: %q", t.StringValue)
	case HeaderName:
		return fmt.Sprintf("header-name: %q", t.StringValue)
	case Identifier:
		return fmt.Sprintf("identifier: %s", t.Name)
	case EOF:
		return "eof"
	default:
		return t.Type.String()
	}
}
