package token

type TokenType string

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
	PLUS_EQ     TokenType = "+=" // todo
	MINUS_EQ    TokenType = "-=" // todo
	ASTERISK_EQ TokenType = "*=" // todo
	SLASH_EQ    TokenType = "/=" // todo
	MODULO_EQ   TokenType = "%=" // todo

	// Arithmetic operators

	PLUS     TokenType = "+" // addition
	MINUS    TokenType = "-" // subtraction
	ASTERISK TokenType = "*" // multiplication
	SLASH    TokenType = "/" // division
	MODULO   TokenType = "%" // todo

	// Comparison operators

	EQ     TokenType = "==" // equality
	NOT_EQ TokenType = "!=" // inequality
	LT     TokenType = "<"  // less than
	LT_EQ  TokenType = "<=" // less than or equal to
	GT     TokenType = ">"  // greater than
	GT_EQ  TokenType = ">=" // greater than or equal to

	// Logical operators

	AND  TokenType = "&&" // todo
	OR   TokenType = "||" // todo
	BANG TokenType = "!"  // negation

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
	BREAK    TokenType = "BREAK"    // break statement
)

type Token struct {
	Type    TokenType // Type of token
	Literal string    // Literal value of token
}

// Map of Keywords to their TokenType constants.
var Keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"for":    FOR,
	"break":  BREAK,
}

/*
Checks the keywords table to see if the given identifier is a keyword.

If it is, it returns the keyword's TokenType constant. If not, it's an identifier.
*/
func LookupIdent(ident string) TokenType {
	if tok, ok := Keywords[ident]; ok {
		return tok
	}

	return IDENT
}
