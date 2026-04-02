package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParser_type_enum_decl(t *testing.T) {
	assert := assert.New(t)

	t.Run("x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type Color enum {
  Red;Blue;Green;Yellow
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  EnumDecl:
   Type: "type" @3:1 (kind=26)
    Name: "Color" @3:6 (kind=3)
    Public: true
    Enum: "enum" @3:12 (kind=74)
    LBrace: "{" @3:17 (kind=41)
     Variants
      Ident: "Red" @4:3 (kind=3)
      Ident: "Blue" @4:7 (kind=3)
      Ident: "Green" @4:12 (kind=3)
      Ident: "Yellow" @4:18 (kind=3)
    RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type Color enum {
  Red
	Blue
	Green
	Yellow // comment
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  EnumDecl:
   Type: "type" @3:1 (kind=26)
    Name: "Color" @3:6 (kind=3)
    Public: true
    Enum: "enum" @3:12 (kind=74)
    LBrace: "{" @3:17 (kind=41)
     Variants
      Ident: "Red" @4:3 (kind=3)
      Ident: "Blue" @5:2 (kind=3)
      Ident: "Green" @6:2 (kind=3)
      Ident: "Yellow" @7:2 (kind=3)
    RBrace: "}" @8:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("bad_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type Color enum {}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type Color enum {
  string
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_x3", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type Color enum {
  Red Green
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_x4", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type _Color enum {
  Red;Blue;Green;Yellow
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}
