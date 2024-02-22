package lexer

import (
	"cidoka/token"
	"testing"
)

type ExpectedToken struct {
	expectedType    token.TokenType
	expectedLiteral string
}

func TestNextToken(t *testing.T) {
	tests := []struct {
		input    string
		expected []ExpectedToken
	}{
		{
			input: `let five = 5;
			five = 10;`,
			expected: []ExpectedToken{
				{token.LET, "let"},
				{token.IDENT, "five"},
				{token.ASSIGN, "="},
				{token.INT, "5"},
				{token.SEMICOLON, ";"},
				{token.IDENT, "five"},
				{token.ASSIGN, "="},
				{token.INT, "10"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		{
			input: `let add = fn(x, y) { x + y; };
			let result = add(five, ten);`,
			expected: []ExpectedToken{
				{token.LET, "let"},
				{token.IDENT, "add"},
				{token.ASSIGN, "="},
				{token.FUNCTION, "fn"},
				{token.LPAREN, "("},
				{token.IDENT, "x"},
				{token.COMMA, ","},
				{token.IDENT, "y"},
				{token.RPAREN, ")"},
				{token.LBRACE, "{"},
				{token.IDENT, "x"},
				{token.PLUS, "+"},
				{token.IDENT, "y"},
				{token.SEMICOLON, ";"},
				{token.RBRACE, "}"},
				{token.SEMICOLON, ";"},
				{token.LET, "let"},
				{token.IDENT, "result"},
				{token.ASSIGN, "="},
				{token.IDENT, "add"},
				{token.LPAREN, "("},
				{token.IDENT, "five"},
				{token.COMMA, ","},
				{token.IDENT, "ten"},
				{token.RPAREN, ")"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		{
			input: `!-+/*
			< <= > >= == !=
			+= -= *= /=
			% %=`,
			expected: []ExpectedToken{
				{token.BANG, "!"},
				{token.MINUS, "-"},
				{token.PLUS, "+"},
				{token.SLASH, "/"},
				{token.ASTERISK, "*"},
				{token.LT, "<"},
				{token.LT_EQ, "<="},
				{token.GT, ">"},
				{token.GT_EQ, ">="},
				{token.EQ, "=="},
				{token.NOT_EQ, "!="},
				{token.PLUS_EQ, "+="},
				{token.MINUS_EQ, "-="},
				{token.ASTERISK_EQ, "*="},
				{token.SLASH_EQ, "/="},
				{token.MODULO, "%"},
				{token.MODULO_EQ, "%="},
				{token.EOF, ""},
			},
		},
		{
			input: `if (5 < 10) {
				return true;
			} else {
				return false;
			}`,
			expected: []ExpectedToken{
				{token.IF, "if"},
				{token.LPAREN, "("},
				{token.INT, "5"},
				{token.LT, "<"},
				{token.INT, "10"},
				{token.RPAREN, ")"},
				{token.LBRACE, "{"},
				{token.RETURN, "return"},
				{token.TRUE, "true"},
				{token.SEMICOLON, ";"},
				{token.RBRACE, "}"},
				{token.ELSE, "else"},
				{token.LBRACE, "{"},
				{token.RETURN, "return"},
				{token.FALSE, "false"},
				{token.SEMICOLON, ";"},
				{token.RBRACE, "}"},
				{token.EOF, ""},
			},
		},
		{
			input: `"foobar"
			"foo bar"`,
			expected: []ExpectedToken{
				{token.STRING, "foobar"},
				{token.STRING, "foo bar"},
				{token.EOF, ""},
			},
		},
		{
			input: `[1, 2];
			{"foo": "bar"}`,
			expected: []ExpectedToken{
				{token.LBRACKET, "["},
				{token.INT, "1"},
				{token.COMMA, ","},
				{token.INT, "2"},
				{token.RBRACKET, "]"},
				{token.SEMICOLON, ";"},
				{token.LBRACE, "{"},
				{token.STRING, "foo"},
				{token.COLON, ":"},
				{token.STRING, "bar"},
				{token.RBRACE, "}"},
				{token.EOF, ""},
			},
		},
		{
			input: `123.456
			.123
			0.456`,
			expected: []ExpectedToken{
				{token.FLOAT, "123.456"},
				{token.FLOAT, ".123"},
				{token.FLOAT, "0.456"},
				{token.EOF, ""},
			},
		},
		{
			input: `for (let i = 0; i < 10; i = i + 1) {
				break;
			}`,
			expected: []ExpectedToken{
				{token.FOR, "for"},
				{token.LPAREN, "("},
				{token.LET, "let"},
				{token.IDENT, "i"},
				{token.ASSIGN, "="},
				{token.INT, "0"},
				{token.SEMICOLON, ";"},
				{token.IDENT, "i"},
				{token.LT, "<"},
				{token.INT, "10"},
				{token.SEMICOLON, ";"},
				{token.IDENT, "i"},
				{token.ASSIGN, "="},
				{token.IDENT, "i"},
				{token.PLUS, "+"},
				{token.INT, "1"},
				{token.RPAREN, ")"},
				{token.LBRACE, "{"},
				{token.BREAK, "break"},
				{token.SEMICOLON, ";"},
				{token.RBRACE, "}"},
				{token.EOF, ""},
			},
		},
		{
			input: `let x = 12.34.56`,
			expected: []ExpectedToken{
				{token.LET, "let"},
				{token.IDENT, "x"},
				{token.ASSIGN, "="},
				{token.ILLEGAL, "12.34.56"},
				{token.EOF, ""},
			},
		},
		{
			input: `let x = 12.34;56;`,
			expected: []ExpectedToken{
				{token.LET, "let"},
				{token.IDENT, "x"},
				{token.ASSIGN, "="},
				{token.FLOAT, "12.34"},
				{token.SEMICOLON, ";"},
				{token.INT, "56"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
		{
			input: `let x = 12.`,
			expected: []ExpectedToken{
				{token.LET, "let"},
				{token.IDENT, "x"},
				{token.ASSIGN, "="},
				{token.FLOAT, "12."},
				{token.EOF, ""},
			},
		},
		{
			input: `while (true) {
				continue;
			}`,
			expected: []ExpectedToken{
				{token.WHILE, "while"},
				{token.LPAREN, "("},
				{token.TRUE, "true"},
				{token.RPAREN, ")"},
				{token.LBRACE, "{"},
				{token.CONTINUE, "continue"},
				{token.SEMICOLON, ";"},
				{token.RBRACE, "}"},
				{token.EOF, ""},
			},
		},
		{
			input: `&& || & |`,
			expected: []ExpectedToken{
				{token.AND, "&&"},
				{token.OR, "||"},
				{token.ILLEGAL, "&"},
				{token.ILLEGAL, "|"},
				{token.EOF, ""},
			},
		},
		{
			input: `foo++ bar-- foo + bar`,
			expected: []ExpectedToken{
				{token.IDENT, "foo"},
				{token.INCREMENT, "++"},
				{token.IDENT, "bar"},
				{token.DECREMENT, "--"},
				{token.IDENT, "foo"},
				{token.PLUS, "+"},
				{token.IDENT, "bar"},
				{token.EOF, ""},
			},
		},
		{
			input: `let x int = 5;
			let y float = 10.5;
			let z int[] = [1, 2, 3];`,
			expected: []ExpectedToken{
				{token.LET, "let"},
				{token.IDENT, "x"},
				{token.INT_TYPE, "int"},
				{token.ASSIGN, "="},
				{token.INT, "5"},
				{token.SEMICOLON, ";"},
				{token.LET, "let"},
				{token.IDENT, "y"},
				{token.FLOAT_TYPE, "float"},
				{token.ASSIGN, "="},
				{token.FLOAT, "10.5"},
				{token.SEMICOLON, ";"},
				{token.LET, "let"},
				{token.IDENT, "z"},
				{token.INT_TYPE, "int"},
				{token.LBRACKET, "["},
				{token.RBRACKET, "]"},
				{token.ASSIGN, "="},
				{token.LBRACKET, "["},
				{token.INT, "1"},
				{token.COMMA, ","},
				{token.INT, "2"},
				{token.COMMA, ","},
				{token.INT, "3"},
				{token.RBRACKET, "]"},
				{token.SEMICOLON, ";"},
				{token.EOF, ""},
			},
		},
	}

	for _, tt := range tests {
		l := New(tt.input)

		for i, ts := range tt.expected {
			tok := l.NextToken()

			if tok.Type != ts.expectedType {
				t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
					i, ts.expectedType, tok.Type)
			}

			if tok.Literal != ts.expectedLiteral {
				t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
					i, ts.expectedLiteral, tok.Literal)
			}
		}
	}
}
