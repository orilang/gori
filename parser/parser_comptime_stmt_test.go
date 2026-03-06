package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
	"github.com/stretchr/testify/assert"
)

func TestParser_comptime_stmt(t *testing.T) {
	assert := assert.New(t)

	t.Run("const_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWComptime, Value: "comptime", Line: 3, Column: 1},
			{Kind: token.KWConst, Value: "const", Line: 3, Column: 10},
			{Kind: token.Ident, Value: "a", Line: 3, Column: 16},
			{Kind: token.KWFloat, Value: "float", Line: 3, Column: 18},
			{Kind: token.Assign, Value: "=", Line: 3, Column: 24},
			{Kind: token.FloatLit, Value: "3.14", Line: 3, Column: 26},
			{Kind: token.EOF, Value: "", Line: 3, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 ComptimeStmt
  Comptime: "comptime" @3:1 (kind=78)
  ConstDecl
   Const: "const" @3:10 (kind=23)
   Name: "a" @3:16 (kind=3)
   Type
    NameType
     Name: "float" @3:18 (kind=20)
   Eq: "=" @3:24 (kind=49)
   Init
    FloatLitExpr
     Value: "3.14" @3:26 (kind=5)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("func_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWComptime, Value: "comptime", Line: 3, Column: 1},
			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 10},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 12},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 13},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 14},
			{Kind: token.LBracket, Value: "[", Line: 3, Column: 15},
			{Kind: token.RBracket, Value: "]", Line: 3, Column: 16},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 19},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 20},
			{Kind: token.RBrace, Value: "}", Line: 3, Column: 21},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 ComptimeStmt
  Comptime: "comptime" @3:1 (kind=78)
  FuncDecl
   Function: "func" @3:10 (kind=10)
   Name: "x" @3:12 (kind=3)
   Params
    (none)
   Results
     Param
      Type
         LBracket: "[" @3:15 (kind=43)
         RBracket: "]" @3:16 (kind=44)
         Ident: "int" @3:19 (kind=12)
   Body
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("bad_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWComptime, Value: "comptime", Line: 3, Column: 1},
			{Kind: token.KWVar, Value: "var", Line: 3, Column: 10},
			{Kind: token.Ident, Value: "a", Line: 3, Column: 16},
			{Kind: token.KWFloat, Value: "float", Line: 3, Column: 18},
			{Kind: token.Assign, Value: "=", Line: 3, Column: 24},
			{Kind: token.FloatLit, Value: "3.14", Line: 3, Column: 26},
			{Kind: token.EOF, Value: "", Line: 3, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_x2", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWConst, Value: "const", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "a", Line: 3, Column: 6},
			{Kind: token.KWFloat, Value: "float", Line: 3, Column: 13},
			{Kind: token.Assign, Value: "=", Line: 3, Column: 14},
			{Kind: token.KWComptime, Value: "comptime", Line: 3, Column: 16},
			{Kind: token.FloatLit, Value: "3.14", Line: 3, Column: 26},
			{Kind: token.EOF, Value: "", Line: 3, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}
