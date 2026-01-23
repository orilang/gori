package ast

import (
	"bytes"
	"testing"

	"github.com/orilang/gori/token"
	"github.com/stretchr/testify/assert"
)

func TestAst_dump(t *testing.T) {
	assert := assert.New(t)

	t.Run("default_nil", func(t *testing.T) {
		result := Dump(nil)
		assert.NotNil(result)
	})

	t.Run("default_not_nil", func(t *testing.T) {
		result := Dump("z")
		assert.NotNil(result)
	})

	t.Run("nil_func_decl_body", func(t *testing.T) {
		f := &FuncDecl{
			Body: nil,
		}
		result := Dump(f)
		assert.NotNil(result)
	})

	t.Run("nil_const_type", func(t *testing.T) {
		f := &ConstDeclStmt{
			ConstKW: token.Token{Kind: token.KWConst, Value: "const"},
			Name:    token.Token{Kind: token.Ident, Value: "x"},
			Type:    nil,
		}
		result := Dump(f)
		assert.NotNil(result)
	})

	t.Run("nil_var_type", func(t *testing.T) {
		f := &VarDeclStmt{
			VarKW: token.Token{Kind: token.KWVar, Value: "var"},
			Name:  token.Token{Kind: token.Ident, Value: "x"},
			Type:  nil,
		}
		result := Dump(f)
		assert.NotNil(result)
	})

	t.Run("nil_parent_expr", func(t *testing.T) {
		f := &ParenExpr{
			Left:  token.Token{Kind: token.LParen, Value: "("},
			Right: token.Token{Kind: token.RParen, Value: ")"},
			Inner: nil,
		}
		result := Dump(f)
		assert.NotNil(result)
	})

	t.Run("bad_type_from", func(t *testing.T) {
		f := &ConstDeclStmt{
			ConstKW: token.Token{Kind: token.KWConst, Value: "const"},
			Name:    token.Token{Kind: token.Ident, Value: "x"},
			Type:    &BadType{From: token.Token{Kind: token.Not, Line: 1, Column: 9, Value: "!"}},
		}
		result := Dump(f)
		assert.NotNil(result)
		assert.Contains(result, "BadType")
	})

	t.Run("bad_type_from_to", func(t *testing.T) {
		f := &ConstDeclStmt{
			ConstKW: token.Token{Kind: token.KWConst, Value: "const"},
			Name:    token.Token{Kind: token.Ident, Value: "x"},
			Type: &BadType{
				From: token.Token{Kind: token.Not, Line: 1, Column: 9, Value: "!"},
				To:   token.Token{Kind: token.Not, Line: 1, Column: 11, Value: "!"},
			},
		}
		result := Dump(f)
		assert.NotNil(result)
		assert.Contains(result, "BadType")
	})

	t.Run("token", func(t *testing.T) {
		result := Dump(token.Token{Kind: token.KWConst, Value: "const"})
		assert.NotNil(result)
	})

	t.Run("decl_nil", func(t *testing.T) {
		var b bytes.Buffer
		d := dumper{w: &b}
		d.decl(nil, 0)
		assert.Equal("(nil decl)\n", b.String())
	})

	t.Run("decl_not_nil", func(t *testing.T) {
		var b bytes.Buffer
		d := dumper{w: &b}
		d.decl(&dumpType{}, 0)
		assert.Equal("<<unhandled decl *ast.dumpType>>\n", b.String())
	})

	t.Run("typ_nil", func(t *testing.T) {
		var b bytes.Buffer
		d := dumper{w: &b}
		d.typ(nil, 0)
		assert.Equal("(nil type)\n", b.String())
	})

	t.Run("typ_not_nil", func(t *testing.T) {
		var b bytes.Buffer
		d := dumper{w: &b}
		d.typ(&dumpType{}, 0)
		assert.Equal("<<unhandled type *ast.dumpType>>\n", b.String())
	})

	t.Run("stmt_nil", func(t *testing.T) {
		var b bytes.Buffer
		d := dumper{w: &b}
		d.stmt(nil, 0)
		assert.Equal("(nil stmt)\n", b.String())
	})

	t.Run("stmt_not_nil", func(t *testing.T) {
		var b bytes.Buffer
		d := dumper{w: &b}
		d.stmt(&dumpType{}, 0)
		assert.Equal("<<unhandled stmt *ast.dumpType>>\n", b.String())
	})

	t.Run("expr_nil", func(t *testing.T) {
		var b bytes.Buffer
		d := dumper{w: &b}
		d.expr(nil, 0)
		assert.Equal("(nil expr)\n", b.String())
	})

	t.Run("expr_not_nil", func(t *testing.T) {
		var b bytes.Buffer
		d := dumper{w: &b}
		d.expr(&dumpType{}, 0)
		assert.Equal("<<unhandled expr *ast.dumpType>>\n", b.String())
	})

	t.Run("bad_expr_from", func(t *testing.T) {
		f := &BadExpr{
			From:   token.Token{Kind: token.LParen, Value: "(", Line: 1, Column: 1},
			Reason: "unexpected expression",
		}
		result := Dump(f)
		assert.NotNil(result)
		assert.Contains(result, "BadExpr")
	})

	t.Run("bad_expr_from_to", func(t *testing.T) {
		f := &BadExpr{
			From:   token.Token{Kind: token.LParen, Value: "(", Line: 1, Column: 1},
			To:     token.Token{Kind: token.RParen, Value: ")", Line: 1, Column: 2},
			Reason: "expected expression inside parentheses",
		}
		result := Dump(f)
		assert.NotNil(result)
		assert.Contains(result, "BadExpr")
	})

	t.Run("bad_stmt_from", func(t *testing.T) {
		f := &BadStmt{
			From:   token.Token{Kind: token.Not, Value: "!", Line: 1, Column: 1},
			Reason: "unexpected statement",
		}
		result := Dump(f)
		assert.NotNil(result)
		assert.Contains(result, "BadStmt")
	})

	t.Run("bad_stmt_from_to", func(t *testing.T) {
		f := &BadStmt{
			From:   token.Token{Kind: token.Not, Value: "!", Line: 1, Column: 1},
			To:     token.Token{Kind: token.Assign, Value: "=", Line: 1, Column: 5},
			Reason: "unexpected statement",
		}
		result := Dump(f)
		assert.NotNil(result)
		assert.Contains(result, "BadStmt")
	})

	t.Run("bad_decl_from", func(t *testing.T) {
		f := &BadDecl{
			From:   token.Token{Kind: token.Not, Value: "!", Line: 1, Column: 1},
			Reason: "unexpected statement",
		}
		result := Dump(f)
		assert.NotNil(result)
		assert.Contains(result, "BadDecl")
	})

	t.Run("bad_decl_from_to", func(t *testing.T) {
		f := &BadDecl{
			From:   token.Token{Kind: token.Not, Value: "!", Line: 1, Column: 1},
			To:     token.Token{Kind: token.Assign, Value: "=", Line: 1, Column: 5},
			Reason: "unexpected statement",
		}
		result := Dump(f)
		assert.NotNil(result)
		assert.Contains(result, "BadDecl")
	})
}
