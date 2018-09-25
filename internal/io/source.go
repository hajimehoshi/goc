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

package io

import (
	"bufio"
	"bytes"
)

type Source interface {
	ReadByte() (byte, error)
	Peek(int) ([]byte, error)
	Filename() string
	LineNo() int
	ByteNo() int
}

type source struct {
	r        *bufio.Reader
	filename string
	lineno   int
	byteno   int
}

func NewSource(src []byte, filename string) Source {
	if len(src) == 0 || src[len(src)-1] != '\n' {
		src = append(src, '\n')
	}
	r := NewBackslashNewLineStripper(bytes.NewReader(src))
	return &source{
		r:        bufio.NewReader(r),
		filename: filename,
	}
}

func (s *source) ReadByte() (byte, error) {
	b, err := s.r.ReadByte()
	if err != nil {
		return 0, err
	}
	if b == '\n' {
		s.lineno++
	}
	s.byteno++
	return b, nil
}

func (s *source) Peek(n int) ([]byte, error) {
	return s.r.Peek(n)
}

func (s *source) Filename() string {
	return s.filename
}

func (s *source) LineNo() int {
	return s.lineno
}

func (s *source) ByteNo() int {
	return s.byteno
}

type BufSource struct {
	src Source
	raw []byte
}

func NewBufSource(src Source) *BufSource {
	return &BufSource{
		src: src,
	}
}

func (s *BufSource) Buf() string {
	return string(s.raw)
}

func (s *BufSource) ReadByte() (byte, error) {
	b, err := s.src.ReadByte()
	if err != nil {
		return 0, err
	}
	s.raw = append(s.raw, b)
	return b, nil
}

func (s *BufSource) Peek(n int) ([]byte, error) {
	return s.src.Peek(n)
}

func (s *BufSource) Filename() string {
	return s.src.Filename()
}

func (s *BufSource) LineNo() int {
	return s.src.LineNo()
}

func (s *BufSource) ByteNo() int {
	return s.src.ByteNo()
}
