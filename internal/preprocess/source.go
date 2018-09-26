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
	"io"
)

type Source struct {
	src []byte
	pos int

	filename string
	lineno   int
}

func NewSource(src []byte, filename string) *Source {
	if len(src) == 0 || src[len(src)-1] != '\n' {
		src = append(src, '\n')
	}
	return &Source{
		src:      src,
		filename: filename,
	}
}

func (s *Source) ReadByte() (byte, error) {
	for {
		if len(s.src) == 0 {
			return 0, io.EOF
		}
		b := s.src[0]
		s.src = s.src[1:]
		s.pos++
		if b == '\n' {
			s.lineno++
		}

		if b != '\\' {
			return b, nil
		}
		if len(s.src) == 0 {
			return b, nil
		}
		if s.src[0] != '\n' {
			return b, nil
		}
		s.src = s.src[1:]
		s.pos++
		s.lineno++
	}
}

func (s *Source) Peek(n int) ([]byte, error) {
	bs := []byte{}
	for i := 0; len(bs) < n && i < len(s.src); i++ {
		b := s.src[i]
		if b != '\\' {
			bs = append(bs, b)
			continue
		}
		if len(s.src) <= i+1 {
			bs = append(bs, b)
			continue
		}
		if s.src[i+1] != '\n' {
			bs = append(bs, b)
			continue
		}
		i++
	}
	if len(bs) < n {
		return bs, io.EOF
	}
	return bs, nil
}

func (s *Source) Filename() string {
	return s.filename
}

func (s *Source) LineNo() int {
	return s.lineno
}

func (s *Source) Pos() int {
	return s.pos
}

type bufSource struct {
	src *Source
	raw []byte
}

func newBufSource(src *Source) *bufSource {
	return &bufSource{
		src: src,
	}
}

func (s *bufSource) Buf() string {
	return string(s.raw)
}

func (s *bufSource) ReadByte() (byte, error) {
	b, err := s.src.ReadByte()
	if err != nil {
		return 0, err
	}
	s.raw = append(s.raw, b)
	return b, nil
}

func (s *bufSource) Peek(n int) ([]byte, error) {
	return s.src.Peek(n)
}
