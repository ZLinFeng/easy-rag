package core

import (
	"fmt"
	"testing"
)

func Test_Tokenizer(t *testing.T) {
	l := &Lexer{input: "$.name.lastName[0].age.children[*]"}
	for {
		tok, err := l.NextToken()
		if err != nil {
			fmt.Println(err)
		}
		if tok.Type == EOF {
			break
		}
		fmt.Printf("%+v\n", tok)
	}
}
