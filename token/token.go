package token

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL" // unknown token
	EOF     TokenType = "EOF"     // end of file

	// Identifiers + literals
	IDENT  TokenType = "IDENT"  // add, foobar, x, y, ...
	INT    TokenType = "INT"    // 1234567890
	FLOAT  TokenType = "FLOAT"  // 123.456
	STRING TokenType = "STRING" // "foobar"

	// Assignment operators
	ASSIGN      TokenType = "="
	PLUS_EQ     TokenType = "+="
	MINUS_EQ    TokenType = "-="
	ASTERISK_EQ TokenType = "*="
	SLASH_EQ    TokenType = "/="
	MODULO_EQ   TokenType = "%=" // todo

	// Arithmetic operators
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"
	MODULO   TokenType = "%" // todo

	// Comparison operators
	EQ     TokenType = "=="
	NOT_EQ TokenType = "!="
	LT     TokenType = "<"
	LT_EQ  TokenType = "<=" // todo
	GT     TokenType = ">"
	GT_EQ  TokenType = ">=" // todo

	// Logical operators
	AND  TokenType = "&&" // todo
	OR   TokenType = "||" // todo
	BANG TokenType = "!"

	// Delimiters
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"
	COLON     TokenType = ":"

	// Brackets
	LPAREN   TokenType = "("
	RPAREN   TokenType = ")"
	LBRACE   TokenType = "{"
	RBRACE   TokenType = "}"
	LBRACKET TokenType = "["
	RBRACKET TokenType = "]"

	// Keywords
	FUNCTION TokenType = "FUNCTION"
	LET      TokenType = "LET"
	TRUE     TokenType = "TRUE"
	FALSE    TokenType = "FALSE"
	IF       TokenType = "IF"
	ELSE     TokenType = "ELSE"
	RETURN   TokenType = "RETURN"
)

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

/*
LookupIdent checks the keywords table to see if the given identifier is a keyword.

If it is, it returns the keyword's TokenType constant. If not, it's an identifier.
*/
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}
