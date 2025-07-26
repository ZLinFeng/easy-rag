package core

import (
	"fmt"
	"testing"
)

func Test_Tokenizer(t *testing.T) {
	l := &JsonLexer{input: "$.*.age.children[12]"}
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

func TestJsonParser_Parse(t *testing.T) {
	p := NewParser("$.*.age.children[12]")
	if ast, err := p.Parse(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(ast)
	}
}
