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

package lex

import (
	"fmt"
	"io"
)

type BytePeeker interface {
	Peek(int) ([]byte, error)
}

type ByteReadPeeker interface {
	io.ByteReader
	BytePeeker
}

func shouldPeekByte(src BytePeeker) (byte, error) {
	bs, err := shouldPeek(src, 1)
	if err != nil {
		return 0, err
	}
	return bs[0], nil
}

func shouldPeek(src BytePeeker, num int) ([]byte, error) {
	bs, err := src.Peek(num)
	if err != nil && err != io.EOF {
		return nil, err
	}
	if len(bs) < num {
		return nil, fmt.Errorf("lex: unexpected EOF")
	}
	return bs, nil
}

func shouldReadByte(src io.ByteReader) (byte, error) {
	b, err := src.ReadByte()
	if err != nil {
		if err == io.EOF {
			return 0, fmt.Errorf("lex: unexpected EOF")
		}
		return 0, err
	}
	return b, nil
}

func shouldRead(src io.ByteReader, expected byte) error {
	b, err := src.ReadByte()
	if err != nil {
		if err == io.EOF {
			return fmt.Errorf("lex: unexpected EOF")
		}
		return err
	}
	if b != expected {
		return fmt.Errorf("lex: expected %q but %q", expected, b)
	}

	return nil
}

func mustDiscard(src io.ByteReader, n int) {
	for i := 0; i < n; i++ {
		if _, err := src.ReadByte(); err != nil {
			panic("not reached: " + err.Error())
		}
	}
}
