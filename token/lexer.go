package token

type Lexer struct {
    input     string
    pos       int
    toReadPos int
    char      byte
}

func NewLexer(input string) *Lexer {
    lexer := &Lexer{input: input}
    lexer.readChar()

    return lexer
}

func (lexer *Lexer) HasNext() bool {
    return lexer.pos != len(lexer.input)
}

func (lexer *Lexer) NextToken() *Token {
    lexer.skipWhitespaces()

    var token *Token
    currentChar := lexer.char
    switch currentChar {
    case '=':
        if lexer.peekChar() == '=' {
            lexer.readChar()
            token = newToken(Eq, "==")
        } else {
            token = newToken(Assign, "=")
        }
    case ';':
        token = newToken(Semicolon, ";")
    case ',':
        token = newToken(Comma, ",")
    case '(':
        token = newToken(Lparen, "(")
    case ')':
        token = newToken(Rparen, ")")
    case '{':
        token = newToken(Lbrace, "{")
    case '}':
        token = newToken(Rbrace, "}")
    case '+':
        token = newToken(Plus, "+")
    case '-':
        token = newToken(Minus, "-")
    case '!':
        if lexer.peekChar() == '=' {
            lexer.readChar()
            token = newToken(Ne, "!=")
        } else {
            token = newToken(Bang, "!")
        }
    case '*':
        token = newToken(Asterisk, "*")
    case '/':
        token = newToken(Slash, "/")
    case '<':
        token = newToken(Lt, "<")
    case '>':
        token = newToken(Gt, ">")
    case 0:
        token = newToken(Eof, "")
    default:
        if isLetter(currentChar) {
            word := lexer.readWord()
            token = newToken(Ident, word)
            if ok, tokenType := isKeyword(word); ok {
                token.Type = tokenType
            }
            return token
        } else if isDigit(currentChar) {
            number := lexer.readNumber()
            token = newToken(Int, number)
            return token
        } else {
            token = newToken(Illegal, "")
        }
    }

    lexer.readChar()
    return token
}

func (lexer *Lexer) readChar() byte {
    if lexer.toReadPos < len(lexer.input) {
        lexer.char = lexer.input[lexer.toReadPos]
        lexer.pos = lexer.toReadPos
        lexer.toReadPos++
    } else {
        lexer.char = 0
        lexer.pos = lexer.toReadPos
    }
    return lexer.char
}

func (lexer *Lexer) peekChar() byte {
    if lexer.HasNext() {
        return lexer.input[lexer.toReadPos]
    } else {
        return 0
    }
}

func (lexer *Lexer) readWord() string {
    pos := lexer.pos
    for isLetter(lexer.char) {
        lexer.readChar()
    }
    return lexer.input[pos:lexer.pos]
}

func (lexer *Lexer) readNumber() string {
    pos := lexer.pos
    for isDigit(lexer.char) {
        lexer.readChar()
    }
    return lexer.input[pos:lexer.pos]
}

func (lexer *Lexer) skipWhitespaces() {
    for lexer.char == ' ' || lexer.char == '\t' || lexer.char == '\r' || lexer.char == '\n' {
        lexer.readChar()
    }
}
