package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParser_return_stmt(t *testing.T) {
	assert := assert.New(t)

	t.Run("function_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  return
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
   Name: "x" @3:6 (kind=3)
   Params
    (none)
   Body
    BlockStmt
     LBrace: "{" @3:9 (kind=41)
     Stmts
      ReturnStmt
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  return a
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
   Name: "x" @3:6 (kind=3)
   Params
    (none)
   Body
    BlockStmt
     LBrace: "{" @3:9 (kind=41)
     Stmts
      ReturnStmt
       Values
        IdentExpr
         Name: "a" @4:10 (kind=3)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_x3", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  a:="test"
  return a
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
   Name: "x" @3:6 (kind=3)
   Params
    (none)
   Body
    BlockStmt
     LBrace: "{" @3:9 (kind=41)
     Stmts
      AssignStmt
       Left
        IdentExpr
         Name: "a" @4:3 (kind=3)
       Operator: ":=" @4:4 (kind=50)
       Right
        StringLitExpr
         Value: "test" @4:6 (kind=6)
      ReturnStmt
       Values
        IdentExpr
         Name: "a" @5:10 (kind=3)
     RBrace: "}" @6:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_x4", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  return a,b
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
   Name: "x" @3:6 (kind=3)
   Params
    (none)
   Body
    BlockStmt
     LBrace: "{" @3:9 (kind=41)
     Stmts
      ReturnStmt
       Values
        IdentExpr
         Name: "a" @4:10 (kind=3)
        IdentExpr
         Name: "b" @4:12 (kind=3)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("bad_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  return a b
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

func x(){
  return a,b,
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

func x(){
  return ,b
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

func x(){
 return _
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})
}
