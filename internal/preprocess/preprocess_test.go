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

func outputPreprocessedTokens(path string, srcs map[string]string) {
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
	outputPreprocessedTokens("main.c", map[string]string{
		"main.c": `#`,
	})
	// Output:
}

func ExampleIncludeSimple() {
	outputPreprocessedTokens("main.c", map[string]string{
		"main.c": `#include <stdio.h>
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
	outputPreprocessedTokens("main.c", map[string]string{
		"main.c":  `#include <stdio.h>`,
		"stdio.h": `#include <main.c>`,
	})
	// Output:
	// error
}

func ExampleDefineObjLike() {
	outputPreprocessedTokens("main.c", map[string]string{
		"main.c": `#define FOO
#define BAR (1)
FOO
BAR
BAZ`,
	})
	// Output:
	// (
	// number: 1 (int)
	// )
	// ident: BAZ
}

func ExampleDefineFuncLike() {
	outputPreprocessedTokens("main.c", map[string]string{
		"main.c": `#define FOO
#define BAR(X, Y) (Y + X + Y)
FOO(1)
BAR(1, 2)
BAR((1, 2), 3)
BAZ`,
	})
	// Output:
	// (
	// number: 1 (int)
	// )
	// (
	// number: 2 (int)
	// +
	// number: 1 (int)
	// +
	// number: 2 (int)
	// )
	// (
	// number: 3 (int)
	// +
	// (
	// number: 1 (int)
	// ,
	// number: 2 (int)
	// )
	// +
	// number: 3 (int)
	// )
	// ident: BAZ
}

func ExampleUndef() {
	outputPreprocessedTokens("main.c", map[string]string{
		"main.c": `#define FOO 1
FOO
#undef FOO
FOO`,
	})
	// Output:
	// number: 1 (int)
	// ident: FOO
}

func ExampleUndefIgnored() {
	outputPreprocessedTokens("main.c", map[string]string{
		"main.c": `#define FOO 1
#undef BAR`,
	})
	// Output:
}

func ExampleUndefError() {
	outputPreprocessedTokens("main.c", map[string]string{
		"main.c": `#define FOO 1
#undef FOO BAR`,
	})
	// Output:
	// error
}

func ExampleDefineRescan() {
	// 0. plus(plus(a, b), c)
	// 1. add(c, plus(a, b))
	// 2. ((c) + (plus(a, b)))
	// 3. ((c) + (add(b, a)))
	// 4. ((c) + (((b) + (a)))
	outputPreprocessedTokens("main.c", map[string]string{
		"main.c": `#define plus(x, y) add(y, x)
#define add(x, y) ((x)+(y))
plus(plus(a, b), c)
`,
	})
	// Output:
	// (
	// (
	// ident: c
	// )
	// +
	// (
	// (
	// (
	// ident: b
	// )
	// +
	// (
	// ident: a
	// )
	// )
	// )
	// )
}

func ExampleDefineRescanRecursive() {
	outputPreprocessedTokens("main.c", map[string]string{
		"main.c": `#define a b
#define b a
a`,
	})
	// Output:
	// ident: a
}

func ExampleDefineRescanRecursive2() {
	outputPreprocessedTokens("main.c", map[string]string{
		"main.c": `#define a a b
a`,
	})
	// Output:
	// ident: a
	// ident: b
}

func ExampleDefineKeyword() {
	outputPreprocessedTokens("main.c", map[string]string{
		"main.c": `#define char unsigned char
#define foo(long) long
char x
foo(y)
long z`,
	})
	// Output:
	// unsigned
	// char
	// ident: x
	// ident: y
	// long
	// ident: z
}
