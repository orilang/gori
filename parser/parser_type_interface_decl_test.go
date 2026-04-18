package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParser_type_interface_decl(t *testing.T) {
	assert := assert.New(t)

	t.Run("empty_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test interface{}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  InterfaceDecl:
   Type: "type" @3:1 (kind=26)
   Name: "test" @3:6 (kind=3)
   Interface: "interface" @3:11 (kind=28)
   LBrace: "{" @3:20 (kind=41)
   RBrace: "}" @3:21 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("empty_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type Test interface{}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  InterfaceDecl:
   Type: "type" @3:1 (kind=26)
   Name: "Test" @3:6 (kind=3)
   Interface: "interface" @3:11 (kind=28)
   Public: true
   LBrace: "{" @3:20 (kind=41)
   RBrace: "}" @3:21 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test interface{
  X() error
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  InterfaceDecl:
   Type: "type" @3:1 (kind=26)
   Name: "test" @3:6 (kind=3)
   Interface: "interface" @3:11 (kind=28)
   LBrace: "{" @3:20 (kind=41)
   Name: "X" @4:3 (kind=3)
   Params
    (none)
   Results
     Param
      Type
       NamedType
        Ident: "error" @4:7 (kind=3)
   RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test interface{ X() error }
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  InterfaceDecl:
   Type: "type" @3:1 (kind=26)
   Name: "test" @3:6 (kind=3)
   Interface: "interface" @3:11 (kind=28)
   LBrace: "{" @3:20 (kind=41)
   Name: "X" @3:22 (kind=3)
   Params
    (none)
   Results
     Param
      Type
       NamedType
        Ident: "error" @3:26 (kind=3)
   RBrace: "}" @3:32 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("x3", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test interface{ X() error;Y(a int) error}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  InterfaceDecl:
   Type: "type" @3:1 (kind=26)
   Name: "test" @3:6 (kind=3)
   Interface: "interface" @3:11 (kind=28)
   LBrace: "{" @3:20 (kind=41)
   Name: "X" @3:22 (kind=3)
   Params
    (none)
   Results
     Param
      Type
       NamedType
        Ident: "error" @3:26 (kind=3)
   Name: "Y" @3:32 (kind=3)
   Params
    Param
     Ident: "a" @3:34 (kind=3)
     Type
      NamedType
       Ident: "int" @3:36 (kind=12)
   Results
     Param
      Type
       NamedType
        Ident: "error" @3:41 (kind=3)
   RBrace: "}" @3:46 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("x4", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test interface{
  X() error
  Y(a int) error
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  InterfaceDecl:
   Type: "type" @3:1 (kind=26)
   Name: "test" @3:6 (kind=3)
   Interface: "interface" @3:11 (kind=28)
   LBrace: "{" @3:20 (kind=41)
   Name: "X" @4:3 (kind=3)
   Params
    (none)
   Results
     Param
      Type
       NamedType
        Ident: "error" @4:7 (kind=3)
   Name: "Y" @5:3 (kind=3)
   Params
    Param
     Ident: "a" @5:5 (kind=3)
     Type
      NamedType
       Ident: "int" @5:7 (kind=12)
   Results
     Param
      Type
       NamedType
        Ident: "error" @5:12 (kind=3)
   RBrace: "}" @6:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("x5", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test interface{
  X() error
  Y(a int) error // comment
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  InterfaceDecl:
   Type: "type" @3:1 (kind=26)
   Name: "test" @3:6 (kind=3)
   Interface: "interface" @3:11 (kind=28)
   LBrace: "{" @3:20 (kind=41)
   Name: "X" @4:3 (kind=3)
   Params
    (none)
   Results
     Param
      Type
       NamedType
        Ident: "error" @4:7 (kind=3)
   Name: "Y" @5:3 (kind=3)
   Params
    Param
     Ident: "a" @5:5 (kind=3)
     Type
      NamedType
       Ident: "int" @5:7 (kind=12)
   Results
     Param
      Type
       NamedType
        Ident: "error" @5:12 (kind=3)
   RBrace: "}" @6:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("x6", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test interface{
  X()(y int,z int)
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  InterfaceDecl:
   Type: "type" @3:1 (kind=26)
   Name: "test" @3:6 (kind=3)
   Interface: "interface" @3:11 (kind=28)
   LBrace: "{" @3:20 (kind=41)
   Name: "X" @4:3 (kind=3)
   Params
    (none)
   Results
    LParent: "(" @4:6 (kind=39)
     Param
      Ident: "y" @4:7 (kind=3)
      Type
       NamedType
        Ident: "int" @4:9 (kind=12)
     Param
      Ident: "z" @4:13 (kind=3)
      Type
       NamedType
        Ident: "int" @4:15 (kind=12)
    RParent: ")" @4:18 (kind=40)
   RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("x7", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test interface{
  X() error;Y(a int)error // comment
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  InterfaceDecl:
   Type: "type" @3:1 (kind=26)
   Name: "test" @3:6 (kind=3)
   Interface: "interface" @3:11 (kind=28)
   LBrace: "{" @3:20 (kind=41)
   Name: "X" @4:3 (kind=3)
   Params
    (none)
   Results
     Param
      Type
       NamedType
        Ident: "error" @4:7 (kind=3)
   Name: "Y" @4:13 (kind=3)
   Params
    Param
     Ident: "a" @4:15 (kind=3)
     Type
      NamedType
       Ident: "int" @4:17 (kind=12)
   Results
     Param
      Type
       NamedType
        Ident: "error" @4:21 (kind=3)
   RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("x8", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test interface{
  X()(a, int) // comment
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  InterfaceDecl:
   Type: "type" @3:1 (kind=26)
   Name: "test" @3:6 (kind=3)
   Interface: "interface" @3:11 (kind=28)
   LBrace: "{" @3:20 (kind=41)
   Name: "X" @4:3 (kind=3)
   Params
    (none)
   Results
    LParent: "(" @4:6 (kind=39)
     Param
      Type
       NamedType
        Ident: "a" @4:7 (kind=3)
     Param
      Type
       NamedType
        Ident: "int" @4:10 (kind=12)
    RParent: ")" @4:13 (kind=40)
   RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("x9", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test interface{
  X.Y // comment
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  InterfaceDecl:
   Type: "type" @3:1 (kind=26)
   Name: "test" @3:6 (kind=3)
   Interface: "interface" @3:11 (kind=28)
   LBrace: "{" @3:20 (kind=41)
    Embeds
     NamedType
      Ident: "X" @4:3 (kind=3)
      Dot: "." @4:4 (kind=48)
      Ident: "Y" @4:5 (kind=3)
   RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("bad_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test interface{
  x x
}
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

type test interface{
  x,
}
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

type test interface{
  X() error
  Y(a int) error,
}
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

type test interface{
  X(,) error
  Y(a int) error
}
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

type test interface{
  X(a int,) error
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_x6", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test interface{
  X(a int x) error
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_x7", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test interface{
  X(a struct) error
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_x8", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test interface{
  X(a int)()
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_x9", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test interface{
  X(a int)(,)
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_x10", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test interface{
  X(a int)(a b,c)
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_x11", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test interface{
  X(a int)(a b,)
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_x12", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test interface{
  X()(a, int int)
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_x13", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test interface{
  X()(a,int,)
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_x14", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test interface{
  X(),
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_x15", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test interface{
  X.
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_x16", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test interface{
  _() error
  Y(a int) error // comment
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})
}
