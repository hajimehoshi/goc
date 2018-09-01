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

package ioutil

import (
	"bufio"
	"fmt"
	"io"
)

func ShouldPeekByte(src *bufio.Reader) (byte, error) {
	bs, err := ShouldPeek(src, 1)
	if err != nil {
		return 0, err
	}
	return bs[0], nil
}

func ShouldPeek(src *bufio.Reader, num int) ([]byte, error) {
	bs, err := src.Peek(num)
	if err != nil && err != io.EOF {
		return nil, err
	}
	if len(bs) < num {
		return nil, fmt.Errorf("literal: unexpected EOF")
	}
	return bs, nil
}

func ShouldReadByte(src *bufio.Reader) (byte, error) {
	b, err := src.ReadByte()
	if err != nil {
		if err == io.EOF {
			return 0, fmt.Errorf("literal: unexpected EOF")
		}
		return 0, err
	}
	return b, nil
}

func ShouldRead(src *bufio.Reader, expected byte) error {
	b, err := src.ReadByte()
	if err != nil {
		if err == io.EOF {
			return fmt.Errorf("literal: unexpected EOF")
		}
		return err
	}
	if b != expected {
		return fmt.Errorf("literal: expected %q but %q", expected, b)
	}

	return nil
}
