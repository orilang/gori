package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParser_type_defined_decl(t *testing.T) {
	assert := assert.New(t)

	t.Run("x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type a int
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  DefinedTypeDecl:
   TypeDecl: "type" @3:1 (kind=26)
    Name: "a" @3:6 (kind=3)
    Type
     NamedType
      Ident: "int" @3:8 (kind=12)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type a int
type b int
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  DefinedTypeDecl:
   TypeDecl: "type" @3:1 (kind=26)
    Name: "a" @3:6 (kind=3)
    Type
     NamedType
      Ident: "int" @3:8 (kind=12)
  DefinedTypeDecl:
   TypeDecl: "type" @4:1 (kind=26)
    Name: "b" @4:6 (kind=3)
    Type
     NamedType
      Ident: "int" @4:8 (kind=12)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("x3", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type a int;type b int
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  DefinedTypeDecl:
   TypeDecl: "type" @3:1 (kind=26)
    Name: "a" @3:6 (kind=3)
    Type
     NamedType
      Ident: "int" @3:8 (kind=12)
  DefinedTypeDecl:
   TypeDecl: "type" @3:12 (kind=26)
    Name: "b" @3:17 (kind=3)
    Type
     NamedType
      Ident: "int" @3:19 (kind=12)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("x4", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func main(){
	type a int
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  FuncDecl
   Function: "func" @3:1 (kind=10)
   Name: "main" @3:6 (kind=3)
   Params
    (none)
   Body
    BlockStmt
     LBrace: "{" @3:12 (kind=41)
     Stmts
      DefinedTypeDecl:
       TypeDecl: "type" @4:2 (kind=26)
        Name: "a" @4:7 (kind=3)
        Type
         NamedType
          Ident: "int" @4:9 (kind=12)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})
}
