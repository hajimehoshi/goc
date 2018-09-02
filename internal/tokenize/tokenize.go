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

package tokenize

import (
	"fmt"
	"bufio"
	"io"

	"github.com/hajimehoshi/goc/internal/literal"
	"github.com/hajimehoshi/goc/internal/ioutil"
	"github.com/hajimehoshi/goc/internal/token"
)

func isIdentFirstChar(c byte) bool {
	if 'A' <= c && c <= 'Z' {
		return true
	}
	if 'a' <= c && c <= 'z' {
		return true
	}
	if c == '_' {
		return true
	}
	return false
}

func isIdentChar(c byte) bool {
	if isIdentFirstChar(c) {
		return true
	}
	if '0' <= c && c <= '9' {
		return true
	}
	return false
}

func nextToken(src *bufio.Reader) (*token.Token, error) {
	// TODO: Can this read runes intead of bytes?
	bs, err := src.Peek(3)
	if err != nil && err != io.EOF {
		return nil, err
	}
	if len(bs) == 0 {
		if err != io.EOF {
			panic("not reached")
		}
		return nil, err
	}
	switch b := bs[0]; b {
	case '\n':
		// New line; preprocessor uses this.
		src.Discard(1)
		return &token.Token{
			Type: token.Type(b),
		}, nil
	case ' ', '\t', '\v', '\f', '\r':
		// Space
		src.Discard(1)
		return nil, nil
	case '+':
		if len(bs) >= 2 {
			switch bs[1] {
			case '+':
				src.Discard(2)
				return &token.Token{
					Type: token.Inc,
				}, nil
			case '=':
				src.Discard(2)
				return &token.Token{
					Type: token.AddEq,
				}, nil
			}
		}
	case '-':
		if len(bs) >= 2 {
			switch bs[1] {
			case '-':
				src.Discard(2)
				return &token.Token{
					Type: token.Dec,
				}, nil
			case '=':
				src.Discard(2)
				return &token.Token{
					Type: token.SubEq,
				}, nil
			case '>':
				src.Discard(2)
				return &token.Token{
					Type: token.Arrow,
				}, nil
			}
		}
	case '*':
		if len(bs) >= 2 && bs[1] == '=' {
			src.Discard(2)
			return &token.Token{
				Type: token.MulEq,
			}, nil
		}
	case '/':
		if len(bs) >= 2 {
			switch bs[1] {
			case '/':
				// Line comment
				src.Discard(2)
				for {
					b, err := src.ReadByte()
					if err != nil && err != io.EOF {
						return nil, err
					}
					if b == '\n' {
						break
					}
					if err == io.EOF {
						break
					}
				}
				return nil, nil
			case '*':
				// Block comment
				src.Discard(2)
				for {
					bs, err := src.Peek(2)
					if err != nil && err != io.EOF {
						return nil, err
					}
					if len(bs) <= 1 {
						return nil, fmt.Errorf("tokenizer: unclosed block comment")
					}
					if bs[0] == '*' && bs[1] == '/' {
						src.Discard(2)
						break
					}
					src.Discard(1)
				}
				return nil, nil
			case '=':
				src.Discard(2)
				return &token.Token{
					Type: token.DivEq,
				}, nil
			}
		}
	case '%':
		if len(bs) >= 2 && bs[1] == '=' {
			src.Discard(2)
			return &token.Token{
				Type: token.ModEq,
			}, nil
		}
	case '=':
		if len(bs) >= 2 && bs[1] == '=' {
			src.Discard(2)
			return &token.Token{
				Type: token.Eq,
			}, nil
		}
	case '<':
		if len(bs) >= 2 && bs[1] == '<' {
			if len(bs) >= 3 && bs[2] == '=' {
				src.Discard(3)
				return &token.Token{
					Type: token.ShlEq,
				}, nil
			}
			src.Discard(2)
			return &token.Token{
				Type: token.Shl,
			}, nil
		}
	case '>':
		if len(bs) >= 2 && bs[1] == '>' {
			if len(bs) >= 3 && bs[2] == '=' {
				src.Discard(3)
				return &token.Token{
					Type: token.ShrEq,
				}, nil
			}
			src.Discard(2)
			return &token.Token{
				Type: token.Shr,
			}, nil
		}
	case '&':
		if len(bs) >= 2 {
			switch bs[1] {
			case '&':
				src.Discard(2)
				return &token.Token{
					Type: token.AndAnd,
				}, nil
			case '=':
				src.Discard(2)
				return &token.Token{
					Type: token.AndEq,
				}, nil
			}
		}
	case '|':
		if len(bs) >= 2 {
			switch bs[1] {
			case '|':
				src.Discard(2)
				return &token.Token{
					Type: token.OrOr,
				}, nil
			case '=':
				src.Discard(2)
				return &token.Token{
					Type: token.OrEq,
				}, nil
			}
		}
	case '!':
		if len(bs) >= 2 && bs[1] == '=' {
			src.Discard(2)
			return &token.Token{
				Type: token.Ne,
			}, nil
		}
	case '^':
		if len(bs) >= 2 && bs[1] == '=' {
			src.Discard(2)
			return &token.Token{
				Type: token.XorEq,
			}, nil
		}
	case '\'':
		// Char literal
		n, err := literal.ReadChar(src)
		if err != nil {
			return nil, err
		}
		return &token.Token{
			Type:        token.NumberLiteral,
			NumberValue: n,
		}, nil
	case '"':
		// String literal
		s, err := literal.ReadString(src)
		if err != nil {
			return nil, err
		}
		return &token.Token{
			Type:        token.StringLiteral,
			StringValue: s,
		}, nil
	case '.':
		if len(bs) >= 2 && '0' <= bs[1] && bs[1] <= '9' {
			n, err := literal.ReadNumber(src)
			if err != nil {
				return nil, err
			}
			return &token.Token{
				Type:        token.NumberLiteral,
				NumberValue: n,
			}, nil
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		n, err := literal.ReadNumber(src)
		if err != nil {
			return nil, err
		}
		return &token.Token{
			Type:        token.NumberLiteral,
			NumberValue: n,
		}, nil
	case ';', '(', ')', ',', '{', '}', '[', ']', '?', ':', '~', '#':
		// Single char token
	default:
		if isIdentFirstChar(b) {
			name := []byte{b}
			src.Discard(1)
			for {
				bs, err := src.Peek(1)
				if err != nil && err != io.EOF {
					return nil, err
				}
				if len(bs) < 1 {
					break
				}
				if !isIdentChar(bs[0]) {
					break
				}
				src.Discard(1)
				name = append(name, bs[0])
			}
			if t, ok := token.KeywordToType(string(name)); ok {
				return &token.Token{
					Type: t,
				}, nil
			}
			return &token.Token{
				Type: token.Ident,
				Name: string(name),
			}, nil
		}
		// Invalid
		return nil, fmt.Errorf("tokenizer: invalid token: %s", string(b))
	}

	src.Discard(1)
	return &token.Token{
		Type: token.Type(bs[0]),
	}, nil
}

func scan(src io.Reader) ([]*token.Token, error) {
	buf := bufio.NewReader(src)
	ts := []*token.Token{}
	for {
		t, err := nextToken(buf)
		if t != nil {
			ts = append(ts, t)
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}
	return ts, nil
}

func stripNewLineTokens(tokens []*token.Token) []*token.Token {
	r := []*token.Token{}
	for _, t := range tokens {
		if t.Type == '\n' {
			continue
		}
		r = append(r, t)
	}
	return r
}

func joinStringLiterals(tokens []*token.Token) []*token.Token {
	r := []*token.Token{}
	for _, t := range tokens {
		var last *token.Token
		if len(r) > 0 {
			last = r[len(r)-1]
		}
		if last != nil && last.Type == token.StringLiteral && t.Type == token.StringLiteral {
			last.StringValue += t.StringValue
			continue
		}
		r = append(r, t)
	}
	return r
}

func Tokenize(src io.Reader) ([]*token.Token, error) {
	// TODO: Add TokenReader instead of using slices
	// TODO: Count line numbers
	// TODO: Preprocessor
	tokens, err := scan(ioutil.NewBackslashNewLineStripper(src))
	if err != nil {
		return nil, err
	}
	tokens = stripNewLineTokens(tokens)
	tokens = joinStringLiterals(tokens)
	return tokens, nil
}
