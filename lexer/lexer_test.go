package lexer

import (
	"dara/token"
	"testing"
)

type tokenTest struct {
	expectedType    token.TokenType
	expectedLiteral string
}

func TestMatch(t *testing.T) {
	input := `=`
	tests := []tokenTest{
		{token.ASSIGN, "="},
		{token.EOF, ""},
	}
	testRunner(t, input, tests)

	input = `==`
	tests = []tokenTest{
		{token.EQ, "=="},
		{token.EOF, ""},
	}
	testRunner(t, input, tests)
}

func TestComment(t *testing.T) {
	input := `// Some comment
`
	tests := []tokenTest{
		{token.COMMENT, "Some comment"},
		{token.EOF, ""},
	}
	testRunner(t, input, tests)
}

func TestNoLineBreakComment(t *testing.T) {
	input := `// Some comment`
	tests := []tokenTest{
		{token.COMMENT, "Some comment"},
		{token.EOF, ""},
	}
	testRunner(t, input, tests)
}

func TestCStyleComment(t *testing.T) {
	input := `/* Some comment */`
	tests := []tokenTest{
		{token.COMMENT, "Some comment"},
		{token.EOF, ""},
	}
	testRunner(t, input, tests)
}

func TestCStyleComment2(t *testing.T) {
	input := `/*
Some comment
And more
*/`
	tests := []tokenTest{
		{token.COMMENT, "Some comment\nAnd more"},
		{token.EOF, ""},
	}
	testRunner(t, input, tests)
}

func TestNextToken(t *testing.T) {
	input := `let five5 = 5;
let ten = 10.0;

let add = fn(x, y) {
  x + y;
};

let result = add(five, ten);
!-*/5;
5 < 10 > 5 >= 10;

if (5 <= 10) {
    return true;
} else {
    return false;
}

// 10 == 10; 10 != 9;

/* comment */10

10 == "10"; 10 != '9';`
	tests := []tokenTest{
		// 0
		{token.LET, "let"},
		{token.IDENT, "five5"},
		{token.ASSIGN, "="},
		{token.NUMBER, "5"},
		{token.SEMICOLON, ";"},

		// 5
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.NUMBER, "10.0"},
		{token.SEMICOLON, ";"},

		// 10
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

		// 26
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

		// 36
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.ASTERISK, "*"},
		{token.SLASH, "/"},
		{token.NUMBER, "5"},
		{token.SEMICOLON, ";"},

		// 42
		{token.NUMBER, "5"},
		{token.LT, "<"},
		{token.NUMBER, "10"},
		{token.GT, ">"},
		{token.NUMBER, "5"},
		{token.GT_EQ, ">="},
		{token.NUMBER, "10"},
		{token.SEMICOLON, ";"},

		// 50
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.NUMBER, "5"},
		{token.LT_EQ, "<="},
		{token.NUMBER, "10"},
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

		// 67
		{token.COMMENT, "10 == 10; 10 != 9;"},

		// 68
		{token.COMMENT, "comment"},
		{token.NUMBER, "10"},

		// 68
		{token.NUMBER, "10"},
		{token.EQ, "=="},
		{token.STRING, "10"},
		{token.SEMICOLON, ";"},
		{token.NUMBER, "10"},
		{token.NOT_EQ, "!="},
		{token.STRING, "9"},
		{token.SEMICOLON, ";"},

		{token.EOF, ""},
	}
	testRunner(t, input, tests)
}

func TestScan(t *testing.T) {
	input := `let test = 5.2;`
	tests := []tokenTest{
		{token.LET, "let"},
		{token.IDENT, "test"},
		{token.ASSIGN, "="},
		{token.NUMBER, "5.2"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	l := New(input)
	tokens := l.Scan()
	for i, tt := range tests {
		token := tokens[i]
		if token.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, token.Type)
		}
		if token.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, token.Literal)
		}
	}
	testRunner(t, input, tests)
}

func testRunner(t *testing.T, input string, tests []tokenTest) {
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
