package ast

import (
    "errors"
    "fmt"
    "log"
    token2 "monkey/token"
    "strconv"
)

const (
    _ = iota
    Lowest
    Equals
    LessGreater
    Sum
    Product
    Prefix
    Call
)

var Precedences = map[token2.Type]int{
    token2.Eq:       Equals,
    token2.Ne:       Equals,
    token2.Lt:       LessGreater,
    token2.Gt:       LessGreater,
    token2.Plus:     Sum,
    token2.Minus:    Sum,
    token2.Asterisk: Product,
    token2.Slash:    Product,
    token2.Lparen:   Call,
}

type (
    PrefixExpressionResolver = func() Expression
    InfixExpressionResolver  = func(Expression) Expression
)

type Parser struct {
    lexer               *token2.Lexer
    errors              []error
    currentToken        *token2.Token
    peekToken           *token2.Token
    prefixExprResolvers map[string]PrefixExpressionResolver
    infixExprResolvers  map[string]InfixExpressionResolver
}

func NewParser(l *token2.Lexer) *Parser {
    parser := &Parser{lexer: l}
    parser.nextToken()
    parser.nextToken()

    parser.prefixExprResolvers = make(map[string]PrefixExpressionResolver)
    parser.registerPrefix(token2.Ident, parser.parseIdentifier)
    parser.registerPrefix(token2.Int, parser.parseInteger)
    parser.registerPrefix(token2.Bang, parser.parsePrefixExpression)
    parser.registerPrefix(token2.Minus, parser.parsePrefixExpression)
    parser.registerPrefix(token2.True, parser.parseBoolean)
    parser.registerPrefix(token2.False, parser.parseBoolean)
    parser.registerPrefix(token2.Lparen, parser.parseGroupedExpression)
    parser.registerPrefix(token2.If, parser.parseIfExpression)
    parser.registerPrefix(token2.Function, parser.parseFunction)

    parser.infixExprResolvers = make(map[string]InfixExpressionResolver)
    parser.registerInfix(token2.Plus, parser.parseInfixExpression)
    parser.registerInfix(token2.Minus, parser.parseInfixExpression)
    parser.registerInfix(token2.Asterisk, parser.parseInfixExpression)
    parser.registerInfix(token2.Slash, parser.parseInfixExpression)
    parser.registerInfix(token2.Eq, parser.parseInfixExpression)
    parser.registerInfix(token2.Ne, parser.parseInfixExpression)
    parser.registerInfix(token2.Lt, parser.parseInfixExpression)
    parser.registerInfix(token2.Gt, parser.parseInfixExpression)
    parser.registerInfix(token2.Lparen, parser.parseCallExpression)

    return parser
}

func (parser *Parser) nextToken() {
    parser.currentToken = parser.peekToken
    parser.peekToken = parser.lexer.NextToken()
}

func (parser *Parser) Parse() *Program {
    program := &Program{Statements: make([]Statement, 2)}
    for parser.currentToken.Type != token2.Eof {
        statement := parser.parseStatement()
        program.Statements = append(program.Statements, statement)
        parser.nextToken()
    }
    if len(parser.errors) != 0 {
        log.Fatal(parser.errors)
    }
    return program
}

func (parser *Parser) parseStatement() Statement {
    switch parser.currentToken.Type {
    case token2.Let:
        return parser.parseLetStatement()
    case token2.Return:
        return parser.parseReturnStatement()
    default:
        return parser.parseExpressionStatement()
    }
}

// let v = 1;
func (parser *Parser) parseLetStatement() (letStmt *LetStatement) {
    letStmt = &LetStatement{
        Token: parser.currentToken,
        Name:  nil,
        Value: nil,
    }

    if !parser.assertPeekTokenIs(token2.Ident) {
        return nil
    }

    letStmt.Name = &Identifier{
        Token: parser.currentToken,
        Value: parser.currentToken.Literal,
    }

    if !parser.assertPeekTokenIs(token2.Assign) {
        return nil
    }

    parser.nextToken()
    letStmt.Value = parser.parseExpression(Lowest)

    if !parser.currentTokenIs(token2.Semicolon) {
        parser.nextToken()
    }
    return letStmt
}

// return v;
func (parser *Parser) parseReturnStatement() (stmt *ReturnStatement) {
    stmt = &ReturnStatement{Token: parser.currentToken, ReturnValue: nil}

    parser.nextToken()
    stmt.ReturnValue = parser.parseExpression(Lowest)

    if !parser.currentTokenIs(token2.Semicolon) {
        parser.nextToken()
    }
    return stmt
}

//  x + y;
func (parser *Parser) parseExpressionStatement() (stmt *ExpressionStatement) {
    stmt = &ExpressionStatement{
        Token:      parser.currentToken,
        Expression: parser.parseExpression(Lowest),
    }

    if parser.peekTokenIs(token2.Semicolon) {
        parser.nextToken()
    }
    return stmt
}

func (parser *Parser) parseExpression(precedence int) Expression {
    prefix := parser.prefixExprResolvers[parser.currentToken.Type]
    if prefix == nil {
        parser.error(errors.New(fmt.Sprintf("no prefix for %s found", parser.currentToken.Type)))
        return nil
    }
    left := prefix()

    for precedence < parser.peekPrecedence() {
        t := parser.peekToken.Type
        infix := parser.infixExprResolvers[t]
        if infix == nil {
            return left
        }
        parser.nextToken()
        left = infix(left)
    }
    return left
}

func (parser *Parser) parseGroupedExpression() Expression {
    parser.nextToken()
    groupedExpr := parser.parseExpression(Lowest)
    if !parser.assertPeekTokenIs(token2.Rparen) {
        return nil
    }
    return groupedExpr
}

