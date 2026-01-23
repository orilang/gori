package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
	"github.com/stretchr/testify/assert"
)

func TestParser_expr(t *testing.T) {
	assert := assert.New(t)
	t.Run("grouping", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.LParen, Value: "(", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "a", Line: 1, Column: 2},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 3},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
		pr := parser.parseExpr(LOWEST)
		result := `ParenExpr
 IdentExpr
  Name: "a" @1:2 (kind=3)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("grouping_bad_expr_operator_indent", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.LParen, Value: "(", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "a", Line: 1, Column: 2},
			{Kind: token.Ident, Value: "b", Line: 1, Column: 4},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 5},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
		pr := parser.parseExpr(LOWEST)
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("grouping_bad_expr_operator_int", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.LParen, Value: "(", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "a", Line: 1, Column: 2},
			{Kind: token.IntLit, Value: "1", Line: 1, Column: 4},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 5},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
		pr := parser.parseExpr(LOWEST)
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("error_unclosed", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.LParen, Value: "(", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "a", Line: 1, Column: 2},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
		pr := parser.parseExpr(LOWEST)
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("error_empty", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.LParen, Value: "(", Line: 1, Column: 1},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 2},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
		pr := parser.parseExpr(LOWEST)
		assert.Contains(ast.Dump(pr), "BadExpr")
		assert.Contains(ast.Dump(pr), "expected expression inside parentheses")
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("additive_plus", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.IntLit, Value: "1", Line: 1, Column: 1},
			{Kind: token.Plus, Value: "+", Line: 1, Column: 2},
			{Kind: token.IntLit, Value: "2", Line: 1, Column: 3},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
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
		input := []token.Token{
			{Kind: token.IntLit, Value: "1", Line: 1, Column: 1},
			{Kind: token.Minus, Value: "-", Line: 1, Column: 2},
			{Kind: token.IntLit, Value: "2", Line: 1, Column: 3},
			{Kind: token.Minus, Value: "-", Line: 1, Column: 4},
			{Kind: token.IntLit, Value: "3", Line: 1, Column: 5},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
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
		input := []token.Token{
			{Kind: token.IntLit, Value: "1", Line: 1, Column: 2},
			{Kind: token.Star, Value: "*", Line: 1, Column: 3},
			{Kind: token.IntLit, Value: "2", Line: 1, Column: 4},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
		pr := parser.parseExpr(LOWEST)
		result := `BinaryExpr
 IntLitExpr
  Value: "1" @1:2 (kind=4)
 Operator: "*" @1:3 (kind=57)
 IntLitExpr
  Value: "2" @1:4 (kind=4)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("multiplicative_precedence", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.IntLit, Value: "2", Line: 1, Column: 1},
			{Kind: token.Plus, Value: "+", Line: 1, Column: 2},
			{Kind: token.IntLit, Value: "3", Line: 1, Column: 3},
			{Kind: token.Star, Value: "*", Line: 1, Column: 4},
			{Kind: token.IntLit, Value: "4", Line: 1, Column: 5},
			{Kind: token.Plus, Value: "+", Line: 1, Column: 6},
			{Kind: token.IntLit, Value: "5", Line: 1, Column: 7},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
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
		input := []token.Token{
			{Kind: token.IntLit, Value: "1", Line: 1, Column: 1},
			{Kind: token.Star, Value: "*", Line: 1, Column: 2},
			{Kind: token.LParen, Value: "(", Line: 1, Column: 3},
			{Kind: token.IntLit, Value: "2", Line: 1, Column: 4},
			{Kind: token.Plus, Value: "+", Line: 1, Column: 5},
			{Kind: token.IntLit, Value: "3", Line: 1, Column: 6},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 7},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
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
		input := []token.Token{
			{Kind: token.LParen, Value: "(", Line: 1, Column: 1},
			{Kind: token.IntLit, Value: "1", Line: 1, Column: 2},
			{Kind: token.Plus, Value: "+", Line: 1, Column: 3},
			{Kind: token.IntLit, Value: "2", Line: 1, Column: 4},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 5},
			{Kind: token.Star, Value: "*", Line: 1, Column: 6},
			{Kind: token.IntLit, Value: "3", Line: 1, Column: 7},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
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
		input := []token.Token{
			{Kind: token.LParen, Value: "(", Line: 1, Column: 1},
			{Kind: token.IntLit, Value: "1", Line: 1, Column: 2},
			{Kind: token.Plus, Value: "+", Line: 1, Column: 3},
			{Kind: token.IntLit, Value: "2", Line: 1, Column: 4},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 5},
			{Kind: token.Star, Value: "*", Line: 1, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 1, Column: 7},
			{Kind: token.IntLit, Value: "3", Line: 1, Column: 8},
			{Kind: token.Plus, Value: "+", Line: 1, Column: 9},
			{Kind: token.IntLit, Value: "4", Line: 1, Column: 10},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 11},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
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
		input := []token.Token{
			{Kind: token.IntLit, Value: "8", Line: 1, Column: 1},
			{Kind: token.Slash, Value: "/", Line: 1, Column: 2},
			{Kind: token.IntLit, Value: "4", Line: 1, Column: 3},
			{Kind: token.Slash, Value: "/", Line: 1, Column: 4},
			{Kind: token.IntLit, Value: "2", Line: 1, Column: 5},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
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
		input := []token.Token{
			{Kind: token.LParen, Value: "(", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "a", Line: 1, Column: 2},
			{Kind: token.Plus, Value: "+", Line: 1, Column: 3},
			{Kind: token.Ident, Value: "b", Line: 1, Column: 4},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 5},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
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
		input := []token.Token{
			{Kind: token.Ident, Value: "a", Line: 1, Column: 1},
			{Kind: token.Plus, Value: "+", Line: 1, Column: 2},
			{Kind: token.LParen, Value: "(", Line: 1, Column: 3},
			{Kind: token.Ident, Value: "b", Line: 1, Column: 4},
			{Kind: token.Plus, Value: "+", Line: 1, Column: 5},
			{Kind: token.IntLit, Value: "1", Line: 1, Column: 6},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 7},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
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
		input := []token.Token{
			{Kind: token.Minus, Value: "-", Line: 1, Column: 1},
			{Kind: token.IntLit, Value: "1", Line: 1, Column: 2},
			{Kind: token.Star, Value: "*", Line: 1, Column: 3},
			{Kind: token.IntLit, Value: "2", Line: 1, Column: 4},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
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
		input := []token.Token{
			{Kind: token.Minus, Value: "-", Line: 1, Column: 1},
			{Kind: token.LParen, Value: "(", Line: 1, Column: 2},
			{Kind: token.IntLit, Value: "1", Line: 1, Column: 3},
			{Kind: token.Plus, Value: "+", Line: 1, Column: 4},
			{Kind: token.IntLit, Value: "2", Line: 1, Column: 5},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 6},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
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
		input := []token.Token{
			{Kind: token.Not, Value: "!", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "a", Line: 1, Column: 2},
			{Kind: token.And, Value: "&&", Line: 1, Column: 3},
			{Kind: token.Ident, Value: "b", Line: 1, Column: 4},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
		pr := parser.parseExpr(LOWEST)
		result := `BinaryExpr
 UnaryExpr
  Operator: "!" @1:1 (kind=70)
  IdentExpr
   Name: "a" @1:2 (kind=3)
 Operator: "&&" @1:3 (kind=68)
 IdentExpr
  Name: "b" @1:4 (kind=3)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("unary_one", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.Not, Value: "-", Line: 1, Column: 1},
			{Kind: token.IntLit, Value: "1", Line: 1, Column: 2},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
		pr := parser.parseExpr(LOWEST)
		result := `UnaryExpr
 Operator: "-" @1:1 (kind=70)
 IntLitExpr
  Value: "1" @1:2 (kind=4)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("unary_minus_one", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.Minus, Value: "-", Line: 1, Column: 1},
			{Kind: token.Minus, Value: "-", Line: 1, Column: 2},
			{Kind: token.IntLit, Value: "1", Line: 1, Column: 3},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
		pr := parser.parseExpr(LOWEST)
		result := `UnaryExpr
 Operator: "-" @1:1 (kind=54)
 UnaryExpr
  Operator: "-" @1:2 (kind=54)
  IntLitExpr
   Value: "1" @1:3 (kind=4)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("multiplicative_unary", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.IntLit, Value: "1", Line: 1, Column: 1},
			{Kind: token.Star, Value: "*", Line: 1, Column: 2},
			{Kind: token.Minus, Value: "-", Line: 1, Column: 3},
			{Kind: token.IntLit, Value: "2", Line: 1, Column: 4},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
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
		input := []token.Token{
			{Kind: token.Not, Value: "!", Line: 1, Column: 1},
			{Kind: token.Not, Value: "!", Line: 1, Column: 2},
			{Kind: token.Ident, Value: "a", Line: 1, Column: 3},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
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
		input := []token.Token{
			{Kind: token.Not, Value: "!", Line: 1, Column: 1},
			{Kind: token.LParen, Value: "(", Line: 1, Column: 2},
			{Kind: token.Ident, Value: "a", Line: 1, Column: 3},
			{Kind: token.Dot, Value: ".", Line: 1, Column: 4},
			{Kind: token.Ident, Value: "b", Line: 1, Column: 5},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 6},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
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
		input := []token.Token{
			{Kind: token.Ident, Value: "a", Line: 1, Column: 1},
			{Kind: token.Lt, Value: "<", Line: 1, Column: 2},
			{Kind: token.Ident, Value: "b", Line: 1, Column: 3},
			{Kind: token.And, Value: "&&", Line: 1, Column: 4},
			{Kind: token.Ident, Value: "c", Line: 1, Column: 5},
			{Kind: token.Lt, Value: "<", Line: 1, Column: 6},
			{Kind: token.Ident, Value: "d", Line: 1, Column: 7},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
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
   Name: "c" @1:5 (kind=3)
  Operator: "<" @1:6 (kind=64)
  IdentExpr
   Name: "d" @1:7 (kind=3)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("comparison_x2", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.Ident, Value: "a", Line: 1, Column: 1},
			{Kind: token.Eq, Value: "==", Line: 1, Column: 2},
			{Kind: token.Ident, Value: "b", Line: 1, Column: 3},
			{Kind: token.Or, Value: "||", Line: 1, Column: 4},
			{Kind: token.Ident, Value: "c", Line: 1, Column: 5},
			{Kind: token.Eq, Value: "==", Line: 1, Column: 6},
			{Kind: token.Ident, Value: "d", Line: 1, Column: 7},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
		pr := parser.parseExpr(LOWEST)
		result := `BinaryExpr
 BinaryExpr
  IdentExpr
   Name: "a" @1:1 (kind=3)
  Operator: "==" @1:2 (kind=62)
  IdentExpr
   Name: "b" @1:3 (kind=3)
 Operator: "||" @1:4 (kind=69)
 BinaryExpr
  IdentExpr
   Name: "c" @1:5 (kind=3)
  Operator: "==" @1:6 (kind=62)
  IdentExpr
   Name: "d" @1:7 (kind=3)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("comparison_chaining_error", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.Ident, Value: "a", Line: 1, Column: 1},
			{Kind: token.Lt, Value: "<", Line: 1, Column: 2},
			{Kind: token.Ident, Value: "b", Line: 1, Column: 3},
			{Kind: token.Lt, Value: "<", Line: 1, Column: 4},
			{Kind: token.Ident, Value: "c", Line: 1, Column: 5},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
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
		input := []token.Token{
			{Kind: token.Ident, Value: "a", Line: 1, Column: 1},
			{Kind: token.Neq, Value: "!=", Line: 1, Column: 2},
			{Kind: token.Ident, Value: "b", Line: 1, Column: 3},
			{Kind: token.Lt, Value: "<", Line: 1, Column: 4},
			{Kind: token.Ident, Value: "c", Line: 1, Column: 5},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
		pr := parser.parseExpr(LOWEST)
		result := `BinaryExpr
 IdentExpr
  Name: "a" @1:1 (kind=3)
 Operator: "!=" @1:2 (kind=63)
 BinaryExpr
  IdentExpr
   Name: "b" @1:3 (kind=3)
  Operator: "<" @1:4 (kind=64)
  IdentExpr
   Name: "c" @1:5 (kind=3)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("postfix_selector_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.Ident, Value: "a", Line: 1, Column: 1},
			{Kind: token.Dot, Value: ".", Line: 1, Column: 2},
			{Kind: token.Ident, Value: "b", Line: 1, Column: 3},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
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
		input := []token.Token{
			{Kind: token.Ident, Value: "a", Line: 1, Column: 1},
			{Kind: token.Dot, Value: ".", Line: 1, Column: 2},
			{Kind: token.Ident, Value: "b", Line: 1, Column: 3},
			{Kind: token.Dot, Value: ".", Line: 1, Column: 4},
			{Kind: token.Ident, Value: "c", Line: 1, Column: 5},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
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
		input := []token.Token{
			{Kind: token.Ident, Value: "a", Line: 1, Column: 1},
			{Kind: token.LBracket, Value: "[", Line: 1, Column: 2},
			{Kind: token.IntLit, Value: "0", Line: 1, Column: 3},
			{Kind: token.RBracket, Value: "]", Line: 1, Column: 4},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
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
		input := []token.Token{
			{Kind: token.Ident, Value: "a", Line: 1, Column: 1},
			{Kind: token.LBracket, Value: "[", Line: 1, Column: 2},
			{Kind: token.Ident, Value: "x", Line: 1, Column: 3},
			{Kind: token.RBracket, Value: "]", Line: 1, Column: 4},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
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
		input := []token.Token{
			{Kind: token.Ident, Value: "a", Line: 1, Column: 1},
			{Kind: token.LBracket, Value: "[", Line: 1, Column: 2},
			{Kind: token.IntLit, Value: "1", Line: 1, Column: 3},
			{Kind: token.Plus, Value: "+", Line: 1, Column: 3},
			{Kind: token.IntLit, Value: "2", Line: 1, Column: 3},
			{Kind: token.RBracket, Value: "]", Line: 1, Column: 4},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
		pr := parser.parseExpr(LOWEST)
		result := `IndexExpr
 X:
 IdentExpr
  Name: "a" @1:1 (kind=3)
 LBracket: "[" @1:2 (kind=43)
  BinaryExpr
   IntLitExpr
    Value: "1" @1:3 (kind=4)
   Operator: "+" @1:3 (kind=51)
   IntLitExpr
    Value: "2" @1:3 (kind=4)
 RBracket: "]" @1:4 (kind=44)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("postfix_index_selector_error_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.Ident, Value: "a", Line: 1, Column: 1},
			{Kind: token.LBracket, Value: "[", Line: 1, Column: 2},
			{Kind: token.Ident, Value: "x", Line: 1, Column: 3},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
		pr := parser.parseExpr(LOWEST)
		assert.NotNil(ast.Dump(pr))
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("postfix_func_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.Ident, Value: "f", Line: 1, Column: 1},
			{Kind: token.LParen, Value: "(", Line: 1, Column: 2},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 3},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
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
		input := []token.Token{
			{Kind: token.Ident, Value: "f", Line: 1, Column: 1},
			{Kind: token.LParen, Value: "(", Line: 1, Column: 2},
			{Kind: token.IntLit, Value: "1", Line: 1, Column: 3},
			{Kind: token.Comma, Value: ",", Line: 1, Column: 4},
			{Kind: token.IntLit, Value: "2", Line: 1, Column: 5},
			{Kind: token.Plus, Value: "+", Line: 1, Column: 6},
			{Kind: token.IntLit, Value: "3", Line: 1, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 8},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
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
		input := []token.Token{
			{Kind: token.Ident, Value: "f", Line: 1, Column: 1},
			{Kind: token.LParen, Value: "(", Line: 1, Column: 2},
			{Kind: token.IntLit, Value: "1", Line: 1, Column: 3},
			{Kind: token.IntLit, Value: "2", Line: 1, Column: 5},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 4},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
		pr := parser.parseExpr(LOWEST)
		assert.NotNil(ast.Dump(pr))
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("postfix_func_x4", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.Ident, Value: "f", Line: 1, Column: 1},
			{Kind: token.LParen, Value: "(", Line: 1, Column: 2},
			{Kind: token.IntLit, Value: "1", Line: 1, Column: 3},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
		pr := parser.parseExpr(LOWEST)
		assert.NotNil(ast.Dump(pr))
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("postfix_func_selector_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.Ident, Value: "a", Line: 1, Column: 1},
			{Kind: token.Dot, Value: ".", Line: 1, Column: 2},
			{Kind: token.Ident, Value: "b", Line: 1, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 1, Column: 4},
			{Kind: token.IntLit, Value: "1", Line: 1, Column: 5},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 6},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
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
		input := []token.Token{
			{Kind: token.Ident, Value: "a", Line: 1, Column: 1},
			{Kind: token.Dot, Value: ".", Line: 1, Column: 2},
			{Kind: token.Ident, Value: "b", Line: 1, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 1, Column: 4},
			{Kind: token.IntLit, Value: "1", Line: 1, Column: 5},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 6},
			{Kind: token.LBracket, Value: "[", Line: 1, Column: 7},
			{Kind: token.IntLit, Value: "2", Line: 1, Column: 8},
			{Kind: token.RBracket, Value: "]", Line: 1, Column: 9},
			{Kind: token.Dot, Value: ".", Line: 1, Column: 10},
			{Kind: token.Ident, Value: "c", Line: 1, Column: 11},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
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
		input := []token.Token{
			{Kind: token.Ident, Value: "x", Line: 1, Column: 1},
			{Kind: token.Plus, Value: "+", Line: 1, Column: 2},
			{Kind: token.Ident, Value: "f", Line: 1, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 1, Column: 4},
			{Kind: token.IntLit, Value: "1", Line: 1, Column: 5},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 6},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
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
		input := []token.Token{
			{Kind: token.Ident, Value: "x", Line: 1, Column: 1},
			{Kind: token.Star, Value: "*", Line: 1, Column: 2},
			{Kind: token.Ident, Value: "a", Line: 1, Column: 3},
			{Kind: token.Dot, Value: ".", Line: 1, Column: 4},
			{Kind: token.Ident, Value: "b", Line: 1, Column: 5},
			{Kind: token.LParen, Value: "(", Line: 1, Column: 6},
			{Kind: token.IntLit, Value: "1", Line: 1, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 8},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
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
		input := []token.Token{
			{Kind: token.Ident, Value: "a", Line: 1, Column: 1},
			{Kind: token.And, Value: "&&", Line: 1, Column: 2},
			{Kind: token.Ident, Value: "f", Line: 1, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 1, Column: 6},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 8},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
		pr := parser.parseExpr(LOWEST)
		result := `CallExpr
 Callee
  BinaryExpr
   IdentExpr
    Name: "a" @1:1 (kind=3)
   Operator: "&&" @1:2 (kind=68)
   IdentExpr
    Name: "f" @1:3 (kind=3)
 LParent: "(" @1:6 (kind=39)
 RParent: ")" @1:8 (kind=40)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("postfix_func_selector_x6", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.Minus, Value: "-", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "a", Line: 1, Column: 2},
			{Kind: token.Dot, Value: ".", Line: 1, Column: 3},
			{Kind: token.Ident, Value: "b", Line: 1, Column: 4},
			{Kind: token.LParen, Value: "(", Line: 1, Column: 5},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 6},
			{Kind: token.LBracket, Value: "[", Line: 1, Column: 7},
			{Kind: token.IntLit, Value: "0", Line: 1, Column: 8},
			{Kind: token.RBracket, Value: "]", Line: 1, Column: 9},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
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
		input := []token.Token{
			{Kind: token.Minus, Value: "-", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "a", Line: 1, Column: 2},
			{Kind: token.Dot, Value: ".", Line: 1, Column: 3},
			{Kind: token.Ident, Value: "b", Line: 1, Column: 4},
			{Kind: token.LParen, Value: "(", Line: 1, Column: 5},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 6},
			{Kind: token.LBracket, Value: "[", Line: 1, Column: 7},
			{Kind: token.IntLit, Value: "0", Line: 1, Column: 8},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
		pr := parser.parseExpr(LOWEST)
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("postfix_func_call_on_call", func(t *testing.T) {
		// https://stackoverflow.com/questions/48289821/golang-returning-functions
		// f()(1) can be read as g := f(); g(1)
		// If we don't want to support it, it will be the job of the linter
		input := []token.Token{
			{Kind: token.Ident, Value: "f", Line: 1, Column: 1},
			{Kind: token.LParen, Value: "(", Line: 1, Column: 2},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 1, Column: 4},
			{Kind: token.IntLit, Value: "1", Line: 1, Column: 5},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 6},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
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
		input := []token.Token{
			{Kind: token.Ident, Value: "x", Line: 1, Column: 1},
			{Kind: token.Plus, Value: "+", Line: 1, Column: 2},
			{Kind: token.Ident, Value: "f", Line: 1, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 1, Column: 4},
			{Kind: token.Comma, Value: ",", Line: 1, Column: 5},
			{Kind: token.IntLit, Value: "1", Line: 1, Column: 6},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 8},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
		pr := parser.parseExpr(LOWEST)
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("postfix_func_bad_x2", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.Ident, Value: "x", Line: 1, Column: 1},
			{Kind: token.Plus, Value: "+", Line: 1, Column: 2},
			{Kind: token.Ident, Value: "f", Line: 1, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 1, Column: 4},
			{Kind: token.IntLit, Value: "1", Line: 1, Column: 5},
			{Kind: token.Comma, Value: ",", Line: 1, Column: 6},
			{Kind: token.Comma, Value: ",", Line: 1, Column: 7},
			{Kind: token.IntLit, Value: "1", Line: 1, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 9},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
		pr := parser.parseExpr(LOWEST)
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("postfix_func_bad_x3", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.Ident, Value: "x", Line: 1, Column: 1},
			{Kind: token.Plus, Value: "+", Line: 1, Column: 2},
			{Kind: token.Ident, Value: "f", Line: 1, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 1, Column: 4},
			{Kind: token.IntLit, Value: "1", Line: 1, Column: 5},
			{Kind: token.Comma, Value: ",", Line: 1, Column: 6},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 7},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
		pr := parser.parseExpr(LOWEST)
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("postfix_func_bad_x4", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.Ident, Value: "x", Line: 1, Column: 1},
			{Kind: token.Plus, Value: "+", Line: 1, Column: 2},
			{Kind: token.Ident, Value: "f", Line: 1, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 1, Column: 4},
			{Kind: token.IntLit, Value: "1", Line: 1, Column: 5},
			{Kind: token.Comma, Value: ",", Line: 1, Column: 6},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 7},
			{Kind: token.Plus, Value: "+", Line: 1, Column: 8},
			{Kind: token.IntLit, Value: "2", Line: 1, Column: 9},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
		pr := parser.parseExpr(LOWEST)
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}
