package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
	"github.com/stretchr/testify/assert"
)

func TestParser_return_stmt(t *testing.T) {
	assert := assert.New(t)

	t.Run("function_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWReturn, Value: "return", Line: 4, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 1},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
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
     LBrace: "{" @3:10 (kind=41)
     Stmts
      ReturnStmt
     RBrace: "}" @4:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_x2", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWReturn, Value: "return", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "a", Line: 3, Column: 5},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 1},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
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
     LBrace: "{" @3:10 (kind=41)
     Stmts
      ReturnStmt
       Values
        IdentExpr
         Name: "a" @3:5 (kind=3)
     RBrace: "}" @4:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_x3", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 3},
			{Kind: token.Define, Value: ":=", Line: 4, Column: 4},
			{Kind: token.StringLit, Value: "test", Line: 4, Column: 6},
			{Kind: token.KWReturn, Value: "return", Line: 5, Column: 3},
			{Kind: token.Ident, Value: "a", Line: 5, Column: 10},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 1},
			{Kind: token.EOF, Value: "", Line: 7, Column: 1},
		}

		parser := New(input)
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
     LBrace: "{" @3:10 (kind=41)
     Stmts
      AssignStmt
       Left
        IdentExpr
         Name: "a" @4:3 (kind=3)
       Operator: ":=" @4:4 (kind=50)
       Right
        StringLitExpr
         Value: "test" @4:6 (kind=6)
      ReturnStmt
       Values
        IdentExpr
         Name: "a" @5:10 (kind=3)
     RBrace: "}" @6:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_x4", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWReturn, Value: "return", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "a", Line: 3, Column: 5},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 6},
			{Kind: token.Ident, Value: "b", Line: 3, Column: 7},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 1},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
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
     LBrace: "{" @3:10 (kind=41)
     Stmts
      ReturnStmt
       Values
        IdentExpr
         Name: "a" @3:5 (kind=3)
        IdentExpr
         Name: "b" @3:7 (kind=3)
     RBrace: "}" @4:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("bad_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWReturn, Value: "return", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "a", Line: 3, Column: 5},
			{Kind: token.Ident, Value: "b", Line: 3, Column: 7},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 1},
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
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWReturn, Value: "return", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "a", Line: 3, Column: 5},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 6},
			{Kind: token.Ident, Value: "b", Line: 3, Column: 7},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 8},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 1},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}
