package ast

import (
	"boludolang/token"
	"bytes"
	"strings"
)

// The base Node interface
type Node interface {
	TokenLiteral() string
	String() string
}

// All statement nodes implement this
type Statement interface {
	Node
	statementNode()
}

// All expression nodes implement this
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

type Boolean struct {
	Token token.Token // the token.TRUE or token.FALSE token
	Value bool
}

func (boolExpr *Boolean) expressionNode()      {}
func (boolExpr *Boolean) TokenLiteral() string { return boolExpr.Token.Literal }
func (boolExpr *Boolean) String() string       { return boolExpr.Token.Literal }

type BlockStatement struct {
	Token      token.Token // the token.LBRACE token
	Statements []Statement
}

func (blockStmt *BlockStatement) statementNode()       {}
func (blockStmt *BlockStatement) TokenLiteral() string { return blockStmt.Token.Literal }
func (blockStmt *BlockStatement) String() string {
	var out bytes.Buffer

	for _, stmt := range blockStmt.Statements {
		out.WriteString(stmt.String())
	}

	return out.String()
}

type IfExpression struct {
	Token       token.Token // the token.IF token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ifExpr *IfExpression) expressionNode()      {}
func (ifExpr *IfExpression) TokenLiteral() string { return ifExpr.Token.Literal }
func (ifExpr *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ifExpr.Condition.String())
	out.WriteString(" ")
	out.WriteString(ifExpr.Consequence.String())

	if ifExpr.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ifExpr.Alternative.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token // the token.FUNCTION token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (funcLit *FunctionLiteral) expressionNode()      {}
func (funcLit *FunctionLiteral) TokenLiteral() string { return funcLit.Token.Literal }
func (funcLit *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}

	for _, param := range funcLit.Parameters {
		params = append(params, param.String())
	}

	out.WriteString(funcLit.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(funcLit.Body.String())

	return out.String()
}

type CallExpression struct {
	Token     token.Token // the token.LPAREN token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (callExpr *CallExpression) expressionNode()      {}
func (callExpr *CallExpression) TokenLiteral() string { return callExpr.Token.Literal }
func (callExpr *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}

	for _, arg := range callExpr.Arguments {
		args = append(args, arg.String())
	}

	out.WriteString(callExpr.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type StringLiteral struct {
	Token token.Token // the token.STRING token
	Value string
}

func (strLit *StringLiteral) expressionNode()      {}
func (strLit *StringLiteral) TokenLiteral() string { return strLit.Token.Literal }
func (strLit *StringLiteral) String() string       { return strLit.Token.Literal }

type ArrayLiteral struct {
	Token    token.Token // the token.LBRACKET token
	Elements []Expression
}

func (arrLit *ArrayLiteral) expressionNode()      {}
func (arrLit *ArrayLiteral) TokenLiteral() string { return arrLit.Token.Literal }
func (arrLit *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, elem := range arrLit.Elements {
		elements = append(elements, elem.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type IndexExpression struct {
	Token token.Token // the token.LBRACKET token
	Left  Expression
	Index Expression
}

func (indexExpr *IndexExpression) expressionNode()      {}
func (indexExpr *IndexExpression) TokenLiteral() string { return indexExpr.Token.Literal }
func (indexExpr *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(indexExpr.Left.String())
	out.WriteString("[")
	out.WriteString(indexExpr.Index.String())
	out.WriteString("])")

	return out.String()
}

type HashLiteral struct {
	Token token.Token // the token.LBRACE token
	Pairs map[Expression]Expression
}

func (hashLit *HashLiteral) expressionNode()      {}
func (hashLit *HashLiteral) TokenLiteral() string { return hashLit.Token.Literal }
func (hashLit *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range hashLit.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
