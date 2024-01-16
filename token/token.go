package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL" // unknown token
	EOF     = "EOF"     // end of file

	// Identifiers + literals
	IDENT = "IDENT" // add, foobar, x, y, ...
	INT   = "INT"   // 1234567890

	// Assignment operators
	ASSIGN  = "="
	PLUSEQ  = "+="
	MINUSEQ = "-="
	MULTEQ  = "*="
	DIVEQ   = "/="
	MODEQ   = "%="

	// Arithmetic operators
	PLUS  = "+"
	MINUS = "-"
	MULT  = "*"
	DIV   = "/"
	MOD   = "%"

	// Comparison operators
	EQ  = "=="
	NEQ = "!="
	LT  = "<"
	LTE = "<="
	GT  = ">"
	GTE = ">="

	// Logical operators
	AND = "&&"
	OR  = "||"
	NOT = "!"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	// Brackets
	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
	LBRACK = "["
	RBRACK = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

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
