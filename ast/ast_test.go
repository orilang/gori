package ast

import (
	"testing"

	"github.com/orilang/gori/token"
	"github.com/stretchr/testify/assert"
)

func TestAst_position(t *testing.T) {
	assert := assert.New(t)

	t.Run("ident_expr", func(t *testing.T) {
		z := token.Token{
			Kind:   token.Ident,
			Value:  "a",
			Line:   1,
			Column: 1,
		}
		x := &IdentExpr{z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("name_type", func(t *testing.T) {
		z := token.Token{
			Kind:   token.KWInt,
			Value:  "int",
			Line:   1,
			Column: 1,
		}
		x := &NameType{z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("int_lit_expr", func(t *testing.T) {
		z := token.Token{
			Kind:   token.IntLit,
			Value:  "1",
			Line:   1,
			Column: 1,
		}
		x := &IntLitExpr{z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("float_lit_expr", func(t *testing.T) {
		z := token.Token{
			Kind:   token.FloatLit,
			Value:  "3.14",
			Line:   1,
			Column: 1,
		}
		x := &FloatLitExpr{z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("bool_lit_expr", func(t *testing.T) {
		z := token.Token{
			Kind:   token.BoolLit,
			Value:  "true",
			Line:   1,
			Column: 1,
		}
		x := &BoolLitExpr{z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("string_lit_expr", func(t *testing.T) {
		z := token.Token{
			Kind:   token.StringLit,
			Value:  "string",
			Line:   1,
			Column: 1,
		}
		x := &StringLitExpr{z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("parent_expr", func(t *testing.T) {
		z := token.Token{
			Kind:   token.Ident,
			Value:  "a",
			Line:   1,
			Column: 1,
		}
		x := &ParenExpr{Left: z, Right: z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("binary_expr", func(t *testing.T) {
		z := token.Token{
			Kind:   token.Ident,
			Value:  "a",
			Line:   1,
			Column: 1,
		}
		x := &BinaryExpr{Left: &StringLitExpr{z}, Right: &StringLitExpr{z}}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("unary_expr", func(t *testing.T) {
		z := token.Token{
			Kind:   token.Plus,
			Value:  "+",
			Line:   1,
			Column: 1,
		}
		x := &UnaryExpr{Operator: z, Right: &StringLitExpr{z}}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("selector_expr", func(t *testing.T) {
		z := token.Token{
			Kind:   token.Ident,
			Value:  "a",
			Line:   1,
			Column: 1,
		}
		x := &SelectorExpr{X: &IdentExpr{z}, Selector: z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("index_expr", func(t *testing.T) {
		z := token.Token{
			Kind:   token.Ident,
			Value:  "a",
			Line:   1,
			Column: 1,
		}
		x := &IndexExpr{X: &IdentExpr{z}, RBracket: z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("call_expr", func(t *testing.T) {
		z := token.Token{
			Kind:   token.Ident,
			Value:  "a",
			Line:   1,
			Column: 1,
		}
		x := &CallExpr{Callee: &IdentExpr{z}, RParen: z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("expr_stmt", func(t *testing.T) {
		z := token.Token{
			Kind:   token.Ident,
			Value:  "a",
			Line:   1,
			Column: 1,
		}
		x := &ExprStmt{&IdentExpr{z}}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("bad_expr", func(t *testing.T) {
		z := token.Token{
			Kind:   token.Ident,
			Value:  "a",
			Line:   1,
			Column: 1,
		}
		x := &BadExpr{From: z, To: z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("dump_type", func(t *testing.T) {
		x := &dumpType{S: ""}
		assert.Equal(token.Token{}, x.Start())
		assert.Equal(token.Token{}, x.End())
	})
}
