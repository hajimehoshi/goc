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
)

func outputPreprocessedTokens(path string, srcs map[string]string) {
	files := map[string]PPTokenReader{}
	for path, src := range srcs {
		files[path] = Tokenize(bytes.NewReader([]byte(src)))
	}

	pptokens, err := Preprocess(path, files)
	if err != nil {
		fmt.Println("error")
		return
	}
	for {
		t, err := pptokens.NextPPToken()
		if err != nil {
			fmt.Println("error")
			break
		}
		if t.Type == EOF {
			break
		}
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
	// foo
	// bar
	// baz
	// qux
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
	// 1
	// )
	// BAZ
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
	// 1
	// )
	// (
	// 2
	// +
	// 1
	// +
	// 2
	// )
	// (
	// 3
	// +
	// (
	// 1
	// ,
	// 2
	// )
	// +
	// 3
	// )
	// BAZ
}

func ExampleUndef() {
	outputPreprocessedTokens("main.c", map[string]string{
		"main.c": `#define FOO 1
FOO
#undef FOO
FOO`,
	})
	// Output:
	// 1
	// FOO
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

func ExampleDefineFunctionLike() {
	outputPreprocessedTokens("main.c", map[string]string{
		"main.c": `#define foo(x) x
foo(1)
foo(.11dd)
foo("")
foo(\)
foo(<<<<<)`,
	})
	// Output:
	// 1
	// .11dd
	// ""
	// \
	// <<
	// <<
	// <
}

func ExampleDefineFunctionLikeNotEnded() {
	outputPreprocessedTokens("main.c", map[string]string{
		"main.c": `#define foo(x) x
foo(1`,
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
	// c
	// )
	// +
	// (
	// (
	// (
	// b
	// )
	// +
	// (
	// a
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
	// a
}

func ExampleDefineRescanRecursive2() {
	outputPreprocessedTokens("main.c", map[string]string{
		"main.c": `#define a a b
a`,
	})
	// Output:
	// a
	// b
}

func ExampleDefineRescanRecursive3() {
	outputPreprocessedTokens("main.c", map[string]string{
		"main.c": `#define b a
#define a b
a
#undef b
a
#define b c
a`,
	})
	// Output:
	// a
	// b
	// c
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
	// x
	// y
	// long
	// z
}

func ExampleHash() {
	outputPreprocessedTokens("main.c", map[string]string{
		"main.c": `#define str(x) #x
str(ddd    eeeee)
str(111ddd)
str(  <<<<< @@  )
str(\n)
str("")
str("\"")
str("\n")
str(str(a))`,
	})
	// Output:
	// "ddd eeeee"
	// "111ddd"
	// "<<<<< @@"
	// "\n"
	// "\"\""
	// "\"\\\"\""
	// "\"\\n\""
	// "str(a)"
}

func ExampleHashError() {
	outputPreprocessedTokens("main.c", map[string]string{
		"main.c": `#define str(x) #x
str(\) // Syntax error`,
	})
	// Output:
	// error
}

func ExampleHashPositionError() {
	outputPreprocessedTokens("main.c", map[string]string{
		"main.c": `#define str(x) #x #`,
	})
	// Output:
	// error
}

func ExampleHashPositionError2() {
	outputPreprocessedTokens("main.c", map[string]string{
		"main.c": `#define str(x) # #x`,
	})
	// Output:
	// error
}

func ExampleHashPositionError3() {
	outputPreprocessedTokens("main.c", map[string]string{
		"main.c": `#define str(x) #y #x`,
	})
	// Output:
	// error
}
