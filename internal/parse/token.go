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

package parse

import (
	"fmt"

	"github.com/hajimehoshi/goc/internal/ctype"
	"github.com/hajimehoshi/goc/internal/io"
	"github.com/hajimehoshi/goc/internal/lex"
	"github.com/hajimehoshi/goc/internal/preprocess"
)

type Token struct {
	Type TokenType

	IntegerValue ctype.IntegerValue
	FloatValue   ctype.FloatValue
	StringValue  string

	Name string
}

type TokenReader interface {
	NextToken() (*Token, error)
}

type tokenReader struct {
	src preprocess.PPTokenReader
}

func (t *tokenReader) NextToken() (*Token, error) {
	p, err := t.src.NextPPToken()
	if err != nil {
		return nil, err
	}

	if p.Type < 128 && lex.IsSingleCharPunctuator(byte(p.Type)) {
		return &Token{
			Type: TokenType(p.Type),
		}, nil
	}
	switch p.Type {
	case preprocess.Arrow:
		return &Token{
			Type: TokenType(Arrow),
		}, nil
	case preprocess.Inc:
		return &Token{
			Type: TokenType(Inc),
		}, nil
	case preprocess.Dec:
		return &Token{
			Type: TokenType(Dec),
		}, nil
	case preprocess.Shl:
		return &Token{
			Type: TokenType(Shl),
		}, nil
	case preprocess.Shr:
		return &Token{
			Type: TokenType(Shr),
		}, nil
	case preprocess.Le:
		return &Token{
			Type: TokenType(Le),
		}, nil
	case preprocess.Ge:
		return &Token{
			Type: TokenType(Ge),
		}, nil
	case preprocess.Eq:
		return &Token{
			Type: TokenType(Eq),
		}, nil
	case preprocess.Ne:
		return &Token{
			Type: TokenType(Ne),
		}, nil
	case preprocess.AndAnd:
		return &Token{
			Type: TokenType(AndAnd),
		}, nil
	case preprocess.OrOr:
		return &Token{
			Type: TokenType(OrOr),
		}, nil
	case preprocess.DotDotDot:
		return &Token{
			Type: TokenType(DotDotDot),
		}, nil
	case preprocess.MulEq:
		return &Token{
			Type: TokenType(MulEq),
		}, nil
	case preprocess.DivEq:
		return &Token{
			Type: TokenType(DivEq),
		}, nil
	case preprocess.ModEq:
		return &Token{
			Type: TokenType(ModEq),
		}, nil
	case preprocess.AddEq:
		return &Token{
			Type: TokenType(AddEq),
		}, nil
	case preprocess.SubEq:
		return &Token{
			Type: TokenType(SubEq),
		}, nil
	case preprocess.ShlEq:
		return &Token{
			Type: TokenType(ShlEq),
		}, nil
	case preprocess.ShrEq:
		return &Token{
			Type: TokenType(ShrEq),
		}, nil
	case preprocess.AndEq:
		return &Token{
			Type: TokenType(AndEq),
		}, nil
	case preprocess.XorEq:
		return &Token{
			Type: TokenType(XorEq),
		}, nil
	case preprocess.OrEq:
		return &Token{
			Type: TokenType(OrEq),
		}, nil
	case preprocess.HashHash:
		return &Token{
			Type: TokenType(HashHash),
		}, nil
	case preprocess.Identifier:
		if t, ok := KeywordToTokenType(p.Val); ok {
			return &Token{
				Type: t,
			}, nil
		}
		return &Token{
			Type: Identifier,
			Name: p.Val,
		}, nil
	case preprocess.PPNumber:
		v, err := lex.ReadNumber(io.NewByteSource([]byte(p.Raw), ""))
		if err != nil {
			return nil, err
		}
		/*if bs.Len() > 0 {
			return nil, fmt.Errorf("token: invalid token: %q", p.Raw)
		}*/
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

func Tokenize(src preprocess.PPTokenReader) TokenReader {
	return &tokenReader{
		src: src,
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

type tokenReadPeeker struct {
	r   TokenReader
	buf *Token
}

func (t *tokenReadPeeker) NextToken() (*Token, error) {
	if t.buf != nil {
		tk := t.buf
		t.buf = nil
		return tk, nil
	}
	return t.r.NextToken()
}

func (t *tokenReadPeeker) peekToken() (*Token, error) {
	if t.buf != nil {
		return t.buf, nil
	}

	tk, err := t.r.NextToken()
	if err != nil {
		return nil, err
	}
	t.buf = tk
	return tk, nil
}
