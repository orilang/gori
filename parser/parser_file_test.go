package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParser_parse_file(t *testing.T) {
	assert := assert.New(t)

	t.Run("error_global_var_forbidden", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `var a int = 0

func main(){}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("function", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main

const a float = 3.14
func main(){
  const ab float = 3.14
	var a int = 0
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
   Name: "a" @3:7 (kind=3)
   Type
    NamedType
     Ident: "float" @3:9 (kind=20)
   Eq: "=" @3:15 (kind=49)
   Init
    FloatLitExpr
     Value: "3.14" @3:17 (kind=5)
  FuncDecl
   Function: "func" @4:1 (kind=10)
   Name: "main" @4:6 (kind=3)
   Params
    (none)
   Body
    BlockStmt
     LBrace: "{" @4:12 (kind=41)
     Stmts
      ConstDecl
       Const: "const" @5:3 (kind=23)
       Name: "ab" @5:9 (kind=3)
       Type
        NamedType
         Ident: "float" @5:12 (kind=20)
       Eq: "=" @5:18 (kind=49)
       Init
        FloatLitExpr
         Value: "3.14" @5:20 (kind=5)
      VarDecl
       Var: "var" @6:2 (kind=11)
       Name: "a" @6:6 (kind=3)
       Type
        NamedType
         Ident: "int" @6:8 (kind=12)
       Eq: "=" @6:12 (kind=49)
       Init
        IntLitExpr
         Value: "0" @6:14 (kind=4)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("decls_none", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package main
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_params", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package dummy

func x(a int, b string, c string){}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "dummy" @1:9 (kind=3)
 Decls
  FuncDecl
   Function: "func" @3:1 (kind=10)
   Name: "x" @3:6 (kind=3)
   Params
    Param
     Ident: "a" @3:8 (kind=3)
     Type
      NamedType
       Ident: "int" @3:10 (kind=12)
    Param
     Ident: "b" @3:15 (kind=3)
     Type
      NamedType
       Ident: "string" @3:17 (kind=24)
    Param
     Ident: "c" @3:25 (kind=3)
     Type
      NamedType
       Ident: "string" @3:27 (kind=24)
   Body
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("bad_x1", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package dummy

func x(, a string){}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_x2", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package dummy

func x(a int,){}
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
		data := `package dummy

func x(a int,,,, a int){}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_x4", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `package dummy

func x(a int a){}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}
