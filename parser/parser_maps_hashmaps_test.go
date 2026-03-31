package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParser_maps_hashmaps(t *testing.T) {
	assert := assert.New(t)

	t.Run("map", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func main(){
  var x map[string]string=make(map[string]string)
  var y hashmap[string]string=make(hashmap[string]string)
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
        Map: "map" @4:9 (kind=79)
        LBracket: "[" @4:12 (kind=43)
        KeyType:
         Name: "string" @4:13 (kind=24)
        RBracket: "]" @4:19 (kind=44)
        ValueType:
         Name: "string" @4:13 (kind=24)
       Eq: "=" @4:26 (kind=49)
       Init
        Make: "make" @4:27 (kind=3)
        LParen: "(" @4:31 (kind=39)
        Map: "map" @4:32 (kind=79)
        LBracket: "[" @4:35 (kind=43)
        KeyType:
         Name: "string" @4:36 (kind=24)
        RBracket: "]" @4:42 (kind=44)
        ValueType:
         Name: "string" @4:36 (kind=24)
        RParen: ")" @4:49 (kind=40)
      VarDeclStmt
       Var: "var" @5:3 (kind=11)
       Name: "y" @5:7 (kind=3)
       Type
        Hashmap: "hashmap" @5:9 (kind=80)
        LBracket: "[" @5:16 (kind=43)
        KeyType:
         Name: "string" @5:17 (kind=24)
        RBracket: "]" @5:23 (kind=44)
        ValueType:
         Name: "string" @5:17 (kind=24)
       Eq: "=" @5:30 (kind=49)
       Init
        Make: "make" @5:31 (kind=3)
        LParen: "(" @5:35 (kind=39)
        Hashmap: "hashmap" @5:36 (kind=80)
        LBracket: "[" @5:43 (kind=43)
        KeyType:
         Name: "string" @5:44 (kind=24)
        RBracket: "]" @5:50 (kind=44)
        ValueType:
         Name: "string" @5:44 (kind=24)
        RParen: ")" @5:57 (kind=40)
     RBrace: "}" @6:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("bad_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

func main(){
  var x map[struct]string=make(map[struct]string)
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

func main(){
  var x map[string]struct=make(map[struct]string)
}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}
