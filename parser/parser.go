package parser

import (
	"cidoka/ast"
	"cidoka/lexer"
	"cidoka/token"
	"fmt"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // == or !=
	LESSGREATER // <, >, <=, or >=
	SUM         // + or -
	PRODUCT     // *, / or %
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
)

/*
Maps token types to their respective precedences

The precedence of an operator determines how the operator is grouped with its operands.
For example, the expression 5 + 10 * 2 is grouped as 5 + (10 * 2) because the * operator
has a higher precedence than the + operator.
*/
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.LT_EQ:    LESSGREATER,
	token.GT_EQ:    LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.MODULO:   PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	lex    *lexer.Lexer // lexer instance
	errors []string     // parsing errors

	curToken  token.Token // current token
	peekToken token.Token // next token

	prefixParseFns map[token.TokenType]prefixParseFn // prefix parse functions
	infixParseFns  map[token.TokenType]infixParseFn  // infix parse functions
}

// ----------------------------------------------------------------------------
// 							 General methods
// ----------------------------------------------------------------------------

/* New returns a new Parser instance fully initialized */
func New(lex *lexer.Lexer) *Parser {
	parser := &Parser{
		lex:    lex,
		errors: []string{},
	}

	// Prefix parse functions
	parser.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	parser.registerPrefix(token.IDENT, parser.parseIdentifier)
	parser.registerPrefix(token.INT, parser.parseIntegerLiteral)
	parser.registerPrefix(token.FLOAT, parser.parseFloatLiteral)
	parser.registerPrefix(token.TRUE, parser.parseBoolean)
	parser.registerPrefix(token.FALSE, parser.parseBoolean)
	parser.registerPrefix(token.STRING, parser.parseStringLiteral)

	parser.registerPrefix(token.BANG, parser.parsePrefixExpression)
	parser.registerPrefix(token.MINUS, parser.parsePrefixExpression)

	parser.registerPrefix(token.IF, parser.parseIfExpression)
	parser.registerPrefix(token.FUNCTION, parser.parseFunctionLiteral)

	parser.registerPrefix(token.LPAREN, parser.parseGroupedExpression)
	parser.registerPrefix(token.LBRACKET, parser.parseArrayLiteral)
	parser.registerPrefix(token.LBRACE, parser.parseHashLiteral)

	// Infix parse functions
	parser.infixParseFns = make(map[token.TokenType]infixParseFn)
	parser.registerInfix(token.PLUS, parser.parseInfixExpression)
	parser.registerInfix(token.MINUS, parser.parseInfixExpression)
	parser.registerInfix(token.SLASH, parser.parseInfixExpression)
	parser.registerInfix(token.ASTERISK, parser.parseInfixExpression)
	parser.registerInfix(token.MODULO, parser.parseInfixExpression)

	parser.registerInfix(token.EQ, parser.parseInfixExpression)
	parser.registerInfix(token.NOT_EQ, parser.parseInfixExpression)
	parser.registerInfix(token.LT, parser.parseInfixExpression)
	parser.registerInfix(token.LT_EQ, parser.parseInfixExpression)
	parser.registerInfix(token.GT, parser.parseInfixExpression)
	parser.registerInfix(token.GT_EQ, parser.parseInfixExpression)

	parser.registerInfix(token.LPAREN, parser.parseCallExpression)
	parser.registerInfix(token.LBRACKET, parser.parseIndexExpression)

	// Read two tokens, so curToken and peekToken are both set
	parser.nextToken()
	parser.nextToken()

	return parser
}

/* Moves the parser to the next token */
func (parser *Parser) nextToken() {
	parser.curToken = parser.peekToken
	parser.peekToken = parser.lex.NextToken()
}

/* Returns True if the current token is of type t, False otherwise */
func (parser *Parser) curTokenIs(t token.TokenType) bool {
	return parser.curToken.Type == t
}

/* Returns True if the next token is of type t, False otherwise */
func (parser *Parser) peekTokenIs(t token.TokenType) bool {
	return parser.peekToken.Type == t
}

/*
Moves the parser to the next token if the next token is of type t and returns True.

Otherwise it appends an error message to the parser's errors list and returns False
*/
func (parser *Parser) expectPeek(t token.TokenType) bool {
	if parser.peekTokenIs(t) {
		parser.nextToken()
		return true
	} else {
		parser.peekError(t)
		return false
	}
}

// ----------------------------------------------------------------------------
// 							Error handling methods
// ----------------------------------------------------------------------------

/* Returns a slice of the parsing errors as strings */
func (parser *Parser) Errors() []string {
	return parser.errors
}

/* Appends an error message to the parser's errors list when the peekToken pointer didn't match the expected next token */
func (parser *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, parser.peekToken.Type)
	parser.errors = append(parser.errors, msg)
}

/* Appends an error message to the parser's errors list when no prefix parse function was found for a token type */
func (parser *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	parser.errors = append(parser.errors, msg)
}

func (parser *Parser) forLoopInvalidInitializerError() {
	msg := "FOR loop: expected let, identifier or empty initializer followed by a semicolon"
	parser.errors = append(parser.errors, msg)
}

