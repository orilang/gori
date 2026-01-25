package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
	"github.com/stretchr/testify/assert"
)

func TestParser_if_stmt(t *testing.T) {
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
			{Kind: token.KWIf, Value: "if", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 6},
			{Kind: token.Lt, Value: "<", Line: 4, Column: 7},
			{Kind: token.Ident, Value: "b", Line: 4, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 10},
			{Kind: token.KWReturn, Value: "return", Line: 5, Column: 5},
			{Kind: token.Ident, Value: "b", Line: 5, Column: 12},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 1},
			{Kind: token.KWElse, Value: "else", Line: 6, Column: 5},
			{Kind: token.KWIf, Value: "if", Line: 6, Column: 10},
			{Kind: token.Ident, Value: "a", Line: 6, Column: 13},
			{Kind: token.Lt, Value: "<", Line: 6, Column: 14},
			{Kind: token.Ident, Value: "c", Line: 6, Column: 15},
			{Kind: token.LBrace, Value: "{", Line: 6, Column: 17},
			{Kind: token.KWReturn, Value: "return", Line: 7, Column: 5},
			{Kind: token.Ident, Value: "c", Line: 7, Column: 12},
			{Kind: token.RBrace, Value: "}", Line: 8, Column: 1},
			{Kind: token.KWElse, Value: "else", Line: 8, Column: 5},
			{Kind: token.LBrace, Value: "{", Line: 8, Column: 17},
			{Kind: token.KWReturn, Value: "return", Line: 9, Column: 5},
			{Kind: token.Ident, Value: "a", Line: 9, Column: 12},
			{Kind: token.RBrace, Value: "}", Line: 10, Column: 1},
			{Kind: token.RBrace, Value: "}", Line: 12, Column: 1},
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
      IfStmt
       Condition
        BinaryExpr
         IdentExpr
          Name: "a" @4:6 (kind=3)
         Operator: "<" @4:7 (kind=64)
         IdentExpr
          Name: "b" @4:8 (kind=3)
       Then
        BlockStmt
         LBrace: "{" @4:10 (kind=41)
         Stmts
          ReturnStmt
           Values
            IdentExpr
             Name: "b" @5:12 (kind=3)
         RBrace: "}" @6:1 (kind=42)
       Else
        IfStmt
         Condition
          BinaryExpr
           IdentExpr
            Name: "a" @6:13 (kind=3)
           Operator: "<" @6:14 (kind=64)
           IdentExpr
            Name: "c" @6:15 (kind=3)
         Then
          BlockStmt
           LBrace: "{" @6:17 (kind=41)
           Stmts
            ReturnStmt
             Values
              IdentExpr
               Name: "c" @7:12 (kind=3)
           RBrace: "}" @8:1 (kind=42)
         Else
          BlockStmt
           LBrace: "{" @8:17 (kind=41)
           Stmts
            ReturnStmt
             Values
              IdentExpr
               Name: "a" @9:12 (kind=3)
           RBrace: "}" @10:1 (kind=42)
     RBrace: "}" @12:1 (kind=42)
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
			{Kind: token.KWIf, Value: "if", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 6},
			{Kind: token.Lt, Value: "<", Line: 4, Column: 7},
			{Kind: token.Ident, Value: "b", Line: 4, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 10},
			{Kind: token.KWIf, Value: "if", Line: 5, Column: 5},
			{Kind: token.Ident, Value: "a", Line: 5, Column: 8},
			{Kind: token.Lt, Value: "<", Line: 5, Column: 9},
			{Kind: token.Ident, Value: "c", Line: 5, Column: 10},
			{Kind: token.LBrace, Value: "{", Line: 5, Column: 12},
			{Kind: token.KWReturn, Value: "return", Line: 6, Column: 7},
			{Kind: token.Ident, Value: "c", Line: 6, Column: 14},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 8, Column: 1},
			{Kind: token.RBrace, Value: "}", Line: 9, Column: 1},
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
      IfStmt
       Condition
        BinaryExpr
         IdentExpr
          Name: "a" @4:6 (kind=3)
         Operator: "<" @4:7 (kind=64)
         IdentExpr
          Name: "b" @4:8 (kind=3)
       Then
        BlockStmt
         LBrace: "{" @4:10 (kind=41)
         Stmts
          IfStmt
           Condition
            BinaryExpr
             IdentExpr
              Name: "a" @5:8 (kind=3)
             Operator: "<" @5:9 (kind=64)
             IdentExpr
              Name: "c" @5:10 (kind=3)
           Then
            BlockStmt
             LBrace: "{" @5:12 (kind=41)
             Stmts
              ReturnStmt
               Values
                IdentExpr
                 Name: "c" @6:14 (kind=3)
             RBrace: "}" @7:3 (kind=42)
         RBrace: "}" @8:1 (kind=42)
     RBrace: "}" @9:1 (kind=42)
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
			{Kind: token.KWIf, Value: "if", Line: 4, Column: 6},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 7},
			{Kind: token.Lt, Value: "<", Line: 4, Column: 8},
			{Kind: token.Ident, Value: "b", Line: 4, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 5, Column: 1},
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
			{Kind: token.KWIf, Value: "if", Line: 4, Column: 3},
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
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWIf, Value: "if", Line: 4, Column: 3},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 5},
			{Kind: token.RBrace, Value: "}", Line: 3, Column: 6},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_x4", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWElse, Value: "else", Line: 4, Column: 3},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 5},
			{Kind: token.RBrace, Value: "}", Line: 3, Column: 6},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_x5", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWIf, Value: "if", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 5},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 6},
			{Kind: token.KWElse, Value: "else", Line: 5, Column: 3},
			{Kind: token.LBrace, Value: "{", Line: 5, Column: 5},
			{Kind: token.RBrace, Value: "}", Line: 5, Column: 6},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 6},
			{Kind: token.EOF, Value: "", Line: 7, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_x6", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWIf, Value: "if", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 5},
			{Kind: token.Assign, Value: "=", Line: 4, Column: 6},
			{Kind: token.IntLit, Value: "1", Line: 4, Column: 7},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 8},
			{Kind: token.RBrace, Value: "}", Line: 5, Column: 6},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 6},
			{Kind: token.EOF, Value: "", Line: 7, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_x7", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWIf, Value: "if", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 6},
			{Kind: token.Lt, Value: "<", Line: 4, Column: 7},
			{Kind: token.Ident, Value: "b", Line: 4, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 10},
			{Kind: token.KWReturn, Value: "return", Line: 5, Column: 5},
			{Kind: token.Ident, Value: "b", Line: 5, Column: 12},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 1},
			{Kind: token.KWElse, Value: "else", Line: 6, Column: 5},
			{Kind: token.KWReturn, Value: "return", Line: 6, Column: 10},
			{Kind: token.Ident, Value: "c", Line: 6, Column: 17},
			{Kind: token.RBrace, Value: "}", Line: 12, Column: 1},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}