func (parser *Parser) parsePrefixExpression() Expression {
    prefixExpr := &PrefixExpression{
        Token:    parser.currentToken,
        Operator: parser.currentToken.Literal,
        Right:    nil,
    }
    parser.nextToken()
    prefixExpr.Right = parser.parseExpression(Prefix)
    return prefixExpr
}

func (parser *Parser) parseInfixExpression(leftExpr Expression) Expression {
    infixExpr := &InfixExpression{
        Token:    parser.currentToken,
        Left:     leftExpr,
        Operator: parser.currentToken.Literal,
        Right:    nil,
    }

    precedence := parser.currentPrecedence()
    parser.nextToken()
    infixExpr.Right = parser.parseExpression(precedence)
    return infixExpr
}

func (parser *Parser) parseIfExpression() Expression {
    ifExpr := &IfExpression{
        Token:       nil,
        Condition:   nil,
        Consequence: nil,
        Alternative: nil,
    }

    if !parser.assertPeekTokenIs(token2.Lparen) {
        return nil
    }

    parser.nextToken()
    ifExpr.Condition = parser.parseExpression(Lowest)

    if !parser.assertPeekTokenIs(token2.Rparen) {
        return nil
    }

    if !parser.assertPeekTokenIs(token2.Lbrace) {
        return nil
    }

    ifExpr.Consequence = parser.parseBlockStatement()

    if parser.peekTokenIs(token2.Else) {
        parser.nextToken()
        if !parser.assertPeekTokenIs(token2.Lbrace) {
            return nil
        }
        ifExpr.Alternative = parser.parseBlockStatement()
    }

    return ifExpr
}

func (parser *Parser) parseFunction() Expression {
    function := Function{
        Token:  parser.currentToken,
        Params: []*Identifier{},
        Body:   nil,
    }
    if !parser.assertPeekTokenIs(token2.Lparen) {
        return nil
    }

    // parse params
    if !parser.peekTokenIs(token2.Rparen) {
        for {
            parser.assertPeekTokenIs(token2.Ident)
            function.Params = append(function.Params, &Identifier{
                Token: parser.currentToken,
                Value: parser.currentToken.Literal,
            })

            parser.nextToken()
            if parser.currentTokenIs(token2.Comma) {
                continue
            } else if parser.currentTokenIs(token2.Rparen) {
                break
            }
        }
    }

    if !parser.assertPeekTokenIs(token2.Lbrace) {
        return nil
    }

    function.Body = parser.parseBlockStatement()
    return &function
}

func (parser *Parser) parseCallExpression(function Expression) Expression {
    callExpr := &CallExpression{
        Token:     parser.currentToken,
        Function:  function,
        Arguments: []Expression{},
    }

    parser.nextToken()
    // parse params
    if !parser.currentTokenIs(token2.Rparen) {
        callExpr.Arguments = append(callExpr.Arguments, parser.parseExpression(Lowest))
        for parser.peekTokenIs(token2.Comma) {
            parser.nextToken()
            parser.nextToken()
            callExpr.Arguments = append(callExpr.Arguments, parser.parseExpression(Lowest))
        }
        if !parser.assertPeekTokenIs(token2.Rparen) {
            return nil
        }
    }

    return callExpr
}

func (parser *Parser) parseBlockStatement() *BlockStatement {
    blockStmt := BlockStatement{
        Token:      parser.currentToken,
        Statements: []Statement{},
    }

    parser.nextToken()
    for !parser.currentTokenIs(token2.Rbrace) && !parser.currentTokenIs(token2.Eof) {
        stmt := parser.parseStatement()
        if stmt != nil {
            blockStmt.Statements = append(blockStmt.Statements, stmt)
        }
        parser.nextToken()
    }

    return &blockStmt
}

// varName
func (parser *Parser) parseIdentifier() Expression {
    return &Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}
}

// 5
func (parser *Parser) parseInteger() Expression {
    integer := &Integer{
        Token: parser.currentToken,
        Value: 0,
    }

    value, err := strconv.ParseInt(parser.currentToken.Literal, 0, 64)
    if err != nil {
        parser.error(errors.New(fmt.Sprintf("could not parse %q as integer", integer.Token)))
    }

    integer.Value = value
    return integer
}

func (parser *Parser) parseBoolean() Expression {
    return &Boolean{
        Token: parser.currentToken,
        Value: parser.currentTokenIs(token2.True),
    }
}

func (parser *Parser) assertPeekTokenIs(tokenType token2.Type) bool {
    if parser.peekTokenIs(tokenType) {
        parser.nextToken()
        return true
    }
    parser.error(errors.New(fmt.Sprintf("expected next token is %s, but got %s", tokenType, parser.peekToken.Type)))
    return false
}

func (parser *Parser) currentTokenIs(tokenType token2.Type) bool {
    return parser.currentToken.Type == tokenType
}

func (parser *Parser) peekTokenIs(tokenType token2.Type) bool {
    return parser.peekToken.Type == tokenType
}

func (parser *Parser) registerPrefix(prefix string, resolver PrefixExpressionResolver) {
    parser.prefixExprResolvers[prefix] = resolver
}

func (parser *Parser) registerInfix(prefix string, resolver InfixExpressionResolver) {
    parser.infixExprResolvers[prefix] = resolver
}

func (parser *Parser) error(err error) {
    parser.errors = append(parser.errors, err)
}

func (parser *Parser) currentPrecedence() int {
    if precedence, ok := Precedences[parser.currentToken.Type]; ok {
        return precedence
    }
    return Lowest
}

func (parser *Parser) peekPrecedence() int {
    if precedence, ok := Precedences[parser.peekToken.Type]; ok {
        return precedence
    }
    return Lowest
}
