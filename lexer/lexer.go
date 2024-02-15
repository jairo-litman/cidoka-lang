package lexer

import (
	"cidoka/token"
)

type Lexer struct {
	input        string // input to be tokenized
	position     int    // current position in input (points to current char)
	readPosition int    // current reading position in input (after current char)
	ch           byte   // current char under examination
}

/* Returns a new Lexer instance fully initialized */
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

/* Returns the next token from the input */
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = l.compundableAssignment(token.ASSIGN, token.EQ)
	case '+':
		tok = l.compundableAssignment(token.PLUS, token.PLUS_EQ)
	case '-':
		tok = l.compundableAssignment(token.MINUS, token.MINUS_EQ)
	case '!':
		tok = l.compundableAssignment(token.BANG, token.NOT_EQ)
	case '*':
		tok = l.compundableAssignment(token.ASTERISK, token.ASTERISK_EQ)
	case '/':
		tok = l.compundableAssignment(token.SLASH, token.SLASH_EQ)
	case '%':
		tok = l.compundableAssignment(token.MODULO, token.MODULO_EQ)
	case '<':
		tok = l.compundableAssignment(token.LT, token.LT_EQ)
	case '>':
		tok = l.compundableAssignment(token.GT, token.GT_EQ)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		switch {
		// if it's a letter, it's an identifier or a keyword
		case isLetter(l.ch):
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok

		// if it's a digit or a dot followed by a digit, it's a number
		// THIS CAN ALSO THROW AN ILLEGAL TOKEN DUE TO MALFORMED NUMBERS
		case isDigit(l.ch) || (l.ch == '.' && isDigit(l.peekChar())):
			tok.Literal, tok.Type = l.readNumber()

			return tok

		// if it's none of the above, it's an illegal token
		default:
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

/*
Skips any whitespace characters in the input

These are ' ', '\t', '\n' and '\r'
*/
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

/*
Reads the next character in the input and advances the position
and readPosition pointers in the input string
*/
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

/*
Returns the next character in the input without advancing the
position and readPosition pointers in the input string
*/
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

/*
Reads the next identifier in the input and advances the
position and readPosition pointers in the input string to the end of the
identifier

It returns the identifier as a string
*/
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

/*
Reads the next number in the input and advances the position and
readPosition pointers in the input string to the end of the number

It returns the number as a string and its type (INT or FLOAT)

It can return an ILLEGAL token if the number is malformed (e.g. 12.34.56)
*/
func (l *Lexer) readNumber() (string, token.TokenType) {
	position := l.position
	isFloat := false
	dotCount := 0
	illegal := false

	for isDigit(l.ch) || l.ch == '.' {
		if l.ch == '.' {
			dotCount++
			if dotCount > 1 {
				illegal = true
				break
			}
			isFloat = true
		}
		l.readChar()
	}

	for illegal && !isPossibleTerminator(l.ch) {
		l.readChar()
	}

	if illegal {
		return l.input[position:l.position], token.ILLEGAL
	}

	if isFloat {
		return l.input[position:l.position], token.FLOAT
	} else {
		return l.input[position:l.position], token.INT
	}
}

/*
Reads the next compoundable assignment operator in the input and advances the position and
readPosition pointers in the input string to the end of the operator

If the next character is not an '=', it returns the single operator token.
If it is, it returns the compound operator token
*/
func (l *Lexer) compundableAssignment(single token.TokenType, compound token.TokenType) token.Token {
	ch := l.ch
	if l.peekChar() != '=' {
		return newToken(single, ch)
	}

	l.readChar()
	literal := string(ch) + string(l.ch)
	return token.Token{Type: compound, Literal: literal}
}

/*
Reads the next string in the input and advances the position and
readPosition pointers in the input string to the end of the string's closing ' " '

It returns the string as a string literal
*/
func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

/* Returns true if the given byte is a letter */
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

/* Returns true if the given byte is a digit */
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

/* Returns a new token with the given type and literal */
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

/*
Returns true if the given byte can be a terminator

These are ' ', 0, '\t', '\n', '\r', ';', ')', '}', ']', ','
*/
func isPossibleTerminator(ch byte) bool {
	return ch == ' ' || ch == 0 || ch == '\t' || ch == '\n' || ch == '\r' || ch == ';' || ch == ')' || ch == '}' || ch == ']' || ch == ','
}
