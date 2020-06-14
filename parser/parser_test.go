package parser

import (
	"dara/ast"
	"dara/lexer"
	"fmt"
	"testing"
)

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return y;", "y"},
	}

	for _, tt := range tests {
		var (
			l       = lexer.New(tt.input)
			p       = New(l)
			program = p.ParseProgram()
		)

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		stmt := program.Statements[0]
		if _, ok := stmt.(*ast.ReturnStatement); !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}

		val := stmt.(*ast.ReturnStatement).ReturnValue
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestDeclareStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"x := 5;", "x", 5},
		{"y;", "y", nil},
		{"foobar := y;", "foobar", "y"},
	}

	for _, tt := range tests {
		var (
			l       = lexer.New(tt.input)
			p       = New(l)
			program = p.ParseProgram()
		)

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testDeclareStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.DeclareStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestAssignStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"x = 5;", "x", 5},
		{"y = nil;", "y", nil},
		{"foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		var (
			l       = lexer.New(tt.input)
			p       = New(l)
			program = p.ParseProgram()
		)

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testAssignStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.AssignStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	input := "true;"

	var (
		l       = lexer.New(input)
		p       = New(l)
		program = p.ParseProgram()
	)

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("not enough statements in program. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Boolean)
	if !ok {
		t.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
	}

	if ident.Value != true {
		t.Errorf("ident.Value not %t. got=%t", true, ident.Value)
	}

	if ident.TokenLiteral() != "true" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "true", ident.TokenLiteral())
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	var (
		l       = lexer.New(input)
		p       = New(l)
		program = p.ParseProgram()
	)

	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}

func TestIfExpression(t *testing.T) {
	input := `if x < y { x }`

	var (
		l       = lexer.New(input)
		p       = New(l)
		program = p.ParseProgram()
	)

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.IfStatement. got=%T", program.Statements[0])
	}

	if !testInfixExpression(t, stmt.Condition, "x", "<", "y") {
		return
	}

	if len(stmt.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(stmt.Consequence.Statements))
	}

	consequence, ok := stmt.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", stmt.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if stmt.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", stmt.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if x < y { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	if !testInfixExpression(t, stmt.Condition, "x", "<", "y") {
		return
	}

	if len(stmt.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(stmt.Consequence.Statements))
	}

	consequence, ok := stmt.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", stmt.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	alternative, ok := stmt.Alternative.(*ast.BlockStatement)
	if !ok {
		t.Fatalf("stmt.Alternative is not ast.BlockStatement. got=%T", stmt.Alternative)
	}

	altStmt, ok := alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("alternative.Statements[0] is not ast.ExpressionStatement. got=%T", alternative.Statements[0])
	}

	if !testIdentifier(t, altStmt.Expression, "y") {
		return
	}
}

