package token

import "testing"

var input = `
###
`

var s = `1`

func Test(t *testing.T) {
    lexer := NewLexer(s)
    for lexer.HasNext() {
        println(lexer.NextToken().String())
    }
}
