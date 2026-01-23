package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
	"github.com/stretchr/testify/assert"
)

func TestParser_parse_file(t *testing.T) {
	assert := assert.New(t)

	t.Run("error_global_var_forbidden", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "a"},
			{Kind: token.KWInt, Value: "int"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "0"},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 10},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 13},
			{Kind: token.RBrace, Value: "}", Line: 8, Column: 1},
			{Kind: token.EOF, Value: "", Line: 9, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("function", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWConst, Value: "const"},
			{Kind: token.Ident, Value: "a"},
			{Kind: token.KWFloat, Value: "float"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.FloatLit, Value: "3.14"},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 10},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 13},

			{Kind: token.KWConst, Value: "const"},
			{Kind: token.Ident, Value: "ab"},
			{Kind: token.KWFloat, Value: "float"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.FloatLit, Value: "3.14"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "a"},
			{Kind: token.KWInt, Value: "int"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "0"},

			{Kind: token.RBrace, Value: "}", Line: 8, Column: 1},
			{Kind: token.EOF, Value: "", Line: 9, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 ConstDecls
  ConstDecl
   Const: "const" @0:0 (kind=23)
   Name: "a" @0:0 (kind=3)
   Type
    NameType
     Name: "float" @0:0 (kind=20)
   Eq: "=" @0:0 (kind=49)
   Init
    FloatLitExpr
     Value: "3.14" @0:0 (kind=5)
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
      ConstDecl
       Const: "const" @0:0 (kind=23)
       Name: "ab" @0:0 (kind=3)
       Type
        NameType
         Name: "float" @0:0 (kind=20)
       Eq: "=" @0:0 (kind=49)
       Init
        FloatLitExpr
         Value: "3.14" @0:0 (kind=5)
      VarDeclStmt
       Var: "var" @0:0 (kind=11)
       Name: "a" @0:0 (kind=3)
       Type
        NameType
         Name: "int" @0:0 (kind=12)
       Eq: "=" @0:0 (kind=49)
       Init
        IntLitExpr
         Value: "0" @0:0 (kind=4)
     RBrace: "}" @8:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("decls_none", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  (none)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_params", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "dummy", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},

			{Kind: token.Ident, Value: "a", Line: 3, Column: 8},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 10},

			{Kind: token.Comma, Value: ",", Line: 3, Column: 13},

			{Kind: token.Ident, Value: "b", Line: 3, Column: 15},
			{Kind: token.KWString, Value: "string", Line: 3, Column: 17},

			{Kind: token.Comma, Value: ",", Line: 3, Column: 23},

			{Kind: token.Ident, Value: "c", Line: 3, Column: 25},
			{Kind: token.Ident, Value: "string", Line: 3, Column: 27},

			{Kind: token.RParen, Value: ")", Line: 3, Column: 28},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 29},

			{Kind: token.RBrace, Value: "}", Line: 3, Column: 30},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
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
     Function: "a" @3:8 (kind=3)
     Type
      NameType
       Name: "int" @3:10 (kind=12)
    Param
     Function: "b" @3:15 (kind=3)
     Type
      NameType
       Name: "string" @3:17 (kind=24)
    Param
     Function: "c" @3:25 (kind=3)
     Type
      NameType
       Name: "string" @3:27 (kind=3)
   Body
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("bad_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "dummy", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},

			{Kind: token.Comma, Value: ",", Line: 3, Column: 8},
			{Kind: token.Ident, Value: "a", Line: 3, Column: 9},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 10},

			{Kind: token.RParen, Value: ")", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 12},

			{Kind: token.RBrace, Value: "}", Line: 3, Column: 30},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_x2", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "dummy", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},

			{Kind: token.Ident, Value: "a", Line: 3, Column: 8},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 9},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 10},

			{Kind: token.RParen, Value: ")", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 12},

			{Kind: token.RBrace, Value: "}", Line: 3, Column: 30},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_x3", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "dummy", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},

			{Kind: token.Ident, Value: "a", Line: 3, Column: 8},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 9},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 10},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 11},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 12},
			{Kind: token.Ident, Value: "a", Line: 3, Column: 13},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 14},

			{Kind: token.RParen, Value: ")", Line: 3, Column: 15},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 16},

			{Kind: token.RBrace, Value: "}", Line: 3, Column: 30},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_x3", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "dummy", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},

			{Kind: token.Ident, Value: "a", Line: 3, Column: 8},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 9},
			{Kind: token.Ident, Value: "a", Line: 3, Column: 10},

			{Kind: token.RParen, Value: ")", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 12},

			{Kind: token.RBrace, Value: "}", Line: 3, Column: 30},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}
