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
	"bufio"
	"bytes"
	"fmt"

	"github.com/hajimehoshi/goc/internal/ctype"
	"github.com/hajimehoshi/goc/internal/lex"
	"github.com/hajimehoshi/goc/internal/preprocess"
)

type Token struct {
	Type Type

	IntegerValue ctype.IntegerValue
	FloatValue   ctype.FloatValue
	StringValue  string

	Name string
}

func FromPPToken(p *preprocess.Token) (*Token, error) {
	if p.Type < 128 && lex.IsSingleCharPunctuator(byte(p.Type)) {
		return &Token{
			Type: Type(p.Type),
		}, nil
	}
	switch p.Type {
	case preprocess.Arrow:
		return &Token{
			Type: Type(Arrow),
		}, nil
	case preprocess.Inc:
		return &Token{
			Type: Type(Inc),
		}, nil
	case preprocess.Dec:
		return &Token{
			Type: Type(Dec),
		}, nil
	case preprocess.Shl:
		return &Token{
			Type: Type(Shl),
		}, nil
	case preprocess.Shr:
		return &Token{
			Type: Type(Shr),
		}, nil
	case preprocess.Le:
		return &Token{
			Type: Type(Le),
		}, nil
	case preprocess.Ge:
		return &Token{
			Type: Type(Ge),
		}, nil
	case preprocess.Eq:
		return &Token{
			Type: Type(Eq),
		}, nil
	case preprocess.Ne:
		return &Token{
			Type: Type(Ne),
		}, nil
	case preprocess.AndAnd:
		return &Token{
			Type: Type(AndAnd),
		}, nil
	case preprocess.OrOr:
		return &Token{
			Type: Type(OrOr),
		}, nil
	case preprocess.DotDotDot:
		return &Token{
			Type: Type(DotDotDot),
		}, nil
	case preprocess.MulEq:
		return &Token{
			Type: Type(MulEq),
		}, nil
	case preprocess.DivEq:
		return &Token{
			Type: Type(DivEq),
		}, nil
	case preprocess.ModEq:
		return &Token{
			Type: Type(ModEq),
		}, nil
	case preprocess.AddEq:
		return &Token{
			Type: Type(AddEq),
		}, nil
	case preprocess.SubEq:
		return &Token{
			Type: Type(SubEq),
		}, nil
	case preprocess.ShlEq:
		return &Token{
			Type: Type(ShlEq),
		}, nil
	case preprocess.ShrEq:
		return &Token{
			Type: Type(ShrEq),
		}, nil
	case preprocess.AndEq:
		return &Token{
			Type: Type(AndEq),
		}, nil
	case preprocess.XorEq:
		return &Token{
			Type: Type(XorEq),
		}, nil
	case preprocess.OrEq:
		return &Token{
			Type: Type(OrEq),
		}, nil
	case preprocess.HashHash:
		return &Token{
			Type: Type(HashHash),
		}, nil
	case preprocess.Identifier:
		if t, ok := KeywordToType(p.Val); ok {
			return &Token{
				Type: t,
			}, nil
		}
		return &Token{
			Type: Identifier,
			Name: p.Val,
		}, nil
	case preprocess.PPNumber:
		bs := bytes.NewReader([]byte(p.Raw))
		v, err := lex.ReadNumber(bufio.NewReader(bs))
		if err != nil {
			return nil, err
		}
		if bs.Len() > 0 {
			println("!?: ", fmt.Sprintf("%q", p.Raw))
			return nil, fmt.Errorf("token: invalid token: %q", p.Raw)
		}
		return &Token{
			Type:         IntegerLiteral,
			IntegerValue: v,
		}, nil
	case preprocess.CharacterConstant:
		return &Token{
			Type: IntegerLiteral,
			IntegerValue: ctype.IntegerValue{
				Type:  ctype.Int,
				Value: int64(p.Val[0]),
			},
		}, nil
	case preprocess.StringLiteral:
		return &Token{
			Type:        StringLiteral,
			StringValue: p.Val,
		}, nil
	case preprocess.EOF:
		return &Token{
			Type: EOF,
		}, nil
	default:
		return nil, fmt.Errorf("token: invalid token: %q", p.Raw)
	}
}

func (t *Token) String() string {
	switch t.Type {
	case IntegerLiteral:
		return fmt.Sprintf("integer: %v (%s)", t.IntegerValue.Value, t.IntegerValue.Type)
	case FloatLiteral:
		return fmt.Sprintf("float: %v (%s)", t.FloatValue.Value, t.FloatValue.Type)
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
