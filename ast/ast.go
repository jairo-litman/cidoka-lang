package ast

import (
	"boludolang/token"
	"bytes"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (prog *Program) TokenLiteral() string {
	if len(prog.Statements) > 0 {
		return prog.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (prog *Program) String() string {
	var out bytes.Buffer

	for _, stmt := range prog.Statements {
		out.WriteString(stmt.String())
	}

	return out.String()
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (ident *Identifier) expressionNode()      {}
func (ident *Identifier) TokenLiteral() string { return ident.Token.Literal }
func (ident *Identifier) String() string       { return ident.Value }

type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (letStmt *LetStatement) statementNode()       {}
func (letStmt *LetStatement) TokenLiteral() string { return letStmt.Token.Literal }
func (letStmt *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(letStmt.TokenLiteral() + " ")
	out.WriteString(letStmt.Name.String())
	out.WriteString(" = ")

	if letStmt.Value != nil {
		out.WriteString(letStmt.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type ReturnStatement struct {
	Token       token.Token // the token.RETURN token
	ReturnValue Expression
}

func (returnStmt *ReturnStatement) statementNode()       {}
func (returnStmt *ReturnStatement) TokenLiteral() string { return returnStmt.Token.Literal }
func (returnStmt *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(returnStmt.TokenLiteral() + " ")

	if returnStmt.ReturnValue != nil {
		out.WriteString(returnStmt.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (exprStmt *ExpressionStatement) statementNode()       {}
func (exprStmt *ExpressionStatement) TokenLiteral() string { return exprStmt.Token.Literal }
func (exprStmt *ExpressionStatement) String() string {
	if exprStmt.Expression != nil {
		return exprStmt.Expression.String()
	} else {
		return ""
	}
}

type IntegerLiteral struct {
	Token token.Token // the token.INT token
	Value int64
}

func (intLit *IntegerLiteral) expressionNode()      {}
func (intLit *IntegerLiteral) TokenLiteral() string { return intLit.Token.Literal }
func (intLit *IntegerLiteral) String() string       { return intLit.Token.Literal }

type PrefixExpression struct {
	Token    token.Token // the prefix token, e.g. !
	Operator string
	Right    Expression
}

func (prefixExpr *PrefixExpression) expressionNode()      {}
func (prefixExpr *PrefixExpression) TokenLiteral() string { return prefixExpr.Token.Literal }
func (prefixExpr *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(prefixExpr.Operator)
	out.WriteString(prefixExpr.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token // the infix token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (infixExpr *InfixExpression) expressionNode()      {}
func (infixExpr *InfixExpression) TokenLiteral() string { return infixExpr.Token.Literal }
func (infixExpr *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(infixExpr.Left.String())
	out.WriteString(" " + infixExpr.Operator + " ")
	out.WriteString(infixExpr.Right.String())
	out.WriteString(")")

	return out.String()
}
