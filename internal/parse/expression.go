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
	"strings"
)

func peekExpected(src *tokenReadPeeker, expected ...TokenType) (*Token, error) {
	tk, err := src.peekToken()
	if err != nil {
		return nil, err
	}
	for _, e := range expected {
		if tk.Type == e {
			return tk, nil
		}
	}
	return nil, nil
}

func nextExpected(src *tokenReadPeeker, expected ...TokenType) (*Token, error) {
	tk, err := src.peekToken()
	if err != nil {
		return nil, err
	}
	for _, e := range expected {
		if tk.Type == e {
			if _, err := src.NextToken(); err != nil {
				return nil, err
			}
			return tk, nil
		}
	}

	s := []string{}
	for _, e := range expected {
		s = append(s, e.String())
	}
	return nil, fmt.Errorf("parse: expected %s but %s", strings.Join(s, ","), tk.Type)
}

type Parser struct {
	src    *tokenReadPeeker
	errors []error
}

func (p *Parser) appendError(err error) {
	p.errors = append(p.errors, err)
}

type Expression interface {
	IsUnaryExpression() bool
}

type UnaryExpression struct{}

func (*UnaryExpression) IsUnaryExpression() bool {
	return true
}

type BiOpExpression struct {
	Op  TokenType
	Lhs Expression
	Rhs Expression
}

func (*BiOpExpression) IsUnaryExpression() bool {
	return false
}

type TriOpExpression struct {
	Op   TokenType
	Exp1 Expression
	Exp2 Expression
	Exp3 Expression
}

func (*TriOpExpression) IsUnaryExpression() bool {
	return false
}

func (p *Parser) ParseLogicalOrExpression() Expression {
	return nil
}

func (p *Parser) ParseConditionalExpression() Expression {
	exp1 := p.ParseLogicalOrExpression()

	t, err := peekExpected(p.src, '?')
	if err != nil {
		p.appendError(err)
		return exp1
	}
	if t == nil {
		return exp1
	}
	if _, err := p.src.NextToken(); err != nil {
		p.appendError(err)
		return exp1
	}

	exp2 := p.ParseExpression()

	if _, err := nextExpected(p.src, ':'); err != nil {
		p.appendError(err)
		return exp1
	}

	exp3 := p.ParseConditionalExpression()

	return &TriOpExpression{
		Op:   '?',
		Exp1: exp1,
		Exp2: exp2,
		Exp3: exp3,
	}
}

func (p *Parser) ParseAssignmentExpression() Expression {
	lhs := p.ParseConditionalExpression()

	if !lhs.IsUnaryExpression() {
		return lhs
	}

	t, err := peekExpected(p.src, '=', MulEq, DivEq, ModEq, AddEq, SubEq, ShlEq, ShrEq, AndEq, XorEq, OrEq)
	if err != nil {
		p.appendError(err)
		return lhs
	}
	if t == nil {
		return lhs
	}
	if _, err := p.src.NextToken(); err != nil {
		p.appendError(err)
		return lhs
	}

	return &BiOpExpression{
		Op:  t.Type,
		Lhs: lhs,
		Rhs: p.ParseAssignmentExpression(),
	}
}

func (p *Parser) ParseExpression() Expression {
	lhs := p.ParseAssignmentExpression()

	t, err := peekExpected(p.src, ',')
	if err != nil {
		p.appendError(err)
		return lhs
	}
	if t == nil {
		return lhs
	}
	if _, err := p.src.NextToken(); err != nil {
		p.appendError(err)
		return lhs
	}

	return &BiOpExpression{
		Op:  ',',
		Lhs: lhs,
		Rhs: p.ParseExpression(),
	}
}