func TestIfElseIfElseExpression(t *testing.T) {
	input := `if x < y { x } else if x > y { y } else { z }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	if !testInfixExpression(t, stmt.Condition, "x", "<", "y") {
		return
	}

	if len(stmt.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(stmt.Consequence.Statements))
	}

	consequence, ok := stmt.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", stmt.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	altStmt, ok := stmt.Alternative.(*ast.IfStatement)
	if !ok {
		t.Fatalf("stmt.Alternative is not ast.IfStatement. got=%T", stmt.Alternative)
	}

	if !testInfixExpression(t, altStmt.Condition, "x", ">", "y") {
		return
	}

	if len(altStmt.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(altStmt.Consequence.Statements))
	}

	altConsequence, ok := altStmt.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", altStmt.Consequence.Statements[0])
	}

	if !testIdentifier(t, altConsequence.Expression, "y") {
		return
	}

	altAlternative, ok := altStmt.Alternative.(*ast.BlockStatement)
	if !ok {
		t.Fatalf("altStmt.Alternative is not ast.BlockStatement. got=%T", altStmt.Alternative)
	}

	altExpression, ok := altAlternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("altAlternative.Statements[0] is not ast.ExpressionStatement. got=%T", altAlternative.Statements[0])
	}

	if !testIdentifier(t, altExpression.Expression, "z") {
		return
	}
}

func TestNumberLiteralExpression(t *testing.T) {
	input := `5.4;`

	var (
		l       = lexer.New(input)
		p       = New(l)
		program = p.ParseProgram()
	)

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("not enough statements in program. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.NumberLiteral)
	if !ok {
		t.Errorf("exp not *ast.NumberLiteral. got=%T", stmt.Expression)
	}
	if ident.Value != 5.4 {
		t.Errorf("ident.Value not %v. got=%v", 5, ident.Value)
	}
	if ident.TokenLiteral() != "5.4" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "5.4", ident.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input       string
		operator    string
		numberValue interface{}
	}{
		{"!5;", "!", 5},
		{"-15.2;", "-", 15.2},
		{"-true;", "-", true},
		{"-false;", "-", false},
	}

	for _, tt := range prefixTests {
		var (
			l       = lexer.New(tt.input)
			p       = New(l)
			program = p.ParseProgram()
		)

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %q. got=%s", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.numberValue) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 % 5;", 5, "%", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"5 >= 5;", 5, ">=", 5},
		{"5 <= 5;", 5, "<=", 5},
		{"5 && 5;", 5, "&&", 5},
		{"5 || 5;", 5, "||", 5},
		{"true != false", true, "!=", false},
	}

	for _, tt := range infixTests {
		var (
			l       = lexer.New(tt.input)
			p       = New(l)
			program = p.ParseProgram()
		)

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
		}

		if !testLiteralExpression(t, exp.Left, tt.leftValue) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 <= 3 * 1 + 4 * 5", "((3 + (4 * 5)) <= ((3 * 1) + (4 * 5)))"},
		{"-3 + 4 % 5 >= 3 * -1 + 4 * 5", "(((-3) + (4 % 5)) >= ((3 * (-1)) + (4 * 5)))"},
		{"-a * b || a == b", "(((-a) * b) || (a == b))"},
		{"-a * b && a == b", "(((-a) * b) && (a == b))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
	}

	for _, tt := range tests {
		var (
			l       = lexer.New(tt.input)
			p       = New(l)
			program = p.ParseProgram()
		)

		checkParserErrors(t, p)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	var (
		l       = lexer.New(input)
		p       = New(l)
		program = p.ParseProgram()
	)

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has more than 1 statement. got=%d\n", len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body statement is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		var (
			l       = lexer.New(tt.input)
			p       = New(l)
			program = p.ParseProgram()
		)

		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want=%d, got=%d\n", len(tt.expectedParams), len(function.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5);`

	var (
		l       = lexer.New(input)
		p       = New(l)
		program = p.ParseProgram()
	)

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testDeclareStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != ":=" {
		t.Errorf("s.TokenLiteral not ':='. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.DeclareStatement)
	if !ok {
		t.Errorf("s not *ast.DeclareStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name)
		return false
	}

	return true
}

func testAssignStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "=" {
		t.Errorf("s.TokenLiteral not '='. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.AssignStatement)
	if !ok {
		t.Errorf("s not *ast.AssignStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name)
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	iExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, iExp.Left, left) {
		return false
	}

	if iExp.Operator != operator {
		t.Errorf("exp.Operator is not %q. got=%q", operator, iExp.Operator)
		return false
	}

	if !testLiteralExpression(t, iExp.Right, right) {
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testNumberLiteral(t, exp, float64(v))
	case float64:
		return testNumberLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	case nil:
		return testNilLiteral(t, exp)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testNumberLiteral(t *testing.T, nl ast.Expression, value float64) bool {
	fl, ok := nl.(*ast.NumberLiteral)
	if !ok {
		t.Errorf("nl not *ast.NumberLiteral. got=%T", nl)
		return false
	}

	if fl.Value != value {
		t.Errorf("fl.Value not %v. got=%v", value, fl.Value)
		return false
	}

	// TODO: this test is going to fail since `5.00` converts to string `"5"`
	if fl.TokenLiteral() != fmt.Sprintf("%v", value) {
		t.Errorf("fl.TokenLiteral not %v. got=%s", value, fl.TokenLiteral())
		return false
	}
	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s", value, bo.TokenLiteral())
		return false
	}

	return true
}

func testNilLiteral(t *testing.T, exp ast.Expression) bool {
	if _, ok := exp.(*ast.Nil); !ok {
		t.Errorf("exp not *ast.Nil. got=%T", exp)
		return false
	}

	return true
}
