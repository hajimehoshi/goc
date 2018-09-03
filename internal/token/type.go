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

type Type int

const (
	NumberLiteral Type = iota + 128
	StringLiteral
	Ident

	Auto
	Bool
	Break
	Case
	Char
	Complex
	Const
	Continue
	Default
	Do
	Double
	Else
	Enum
	Extern
	Float
	For
	Goto
	If
	Imaginary
	Inline
	Int
	Long
	Register
	Restrict
	Return
	Short
	Signed
	Sizeof
	Static
	Struct
	Switch
	Typedef
	Union
	Unsigned
	Void
	Volatile
	While

	// 6.4.6 Punctuators
	Arrow    // ->
	Inc      // ++
	Dec      // --
	Shl      // <<
	Shr      // >>
	Le       // <=
	Ge       // >=
	Eq       // ==
	Ne       // !=
	AndAnd   // &&
	OrOr     // ||
	MulEq    // *=
	DivEq    // /=
	ModEq    // %=
	AddEq    // +=
	SubEq    // -=
	ShlEq    // <<=
	ShrEq    // >>=
	AndEq    // &=
	XorEq    // ^=
	OrEq     // |=
	HashHash // ##

	// TODO: Define these punctuators
	// <:
	// :>
	// <%
	// %>
	// %:
	// %:%:
)

func (t Type) String() string {
	switch t {
	case NumberLiteral:
		return "(number)"
	case StringLiteral:
		return "(string)"
	case Ident:
		return "(ident)"
	case Auto:
		return "auto"
	case Bool:
		return "_Bool"
	case Break:
		return "break"
	case Case:
		return "case"
	case Char:
		return "char"
	case Complex:
		return "_Complex"
	case Const:
		return "const"
	case Continue:
		return "continue"
	case Default:
		return "default"
	case Do:
		return "do"
	case Double:
		return "double"
	case Else:
		return "else"
	case Enum:
		return "enum"
	case Extern:
		return "extern"
	case Float:
		return "float"
	case For:
		return "fot"
	case Goto:
		return "goto"
	case If:
		return "if"
	case Imaginary:
		return "_Imaginary"
	case Inline:
		return "inline"
	case Int:
		return "int"
	case Long:
		return "long"
	case Register:
		return "register"
	case Restrict:
		return "restrict"
	case Return:
		return "return"
	case Short:
		return "short"
	case Signed:
		return "signed"
	case Sizeof:
		return "sizeof"
	case Static:
		return "static"
	case Struct:
		return "struct"
	case Switch:
		return "switch"
	case Typedef:
		return "typedef"
	case Union:
		return "union"
	case Unsigned:
		return "unsigned"
	case Void:
		return "void"
	case Volatile:
		return "volatile"
	case While:
		return "while"
	case AddEq:
		return "+="
	case SubEq:
		return "-="
	case MulEq:
		return "*="
	case DivEq:
		return "/="
	case ModEq:
		return "%="
	case ShlEq:
		return "<<="
	case ShrEq:
		return ">>="
	case AndEq:
		return "&="
	case XorEq:
		return "^="
	case OrEq:
		return "|="
	case Arrow:
		return "->"
	case Inc:
		return "++"
	case Dec:
		return "--"
	case Shl:
		return "<<"
	case Shr:
		return ">>"
	case Le:
		return "<="
	case Ge:
		return ">="
	case Eq:
		return "=="
	case Ne:
		return "!="
	case AndAnd:
		return "&&"
	case OrOr:
		return "||"
	default:
		if 0 <= t && t <= 127 {
			return string(t)
		}
		return fmt.Sprintf("invalid: %d", t)
	}
}

var keywordToType = map[string]Type{
	"auto":       Auto,
	"_Bool":      Bool,
	"break":      Break,
	"case":       Case,
	"char":       Char,
	"_Complex":   Complex,
	"const":      Const,
	"continue":   Continue,
	"default":    Default,
	"do":         Do,
	"double":     Double,
	"else":       Else,
	"enum":       Enum,
	"extern":     Extern,
	"float":      Float,
	"for":        For,
	"goto":       Goto,
	"if":         If,
	"_Imaginary": Imaginary,
	"inline":     Inline,
	"int":        Int,
	"long":       Long,
	"register":   Register,
	"restrict":   Restrict,
	"return":     Return,
	"short":      Short,
	"signed":     Signed,
	"sizeof":     Sizeof,
	"static":     Static,
	"struct":     Struct,
	"switch":     Switch,
	"typeof":     Typedef,
	"union":      Union,
	"unsigned":   Unsigned,
	"void":       Void,
	"volatile":   Volatile,
	"while":      While,
}

func KeywordToType(keyword string) (Type, bool) {
	t, ok := keywordToType[keyword]
	return t, ok
}
