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
	a []int
	b map[string]string
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
    Name: "a" @6:2 (kind=3)
    Type:
     SliceType:
      LBracket: "[" @6:4 (kind=43)
      RBracket: "]" @6:5 (kind=44)
      NamedType
       Ident: "int" @6:6 (kind=12)
    Name: "b" @7:2 (kind=3)
    Type:
     MapType:
      Map: "map" @7:4 (kind=79)
      LBracket: "[" @7:7 (kind=43)
      KeyType:
       NamedType
        Ident: "string" @7:8 (kind=24)
      RBracket: "]" @7:14 (kind=44)
      ValueType:
       NamedType
        Ident: "string" @7:15 (kind=24)
   RBrace: "}" @8:1 (kind=42)
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
     LParen: "(" @4:15 (kind=39)
     Args:
      IdentExpr
       Name: "a" @4:16 (kind=3)
      IdentExpr
       Name: "b" @4:19 (kind=3)
     RParen: ")" @4:20 (kind=40)
   RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("x9", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test struct{x int}

func main(){
  x := test{x: int(0)}
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
   RBrace: "}" @3:23 (kind=42)
  FuncDecl
   Function: "func" @5:1 (kind=10)
   Name: "main" @5:6 (kind=3)
   Params
    (none)
   Body
    BlockStmt
     LBrace: "{" @5:12 (kind=41)
     Stmts
      AssignStmt
       Left
        IdentExpr
         Name: "x" @6:3 (kind=3)
       Operator: ":=" @6:5 (kind=50)
       Right
         Name: "test" @6:8 (kind=3)
          LBrace: "{" @6:12 (kind=41)
           Fields:
            Field:
             Name: "x" @6:13 (kind=3)
             Value:
              CallExpr
               Callee
                IdentExpr
                 Name: "int" @6:16 (kind=3)
               LParen: "(" @6:19 (kind=39)
               Args:
                IntLitExpr
                 Value: "0" @6:20 (kind=4)
               RParen: ")" @6:21 (kind=40)
          RBrace: "}" @6:22 (kind=42)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("x10", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test struct{x int}

func main(){
  x := test{x: int(0), b: "toto"}
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
   RBrace: "}" @3:23 (kind=42)
  FuncDecl
   Function: "func" @5:1 (kind=10)
   Name: "main" @5:6 (kind=3)
   Params
    (none)
   Body
    BlockStmt
     LBrace: "{" @5:12 (kind=41)
     Stmts
      AssignStmt
       Left
        IdentExpr
         Name: "x" @6:3 (kind=3)
       Operator: ":=" @6:5 (kind=50)
       Right
         Name: "test" @6:8 (kind=3)
          LBrace: "{" @6:12 (kind=41)
           Fields:
            Field:
             Name: "x" @6:13 (kind=3)
             Value:
              CallExpr
               Callee
                IdentExpr
                 Name: "int" @6:16 (kind=3)
               LParen: "(" @6:19 (kind=39)
               Args:
                IntLitExpr
                 Value: "0" @6:20 (kind=4)
               RParen: ")" @6:21 (kind=40)
            Field:
             Name: "b" @6:24 (kind=3)
             Value:
              StringLitExpr
               Value: "toto" @6:27 (kind=6)
          RBrace: "}" @6:33 (kind=42)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("x11", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test struct{x int}

func main(){
  x := test{}
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
   RBrace: "}" @3:23 (kind=42)
  FuncDecl
   Function: "func" @5:1 (kind=10)
   Name: "main" @5:6 (kind=3)
   Params
    (none)
   Body
    BlockStmt
     LBrace: "{" @5:12 (kind=41)
     Stmts
      AssignStmt
       Left
        IdentExpr
         Name: "x" @6:3 (kind=3)
       Operator: ":=" @6:5 (kind=50)
       Right
         Name: "test" @6:8 (kind=3)
          LBrace: "{" @6:12 (kind=41)
          RBrace: "}" @6:13 (kind=42)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("x12", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test struct{x int}

func main(){
  var x test = test{x: int(0), b: "toto"}
  const x test = test{x: int64(0), b: "tata"}
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
   RBrace: "}" @3:23 (kind=42)
  FuncDecl
   Function: "func" @5:1 (kind=10)
   Name: "main" @5:6 (kind=3)
   Params
    (none)
   Body
    BlockStmt
     LBrace: "{" @5:12 (kind=41)
     Stmts
      VarDecl
       Var: "var" @6:3 (kind=11)
       Name: "x" @6:7 (kind=3)
       Type
        NamedType
         Ident: "test" @6:9 (kind=3)
       Eq: "=" @6:14 (kind=49)
       Init
         Name: "test" @6:16 (kind=3)
          LBrace: "{" @6:20 (kind=41)
           Fields:
            Field:
             Name: "x" @6:21 (kind=3)
             Value:
              CallExpr
               Callee
                IdentExpr
                 Name: "int" @6:24 (kind=3)
               LParen: "(" @6:27 (kind=39)
               Args:
                IntLitExpr
                 Value: "0" @6:28 (kind=4)
               RParen: ")" @6:29 (kind=40)
            Field:
             Name: "b" @6:32 (kind=3)
             Value:
              StringLitExpr
               Value: "toto" @6:35 (kind=6)
          RBrace: "}" @6:41 (kind=42)
      ConstDecl
       Const: "const" @7:3 (kind=23)
       Name: "x" @7:9 (kind=3)
       Type
        NamedType
         Ident: "test" @7:11 (kind=3)
       Eq: "=" @7:16 (kind=49)
       Init
         Name: "test" @7:18 (kind=3)
          LBrace: "{" @7:22 (kind=41)
           Fields:
            Field:
             Name: "x" @7:23 (kind=3)
             Value:
              CallExpr
               Callee
                IdentExpr
                 Name: "int64" @7:26 (kind=3)
               LParen: "(" @7:31 (kind=39)
               Args:
                IntLitExpr
                 Value: "0" @7:32 (kind=4)
               RParen: ")" @7:33 (kind=40)
            Field:
             Name: "b" @7:36 (kind=3)
             Value:
              StringLitExpr
               Value: "tata" @7:39 (kind=6)
          RBrace: "}" @7:45 (kind=42)
     RBrace: "}" @8:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("x13", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type User struct {
  name string
}

func f() {
  User{name: "a"}.name = "b"
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
   Name: "User" @3:6 (kind=3)
   Struct: "struct" @3:11 (kind=27)
   Public: true
   LBrace: "{" @3:18 (kind=41)
    Name: "name" @4:3 (kind=3)
    Type:
     NamedType
      Ident: "string" @4:8 (kind=24)
   RBrace: "}" @5:1 (kind=42)
  FuncDecl
   Function: "func" @7:1 (kind=10)
   Name: "f" @7:6 (kind=3)
   Params
    (none)
   Body
    BlockStmt
     LBrace: "{" @7:10 (kind=41)
     Stmts
      AssignStmt
       Left
        SelectorExpr
         X:
           Name: "User" @8:3 (kind=3)
            LBrace: "{" @8:7 (kind=41)
             Fields:
              Field:
               Name: "name" @8:8 (kind=3)
               Value:
                StringLitExpr
                 Value: "a" @8:14 (kind=6)
            RBrace: "}" @8:17 (kind=42)
         Dot: "." @8:18 (kind=48)
         Selector: "name" @8:19 (kind=3)
       Operator: "=" @8:24 (kind=49)
       Right
        StringLitExpr
         Value: "b" @8:26 (kind=6)
     RBrace: "}" @9:1 (kind=42)
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

	t.Run("bad_x6", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

type test struct{x int}

func main(){
  x := test{x:}
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

type test struct{x int}

func main(){
  x := test{x: 5, y: 10 1}
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})
}
