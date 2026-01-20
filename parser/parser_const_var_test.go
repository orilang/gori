package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
	"github.com/stretchr/testify/assert"
)

func TestParser_parse_const(t *testing.T) {
	assert := assert.New(t)
	t.Run("float_lit", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWConst, Value: "const"},
			{Kind: token.Ident, Value: "a"},
			{Kind: token.KWFloat, Value: "float"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.FloatLit, Value: "3.14"},
			{Kind: token.EOF, Value: ""},
		}

		parser := New(input)
		pr := parser.parseConstDecl()
		result := `ConstDecl
 Const: "const" @0:0 (kind=23)
 Name: "a" @0:0 (kind=3)
 Type
  NameType
   Name: "float" @0:0 (kind=20)
 Eq: "=" @0:0 (kind=49)
 Init
  FloatLitExpr
   Value: "3.14" @0:0 (kind=5)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("const_bad_type", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWConst, Value: "const"},
			{Kind: token.Ident, Value: "a"},
			{Kind: token.KWFunc, Value: "func"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.FloatLit, Value: "3.14"},
			{Kind: token.EOF, Value: ""},
		}

		parser := New(input)
		pr := parser.parseConstDecl()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}
func TestParser_parse_var(t *testing.T) {
	assert := assert.New(t)

	t.Run("int_lit", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "a"},
			{Kind: token.KWInt, Value: "int"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "0"},
			{Kind: token.EOF, Value: ""},
		}

		parser := New(input)
		pr := parser.parseVarDecl()
		result := `VarDeclStmt
 Var: "var" @0:0 (kind=11)
 Name: "a" @0:0 (kind=3)
 Type
  NameType
   Name: "int" @0:0 (kind=12)
 Eq: "=" @0:0 (kind=49)
 Init
  IntLitExpr
   Value: "0" @0:0 (kind=4)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("float_lit", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "a"},
			{Kind: token.KWFloat, Value: "float"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.FloatLit, Value: "3.14"},
			{Kind: token.EOF, Value: ""},
		}

		parser := New(input)
		pr := parser.parseVarDecl()
		result := `VarDeclStmt
 Var: "var" @0:0 (kind=11)
 Name: "a" @0:0 (kind=3)
 Type
  NameType
   Name: "float" @0:0 (kind=20)
 Eq: "=" @0:0 (kind=49)
 Init
  FloatLitExpr
   Value: "3.14" @0:0 (kind=5)
`
		assert.Equal(result, ast.Dump(pr))
	})

	t.Run("var_bad_type", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "a"},
			{Kind: token.KWFunc, Value: "func"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.FloatLit, Value: "3.14"},
			{Kind: token.EOF, Value: ""},
		}

		parser := New(input)
		pr := parser.parseVarDecl()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bool_lit", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "a"},
			{Kind: token.KWBool, Value: "bool"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.BoolLit, Value: "true"},
			{Kind: token.EOF, Value: ""},
		}

		parser := New(input)
		pr := parser.parseVarDecl()
		result := `VarDeclStmt
 Var: "var" @0:0 (kind=11)
 Name: "a" @0:0 (kind=3)
 Type
  NameType
   Name: "bool" @0:0 (kind=25)
 Eq: "=" @0:0 (kind=49)
 Init
  BoolLitExpr
   Value: "true" @0:0 (kind=7)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("string_lit", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "a"},
			{Kind: token.KWString, Value: "string"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.StringLit, Value: "ok"},
			{Kind: token.EOF, Value: ""},
		}

		parser := New(input)
		pr := parser.parseVarDecl()
		result := `VarDeclStmt
 Var: "var" @0:0 (kind=11)
 Name: "a" @0:0 (kind=3)
 Type
  NameType
   Name: "string" @0:0 (kind=24)
 Eq: "=" @0:0 (kind=49)
 Init
  StringLitExpr
   Value: "ok" @0:0 (kind=6)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("indent_lit", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "a"},
			{Kind: token.KWFloat, Value: "float"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.Ident, Value: "x"},
			{Kind: token.EOF, Value: ""},
		}

		parser := New(input)
		pr := parser.parseVarDecl()
		result := `VarDeclStmt
 Var: "var" @0:0 (kind=11)
 Name: "a" @0:0 (kind=3)
 Type
  NameType
   Name: "float" @0:0 (kind=20)
 Eq: "=" @0:0 (kind=49)
 Init
  IdentExpr
   Name: "x" @0:0 (kind=3)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})
}
