package lexer

import "litlang/token"

type Lexer struct {
	input        string // input string to lex
	position     int    // current position in input (points to current char)
	readPosition int    // current reading position in input (after current char)
	char         byte   // current char under examination
}

func New(input string) *Lexer {
	lex := &Lexer{input: input}
	lex.readChar() // initialize lex.char and lex.readPosition
	return lex
}

func (lex *Lexer) readChar() {
	if lex.readPosition >= len(lex.input) { // check if we reached the end of input
		lex.char = 0 // ASCII code for "NUL" character
	} else {
		lex.char = lex.input[lex.readPosition] // read next character
	}
	lex.position = lex.readPosition // update position
	lex.readPosition += 1           // update readPosition
}

func (lex *Lexer) NextToken() token.Token {
	var tok token.Token

	lex.skipWhitespace()

	switch lex.char {
	case '=':
		tok = newToken(token.ASSIGN, lex.char)
	case ';':
		tok = newToken(token.SEMICOLON, lex.char)
	case '(':
		tok = newToken(token.LPAREN, lex.char)
	case ')':
		tok = newToken(token.RPAREN, lex.char)
	case ',':
		tok = newToken(token.COMMA, lex.char)
	case '+':
		tok = newToken(token.PLUS, lex.char)
	case '{':
		tok = newToken(token.LBRACE, lex.char)
	case '}':
		tok = newToken(token.RBRACE, lex.char)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(lex.char) {
			tok.Literal = lex.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(lex.char) {
			tok.Type = token.INT
			tok.Literal = lex.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, lex.char)
		}
	}

	lex.readChar()
	return tok
}

func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(char)}
}

func (lex *Lexer) readIdentifier() string {
	position := lex.position
	for isLetter(lex.char) {
		lex.readChar()
	}
	return lex.input[position:lex.position]
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func (lex *Lexer) skipWhitespace() {
	for lex.char == ' ' || lex.char == '\t' || lex.char == '\n' || lex.char == '\r' {
		lex.readChar()
	}
}

func (lex *Lexer) readNumber() string {
	position := lex.position
	for isDigit(lex.char) {
		lex.readChar()
	}
	return lex.input[position:lex.position]
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}
