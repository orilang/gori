package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParser_type_struct_decl(t *testing.T) {
	assert := assert.New(t)

	t.Run("empty_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test struct{}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  StructDecl:
   Type: "type" @3:1 (kind=26)
   Name: "test" @3:6 (kind=3)
   Struct: "struct" @3:11 (kind=27)
   LBrace: "{" @3:17 (kind=41)
   RBrace: "}" @3:18 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("empty_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type Test struct{}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  StructDecl:
   Type: "type" @3:1 (kind=26)
   Name: "Test" @3:6 (kind=3)
   Struct: "struct" @3:11 (kind=27)
   Public: true
   LBrace: "{" @3:17 (kind=41)
   RBrace: "}" @3:18 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test struct{
  x int
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  StructDecl:
   Type: "type" @3:1 (kind=26)
   Name: "test" @3:6 (kind=3)
   Struct: "struct" @3:11 (kind=27)
   LBrace: "{" @3:17 (kind=41)
    Name: "x" @4:3 (kind=3)
    Type:
     NamedType
      Ident: "int" @4:5 (kind=12)
   RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test struct{x int}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  StructDecl:
   Type: "type" @3:1 (kind=26)
   Name: "test" @3:6 (kind=3)
   Struct: "struct" @3:11 (kind=27)
   LBrace: "{" @3:17 (kind=41)
    Name: "x" @3:18 (kind=3)
    Type:
     NamedType
      Ident: "int" @3:20 (kind=12)
   RBrace: "}" @3:23 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("x3", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test struct{x int;Y string}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  StructDecl:
   Type: "type" @3:1 (kind=26)
   Name: "test" @3:6 (kind=3)
   Struct: "struct" @3:11 (kind=27)
   LBrace: "{" @3:17 (kind=41)
    Name: "x" @3:18 (kind=3)
    Type:
     NamedType
      Ident: "int" @3:20 (kind=12)
    Name: "Y" @3:24 (kind=3)
    Public: true
    Type:
     NamedType
      Ident: "string" @3:26 (kind=24)
   RBrace: "}" @3:32 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("x4", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test struct{x int
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  StructDecl:
   Type: "type" @3:1 (kind=26)
   Name: "test" @3:6 (kind=3)
   Struct: "struct" @3:11 (kind=27)
   LBrace: "{" @3:17 (kind=41)
    Name: "x" @3:18 (kind=3)
    Type:
     NamedType
      Ident: "int" @3:20 (kind=12)
   RBrace: "}" @4:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("x5", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test struct { x int // comment 
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  StructDecl:
   Type: "type" @3:1 (kind=26)
   Name: "test" @3:6 (kind=3)
   Struct: "struct" @3:11 (kind=27)
   LBrace: "{" @3:18 (kind=41)
    Name: "x" @3:20 (kind=3)
    Type:
     NamedType
      Ident: "int" @3:22 (kind=12)
   RBrace: "}" @4:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("x6", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test struct{x int = 5}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  StructDecl:
   Type: "type" @3:1 (kind=26)
   Name: "test" @3:6 (kind=3)
   Struct: "struct" @3:11 (kind=27)
   LBrace: "{" @3:17 (kind=41)
    Name: "x" @3:18 (kind=3)
    Type:
     NamedType
      Ident: "int" @3:20 (kind=12)
    Eq: "=" @3:24 (kind=49)
    IntLitExpr
     Value: "5" @3:26 (kind=4)
   RBrace: "}" @3:27 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("x7", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test struct{
  x int = 5
  y int = 5
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  StructDecl:
   Type: "type" @3:1 (kind=26)
   Name: "test" @3:6 (kind=3)
   Struct: "struct" @3:11 (kind=27)
   LBrace: "{" @3:17 (kind=41)
    Name: "x" @4:3 (kind=3)
    Type:
     NamedType
      Ident: "int" @4:5 (kind=12)
    Eq: "=" @4:9 (kind=49)
    IntLitExpr
     Value: "5" @4:11 (kind=4)
    Name: "y" @5:3 (kind=3)
    Type:
     NamedType
      Ident: "int" @5:5 (kind=12)
    Eq: "=" @5:9 (kind=49)
    IntLitExpr
     Value: "5" @5:11 (kind=4)
   RBrace: "}" @6:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("x8", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test struct{
  x int = join(a, b)
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  StructDecl:
   Type: "type" @3:1 (kind=26)
   Name: "test" @3:6 (kind=3)
   Struct: "struct" @3:11 (kind=27)
   LBrace: "{" @3:17 (kind=41)
    Name: "x" @4:3 (kind=3)
    Type:
     NamedType
      Ident: "int" @4:5 (kind=12)
    Eq: "=" @4:9 (kind=49)
    CallExpr
     Callee
      IdentExpr
       Name: "join" @4:11 (kind=3)
     LParent: "(" @4:15 (kind=39)
     Args:
      IdentExpr
       Name: "a" @4:16 (kind=3)
      IdentExpr
       Name: "b" @4:19 (kind=3)
     RParent: ")" @4:20 (kind=40)
   RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("bad_type", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test structt{}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test struct{x}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test struct{x int,}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_x3", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test struct{x int;x}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_x4", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test struct{x int =}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_x5", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type _test struct {}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})
}
