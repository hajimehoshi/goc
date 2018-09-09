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

func outputTokens(src string) {
	tks := Tokenize(bytes.NewReader([]byte(src)))
	
	for {
		t, err := tks.NextPPToken()
		if err != nil {
			fmt.Println("error")
			return
		}
		if t.Type == EOF {
			break
		}
		fmt.Println(t.String())
	}
}

func ExampleTokenizeEmpty() {
	outputTokens("")
	// Output:
}

func ExampleTokenizeHash() {
	outputTokens("#")
	// Output:
	// #
}

func ExampleTokenizeHashHash() {
	outputTokens("##")
	// Output:
	// ##
}

func ExampleTokenizeUnknownToken() {
	outputTokens("@@ @@@")
	// Output:
	// @@
	// @@@
}

func ExampleTokenizeBackslash() {
	outputTokens("\\")
	// Output:
	// \
}

func ExampleTokenizeCalc() {
	outputTokens("1+1=2")
	// Output:
	// 1
	// +
	// 1
	// =
	// 2
}

func ExampleTokenizeStrings() {
	outputTokens(`"a""b""c"`)
	// Output:
	// "a"
	// "b"
	// "c"
}

func ExampleTokenizeHelloWorld() {
	outputTokens(`int main() {
  printf("Hello, World!\n");
  return 0;
}`)
	// Output:
	// int
	// main
	// (
	// )
	// {
	// (\n)
	// printf
	// (
	// "Hello, World!\n"
	// )
	// ;
	// (\n)
	// return
	// 0
	// ;
	// (\n)
	// }
}

func ExampleTokenizeNewLines() {
	outputTokens(`foo \
bar`)
	// Output:
	// foo
	// bar
}

func ExampleTokenizeBackslashNewLine() {
	outputTokens(`i\
f ("foo\
bar") el\
se
\
`)
	// Output:
	// if
	// (
	// "foobar"
	// )
	// else
	// (\n)
}

func ExampleTokenizeInc() {
	outputTokens(`c+++++c`)
	// Output:
	// c
	// ++
	// ++
	// +
	// c
}

func ExampleTokenizeLineComment() {
	outputTokens(`int main() { // ABC
  return 0;
} // DEF`)
	// Output:
	// int
	// main
	// (
	// )
	// {
	// (\n)
	// return
	// 0
	// ;
	// (\n)
	// }
}

func ExampleTokenizeBlockComment() {
	outputTokens(`int main() {
  /*
    hi
  */
  return /* hihi */ 0;
}`)
	// Output:
	// int
	// main
	// (
	// )
	// {
	// (\n)
	// (\n)
	// return
	// 0
	// ;
	// (\n)
	// }
}

func ExampleTokenizeComplexComment() {
	outputTokens(`/**/*/*"*/*/*"//*//**/*/`)
	// Output:
	// *
	// *
	// *
	// /
}

func ExampleTokenizeInclude() {
	outputTokens(`#include <abc>
# <abc>
#foo <abc>
abc <abc>
#include "abc"
"abc"`)
	// Output:
	// #
	// include
	// <abc>
	// (\n)
	// #
	// <
	// abc
	// >
	// (\n)
	// #
	// foo
	// <
	// abc
	// >
	// (\n)
	// abc
	// <
	// abc
	// >
	// (\n)
	// #
	// include
	// "abc"
	// (\n)
	// "abc"
}

func ExampleTokenizeIncludeWithBackslash() {
	outputTokens(`#include <ab\c>
#include "ab\c"`)
	// Output:
	// #
	// include
	// <ab\c>
	// (\n)
	// #
	// include
	// "ab\c"
}
