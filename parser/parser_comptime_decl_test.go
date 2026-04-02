package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParser_comptime_decl(t *testing.T) {
	assert := assert.New(t)

	t.Run("const_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

comptime const a float = 3.14
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  CompTimeBlockDecl:
   Comptime: "comptime" @3:1 (kind=78)
    ConstDecl
     Const: "const" @3:10 (kind=23)
     Name: "a" @3:16 (kind=3)
     Type
      NamedType
       Ident: "float" @3:18 (kind=20)
     Eq: "=" @3:24 (kind=49)
     Init
      FloatLitExpr
       Value: "3.14" @3:26 (kind=5)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("func_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

comptime func x()[]int{}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  CompTimeBlockDecl:
   Comptime: "comptime" @3:1 (kind=78)
    FuncDecl
     Function: "func" @3:10 (kind=10)
     Name: "x" @3:15 (kind=3)
     Params
      (none)
     Results
       Param
        Type
         SliceType:
          LBracket: "[" @3:18 (kind=43)
          RBracket: "]" @3:19 (kind=44)
          NamedType
           Ident: "int" @3:20 (kind=12)
     Body
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("bad_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

comptime var a float = 3.14
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

const a float = comptime 3.14
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}
