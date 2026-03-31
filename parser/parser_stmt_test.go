package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParser_stmt(t *testing.T) {
	assert := assert.New(t)

	t.Run("function_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  a=1
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
       Operator: "=" @4:4 (kind=49)
       Right
        IntLitExpr
         Value: "1" @4:5 (kind=4)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(a z){
  a.b+=1
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
    Param
     Ident: "a" @3:8 (kind=3)
     Type
      NameType
       Name: "z" @3:10 (kind=3)
   Body
    BlockStmt
     LBrace: "{" @3:12 (kind=41)
     Stmts
      AssignStmt
       Left
        SelectorExpr
         X:
          IdentExpr
           Name: "a" @4:3 (kind=3)
         Dot: "." @4:4 (kind=48)
         Selector: "b" @4:5 (kind=3)
       Operator: "+=" @4:6 (kind=52)
       Right
        IntLitExpr
         Value: "1" @4:8 (kind=4)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_x3", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(a z){
  a[i]=3
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
    Param
     Ident: "a" @3:8 (kind=3)
     Type
      NameType
       Name: "z" @3:10 (kind=3)
   Body
    BlockStmt
     LBrace: "{" @3:12 (kind=41)
     Stmts
      AssignStmt
       Left
        IndexExpr
         X:
         IdentExpr
          Name: "a" @4:3 (kind=3)
         LBracket: "[" @4:4 (kind=43)
          IdentExpr
           Name: "i" @4:5 (kind=3)
         RBracket: "]" @4:6 (kind=44)
       Operator: "=" @4:7 (kind=49)
       Right
        IntLitExpr
         Value: "3" @4:8 (kind=4)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_x4", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  a=f()
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
       Operator: "=" @4:4 (kind=49)
       Right
        CallExpr
         Callee
          IdentExpr
           Name: "f" @4:5 (kind=3)
         LParent: "(" @4:6 (kind=39)
         RParent: ")" @4:7 (kind=40)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_x5", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  f()
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
      CallExpr
       Callee
        IdentExpr
         Name: "f" @4:3 (kind=3)
       LParent: "(" @4:4 (kind=39)
       RParent: ")" @4:5 (kind=40)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_x6", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(a []int){
  f()
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
    Param
     Ident: "a" @3:8 (kind=3)
     Type
        LBracket: "[" @3:10 (kind=43)
        RBracket: "]" @3:11 (kind=44)
        Ident: "int" @3:12 (kind=12)
   Body
    BlockStmt
     LBrace: "{" @3:16 (kind=41)
     Stmts
      CallExpr
       Callee
        IdentExpr
         Name: "f" @4:3 (kind=3)
       LParent: "(" @4:4 (kind=39)
       RParent: ")" @4:5 (kind=40)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_x7", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(a [5]int){
  f()
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
    Param
     Ident: "a" @3:8 (kind=3)
     Type
        LBracket: "[" @3:10 (kind=43)
        Size: "5" @3:11 (kind=4)
        RBracket: "]" @3:12 (kind=44)
        Ident: "int" @3:13 (kind=12)
   Body
    BlockStmt
     LBrace: "{" @3:17 (kind=41)
     Stmts
      CallExpr
       Callee
        IdentExpr
         Name: "f" @4:3 (kind=3)
       LParent: "(" @4:4 (kind=39)
       RParent: ")" @4:5 (kind=40)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_x8", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(_ z){
  a[i]=3
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
    Param
     Ident: "_" @3:8 (kind=3)
     Type
      NameType
       Name: "z" @3:10 (kind=3)
   Body
    BlockStmt
     LBrace: "{" @3:12 (kind=41)
     Stmts
      AssignStmt
       Left
        IndexExpr
         X:
         IdentExpr
          Name: "a" @4:3 (kind=3)
         LBracket: "[" @4:4 (kind=43)
          IdentExpr
           Name: "i" @4:5 (kind=3)
         RBracket: "]" @4:6 (kind=44)
       Operator: "=" @4:7 (kind=49)
       Right
        IntLitExpr
         Value: "3" @4:8 (kind=4)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_bad_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  a+b
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("function_bad_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  !a
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}
