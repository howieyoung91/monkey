package token

const (
	Illegal   = "ILLEGAL"
	Eof       = "EOF"
	Ident     = "IDENT"
	Int       = "INT"
	Assign    = "="
	Plus      = "+"
	Minus     = "-"
	Bang      = "!"
	Asterisk  = "*"
	Slash     = "/"
	Lt        = "<"
	Gt        = ">"
	Eq        = "=="
	Ne        = "!="
	Comma     = ","
	Semicolon = ";"
	Lparen    = "("
	Rparen    = ")"
	Lbrace    = "{"
	Rbrace    = "}"
	Function  = "FUNCTION"
	Let       = "LET"
	If        = "IF"
	Else      = "ELSE"
	True      = "TRUE"
	False     = "FALSE"
	Return    = "RETURN"
)

var KEYWORDS = map[string]Type{
	"fn":     Function,
	"let":    Let,
	"if":     If,
	"else":   Else,
	"true":   True,
	"false":  False,
	"return": Return,
}

type Type = string

type Token struct {
	Type    Type
	Literal string
}

func (token *Token) String() string {
	return "<" + token.Type + ", " + token.Literal + ">"
}

func newToken(tokenType Type, literal string) *Token {
	return &Token{
		Type:    tokenType,
		Literal: literal,
	}
}

func isKeyword(word string) (ok bool, tokenType Type) {
	tokenType, ok = KEYWORDS[word]
	return
}

func isLetter(char byte) bool {
	return ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z')
}

func isDigit(value byte) bool {
	return '0' <= value && value <= '9'
}
