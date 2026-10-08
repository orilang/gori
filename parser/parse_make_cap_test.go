package parser

import (
	"fmt"
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParser_parser_make_cap(t *testing.T) {
	assert := assert.New(t)

	t.Run("map_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func main() {
  var x map[string]string = make(map[string]string) // comment
  var y hashmap[string]string = make(hashmap[string]string)
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
     LBrace: "{" @3:13 (kind=41)
     Stmts
      VarDecl
       Var: "var" @4:3 (kind=11)
       Name: "x" @4:7 (kind=3)
       Type
        MapType:
         Map: "map" @4:9 (kind=79)
         LBracket: "[" @4:12 (kind=43)
         KeyType:
          NamedType
           Ident: "string" @4:13 (kind=24)
         RBracket: "]" @4:19 (kind=44)
         ValueType:
          NamedType
           Ident: "string" @4:20 (kind=24)
       Eq: "=" @4:27 (kind=49)
       Init
        MakeExpr:
         Make: "make" @4:29 (kind=3)
         LParen: "(" @4:33 (kind=39)
         MapType:
          Map: "map" @4:34 (kind=79)
          LBracket: "[" @4:37 (kind=43)
          KeyType:
           NamedType
            Ident: "string" @4:38 (kind=24)
          RBracket: "]" @4:44 (kind=44)
          ValueType:
           NamedType
            Ident: "string" @4:45 (kind=24)
         RParen: ")" @4:51 (kind=40)
      VarDecl
       Var: "var" @5:3 (kind=11)
       Name: "y" @5:7 (kind=3)
       Type
        MapType:
         Hashmap: "hashmap" @5:9 (kind=80)
         LBracket: "[" @5:16 (kind=43)
         KeyType:
          NamedType
           Ident: "string" @5:17 (kind=24)
         RBracket: "]" @5:23 (kind=44)
         ValueType:
          NamedType
           Ident: "string" @5:24 (kind=24)
       Eq: "=" @5:31 (kind=49)
       Init
        MakeExpr:
         Make: "make" @5:33 (kind=3)
         LParen: "(" @5:37 (kind=39)
         MapType:
          Hashmap: "hashmap" @5:38 (kind=80)
          LBracket: "[" @5:45 (kind=43)
          KeyType:
           NamedType
            Ident: "string" @5:46 (kind=24)
          RBracket: "]" @5:52 (kind=44)
          ValueType:
           NamedType
            Ident: "string" @5:53 (kind=24)
         RParen: ")" @5:59 (kind=40)
     RBrace: "}" @6:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("map_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func main(){
  var x map[string]string = make(map[string]string,10)
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
      VarDecl
       Var: "var" @4:3 (kind=11)
       Name: "x" @4:7 (kind=3)
       Type
        MapType:
         Map: "map" @4:9 (kind=79)
         LBracket: "[" @4:12 (kind=43)
         KeyType:
          NamedType
           Ident: "string" @4:13 (kind=24)
         RBracket: "]" @4:19 (kind=44)
         ValueType:
          NamedType
           Ident: "string" @4:20 (kind=24)
       Eq: "=" @4:27 (kind=49)
       Init
        MakeExpr:
         Make: "make" @4:29 (kind=3)
         LParen: "(" @4:33 (kind=39)
         MapType:
          Map: "map" @4:34 (kind=79)
          LBracket: "[" @4:37 (kind=43)
          KeyType:
           NamedType
            Ident: "string" @4:38 (kind=24)
          RBracket: "]" @4:44 (kind=44)
          ValueType:
           NamedType
            Ident: "string" @4:45 (kind=24)
         Size:
          IntLitExpr
           Value: "10" @4:52 (kind=4)
         RParen: ")" @4:54 (kind=40)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("map_x3", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main
    
const x map[string]string = make(map[string]string) // comment
func main() {
  var y hashmap[string]string = make(hashmap[string]string)
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  ConstDecl
   Const: "const" @3:1 (kind=23)
   Name: "x" @3:7 (kind=3)
   Type
    MapType:
     Map: "map" @3:9 (kind=79)
     LBracket: "[" @3:12 (kind=43)
     KeyType:
      NamedType
       Ident: "string" @3:13 (kind=24)
     RBracket: "]" @3:19 (kind=44)
     ValueType:
      NamedType
       Ident: "string" @3:20 (kind=24)
   Eq: "=" @3:27 (kind=49)
   Init
    MakeExpr:
     Make: "make" @3:29 (kind=3)
     LParen: "(" @3:33 (kind=39)
     MapType:
      Map: "map" @3:34 (kind=79)
      LBracket: "[" @3:37 (kind=43)
      KeyType:
       NamedType
        Ident: "string" @3:38 (kind=24)
      RBracket: "]" @3:44 (kind=44)
      ValueType:
       NamedType
        Ident: "string" @3:45 (kind=24)
     RParen: ")" @3:51 (kind=40)
  FuncDecl
   Function: "func" @4:1 (kind=10)
   Name: "main" @4:6 (kind=3)
   Params
    (none)
   Body
    BlockStmt
     LBrace: "{" @4:13 (kind=41)
     Stmts
      VarDecl
       Var: "var" @5:3 (kind=11)
       Name: "y" @5:7 (kind=3)
       Type
        MapType:
         Hashmap: "hashmap" @5:9 (kind=80)
         LBracket: "[" @5:16 (kind=43)
         KeyType:
          NamedType
           Ident: "string" @5:17 (kind=24)
         RBracket: "]" @5:23 (kind=44)
         ValueType:
          NamedType
           Ident: "string" @5:24 (kind=24)
       Eq: "=" @5:31 (kind=49)
       Init
        MakeExpr:
         Make: "make" @5:33 (kind=3)
         LParen: "(" @5:37 (kind=39)
         MapType:
          Hashmap: "hashmap" @5:38 (kind=80)
          LBracket: "[" @5:45 (kind=43)
          KeyType:
           NamedType
            Ident: "string" @5:46 (kind=24)
          RBracket: "]" @5:52 (kind=44)
          ValueType:
           NamedType
            Ident: "string" @5:53 (kind=24)
         RParen: ")" @5:59 (kind=40)
     RBrace: "}" @6:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("map_x4", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func main() {
   x := make(hashmap[string]string)
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
     LBrace: "{" @3:13 (kind=41)
     Stmts
      AssignStmt
       Left
        IdentExpr
         Name: "x" @4:4 (kind=3)
       Operator: ":=" @4:6 (kind=50)
       Right
        MakeExpr:
         Make: "make" @4:9 (kind=3)
         LParen: "(" @4:13 (kind=39)
         MapType:
          Hashmap: "hashmap" @4:14 (kind=80)
          LBracket: "[" @4:21 (kind=43)
          KeyType:
           NamedType
            Ident: "string" @4:22 (kind=24)
          RBracket: "]" @4:28 (kind=44)
          ValueType:
           NamedType
            Ident: "string" @4:29 (kind=24)
         RParen: ")" @4:35 (kind=40)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("map_x5", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main
type UserID []string
func f() UserID {
  m := make(UserID, int(5))
  return m
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  DefinedTypeDecl:
   TypeDecl: "type" @2:1 (kind=26)
    Name: "UserID" @2:6 (kind=3)
    Type
     SliceType:
      LBracket: "[" @2:13 (kind=43)
      RBracket: "]" @2:14 (kind=44)
      NamedType
       Ident: "string" @2:15 (kind=24)
  FuncDecl
   Function: "func" @3:1 (kind=10)
   Name: "f" @3:6 (kind=3)
   Params
    (none)
   Results
     Param
      Type
       NamedType
        Ident: "UserID" @3:10 (kind=3)
   Body
    BlockStmt
     LBrace: "{" @3:17 (kind=41)
     Stmts
      AssignStmt
       Left
        IdentExpr
         Name: "m" @4:3 (kind=3)
       Operator: ":=" @4:5 (kind=50)
       Right
        MakeExpr:
         Make: "make" @4:8 (kind=3)
         LParen: "(" @4:12 (kind=39)
         NamedType
          Ident: "UserID" @4:13 (kind=3)
         Size:
          CallExpr
           Callee
            IdentExpr
             Name: "int" @4:21 (kind=3)
           LParen: "(" @4:24 (kind=39)
           Args:
            IntLitExpr
             Value: "5" @4:25 (kind=4)
           RParen: ")" @4:26 (kind=40)
         RParen: ")" @4:27 (kind=40)
      ReturnStmt
       Values
        IdentExpr
         Name: "m" @5:10 (kind=3)
     RBrace: "}" @6:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("slice_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func main() {
  var x []string = make([]string,10);
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
     LBrace: "{" @3:13 (kind=41)
     Stmts
      VarDecl
       Var: "var" @4:3 (kind=11)
       Name: "x" @4:7 (kind=3)
       Type
        SliceType:
         LBracket: "[" @4:9 (kind=43)
         RBracket: "]" @4:10 (kind=44)
         NamedType
          Ident: "string" @4:11 (kind=24)
       Eq: "=" @4:18 (kind=49)
       Init
        MakeExpr:
         Make: "make" @4:20 (kind=3)
         LParen: "(" @4:24 (kind=39)
         SliceType:
          LBracket: "[" @4:25 (kind=43)
          RBracket: "]" @4:26 (kind=44)
          NamedType
           Ident: "string" @4:27 (kind=24)
         Size:
          IntLitExpr
           Value: "10" @4:34 (kind=4)
         RParen: ")" @4:36 (kind=40)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("slice_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func main() {
  var x []string = make([]string,10,10);
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
     LBrace: "{" @3:13 (kind=41)
     Stmts
      VarDecl
       Var: "var" @4:3 (kind=11)
       Name: "x" @4:7 (kind=3)
       Type
        SliceType:
         LBracket: "[" @4:9 (kind=43)
         RBracket: "]" @4:10 (kind=44)
         NamedType
          Ident: "string" @4:11 (kind=24)
       Eq: "=" @4:18 (kind=49)
       Init
        MakeExpr:
         Make: "make" @4:20 (kind=3)
         LParen: "(" @4:24 (kind=39)
         SliceType:
          LBracket: "[" @4:25 (kind=43)
          RBracket: "]" @4:26 (kind=44)
          NamedType
           Ident: "string" @4:27 (kind=24)
         Size:
          IntLitExpr
           Value: "10" @4:34 (kind=4)
         Cap:
          IntLitExpr
           Value: "10" @4:37 (kind=4)
         RParen: ")" @4:39 (kind=40)
     RBrace: "}" @5:1 (kind=42)
`
		fmt.Println(ast.Dump(pr))
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.Errors))
	})

	t.Run("bad_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func main() {
  var x map[string]string = make(map[string]string,10,10,10)
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

func main() {
  var x map[string]string = make(mmap[string]string)
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

func main() {
  var x test[string]string = make(mmap[string]string 10,10,10,10 "string")
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})
}
