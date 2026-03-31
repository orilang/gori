package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParser_parse_const(t *testing.T) {
	assert := assert.New(t)

	t.Run("float_lit", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `const a float = 3.14
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseConstDecl()
		result := `ConstDecl
 Const: "const" @1:1 (kind=23)
 Name: "a" @1:7 (kind=3)
 Type
  NameType
   Name: "float" @1:9 (kind=20)
 Eq: "=" @1:15 (kind=49)
 Init
  FloatLitExpr
   Value: "3.14" @1:17 (kind=5)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("const_bad_type", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `const a func = 3.14
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseConstDecl()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("const_bad_type_blank_identifier_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `const _ func = 3.14
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseConstDecl()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("const_bad_type_blank_identifier_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `const a@ string = "YES"
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseConstDecl()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}

func TestParser_parse_var(t *testing.T) {
	assert := assert.New(t)

	t.Run("int_lit", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `var a int = 0
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseVarDecl()
		result := `VarDeclStmt
 Var: "var" @1:1 (kind=11)
 Name: "a" @1:5 (kind=3)
 Type
  NameType
   Name: "int" @1:7 (kind=12)
 Eq: "=" @1:11 (kind=49)
 Init
  IntLitExpr
   Value: "0" @1:13 (kind=4)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("float_lit", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `var a float = 3.14
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseVarDecl()
		result := `VarDeclStmt
 Var: "var" @1:1 (kind=11)
 Name: "a" @1:5 (kind=3)
 Type
  NameType
   Name: "float" @1:7 (kind=20)
 Eq: "=" @1:13 (kind=49)
 Init
  FloatLitExpr
   Value: "3.14" @1:15 (kind=5)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("var_bad_type", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `var a func = 3.14
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseVarDecl()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("var_bad_type_blank_identifier", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `var _ func = 3.14
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseVarDecl()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bool_lit", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `var a bool = true
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseVarDecl()
		result := `VarDeclStmt
 Var: "var" @1:1 (kind=11)
 Name: "a" @1:5 (kind=3)
 Type
  NameType
   Name: "bool" @1:7 (kind=25)
 Eq: "=" @1:12 (kind=49)
 Init
  BoolLitExpr
   Value: "true" @1:14 (kind=7)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("string_lit", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `var a string = "ok"
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseVarDecl()
		result := `VarDeclStmt
 Var: "var" @1:1 (kind=11)
 Name: "a" @1:5 (kind=3)
 Type
  NameType
   Name: "string" @1:7 (kind=24)
 Eq: "=" @1:14 (kind=49)
 Init
  StringLitExpr
   Value: "ok" @1:16 (kind=6)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("indent_lit", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `var a string = x
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseVarDecl()
		result := `VarDeclStmt
 Var: "var" @1:1 (kind=11)
 Name: "a" @1:5 (kind=3)
 Type
  NameType
   Name: "string" @1:7 (kind=24)
 Eq: "=" @1:14 (kind=49)
 Init
  IdentExpr
   Name: "x" @1:16 (kind=3)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})
}
