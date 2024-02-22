package token

type TokenType string

type Token struct {
	Type    TokenType // Type of token
	Literal string    // Literal value of token
}

const (
	// Special tokens

	ILLEGAL TokenType = "ILLEGAL" // unknown token
	EOF     TokenType = "EOF"     // end of file

	// Identifiers + literals

	IDENT  TokenType = "IDENT"  // add, foobar, x, y, ...
	INT    TokenType = "INT"    // 1234567890
	FLOAT  TokenType = "FLOAT"  // 123456.0987
	STRING TokenType = "STRING" // "foobar"

	// Assignment operators

	ASSIGN      TokenType = "="  // assignment
	PLUS_EQ     TokenType = "+=" // addition assignment
	MINUS_EQ    TokenType = "-=" // subtraction assignment
	ASTERISK_EQ TokenType = "*=" // multiplication assignment
	SLASH_EQ    TokenType = "/=" // division assignment
	MODULO_EQ   TokenType = "%=" // modulo assignment

	// Arithmetic operators

	PLUS     TokenType = "+" // addition
	MINUS    TokenType = "-" // subtraction
	ASTERISK TokenType = "*" // multiplication
	SLASH    TokenType = "/" // division
	MODULO   TokenType = "%" // modulo

	// Comparison operators

	EQ     TokenType = "==" // equality
	NOT_EQ TokenType = "!=" // inequality
	LT     TokenType = "<"  // less than
	LT_EQ  TokenType = "<=" // less than or equal to
	GT     TokenType = ">"  // greater than
	GT_EQ  TokenType = ">=" // greater than or equal to

	// Logical operators

	AND  TokenType = "&&" // and
	OR   TokenType = "||" // or
	BANG TokenType = "!"  // negation

	// Postfix operators

	INCREMENT TokenType = "++" // increment
	DECREMENT TokenType = "--" // decrement

	// Delimiters

	COMMA     TokenType = "," // separator
	SEMICOLON TokenType = ";" // terminator
	COLON     TokenType = ":" // separator

	// Brackets

	LPAREN   TokenType = "(" // left parenthesis
	RPAREN   TokenType = ")" // right parenthesis
	LBRACE   TokenType = "{" // left brace
	RBRACE   TokenType = "}" // right brace
	LBRACKET TokenType = "[" // left bracket
	RBRACKET TokenType = "]" // right bracket

	// Keywords

	FUNCTION TokenType = "FUNCTION" // function
	LET      TokenType = "LET"      // variable declaration
	TRUE     TokenType = "TRUE"     // boolean true
	FALSE    TokenType = "FALSE"    // boolean false
	IF       TokenType = "IF"       // if statement
	ELSE     TokenType = "ELSE"     // else statement
	RETURN   TokenType = "RETURN"   // return statement
	FOR      TokenType = "FOR"      // for loop
	WHILE    TokenType = "WHILE"    // while loop
	BREAK    TokenType = "BREAK"    // break statement
	CONTINUE TokenType = "CONTINUE" // continue statement

	// Types

	INT_TYPE      TokenType = "INT_TYPE"      // int type
	FLOAT_TYPE    TokenType = "FLOAT_TYPE"    // float type
	STRING_TYPE   TokenType = "STRING_TYPE"   // string type
	BOOL_TYPE     TokenType = "BOOL_TYPE"     // bool type
	FUNCTION_TYPE TokenType = "FUNCTION_TYPE" // function type
)

// Map of AssignmentOperators.
var AssignmentOperators = map[TokenType]bool{
	ASSIGN:      true,
	PLUS_EQ:     true,
	MINUS_EQ:    true,
	ASTERISK_EQ: true,
	SLASH_EQ:    true,
	MODULO_EQ:   true,
}

// Map of Keywords to their TokenType constants.
var Keywords = map[string]TokenType{
	"fn":       FUNCTION,
	"let":      LET,
	"true":     TRUE,
	"false":    FALSE,
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
	"for":      FOR,
	"while":    WHILE,
	"break":    BREAK,
	"continue": CONTINUE,
}

// Map of Types to their TokenType constants.
var Types = map[string]TokenType{
	"int":      INT_TYPE,
	"float":    FLOAT_TYPE,
	"string":   STRING_TYPE,
	"bool":     BOOL_TYPE,
	"function": FUNCTION_TYPE,
}

/*
Checks the keywords table to see if the given identifier is a keyword/type.

If it is, it returns the keyword's TokenType constant. If not, it's an identifier.
*/
func LookupIdent(ident string) TokenType {
	if tok, ok := Keywords[ident]; ok {
		return tok
	}

	if tok, ok := Types[ident]; ok {
		return tok
	}

	return IDENT
}