func (parser *Parser) integerParseError() {
	msg := fmt.Sprintf("could not parse %q as integer", parser.curToken.Literal)
	parser.errors = append(parser.errors, msg)
}

func (parser *Parser) floatParseError() {
	msg := fmt.Sprintf("could not parse %q as float", parser.curToken.Literal)
	parser.errors = append(parser.errors, msg)
}

// ----------------------------------------------------------------------------
// 							  Parsing methods
// ----------------------------------------------------------------------------

/* Parses the input and returns the resulting AST */
func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for parser.curToken.Type != token.EOF {
		stmt := parser.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		parser.nextToken()
	}

	return program
}

// ----------------------------------------------------------------------------
// 								Statements
// ----------------------------------------------------------------------------

/* Parses a statement and returns the resulting AST node */
func (parser *Parser) parseStatement() ast.Statement {
	switch parser.curToken.Type {
	case token.LET:
		return parser.parseLetStatement()
	case token.RETURN:
		return parser.parseReturnStatement()
	case token.FOR:
		return parser.parseForLoopStatement()
	case token.BREAK:
		return parser.parseBreakStatement()
	case token.IDENT:
		if token.AssignmentOperators[parser.peekToken.Type] {
			return parser.parseAssignStatement()
		}
		fallthrough // if it's not an assign statement, parse it as an expression statement
	default:
		expr := parser.parseExpressionStatement()
		if expr != nil && expr.Expression != nil {
			return expr
		}
	}
	return nil
}

/* Parses a let statement and returns the resulting AST node */
func (parser *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: parser.curToken}

	if !parser.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}

	if !parser.expectPeek(token.ASSIGN) {
		return nil
	}

	parser.nextToken()

	stmt.Value = parser.parseExpression(LOWEST)

	if fl, ok := stmt.Value.(*ast.FunctionLiteral); ok {
		fl.Name = stmt.Name.Value
	}

	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return stmt
}

/* Parses a return statement and returns the resulting AST node */
func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: parser.curToken}

	parser.nextToken()

	stmt.ReturnValue = parser.parseExpression(LOWEST)

	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return stmt
}

/* Parses an assign statement and returns the resulting AST node */
func (parser *Parser) parseAssignStatement() *ast.AssignStatement {
	stmt := &ast.AssignStatement{}
	stmt.Name = &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}

	parser.nextToken()
	stmt.Token = parser.curToken

	parser.nextToken()

	stmt.Value = parser.parseExpression(LOWEST)

	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return stmt
}

/* Parses an expression statement and returns the resulting AST node */
func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: parser.curToken}

	stmt.Expression = parser.parseExpression(LOWEST)

	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return stmt
}

/* Parses a for loop statement and returns the resulting AST node */
func (parser *Parser) parseForLoopStatement() ast.Statement {
	stmt := &ast.ForLoopStatement{Token: parser.curToken}

	if !parser.expectPeek(token.LPAREN) {
		return nil
	}

	parser.nextToken()
	switch parser.curToken.Type {
	case token.LET:
		stmt.Initializer = parser.parseLetStatement()
	case token.IDENT:
		stmt.Initializer = parser.parseAssignStatement()
	case token.SEMICOLON:
		stmt.Initializer = nil
	default:
		parser.forLoopInvalidInitializerError()
		return nil
	}

	parser.nextToken()
	switch parser.curToken.Type {
	case token.SEMICOLON:
		stmt.Condition = nil
	default:
		stmt.Condition = parser.parseExpression(LOWEST)

		if !parser.expectPeek(token.SEMICOLON) {
			return nil
		}
	}

	parser.nextToken()
	switch parser.curToken.Type {
	case token.RPAREN:
		stmt.Update = nil
	default:
		stmt.Update = parser.parseAssignStatement()

		if !parser.expectPeek(token.RPAREN) {
			return nil
		}
	}

	if !parser.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = parser.parseBlockStatement()

	return stmt
}

/* Parses a break statement and returns the resulting AST node */
func (parser *Parser) parseBreakStatement() *ast.BreakStatement {
	stmt := &ast.BreakStatement{Token: parser.curToken}

	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return stmt
}

// ----------------------------------------------------------------------------
// 								Expressions
// ----------------------------------------------------------------------------

/* Parses an expression and returns the resulting AST node */
func (parser *Parser) parseExpression(precedence int) ast.Expression {
	prefix := parser.prefixParseFns[parser.curToken.Type]
	if prefix == nil {
		parser.noPrefixParseFnError(parser.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !parser.peekTokenIs(token.SEMICOLON) && precedence < parser.peekPrecedence() {
		infix := parser.infixParseFns[parser.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		parser.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

/* Parses an identifier and returns the resulting AST node */
func (parser *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}
}

/* Parses an integer literal and returns the resulting AST node */
func (parser *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: parser.curToken}

	value, err := strconv.ParseInt(parser.curToken.Literal, 0, 64)
	if err != nil {
		parser.integerParseError()
		return nil
	}

	lit.Value = value

	return lit
}

/* Parses a float literal and returns the resulting AST node */
func (parser *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: parser.curToken}

	value, err := strconv.ParseFloat(parser.curToken.Literal, 64)
	if err != nil {
		parser.floatParseError()
		return nil
	}

	lit.Value = value

	return lit
}

/* Parses a boolean and returns the resulting AST node */
func (parser *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: parser.curToken, Value: parser.curTokenIs(token.TRUE)}
}

/* Parses a string literal and returns the resulting AST node */
func (parser *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: parser.curToken, Value: parser.curToken.Literal}
}

