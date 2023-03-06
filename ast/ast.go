package ast

import (
    "fmt"
    "monkey/token"
    "strings"
)

type Node interface {
    fmt.Stringer
    Literal() string
}

type Statement interface {
    Node
    statementNode()
}

type Expression interface {
    Node
    expressionNode()
}

// ===========================================   Program   ============================================

type Program struct {
    Statements []Statement
}

func (program *Program) String() string {
    var s = ""
    for i := range program.Statements {
        stmt := program.Statements[i]
        if stmt != nil {
            s += stmt.String()
            s += "\n"
        }
    }
    return s
}

func (program *Program) Literal() string {
    if len(program.Statements) > 0 {
        return program.Statements[0].Literal()
    }
    return ""
}

// ====================================================================================================

// ========================================   LetStatement   ==========================================

type LetStatement struct {
    Token *token.Token
    Name  *Identifier
    Value Expression
}

func (letStmt LetStatement) String() string {
    var builder strings.Builder
    builder.WriteString(letStmt.Literal() + " ")
    builder.WriteString(letStmt.Name.String())
    builder.WriteString(" = ")
    if letStmt.Value == nil {
        builder.WriteString("nil")
    } else {
        builder.WriteString(letStmt.Value.String())
    }
    builder.WriteString("\n")
    return builder.String()
}

func (letStmt LetStatement) Literal() string {
    return letStmt.Token.Literal
}

func (letStmt LetStatement) statementNode() {

}

// =====================================================================================================

// ========================================   ReturnStatement  =========================================

type ReturnStatement struct {
    Token       *token.Token
    ReturnValue Expression
}

func (retStmt *ReturnStatement) statementNode() {

}

func (retStmt *ReturnStatement) String() string {
    return fmt.Sprintf("return %s", retStmt.ReturnValue.String())
}

func (retStmt *ReturnStatement) Literal() string {
    return retStmt.ReturnValue.Literal()
}

// =====================================================================================================

type BlockStatement struct {
    Token      *token.Token
    Statements []Statement
}

func (blockStmt *BlockStatement) String() string {
    var builder strings.Builder
    builder.WriteString("{\n")
    for _, stmt := range blockStmt.Statements {
        builder.WriteString("\t")
        builder.WriteString(stmt.String())
    }
    builder.WriteString("\n}\n")
    return builder.String()
}

func (blockStmt *BlockStatement) Literal() string {
    return blockStmt.Token.Literal
}

func (blockStmt *BlockStatement) statementNode() {

}

// ======================================= ExpressionStatement =========================================

type ExpressionStatement struct {
    Token      *token.Token
    Expression Expression
}

func (exprStmt *ExpressionStatement) String() string {
    return exprStmt.Expression.String()
}

func (exprStmt *ExpressionStatement) Literal() string {
    return exprStmt.Token.Literal
}

func (exprStmt *ExpressionStatement) statementNode() {
}

// =====================================================================================================

// ========================================  PrefixExpression  =========================================

type PrefixExpression struct {
    Token    *token.Token
    Operator string
    Right    Expression
}

func (prefixExpr *PrefixExpression) String() string {
    var builder strings.Builder
    builder.WriteString("(")
    builder.WriteString(prefixExpr.Operator)
    builder.WriteString(prefixExpr.Right.String())
    builder.WriteString(")")
    return builder.String()
}

func (prefixExpr *PrefixExpression) Literal() string {
    return prefixExpr.Token.Literal
}

func (prefixExpr *PrefixExpression) expressionNode() {

}

// =====================================================================================================

// ========================================   InfixExpression  =========================================

type InfixExpression struct {
    Token    *token.Token
    Left     Expression
    Operator string
    Right    Expression
}

func (infixExpr *InfixExpression) String() string {
    var builder strings.Builder
    builder.WriteString("(")
    builder.WriteString(infixExpr.Left.String())
    builder.WriteString(infixExpr.Operator)
    builder.WriteString(infixExpr.Right.String())
    builder.WriteString(")")
    return builder.String()
}

func (infixExpr *InfixExpression) Literal() string {
    return infixExpr.Token.Literal
}

func (infixExpr *InfixExpression) expressionNode() {
}

// =====================================================================================================

// =========================================   IfExpression   ==========================================

type IfExpression struct {
    Token       *token.Token
    Condition   Expression
    Consequence *BlockStatement
    Alternative *BlockStatement
}

func (infixExpr *IfExpression) String() string {
    var builder strings.Builder
    builder.WriteString("if ")
    builder.WriteString(infixExpr.Condition.String())
    builder.WriteString(" ")
    builder.WriteString(infixExpr.Consequence.String())
    if infixExpr.Alternative != nil {
        builder.WriteString("else ")
        builder.WriteString(infixExpr.Alternative.String())
    }
    builder.WriteString("\n")
    return builder.String()
}

func (infixExpr *IfExpression) Literal() string {
    return infixExpr.Token.Literal
}

func (infixExpr *IfExpression) expressionNode() {
}

// =====================================================================================================

// =======================================   CallExpression   ==========================================

type CallExpression struct {
    Token     *token.Token
    Function  Expression
    Arguments []Expression
}

func (callExpr *CallExpression) String() string {
    var builder strings.Builder
    builder.WriteString(callExpr.Function.String())

    builder.WriteString("(")
    var t []string
    for _, args := range callExpr.Arguments {
        t = append(t, args.String())
    }
    builder.WriteString(strings.Join(t, ","))
    builder.WriteString(")")

    return builder.String()
}

func (callExpr *CallExpression) Literal() string {
    return callExpr.Token.Literal
}

func (callExpr *CallExpression) expressionNode() {
}

// =====================================================================================================

// ==========================================   Identifier   ===========================================

type Identifier struct {
    Token *token.Token
    Value string
}

func (id Identifier) expressionNode() {

}

func (id Identifier) String() string {
    return id.Value
}

func (id Identifier) Literal() string {
    return id.Token.Literal
}

// =====================================================================================================

type Integer struct {
    Token *token.Token
    Value int64
}

func (integer *Integer) expressionNode() {

}

func (integer *Integer) String() string {
    return integer.Token.Literal
}

func (integer *Integer) Literal() string {
    return integer.Token.Literal
}

type Boolean struct {
    Token *token.Token
    Value bool
}

func (boolean *Boolean) String() string {
    if boolean.Value {
        return "true"
    }
    return "false"
}

func (boolean *Boolean) Literal() string {
    return boolean.Token.Literal
}

func (boolean *Boolean) expressionNode() {

}

type Function struct {
    Token  *token.Token
    Params []*Identifier
    Body   *BlockStatement
}

func (function *Function) String() string {
    var builder strings.Builder
    builder.WriteString("(")

    var t []string
    for _, param := range function.Params {
        t = append(t, param.String())
    }
    builder.WriteString(strings.Join(t, ","))

    builder.WriteString(") ")
    builder.WriteString(function.Body.String())
    return builder.String()
}

func (function *Function) Literal() string {
    return function.Token.Literal
}

func (function *Function) expressionNode() {
}
