package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
	"github.com/stretchr/testify/assert"
)

func TestParser_stmt(t *testing.T) {
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
			{Kind: token.Ident, Value: "a", Line: 4, Column: 3},
			{Kind: token.Assign, Value: "=", Line: 4, Column: 5},
			{Kind: token.IntLit, Value: "1", Line: 4, Column: 7},
			{Kind: token.RBrace, Value: "}", Line: 5, Column: 1},
			{Kind: token.EOF, Value: "", Line: 9, Column: 1},
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
       Operator: "=" @4:5 (kind=49)
       Right
        IntLitExpr
         Value: "1" @4:7 (kind=4)
     RBrace: "}" @5:1 (kind=42)
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
			{Kind: token.Ident, Value: "a", Line: 3, Column: 8},
			{Kind: token.Ident, Value: "z", Line: 3, Column: 10},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 13},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 3},
			{Kind: token.Dot, Value: ".", Line: 4, Column: 4},
			{Kind: token.Ident, Value: "b", Line: 4, Column: 5},
			{Kind: token.PlusEq, Value: "+=", Line: 4, Column: 7},
			{Kind: token.IntLit, Value: "1", Line: 4, Column: 10},
			{Kind: token.RBrace, Value: "}", Line: 5, Column: 1},
			{Kind: token.EOF, Value: "", Line: 6, Column: 1},
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
    Param
     Function: "a" @3:8 (kind=3)
     Type
      NameType
       Name: "z" @3:10 (kind=3)
   Body
    BlockStmt
     LBrace: "{" @3:13 (kind=41)
     Stmts
      AssignStmt
       Left
        SelectorExpr
         X:
          IdentExpr
           Name: "a" @4:3 (kind=3)
         Dot: "." @4:4 (kind=48)
         Selector: "b" @4:5 (kind=3)
       Operator: "+=" @4:7 (kind=52)
       Right
        IntLitExpr
         Value: "1" @4:10 (kind=4)
     RBrace: "}" @5:1 (kind=42)
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
			{Kind: token.Ident, Value: "a", Line: 3, Column: 8},
			{Kind: token.Ident, Value: "z", Line: 3, Column: 10},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 13},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 3},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 4},
			{Kind: token.Ident, Value: "i", Line: 4, Column: 5},
			{Kind: token.RBracket, Value: "]", Line: 4, Column: 6},
			{Kind: token.Assign, Value: "=", Line: 4, Column: 7},
			{Kind: token.IntLit, Value: "3", Line: 4, Column: 8},
			{Kind: token.RBrace, Value: "}", Line: 5, Column: 1},
			{Kind: token.EOF, Value: "", Line: 6, Column: 1},
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
    Param
     Function: "a" @3:8 (kind=3)
     Type
      NameType
       Name: "z" @3:10 (kind=3)
   Body
    BlockStmt
     LBrace: "{" @3:13 (kind=41)
     Stmts
      AssignStmt
       Left
        IndexExpr
         X:
         IdentExpr
          Name: "a" @4:3 (kind=3)
         LBracket: "[" @4:4 (kind=43)
          IdentExpr
           Name: "i" @4:5 (kind=3)
         RBracket: "]" @4:6 (kind=44)
       Operator: "=" @4:7 (kind=49)
       Right
        IntLitExpr
         Value: "3" @4:8 (kind=4)
     RBrace: "}" @5:1 (kind=42)
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
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 3},
			{Kind: token.Assign, Value: "=", Line: 4, Column: 5},
			{Kind: token.Ident, Value: "f", Line: 4, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 9},
			{Kind: token.RBrace, Value: "}", Line: 5, Column: 1},
			{Kind: token.EOF, Value: "", Line: 6, Column: 1},
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
     LBrace: "{" @3:9 (kind=41)
     Stmts
      AssignStmt
       Left
        IdentExpr
         Name: "a" @4:3 (kind=3)
       Operator: "=" @4:5 (kind=49)
       Right
        CallExpr
         Callee
          IdentExpr
           Name: "f" @4:7 (kind=3)
         LParent: "(" @4:8 (kind=39)
         RParent: ")" @4:9 (kind=40)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_x5", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.Ident, Value: "f", Line: 4, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 4},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 5},
			{Kind: token.RBrace, Value: "}", Line: 5, Column: 1},
			{Kind: token.EOF, Value: "", Line: 6, Column: 1},
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
     LBrace: "{" @3:9 (kind=41)
     Stmts
      CallExpr
       Callee
        IdentExpr
         Name: "f" @4:3 (kind=3)
       LParent: "(" @4:4 (kind=39)
       RParent: ")" @4:5 (kind=40)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})
}