/* Parses a prefix expression and returns the resulting AST node */
func (parser *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    parser.curToken,
		Operator: parser.curToken.Literal,
	}

	parser.nextToken()

	expression.Right = parser.parseExpression(PREFIX)

	return expression
}

/* Parses an infix expression and returns the resulting AST node */
func (parser *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    parser.curToken,
		Operator: parser.curToken.Literal,
		Left:     left,
	}

	precedence := parser.curPrecedence()
	parser.nextToken()
	expression.Right = parser.parseExpression(precedence)

	return expression
}

/* Parses a grouped expression and returns the resulting AST node */
func (parser *Parser) parseGroupedExpression() ast.Expression {
	parser.nextToken()

	exp := parser.parseExpression(LOWEST)

	if !parser.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

/* Parses an if expression and returns the resulting AST node */
func (parser *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: parser.curToken}

	if !parser.expectPeek(token.LPAREN) {
		return nil
	}

	parser.nextToken()
	expression.Condition = parser.parseExpression(LOWEST)

	if !parser.expectPeek(token.RPAREN) {
		return nil
	}

	if !parser.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = parser.parseBlockStatement()

	if parser.peekTokenIs(token.ELSE) {
		parser.nextToken()

		if parser.peekTokenIs(token.IF) {
			parser.nextToken()

			expression.Alternative = parser.parseExpressionStatement()
		} else {
			if !parser.expectPeek(token.LBRACE) {
				return nil
			}

			expression.Alternative = parser.parseBlockStatement()
		}
	}

	return expression
}

/* Parses a block statement and returns the resulting AST node */
func (parser *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: parser.curToken}
	block.Statements = []ast.Statement{}

	parser.nextToken()

	for !parser.curTokenIs(token.RBRACE) && !parser.curTokenIs(token.EOF) {
		stmt := parser.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		parser.nextToken()
	}

	return block
}

/* Parses a function literal and returns the resulting AST node */
func (parser *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: parser.curToken}

	if !parser.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = parser.parseFunctionParameters()

	if !parser.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = parser.parseBlockStatement()

	return lit
}

/* Parses function parameters and returns a slice of the resulting AST nodes */
func (parser *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if parser.peekTokenIs(token.RPAREN) {
		parser.nextToken()
		return identifiers
	}

	parser.nextToken()

	ident := &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}
	identifiers = append(identifiers, ident)

	for parser.peekTokenIs(token.COMMA) {
		parser.nextToken()
		parser.nextToken()

		ident := &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !parser.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

/* Parses a call expression and returns the resulting AST node */
func (parser *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: parser.curToken, Function: function}
	exp.Arguments = parser.parseExpressionList(token.RPAREN)
	return exp
}

/* Parses an array literal and returns the resulting AST node */
func (parser *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: parser.curToken}
	array.Elements = parser.parseExpressionList(token.RBRACKET)
	return array
}

/* Parses a list of expressions and returns a slice of the resulting AST nodes */
func (parser *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if parser.peekTokenIs(end) {
		parser.nextToken()
		return list
	}

	parser.nextToken()
	list = append(list, parser.parseExpression(LOWEST))

	for parser.peekTokenIs(token.COMMA) {
		parser.nextToken()
		parser.nextToken()
		list = append(list, parser.parseExpression(LOWEST))
	}

	if !parser.expectPeek(end) {
		return nil
	}

	return list
}

/* Parses an index expression and returns the resulting AST node */
func (parser *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: parser.curToken, Left: left}

	parser.nextToken()
	exp.Index = parser.parseExpression(LOWEST)

	if !parser.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

/* Parses a hash literal and returns the resulting AST node */
func (parser *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: parser.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !parser.peekTokenIs(token.RBRACE) {
		parser.nextToken()
		key := parser.parseExpression(LOWEST)

		if !parser.expectPeek(token.COLON) {
			return nil
		}

		parser.nextToken()
		value := parser.parseExpression(LOWEST)

		hash.Pairs[key] = value

		if !parser.peekTokenIs(token.RBRACE) && !parser.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !parser.expectPeek(token.RBRACE) {
		return nil
	}

	return hash
}

// ----------------------------------------------------------------------------
// 								Helper methods
// ----------------------------------------------------------------------------

/* Returns the precedence of the next token */
func (parser *Parser) peekPrecedence() int {
	if p, ok := precedences[parser.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

/* Returns the precedence of the current token */
func (parser *Parser) curPrecedence() int {
	if p, ok := precedences[parser.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

/* Registers a prefix parse function for a token type */
func (parser *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	parser.prefixParseFns[tokenType] = fn
}

/* Registers an infix parse function for a token type */
func (parser *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	parser.infixParseFns[tokenType] = fn
}
