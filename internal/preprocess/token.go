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

package preprocess

type TokenType int

// "6.4 Lexical elements" [spec]

const (
	HeaderName TokenType = iota + 128
	Identifier
	PPNumber
	CharacterConstant
	StringLiteral

	// "6.4.6 Punctuators" [spec]
	Arrow     // ->
	Inc       // ++
	Dec       // --
	Shl       // <<
	Shr       // >>
	Le        // <=
	Ge        // >=
	Eq        // ==
	Ne        // !=
	AndAnd    // &&
	OrOr      // ||
	DotDotDot // ...
	MulEq     // *=
	DivEq     // /=
	ModEq     // %=
	AddEq     // +=
	SubEq     // -=
	ShlEq     // <<=
	ShrEq     // >>=
	AndEq     // &=
	XorEq     // ^=
	OrEq      // |=
	HashHash  // ##

	// "each non-white-space character that cannot be one of the above" [spec]
	Other

	EOF
)

type Token struct {
	Type     TokenType
	Val      string
	Raw      string
	Adjacent bool
}

func (t *Token) String() string {
	switch t.Type {
	case '\n':
		return `(\n)`
	case EOF:
		return `(eof)`
	default:
		return t.Raw
	}
}
