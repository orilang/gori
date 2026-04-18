package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParser_for_stmt(t *testing.T) {
	assert := assert.New(t)

	t.Run("function_infinite_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for {
   return false
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
      ForStmt
       For: "for" @4:2 (kind=31)
        BlockStmt
         LBrace: "{" @4:6 (kind=41)
         Stmts
          ReturnStmt
           Values
            BoolLitExpr
             Value: "false" @5:11 (kind=7)
         RBrace: "}" @6:2 (kind=42)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_infinite_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for {}
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
      ForStmt
       For: "for" @4:2 (kind=31)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_condition_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for a < b {
   return false
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
      ForStmt
       For: "for" @4:2 (kind=31)
       Condition
        BinaryExpr
         IdentExpr
          Name: "a" @4:6 (kind=3)
         Operator: "<" @4:8 (kind=64)
         IdentExpr
          Name: "b" @4:10 (kind=3)
        BlockStmt
         LBrace: "{" @4:12 (kind=41)
         Stmts
          ReturnStmt
           Values
            BoolLitExpr
             Value: "false" @5:11 (kind=7)
         RBrace: "}" @6:2 (kind=42)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_condition_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for a < b {}
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
      ForStmt
       For: "for" @4:2 (kind=31)
       Condition
        BinaryExpr
         IdentExpr
          Name: "a" @4:6 (kind=3)
         Operator: "<" @4:8 (kind=64)
         IdentExpr
          Name: "b" @4:10 (kind=3)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_all_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for a := 0;a<5;a+=1 {
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
      ForStmt
       For: "for" @4:2 (kind=31)
       Init
        AssignStmt
         Left
          IdentExpr
           Name: "a" @4:6 (kind=3)
         Operator: ":=" @4:8 (kind=50)
         Right
          IntLitExpr
           Value: "0" @4:11 (kind=4)
       Condition
        BinaryExpr
         IdentExpr
          Name: "a" @4:13 (kind=3)
         Operator: "<" @4:14 (kind=64)
         IntLitExpr
          Value: "5" @4:15 (kind=4)
       Post
        AssignStmt
         Left
          IdentExpr
           Name: "a" @4:17 (kind=3)
         Operator: "+=" @4:18 (kind=52)
         Right
          IntLitExpr
           Value: "1" @4:20 (kind=4)
        BlockStmt
         LBrace: "{" @4:22 (kind=41)
         Stmts
          ReturnStmt
           Values
            IdentExpr
             Name: "a" @5:11 (kind=3)
         RBrace: "}" @6:2 (kind=42)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_all_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for a := 0;a<5;a+=1 {}
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
      ForStmt
       For: "for" @4:2 (kind=31)
       Init
        AssignStmt
         Left
          IdentExpr
           Name: "a" @4:6 (kind=3)
         Operator: ":=" @4:8 (kind=50)
         Right
          IntLitExpr
           Value: "0" @4:11 (kind=4)
       Condition
        BinaryExpr
         IdentExpr
          Name: "a" @4:13 (kind=3)
         Operator: "<" @4:14 (kind=64)
         IntLitExpr
          Value: "5" @4:15 (kind=4)
       Post
        AssignStmt
         Left
          IdentExpr
           Name: "a" @4:17 (kind=3)
         Operator: "+=" @4:18 (kind=52)
         Right
          IntLitExpr
           Value: "1" @4:20 (kind=4)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_all_x3", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for a:=0;a<5;a++ {
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
      ForStmt
       For: "for" @4:2 (kind=31)
       Init
        AssignStmt
         Left
          IdentExpr
           Name: "a" @4:6 (kind=3)
         Operator: ":=" @4:7 (kind=50)
         Right
          IntLitExpr
           Value: "0" @4:9 (kind=4)
       Condition
        BinaryExpr
         IdentExpr
          Name: "a" @4:11 (kind=3)
         Operator: "<" @4:12 (kind=64)
         IntLitExpr
          Value: "5" @4:13 (kind=4)
       Post
        IncDecStmt
         X:
          IdentExpr
           Name: "a" @4:15 (kind=3)
         Operator: "++" @4:16 (kind=53)
        BlockStmt
         LBrace: "{" @4:19 (kind=41)
         Stmts
          ReturnStmt
           Values
            IdentExpr
             Name: "a" @5:11 (kind=3)
         RBrace: "}" @6:2 (kind=42)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_range_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for i := range 5 {
   z = i
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
      RangeStmt
       For: "for" @4:2 (kind=31)
       Key
        IdentExpr
         Name: "i" @4:6 (kind=3)
       Op: ":=" @4:8 (kind=50)
       Range: "range" @4:11 (kind=71)
        IntLitExpr
         Value: "5" @4:17 (kind=4)
        BlockStmt
         LBrace: "{" @4:19 (kind=41)
         Stmts
          AssignStmt
           Left
            IdentExpr
             Name: "z" @5:4 (kind=3)
           Operator: "=" @5:6 (kind=49)
           Right
            IdentExpr
             Name: "i" @5:8 (kind=3)
         RBrace: "}" @6:2 (kind=42)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_range_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for k,v := range 5 {
   z = i
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
      RangeStmt
       For: "for" @4:2 (kind=31)
       Key
        IdentExpr
         Name: "k" @4:6 (kind=3)
       Condition
        IdentExpr
         Name: "v" @4:8 (kind=3)
       Op: ":=" @4:10 (kind=50)
       Range: "range" @4:13 (kind=71)
        IntLitExpr
         Value: "5" @4:19 (kind=4)
        BlockStmt
         LBrace: "{" @4:21 (kind=41)
         Stmts
          AssignStmt
           Left
            IdentExpr
             Name: "z" @5:4 (kind=3)
           Operator: "=" @5:6 (kind=49)
           Right
            IdentExpr
             Name: "i" @5:8 (kind=3)
         RBrace: "}" @6:2 (kind=42)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_range_x3", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for i := range 5 {}
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
      RangeStmt
       For: "for" @4:2 (kind=31)
       Key
        IdentExpr
         Name: "i" @4:6 (kind=3)
       Op: ":=" @4:8 (kind=50)
       Range: "range" @4:11 (kind=71)
        IntLitExpr
         Value: "5" @4:17 (kind=4)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_range_x4", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for k,v := range 5 {}
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
      RangeStmt
       For: "for" @4:2 (kind=31)
       Key
        IdentExpr
         Name: "k" @4:6 (kind=3)
       Condition
        IdentExpr
         Name: "v" @4:8 (kind=3)
       Op: ":=" @4:10 (kind=50)
       Range: "range" @4:13 (kind=71)
        IntLitExpr
         Value: "5" @4:19 (kind=4)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_range_x5", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for range 5 {}
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
      RangeStmt
       For: "for" @4:2 (kind=31)
       Op: "" @0:0 (kind=0)
       Range: "range" @4:6 (kind=71)
        IntLitExpr
         Value: "5" @4:12 (kind=4)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_range_x6", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for range ch {
   return x
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
      RangeStmt
       For: "for" @4:2 (kind=31)
       Op: "" @0:0 (kind=0)
       Range: "range" @4:6 (kind=71)
        IdentExpr
         Name: "ch" @4:12 (kind=3)
        BlockStmt
         LBrace: "{" @4:15 (kind=41)
         Stmts
          ReturnStmt
           Values
            IdentExpr
             Name: "x" @5:11 (kind=3)
         RBrace: "}" @6:2 (kind=42)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_range_x7", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for k = range 5 {}
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
      RangeStmt
       For: "for" @4:2 (kind=31)
       Key
        IdentExpr
         Name: "k" @4:6 (kind=3)
       Op: "=" @4:8 (kind=49)
       Range: "range" @4:10 (kind=71)
        IntLitExpr
         Value: "5" @4:16 (kind=4)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_range_x8", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for _,v := range 5 {
   z = i
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
      RangeStmt
       For: "for" @4:2 (kind=31)
       Key
        IdentExpr
         Name: "_" @4:6 (kind=3)
       Condition
        IdentExpr
         Name: "v" @4:8 (kind=3)
       Op: ":=" @4:10 (kind=50)
       Range: "range" @4:13 (kind=71)
        IntLitExpr
         Value: "5" @4:19 (kind=4)
        BlockStmt
         LBrace: "{" @4:21 (kind=41)
         Stmts
          AssignStmt
           Left
            IdentExpr
             Name: "z" @5:4 (kind=3)
           Operator: "=" @5:6 (kind=49)
           Right
            IdentExpr
             Name: "i" @5:8 (kind=3)
         RBrace: "}" @6:2 (kind=42)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("bad_infinite_x1", func(t *testing.T) {
		// infinity loop
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for [ }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("function_break_continue_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for {
   break
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
      ForStmt
       For: "for" @4:2 (kind=31)
        BlockStmt
         LBrace: "{" @4:6 (kind=41)
         Stmts
          BreakStmt
           Break: "break" @5:4 (kind=32)
         RBrace: "}" @6:2 (kind=42)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_break_continue_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for {
   continue
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
      ForStmt
       For: "for" @4:2 (kind=31)
        BlockStmt
         LBrace: "{" @4:6 (kind=41)
         Stmts
          ContinueStmt
           Continue: "continue" @5:4 (kind=33)
         RBrace: "}" @6:2 (kind=42)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_break_continue_x3", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for {
   if a {
     break
	}
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
      ForStmt
       For: "for" @4:2 (kind=31)
        BlockStmt
         LBrace: "{" @4:6 (kind=41)
         Stmts
          IfStmt
           Condition
            IdentExpr
             Name: "a" @5:7 (kind=3)
           Then
            BlockStmt
             LBrace: "{" @5:9 (kind=41)
             Stmts
              BreakStmt
               Break: "break" @6:6 (kind=32)
             RBrace: "}" @7:2 (kind=42)
         RBrace: "}" @8:2 (kind=42)
     RBrace: "}" @9:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_break_continue_x4", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 if a {
   for {
     break
	 }
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
        IdentExpr
         Name: "a" @4:5 (kind=3)
       Then
        BlockStmt
         LBrace: "{" @4:7 (kind=41)
         Stmts
          ForStmt
           For: "for" @5:4 (kind=31)
            BlockStmt
             LBrace: "{" @5:8 (kind=41)
             Stmts
              BreakStmt
               Break: "break" @6:6 (kind=32)
             RBrace: "}" @7:3 (kind=42)
         RBrace: "}" @8:2 (kind=42)
     RBrace: "}" @9:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_break_continue_x5", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for {
   if a {
     continue
	 }
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
      ForStmt
       For: "for" @4:2 (kind=31)
        BlockStmt
         LBrace: "{" @4:6 (kind=41)
         Stmts
          IfStmt
           Condition
            IdentExpr
             Name: "a" @5:7 (kind=3)
           Then
            BlockStmt
             LBrace: "{" @5:9 (kind=41)
             Stmts
              ContinueStmt
               Continue: "continue" @6:6 (kind=33)
             RBrace: "}" @7:3 (kind=42)
         RBrace: "}" @8:2 (kind=42)
     RBrace: "}" @9:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_break_continue_x6", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for {
   if a {
     continue
		 w()
	 }
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
      ForStmt
       For: "for" @4:2 (kind=31)
        BlockStmt
         LBrace: "{" @4:6 (kind=41)
         Stmts
          IfStmt
           Condition
            IdentExpr
             Name: "a" @5:7 (kind=3)
           Then
            BlockStmt
             LBrace: "{" @5:9 (kind=41)
             Stmts
              ContinueStmt
               Continue: "continue" @6:6 (kind=33)
              CallExpr
               Callee
                IdentExpr
                 Name: "w" @7:4 (kind=3)
               LParent: "(" @7:5 (kind=39)
               RParent: ")" @7:6 (kind=40)
             RBrace: "}" @8:3 (kind=42)
         RBrace: "}" @9:2 (kind=42)
     RBrace: "}" @10:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_break_continue_x7", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for {
   if a {
     continue;w()
	 }
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
      ForStmt
       For: "for" @4:2 (kind=31)
        BlockStmt
         LBrace: "{" @4:6 (kind=41)
         Stmts
          IfStmt
           Condition
            IdentExpr
             Name: "a" @5:7 (kind=3)
           Then
            BlockStmt
             LBrace: "{" @5:9 (kind=41)
             Stmts
              ContinueStmt
               Continue: "continue" @6:6 (kind=33)
              CallExpr
               Callee
                IdentExpr
                 Name: "w" @6:15 (kind=3)
               LParent: "(" @6:16 (kind=39)
               RParent: ")" @6:17 (kind=40)
             RBrace: "}" @7:3 (kind=42)
         RBrace: "}" @8:2 (kind=42)
     RBrace: "}" @9:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_break_continue_x8", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for {
   if a {
     break
		 w()
	 }
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
      ForStmt
       For: "for" @4:2 (kind=31)
        BlockStmt
         LBrace: "{" @4:6 (kind=41)
         Stmts
          IfStmt
           Condition
            IdentExpr
             Name: "a" @5:7 (kind=3)
           Then
            BlockStmt
             LBrace: "{" @5:9 (kind=41)
             Stmts
              BreakStmt
               Break: "break" @6:6 (kind=32)
              CallExpr
               Callee
                IdentExpr
                 Name: "w" @7:4 (kind=3)
               LParent: "(" @7:5 (kind=39)
               RParent: ")" @7:6 (kind=40)
             RBrace: "}" @8:3 (kind=42)
         RBrace: "}" @9:2 (kind=42)
     RBrace: "}" @10:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_break_continue_x9", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for {
   if a {
     break;w()
	 }
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
      ForStmt
       For: "for" @4:2 (kind=31)
        BlockStmt
         LBrace: "{" @4:6 (kind=41)
         Stmts
          IfStmt
           Condition
            IdentExpr
             Name: "a" @5:7 (kind=3)
           Then
            BlockStmt
             LBrace: "{" @5:9 (kind=41)
             Stmts
              BreakStmt
               Break: "break" @6:6 (kind=32)
              CallExpr
               Callee
                IdentExpr
                 Name: "w" @6:12 (kind=3)
               LParent: "(" @6:13 (kind=39)
               RParent: ")" @6:14 (kind=40)
             RBrace: "}" @7:3 (kind=42)
         RBrace: "}" @8:2 (kind=42)
     RBrace: "}" @9:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_break_continue_x10", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for {
   if a {
     break // comment
		 ;w()
	 }
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
      ForStmt
       For: "for" @4:2 (kind=31)
        BlockStmt
         LBrace: "{" @4:6 (kind=41)
         Stmts
          IfStmt
           Condition
            IdentExpr
             Name: "a" @5:7 (kind=3)
           Then
            BlockStmt
             LBrace: "{" @5:9 (kind=41)
             Stmts
              BreakStmt
               Break: "break" @6:6 (kind=32)
              CallExpr
               Callee
                IdentExpr
                 Name: "w" @7:5 (kind=3)
               LParent: "(" @7:6 (kind=39)
               RParent: ")" @7:7 (kind=40)
             RBrace: "}" @8:3 (kind=42)
         RBrace: "}" @9:2 (kind=42)
     RBrace: "}" @10:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_break_continue_x11", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for {
   if a {
     break /*
		 test
		 */
		 ;w()
	 }
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
      ForStmt
       For: "for" @4:2 (kind=31)
        BlockStmt
         LBrace: "{" @4:6 (kind=41)
         Stmts
          IfStmt
           Condition
            IdentExpr
             Name: "a" @5:7 (kind=3)
           Then
            BlockStmt
             LBrace: "{" @5:9 (kind=41)
             Stmts
              BreakStmt
               Break: "break" @6:6 (kind=32)
              CallExpr
               Callee
                IdentExpr
                 Name: "w" @9:5 (kind=3)
               LParent: "(" @9:6 (kind=39)
               RParent: ")" @9:7 (kind=40)
             RBrace: "}" @10:3 (kind=42)
         RBrace: "}" @11:2 (kind=42)
     RBrace: "}" @12:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_break_continue_x12", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for {
   if a {
     continue // comment
		 ;w()
	 }
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
      ForStmt
       For: "for" @4:2 (kind=31)
        BlockStmt
         LBrace: "{" @4:6 (kind=41)
         Stmts
          IfStmt
           Condition
            IdentExpr
             Name: "a" @5:7 (kind=3)
           Then
            BlockStmt
             LBrace: "{" @5:9 (kind=41)
             Stmts
              ContinueStmt
               Continue: "continue" @6:6 (kind=33)
              CallExpr
               Callee
                IdentExpr
                 Name: "w" @7:5 (kind=3)
               LParent: "(" @7:6 (kind=39)
               RParent: ")" @7:7 (kind=40)
             RBrace: "}" @8:3 (kind=42)
         RBrace: "}" @9:2 (kind=42)
     RBrace: "}" @10:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("function_break_continue_x13", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for {
   if a {
     continue /*
		 test
		 */
		 ;w()
	 }
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
      ForStmt
       For: "for" @4:2 (kind=31)
        BlockStmt
         LBrace: "{" @4:6 (kind=41)
         Stmts
          IfStmt
           Condition
            IdentExpr
             Name: "a" @5:7 (kind=3)
           Then
            BlockStmt
             LBrace: "{" @5:9 (kind=41)
             Stmts
              ContinueStmt
               Continue: "continue" @6:6 (kind=33)
              CallExpr
               Callee
                IdentExpr
                 Name: "w" @9:5 (kind=3)
               LParent: "(" @9:6 (kind=39)
               RParent: ")" @9:7 (kind=40)
             RBrace: "}" @10:3 (kind=42)
         RBrace: "}" @11:2 (kind=42)
     RBrace: "}" @12:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("bad_condition_x1", func(t *testing.T) {
		// condition
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for ,a < b;{
   return false
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_condition_x2", func(t *testing.T) {
		// condition
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for a < b;{
   return false
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_condition_x3", func(t *testing.T) {
		// condition
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for a < b [
   return false
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_condition_x4", func(t *testing.T) {
		// condition
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for a:=0; a < 5; a++ b{
   return false
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_range_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for i:= range {
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_range_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for ,i :=range {
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_range_x3", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for true :=range {
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_range_x4", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for i, true :=range {
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_range_x5", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for k,v :=range 5 5 {
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_range_x6", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for k +=range 5{
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_range_x7", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for range ch [
   return x
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_range_x8", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for range {
   return x
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_range_x9", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for k,v +=range 5 {
   return false
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_range_x10", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for k,v +=range {
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_range_x11", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for k,v =range {
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_range_x12", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for i :=range [
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_range_x13", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 break
 for {
   return a
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_all_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for a :=0;a<5;+=1 [
   return a
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_all_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for a:=0;a<5; {
   return a
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_all_x3", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for a()++;a<5; {
   return a
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_all_x4", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for a :=0, a<5; {
   return a
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_all_x5", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for a :=0;a<5, {
   return a
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_all_x6", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for a :=0;a<5;5 {
   return a
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_all_x7", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for a :=0;a<5;var {
   return a
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_break_continue_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 continue
 for {
   return false
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_break_continue_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 if a {
   break
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_break_continue_x3", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 if a {
   continue
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_break_continue_x4", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for {
   func inner(){
	   break
	 }
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_break_continue_x5", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for {
   if a {
	   break
		 x
	 }
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_break_continue_x6", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for {
   if a {
	   continue
		 x
	 }
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_break_continue_x7", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for {
   break /*test*/ x = 1
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_break_continue_x8", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for {
   continue /*test*/ x = 1
 }
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_identifier_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
 for 123abcd,v := range 5 {}
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})
}
