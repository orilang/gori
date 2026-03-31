package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParser_slice_array_stmt(t *testing.T) {
	assert := assert.New(t)

	t.Run("const_slice_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

const x []int=[]int{1,2,3}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 ConstDecls
  ConstDecl
   Const: "const" @3:1 (kind=23)
   Name: "x" @3:7 (kind=3)
   Type
      LBracket: "[" @3:9 (kind=43)
      RBracket: "]" @3:10 (kind=44)
      Ident: "int" @3:11 (kind=12)
   Eq: "=" @3:14 (kind=49)
   Init
      LBracket: "[" @3:15 (kind=43)
      RBracket: "]" @3:16 (kind=44)
      Ident: "int" @3:17 (kind=12)
    LBrace: "{" @3:20 (kind=41)
     Elements
      IntLitExpr
       Value: "1" @3:21 (kind=4)
      IntLitExpr
       Value: "2" @3:23 (kind=4)
      IntLitExpr
       Value: "3" @3:25 (kind=4)
    RBrace: "}" @3:26 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("const_slice_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func main(){
  const x []int=[]int{1,2,3}
  const y []int=[]int{1,2,3}
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
      ConstDecl
       Const: "const" @4:3 (kind=23)
       Name: "x" @4:9 (kind=3)
       Type
          LBracket: "[" @4:11 (kind=43)
          RBracket: "]" @4:12 (kind=44)
          Ident: "int" @4:13 (kind=12)
       Eq: "=" @4:16 (kind=49)
       Init
          LBracket: "[" @4:17 (kind=43)
          RBracket: "]" @4:18 (kind=44)
          Ident: "int" @4:19 (kind=12)
        LBrace: "{" @4:22 (kind=41)
         Elements
          IntLitExpr
           Value: "1" @4:23 (kind=4)
          IntLitExpr
           Value: "2" @4:25 (kind=4)
          IntLitExpr
           Value: "3" @4:27 (kind=4)
        RBrace: "}" @4:28 (kind=42)
      ConstDecl
       Const: "const" @5:3 (kind=23)
       Name: "y" @5:9 (kind=3)
       Type
          LBracket: "[" @5:11 (kind=43)
          RBracket: "]" @5:12 (kind=44)
          Ident: "int" @5:13 (kind=12)
       Eq: "=" @5:16 (kind=49)
       Init
          LBracket: "[" @5:17 (kind=43)
          RBracket: "]" @5:18 (kind=44)
          Ident: "int" @5:19 (kind=12)
        LBrace: "{" @5:22 (kind=41)
         Elements
          IntLitExpr
           Value: "1" @5:23 (kind=4)
          IntLitExpr
           Value: "2" @5:25 (kind=4)
          IntLitExpr
           Value: "3" @5:27 (kind=4)
        RBrace: "}" @5:28 (kind=42)
     RBrace: "}" @6:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("var_slice_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func main(){
  var x []int=make([]int,10);var y []int=make([]int,10);
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
      VarDeclStmt
       Var: "var" @4:3 (kind=11)
       Name: "x" @4:7 (kind=3)
       Type
          LBracket: "[" @4:9 (kind=43)
          RBracket: "]" @4:10 (kind=44)
          Ident: "int" @4:11 (kind=12)
       Eq: "=" @4:14 (kind=49)
       Init
        Make: "make" @4:15 (kind=3)
        LParen: "(" @4:19 (kind=39)
          LBracket: "[" @4:20 (kind=43)
          RBracket: "]" @4:21 (kind=44)
          Ident: "int" @4:22 (kind=12)
        Size:
         IntLitExpr
          Value: "10" @4:26 (kind=4)
        RParen: ")" @4:28 (kind=40)
      VarDeclStmt
       Var: "var" @4:30 (kind=11)
       Name: "y" @4:34 (kind=3)
       Type
          LBracket: "[" @4:36 (kind=43)
          RBracket: "]" @4:37 (kind=44)
          Ident: "int" @4:38 (kind=12)
       Eq: "=" @4:41 (kind=49)
       Init
        Make: "make" @4:42 (kind=3)
        LParen: "(" @4:46 (kind=39)
          LBracket: "[" @4:47 (kind=43)
          RBracket: "]" @4:48 (kind=44)
          Ident: "int" @4:49 (kind=12)
        Size:
         IntLitExpr
          Value: "10" @4:53 (kind=4)
        RParen: ")" @4:55 (kind=40)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("var_slice_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func main(){
  var x []int=[]int{1,2,3}
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
      VarDeclStmt
       Var: "var" @4:3 (kind=11)
       Name: "x" @4:7 (kind=3)
       Type
          LBracket: "[" @4:9 (kind=43)
          RBracket: "]" @4:10 (kind=44)
          Ident: "int" @4:11 (kind=12)
       Eq: "=" @4:14 (kind=49)
       Init
          LBracket: "[" @4:15 (kind=43)
          RBracket: "]" @4:16 (kind=44)
          Ident: "int" @4:17 (kind=12)
        LBrace: "{" @4:20 (kind=41)
         Elements
          IntLitExpr
           Value: "1" @4:21 (kind=4)
          IntLitExpr
           Value: "2" @4:23 (kind=4)
          IntLitExpr
           Value: "3" @4:25 (kind=4)
        RBrace: "}" @4:26 (kind=42)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("var_slice_x3", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func main(){
  var x []c.d=make([]c.d,10)
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
      VarDeclStmt
       Var: "var" @4:3 (kind=11)
       Name: "x" @4:7 (kind=3)
       Type
          LBracket: "[" @4:9 (kind=43)
          RBracket: "]" @4:10 (kind=44)
          Ident: "c" @4:11 (kind=3)
          Dot: "." @4:12 (kind=48)
          Ident: "d" @4:13 (kind=3)
       Eq: "=" @4:14 (kind=49)
       Init
        Make: "make" @4:15 (kind=3)
        LParen: "(" @4:19 (kind=39)
          LBracket: "[" @4:20 (kind=43)
          RBracket: "]" @4:21 (kind=44)
          Ident: "c" @4:22 (kind=3)
          Dot: "." @4:23 (kind=48)
          Ident: "d" @4:24 (kind=3)
        Size:
         IntLitExpr
          Value: "10" @4:26 (kind=4)
        RParen: ")" @4:28 (kind=40)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("var_slice_x4", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func main(){
  var x []int=[]int{1,2,3}
  var y []int=x[:]
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
      VarDeclStmt
       Var: "var" @4:3 (kind=11)
       Name: "x" @4:7 (kind=3)
       Type
          LBracket: "[" @4:9 (kind=43)
          RBracket: "]" @4:10 (kind=44)
          Ident: "int" @4:11 (kind=12)
       Eq: "=" @4:14 (kind=49)
       Init
          LBracket: "[" @4:15 (kind=43)
          RBracket: "]" @4:16 (kind=44)
          Ident: "int" @4:17 (kind=12)
        LBrace: "{" @4:20 (kind=41)
         Elements
          IntLitExpr
           Value: "1" @4:21 (kind=4)
          IntLitExpr
           Value: "2" @4:23 (kind=4)
          IntLitExpr
           Value: "3" @4:25 (kind=4)
        RBrace: "}" @4:26 (kind=42)
      VarDeclStmt
       Var: "var" @5:3 (kind=11)
       Name: "y" @5:7 (kind=3)
       Type
          LBracket: "[" @5:9 (kind=43)
          RBracket: "]" @5:10 (kind=44)
          Ident: "int" @5:11 (kind=12)
       Eq: "=" @5:14 (kind=49)
       Init
         IdentExpr
          Name: "x" @5:15 (kind=3)
          LBracket: "[" @5:16 (kind=43)
          Colon: ":" @5:17 (kind=47)
          RBracket: "]" @5:18 (kind=44)
     RBrace: "}" @6:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("var_slice_x5", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func main(){
  var x []int=[]int{1,2,3}
  var y view []int=x[3:6];
  var z []int=[]int{1,2,3}
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
      VarDeclStmt
       Var: "var" @4:3 (kind=11)
       Name: "x" @4:7 (kind=3)
       Type
          LBracket: "[" @4:9 (kind=43)
          RBracket: "]" @4:10 (kind=44)
          Ident: "int" @4:11 (kind=12)
       Eq: "=" @4:14 (kind=49)
       Init
          LBracket: "[" @4:15 (kind=43)
          RBracket: "]" @4:16 (kind=44)
          Ident: "int" @4:17 (kind=12)
        LBrace: "{" @4:20 (kind=41)
         Elements
          IntLitExpr
           Value: "1" @4:21 (kind=4)
          IntLitExpr
           Value: "2" @4:23 (kind=4)
          IntLitExpr
           Value: "3" @4:25 (kind=4)
        RBrace: "}" @4:26 (kind=42)
      VarDeclStmt
       Var: "var" @5:3 (kind=11)
       Name: "y" @5:7 (kind=3)
       View: "view" @5:9 (kind=76)
       Type
          LBracket: "[" @5:14 (kind=43)
          RBracket: "]" @5:15 (kind=44)
          Ident: "int" @5:16 (kind=12)
       Eq: "=" @5:19 (kind=49)
       Init
         IdentExpr
          Name: "x" @5:20 (kind=3)
          LBracket: "[" @5:21 (kind=43)
         IntLitExpr
          Value: "3" @5:22 (kind=4)
          Colon: ":" @5:23 (kind=47)
         IntLitExpr
          Value: "6" @5:24 (kind=4)
          RBracket: "]" @5:25 (kind=44)
      VarDeclStmt
       Var: "var" @6:3 (kind=11)
       Name: "z" @6:7 (kind=3)
       Type
          LBracket: "[" @6:9 (kind=43)
          RBracket: "]" @6:10 (kind=44)
          Ident: "int" @6:11 (kind=12)
       Eq: "=" @6:14 (kind=49)
       Init
          LBracket: "[" @6:15 (kind=43)
          RBracket: "]" @6:16 (kind=44)
          Ident: "int" @6:17 (kind=12)
        LBrace: "{" @6:20 (kind=41)
         Elements
          IntLitExpr
           Value: "1" @6:21 (kind=4)
          IntLitExpr
           Value: "2" @6:23 (kind=4)
          IntLitExpr
           Value: "3" @6:25 (kind=4)
        RBrace: "}" @6:26 (kind=42)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("var_slice_x6", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func main(){
  var y view []int=x[:6]
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
      VarDeclStmt
       Var: "var" @4:3 (kind=11)
       Name: "y" @4:7 (kind=3)
       View: "view" @4:9 (kind=76)
       Type
          LBracket: "[" @4:14 (kind=43)
          RBracket: "]" @4:15 (kind=44)
          Ident: "int" @4:16 (kind=12)
       Eq: "=" @4:19 (kind=49)
       Init
         IdentExpr
          Name: "x" @4:20 (kind=3)
          LBracket: "[" @4:21 (kind=43)
          Colon: ":" @4:22 (kind=47)
         IntLitExpr
          Value: "6" @4:23 (kind=4)
          RBracket: "]" @4:24 (kind=44)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("var_slice_x7", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func main(){
  var y view []int=x[6]
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
      VarDeclStmt
       Var: "var" @4:3 (kind=11)
       Name: "y" @4:7 (kind=3)
       View: "view" @4:9 (kind=76)
       Type
          LBracket: "[" @4:14 (kind=43)
          RBracket: "]" @4:15 (kind=44)
          Ident: "int" @4:16 (kind=12)
       Eq: "=" @4:19 (kind=49)
       Init
         IdentExpr
          Name: "x" @4:20 (kind=3)
          LBracket: "[" @4:21 (kind=43)
         IntLitExpr
          Value: "6" @4:22 (kind=4)
          RBracket: "]" @4:23 (kind=44)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("const_array_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

const x [3]int=[]int{1,2,3}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 ConstDecls
  ConstDecl
   Const: "const" @3:1 (kind=23)
   Name: "x" @3:7 (kind=3)
   Type
      LBracket: "[" @3:9 (kind=43)
      Size: "3" @3:10 (kind=4)
      RBracket: "]" @3:11 (kind=44)
      Ident: "int" @3:12 (kind=12)
   Eq: "=" @3:15 (kind=49)
   Init
      LBracket: "[" @3:16 (kind=43)
      RBracket: "]" @3:17 (kind=44)
      Ident: "int" @3:18 (kind=12)
    LBrace: "{" @3:21 (kind=41)
     Elements
      IntLitExpr
       Value: "1" @3:22 (kind=4)
      IntLitExpr
       Value: "2" @3:24 (kind=4)
      IntLitExpr
       Value: "3" @3:26 (kind=4)
    RBrace: "}" @3:27 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("const_array_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func main(){
  const x [5]int=[]int{1,2,3};const y []int=[]int{1,2,3}
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
      ConstDecl
       Const: "const" @4:3 (kind=23)
       Name: "x" @4:9 (kind=3)
       Type
          LBracket: "[" @4:11 (kind=43)
          Size: "5" @4:12 (kind=4)
          RBracket: "]" @4:13 (kind=44)
          Ident: "int" @4:14 (kind=12)
       Eq: "=" @4:17 (kind=49)
       Init
          LBracket: "[" @4:18 (kind=43)
          RBracket: "]" @4:19 (kind=44)
          Ident: "int" @4:20 (kind=12)
        LBrace: "{" @4:23 (kind=41)
         Elements
          IntLitExpr
           Value: "1" @4:24 (kind=4)
          IntLitExpr
           Value: "2" @4:26 (kind=4)
          IntLitExpr
           Value: "3" @4:28 (kind=4)
        RBrace: "}" @4:29 (kind=42)
      ConstDecl
       Const: "const" @4:31 (kind=23)
       Name: "y" @4:37 (kind=3)
       Type
          LBracket: "[" @4:39 (kind=43)
          RBracket: "]" @4:40 (kind=44)
          Ident: "int" @4:41 (kind=12)
       Eq: "=" @4:44 (kind=49)
       Init
          LBracket: "[" @4:45 (kind=43)
          RBracket: "]" @4:46 (kind=44)
          Ident: "int" @4:47 (kind=12)
        LBrace: "{" @4:50 (kind=41)
         Elements
          IntLitExpr
           Value: "1" @4:51 (kind=4)
          IntLitExpr
           Value: "2" @4:53 (kind=4)
          IntLitExpr
           Value: "3" @4:55 (kind=4)
        RBrace: "}" @4:56 (kind=42)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("var_array_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func main(){
  var x [5]int=[]int{1,2,3}
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
      VarDeclStmt
       Var: "var" @4:3 (kind=11)
       Name: "x" @4:7 (kind=3)
       Type
          LBracket: "[" @4:9 (kind=43)
          Size: "5" @4:10 (kind=4)
          RBracket: "]" @4:11 (kind=44)
          Ident: "int" @4:12 (kind=12)
       Eq: "=" @4:15 (kind=49)
       Init
          LBracket: "[" @4:16 (kind=43)
          RBracket: "]" @4:17 (kind=44)
          Ident: "int" @4:18 (kind=12)
        LBrace: "{" @4:21 (kind=41)
         Elements
          IntLitExpr
           Value: "1" @4:22 (kind=4)
          IntLitExpr
           Value: "2" @4:24 (kind=4)
          IntLitExpr
           Value: "3" @4:26 (kind=4)
        RBrace: "}" @4:27 (kind=42)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("var_array_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func main(){
  var x [5]int=[]int{1,2,3}
	var y view []int=x[:]
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
      VarDeclStmt
       Var: "var" @4:3 (kind=11)
       Name: "x" @4:7 (kind=3)
       Type
          LBracket: "[" @4:9 (kind=43)
          Size: "5" @4:10 (kind=4)
          RBracket: "]" @4:11 (kind=44)
          Ident: "int" @4:12 (kind=12)
       Eq: "=" @4:15 (kind=49)
       Init
          LBracket: "[" @4:16 (kind=43)
          RBracket: "]" @4:17 (kind=44)
          Ident: "int" @4:18 (kind=12)
        LBrace: "{" @4:21 (kind=41)
         Elements
          IntLitExpr
           Value: "1" @4:22 (kind=4)
          IntLitExpr
           Value: "2" @4:24 (kind=4)
          IntLitExpr
           Value: "3" @4:26 (kind=4)
        RBrace: "}" @4:27 (kind=42)
      VarDeclStmt
       Var: "var" @5:2 (kind=11)
       Name: "y" @5:6 (kind=3)
       View: "view" @5:8 (kind=76)
       Type
          LBracket: "[" @5:13 (kind=43)
          RBracket: "]" @5:14 (kind=44)
          Ident: "int" @5:15 (kind=12)
       Eq: "=" @5:18 (kind=49)
       Init
         IdentExpr
          Name: "x" @5:19 (kind=3)
          LBracket: "[" @5:20 (kind=43)
          Colon: ":" @5:21 (kind=47)
          RBracket: "]" @5:22 (kind=44)
     RBrace: "}" @6:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("var_array_x3", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func main(){
  var x [5]int=[]int{1,2,3}
	var y view []int=x[3:6];
  var z [5]int=[]int{1,2,3}
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
      VarDeclStmt
       Var: "var" @4:3 (kind=11)
       Name: "x" @4:7 (kind=3)
       Type
          LBracket: "[" @4:9 (kind=43)
          Size: "5" @4:10 (kind=4)
          RBracket: "]" @4:11 (kind=44)
          Ident: "int" @4:12 (kind=12)
       Eq: "=" @4:15 (kind=49)
       Init
          LBracket: "[" @4:16 (kind=43)
          RBracket: "]" @4:17 (kind=44)
          Ident: "int" @4:18 (kind=12)
        LBrace: "{" @4:21 (kind=41)
         Elements
          IntLitExpr
           Value: "1" @4:22 (kind=4)
          IntLitExpr
           Value: "2" @4:24 (kind=4)
          IntLitExpr
           Value: "3" @4:26 (kind=4)
        RBrace: "}" @4:27 (kind=42)
      VarDeclStmt
       Var: "var" @5:2 (kind=11)
       Name: "y" @5:6 (kind=3)
       View: "view" @5:8 (kind=76)
       Type
          LBracket: "[" @5:13 (kind=43)
          RBracket: "]" @5:14 (kind=44)
          Ident: "int" @5:15 (kind=12)
       Eq: "=" @5:18 (kind=49)
       Init
         IdentExpr
          Name: "x" @5:19 (kind=3)
          LBracket: "[" @5:20 (kind=43)
         IntLitExpr
          Value: "3" @5:21 (kind=4)
          Colon: ":" @5:22 (kind=47)
         IntLitExpr
          Value: "6" @5:23 (kind=4)
          RBracket: "]" @5:24 (kind=44)
      VarDeclStmt
       Var: "var" @6:3 (kind=11)
       Name: "z" @6:7 (kind=3)
       Type
          LBracket: "[" @6:9 (kind=43)
          Size: "5" @6:10 (kind=4)
          RBracket: "]" @6:11 (kind=44)
          Ident: "int" @6:12 (kind=12)
       Eq: "=" @6:15 (kind=49)
       Init
          LBracket: "[" @6:16 (kind=43)
          RBracket: "]" @6:17 (kind=44)
          Ident: "int" @6:18 (kind=12)
        LBrace: "{" @6:21 (kind=41)
         Elements
          IntLitExpr
           Value: "1" @6:22 (kind=4)
          IntLitExpr
           Value: "2" @6:24 (kind=4)
          IntLitExpr
           Value: "3" @6:26 (kind=4)
        RBrace: "}" @6:27 (kind=42)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("var_array_x4", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func main(){
	var y view [5]int=x[:6]
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
      VarDeclStmt
       Var: "var" @4:2 (kind=11)
       Name: "y" @4:6 (kind=3)
       View: "view" @4:8 (kind=76)
       Type
          LBracket: "[" @4:13 (kind=43)
          Size: "5" @4:14 (kind=4)
          RBracket: "]" @4:15 (kind=44)
          Ident: "int" @4:16 (kind=12)
       Eq: "=" @4:19 (kind=49)
       Init
         IdentExpr
          Name: "x" @4:20 (kind=3)
          LBracket: "[" @4:21 (kind=43)
          Colon: ":" @4:22 (kind=47)
         IntLitExpr
          Value: "6" @4:23 (kind=4)
          RBracket: "]" @4:24 (kind=44)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("var_array_x5", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func main(){
	var y view []int=x[6]
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
      VarDeclStmt
       Var: "var" @4:2 (kind=11)
       Name: "y" @4:6 (kind=3)
       View: "view" @4:8 (kind=76)
       Type
          LBracket: "[" @4:13 (kind=43)
          RBracket: "]" @4:14 (kind=44)
          Ident: "int" @4:15 (kind=12)
       Eq: "=" @4:18 (kind=49)
       Init
         IdentExpr
          Name: "x" @4:19 (kind=3)
          LBracket: "[" @4:20 (kind=43)
         IntLitExpr
          Value: "6" @4:21 (kind=4)
          RBracket: "]" @4:22 (kind=44)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("bad_slice_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

const x []int
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_slice_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

var x []int
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_slice_x3", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

const x []struct=[]int{1,2,3}
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_array_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

const x [5]int
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}
