package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParser_func_decl(t *testing.T) {
	assert := assert.New(t)

	t.Run("return_types_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x()int{
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
   Results
     Param
      Type
       NamedType
        Ident: "int" @3:9 (kind=12)
   Body
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("return_types_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x()(int,int){
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
   Results
    LParent: "(" @3:9 (kind=39)
     Param
      Type
       NamedType
        Ident: "int" @3:10 (kind=12)
     Param
      Type
       NamedType
        Ident: "int" @3:14 (kind=12)
    RParent: ")" @3:17 (kind=40)
   Body
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("return_types_x3", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x()(a int,b int){}
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
   Results
    LParent: "(" @3:9 (kind=39)
     Param
      Ident: "a" @3:10 (kind=3)
      Type
       NamedType
        Ident: "int" @3:12 (kind=12)
     Param
      Ident: "b" @3:16 (kind=3)
      Type
       NamedType
        Ident: "int" @3:18 (kind=12)
    RParent: ")" @3:21 (kind=40)
   Body
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("return_types_x4", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){}
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
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("return_types_x5", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x()(a z, b z){}
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
   Results
    LParent: "(" @3:9 (kind=39)
     Param
      Ident: "a" @3:10 (kind=3)
      Type
       NamedType
        Ident: "z" @3:12 (kind=3)
     Param
      Ident: "b" @3:15 (kind=3)
      Type
       NamedType
        Ident: "z" @3:17 (kind=3)
    RParent: ")" @3:18 (kind=40)
   Body
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("return_types_x6", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x()(a z){}
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
   Results
    LParent: "(" @3:9 (kind=39)
     Param
      Ident: "a" @3:10 (kind=3)
      Type
       NamedType
        Ident: "z" @3:12 (kind=3)
    RParent: ")" @3:13 (kind=40)
   Body
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("return_types_x7", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x()(a []int){}
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
   Results
    LParent: "(" @3:9 (kind=39)
     Param
      Ident: "a" @3:10 (kind=3)
      Type
       SliceType:
        LBracket: "[" @3:12 (kind=43)
        RBracket: "]" @3:13 (kind=44)
        NamedType
         Ident: "int" @3:14 (kind=12)
    RParent: ")" @3:17 (kind=40)
   Body
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("return_types_x8", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x()[]int{}
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
   Results
     Param
      Type
       SliceType:
        LBracket: "[" @3:9 (kind=43)
        RBracket: "]" @3:10 (kind=44)
        NamedType
         Ident: "int" @3:11 (kind=12)
   Body
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("return_types_x9", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x()([]int,[]string){}
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
   Results
    LParent: "(" @3:9 (kind=39)
     Param
      Type
       SliceType:
        LBracket: "[" @3:10 (kind=43)
        RBracket: "]" @3:11 (kind=44)
        NamedType
         Ident: "int" @3:12 (kind=12)
     Param
      Type
       SliceType:
        LBracket: "[" @3:16 (kind=43)
        RBracket: "]" @3:17 (kind=44)
        NamedType
         Ident: "string" @3:18 (kind=24)
    RParent: ")" @3:24 (kind=40)
   Body
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("return_types_x10", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x1(m map[string]string)([]int,map[string]string){}
func x2(m map[string]string)(a []int, b map[string]string){}
func x3(m map[string]string) map[string]string {}
func x4(m map[string]string) map[string]string {}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  FuncDecl
   Function: "func" @3:1 (kind=10)
   Name: "x1" @3:6 (kind=3)
   Params
    Param
     Ident: "m" @3:9 (kind=3)
     Type
      MapType:
       Map: "map" @3:11 (kind=79)
       LBracket: "[" @3:14 (kind=43)
       KeyType:
        NamedType
         Ident: "string" @3:15 (kind=24)
       RBracket: "]" @3:21 (kind=44)
       ValueType:
        NamedType
         Ident: "string" @3:22 (kind=24)
   Results
    LParent: "(" @3:29 (kind=39)
     Param
      Type
       SliceType:
        LBracket: "[" @3:30 (kind=43)
        RBracket: "]" @3:31 (kind=44)
        NamedType
         Ident: "int" @3:32 (kind=12)
     Param
      Type
       MapType:
        Map: "map" @3:36 (kind=79)
        LBracket: "[" @3:39 (kind=43)
        KeyType:
         NamedType
          Ident: "string" @3:40 (kind=24)
        RBracket: "]" @3:46 (kind=44)
        ValueType:
         NamedType
          Ident: "string" @3:47 (kind=24)
    RParent: ")" @3:53 (kind=40)
   Body
  FuncDecl
   Function: "func" @4:1 (kind=10)
   Name: "x2" @4:6 (kind=3)
   Params
    Param
     Ident: "m" @4:9 (kind=3)
     Type
      MapType:
       Map: "map" @4:11 (kind=79)
       LBracket: "[" @4:14 (kind=43)
       KeyType:
        NamedType
         Ident: "string" @4:15 (kind=24)
       RBracket: "]" @4:21 (kind=44)
       ValueType:
        NamedType
         Ident: "string" @4:22 (kind=24)
   Results
    LParent: "(" @4:29 (kind=39)
     Param
      Ident: "a" @4:30 (kind=3)
      Type
       SliceType:
        LBracket: "[" @4:32 (kind=43)
        RBracket: "]" @4:33 (kind=44)
        NamedType
         Ident: "int" @4:34 (kind=12)
     Param
      Ident: "b" @4:39 (kind=3)
      Type
       MapType:
        Map: "map" @4:41 (kind=79)
        LBracket: "[" @4:44 (kind=43)
        KeyType:
         NamedType
          Ident: "string" @4:45 (kind=24)
        RBracket: "]" @4:51 (kind=44)
        ValueType:
         NamedType
          Ident: "string" @4:52 (kind=24)
    RParent: ")" @4:58 (kind=40)
   Body
  FuncDecl
   Function: "func" @5:1 (kind=10)
   Name: "x3" @5:6 (kind=3)
   Params
    Param
     Ident: "m" @5:9 (kind=3)
     Type
      MapType:
       Map: "map" @5:11 (kind=79)
       LBracket: "[" @5:14 (kind=43)
       KeyType:
        NamedType
         Ident: "string" @5:15 (kind=24)
       RBracket: "]" @5:21 (kind=44)
       ValueType:
        NamedType
         Ident: "string" @5:22 (kind=24)
   Results
     Param
      Type
       MapType:
        Map: "map" @5:30 (kind=79)
        LBracket: "[" @5:33 (kind=43)
        KeyType:
         NamedType
          Ident: "string" @5:34 (kind=24)
        RBracket: "]" @5:40 (kind=44)
        ValueType:
         NamedType
          Ident: "string" @5:41 (kind=24)
   Body
  FuncDecl
   Function: "func" @6:1 (kind=10)
   Name: "x4" @6:6 (kind=3)
   Params
    Param
     Ident: "m" @6:9 (kind=3)
     Type
      MapType:
       Map: "map" @6:11 (kind=79)
       LBracket: "[" @6:14 (kind=43)
       KeyType:
        NamedType
         Ident: "string" @6:15 (kind=24)
       RBracket: "]" @6:21 (kind=44)
       ValueType:
        NamedType
         Ident: "string" @6:22 (kind=24)
   Results
     Param
      Type
       MapType:
        Map: "map" @6:30 (kind=79)
        LBracket: "[" @6:33 (kind=43)
        KeyType:
         NamedType
          Ident: "string" @6:34 (kind=24)
        RBracket: "]" @6:40 (kind=44)
        ValueType:
         NamedType
          Ident: "string" @6:41 (kind=24)
   Body
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("struct_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  type test struct {
	  x int
		m map[string]string
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
      StructDecl:
       Type: "type" @4:3 (kind=26)
       Name: "test" @4:8 (kind=3)
       Struct: "struct" @4:13 (kind=27)
       LBrace: "{" @4:20 (kind=41)
        Name: "x" @5:4 (kind=3)
        Type:
         NamedType
          Ident: "int" @5:6 (kind=12)
        Name: "m" @6:3 (kind=3)
        Type:
         MapType:
          Map: "map" @6:5 (kind=79)
          LBracket: "[" @6:8 (kind=43)
          KeyType:
           NamedType
            Ident: "string" @6:9 (kind=24)
          RBracket: "]" @6:15 (kind=44)
          ValueType:
           NamedType
            Ident: "string" @6:16 (kind=24)
       RBrace: "}" @7:2 (kind=42)
     RBrace: "}" @8:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("interface_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  type test interface {}
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
      InterfaceDecl:
       Type: "type" @4:3 (kind=26)
       Name: "test" @4:8 (kind=3)
       Interface: "interface" @4:13 (kind=28)
       LBrace: "{" @4:23 (kind=41)
       RBrace: "}" @4:24 (kind=42)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("enum_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  type Color enum {
	  Red
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
      EnumDecl:
       Type: "type" @4:3 (kind=26)
        Name: "Color" @4:8 (kind=3)
        Public: true
        Enum: "enum" @4:14 (kind=74)
        LBrace: "{" @4:19 (kind=41)
         Variants
          Ident: "Red" @5:4 (kind=3)
        RBrace: "}" @6:2 (kind=42)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("sum_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(){
  type Shape sum {
	  Circle(radius float)
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
      SumDecl:
       Type: "type" @4:3 (kind=26)
        Name: "Shape" @4:8 (kind=3)
        Sum: "sum" @4:14 (kind=75)
       Public: true
       LBrace: "{" @4:18 (kind=41)
        Variants
         SumVariant: "Circle" @5:4 (kind=3)
          Params
           Param
            Ident: "radius" @5:11 (kind=3)
            Type
             NamedType
              Ident: "float" @5:18 (kind=20)
       RBrace: "}" @6:2 (kind=42)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("bad_return_types_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x(),int{}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_return_types_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x()(int,,int){}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_return_types_x3", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x()(a int,b){}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_return_types_x4", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x()(){}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_return_types_x5", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x()(,a int,b){}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_return_types_x6", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x()int,{}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_return_types_x7", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x()(a int,b,){}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_return_types_x8", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x()(a int,b b b){}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_return_types_x9", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x()(a int,b b,){}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_return_types_x10", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x()(a int,b return){}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_return_types_x11", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x()(a int,b return){}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_return_types_x12", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x()(a return,b return){}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_return_types_x13", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x() struct {}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_return_types_x14", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x()(int, struct){}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_return_types_x15", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x()(int struct){}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_return_types_x16", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x()(int,int,){}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_return_types_x17", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x()(a, b z){}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_return_types_x18", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x()(_ b){}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("bad_return_types_x19", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func x()(a _){}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})
}
