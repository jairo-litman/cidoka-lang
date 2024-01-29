package token

type TokenType string

const (
	ILLEGAL = "ILLEGAL" // unknown token
	EOF     = "EOF"     // end of file

	// Identifiers + literals
	IDENT  = "IDENT"  // add, foobar, x, y, ...
	INT    = "INT"    // 1234567890
	STRING = "STRING" // "foobar"

	// Assignment operators
	ASSIGN      = "="
	PLUS_EQ     = "+="
	MINUS_EQ    = "-="
	ASTERISK_EQ = "*="
	SLASH_EQ    = "/="
	MODULO_EQ   = "%=" // todo

	// Arithmetic operators
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	MODULO   = "%" // todo

	// Comparison operators
	EQ     = "=="
	NOT_EQ = "!="
	LT     = "<"
	LT_EQ  = "<=" // todo
	GT     = ">"
	GT_EQ  = ">=" // todo

	// Logical operators
	AND  = "&&" // todo
	OR   = "||" // todo
	BANG = "!"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	// Brackets
	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
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
