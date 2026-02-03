package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
	"github.com/stretchr/testify/assert"
)

func TestParser_switch_stmt(t *testing.T) {
	assert := assert.New(t)

	t.Run("empty_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 10},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 11},
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       LBrace: "{" @4:10 (kind=41)
       RBrace: "}" @4:11 (kind=42)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("no_tag_empty_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 10},
			{Kind: token.KWDefault, Value: "default", Line: 5, Column: 5},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 12},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 8, Column: 1},
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       LBrace: "{" @4:10 (kind=41)
       Case: "default" @5:5 (kind=36)
       Colon: ":" @5:12 (kind=47)
       RBrace: "}" @6:3 (kind=42)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("no_tag_empty_x2", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 10},
			{Kind: token.KWDefault, Value: "default", Line: 5, Column: 5},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 12},
			{Kind: token.KWCase, Value: "case", Line: 6, Column: 5},
			{Kind: token.IntLit, Value: "1", Line: 6, Column: 6},
			{Kind: token.Colon, Value: ":", Line: 6, Column: 7},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 8, Column: 1},
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
     LBrace: "{" @3:9 (kind=41)
     Stmts
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       LBrace: "{" @4:10 (kind=41)
       Case: "default" @5:5 (kind=36)
       Colon: ":" @5:12 (kind=47)
       Case: "case" @6:5 (kind=35)
        Values:
         IntLitExpr
          Value: "1" @6:6 (kind=4)
       Colon: ":" @6:7 (kind=47)
       RBrace: "}" @7:3 (kind=42)
     RBrace: "}" @8:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("no_tag_empty_x3", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 10},
			{Kind: token.KWCase, Value: "case", Line: 5, Column: 5},
			{Kind: token.IntLit, Value: "1", Line: 5, Column: 6},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 7},
			{Kind: token.KWDefault, Value: "default", Line: 6, Column: 5},
			{Kind: token.Colon, Value: ":", Line: 6, Column: 12},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 8, Column: 1},
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
     LBrace: "{" @3:9 (kind=41)
     Stmts
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       LBrace: "{" @4:10 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         IntLitExpr
          Value: "1" @5:6 (kind=4)
       Colon: ":" @5:7 (kind=47)
       Case: "default" @6:5 (kind=36)
       Colon: ":" @6:12 (kind=47)
       RBrace: "}" @7:3 (kind=42)
     RBrace: "}" @8:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("no_tag_default", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 10},
			{Kind: token.KWDefault, Value: "default", Line: 5, Column: 5},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 12},
			{Kind: token.KWReturn, Value: "return", Line: 6, Column: 7},
			{Kind: token.Ident, Value: "a", Line: 6, Column: 14},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 8, Column: 1},
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
     LBrace: "{" @3:9 (kind=41)
     Stmts
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       LBrace: "{" @4:10 (kind=41)
       Case: "default" @5:5 (kind=36)
       Colon: ":" @5:12 (kind=47)
        Body:
         ReturnStmt
          Values
           IdentExpr
            Name: "a" @6:14 (kind=3)
       RBrace: "}" @7:3 (kind=42)
     RBrace: "}" @8:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("no_tag_cases_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 10},
			{Kind: token.KWCase, Value: "case", Line: 5, Column: 5},
			{Kind: token.IntLit, Value: "1", Line: 5, Column: 10},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 11},
			{Kind: token.Ident, Value: "b", Line: 6, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 6, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 6, Column: 9},
			{Kind: token.KWDefault, Value: "default", Line: 7, Column: 5},
			{Kind: token.Colon, Value: ":", Line: 7, Column: 12},
			{Kind: token.KWReturn, Value: "return", Line: 8, Column: 7},
			{Kind: token.Ident, Value: "a", Line: 8, Column: 14},
			{Kind: token.RBrace, Value: "}", Line: 9, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 10, Column: 1},
			{Kind: token.EOF, Value: "", Line: 11, Column: 1},
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       LBrace: "{" @4:10 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         IntLitExpr
          Value: "1" @5:10 (kind=4)
       Colon: ":" @5:11 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "b" @6:7 (kind=3)
          LParent: "(" @6:8 (kind=39)
          RParent: ")" @6:9 (kind=40)
       Case: "default" @7:5 (kind=36)
       Colon: ":" @7:12 (kind=47)
        Body:
         ReturnStmt
          Values
           IdentExpr
            Name: "a" @8:14 (kind=3)
       RBrace: "}" @9:3 (kind=42)
     RBrace: "}" @10:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("no_tag_cases_x2", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 10},
			{Kind: token.KWCase, Value: "case", Line: 5, Column: 5},
			{Kind: token.IntLit, Value: "1", Line: 5, Column: 10},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 11},
			{Kind: token.KWReturn, Value: "return", Line: 6, Column: 7},
			{Kind: token.Ident, Value: "b", Line: 6, Column: 14},
			{Kind: token.KWDefault, Value: "default", Line: 7, Column: 5},
			{Kind: token.Colon, Value: ":", Line: 7, Column: 12},
			{Kind: token.KWReturn, Value: "return", Line: 8, Column: 7},
			{Kind: token.Ident, Value: "a", Line: 8, Column: 14},
			{Kind: token.RBrace, Value: "}", Line: 9, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 10, Column: 1},
			{Kind: token.EOF, Value: "", Line: 11, Column: 1},
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       LBrace: "{" @4:10 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         IntLitExpr
          Value: "1" @5:10 (kind=4)
       Colon: ":" @5:11 (kind=47)
        Body:
         ReturnStmt
          Values
           IdentExpr
            Name: "b" @6:14 (kind=3)
       Case: "default" @7:5 (kind=36)
       Colon: ":" @7:12 (kind=47)
        Body:
         ReturnStmt
          Values
           IdentExpr
            Name: "a" @8:14 (kind=3)
       RBrace: "}" @9:3 (kind=42)
     RBrace: "}" @10:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("no_tag_cases_x3", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 10},
			{Kind: token.KWCase, Value: "case", Line: 5, Column: 5},
			{Kind: token.IntLit, Value: "1", Line: 5, Column: 10},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 11},
			{Kind: token.KWCase, Value: "case", Line: 6, Column: 5},
			{Kind: token.IntLit, Value: "2", Line: 6, Column: 10},
			{Kind: token.Colon, Value: ":", Line: 6, Column: 11},
			{Kind: token.Ident, Value: "b", Line: 6, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 6, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 6, Column: 9},
			{Kind: token.Ident, Value: "c", Line: 7, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 7, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 7, Column: 9},
			{Kind: token.KWDefault, Value: "default", Line: 8, Column: 5},
			{Kind: token.Colon, Value: ":", Line: 8, Column: 12},
			{Kind: token.KWReturn, Value: "return", Line: 9, Column: 7},
			{Kind: token.Ident, Value: "a", Line: 9, Column: 14},
			{Kind: token.RBrace, Value: "}", Line: 10, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 11, Column: 1},
			{Kind: token.EOF, Value: "", Line: 12, Column: 1},
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       LBrace: "{" @4:10 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         IntLitExpr
          Value: "1" @5:10 (kind=4)
       Colon: ":" @5:11 (kind=47)
       Case: "case" @6:5 (kind=35)
        Values:
         IntLitExpr
          Value: "2" @6:10 (kind=4)
       Colon: ":" @6:11 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "b" @6:7 (kind=3)
          LParent: "(" @6:8 (kind=39)
          RParent: ")" @6:9 (kind=40)
         CallExpr
          Callee
           IdentExpr
            Name: "c" @7:7 (kind=3)
          LParent: "(" @7:8 (kind=39)
          RParent: ")" @7:9 (kind=40)
       Case: "default" @8:5 (kind=36)
       Colon: ":" @8:12 (kind=47)
        Body:
         ReturnStmt
          Values
           IdentExpr
            Name: "a" @9:14 (kind=3)
       RBrace: "}" @10:3 (kind=42)
     RBrace: "}" @11:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("no_tag_cases_x4", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 10},
			{Kind: token.KWCase, Value: "case", Line: 5, Column: 5},
			{Kind: token.IntLit, Value: "1", Line: 5, Column: 10},
			{Kind: token.Comma, Value: ",", Line: 5, Column: 11},
			{Kind: token.IntLit, Value: "2", Line: 5, Column: 12},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 13},
			{Kind: token.Ident, Value: "b", Line: 6, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 6, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 6, Column: 9},
			{Kind: token.Ident, Value: "c", Line: 7, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 7, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 7, Column: 9},
			{Kind: token.KWDefault, Value: "default", Line: 8, Column: 5},
			{Kind: token.Colon, Value: ":", Line: 8, Column: 12},
			{Kind: token.KWReturn, Value: "return", Line: 9, Column: 7},
			{Kind: token.Ident, Value: "a", Line: 9, Column: 14},
			{Kind: token.RBrace, Value: "}", Line: 10, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 11, Column: 1},
			{Kind: token.EOF, Value: "", Line: 12, Column: 1},
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       LBrace: "{" @4:10 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         IntLitExpr
          Value: "1" @5:10 (kind=4)
         IntLitExpr
          Value: "2" @5:12 (kind=4)
       Colon: ":" @5:13 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "b" @6:7 (kind=3)
          LParent: "(" @6:8 (kind=39)
          RParent: ")" @6:9 (kind=40)
         CallExpr
          Callee
           IdentExpr
            Name: "c" @7:7 (kind=3)
          LParent: "(" @7:8 (kind=39)
          RParent: ")" @7:9 (kind=40)
       Case: "default" @8:5 (kind=36)
       Colon: ":" @8:12 (kind=47)
        Body:
         ReturnStmt
          Values
           IdentExpr
            Name: "a" @9:14 (kind=3)
       RBrace: "}" @10:3 (kind=42)
     RBrace: "}" @11:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("no_tag_cases_x5", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 10},
			{Kind: token.KWCase, Value: "case", Line: 5, Column: 5},
			{Kind: token.IntLit, Value: "1", Line: 5, Column: 10},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 13},
			{Kind: token.Ident, Value: "b", Line: 6, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 6, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 6, Column: 9},
			{Kind: token.KWFallThrough, Value: "fallthrough", Line: 7, Column: 7},
			{Kind: token.KWCase, Value: "case", Line: 8, Column: 5},
			{Kind: token.IntLit, Value: "2", Line: 8, Column: 10},
			{Kind: token.Colon, Value: ":", Line: 8, Column: 13},
			{Kind: token.Ident, Value: "c", Line: 9, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 9, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 9, Column: 9},
			{Kind: token.RBrace, Value: "}", Line: 10, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 11, Column: 1},
			{Kind: token.EOF, Value: "", Line: 12, Column: 1},
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       LBrace: "{" @4:10 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         IntLitExpr
          Value: "1" @5:10 (kind=4)
       Colon: ":" @5:13 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "b" @6:7 (kind=3)
          LParent: "(" @6:8 (kind=39)
          RParent: ")" @6:9 (kind=40)
         FallThrough: "fallthrough" @7:7 (kind=37)
       Case: "case" @8:5 (kind=35)
        Values:
         IntLitExpr
          Value: "2" @8:10 (kind=4)
       Colon: ":" @8:13 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "c" @9:7 (kind=3)
          LParent: "(" @9:8 (kind=39)
          RParent: ")" @9:9 (kind=40)
       RBrace: "}" @10:3 (kind=42)
     RBrace: "}" @11:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("no_tag_cases_x6", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 10},
			{Kind: token.KWCase, Value: "case", Line: 5, Column: 5},
			{Kind: token.IntLit, Value: "1", Line: 5, Column: 10},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 13},
			{Kind: token.Ident, Value: "b", Line: 6, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 6, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 6, Column: 9},
			{Kind: token.KWFallThrough, Value: "fallthrough", Line: 7, Column: 7},
			{Kind: token.Comment, Value: "// fallthrough", Line: 7, Column: 20},
			{Kind: token.KWCase, Value: "case", Line: 8, Column: 5},
			{Kind: token.IntLit, Value: "2", Line: 8, Column: 10},
			{Kind: token.Colon, Value: ":", Line: 8, Column: 13},
			{Kind: token.Ident, Value: "c", Line: 9, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 9, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 9, Column: 9},
			{Kind: token.RBrace, Value: "}", Line: 10, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 11, Column: 1},
			{Kind: token.EOF, Value: "", Line: 12, Column: 1},
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       LBrace: "{" @4:10 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         IntLitExpr
          Value: "1" @5:10 (kind=4)
       Colon: ":" @5:13 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "b" @6:7 (kind=3)
          LParent: "(" @6:8 (kind=39)
          RParent: ")" @6:9 (kind=40)
         FallThrough: "fallthrough" @7:7 (kind=37)
       Case: "case" @8:5 (kind=35)
        Values:
         IntLitExpr
          Value: "2" @8:10 (kind=4)
       Colon: ":" @8:13 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "c" @9:7 (kind=3)
          LParent: "(" @9:8 (kind=39)
          RParent: ")" @9:9 (kind=40)
       RBrace: "}" @10:3 (kind=42)
     RBrace: "}" @11:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("no_tag_cases_x7", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 10},
			{Kind: token.KWCase, Value: "case", Line: 5, Column: 5},
			{Kind: token.Ident, Value: "a", Line: 5, Column: 10},
			{Kind: token.Gt, Value: ">", Line: 5, Column: 11},
			{Kind: token.Ident, Value: "b", Line: 5, Column: 12},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 13},
			{Kind: token.Ident, Value: "b", Line: 6, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 6, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 6, Column: 9},
			{Kind: token.KWFallThrough, Value: "fallthrough", Line: 7, Column: 7},
			{Kind: token.Comment, Value: "// fallthrough", Line: 7, Column: 20},
			{Kind: token.KWCase, Value: "case", Line: 8, Column: 5},
			{Kind: token.IntLit, Value: "2", Line: 8, Column: 10},
			{Kind: token.Colon, Value: ":", Line: 8, Column: 13},
			{Kind: token.Ident, Value: "c", Line: 9, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 9, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 9, Column: 9},
			{Kind: token.RBrace, Value: "}", Line: 10, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 11, Column: 1},
			{Kind: token.EOF, Value: "", Line: 12, Column: 1},
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       LBrace: "{" @4:10 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         BinaryExpr
          IdentExpr
           Name: "a" @5:10 (kind=3)
          Operator: ">" @5:11 (kind=66)
          IdentExpr
           Name: "b" @5:12 (kind=3)
       Colon: ":" @5:13 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "b" @6:7 (kind=3)
          LParent: "(" @6:8 (kind=39)
          RParent: ")" @6:9 (kind=40)
         FallThrough: "fallthrough" @7:7 (kind=37)
       Case: "case" @8:5 (kind=35)
        Values:
         IntLitExpr
          Value: "2" @8:10 (kind=4)
       Colon: ":" @8:13 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "c" @9:7 (kind=3)
          LParent: "(" @9:8 (kind=39)
          RParent: ")" @9:9 (kind=40)
       RBrace: "}" @10:3 (kind=42)
     RBrace: "}" @11:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("tag_cases_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "z", Line: 4, Column: 10},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 11},
			{Kind: token.KWCase, Value: "case", Line: 5, Column: 5},
			{Kind: token.Ident, Value: "a", Line: 5, Column: 10},
			{Kind: token.Gt, Value: ">", Line: 5, Column: 11},
			{Kind: token.Ident, Value: "b", Line: 5, Column: 12},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 13},
			{Kind: token.Ident, Value: "b", Line: 6, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 6, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 6, Column: 9},
			{Kind: token.KWCase, Value: "case", Line: 7, Column: 5},
			{Kind: token.IntLit, Value: "2", Line: 7, Column: 10},
			{Kind: token.Colon, Value: ":", Line: 7, Column: 13},
			{Kind: token.Ident, Value: "c", Line: 8, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 8, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 8, Column: 9},
			{Kind: token.RBrace, Value: "}", Line: 9, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 10, Column: 1},
			{Kind: token.EOF, Value: "", Line: 11, Column: 1},
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       Init:
        IdentExpr
         Name: "z" @4:10 (kind=3)
       LBrace: "{" @4:11 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         BinaryExpr
          IdentExpr
           Name: "a" @5:10 (kind=3)
          Operator: ">" @5:11 (kind=66)
          IdentExpr
           Name: "b" @5:12 (kind=3)
       Colon: ":" @5:13 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "b" @6:7 (kind=3)
          LParent: "(" @6:8 (kind=39)
          RParent: ")" @6:9 (kind=40)
       Case: "case" @7:5 (kind=35)
        Values:
         IntLitExpr
          Value: "2" @7:10 (kind=4)
       Colon: ":" @7:13 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "c" @8:7 (kind=3)
          LParent: "(" @8:8 (kind=39)
          RParent: ")" @8:9 (kind=40)
       RBrace: "}" @9:3 (kind=42)
     RBrace: "}" @10:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("tag_cases_x2", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "z", Line: 4, Column: 10},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 11},
			{Kind: token.KWCase, Value: "case", Line: 5, Column: 5},
			{Kind: token.Ident, Value: "d", Line: 5, Column: 10},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 13},
			{Kind: token.KWCase, Value: "case", Line: 6, Column: 5},
			{Kind: token.Ident, Value: "a", Line: 6, Column: 10},
			{Kind: token.Gt, Value: ">", Line: 6, Column: 11},
			{Kind: token.Ident, Value: "b", Line: 6, Column: 12},
			{Kind: token.Colon, Value: ":", Line: 6, Column: 13},
			{Kind: token.Ident, Value: "b", Line: 7, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 7, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 7, Column: 9},
			{Kind: token.KWCase, Value: "case", Line: 8, Column: 5},
			{Kind: token.IntLit, Value: "2", Line: 8, Column: 10},
			{Kind: token.Colon, Value: ":", Line: 8, Column: 13},
			{Kind: token.Ident, Value: "c", Line: 9, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 9, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 9, Column: 9},
			{Kind: token.RBrace, Value: "}", Line: 10, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 11, Column: 1},
			{Kind: token.EOF, Value: "", Line: 12, Column: 1},
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       Init:
        IdentExpr
         Name: "z" @4:10 (kind=3)
       LBrace: "{" @4:11 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         IdentExpr
          Name: "d" @5:10 (kind=3)
       Colon: ":" @5:13 (kind=47)
       Case: "case" @6:5 (kind=35)
        Values:
         BinaryExpr
          IdentExpr
           Name: "a" @6:10 (kind=3)
          Operator: ">" @6:11 (kind=66)
          IdentExpr
           Name: "b" @6:12 (kind=3)
       Colon: ":" @6:13 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "b" @7:7 (kind=3)
          LParent: "(" @7:8 (kind=39)
          RParent: ")" @7:9 (kind=40)
       Case: "case" @8:5 (kind=35)
        Values:
         IntLitExpr
          Value: "2" @8:10 (kind=4)
       Colon: ":" @8:13 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "c" @9:7 (kind=3)
          LParent: "(" @9:8 (kind=39)
          RParent: ")" @9:9 (kind=40)
       RBrace: "}" @10:3 (kind=42)
     RBrace: "}" @11:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("tag_cases_x3", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "z", Line: 4, Column: 10},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 11},
			{Kind: token.KWDefault, Value: "default", Line: 5, Column: 5},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 13},
			{Kind: token.KWCase, Value: "case", Line: 6, Column: 5},
			{Kind: token.Ident, Value: "a", Line: 6, Column: 10},
			{Kind: token.Gt, Value: ">", Line: 6, Column: 11},
			{Kind: token.Ident, Value: "b", Line: 6, Column: 12},
			{Kind: token.Colon, Value: ":", Line: 6, Column: 13},
			{Kind: token.Ident, Value: "b", Line: 7, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 7, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 7, Column: 9},
			{Kind: token.KWCase, Value: "case", Line: 8, Column: 5},
			{Kind: token.IntLit, Value: "2", Line: 8, Column: 10},
			{Kind: token.Colon, Value: ":", Line: 8, Column: 13},
			{Kind: token.Ident, Value: "c", Line: 9, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 9, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 9, Column: 9},
			{Kind: token.RBrace, Value: "}", Line: 10, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 11, Column: 1},
			{Kind: token.EOF, Value: "", Line: 12, Column: 1},
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       Init:
        IdentExpr
         Name: "z" @4:10 (kind=3)
       LBrace: "{" @4:11 (kind=41)
       Case: "default" @5:5 (kind=36)
       Colon: ":" @5:13 (kind=47)
       Case: "case" @6:5 (kind=35)
        Values:
         BinaryExpr
          IdentExpr
           Name: "a" @6:10 (kind=3)
          Operator: ">" @6:11 (kind=66)
          IdentExpr
           Name: "b" @6:12 (kind=3)
       Colon: ":" @6:13 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "b" @7:7 (kind=3)
          LParent: "(" @7:8 (kind=39)
          RParent: ")" @7:9 (kind=40)
       Case: "case" @8:5 (kind=35)
        Values:
         IntLitExpr
          Value: "2" @8:10 (kind=4)
       Colon: ":" @8:13 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "c" @9:7 (kind=3)
          LParent: "(" @9:8 (kind=39)
          RParent: ")" @9:9 (kind=40)
       RBrace: "}" @10:3 (kind=42)
     RBrace: "}" @11:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("tag_cases_x4", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "z", Line: 4, Column: 10},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 11},
			{Kind: token.KWDefault, Value: "default", Line: 5, Column: 5},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 13},
			{Kind: token.KWCase, Value: "case", Line: 6, Column: 5},
			{Kind: token.Ident, Value: "a", Line: 6, Column: 10},
			{Kind: token.Gt, Value: ">", Line: 6, Column: 11},
			{Kind: token.Ident, Value: "b", Line: 6, Column: 12},
			{Kind: token.Colon, Value: ":", Line: 6, Column: 13},
			{Kind: token.Ident, Value: "b", Line: 7, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 7, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 7, Column: 9},
			{Kind: token.KWCase, Value: "case", Line: 8, Column: 5},
			{Kind: token.IntLit, Value: "2", Line: 8, Column: 10},
			{Kind: token.Colon, Value: ":", Line: 8, Column: 13},
			{Kind: token.KWReturn, Value: "return", Line: 9, Column: 7},
			{Kind: token.Ident, Value: "c", Line: 9, Column: 9},
			{Kind: token.LParen, Value: "(", Line: 9, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 9, Column: 9},
			{Kind: token.RBrace, Value: "}", Line: 10, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 11, Column: 1},
			{Kind: token.EOF, Value: "", Line: 12, Column: 1},
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       Init:
        IdentExpr
         Name: "z" @4:10 (kind=3)
       LBrace: "{" @4:11 (kind=41)
       Case: "default" @5:5 (kind=36)
       Colon: ":" @5:13 (kind=47)
       Case: "case" @6:5 (kind=35)
        Values:
         BinaryExpr
          IdentExpr
           Name: "a" @6:10 (kind=3)
          Operator: ">" @6:11 (kind=66)
          IdentExpr
           Name: "b" @6:12 (kind=3)
       Colon: ":" @6:13 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "b" @7:7 (kind=3)
          LParent: "(" @7:8 (kind=39)
          RParent: ")" @7:9 (kind=40)
       Case: "case" @8:5 (kind=35)
        Values:
         IntLitExpr
          Value: "2" @8:10 (kind=4)
       Colon: ":" @8:13 (kind=47)
        Body:
         ReturnStmt
          Values
           CallExpr
            Callee
             IdentExpr
              Name: "c" @9:9 (kind=3)
            LParent: "(" @9:8 (kind=39)
            RParent: ")" @9:9 (kind=40)
       RBrace: "}" @10:3 (kind=42)
     RBrace: "}" @11:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("tag_cases_x5", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "z", Line: 4, Column: 10},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 11},
			{Kind: token.KWCase, Value: "case", Line: 5, Column: 5},
			{Kind: token.IntLit, Value: "1", Line: 5, Column: 10},
			{Kind: token.Plus, Value: "+", Line: 5, Column: 11},
			{Kind: token.IntLit, Value: "2", Line: 5, Column: 12},
			{Kind: token.Star, Value: "*", Line: 5, Column: 13},
			{Kind: token.IntLit, Value: "3", Line: 5, Column: 14},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 15},
			{Kind: token.KWReturn, Value: "return", Line: 6, Column: 7},
			{Kind: token.IntLit, Value: "7", Line: 6, Column: 9},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 8, Column: 1},
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
     LBrace: "{" @3:9 (kind=41)
     Stmts
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       Init:
        IdentExpr
         Name: "z" @4:10 (kind=3)
       LBrace: "{" @4:11 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         BinaryExpr
          IntLitExpr
           Value: "1" @5:10 (kind=4)
          Operator: "+" @5:11 (kind=51)
          BinaryExpr
           IntLitExpr
            Value: "2" @5:12 (kind=4)
           Operator: "*" @5:13 (kind=57)
           IntLitExpr
            Value: "3" @5:14 (kind=4)
       Colon: ":" @5:15 (kind=47)
        Body:
         ReturnStmt
          Values
           IntLitExpr
            Value: "7" @6:9 (kind=4)
       RBrace: "}" @7:3 (kind=42)
     RBrace: "}" @8:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("tag_cases_x6", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "z", Line: 4, Column: 10},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 11},
			{Kind: token.KWCase, Value: "case", Line: 5, Column: 5},
			{Kind: token.IntLit, Value: "1", Line: 5, Column: 10},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 11},

			{Kind: token.KWSwitch, Value: "switch", Line: 6, Column: 7},
			{Kind: token.Ident, Value: "y", Line: 6, Column: 14},
			{Kind: token.LBrace, Value: "{", Line: 6, Column: 16},
			{Kind: token.KWCase, Value: "case", Line: 7, Column: 9},
			{Kind: token.IntLit, Value: "2", Line: 7, Column: 14},
			{Kind: token.Colon, Value: ":", Line: 7, Column: 15},
			{Kind: token.KWReturn, Value: "return", Line: 8, Column: 7},
			{Kind: token.IntLit, Value: "7", Line: 8, Column: 9},
			{Kind: token.RBrace, Value: "}", Line: 9, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 10, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 11, Column: 1},
			{Kind: token.EOF, Value: "", Line: 12, Column: 1},
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       Init:
        IdentExpr
         Name: "z" @4:10 (kind=3)
       LBrace: "{" @4:11 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         IntLitExpr
          Value: "1" @5:10 (kind=4)
       Colon: ":" @5:11 (kind=47)
        Body:
         SwitchStmt
          Switch: "switch" @6:7 (kind=34)
          Init:
           IdentExpr
            Name: "y" @6:14 (kind=3)
          LBrace: "{" @6:16 (kind=41)
          Case: "case" @7:9 (kind=35)
           Values:
            IntLitExpr
             Value: "2" @7:14 (kind=4)
          Colon: ":" @7:15 (kind=47)
           Body:
            ReturnStmt
             Values
              IntLitExpr
               Value: "7" @8:9 (kind=4)
          RBrace: "}" @9:3 (kind=42)
       RBrace: "}" @10:3 (kind=42)
     RBrace: "}" @11:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("init_tag_cases_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "z", Line: 4, Column: 10},
			{Kind: token.Define, Value: ":=", Line: 4, Column: 11},
			{Kind: token.Ident, Value: "w", Line: 4, Column: 14},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 16},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 17},
			{Kind: token.SemiComma, Value: ";", Line: 4, Column: 18},
			{Kind: token.Ident, Value: "z", Line: 4, Column: 19},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 20},
			{Kind: token.KWCase, Value: "case", Line: 5, Column: 5},
			{Kind: token.Ident, Value: "a", Line: 5, Column: 10},
			{Kind: token.Gt, Value: ">", Line: 5, Column: 11},
			{Kind: token.Ident, Value: "b", Line: 5, Column: 12},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 13},
			{Kind: token.Ident, Value: "b", Line: 6, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 6, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 6, Column: 9},
			{Kind: token.KWCase, Value: "case", Line: 7, Column: 5},
			{Kind: token.IntLit, Value: "2", Line: 7, Column: 10},
			{Kind: token.Colon, Value: ":", Line: 7, Column: 13},
			{Kind: token.Ident, Value: "c", Line: 8, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 8, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 8, Column: 9},
			{Kind: token.RBrace, Value: "}", Line: 10, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 11, Column: 1},
			{Kind: token.EOF, Value: "", Line: 12, Column: 1},
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       Init:
        AssignStmt
         Left
          IdentExpr
           Name: "z" @4:10 (kind=3)
         Operator: ":=" @4:11 (kind=50)
         Right
          CallExpr
           Callee
            IdentExpr
             Name: "w" @4:14 (kind=3)
           LParent: "(" @4:16 (kind=39)
           RParent: ")" @4:17 (kind=40)
       Init:
        IdentExpr
         Name: "z" @4:19 (kind=3)
       LBrace: "{" @4:20 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         BinaryExpr
          IdentExpr
           Name: "a" @5:10 (kind=3)
          Operator: ">" @5:11 (kind=66)
          IdentExpr
           Name: "b" @5:12 (kind=3)
       Colon: ":" @5:13 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "b" @6:7 (kind=3)
          LParent: "(" @6:8 (kind=39)
          RParent: ")" @6:9 (kind=40)
       Case: "case" @7:5 (kind=35)
        Values:
         IntLitExpr
          Value: "2" @7:10 (kind=4)
       Colon: ":" @7:13 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "c" @8:7 (kind=3)
          LParent: "(" @8:8 (kind=39)
          RParent: ")" @8:9 (kind=40)
       RBrace: "}" @10:3 (kind=42)
     RBrace: "}" @11:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("init_tag_cases_x2", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "z", Line: 4, Column: 10},
			{Kind: token.Define, Value: ":=", Line: 4, Column: 11},
			{Kind: token.Ident, Value: "w", Line: 4, Column: 14},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 16},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 17},
			{Kind: token.SemiComma, Value: ";", Line: 4, Column: 18},
			{Kind: token.Ident, Value: "z", Line: 4, Column: 19},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 20},
			{Kind: token.KWDefault, Value: "default", Line: 5, Column: 5},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 6},
			{Kind: token.KWCase, Value: "case", Line: 6, Column: 5},
			{Kind: token.Ident, Value: "a", Line: 6, Column: 10},
			{Kind: token.Gt, Value: ">", Line: 6, Column: 11},
			{Kind: token.Ident, Value: "b", Line: 6, Column: 12},
			{Kind: token.Colon, Value: ":", Line: 6, Column: 13},
			{Kind: token.Ident, Value: "b", Line: 7, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 7, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 7, Column: 9},
			{Kind: token.KWCase, Value: "case", Line: 8, Column: 5},
			{Kind: token.IntLit, Value: "2", Line: 8, Column: 10},
			{Kind: token.Colon, Value: ":", Line: 8, Column: 13},
			{Kind: token.Ident, Value: "c", Line: 9, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 9, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 9, Column: 9},
			{Kind: token.RBrace, Value: "}", Line: 10, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 11, Column: 1},
			{Kind: token.EOF, Value: "", Line: 12, Column: 1},
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       Init:
        AssignStmt
         Left
          IdentExpr
           Name: "z" @4:10 (kind=3)
         Operator: ":=" @4:11 (kind=50)
         Right
          CallExpr
           Callee
            IdentExpr
             Name: "w" @4:14 (kind=3)
           LParent: "(" @4:16 (kind=39)
           RParent: ")" @4:17 (kind=40)
       Init:
        IdentExpr
         Name: "z" @4:19 (kind=3)
       LBrace: "{" @4:20 (kind=41)
       Case: "default" @5:5 (kind=36)
       Colon: ":" @5:6 (kind=47)
       Case: "case" @6:5 (kind=35)
        Values:
         BinaryExpr
          IdentExpr
           Name: "a" @6:10 (kind=3)
          Operator: ">" @6:11 (kind=66)
          IdentExpr
           Name: "b" @6:12 (kind=3)
       Colon: ":" @6:13 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "b" @7:7 (kind=3)
          LParent: "(" @7:8 (kind=39)
          RParent: ")" @7:9 (kind=40)
       Case: "case" @8:5 (kind=35)
        Values:
         IntLitExpr
          Value: "2" @8:10 (kind=4)
       Colon: ":" @8:13 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "c" @9:7 (kind=3)
          LParent: "(" @9:8 (kind=39)
          RParent: ")" @9:9 (kind=40)
       RBrace: "}" @10:3 (kind=42)
     RBrace: "}" @11:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("init_tag_cases_x3", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "z", Line: 4, Column: 10},
			{Kind: token.Define, Value: ":=", Line: 4, Column: 11},
			{Kind: token.Ident, Value: "w", Line: 4, Column: 14},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 16},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 17},
			{Kind: token.SemiComma, Value: ";", Line: 4, Column: 18},
			{Kind: token.Ident, Value: "z", Line: 4, Column: 19},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 20},
			{Kind: token.KWCase, Value: "case", Line: 5, Column: 5},
			{Kind: token.Ident, Value: "a", Line: 5, Column: 10},
			{Kind: token.Gt, Value: ">", Line: 5, Column: 11},
			{Kind: token.Ident, Value: "b", Line: 5, Column: 12},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 13},
			{Kind: token.KWDefault, Value: "default", Line: 6, Column: 5},
			{Kind: token.Colon, Value: ":", Line: 6, Column: 13},
			{Kind: token.Ident, Value: "b", Line: 7, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 7, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 7, Column: 9},
			{Kind: token.KWCase, Value: "case", Line: 8, Column: 5},
			{Kind: token.IntLit, Value: "2", Line: 8, Column: 10},
			{Kind: token.Colon, Value: ":", Line: 8, Column: 13},
			{Kind: token.Ident, Value: "c", Line: 9, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 9, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 9, Column: 9},
			{Kind: token.RBrace, Value: "}", Line: 10, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 11, Column: 1},
			{Kind: token.EOF, Value: "", Line: 12, Column: 1},
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       Init:
        AssignStmt
         Left
          IdentExpr
           Name: "z" @4:10 (kind=3)
         Operator: ":=" @4:11 (kind=50)
         Right
          CallExpr
           Callee
            IdentExpr
             Name: "w" @4:14 (kind=3)
           LParent: "(" @4:16 (kind=39)
           RParent: ")" @4:17 (kind=40)
       Init:
        IdentExpr
         Name: "z" @4:19 (kind=3)
       LBrace: "{" @4:20 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         BinaryExpr
          IdentExpr
           Name: "a" @5:10 (kind=3)
          Operator: ">" @5:11 (kind=66)
          IdentExpr
           Name: "b" @5:12 (kind=3)
       Colon: ":" @5:13 (kind=47)
       Case: "default" @6:5 (kind=36)
       Colon: ":" @6:13 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "b" @7:7 (kind=3)
          LParent: "(" @7:8 (kind=39)
          RParent: ")" @7:9 (kind=40)
       Case: "case" @8:5 (kind=35)
        Values:
         IntLitExpr
          Value: "2" @8:10 (kind=4)
       Colon: ":" @8:13 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "c" @9:7 (kind=3)
          LParent: "(" @9:8 (kind=39)
          RParent: ")" @9:9 (kind=40)
       RBrace: "}" @10:3 (kind=42)
     RBrace: "}" @11:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("init_tag_cases_x4", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "z", Line: 4, Column: 10},
			{Kind: token.Define, Value: ":=", Line: 4, Column: 11},
			{Kind: token.Ident, Value: "w", Line: 4, Column: 14},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 16},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 17},
			{Kind: token.SemiComma, Value: ";", Line: 4, Column: 18},
			{Kind: token.Ident, Value: "z", Line: 4, Column: 19},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 20},
			{Kind: token.KWCase, Value: "case", Line: 5, Column: 5},
			{Kind: token.Ident, Value: "a", Line: 5, Column: 10},
			{Kind: token.Gt, Value: ">", Line: 5, Column: 11},
			{Kind: token.Ident, Value: "b", Line: 5, Column: 12},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 13},
			{Kind: token.KWDefault, Value: "default", Line: 6, Column: 5},
			{Kind: token.Colon, Value: ":", Line: 6, Column: 13},
			{Kind: token.KWCase, Value: "case", Line: 7, Column: 5},
			{Kind: token.IntLit, Value: "2", Line: 7, Column: 10},
			{Kind: token.Colon, Value: ":", Line: 7, Column: 13},
			{Kind: token.Ident, Value: "c", Line: 8, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 8, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 8, Column: 9},
			{Kind: token.RBrace, Value: "}", Line: 9, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 10, Column: 1},
			{Kind: token.EOF, Value: "", Line: 11, Column: 1},
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
      SwitchStmt
       Switch: "switch" @4:3 (kind=34)
       Init:
        AssignStmt
         Left
          IdentExpr
           Name: "z" @4:10 (kind=3)
         Operator: ":=" @4:11 (kind=50)
         Right
          CallExpr
           Callee
            IdentExpr
             Name: "w" @4:14 (kind=3)
           LParent: "(" @4:16 (kind=39)
           RParent: ")" @4:17 (kind=40)
       Init:
        IdentExpr
         Name: "z" @4:19 (kind=3)
       LBrace: "{" @4:20 (kind=41)
       Case: "case" @5:5 (kind=35)
        Values:
         BinaryExpr
          IdentExpr
           Name: "a" @5:10 (kind=3)
          Operator: ">" @5:11 (kind=66)
          IdentExpr
           Name: "b" @5:12 (kind=3)
       Colon: ":" @5:13 (kind=47)
       Case: "default" @6:5 (kind=36)
       Colon: ":" @6:13 (kind=47)
       Case: "case" @7:5 (kind=35)
        Values:
         IntLitExpr
          Value: "2" @7:10 (kind=4)
       Colon: ":" @7:13 (kind=47)
        Body:
         CallExpr
          Callee
           IdentExpr
            Name: "c" @8:7 (kind=3)
          LParent: "(" @8:8 (kind=39)
          RParent: ")" @8:9 (kind=40)
       RBrace: "}" @9:3 (kind=42)
     RBrace: "}" @10:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("bad_no_tag_cases_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 10},
			{Kind: token.KWCase, Value: "case", Line: 5, Column: 5},
			{Kind: token.IntLit, Value: "1", Line: 5, Column: 10},
			{Kind: token.Comma, Value: ",", Line: 5, Column: 11},
			{Kind: token.IntLit, Value: "2", Line: 5, Column: 12},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 13},
			{Kind: token.Ident, Value: "b", Line: 6, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 6, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 6, Column: 9},
			{Kind: token.Ident, Value: "c", Line: 7, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 7, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 7, Column: 9},
			{Kind: token.KWDefault, Value: "default", Line: 8, Column: 5},
			{Kind: token.Colon, Value: ":", Line: 8, Column: 12},
			{Kind: token.KWDefault, Value: "default", Line: 9, Column: 5},
			{Kind: token.Colon, Value: ":", Line: 9, Column: 12},
			{Kind: token.RBrace, Value: "}", Line: 12, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 13, Column: 1},
			{Kind: token.EOF, Value: "", Line: 14, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_no_tag_cases_x2", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 10},
			{Kind: token.KWCase, Value: "case", Line: 5, Column: 5},
			{Kind: token.Comma, Value: ",", Line: 5, Column: 10},
			{Kind: token.IntLit, Value: "1", Line: 5, Column: 11},
			{Kind: token.Comma, Value: ",", Line: 5, Column: 12},
			{Kind: token.IntLit, Value: "2", Line: 5, Column: 13},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 14},
			{Kind: token.Ident, Value: "b", Line: 6, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 6, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 6, Column: 9},
			{Kind: token.Ident, Value: "c", Line: 7, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 7, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 7, Column: 9},
			{Kind: token.KWDefault, Value: "default", Line: 8, Column: 5},
			{Kind: token.Colon, Value: ":", Line: 8, Column: 12},
			{Kind: token.KWReturn, Value: "return", Line: 9, Column: 7},
			{Kind: token.Ident, Value: "a", Line: 9, Column: 14},
			{Kind: token.RBrace, Value: "}", Line: 10, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 11, Column: 1},
			{Kind: token.EOF, Value: "", Line: 12, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_no_tag_cases_x3", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 10},
			{Kind: token.KWCase, Value: "case", Line: 5, Column: 5},
			{Kind: token.IntLit, Value: "1", Line: 5, Column: 10},
			{Kind: token.Comma, Value: ",", Line: 5, Column: 11},
			{Kind: token.IntLit, Value: "2", Line: 5, Column: 12},
			{Kind: token.Comma, Value: ",", Line: 5, Column: 13},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 14},
			{Kind: token.Ident, Value: "b", Line: 6, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 6, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 6, Column: 9},
			{Kind: token.Ident, Value: "c", Line: 7, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 7, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 7, Column: 9},
			{Kind: token.KWDefault, Value: "default", Line: 8, Column: 5},
			{Kind: token.Colon, Value: ":", Line: 8, Column: 12},
			{Kind: token.KWReturn, Value: "return", Line: 9, Column: 7},
			{Kind: token.Ident, Value: "a", Line: 9, Column: 14},
			{Kind: token.RBrace, Value: "}", Line: 10, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 11, Column: 1},
			{Kind: token.EOF, Value: "", Line: 12, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_no_tag_cases_x4", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 10},
			{Kind: token.KWCase, Value: "case", Line: 5, Column: 5},
			{Kind: token.IntLit, Value: "1", Line: 5, Column: 10},
			{Kind: token.Comma, Value: ",", Line: 5, Column: 11},
			{Kind: token.IntLit, Value: "2", Line: 5, Column: 12},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 13},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 14},
			{Kind: token.Ident, Value: "b", Line: 6, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 6, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 6, Column: 9},
			{Kind: token.Ident, Value: "c", Line: 7, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 7, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 7, Column: 9},
			{Kind: token.KWDefault, Value: "default", Line: 8, Column: 5},
			{Kind: token.Colon, Value: ":", Line: 8, Column: 12},
			{Kind: token.KWReturn, Value: "return", Line: 9, Column: 7},
			{Kind: token.Ident, Value: "a", Line: 9, Column: 14},
			{Kind: token.RBrace, Value: "}", Line: 10, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 11, Column: 1},
			{Kind: token.EOF, Value: "", Line: 12, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_no_tag_cases_x5", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 10},
			{Kind: token.KWCase, Value: "case", Line: 5, Column: 5},
			{Kind: token.IntLit, Value: "1", Line: 5, Column: 10},
			{Kind: token.Comma, Value: ",", Line: 5, Column: 11},
			{Kind: token.IntLit, Value: "2", Line: 5, Column: 12},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 13},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 14},
			{Kind: token.Ident, Value: "b", Line: 6, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 6, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 6, Column: 9},
			{Kind: token.Ident, Value: "c", Line: 7, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 7, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 7, Column: 9},
			{Kind: token.KWDefault, Value: "default", Line: 8, Column: 5},
			{Kind: token.Ident, Value: "a", Line: 8, Column: 14},
			{Kind: token.RBrace, Value: "}", Line: 10, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 11, Column: 1},
			{Kind: token.EOF, Value: "", Line: 12, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_no_tag_cases_x6", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 10},
			{Kind: token.KWCase, Value: "case", Line: 5, Column: 5},
			{Kind: token.IntLit, Value: "1", Line: 5, Column: 10},
			{Kind: token.Comma, Value: ",", Line: 5, Column: 11},
			{Kind: token.IntLit, Value: "2", Line: 5, Column: 12},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 13},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 14},
			{Kind: token.Ident, Value: "b", Line: 6, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 6, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 6, Column: 9},
			{Kind: token.Ident, Value: "c", Line: 7, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 7, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 7, Column: 9},
			{Kind: token.KWDefault, Value: "default", Line: 8, Column: 5},
			{Kind: token.RBrace, Value: "}", Line: 9, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 10, Column: 1},
			{Kind: token.EOF, Value: "", Line: 11, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_no_tag_cases_x7", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 10},
			{Kind: token.KWCase, Value: "case", Line: 5, Column: 5},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 13},
			{Kind: token.Ident, Value: "b", Line: 6, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 6, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 6, Column: 9},
			{Kind: token.Ident, Value: "c", Line: 7, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 7, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 7, Column: 9},
			{Kind: token.RBrace, Value: "}", Line: 9, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 10, Column: 1},
			{Kind: token.EOF, Value: "", Line: 11, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("no_tag_cases_x8", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 10},
			{Kind: token.KWCase, Value: "case", Line: 5, Column: 5},
			{Kind: token.IntLit, Value: "1", Line: 5, Column: 10},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 13},
			{Kind: token.Ident, Value: "b", Line: 6, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 6, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 6, Column: 9},
			{Kind: token.KWFallThrough, Value: "fallthrough", Line: 7, Column: 7},
			{Kind: token.KWFallThrough, Value: "fallthrough", Line: 8, Column: 7},
			{Kind: token.KWCase, Value: "case", Line: 9, Column: 5},
			{Kind: token.IntLit, Value: "2", Line: 9, Column: 10},
			{Kind: token.Colon, Value: ":", Line: 9, Column: 13},
			{Kind: token.Ident, Value: "c", Line: 10, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 10, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 10, Column: 9},
			{Kind: token.RBrace, Value: "}", Line: 11, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 12, Column: 1},
			{Kind: token.EOF, Value: "", Line: 13, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_no_tag_cases_x9", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 10},
			{Kind: token.KWCase, Value: "case", Line: 5, Column: 5},
			{Kind: token.IntLit, Value: "1", Line: 5, Column: 10},
			{Kind: token.Comma, Value: ",", Line: 5, Column: 11},
			{Kind: token.IntLit, Value: "2", Line: 5, Column: 12},
			{Kind: token.Comma, Value: ",", Line: 5, Column: 13},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 14},
			{Kind: token.Ident, Value: "b", Line: 6, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 6, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 6, Column: 9},
			{Kind: token.Ident, Value: "c", Line: 7, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 7, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 7, Column: 9},
			{Kind: token.KWDefault, Value: "default", Line: 8, Column: 5},
			{Kind: token.RBrace, Value: "}", Line: 9, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 10, Column: 1},
			{Kind: token.EOF, Value: "", Line: 11, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_no_tag_cases_x10", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 10},
			{Kind: token.KWCase, Value: "case", Line: 5, Column: 5},
			{Kind: token.IntLit, Value: "1", Line: 5, Column: 10},
			{Kind: token.Comma, Value: ",", Line: 5, Column: 11},
			{Kind: token.Comma, Value: ",", Line: 5, Column: 12},
			{Kind: token.IntLit, Value: "2", Line: 5, Column: 13},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 14},
			{Kind: token.Ident, Value: "b", Line: 6, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 6, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 6, Column: 9},
			{Kind: token.Ident, Value: "c", Line: 7, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 7, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 7, Column: 9},
			{Kind: token.KWDefault, Value: "default", Line: 8, Column: 5},
			{Kind: token.RBrace, Value: "}", Line: 9, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 10, Column: 1},
			{Kind: token.EOF, Value: "", Line: 11, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_tag_cases_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "z", Line: 4, Column: 10},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 11},
			{Kind: token.KWCase, Value: "case", Line: 5, Column: 5},
			{Kind: token.Ident, Value: "a", Line: 5, Column: 10},
			{Kind: token.Gt, Value: ">", Line: 5, Column: 11},
			{Kind: token.Ident, Value: "b", Line: 5, Column: 12},
			{Kind: token.Ident, Value: "b", Line: 5, Column: 14},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 13},
			{Kind: token.Ident, Value: "b", Line: 6, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 6, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 6, Column: 9},
			{Kind: token.KWCase, Value: "case", Line: 7, Column: 5},
			{Kind: token.IntLit, Value: "2", Line: 7, Column: 10},
			{Kind: token.Colon, Value: ":", Line: 7, Column: 13},
			{Kind: token.Ident, Value: "c", Line: 8, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 8, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 8, Column: 9},
			{Kind: token.RBrace, Value: "}", Line: 10, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 11, Column: 1},
			{Kind: token.EOF, Value: "", Line: 12, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_init_tag_cases_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "z", Line: 4, Column: 10},
			{Kind: token.Define, Value: ":=", Line: 4, Column: 11},
			{Kind: token.Ident, Value: "w", Line: 4, Column: 14},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 15},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 16},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 17},
			{Kind: token.KWCase, Value: "case", Line: 5, Column: 5},
			{Kind: token.Ident, Value: "a", Line: 5, Column: 10},
			{Kind: token.Gt, Value: ">", Line: 5, Column: 11},
			{Kind: token.Ident, Value: "b", Line: 5, Column: 12},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 13},
			{Kind: token.Ident, Value: "b", Line: 6, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 6, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 6, Column: 9},
			{Kind: token.KWCase, Value: "case", Line: 7, Column: 5},
			{Kind: token.IntLit, Value: "2", Line: 7, Column: 10},
			{Kind: token.Colon, Value: ":", Line: 7, Column: 13},
			{Kind: token.Ident, Value: "c", Line: 8, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 8, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 8, Column: 9},
			{Kind: token.RBrace, Value: "}", Line: 10, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 11, Column: 1},
			{Kind: token.EOF, Value: "", Line: 12, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_init_tag_cases_x2", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWSwitch, Value: "switch", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "z", Line: 4, Column: 10},
			{Kind: token.Define, Value: ":=", Line: 4, Column: 11},
			{Kind: token.Ident, Value: "w", Line: 4, Column: 14},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 16},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 17},
			{Kind: token.SemiComma, Value: ";", Line: 4, Column: 18},
			{Kind: token.Ident, Value: "z", Line: 4, Column: 19},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 20},
			{Kind: token.KWCase, Value: "case", Line: 5, Column: 5},
			{Kind: token.Ident, Value: "a", Line: 5, Column: 10},
			{Kind: token.Gt, Value: ">", Line: 5, Column: 11},
			{Kind: token.Ident, Value: "b", Line: 5, Column: 12},
			{Kind: token.SemiComma, Value: ";", Line: 5, Column: 13},
			{Kind: token.Comma, Value: ",", Line: 5, Column: 14},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 1},
			{Kind: token.Ident, Value: "b", Line: 6, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 6, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 6, Column: 9},
			{Kind: token.KWCase, Value: "case", Line: 7, Column: 5},
			{Kind: token.IntLit, Value: "2", Line: 7, Column: 10},
			{Kind: token.Colon, Value: ":", Line: 7, Column: 13},
			{Kind: token.Ident, Value: "c", Line: 8, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 8, Column: 8},
			{Kind: token.RParen, Value: ")", Line: 8, Column: 9},
			{Kind: token.RBrace, Value: "}", Line: 10, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 11, Column: 1},
			{Kind: token.EOF, Value: "", Line: 12, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}
