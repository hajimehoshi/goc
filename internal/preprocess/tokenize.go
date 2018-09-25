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

import (
	"fmt"
	"io"
	"strings"

	gio "github.com/hajimehoshi/goc/internal/io"
	"github.com/hajimehoshi/goc/internal/lex"
)

type PPTokenReader interface {
	NextPPToken() (*Token, error)
}

type ppTokenReadPeeker interface {
	PPTokenReader
	peekPPToken() (*Token, error)
}

func nextExpected(t PPTokenReader, expected ...TokenType) (*Token, error) {
	tk, err := t.NextPPToken()
	if err != nil {
		return nil, err
	}
	for _, e := range expected {
		if tk.Type == e {
			return tk, nil
		}
	}

	s := []string{}
	for _, e := range expected {
		s = append(s, e.String())
	}
	return nil, fmt.Errorf("preprocess: expected %s but %s", strings.Join(s, ","), tk.Type)
}

type tokenizer struct {
	src gio.Source

	// ppstate represents the current context is in the preprocessor or not.
	// -1 means header-name is no longer expected in the current line.
	// 0 means the start of the new line (just after '\n' or the initial state).
	// 1 means the start of the line of preprocessing (just after '#').
	// 2 means header-name is expected (just after '#include').
	ppstate int

	isSpace  bool
	wasSpace bool
}

func (t *tokenizer) headerNameExpected() bool {
	return t.ppstate == 2
}

func (t *tokenizer) next() (*Token, error) {
	var tk *Token
	for {
		var err error
		tk, err = t.nextImpl(t.src)
		if tk == nil && err == nil {
			continue
		}
		if err != nil {
			if err == io.EOF && tk != nil {
				panic("not reached")
			}
			return nil, err
		}
		break
	}

	tk.Adjacent = !t.wasSpace

	switch tk.Type {
	case '\n':
		t.ppstate = 0
	case '#':
		if t.ppstate == 0 {
			t.ppstate = 1
		} else {
			t.ppstate = -1
		}
	case Identifier:
		if t.ppstate == 1 && tk.Raw == "include" {
			t.ppstate = 2
		} else {
			t.ppstate = -1
		}
	default:
		t.ppstate = -1
	}

	return tk, nil
}

