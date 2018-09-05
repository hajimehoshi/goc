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

package tokenize_test

import (
	"fmt"

	. "github.com/hajimehoshi/goc/internal/tokenize"
)

func outputTokensFS(srcs map[string]string, path string) {
	tokens, err := Tokenize(&mockFileSystem{
		srcs: srcs,
	}, path, true)
	if err != nil {
		fmt.Println("error")
		return
	}
	for _, t := range tokens {
		fmt.Println(t)
	}
}

func ExampleTokenizeIncludeSimple() {
	outputTokensFS(map[string]string{
		"main.c":  `#include <stdio.h>
baz qux`,
		"stdio.h": `foo bar`,
	}, "main.c")
	// Output:
	// ident: foo
	// ident: bar
	// ident: baz
	// ident: qux
}

func ExampleTokenizeIncludeRecursive() {
	outputTokensFS(map[string]string{
		"main.c":  `#include <stdio.h>`,
		"stdio.h": `#include <main.c>`,
	}, "main.c")
	// Output:
	// error
}
