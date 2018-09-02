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
)

type Token struct {
	Type Type

	NumberValue interface{}
	StringValue string

	Name string
}

func (t *Token) String() string {
	switch t.Type {
	case NumberLiteral:
		ts := t.NumberValue.(interface{ TypeString() string }).TypeString()
		return fmt.Sprintf("number: %v (%s)", t.NumberValue, ts)
	case StringLiteral:
		return fmt.Sprintf("string: %q", t.StringValue)
	case Ident:
		return fmt.Sprintf("ident: %s", t.Name)
	default:
		return t.Type.String()
	}
}