func (t *tokenizer) nextImpl(src gio.Source) (*Token, error) {
	bs, err := src.Peek(3)
	if err != nil && err != io.EOF {
		return nil, err
	}
	if len(bs) == 0 {
		if err != io.EOF {
			panic("not reached")
		}
		return &Token{
			Type: EOF,
		}, nil
	}

	t.wasSpace = t.isSpace
	t.isSpace = lex.IsWhitespace(bs[0])

	switch b := bs[0]; b {
	case '\n':
		// New line; preprocessor uses this.
		gio.Discard(src, 1)
		return &Token{
			Type: TokenType(b),
			Val:  string(bs[:1]),
			Raw:  string(bs[:1]),
		}, nil
	case ' ', '\t', '\v', '\f', '\r':
		// Space
		gio.Discard(src, 1)
		return nil, nil
	case '+':
		if len(bs) >= 2 {
			switch bs[1] {
			case '+':
				gio.Discard(src, 2)
				return &Token{
					Type: Inc,
					Val:  string(bs[:2]),
					Raw:  string(bs[:2]),
				}, nil
			case '=':
				gio.Discard(src, 2)
				return &Token{
					Type: AddEq,
					Val:  string(bs[:2]),
					Raw:  string(bs[:2]),
				}, nil
			}
		}
	case '-':
		if len(bs) >= 2 {
			switch bs[1] {
			case '-':
				gio.Discard(src, 2)
				return &Token{
					Type: Dec,
					Val:  string(bs[:2]),
					Raw:  string(bs[:2]),
				}, nil
			case '=':
				gio.Discard(src, 2)
				return &Token{
					Type: SubEq,
					Val:  string(bs[:2]),
					Raw:  string(bs[:2]),
				}, nil
			case '>':
				gio.Discard(src, 2)
				return &Token{
					Type: Arrow,
					Val:  string(bs[:2]),
					Raw:  string(bs[:2]),
				}, nil
			}
		}
	case '*':
		if len(bs) >= 2 && bs[1] == '=' {
			gio.Discard(src, 2)
			return &Token{
				Type: MulEq,
				Val:  string(bs[:2]),
				Raw:  string(bs[:2]),
			}, nil
		}
	case '/':
		if len(bs) >= 2 {
			switch bs[1] {
			case '/':
				// Line comment
				gio.Discard(src, 2)
				for {
					bs, err := src.Peek(1)
					if err != nil && err != io.EOF {
						return nil, err
					}
					if len(bs) < 1 {
						break
					}
					if bs[0] == '\n' {
						break
					}
					gio.Discard(src, 1)
				}
				return nil, nil
			case '*':
				// Block comment
				gio.Discard(src, 2)
				for {
					bs, err := src.Peek(2)
					if err != nil && err != io.EOF {
						return nil, err
					}
					if len(bs) <= 1 {
						return nil, fmt.Errorf("preprocess: unclosed block comment")
					}
					if bs[0] == '*' && bs[1] == '/' {
						gio.Discard(src, 2)
						break
					}
					gio.Discard(src, 1)
				}
				return nil, nil
			case '=':
				gio.Discard(src, 2)
				return &Token{
					Type: DivEq,
					Val:  string(bs[:2]),
					Raw:  string(bs[:2]),
				}, nil
			}
		}
	case '%':
		if len(bs) >= 2 && bs[1] == '=' {
			gio.Discard(src, 2)
			return &Token{
				Type: ModEq,
				Val:  string(bs[:2]),
				Raw:  string(bs[:2]),
			}, nil
		}
	case '=':
		if len(bs) >= 2 && bs[1] == '=' {
			gio.Discard(src, 2)
			return &Token{
				Type: Eq,
				Val:  string(bs[:2]),
				Raw:  string(bs[:2]),
			}, nil
		}
	case '<':
		if t.headerNameExpected() {
			buf := gio.NewBufSource(src)
			val, err := lex.ReadHeaderName(buf)
			if err != nil {
				return nil, err
			}
			return &Token{
				Type: HeaderName,
				Val:  val,
				Raw:  buf.Buf(),
			}, nil
		}
		if len(bs) >= 2 && bs[1] == '<' {
			if len(bs) >= 3 && bs[2] == '=' {
				gio.Discard(src, 3)
				return &Token{
					Type: ShlEq,
					Val:  string(bs[:3]),
					Raw:  string(bs[:3]),
				}, nil
			}
			gio.Discard(src, 2)
			return &Token{
				Type: Shl,
				Val:  string(bs[:2]),
				Raw:  string(bs[:2]),
			}, nil
		}
	case '>':
		if len(bs) >= 2 && bs[1] == '>' {
			if len(bs) >= 3 && bs[2] == '=' {
				gio.Discard(src, 3)
				return &Token{
					Type: ShrEq,
					Val:  string(bs[:3]),
					Raw:  string(bs[:3]),
				}, nil
			}
			gio.Discard(src, 2)
			return &Token{
				Type: Shr,
				Val:  string(bs[:2]),
				Raw:  string(bs[:2]),
			}, nil
		}
	case '&':
		if len(bs) >= 2 {
			switch bs[1] {
			case '&':
				gio.Discard(src, 2)
				return &Token{
					Type: AndAnd,
					Val:  string(bs[:2]),
					Raw:  string(bs[:2]),
				}, nil
			case '=':
				gio.Discard(src, 2)
				return &Token{
					Type: AndEq,
					Val:  string(bs[:2]),
					Raw:  string(bs[:2]),
				}, nil
			}
		}
	case '|':
		if len(bs) >= 2 {
			switch bs[1] {
			case '|':
				gio.Discard(src, 2)
				return &Token{
					Type: OrOr,
					Val:  string(bs[:2]),
					Raw:  string(bs[:2]),
				}, nil
			case '=':
				gio.Discard(src, 2)
				return &Token{
					Type: OrEq,
					Val:  string(bs[:2]),
					Raw:  string(bs[:2]),
				}, nil
			}
		}
	case '!':
		if len(bs) >= 2 && bs[1] == '=' {
			gio.Discard(src, 2)
			return &Token{
				Type: Ne,
				Val:  string(bs[:2]),
				Raw:  string(bs[:2]),
			}, nil
		}
	case '^':
		if len(bs) >= 2 && bs[1] == '=' {
			gio.Discard(src, 2)
			return &Token{
				Type: XorEq,
				Val:  string(bs[:2]),
				Raw:  string(bs[:2]),
			}, nil
		}
	case '\'':
		// Char literal
		buf := gio.NewBufSource(src)
		val, err := lex.ReadChar(buf)
		if err != nil {
			return nil, err
		}
		return &Token{
			Type: CharacterConstant,
			Val:  string(val),
			Raw:  buf.Buf(),
		}, nil
	case '"':
		if t.headerNameExpected() {
			buf := gio.NewBufSource(src)
			val, err := lex.ReadHeaderName(buf)
			if err != nil {
				return nil, err
			}
			return &Token{
				Type: HeaderName,
				Val:  val,
				Raw:  buf.Buf(),
			}, nil
		}
		// String literal
		buf := gio.NewBufSource(src)
		val, err := lex.ReadString(buf)
		if err != nil {
			return nil, err
		}
		return &Token{
			Type: StringLiteral,
			Val:  val,
			Raw:  buf.Buf(),
		}, nil
	case '.':
		if len(bs) >= 2 {
			if bs[1] == '.' && len(bs) >= 3 && bs[2] == '.' {
				gio.Discard(src, 3)
				return &Token{
					Type: DotDotDot,
					Val:  string(bs[:3]),
					Raw:  string(bs[:3]),
				}, nil
			}
			if lex.IsDigit(bs[1]) {
				buf := gio.NewBufSource(src)
				val, err := lex.ReadPPNumber(buf)
				if err != nil {
					return nil, err
				}
				return &Token{
					Type: PPNumber,
					Val:  val,
					Raw:  buf.Buf(),
				}, nil
			}
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		buf := gio.NewBufSource(src)
		val, err := lex.ReadPPNumber(buf)
		if err != nil {
			return nil, err
		}
		return &Token{
			Type: PPNumber,
			Val:  val,
			Raw:  buf.Buf(),
		}, nil
	case '#':
		if len(bs) >= 2 && bs[1] == '#' {
			gio.Discard(src, 2)
			return &Token{
				Type: HashHash,
				Val:  string(bs[:2]),
				Raw:  string(bs[:2]),
			}, nil
		}
	case ';', '(', ')', ',', '{', '}', '[', ']', '?', ':', '~':
		// Single character token
	default:
		if lex.IsNondigit(b) {
			name, err := lex.ReadIdentifier(src)
			if err != nil {
				return nil, err
			}
			return &Token{
				Type: Identifier,
				Val:  name,
				Raw:  name,
			}, nil
		}

		val := []byte{}
		for {
			bs, err := src.Peek(1)
			if err != nil && err != io.EOF {
				return nil, err
			}
			if len(bs) < 1 {
				break
			}
			b := bs[0]
			if lex.IsSingleCharPunctuator(b) {
				break
			}
			if lex.IsDigit(b) || lex.IsNondigit(b) {
				break
			}
			if lex.IsWhitespace(b) {
				break
			}
			gio.Discard(src, 1)
			val = append(val, b)
		}

		return &Token{
			Type: Other,
			Val:  string(val),
			Raw:  string(val),
		}, nil
	}

	// Single character token
	gio.Discard(src, 1)
	return &Token{
		Type: TokenType(bs[0]),
		Val:  string(bs[:1]),
		Raw:  string(bs[:1]),
	}, nil
}

func (t *tokenizer) NextPPToken() (*Token, error) {
	for {
		tk, err := t.next()
		if err != nil {
			return nil, err
		}
		if tk == nil {
			continue
		}
		return tk, nil
	}
}

func Tokenize(src []byte, filename string) PPTokenReader {
	return &tokenizer{
		src: gio.NewByteSource(src, filename),
	}
}
