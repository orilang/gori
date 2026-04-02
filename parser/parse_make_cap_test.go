package parser

import (
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
  var x map[string]string = make(map[string]string)
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
		assert.Equal(0, len(parser.errors))
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
		assert.Equal(0, len(parser.errors))
	})

	t.Run("slice_x1", func(t *testing.T) {
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
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
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
		assert.Greater(len(parser.errors), 0)
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
		assert.Greater(len(parser.errors), 0)
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
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_x4", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func main() {
  var x map[string]string = make(map[string]string,"plop")
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}
