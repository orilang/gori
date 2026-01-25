package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
	"github.com/stretchr/testify/assert"
)

func TestParser_for_stmt(t *testing.T) {
	assert := assert.New(t)

	t.Run("function_infinite_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 7},
			{Kind: token.KWReturn, Value: "return", Line: 5, Column: 4},
			{Kind: token.BoolLit, Value: "false", Line: 5, Column: 11},
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
     LBrace: "{" @3:10 (kind=41)
     Stmts
      ForStmt
       For: "for" @4:3 (kind=31)
        BlockStmt
         LBrace: "{" @4:7 (kind=41)
         Stmts
          ReturnStmt
           Values
            BoolLitExpr
             Value: "false" @5:11 (kind=7)
         RBrace: "}" @6:3 (kind=42)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_infinite_x2", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 7},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 8},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
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
      ForStmt
       For: "for" @4:3 (kind=31)
     RBrace: "}" @4:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_condition_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 7},
			{Kind: token.Lt, Value: "<", Line: 4, Column: 8},
			{Kind: token.Ident, Value: "b", Line: 4, Column: 9},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 11},
			{Kind: token.KWReturn, Value: "return", Line: 5, Column: 4},
			{Kind: token.BoolLit, Value: "false", Line: 5, Column: 11},
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
     LBrace: "{" @3:10 (kind=41)
     Stmts
      ForStmt
       For: "for" @4:3 (kind=31)
       Condition
        BinaryExpr
         IdentExpr
          Name: "a" @4:7 (kind=3)
         Operator: "<" @4:8 (kind=64)
         IdentExpr
          Name: "b" @4:9 (kind=3)
        BlockStmt
         LBrace: "{" @4:11 (kind=41)
         Stmts
          ReturnStmt
           Values
            BoolLitExpr
             Value: "false" @5:11 (kind=7)
         RBrace: "}" @6:3 (kind=42)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_condition_x2", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 7},
			{Kind: token.Lt, Value: "<", Line: 4, Column: 8},
			{Kind: token.Ident, Value: "b", Line: 4, Column: 9},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 10},
			{Kind: token.RBrace, Value: "}", Line: 5, Column: 11},
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
     LBrace: "{" @3:10 (kind=41)
     Stmts
      ForStmt
       For: "for" @4:3 (kind=31)
       Condition
        BinaryExpr
         IdentExpr
          Name: "a" @4:7 (kind=3)
         Operator: "<" @4:8 (kind=64)
         IdentExpr
          Name: "b" @4:9 (kind=3)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_all_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 7},
			{Kind: token.Define, Value: ":=", Line: 4, Column: 8},
			{Kind: token.IntLit, Value: "0", Line: 4, Column: 10},
			{Kind: token.SemiComma, Value: ";", Line: 4, Column: 11},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 12},
			{Kind: token.Lt, Value: "<", Line: 4, Column: 13},
			{Kind: token.IntLit, Value: "5", Line: 4, Column: 14},
			{Kind: token.SemiComma, Value: ";", Line: 4, Column: 15},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 16},
			{Kind: token.PlusEq, Value: "+=", Line: 4, Column: 17},
			{Kind: token.IntLit, Value: "1", Line: 4, Column: 18},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 20},
			{Kind: token.KWReturn, Value: "return", Line: 5, Column: 4},
			{Kind: token.Ident, Value: "a", Line: 5, Column: 11},
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
     LBrace: "{" @3:10 (kind=41)
     Stmts
      ForStmt
       For: "for" @4:3 (kind=31)
       Init
        AssignStmt
         Left
          IdentExpr
           Name: "a" @4:7 (kind=3)
         Operator: ":=" @4:8 (kind=50)
         Right
          IntLitExpr
           Value: "0" @4:10 (kind=4)
       Condition
        BinaryExpr
         IdentExpr
          Name: "a" @4:12 (kind=3)
         Operator: "<" @4:13 (kind=64)
         IntLitExpr
          Value: "5" @4:14 (kind=4)
       Post
        AssignStmt
         Left
          IdentExpr
           Name: "a" @4:16 (kind=3)
         Operator: "+=" @4:17 (kind=52)
         Right
          IntLitExpr
           Value: "1" @4:18 (kind=4)
        BlockStmt
         LBrace: "{" @4:20 (kind=41)
         Stmts
          ReturnStmt
           Values
            IdentExpr
             Name: "a" @5:11 (kind=3)
         RBrace: "}" @6:3 (kind=42)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_all_x2", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 7},
			{Kind: token.Define, Value: ":=", Line: 4, Column: 8},
			{Kind: token.IntLit, Value: "0", Line: 4, Column: 10},
			{Kind: token.SemiComma, Value: ";", Line: 4, Column: 11},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 12},
			{Kind: token.Lt, Value: "<", Line: 4, Column: 13},
			{Kind: token.IntLit, Value: "5", Line: 4, Column: 14},
			{Kind: token.SemiComma, Value: ";", Line: 4, Column: 15},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 16},
			{Kind: token.PlusEq, Value: "+=", Line: 4, Column: 17},
			{Kind: token.IntLit, Value: "1", Line: 4, Column: 18},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 20},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 21},
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
     LBrace: "{" @3:10 (kind=41)
     Stmts
      ForStmt
       For: "for" @4:3 (kind=31)
       Init
        AssignStmt
         Left
          IdentExpr
           Name: "a" @4:7 (kind=3)
         Operator: ":=" @4:8 (kind=50)
         Right
          IntLitExpr
           Value: "0" @4:10 (kind=4)
       Condition
        BinaryExpr
         IdentExpr
          Name: "a" @4:12 (kind=3)
         Operator: "<" @4:13 (kind=64)
         IntLitExpr
          Value: "5" @4:14 (kind=4)
       Post
        AssignStmt
         Left
          IdentExpr
           Name: "a" @4:16 (kind=3)
         Operator: "+=" @4:17 (kind=52)
         Right
          IntLitExpr
           Value: "1" @4:18 (kind=4)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_all_x3", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 7},
			{Kind: token.Define, Value: ":=", Line: 4, Column: 8},
			{Kind: token.IntLit, Value: "0", Line: 4, Column: 10},
			{Kind: token.SemiComma, Value: ";", Line: 4, Column: 11},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 12},
			{Kind: token.Lt, Value: "<", Line: 4, Column: 13},
			{Kind: token.IntLit, Value: "5", Line: 4, Column: 14},
			{Kind: token.SemiComma, Value: ";", Line: 4, Column: 15},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 16},
			{Kind: token.PPlus, Value: "++", Line: 4, Column: 17},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 20},
			{Kind: token.KWReturn, Value: "return", Line: 5, Column: 4},
			{Kind: token.Ident, Value: "a", Line: 5, Column: 11},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
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
      ForStmt
       For: "for" @4:3 (kind=31)
       Init
        AssignStmt
         Left
          IdentExpr
           Name: "a" @4:7 (kind=3)
         Operator: ":=" @4:8 (kind=50)
         Right
          IntLitExpr
           Value: "0" @4:10 (kind=4)
       Condition
        BinaryExpr
         IdentExpr
          Name: "a" @4:12 (kind=3)
         Operator: "<" @4:13 (kind=64)
         IntLitExpr
          Value: "5" @4:14 (kind=4)
       Post
        IncDecStmt
         X:
          IdentExpr
           Name: "a" @4:16 (kind=3)
         Operator: "++" @4:17 (kind=53)
        BlockStmt
         LBrace: "{" @4:20 (kind=41)
         Stmts
          ReturnStmt
           Values
            IdentExpr
             Name: "a" @5:11 (kind=3)
         RBrace: "}" @6:3 (kind=42)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_range_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "i", Line: 4, Column: 7},
			{Kind: token.Define, Value: ":=", Line: 4, Column: 8},
			{Kind: token.KWRange, Value: "range", Line: 4, Column: 11},
			{Kind: token.IntLit, Value: "5", Line: 4, Column: 17},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 19},
			{Kind: token.Ident, Value: "z", Line: 5, Column: 4},
			{Kind: token.Assign, Value: "=", Line: 5, Column: 6},
			{Kind: token.Ident, Value: "i", Line: 5, Column: 8},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
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
      RangeStmt
       For: "for" @4:3 (kind=31)
       Key
        IdentExpr
         Name: "i" @4:7 (kind=3)
       Op: ":=" @4:8 (kind=50)
       Range: "range" @4:11 (kind=71)
        IntLitExpr
         Value: "5" @4:17 (kind=4)
        BlockStmt
         LBrace: "{" @4:19 (kind=41)
         Stmts
          AssignStmt
           Left
            IdentExpr
             Name: "z" @5:4 (kind=3)
           Operator: "=" @5:6 (kind=49)
           Right
            IdentExpr
             Name: "i" @5:8 (kind=3)
         RBrace: "}" @6:3 (kind=42)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_range_x2", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "k", Line: 4, Column: 7},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 8},
			{Kind: token.Ident, Value: "v", Line: 4, Column: 9},
			{Kind: token.Define, Value: ":=", Line: 4, Column: 10},
			{Kind: token.KWRange, Value: "range", Line: 4, Column: 13},
			{Kind: token.IntLit, Value: "5", Line: 4, Column: 19},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 21},
			{Kind: token.Ident, Value: "z", Line: 5, Column: 4},
			{Kind: token.Assign, Value: "=", Line: 5, Column: 6},
			{Kind: token.Ident, Value: "i", Line: 5, Column: 8},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
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
      RangeStmt
       For: "for" @4:3 (kind=31)
       Key
        IdentExpr
         Name: "k" @4:7 (kind=3)
       Condition
        IdentExpr
         Name: "v" @4:9 (kind=3)
       Op: ":=" @4:10 (kind=50)
       Range: "range" @4:13 (kind=71)
        IntLitExpr
         Value: "5" @4:19 (kind=4)
        BlockStmt
         LBrace: "{" @4:21 (kind=41)
         Stmts
          AssignStmt
           Left
            IdentExpr
             Name: "z" @5:4 (kind=3)
           Operator: "=" @5:6 (kind=49)
           Right
            IdentExpr
             Name: "i" @5:8 (kind=3)
         RBrace: "}" @6:3 (kind=42)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_range_x3", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "i", Line: 4, Column: 7},
			{Kind: token.Define, Value: ":=", Line: 4, Column: 8},
			{Kind: token.KWRange, Value: "range", Line: 4, Column: 11},
			{Kind: token.IntLit, Value: "5", Line: 4, Column: 17},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 19},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 20},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
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
      RangeStmt
       For: "for" @4:3 (kind=31)
       Key
        IdentExpr
         Name: "i" @4:7 (kind=3)
       Op: ":=" @4:8 (kind=50)
       Range: "range" @4:11 (kind=71)
        IntLitExpr
         Value: "5" @4:17 (kind=4)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_range_x4", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "k", Line: 4, Column: 7},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 8},
			{Kind: token.Ident, Value: "v", Line: 4, Column: 9},
			{Kind: token.Define, Value: ":=", Line: 4, Column: 10},
			{Kind: token.KWRange, Value: "range", Line: 4, Column: 11},
			{Kind: token.IntLit, Value: "5", Line: 4, Column: 17},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 18},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 20},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
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
      RangeStmt
       For: "for" @4:3 (kind=31)
       Key
        IdentExpr
         Name: "k" @4:7 (kind=3)
       Condition
        IdentExpr
         Name: "v" @4:9 (kind=3)
       Op: ":=" @4:10 (kind=50)
       Range: "range" @4:11 (kind=71)
        IntLitExpr
         Value: "5" @4:17 (kind=4)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_range_x5", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.KWRange, Value: "range", Line: 4, Column: 7},
			{Kind: token.IntLit, Value: "5", Line: 4, Column: 13},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 15},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 16},
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
     LBrace: "{" @3:10 (kind=41)
     Stmts
      RangeStmt
       For: "for" @4:3 (kind=31)
       Op: "" @0:0 (kind=0)
       Range: "range" @4:7 (kind=71)
        IntLitExpr
         Value: "5" @4:13 (kind=4)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_range_x6", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.KWRange, Value: "range", Line: 4, Column: 7},
			{Kind: token.IntLit, Value: "5", Line: 4, Column: 13},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 15},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 16},
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
     LBrace: "{" @3:10 (kind=41)
     Stmts
      RangeStmt
       For: "for" @4:3 (kind=31)
       Op: "" @0:0 (kind=0)
       Range: "range" @4:7 (kind=71)
        IntLitExpr
         Value: "5" @4:13 (kind=4)
     RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_range_x7", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.KWRange, Value: "range", Line: 4, Column: 7},
			{Kind: token.Ident, Value: "ch", Line: 4, Column: 13},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 16},
			{Kind: token.KWReturn, Value: "return", Line: 5, Column: 4},
			{Kind: token.Ident, Value: "x", Line: 5, Column: 11},
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
     LBrace: "{" @3:10 (kind=41)
     Stmts
      RangeStmt
       For: "for" @4:3 (kind=31)
       Op: "" @0:0 (kind=0)
       Range: "range" @4:7 (kind=71)
        IdentExpr
         Name: "ch" @4:13 (kind=3)
        BlockStmt
         LBrace: "{" @4:16 (kind=41)
         Stmts
          ReturnStmt
           Values
            IdentExpr
             Name: "x" @5:11 (kind=3)
         RBrace: "}" @6:3 (kind=42)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_range_x8", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "k", Line: 4, Column: 7},
			{Kind: token.Assign, Value: "=", Line: 4, Column: 8},
			{Kind: token.KWRange, Value: "range", Line: 4, Column: 11},
			{Kind: token.IntLit, Value: "5", Line: 4, Column: 17},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 18},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 20},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
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
      RangeStmt
       For: "for" @4:3 (kind=31)
       Key
        IdentExpr
         Name: "k" @4:7 (kind=3)
       Op: "=" @4:8 (kind=49)
       Range: "range" @4:11 (kind=71)
        IntLitExpr
         Value: "5" @4:17 (kind=4)
     RBrace: "}" @7:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("bad_infinite_x1", func(t *testing.T) {
		// infinity loop
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 7},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 8},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_condition_x1", func(t *testing.T) {
		// condition
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 7},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 8},
			{Kind: token.Lt, Value: "<", Line: 4, Column: 9},
			{Kind: token.Ident, Value: "b", Line: 4, Column: 10},
			{Kind: token.SemiComma, Value: ";", Line: 4, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 12},
			{Kind: token.KWReturn, Value: "return", Line: 5, Column: 4},
			{Kind: token.BoolLit, Value: "false", Line: 5, Column: 11},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_condition_x2", func(t *testing.T) {
		// condition
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 7},
			{Kind: token.Lt, Value: "<", Line: 4, Column: 8},
			{Kind: token.Ident, Value: "b", Line: 4, Column: 9},
			{Kind: token.SemiComma, Value: ";", Line: 4, Column: 10},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 11},
			{Kind: token.KWReturn, Value: "return", Line: 5, Column: 4},
			{Kind: token.BoolLit, Value: "false", Line: 5, Column: 11},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_condition_x3", func(t *testing.T) {
		// condition
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 7},
			{Kind: token.Lt, Value: "<", Line: 4, Column: 8},
			{Kind: token.Ident, Value: "b", Line: 4, Column: 9},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 10},
			{Kind: token.KWReturn, Value: "return", Line: 5, Column: 4},
			{Kind: token.BoolLit, Value: "false", Line: 5, Column: 11},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_range_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "i", Line: 4, Column: 7},
			{Kind: token.Define, Value: ":=", Line: 4, Column: 8},
			{Kind: token.KWRange, Value: "range", Line: 4, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 17},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 18},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_range_x3", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "i", Line: 4, Column: 7},
			{Kind: token.Define, Value: ":=", Line: 4, Column: 8},
			{Kind: token.KWRange, Value: "range", Line: 4, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 17},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 19},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_range_x4", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 7},
			{Kind: token.Ident, Value: "i", Line: 4, Column: 8},
			{Kind: token.Define, Value: ":=", Line: 4, Column: 9},
			{Kind: token.KWRange, Value: "range", Line: 4, Column: 10},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 17},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 18},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_range_x5", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.BoolLit, Value: "true", Line: 4, Column: 7},
			{Kind: token.Define, Value: ":=", Line: 4, Column: 8},
			{Kind: token.KWRange, Value: "range", Line: 4, Column: 10},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 17},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 18},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_range_x6", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "i", Line: 4, Column: 7},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 8},
			{Kind: token.BoolLit, Value: "true", Line: 4, Column: 9},
			{Kind: token.Define, Value: ":=", Line: 4, Column: 10},
			{Kind: token.KWRange, Value: "range", Line: 4, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 17},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 18},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_range_x7", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "k", Line: 4, Column: 7},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 8},
			{Kind: token.Ident, Value: "v", Line: 4, Column: 9},
			{Kind: token.Define, Value: ":=", Line: 4, Column: 10},
			{Kind: token.KWRange, Value: "range", Line: 4, Column: 11},
			{Kind: token.IntLit, Value: "5", Line: 4, Column: 17},
			{Kind: token.IntLit, Value: "5", Line: 4, Column: 18},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 20},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_range_x8", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "k", Line: 4, Column: 7},
			{Kind: token.PlusEq, Value: "+=", Line: 4, Column: 8},
			{Kind: token.KWRange, Value: "range", Line: 4, Column: 11},
			{Kind: token.IntLit, Value: "5", Line: 4, Column: 17},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 18},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 20},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_range_x9", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.KWRange, Value: "range", Line: 4, Column: 7},
			{Kind: token.Ident, Value: "ch", Line: 4, Column: 13},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 16},
			{Kind: token.KWReturn, Value: "return", Line: 5, Column: 4},
			{Kind: token.Ident, Value: "x", Line: 5, Column: 11},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 8, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_range_x10", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.KWRange, Value: "range", Line: 4, Column: 7},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 13},
			{Kind: token.KWReturn, Value: "return", Line: 5, Column: 4},
			{Kind: token.Ident, Value: "x", Line: 5, Column: 11},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 8, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_range_x11", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "k", Line: 4, Column: 7},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 8},
			{Kind: token.Ident, Value: "v", Line: 4, Column: 9},
			{Kind: token.PlusEq, Value: "+=", Line: 4, Column: 10},
			{Kind: token.KWRange, Value: "range", Line: 4, Column: 11},
			{Kind: token.IntLit, Value: "5", Line: 4, Column: 17},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 18},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_range_x12", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "k", Line: 4, Column: 7},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 8},
			{Kind: token.Ident, Value: "v", Line: 4, Column: 9},
			{Kind: token.PlusEq, Value: "+=", Line: 4, Column: 10},
			{Kind: token.KWRange, Value: "range", Line: 4, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 17},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 18},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_range_x13", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "k", Line: 4, Column: 7},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 8},
			{Kind: token.Ident, Value: "v", Line: 4, Column: 9},
			{Kind: token.Assign, Value: "=", Line: 4, Column: 10},
			{Kind: token.KWRange, Value: "range", Line: 4, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 17},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 18},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_range_x14", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "i", Line: 4, Column: 7},
			{Kind: token.Define, Value: ":=", Line: 4, Column: 8},
			{Kind: token.KWRange, Value: "range", Line: 4, Column: 11},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 17},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 18},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_all_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 7},
			{Kind: token.Define, Value: ":=", Line: 4, Column: 8},
			{Kind: token.IntLit, Value: "0", Line: 4, Column: 10},
			{Kind: token.SemiComma, Value: ";", Line: 4, Column: 11},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 12},
			{Kind: token.Lt, Value: "<", Line: 4, Column: 13},
			{Kind: token.IntLit, Value: "5", Line: 4, Column: 14},
			{Kind: token.SemiComma, Value: ";", Line: 4, Column: 15},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 16},
			{Kind: token.PlusEq, Value: "+=", Line: 4, Column: 17},
			{Kind: token.IntLit, Value: "1", Line: 4, Column: 18},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 20},
			{Kind: token.KWReturn, Value: "return", Line: 5, Column: 4},
			{Kind: token.Ident, Value: "a", Line: 5, Column: 11},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 8, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_all_x2", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 7},
			{Kind: token.Define, Value: ":=", Line: 4, Column: 8},
			{Kind: token.IntLit, Value: "0", Line: 4, Column: 10},
			{Kind: token.SemiComma, Value: ";", Line: 4, Column: 11},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 12},
			{Kind: token.Lt, Value: "<", Line: 4, Column: 13},
			{Kind: token.IntLit, Value: "5", Line: 4, Column: 14},
			{Kind: token.SemiComma, Value: ";", Line: 4, Column: 15},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 16},
			{Kind: token.KWReturn, Value: "return", Line: 5, Column: 4},
			{Kind: token.Ident, Value: "a", Line: 5, Column: 11},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 8, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_all_x3", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 7},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 7},
			{Kind: token.PPlus, Value: "++", Line: 4, Column: 8},
			{Kind: token.SemiComma, Value: ";", Line: 4, Column: 11},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 12},
			{Kind: token.Lt, Value: "<", Line: 4, Column: 13},
			{Kind: token.IntLit, Value: "5", Line: 4, Column: 14},
			{Kind: token.SemiComma, Value: ";", Line: 4, Column: 15},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 16},
			{Kind: token.KWReturn, Value: "return", Line: 5, Column: 4},
			{Kind: token.Ident, Value: "a", Line: 5, Column: 11},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 8, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_all_x4", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 7},
			{Kind: token.Define, Value: ":=", Line: 4, Column: 8},
			{Kind: token.IntLit, Value: "0", Line: 4, Column: 9},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 10},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 12},
			{Kind: token.Lt, Value: "<", Line: 4, Column: 13},
			{Kind: token.IntLit, Value: "5", Line: 4, Column: 14},
			{Kind: token.SemiComma, Value: ";", Line: 4, Column: 15},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 16},
			{Kind: token.KWReturn, Value: "return", Line: 5, Column: 4},
			{Kind: token.Ident, Value: "a", Line: 5, Column: 11},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 8, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_all_x5", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 7},
			{Kind: token.Define, Value: ":=", Line: 4, Column: 8},
			{Kind: token.IntLit, Value: "0", Line: 4, Column: 7},
			{Kind: token.SemiComma, Value: ";", Line: 4, Column: 11},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 12},
			{Kind: token.Lt, Value: "<", Line: 4, Column: 13},
			{Kind: token.IntLit, Value: "5", Line: 4, Column: 14},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 15},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 16},
			{Kind: token.KWReturn, Value: "return", Line: 5, Column: 4},
			{Kind: token.Ident, Value: "a", Line: 5, Column: 11},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 8, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_all_x6", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 7},
			{Kind: token.Define, Value: ":=", Line: 4, Column: 8},
			{Kind: token.IntLit, Value: "0", Line: 4, Column: 9},
			{Kind: token.SemiComma, Value: ";", Line: 4, Column: 11},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 12},
			{Kind: token.Lt, Value: "<", Line: 4, Column: 13},
			{Kind: token.IntLit, Value: "5", Line: 4, Column: 14},
			{Kind: token.SemiComma, Value: ";", Line: 4, Column: 15},
			{Kind: token.IntLit, Value: "5", Line: 4, Column: 16},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 17},
			{Kind: token.KWReturn, Value: "return", Line: 5, Column: 4},
			{Kind: token.Ident, Value: "a", Line: 5, Column: 11},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 8, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_all_x7", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 10},
			{Kind: token.KWFor, Value: "for", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 7},
			{Kind: token.Define, Value: ":=", Line: 4, Column: 8},
			{Kind: token.IntLit, Value: "0", Line: 4, Column: 9},
			{Kind: token.SemiComma, Value: ";", Line: 4, Column: 11},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 12},
			{Kind: token.Lt, Value: "<", Line: 4, Column: 13},
			{Kind: token.IntLit, Value: "5", Line: 4, Column: 14},
			{Kind: token.SemiComma, Value: ";", Line: 4, Column: 15},
			{Kind: token.KWVar, Value: "var", Line: 4, Column: 16},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 17},
			{Kind: token.KWReturn, Value: "return", Line: 5, Column: 4},
			{Kind: token.Ident, Value: "a", Line: 5, Column: 11},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 8, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}
