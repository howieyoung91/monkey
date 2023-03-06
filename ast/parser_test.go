package ast

import (
    "monkey/token"
    "testing"
)

var input = `
daw+a+add(1,2,3);
`
var i = `a;`

func Test(t *testing.T) {
    parser := NewParser(token.NewLexer(input))
    program := parser.Parse()
    println(program.String())
}
