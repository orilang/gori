package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
	"github.com/stretchr/testify/assert"
)

func TestParser_slice_stmt(t *testing.T) {
	assert := assert.New(t)

	t.Run("const_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWConst, Value: "const", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 7},
			{Kind: token.LBracket, Value: "[", Line: 3, Column: 9},
			{Kind: token.RBracket, Value: "]", Line: 3, Column: 10},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 11},
			{Kind: token.Assign, Value: "=", Line: 3, Column: 14},
			{Kind: token.LBracket, Value: "[", Line: 3, Column: 14},
			{Kind: token.RBracket, Value: "]", Line: 3, Column: 16},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 17},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 20},
			{Kind: token.IntLit, Value: "1", Line: 3, Column: 21},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 22},
			{Kind: token.IntLit, Value: "2", Line: 3, Column: 23},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 24},
			{Kind: token.IntLit, Value: "3", Line: 3, Column: 24},
			{Kind: token.RBrace, Value: "}", Line: 3, Column: 26},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 ConstDecls
  ConstDecl
   Const: "const" @3:1 (kind=23)
   Name: "x" @3:7 (kind=3)
   Type:
    LBracket: "[" @3:9 (kind=43)
    RBracket: "]" @3:10 (kind=44)
    Ident: "int" @3:11 (kind=12)
   Eq: "=" @3:14 (kind=49)
    Type:
     LBracket: "[" @3:14 (kind=43)
     RBracket: "]" @3:16 (kind=44)
     Ident: "int" @3:17 (kind=12)
    LBrace: "{" @3:20 (kind=41)
     Elements
      IntLitExpr
       Value: "1" @3:21 (kind=4)
      IntLitExpr
       Value: "2" @3:23 (kind=4)
      IntLitExpr
       Value: "3" @3:24 (kind=4)
    RBrace: "}" @3:26 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("const_x2", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 10},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 12},
			{Kind: token.KWConst, Value: "const", Line: 4, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 4, Column: 7},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 9},
			{Kind: token.RBracket, Value: "]", Line: 4, Column: 10},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 11},
			{Kind: token.Assign, Value: "=", Line: 4, Column: 14},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 14},
			{Kind: token.RBracket, Value: "]", Line: 4, Column: 16},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 17},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 20},
			{Kind: token.IntLit, Value: "1", Line: 4, Column: 21},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 22},
			{Kind: token.IntLit, Value: "2", Line: 4, Column: 23},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 24},
			{Kind: token.IntLit, Value: "3", Line: 4, Column: 24},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 26},
			{Kind: token.SemiComma, Value: ";", Line: 4, Column: 26},

			{Kind: token.KWConst, Value: "const", Line: 4, Column: 30},
			{Kind: token.Ident, Value: "y", Line: 4, Column: 36},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 38},
			{Kind: token.RBracket, Value: "]", Line: 4, Column: 39},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 40},
			{Kind: token.Assign, Value: "=", Line: 4, Column: 43},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 44},
			{Kind: token.RBracket, Value: "]", Line: 4, Column: 45},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 46},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 49},
			{Kind: token.IntLit, Value: "1", Line: 4, Column: 50},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 51},
			{Kind: token.IntLit, Value: "2", Line: 4, Column: 52},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 53},
			{Kind: token.IntLit, Value: "3", Line: 4, Column: 54},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 55},
			{Kind: token.RBrace, Value: "}", Line: 5, Column: 56},
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
   Name: "main" @3:6 (kind=3)
   Params
    (none)
   Body
    BlockStmt
     LBrace: "{" @3:12 (kind=41)
     Stmts
      ConstDecl
       Const: "const" @4:1 (kind=23)
       Name: "x" @4:7 (kind=3)
       Type:
        LBracket: "[" @4:9 (kind=43)
        RBracket: "]" @4:10 (kind=44)
        Ident: "int" @4:11 (kind=12)
       Eq: "=" @4:14 (kind=49)
        Type:
         LBracket: "[" @4:14 (kind=43)
         RBracket: "]" @4:16 (kind=44)
         Ident: "int" @4:17 (kind=12)
        LBrace: "{" @4:20 (kind=41)
         Elements
          IntLitExpr
           Value: "1" @4:21 (kind=4)
          IntLitExpr
           Value: "2" @4:23 (kind=4)
          IntLitExpr
           Value: "3" @4:24 (kind=4)
        RBrace: "}" @4:26 (kind=42)
      ConstDecl
       Const: "const" @4:30 (kind=23)
       Name: "y" @4:36 (kind=3)
       Type:
        LBracket: "[" @4:38 (kind=43)
        RBracket: "]" @4:39 (kind=44)
        Ident: "int" @4:40 (kind=12)
       Eq: "=" @4:43 (kind=49)
        Type:
         LBracket: "[" @4:44 (kind=43)
         RBracket: "]" @4:45 (kind=44)
         Ident: "int" @4:46 (kind=12)
        LBrace: "{" @4:49 (kind=41)
         Elements
          IntLitExpr
           Value: "1" @4:50 (kind=4)
          IntLitExpr
           Value: "2" @4:52 (kind=4)
          IntLitExpr
           Value: "3" @4:54 (kind=4)
        RBrace: "}" @4:55 (kind=42)
     RBrace: "}" @5:56 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("var_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 10},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 12},
			{Kind: token.KWVar, Value: "var", Line: 4, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 4, Column: 7},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 9},
			{Kind: token.RBracket, Value: "]", Line: 4, Column: 10},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 11},
			{Kind: token.RBrace, Value: "}", Line: 5, Column: 26},
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
   Name: "main" @3:6 (kind=3)
   Params
    (none)
   Body
    BlockStmt
     LBrace: "{" @3:12 (kind=41)
     Stmts
      VarDecl
       Var: "var" @4:1 (kind=11)
       Name: "x" @4:7 (kind=3)
       Type:
        LBracket: "[" @4:9 (kind=43)
        RBracket: "]" @4:10 (kind=44)
        Ident: "int" @4:11 (kind=12)
     RBrace: "}" @5:26 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("var_x2", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 10},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 12},
			{Kind: token.KWVar, Value: "var", Line: 4, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 4, Column: 7},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 9},
			{Kind: token.RBracket, Value: "]", Line: 4, Column: 10},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 11},
			{Kind: token.SemiComma, Value: ";", Line: 4, Column: 26},

			{Kind: token.KWVar, Value: "var", Line: 4, Column: 30},
			{Kind: token.Ident, Value: "y", Line: 4, Column: 36},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 38},
			{Kind: token.RBracket, Value: "]", Line: 4, Column: 39},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 40},
			{Kind: token.RBrace, Value: "}", Line: 5, Column: 56},
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
   Name: "main" @3:6 (kind=3)
   Params
    (none)
   Body
    BlockStmt
     LBrace: "{" @3:12 (kind=41)
     Stmts
      VarDecl
       Var: "var" @4:1 (kind=11)
       Name: "x" @4:7 (kind=3)
       Type:
        LBracket: "[" @4:9 (kind=43)
        RBracket: "]" @4:10 (kind=44)
        Ident: "int" @4:11 (kind=12)
      VarDecl
       Var: "var" @4:30 (kind=11)
       Name: "y" @4:36 (kind=3)
       Type:
        LBracket: "[" @4:38 (kind=43)
        RBracket: "]" @4:39 (kind=44)
        Ident: "int" @4:40 (kind=12)
     RBrace: "}" @5:56 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("var_x3", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 10},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 12},
			{Kind: token.KWVar, Value: "var", Line: 4, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 4, Column: 7},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 9},
			{Kind: token.RBracket, Value: "]", Line: 4, Column: 10},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 11},
			{Kind: token.Assign, Value: "=", Line: 4, Column: 14},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 14},
			{Kind: token.RBracket, Value: "]", Line: 4, Column: 16},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 17},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 20},
			{Kind: token.IntLit, Value: "1", Line: 4, Column: 21},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 22},
			{Kind: token.IntLit, Value: "2", Line: 4, Column: 23},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 24},
			{Kind: token.IntLit, Value: "3", Line: 4, Column: 24},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 26},
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
   Name: "main" @3:6 (kind=3)
   Params
    (none)
   Body
    BlockStmt
     LBrace: "{" @3:12 (kind=41)
     Stmts
      VarDecl
       Var: "var" @4:1 (kind=11)
       Name: "x" @4:7 (kind=3)
       Type:
        LBracket: "[" @4:9 (kind=43)
        RBracket: "]" @4:10 (kind=44)
        Ident: "int" @4:11 (kind=12)
       Eq: "=" @4:14 (kind=49)
        Type:
         LBracket: "[" @4:14 (kind=43)
         RBracket: "]" @4:16 (kind=44)
         Ident: "int" @4:17 (kind=12)
        LBrace: "{" @4:20 (kind=41)
         Elements
          IntLitExpr
           Value: "1" @4:21 (kind=4)
          IntLitExpr
           Value: "2" @4:23 (kind=4)
          IntLitExpr
           Value: "3" @4:24 (kind=4)
        RBrace: "}" @4:26 (kind=42)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("var_x4", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 10},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 12},
			{Kind: token.KWVar, Value: "var", Line: 4, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 4, Column: 7},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 9},
			{Kind: token.RBracket, Value: "]", Line: 4, Column: 10},
			{Kind: token.Ident, Value: "c", Line: 4, Column: 11},
			{Kind: token.Dot, Value: ".", Line: 4, Column: 11},
			{Kind: token.Ident, Value: "d", Line: 4, Column: 11},
			{Kind: token.RBrace, Value: "}", Line: 5, Column: 26},
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
   Name: "main" @3:6 (kind=3)
   Params
    (none)
   Body
    BlockStmt
     LBrace: "{" @3:12 (kind=41)
     Stmts
      VarDecl
       Var: "var" @4:1 (kind=11)
       Name: "x" @4:7 (kind=3)
       Type:
        LBracket: "[" @4:9 (kind=43)
        RBracket: "]" @4:10 (kind=44)
        Ident: "c" @4:11 (kind=3)
        Dot: "." @4:11 (kind=48)
        Ident: "d" @4:11 (kind=3)
     RBrace: "}" @5:26 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("var_x5", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 10},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 12},
			{Kind: token.KWVar, Value: "var", Line: 4, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 4, Column: 7},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 9},
			{Kind: token.RBracket, Value: "]", Line: 4, Column: 10},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 11},
			{Kind: token.Assign, Value: "=", Line: 4, Column: 14},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 15},
			{Kind: token.RBracket, Value: "]", Line: 4, Column: 16},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 17},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 20},
			{Kind: token.IntLit, Value: "1", Line: 4, Column: 21},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 22},
			{Kind: token.IntLit, Value: "2", Line: 4, Column: 23},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 24},
			{Kind: token.IntLit, Value: "3", Line: 4, Column: 25},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 26},

			{Kind: token.KWVar, Value: "var", Line: 5, Column: 3},
			{Kind: token.Ident, Value: "y", Line: 5, Column: 7},
			{Kind: token.KWView, Value: "view", Line: 5, Column: 9},
			{Kind: token.LBracket, Value: "[", Line: 5, Column: 14},
			{Kind: token.RBracket, Value: "]", Line: 5, Column: 15},
			{Kind: token.KWInt, Value: "int", Line: 5, Column: 16},
			{Kind: token.Assign, Value: "=", Line: 5, Column: 19},
			{Kind: token.Ident, Value: "x", Line: 5, Column: 20},
			{Kind: token.LBracket, Value: "[", Line: 5, Column: 21},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 22},
			{Kind: token.RBracket, Value: "]", Line: 5, Column: 23},

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
   Name: "main" @3:6 (kind=3)
   Params
    (none)
   Body
    BlockStmt
     LBrace: "{" @3:12 (kind=41)
     Stmts
      VarDecl
       Var: "var" @4:1 (kind=11)
       Name: "x" @4:7 (kind=3)
       Type:
        LBracket: "[" @4:9 (kind=43)
        RBracket: "]" @4:10 (kind=44)
        Ident: "int" @4:11 (kind=12)
       Eq: "=" @4:14 (kind=49)
        Type:
         LBracket: "[" @4:15 (kind=43)
         RBracket: "]" @4:16 (kind=44)
         Ident: "int" @4:17 (kind=12)
        LBrace: "{" @4:20 (kind=41)
         Elements
          IntLitExpr
           Value: "1" @4:21 (kind=4)
          IntLitExpr
           Value: "2" @4:23 (kind=4)
          IntLitExpr
           Value: "3" @4:25 (kind=4)
        RBrace: "}" @4:26 (kind=42)
      VarDecl
       Var: "var" @5:3 (kind=11)
       Name: "y" @5:7 (kind=3)
       View: "view" @5:9 (kind=76)
       Type:
        LBracket: "[" @5:14 (kind=43)
        RBracket: "]" @5:15 (kind=44)
        Ident: "int" @5:16 (kind=12)
       Eq: "=" @5:19 (kind=49)
        IdentExpr
         Name: "x" @5:20 (kind=3)
         LBracket: "[" @5:21 (kind=43)
         Colon: ":" @5:22 (kind=47)
         RBracket: "]" @5:23 (kind=44)
     RBrace: "}" @6:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("var_x6", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 10},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 12},
			{Kind: token.KWVar, Value: "var", Line: 4, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 4, Column: 7},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 9},
			{Kind: token.RBracket, Value: "]", Line: 4, Column: 10},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 11},
			{Kind: token.Assign, Value: "=", Line: 4, Column: 14},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 15},
			{Kind: token.RBracket, Value: "]", Line: 4, Column: 16},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 17},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 20},
			{Kind: token.IntLit, Value: "1", Line: 4, Column: 21},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 22},
			{Kind: token.IntLit, Value: "2", Line: 4, Column: 23},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 24},
			{Kind: token.IntLit, Value: "3", Line: 4, Column: 25},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 26},

			{Kind: token.KWVar, Value: "var", Line: 5, Column: 3},
			{Kind: token.Ident, Value: "y", Line: 5, Column: 7},
			{Kind: token.KWView, Value: "view", Line: 5, Column: 9},
			{Kind: token.LBracket, Value: "[", Line: 5, Column: 14},
			{Kind: token.RBracket, Value: "]", Line: 5, Column: 15},
			{Kind: token.KWInt, Value: "int", Line: 5, Column: 16},
			{Kind: token.Assign, Value: "=", Line: 5, Column: 19},
			{Kind: token.Ident, Value: "x", Line: 5, Column: 20},
			{Kind: token.LBracket, Value: "[", Line: 5, Column: 21},
			{Kind: token.IntLit, Value: "3", Line: 5, Column: 22},
			{Kind: token.Colon, Value: ":", Line: 5, Column: 23},
			{Kind: token.IntLit, Value: "6", Line: 5, Column: 24},
			{Kind: token.RBracket, Value: "]", Line: 5, Column: 25},
			{Kind: token.SemiComma, Value: ";", Line: 5, Column: 26},

			{Kind: token.KWVar, Value: "var", Line: 6, Column: 1},
			{Kind: token.Ident, Value: "z", Line: 6, Column: 7},
			{Kind: token.LBracket, Value: "[", Line: 6, Column: 9},
			{Kind: token.RBracket, Value: "]", Line: 6, Column: 10},
			{Kind: token.KWInt, Value: "int", Line: 6, Column: 11},
			{Kind: token.Assign, Value: "=", Line: 6, Column: 14},
			{Kind: token.LBracket, Value: "[", Line: 6, Column: 15},
			{Kind: token.RBracket, Value: "]", Line: 6, Column: 16},
			{Kind: token.KWInt, Value: "int", Line: 6, Column: 17},
			{Kind: token.LBrace, Value: "{", Line: 6, Column: 20},
			{Kind: token.IntLit, Value: "1", Line: 6, Column: 21},
			{Kind: token.Comma, Value: ",", Line: 6, Column: 22},
			{Kind: token.IntLit, Value: "2", Line: 6, Column: 23},
			{Kind: token.Comma, Value: ",", Line: 6, Column: 24},
			{Kind: token.IntLit, Value: "3", Line: 6, Column: 25},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 26},

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
   Name: "main" @3:6 (kind=3)
   Params
    (none)
   Body
    BlockStmt
     LBrace: "{" @3:12 (kind=41)
     Stmts
      VarDecl
       Var: "var" @4:1 (kind=11)
       Name: "x" @4:7 (kind=3)
       Type:
        LBracket: "[" @4:9 (kind=43)
        RBracket: "]" @4:10 (kind=44)
        Ident: "int" @4:11 (kind=12)
       Eq: "=" @4:14 (kind=49)
        Type:
         LBracket: "[" @4:15 (kind=43)
         RBracket: "]" @4:16 (kind=44)
         Ident: "int" @4:17 (kind=12)
        LBrace: "{" @4:20 (kind=41)
         Elements
          IntLitExpr
           Value: "1" @4:21 (kind=4)
          IntLitExpr
           Value: "2" @4:23 (kind=4)
          IntLitExpr
           Value: "3" @4:25 (kind=4)
        RBrace: "}" @4:26 (kind=42)
      VarDecl
       Var: "var" @5:3 (kind=11)
       Name: "y" @5:7 (kind=3)
       View: "view" @5:9 (kind=76)
       Type:
        LBracket: "[" @5:14 (kind=43)
        RBracket: "]" @5:15 (kind=44)
        Ident: "int" @5:16 (kind=12)
       Eq: "=" @5:19 (kind=49)
        IdentExpr
         Name: "x" @5:20 (kind=3)
         LBracket: "[" @5:21 (kind=43)
        IntLitExpr
         Value: "3" @5:22 (kind=4)
         Colon: ":" @5:23 (kind=47)
        IntLitExpr
         Value: "6" @5:24 (kind=4)
         RBracket: "]" @5:25 (kind=44)
      VarDecl
       Var: "var" @6:1 (kind=11)
       Name: "z" @6:7 (kind=3)
       Type:
        LBracket: "[" @6:9 (kind=43)
        RBracket: "]" @6:10 (kind=44)
        Ident: "int" @6:11 (kind=12)
       Eq: "=" @6:14 (kind=49)
        Type:
         LBracket: "[" @6:15 (kind=43)
         RBracket: "]" @6:16 (kind=44)
         Ident: "int" @6:17 (kind=12)
        LBrace: "{" @6:20 (kind=41)
         Elements
          IntLitExpr
           Value: "1" @6:21 (kind=4)
          IntLitExpr
           Value: "2" @6:23 (kind=4)
          IntLitExpr
           Value: "3" @6:25 (kind=4)
        RBrace: "}" @6:26 (kind=42)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("var_x7", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 10},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 12},

			{Kind: token.KWVar, Value: "var", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "y", Line: 4, Column: 7},
			{Kind: token.KWView, Value: "view", Line: 4, Column: 9},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 14},
			{Kind: token.RBracket, Value: "]", Line: 4, Column: 15},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 16},
			{Kind: token.Assign, Value: "=", Line: 4, Column: 19},
			{Kind: token.Ident, Value: "x", Line: 4, Column: 20},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 21},
			{Kind: token.Colon, Value: ":", Line: 4, Column: 23},
			{Kind: token.IntLit, Value: "6", Line: 4, Column: 24},
			{Kind: token.RBracket, Value: "]", Line: 4, Column: 25},

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
   Name: "main" @3:6 (kind=3)
   Params
    (none)
   Body
    BlockStmt
     LBrace: "{" @3:12 (kind=41)
     Stmts
      VarDecl
       Var: "var" @4:3 (kind=11)
       Name: "y" @4:7 (kind=3)
       View: "view" @4:9 (kind=76)
       Type:
        LBracket: "[" @4:14 (kind=43)
        RBracket: "]" @4:15 (kind=44)
        Ident: "int" @4:16 (kind=12)
       Eq: "=" @4:19 (kind=49)
        IdentExpr
         Name: "x" @4:20 (kind=3)
         LBracket: "[" @4:21 (kind=43)
         Colon: ":" @4:23 (kind=47)
        IntLitExpr
         Value: "6" @4:24 (kind=4)
         RBracket: "]" @4:25 (kind=44)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("var_x8", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 10},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 12},

			{Kind: token.KWVar, Value: "var", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "y", Line: 4, Column: 7},
			{Kind: token.KWView, Value: "view", Line: 4, Column: 9},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 14},
			{Kind: token.RBracket, Value: "]", Line: 4, Column: 15},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 16},
			{Kind: token.Assign, Value: "=", Line: 4, Column: 19},
			{Kind: token.Ident, Value: "x", Line: 4, Column: 20},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 21},
			{Kind: token.IntLit, Value: "6", Line: 4, Column: 24},
			{Kind: token.RBracket, Value: "]", Line: 4, Column: 25},

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
   Name: "main" @3:6 (kind=3)
   Params
    (none)
   Body
    BlockStmt
     LBrace: "{" @3:12 (kind=41)
     Stmts
      VarDecl
       Var: "var" @4:3 (kind=11)
       Name: "y" @4:7 (kind=3)
       View: "view" @4:9 (kind=76)
       Type:
        LBracket: "[" @4:14 (kind=43)
        RBracket: "]" @4:15 (kind=44)
        Ident: "int" @4:16 (kind=12)
       Eq: "=" @4:19 (kind=49)
        IdentExpr
         Name: "x" @4:20 (kind=3)
         LBracket: "[" @4:21 (kind=43)
        IntLitExpr
         Value: "6" @4:24 (kind=4)
         Colon: "" @0:0 (kind=0)
         RBracket: "]" @4:25 (kind=44)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("bad_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWConst, Value: "const", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 7},
			{Kind: token.LBracket, Value: "[", Line: 3, Column: 9},
			{Kind: token.RBracket, Value: "]", Line: 3, Column: 10},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 11},
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

			{Kind: token.KWVar, Value: "var", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 7},
			{Kind: token.LBracket, Value: "[", Line: 3, Column: 9},
			{Kind: token.RBracket, Value: "]", Line: 3, Column: 10},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 11},
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

			{Kind: token.KWConst, Value: "const", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 7},
			{Kind: token.LBracket, Value: "[", Line: 3, Column: 9},
			{Kind: token.RBracket, Value: "]", Line: 3, Column: 10},
			{Kind: token.KWStruct, Value: "struct", Line: 3, Column: 11},
			{Kind: token.Assign, Value: "=", Line: 3, Column: 14},
			{Kind: token.LBracket, Value: "[", Line: 3, Column: 14},
			{Kind: token.RBracket, Value: "]", Line: 3, Column: 16},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 17},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 20},
			{Kind: token.IntLit, Value: "1", Line: 3, Column: 21},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 22},
			{Kind: token.IntLit, Value: "2", Line: 3, Column: 23},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 24},
			{Kind: token.IntLit, Value: "3", Line: 3, Column: 24},
			{Kind: token.RBrace, Value: "}", Line: 3, Column: 26},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}
