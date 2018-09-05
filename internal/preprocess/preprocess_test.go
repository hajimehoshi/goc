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

package preprocess_test

import (
	"bytes"
	"fmt"

	. "github.com/hajimehoshi/goc/internal/preprocess"
	"github.com/hajimehoshi/goc/internal/token"
	"github.com/hajimehoshi/goc/internal/tokenize"
)

func outputTokens(path string, srcs map[string]string) {
	files := map[string][]*token.Token{}
	for path, src := range srcs {
		ts, err := tokenize.Tokenize(bytes.NewReader([]byte(src)))
		if err != nil {
			panic("not reached")
		}
		files[path] = ts
	}

	tokens, err := Preprocess(path, files)
	if err != nil {
		fmt.Println("error")
		return
	}
	tokens = tokenize.FinishTokenize(tokens)
	for _, t := range tokens {
		fmt.Println(t)
	}
}

func ExampleEmpty() {
	outputTokens("main.c", map[string]string{
		"main.c":  `#`,
	})
	// Output:
}

func ExampleIncludeSimple() {
	outputTokens("main.c", map[string]string{
		"main.c":  `#include <stdio.h>
baz qux`,
		"stdio.h": `foo bar`,
	})
	// Output:
	// ident: foo
	// ident: bar
	// ident: baz
	// ident: qux
}

func ExampleIncludeRecursive() {
	outputTokens("main.c", map[string]string{
		"main.c":  `#include <stdio.h>`,
		"stdio.h": `#include <main.c>`,
	})
	// Output:
	// error
}
