package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParser_if_stmt(t *testing.T) {
	assert := assert.New(t)

	t.Run("function_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 if a < b {
   return b
 } else if a < c {
	 return c
 } else {
   return a
 }
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
      IfStmt
       Condition
        BinaryExpr
         IdentExpr
          Name: "a" @4:5 (kind=3)
         Operator: "<" @4:7 (kind=64)
         IdentExpr
          Name: "b" @4:9 (kind=3)
       Then
        BlockStmt
         LBrace: "{" @4:11 (kind=41)
         Stmts
          ReturnStmt
           Values
            IdentExpr
             Name: "b" @5:11 (kind=3)
         RBrace: "}" @6:2 (kind=42)
       Else
        IfStmt
         Condition
          BinaryExpr
           IdentExpr
            Name: "a" @6:12 (kind=3)
           Operator: "<" @6:14 (kind=64)
           IdentExpr
            Name: "c" @6:16 (kind=3)
         Then
          BlockStmt
           LBrace: "{" @6:18 (kind=41)
           Stmts
            ReturnStmt
             Values
              IdentExpr
               Name: "c" @7:10 (kind=3)
           RBrace: "}" @8:2 (kind=42)
         Else
          BlockStmt
           LBrace: "{" @8:9 (kind=41)
           Stmts
            ReturnStmt
             Values
              IdentExpr
               Name: "a" @9:11 (kind=3)
           RBrace: "}" @10:2 (kind=42)
     RBrace: "}" @11:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 if a < b {
   if a < c {
	   return c
	 }
 }
}`
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
      IfStmt
       Condition
        BinaryExpr
         IdentExpr
          Name: "a" @4:5 (kind=3)
         Operator: "<" @4:7 (kind=64)
         IdentExpr
          Name: "b" @4:9 (kind=3)
       Then
        BlockStmt
         LBrace: "{" @4:11 (kind=41)
         Stmts
          IfStmt
           Condition
            BinaryExpr
             IdentExpr
              Name: "a" @5:7 (kind=3)
             Operator: "<" @5:9 (kind=64)
             IdentExpr
              Name: "c" @5:11 (kind=3)
           Then
            BlockStmt
             LBrace: "{" @5:13 (kind=41)
             Stmts
              ReturnStmt
               Values
                IdentExpr
                 Name: "c" @6:12 (kind=3)
             RBrace: "}" @7:3 (kind=42)
         RBrace: "}" @8:2 (kind=42)
     RBrace: "}" @9:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("bad_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 if a < b
}
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

func x(){
 if
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

func x(){
 if {
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

func x(){
 else {
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_x5", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 if a {
   else {}
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_x6", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 if a = 1 {}
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_x7", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 if a < b {
   return b
 }
 else return c
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_x8", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 if _ {
   else {
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}
