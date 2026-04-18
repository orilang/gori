package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParser_switch_stmt(t *testing.T) {
	assert := assert.New(t)

	t.Run("empty_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch {}
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       LBrace: "{" @4:10 (kind=41)
       RBrace: "}" @4:11 (kind=42)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("no_tag_empty_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch {
    default:
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       LBrace: "{" @4:10 (kind=41)
       Case: "default" @5:5 (kind=36)
       Colon: ":" @5:12 (kind=47)
       RBrace: "}" @6:3 (kind=42)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("no_tag_empty_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch {
    default:
    case 1:
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       LBrace: "{" @4:10 (kind=41)
       Case: "default" @5:5 (kind=36)
       Colon: ":" @5:12 (kind=47)
       Case: "case" @6:5 (kind=35)
        Values:
         IntLitExpr
          Value: "1" @6:10 (kind=4)
       Colon: ":" @6:11 (kind=47)
       RBrace: "}" @7:3 (kind=42)
     RBrace: "}" @8:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("no_tag_empty_x3", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch {
    case 1:
    default:
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       LBrace: "{" @4:10 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         IntLitExpr
          Value: "1" @5:10 (kind=4)
       Colon: ":" @5:11 (kind=47)
       Case: "default" @6:5 (kind=36)
       Colon: ":" @6:12 (kind=47)
       RBrace: "}" @7:3 (kind=42)
     RBrace: "}" @8:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("no_tag_default", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch {
    default:
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       LBrace: "{" @4:10 (kind=41)
       Case: "default" @5:5 (kind=36)
       Colon: ":" @5:12 (kind=47)
        Body:
         ReturnStmt
          Values
           IdentExpr
            Name: "a" @6:14 (kind=3)
       RBrace: "}" @7:3 (kind=42)
     RBrace: "}" @8:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("no_tag_cases_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch {
    case 1:
      b()
    default:
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       LBrace: "{" @4:10 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         IntLitExpr
          Value: "1" @5:10 (kind=4)
       Colon: ":" @5:11 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "b" @6:7 (kind=3)
          LParent: "(" @6:8 (kind=39)
          RParent: ")" @6:9 (kind=40)
       Case: "default" @7:5 (kind=36)
       Colon: ":" @7:12 (kind=47)
        Body:
         ReturnStmt
          Values
           IdentExpr
            Name: "a" @8:14 (kind=3)
       RBrace: "}" @9:3 (kind=42)
     RBrace: "}" @10:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("no_tag_cases_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch {
    case 1:
      return b
    default:
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       LBrace: "{" @4:10 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         IntLitExpr
          Value: "1" @5:10 (kind=4)
       Colon: ":" @5:11 (kind=47)
        Body:
         ReturnStmt
          Values
           IdentExpr
            Name: "b" @6:14 (kind=3)
       Case: "default" @7:5 (kind=36)
       Colon: ":" @7:12 (kind=47)
        Body:
         ReturnStmt
          Values
           IdentExpr
            Name: "a" @8:14 (kind=3)
       RBrace: "}" @9:3 (kind=42)
     RBrace: "}" @10:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("no_tag_cases_x3", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch {
    case 1:
    case 2:
      b()
      c()
    default:
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       LBrace: "{" @4:10 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         IntLitExpr
          Value: "1" @5:10 (kind=4)
       Colon: ":" @5:11 (kind=47)
       Case: "case" @6:5 (kind=35)
        Values:
         IntLitExpr
          Value: "2" @6:10 (kind=4)
       Colon: ":" @6:11 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "b" @7:7 (kind=3)
          LParent: "(" @7:8 (kind=39)
          RParent: ")" @7:9 (kind=40)
         CallExpr
          Callee
           IdentExpr
            Name: "c" @8:7 (kind=3)
          LParent: "(" @8:8 (kind=39)
          RParent: ")" @8:9 (kind=40)
       Case: "default" @9:5 (kind=36)
       Colon: ":" @9:12 (kind=47)
        Body:
         ReturnStmt
          Values
           IdentExpr
            Name: "a" @10:14 (kind=3)
       RBrace: "}" @11:3 (kind=42)
     RBrace: "}" @12:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("no_tag_cases_x4", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch {
    case 1,2:
      b()
      c()
    default:
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       LBrace: "{" @4:10 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         IntLitExpr
          Value: "1" @5:10 (kind=4)
         IntLitExpr
          Value: "2" @5:12 (kind=4)
       Colon: ":" @5:13 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "b" @6:7 (kind=3)
          LParent: "(" @6:8 (kind=39)
          RParent: ")" @6:9 (kind=40)
         CallExpr
          Callee
           IdentExpr
            Name: "c" @7:7 (kind=3)
          LParent: "(" @7:8 (kind=39)
          RParent: ")" @7:9 (kind=40)
       Case: "default" @8:5 (kind=36)
       Colon: ":" @8:12 (kind=47)
        Body:
         ReturnStmt
          Values
           IdentExpr
            Name: "a" @9:14 (kind=3)
       RBrace: "}" @10:3 (kind=42)
     RBrace: "}" @11:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("no_tag_cases_x5", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch {
    case 1:
      b()
      fallthrough
    case 2:
      c()
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       LBrace: "{" @4:10 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         IntLitExpr
          Value: "1" @5:10 (kind=4)
       Colon: ":" @5:11 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "b" @6:7 (kind=3)
          LParent: "(" @6:8 (kind=39)
          RParent: ")" @6:9 (kind=40)
         FallThroughStmt
          FallThrough: "fallthrough" @7:7 (kind=37)
       Case: "case" @8:5 (kind=35)
        Values:
         IntLitExpr
          Value: "2" @8:10 (kind=4)
       Colon: ":" @8:11 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "c" @9:7 (kind=3)
          LParent: "(" @9:8 (kind=39)
          RParent: ")" @9:9 (kind=40)
       RBrace: "}" @10:3 (kind=42)
     RBrace: "}" @11:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("no_tag_cases_x6", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch {
    case 1:
      b()
      fallthrough // fallthrough
    case 2:
      c()
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       LBrace: "{" @4:10 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         IntLitExpr
          Value: "1" @5:10 (kind=4)
       Colon: ":" @5:11 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "b" @6:7 (kind=3)
          LParent: "(" @6:8 (kind=39)
          RParent: ")" @6:9 (kind=40)
         FallThroughStmt
          FallThrough: "fallthrough" @7:7 (kind=37)
       Case: "case" @8:5 (kind=35)
        Values:
         IntLitExpr
          Value: "2" @8:10 (kind=4)
       Colon: ":" @8:11 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "c" @9:7 (kind=3)
          LParent: "(" @9:8 (kind=39)
          RParent: ")" @9:9 (kind=40)
       RBrace: "}" @10:3 (kind=42)
     RBrace: "}" @11:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("no_tag_cases_x7", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch {
    case a > b:
      b()
      fallthrough // fallthrough
    case 2:
      c()
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       LBrace: "{" @4:10 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         BinaryExpr
          IdentExpr
           Name: "a" @5:10 (kind=3)
          Operator: ">" @5:12 (kind=66)
          IdentExpr
           Name: "b" @5:14 (kind=3)
       Colon: ":" @5:15 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "b" @6:7 (kind=3)
          LParent: "(" @6:8 (kind=39)
          RParent: ")" @6:9 (kind=40)
         FallThroughStmt
          FallThrough: "fallthrough" @7:7 (kind=37)
       Case: "case" @8:5 (kind=35)
        Values:
         IntLitExpr
          Value: "2" @8:10 (kind=4)
       Colon: ":" @8:11 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "c" @9:7 (kind=3)
          LParent: "(" @9:8 (kind=39)
          RParent: ")" @9:9 (kind=40)
       RBrace: "}" @10:3 (kind=42)
     RBrace: "}" @11:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("tag_cases_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch z {
    case a > b:
      b()
    case 2:
      c()
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       Init:
        IdentExpr
         Name: "z" @4:10 (kind=3)
       LBrace: "{" @4:12 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         BinaryExpr
          IdentExpr
           Name: "a" @5:10 (kind=3)
          Operator: ">" @5:12 (kind=66)
          IdentExpr
           Name: "b" @5:14 (kind=3)
       Colon: ":" @5:15 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "b" @6:7 (kind=3)
          LParent: "(" @6:8 (kind=39)
          RParent: ")" @6:9 (kind=40)
       Case: "case" @7:5 (kind=35)
        Values:
         IntLitExpr
          Value: "2" @7:10 (kind=4)
       Colon: ":" @7:11 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "c" @8:7 (kind=3)
          LParent: "(" @8:8 (kind=39)
          RParent: ")" @8:9 (kind=40)
       RBrace: "}" @9:3 (kind=42)
     RBrace: "}" @10:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("tag_cases_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch z {
    case d:
    case a > b:
      b()
    case 2:
      c()
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       Init:
        IdentExpr
         Name: "z" @4:10 (kind=3)
       LBrace: "{" @4:12 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         IdentExpr
          Name: "d" @5:10 (kind=3)
       Colon: ":" @5:11 (kind=47)
       Case: "case" @6:5 (kind=35)
        Values:
         BinaryExpr
          IdentExpr
           Name: "a" @6:10 (kind=3)
          Operator: ">" @6:12 (kind=66)
          IdentExpr
           Name: "b" @6:14 (kind=3)
       Colon: ":" @6:15 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "b" @7:7 (kind=3)
          LParent: "(" @7:8 (kind=39)
          RParent: ")" @7:9 (kind=40)
       Case: "case" @8:5 (kind=35)
        Values:
         IntLitExpr
          Value: "2" @8:10 (kind=4)
       Colon: ":" @8:11 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "c" @9:7 (kind=3)
          LParent: "(" @9:8 (kind=39)
          RParent: ")" @9:9 (kind=40)
       RBrace: "}" @10:3 (kind=42)
     RBrace: "}" @11:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("tag_cases_x3", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch z {
    default:
    case a > b:
      b()
    case 2:
      c()
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       Init:
        IdentExpr
         Name: "z" @4:10 (kind=3)
       LBrace: "{" @4:12 (kind=41)
       Case: "default" @5:5 (kind=36)
       Colon: ":" @5:12 (kind=47)
       Case: "case" @6:5 (kind=35)
        Values:
         BinaryExpr
          IdentExpr
           Name: "a" @6:10 (kind=3)
          Operator: ">" @6:12 (kind=66)
          IdentExpr
           Name: "b" @6:14 (kind=3)
       Colon: ":" @6:15 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "b" @7:7 (kind=3)
          LParent: "(" @7:8 (kind=39)
          RParent: ")" @7:9 (kind=40)
       Case: "case" @8:5 (kind=35)
        Values:
         IntLitExpr
          Value: "2" @8:10 (kind=4)
       Colon: ":" @8:11 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "c" @9:7 (kind=3)
          LParent: "(" @9:8 (kind=39)
          RParent: ")" @9:9 (kind=40)
       RBrace: "}" @10:3 (kind=42)
     RBrace: "}" @11:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("tag_cases_x4", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch z {
    default:
    case a > b:
      b()
    case 2:
      return c()
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       Init:
        IdentExpr
         Name: "z" @4:10 (kind=3)
       LBrace: "{" @4:12 (kind=41)
       Case: "default" @5:5 (kind=36)
       Colon: ":" @5:12 (kind=47)
       Case: "case" @6:5 (kind=35)
        Values:
         BinaryExpr
          IdentExpr
           Name: "a" @6:10 (kind=3)
          Operator: ">" @6:12 (kind=66)
          IdentExpr
           Name: "b" @6:14 (kind=3)
       Colon: ":" @6:15 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "b" @7:7 (kind=3)
          LParent: "(" @7:8 (kind=39)
          RParent: ")" @7:9 (kind=40)
       Case: "case" @8:5 (kind=35)
        Values:
         IntLitExpr
          Value: "2" @8:10 (kind=4)
       Colon: ":" @8:11 (kind=47)
        Body:
         ReturnStmt
          Values
           CallExpr
            Callee
             IdentExpr
              Name: "c" @9:14 (kind=3)
            LParent: "(" @9:15 (kind=39)
            RParent: ")" @9:16 (kind=40)
       RBrace: "}" @10:3 (kind=42)
     RBrace: "}" @11:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("tag_cases_x5", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch z {
    case 1+2*3:
      return 7
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       Init:
        IdentExpr
         Name: "z" @4:10 (kind=3)
       LBrace: "{" @4:12 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         BinaryExpr
          IntLitExpr
           Value: "1" @5:10 (kind=4)
          Operator: "+" @5:11 (kind=51)
          BinaryExpr
           IntLitExpr
            Value: "2" @5:12 (kind=4)
           Operator: "*" @5:13 (kind=57)
           IntLitExpr
            Value: "3" @5:14 (kind=4)
       Colon: ":" @5:15 (kind=47)
        Body:
         ReturnStmt
          Values
           IntLitExpr
            Value: "7" @6:14 (kind=4)
       RBrace: "}" @7:3 (kind=42)
     RBrace: "}" @8:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("tag_cases_x6", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch z {
    case 1:
      switch y {
        case 2:
          return 7
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       Init:
        IdentExpr
         Name: "z" @4:10 (kind=3)
       LBrace: "{" @4:12 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         IntLitExpr
          Value: "1" @5:10 (kind=4)
       Colon: ":" @5:11 (kind=47)
        Body:
         SwitchStmt
          Switch: "switch" @6:7 (kind=34)
          Init:
           IdentExpr
            Name: "y" @6:14 (kind=3)
          LBrace: "{" @6:16 (kind=41)
          Case: "case" @7:9 (kind=35)
           Values:
            IntLitExpr
             Value: "2" @7:14 (kind=4)
          Colon: ":" @7:15 (kind=47)
           Body:
            ReturnStmt
             Values
              IntLitExpr
               Value: "7" @8:18 (kind=4)
          RBrace: "}" @9:7 (kind=42)
       RBrace: "}" @10:3 (kind=42)
     RBrace: "}" @11:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("init_tag_cases_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch z:=w();z {
    case a>b:
      b()
    case 2:
      c()
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       Init:
        AssignStmt
         Left
          IdentExpr
           Name: "z" @4:10 (kind=3)
         Operator: ":=" @4:11 (kind=50)
         Right
          CallExpr
           Callee
            IdentExpr
             Name: "w" @4:13 (kind=3)
           LParent: "(" @4:14 (kind=39)
           RParent: ")" @4:15 (kind=40)
       Init:
        IdentExpr
         Name: "z" @4:17 (kind=3)
       LBrace: "{" @4:19 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         BinaryExpr
          IdentExpr
           Name: "a" @5:10 (kind=3)
          Operator: ">" @5:11 (kind=66)
          IdentExpr
           Name: "b" @5:12 (kind=3)
       Colon: ":" @5:13 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "b" @6:7 (kind=3)
          LParent: "(" @6:8 (kind=39)
          RParent: ")" @6:9 (kind=40)
       Case: "case" @7:5 (kind=35)
        Values:
         IntLitExpr
          Value: "2" @7:10 (kind=4)
       Colon: ":" @7:11 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "c" @8:7 (kind=3)
          LParent: "(" @8:8 (kind=39)
          RParent: ")" @8:9 (kind=40)
       RBrace: "}" @9:2 (kind=42)
     RBrace: "}" @10:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("init_tag_cases_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch z:=w();z {
    default:
    case a>b:
      b()
    case 2:
      c()
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       Init:
        AssignStmt
         Left
          IdentExpr
           Name: "z" @4:10 (kind=3)
         Operator: ":=" @4:11 (kind=50)
         Right
          CallExpr
           Callee
            IdentExpr
             Name: "w" @4:13 (kind=3)
           LParent: "(" @4:14 (kind=39)
           RParent: ")" @4:15 (kind=40)
       Init:
        IdentExpr
         Name: "z" @4:17 (kind=3)
       LBrace: "{" @4:19 (kind=41)
       Case: "default" @5:5 (kind=36)
       Colon: ":" @5:12 (kind=47)
       Case: "case" @6:5 (kind=35)
        Values:
         BinaryExpr
          IdentExpr
           Name: "a" @6:10 (kind=3)
          Operator: ">" @6:11 (kind=66)
          IdentExpr
           Name: "b" @6:12 (kind=3)
       Colon: ":" @6:13 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "b" @7:7 (kind=3)
          LParent: "(" @7:8 (kind=39)
          RParent: ")" @7:9 (kind=40)
       Case: "case" @8:5 (kind=35)
        Values:
         IntLitExpr
          Value: "2" @8:10 (kind=4)
       Colon: ":" @8:11 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "c" @9:7 (kind=3)
          LParent: "(" @9:8 (kind=39)
          RParent: ")" @9:9 (kind=40)
       RBrace: "}" @10:2 (kind=42)
     RBrace: "}" @11:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("init_tag_cases_x3", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch z:=w();z {
    case a>b:
    default:
      b()
    case 2:
      c()
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       Init:
        AssignStmt
         Left
          IdentExpr
           Name: "z" @4:10 (kind=3)
         Operator: ":=" @4:11 (kind=50)
         Right
          CallExpr
           Callee
            IdentExpr
             Name: "w" @4:13 (kind=3)
           LParent: "(" @4:14 (kind=39)
           RParent: ")" @4:15 (kind=40)
       Init:
        IdentExpr
         Name: "z" @4:17 (kind=3)
       LBrace: "{" @4:19 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         BinaryExpr
          IdentExpr
           Name: "a" @5:10 (kind=3)
          Operator: ">" @5:11 (kind=66)
          IdentExpr
           Name: "b" @5:12 (kind=3)
       Colon: ":" @5:13 (kind=47)
       Case: "default" @6:5 (kind=36)
       Colon: ":" @6:12 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "b" @7:7 (kind=3)
          LParent: "(" @7:8 (kind=39)
          RParent: ")" @7:9 (kind=40)
       Case: "case" @8:5 (kind=35)
        Values:
         IntLitExpr
          Value: "2" @8:10 (kind=4)
       Colon: ":" @8:11 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "c" @9:7 (kind=3)
          LParent: "(" @9:8 (kind=39)
          RParent: ")" @9:9 (kind=40)
       RBrace: "}" @10:2 (kind=42)
     RBrace: "}" @11:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("init_tag_cases_x4", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch z:=w();z {
    case a>b:
    default:
    case 2:
      c()
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       Init:
        AssignStmt
         Left
          IdentExpr
           Name: "z" @4:10 (kind=3)
         Operator: ":=" @4:11 (kind=50)
         Right
          CallExpr
           Callee
            IdentExpr
             Name: "w" @4:13 (kind=3)
           LParent: "(" @4:14 (kind=39)
           RParent: ")" @4:15 (kind=40)
       Init:
        IdentExpr
         Name: "z" @4:17 (kind=3)
       LBrace: "{" @4:19 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         BinaryExpr
          IdentExpr
           Name: "a" @5:10 (kind=3)
          Operator: ">" @5:11 (kind=66)
          IdentExpr
           Name: "b" @5:12 (kind=3)
       Colon: ":" @5:13 (kind=47)
       Case: "default" @6:5 (kind=36)
       Colon: ":" @6:12 (kind=47)
       Case: "case" @7:5 (kind=35)
        Values:
         IntLitExpr
          Value: "2" @7:10 (kind=4)
       Colon: ":" @7:11 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "c" @8:7 (kind=3)
          LParent: "(" @8:8 (kind=39)
          RParent: ")" @8:9 (kind=40)
       RBrace: "}" @9:2 (kind=42)
     RBrace: "}" @10:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("bad_no_tag_cases_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch {
    case 1,2:
      b()
      c()
    default:
    default:
	}
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_no_tag_cases_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch {
    case ,1,2:
      b()
      c()
    default:
      return a
	}
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_no_tag_cases_x3", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch {
    case 1,2,:
      b()
      c()
    default:
      return a
	}
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_no_tag_cases_x4", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch {
    case 1,2::
      b()
      c()
    default:
      return a
	}
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_no_tag_cases_x5", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch {
    case 1,2::
      b()
      c()
    default a
	}
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_no_tag_cases_x6", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch {
    case 1,2::
      b()
      c()
    default
	}
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_no_tag_cases_x7", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch {
    case :
      b()
      c()
	}
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("no_tag_cases_x8", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch {
    case 1:
      b()
      fallthrough
      fallthrough
    case 2:
      c()
	}
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_no_tag_cases_x9", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch {
    case 1,2,:
      b()
      c()
    default
	}
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_no_tag_cases_x10", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch {
    case 1,,2:
      b()
      c()
    default
	}
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_no_tag_cases_x11", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch {
    case _:
      b()
    default:
      return a
	}
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_tag_cases_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch z {
    case a>b b:
      b()
    case 2:
      c()
	}
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_init_tag_cases_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch z:=w() {
    case a>b:
      b()
    case 2:
      c()
	}
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_init_tag_cases_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  switch z:=w();z {
	case a>b;,:
      b()
    case 2:
      c()
	}
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})
}
