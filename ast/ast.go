package ast

import (
	"bytes"
	"cidoka/token"
	"fmt"
	"strings"
)

// The base Node interface
type Node interface {
	TokenLiteral() string // returns the literal value of the token
	String() string       // returns a string representation of the node
}

// All statement nodes implement this
type Statement interface {
	Node
	statementNode() // dummy method to distinguish statements from expressions
}

// All expression nodes implement this
type Expression interface {
	Node
	expressionNode() // dummy method to distinguish expressions from statements
}

// Program is the root node of every AST
type Program struct {
	Statements []Statement // a slice of statements
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

// ----------------------------------------------------------------------------
// 									Statements
// ----------------------------------------------------------------------------

// A let statement, e.g. let x = 5;
type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier // name of the variable
	Value Expression  // expression that evaluates to the value of the variable
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

// An assign statement, e.g. x = 5;
type AssignStatement struct {
	Token token.Token // token.ASSIGN
	Name  *Identifier // name of the variable
	Value Expression  // expression that evaluates to the value of the variable
}

func (assignStmt *AssignStatement) statementNode()       {}
func (assignStmt *AssignStatement) TokenLiteral() string { return assignStmt.Token.Literal }
func (assignStmt *AssignStatement) String() string {
	var out bytes.Buffer

	out.WriteString(assignStmt.Name.String())
	out.WriteString(" = ")

	if assignStmt.Value != nil {
		out.WriteString(assignStmt.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

// A return statement, e.g. return 5;
type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression  // expression that evaluates to the value to return
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

// An expression statement, e.g. 5 + 5;
type ExpressionStatement struct {
	Token      token.Token // first token of the expression
	Expression Expression  // the expression
}

func (exprStmt *ExpressionStatement) statementNode()       {}
func (exprStmt *ExpressionStatement) TokenLiteral() string { return exprStmt.Token.Literal }
func (exprStmt *ExpressionStatement) String() string {
	if exprStmt.Expression != nil {
		return exprStmt.Expression.String()
	}
	return ""
}

// A block statement, e.g. { let x = 5; let y = 10; }
type BlockStatement struct {
	Token      token.Token // token.LBRACE '{'
	Statements []Statement // a slice of statements that make up the block
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

// A for loop statement, e.g. for (let i = 0; i < 10; i = i + 1) { ... }
type ForLoopStatement struct {
	Token       token.Token     // token.FOR
	Initializer Statement       // statement that initializes the loop e.g. let i = 0 // or nil
	Condition   Expression      // expression that evaluates to the condition of the loop e.g. i < 10 // or nil
	Update      Statement       // statement that updates the loop e.g. i = i + 1 // or nil
	Body        *BlockStatement // block statement that makes up the body of the loop
}

func (forLoop *ForLoopStatement) statementNode()       {}
func (forLoop *ForLoopStatement) TokenLiteral() string { return forLoop.Token.Literal }
func (forLoop *ForLoopStatement) String() string {
	var out bytes.Buffer

	out.WriteString("for")
	out.WriteString(" (")

	if forLoop.Initializer != nil {
		out.WriteString(forLoop.Initializer.String())
		out.WriteString(" ")
	} else {
		out.WriteString("; ")
	}

	if forLoop.Condition != nil {
		out.WriteString(forLoop.Condition.String())
	}
	out.WriteString("; ")

	if forLoop.Update != nil {
		out.WriteString(forLoop.Update.String())
	}
	out.WriteString(") ")

	out.WriteString("{ ")
	out.WriteString(forLoop.Body.String())
	out.WriteString(" }")

	return out.String()
}

// A break statement, e.g. break;
type BreakStatement struct {
	Token token.Token // token.BREAK
}

func (breakStmt *BreakStatement) statementNode()       {}
func (breakStmt *BreakStatement) TokenLiteral() string { return breakStmt.Token.Literal }
func (breakStmt *BreakStatement) String() string       { return breakStmt.TokenLiteral() }

// ----------------------------------------------------------------------------
// 								Expressions
// ----------------------------------------------------------------------------

// An identifier expression, e.g. x
type Identifier struct {
	Token token.Token // token.IDENT
	Value string      // name of the identifier // should be the same as the token literal
}

func (ident *Identifier) expressionNode()      {}
func (ident *Identifier) TokenLiteral() string { return ident.Token.Literal }
func (ident *Identifier) String() string       { return ident.Value }

// An integer literal expression, e.g. 5
type IntegerLiteral struct {
	Token token.Token // token.INT
	Value int64       // value of the integer
}

func (intLit *IntegerLiteral) expressionNode()      {}
func (intLit *IntegerLiteral) TokenLiteral() string { return intLit.Token.Literal }
func (intLit *IntegerLiteral) String() string       { return intLit.Token.Literal }

// A float literal expression, e.g. 5.5
type FloatLiteral struct {
	Token token.Token // token.FLOAT
	Value float64     // value of the float
}

func (floatLit *FloatLiteral) expressionNode()      {}
func (floatLit *FloatLiteral) TokenLiteral() string { return floatLit.Token.Literal }
func (floatLit *FloatLiteral) String() string       { return floatLit.Token.Literal }

// A boolean expression, e.g. true or false
type Boolean struct {
	Token token.Token // token.TRUE or token.FALSE
	Value bool        // value of the boolean
}

func (boolExpr *Boolean) expressionNode()      {}
func (boolExpr *Boolean) TokenLiteral() string { return boolExpr.Token.Literal }
func (boolExpr *Boolean) String() string       { return boolExpr.Token.Literal }

// A string literal expression, e.g. "foobar"
type StringLiteral struct {
	Token token.Token // token.STRING
	Value string      // value of the string
}

func (strLit *StringLiteral) expressionNode()      {}
func (strLit *StringLiteral) TokenLiteral() string { return strLit.Token.Literal }
func (strLit *StringLiteral) String() string       { return strLit.Token.Literal }

// A prefix expression, e.g. !5 or -15
type PrefixExpression struct {
	Token    token.Token // the prefix token, e.g. token.BANG or token.MINUS
	Operator string      // operator to be applied to the right expression
	Right    Expression  // right expression to be evaluated
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

// An infix expression, e.g. 5 + 5
type InfixExpression struct {
	Token    token.Token // infix token, e.g. token.PLUS or token.MINUS
	Left     Expression  // left expression to be evaluated
	Operator string      // operator to be applied to the left and right expressions
	Right    Expression  // right expression to be evaluated
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

// An if expression, e.g. if (x < y) { x } else { y }
type IfExpression struct {
	Token       token.Token     // token.IF
	Condition   Expression      // expression that evaluates to the condition of the if statement
	Consequence *BlockStatement // block statement that makes up the body of the if statement
	Alternative *BlockStatement // block statement that makes up the body of the else statement // or nil
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

// A function literal, e.g. fn(x, y) { x + y; }
type FunctionLiteral struct {
	Token      token.Token     // token.FUNCTION
	Parameters []*Identifier   // slice of identifiers that make up the parameters of the function
	Body       *BlockStatement // block statement that makes up the body of the function
	Name       string          // name of the function // should be the same as the token literal // or ""
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
	if funcLit.Name != "" {
		out.WriteString(fmt.Sprintf("<%s>", funcLit.Name))
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(funcLit.Body.String())

	return out.String()
}

// A function call expression, e.g. add(1, 2)
type CallExpression struct {
	Token     token.Token  // token.LPAREN '('
	Function  Expression   // Identifier or FunctionLiteral
	Arguments []Expression // slice of expressions that make up the arguments of the function
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

// An array literal, e.g. [1, 2, 3]
type ArrayLiteral struct {
	Token    token.Token  // token.LBRACKET '['
	Elements []Expression // slice of expressions that make up the elements of the array
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

// A hash literal, e.g. {"one": 1, "two": 2}
type HashLiteral struct {
	Token token.Token               // token.LBRACE '{'
	Pairs map[Expression]Expression // map of key-value pairs that make up the hashmap
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

// An index expression, e.g. array[1]
type IndexExpression struct {
	Token token.Token // token.LBRACKET '['
	Left  Expression  // expression to be indexed e.g. array literal
	Index Expression  // expression that evaluates to the index
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
