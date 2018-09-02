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

package token_test

import (
	"bytes"
	"fmt"

	. "github.com/hajimehoshi/goc/internal/token"
)

func outputTokens(src string) {
	tokens, err := Tokenize(bytes.NewReader([]byte(src)))
	if err != nil {
		fmt.Println("error")
		return
	}
	for _, t := range tokens {
		fmt.Println(t)
	}
}

func ExampleTokenizeEmpty() {
	outputTokens("")
	// Output:
}

func ExampleTokenizeNumberSign() {
	outputTokens("#")
	// Output:
	// #
}

func ExampleTokenizeBackslash() {
	outputTokens("\\")
	// Output:
	// error
}

func ExampleTokenizeCalc() {
	outputTokens("1+1=2")
	// Output:
	// number: 1 (int)
	// +
	// number: 1 (int)
	// =
	// number: 2 (int)
}

func ExampleTokenizeStrings() {
	outputTokens(`"a""b""c"`)
	// Output:
	// string: "abc"
}

func ExampleTokenizeHelloWorld() {
	outputTokens(`int main() {
  printf("Hello, World!\n");
  return 0;
}`)
	// Output:
	// int
	// ident: main
	// (
	// )
	// {
	// ident: printf
	// (
	// string: "Hello, World!\n"
	// )
	// ;
	// return
	// number: 0 (int)
	// ;
	// }
}

func ExampleTokenizeNewLines() {
	outputTokens(`foo \
bar`)
	// Output:
	// ident: foo
	// ident: bar
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
	// string: "foobar"
	// )
	// else
}

func ExampleTokenizeInc() {
	outputTokens(`c+++++c`)
	// Output:
	// ident: c
	// ++
	// ++
	// +
	// ident: c
}

func ExampleTokenizeLineComment() {
	outputTokens(`int main() { // ABC
  return 0;
} // DEF`)
	// Output:
	// int
	// ident: main
	// (
	// )
	// {
	// return
	// number: 0 (int)
	// ;
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
	// ident: main
	// (
	// )
	// {
	// return
	// number: 0 (int)
	// ;
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
