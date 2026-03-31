package parser

import (
	"fmt"
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParser_main_expr(t *testing.T) {
	assert := assert.New(t)

	t.Run("grouping", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `(a)
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `ParenExpr
 IdentExpr
  Name: "a" @1:2 (kind=3)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("grouping_bad_expr_operator_indent", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `(a b)
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("grouping_bad_expr_operator_int", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `(a 1)
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("error_unclosed", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `(a
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("error_empty", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `()
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		assert.Contains(ast.Dump(pr), "BadExpr")
		assert.Contains(ast.Dump(pr), "expected expression inside parentheses")
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("additive_plus", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `1+2
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `BinaryExpr
 IntLitExpr
  Value: "1" @1:1 (kind=4)
 Operator: "+" @1:2 (kind=51)
 IntLitExpr
  Value: "2" @1:3 (kind=4)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("additive_minus", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `1-2-3
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `BinaryExpr
 BinaryExpr
  IntLitExpr
   Value: "1" @1:1 (kind=4)
  Operator: "-" @1:2 (kind=54)
  IntLitExpr
   Value: "2" @1:3 (kind=4)
 Operator: "-" @1:4 (kind=54)
 IntLitExpr
  Value: "3" @1:5 (kind=4)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("multiplicative", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `1*2
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `BinaryExpr
 IntLitExpr
  Value: "1" @1:1 (kind=4)
 Operator: "*" @1:2 (kind=57)
 IntLitExpr
  Value: "2" @1:3 (kind=4)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("multiplicative_precedence", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `2+3*4+5
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `BinaryExpr
 BinaryExpr
  IntLitExpr
   Value: "2" @1:1 (kind=4)
  Operator: "+" @1:2 (kind=51)
  BinaryExpr
   IntLitExpr
    Value: "3" @1:3 (kind=4)
   Operator: "*" @1:4 (kind=57)
   IntLitExpr
    Value: "4" @1:5 (kind=4)
 Operator: "+" @1:6 (kind=51)
 IntLitExpr
  Value: "5" @1:7 (kind=4)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("grouping_precedence_prefix", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `1*(2+3)
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `BinaryExpr
 IntLitExpr
  Value: "1" @1:1 (kind=4)
 Operator: "*" @1:2 (kind=57)
 ParenExpr
  BinaryExpr
   IntLitExpr
    Value: "2" @1:4 (kind=4)
   Operator: "+" @1:5 (kind=51)
   IntLitExpr
    Value: "3" @1:6 (kind=4)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("grouping_precedence_postfix", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `(1+2)*3
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `BinaryExpr
 ParenExpr
  BinaryExpr
   IntLitExpr
    Value: "1" @1:2 (kind=4)
   Operator: "+" @1:3 (kind=51)
   IntLitExpr
    Value: "2" @1:4 (kind=4)
 Operator: "*" @1:6 (kind=57)
 IntLitExpr
  Value: "3" @1:7 (kind=4)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("grouping_of_grouping", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `(1+2)*(3+4)
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `BinaryExpr
 ParenExpr
  BinaryExpr
   IntLitExpr
    Value: "1" @1:2 (kind=4)
   Operator: "+" @1:3 (kind=51)
   IntLitExpr
    Value: "2" @1:4 (kind=4)
 Operator: "*" @1:6 (kind=57)
 ParenExpr
  BinaryExpr
   IntLitExpr
    Value: "3" @1:8 (kind=4)
   Operator: "+" @1:9 (kind=51)
   IntLitExpr
    Value: "4" @1:10 (kind=4)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("divide", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `8/4/2
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `BinaryExpr
 BinaryExpr
  IntLitExpr
   Value: "8" @1:1 (kind=4)
  Operator: "/" @1:2 (kind=59)
  IntLitExpr
   Value: "4" @1:3 (kind=4)
 Operator: "/" @1:4 (kind=59)
 IntLitExpr
  Value: "2" @1:5 (kind=4)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("grouping_binary", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `(a+b)
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `ParenExpr
 BinaryExpr
  IdentExpr
   Name: "a" @1:2 (kind=3)
  Operator: "+" @1:3 (kind=51)
  IdentExpr
   Name: "b" @1:4 (kind=3)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("grouping_binary_extended", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `a+(b+1)
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `BinaryExpr
 IdentExpr
  Name: "a" @1:1 (kind=3)
 Operator: "+" @1:2 (kind=51)
 ParenExpr
  BinaryExpr
   IdentExpr
    Name: "b" @1:4 (kind=3)
   Operator: "+" @1:5 (kind=51)
   IntLitExpr
    Value: "1" @1:6 (kind=4)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("unary_minus", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `-1*2
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `BinaryExpr
 UnaryExpr
  Operator: "-" @1:1 (kind=54)
  IntLitExpr
   Value: "1" @1:2 (kind=4)
 Operator: "*" @1:3 (kind=57)
 IntLitExpr
  Value: "2" @1:4 (kind=4)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("unary_minus_grouping", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `-(1+2)
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `UnaryExpr
 Operator: "-" @1:1 (kind=54)
 ParenExpr
  BinaryExpr
   IntLitExpr
    Value: "1" @1:3 (kind=4)
   Operator: "+" @1:4 (kind=51)
   IntLitExpr
    Value: "2" @1:5 (kind=4)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("comparison_unary", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `!a&&b
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `BinaryExpr
 UnaryExpr
  Operator: "!" @1:1 (kind=70)
  IdentExpr
   Name: "a" @1:2 (kind=3)
 Operator: "&&" @1:3 (kind=68)
 IdentExpr
  Name: "b" @1:5 (kind=3)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("unary_one", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `-1
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `UnaryExpr
 Operator: "-" @1:1 (kind=54)
 IntLitExpr
  Value: "1" @1:2 (kind=4)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("unary_minus_one", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `- -1
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `UnaryExpr
 Operator: "-" @1:1 (kind=54)
 UnaryExpr
  Operator: "-" @1:3 (kind=54)
  IntLitExpr
   Value: "1" @1:4 (kind=4)
`
		for _, v := range parser.errors {
			fmt.Println(v.Error())
		}
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("multiplicative_unary", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `1*-2
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `BinaryExpr
 IntLitExpr
  Value: "1" @1:1 (kind=4)
 Operator: "*" @1:2 (kind=57)
 UnaryExpr
  Operator: "-" @1:3 (kind=54)
  IntLitExpr
   Value: "2" @1:4 (kind=4)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("unary_not_not", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `!!a
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `UnaryExpr
 Operator: "!" @1:1 (kind=70)
 UnaryExpr
  Operator: "!" @1:2 (kind=70)
  IdentExpr
   Name: "a" @1:3 (kind=3)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("unary_not_selector", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `!(a.b)
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `UnaryExpr
 Operator: "!" @1:1 (kind=70)
 ParenExpr
  SelectorExpr
   X:
    IdentExpr
     Name: "a" @1:3 (kind=3)
   Dot: "." @1:4 (kind=48)
   Selector: "b" @1:5 (kind=3)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("comparison_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `a<b&&c<d
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `BinaryExpr
 BinaryExpr
  IdentExpr
   Name: "a" @1:1 (kind=3)
  Operator: "<" @1:2 (kind=64)
  IdentExpr
   Name: "b" @1:3 (kind=3)
 Operator: "&&" @1:4 (kind=68)
 BinaryExpr
  IdentExpr
   Name: "c" @1:6 (kind=3)
  Operator: "<" @1:7 (kind=64)
  IdentExpr
   Name: "d" @1:8 (kind=3)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("comparison_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `a==b||c==d
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `BinaryExpr
 BinaryExpr
  IdentExpr
   Name: "a" @1:1 (kind=3)
  Operator: "==" @1:2 (kind=62)
  IdentExpr
   Name: "b" @1:4 (kind=3)
 Operator: "||" @1:5 (kind=69)
 BinaryExpr
  IdentExpr
   Name: "c" @1:7 (kind=3)
  Operator: "==" @1:8 (kind=62)
  IdentExpr
   Name: "d" @1:10 (kind=3)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("comparison_chaining_error", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `a<b<c
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		assert.NotNil(ast.Dump(pr))
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("comparison_chaining_ok", func(t *testing.T) {
		/*
			This example is:
			- OK because the we only check the left side (if binary expr) and the operator but NOT the right side
			- It's the role of the type checker to valid if this expression is right.
			So a!=b<c is valid but also a!=(b<c)
		*/
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `a!=b<c
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `BinaryExpr
 IdentExpr
  Name: "a" @1:1 (kind=3)
 Operator: "!=" @1:2 (kind=63)
 BinaryExpr
  IdentExpr
   Name: "b" @1:4 (kind=3)
  Operator: "<" @1:5 (kind=64)
  IdentExpr
   Name: "c" @1:6 (kind=3)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("postfix_selector_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `a.b
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `SelectorExpr
 X:
  IdentExpr
   Name: "a" @1:1 (kind=3)
 Dot: "." @1:2 (kind=48)
 Selector: "b" @1:3 (kind=3)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("postfix_selector_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `a.b.c
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `SelectorExpr
 X:
  SelectorExpr
   X:
    IdentExpr
     Name: "a" @1:1 (kind=3)
   Dot: "." @1:2 (kind=48)
   Selector: "b" @1:3 (kind=3)
 Dot: "." @1:4 (kind=48)
 Selector: "c" @1:5 (kind=3)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("postfix_index_selector_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `a[0]
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `IndexExpr
 X:
 IdentExpr
  Name: "a" @1:1 (kind=3)
 LBracket: "[" @1:2 (kind=43)
  IntLitExpr
   Value: "0" @1:3 (kind=4)
 RBracket: "]" @1:4 (kind=44)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("postfix_index_selector_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `a[x]
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `IndexExpr
 X:
 IdentExpr
  Name: "a" @1:1 (kind=3)
 LBracket: "[" @1:2 (kind=43)
  IdentExpr
   Name: "x" @1:3 (kind=3)
 RBracket: "]" @1:4 (kind=44)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("postfix_index_selector_x3", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `a[1+2]
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `IndexExpr
 X:
 IdentExpr
  Name: "a" @1:1 (kind=3)
 LBracket: "[" @1:2 (kind=43)
  BinaryExpr
   IntLitExpr
    Value: "1" @1:3 (kind=4)
   Operator: "+" @1:4 (kind=51)
   IntLitExpr
    Value: "2" @1:5 (kind=4)
 RBracket: "]" @1:6 (kind=44)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("postfix_index_selector_error_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `a[x
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		assert.NotNil(ast.Dump(pr))
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("postfix_func_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `f()
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `CallExpr
 Callee
  IdentExpr
   Name: "f" @1:1 (kind=3)
 LParent: "(" @1:2 (kind=39)
 RParent: ")" @1:3 (kind=40)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("postfix_func_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `f(1,2+3)
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `CallExpr
 Callee
  IdentExpr
   Name: "f" @1:1 (kind=3)
 LParent: "(" @1:2 (kind=39)
 Args:
  IntLitExpr
   Value: "1" @1:3 (kind=4)
  BinaryExpr
   IntLitExpr
    Value: "2" @1:5 (kind=4)
   Operator: "+" @1:6 (kind=51)
   IntLitExpr
    Value: "3" @1:7 (kind=4)
 RParent: ")" @1:8 (kind=40)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("postfix_func_x3", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `f(1 2)
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		assert.NotNil(ast.Dump(pr))
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("postfix_func_x4", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `f(1
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		assert.NotNil(ast.Dump(pr))
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("postfix_func_selector_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `a.b(1)
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `CallExpr
 Callee
  SelectorExpr
   X:
    IdentExpr
     Name: "a" @1:1 (kind=3)
   Dot: "." @1:2 (kind=48)
   Selector: "b" @1:3 (kind=3)
 LParent: "(" @1:4 (kind=39)
 Args:
  IntLitExpr
   Value: "1" @1:5 (kind=4)
 RParent: ")" @1:6 (kind=40)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("postfix_func_selector_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `a.b(1)[2].c
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `SelectorExpr
 X:
  IndexExpr
   X:
   CallExpr
    Callee
     SelectorExpr
      X:
       IdentExpr
        Name: "a" @1:1 (kind=3)
      Dot: "." @1:2 (kind=48)
      Selector: "b" @1:3 (kind=3)
    LParent: "(" @1:4 (kind=39)
    Args:
     IntLitExpr
      Value: "1" @1:5 (kind=4)
    RParent: ")" @1:6 (kind=40)
   LBracket: "[" @1:7 (kind=43)
    IntLitExpr
     Value: "2" @1:8 (kind=4)
   RBracket: "]" @1:9 (kind=44)
 Dot: "." @1:10 (kind=48)
 Selector: "c" @1:11 (kind=3)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("postfix_func_selector_x3", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `x+f(1)
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `CallExpr
 Callee
  BinaryExpr
   IdentExpr
    Name: "x" @1:1 (kind=3)
   Operator: "+" @1:2 (kind=51)
   IdentExpr
    Name: "f" @1:3 (kind=3)
 LParent: "(" @1:4 (kind=39)
 Args:
  IntLitExpr
   Value: "1" @1:5 (kind=4)
 RParent: ")" @1:6 (kind=40)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("postfix_func_selector_x4", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `x*a.b(1)
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `CallExpr
 Callee
  BinaryExpr
   IdentExpr
    Name: "x" @1:1 (kind=3)
   Operator: "*" @1:2 (kind=57)
   SelectorExpr
    X:
     IdentExpr
      Name: "a" @1:3 (kind=3)
    Dot: "." @1:4 (kind=48)
    Selector: "b" @1:5 (kind=3)
 LParent: "(" @1:6 (kind=39)
 Args:
  IntLitExpr
   Value: "1" @1:7 (kind=4)
 RParent: ")" @1:8 (kind=40)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("postfix_func_selector_x5", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `a&&f()
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `CallExpr
 Callee
  BinaryExpr
   IdentExpr
    Name: "a" @1:1 (kind=3)
   Operator: "&&" @1:2 (kind=68)
   IdentExpr
    Name: "f" @1:4 (kind=3)
 LParent: "(" @1:5 (kind=39)
 RParent: ")" @1:6 (kind=40)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("postfix_func_selector_x6", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `-a.b()[0]
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `IndexExpr
 X:
 CallExpr
  Callee
   UnaryExpr
    Operator: "-" @1:1 (kind=54)
    SelectorExpr
     X:
      IdentExpr
       Name: "a" @1:2 (kind=3)
     Dot: "." @1:3 (kind=48)
     Selector: "b" @1:4 (kind=3)
  LParent: "(" @1:5 (kind=39)
  RParent: ")" @1:6 (kind=40)
 LBracket: "[" @1:7 (kind=43)
  IntLitExpr
   Value: "0" @1:8 (kind=4)
 RBracket: "]" @1:9 (kind=44)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("postfix_func_selector_x7", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `-a.b()[0
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("postfix_func_call_on_call", func(t *testing.T) {
		// https://stackoverflow.com/questions/48289821/golang-returning-functions
		// f()(1) can be read as g := f(); g(1)
		// If we don't want to support it, it will be the job of the linter
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `f()(1)
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		result := `CallExpr
 Callee
  CallExpr
   Callee
    IdentExpr
     Name: "f" @1:1 (kind=3)
   LParent: "(" @1:2 (kind=39)
   RParent: ")" @1:3 (kind=40)
 LParent: "(" @1:4 (kind=39)
 Args:
  IntLitExpr
   Value: "1" @1:5 (kind=4)
 RParent: ")" @1:6 (kind=40)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("postfix_func_bad_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `x+f(,1)
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("postfix_func_bad_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `x+f(1,,1)
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("postfix_func_bad_x3", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `x+f(1,)
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("postfix_func_bad_x4", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `x+f(1,)+2
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}
