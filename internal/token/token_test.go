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

	"github.com/hajimehoshi/goc/internal/preprocess"
	. "github.com/hajimehoshi/goc/internal/token"
)

func outputTokens(path string, srcs map[string]string) {
	files := map[string]preprocess.PPTokenReader{}
	for path, src := range srcs {
		files[path] = preprocess.Tokenize(bytes.NewReader([]byte(src)))
	}

	pptokens, err := preprocess.Preprocess(path, files)
	if err != nil {
		fmt.Println("error")
		return
	}

	tokens := []*Token{}
	for _, pt := range pptokens {
		t, err := FromPPToken(pt)
		if err != nil {
			fmt.Println("error")
			return
		}
		tokens = append(tokens, t)
	}
	
	for _, t := range tokens {
		fmt.Println(t)
	}
}

func ExampleEmpty() {
	outputTokens("main.c", map[string]string{
		"main.c": `#`,
	})
	// Output:
}

func ExampleFunc() {
	outputTokens("main.c", map[string]string{
		"main.c": `int main() {
  printf("Hello, World!\n");
  return 0;
}`,
	})
	// Output:
	// int
	// identifier: main
	// (
	// )
	// {
	// identifier: printf
	// (
	// string: "Hello, World!\n"
	// )
	// ;
	// return
	// integer: 0 (int)
	// ;
	// }
}

