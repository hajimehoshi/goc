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

type stringConcatter struct {
	src PPTokenReader
	buf *Token
}

func (s *stringConcatter) NextPPToken() (*Token, error) {
	if s.buf != nil {
		t := s.buf
		s.buf = nil
		return t, nil
	}

	t, err := s.src.NextPPToken()
	if err != nil {
		return nil, err
	}
	if t.Type != StringLiteral {
		return t, nil
	}

	str := t
	for {
		t, err := s.src.NextPPToken()
		if err != nil {
			return nil, err
		}
		if t.Type != StringLiteral {
			s.buf = t
			return str, nil
		}
		str.Val += t.Val
		if str.Raw == "" {
			str.Raw += t.Raw
		} else {
			str.Raw += " " + t.Raw
		}
	}
}
