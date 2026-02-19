package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
	"github.com/stretchr/testify/assert"
)

func TestParser_func_decl(t *testing.T) {
	assert := assert.New(t)

	t.Run("return_types_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 10},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 14},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 15},
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
   Results
     Param
      Type
       NameType
        Name: "int" @3:10 (kind=12)
   Body
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("return_types_x2", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 9},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 10},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 14},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 15},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 19},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 20},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 15},
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
   Results
    LParent: "(" @3:9 (kind=39)
     Param
      Type
       NameType
        Name: "int" @3:10 (kind=12)
     Param
      Type
       NameType
        Name: "int" @3:15 (kind=12)
    RParent: ")" @3:19 (kind=40)
   Body
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("return_types_x3", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 9},
			{Kind: token.Ident, Value: "a", Line: 3, Column: 10},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 12},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 16},
			{Kind: token.Ident, Value: "b", Line: 3, Column: 17},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 19},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 20},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 21},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 15},
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
   Results
    LParent: "(" @3:9 (kind=39)
     Param
      Ident: "a" @3:10 (kind=3)
      Type
       NameType
        Name: "int" @3:12 (kind=12)
     Param
      Ident: "b" @3:17 (kind=3)
      Type
       NameType
        Name: "int" @3:19 (kind=12)
    RParent: ")" @3:20 (kind=40)
   Body
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("return_types_x4", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 11},
			{Kind: token.RBrace, Value: "}", Line: 3, Column: 12},
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
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("return_types_x5", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 9},
			{Kind: token.Ident, Value: "a", Line: 3, Column: 10},
			{Kind: token.Ident, Value: "z", Line: 3, Column: 12},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 16},
			{Kind: token.Ident, Value: "b", Line: 3, Column: 17},
			{Kind: token.Ident, Value: "z", Line: 3, Column: 19},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 20},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 21},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 15},
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
   Results
    LParent: "(" @3:9 (kind=39)
     Param
      Ident: "a" @3:10 (kind=3)
      Type
       NameType
        Name: "z" @3:12 (kind=3)
     Param
      Ident: "b" @3:17 (kind=3)
      Type
       NameType
        Name: "z" @3:19 (kind=3)
    RParent: ")" @3:20 (kind=40)
   Body
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("return_types_x6", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 9},
			{Kind: token.Ident, Value: "a", Line: 3, Column: 10},
			{Kind: token.Ident, Value: "z", Line: 3, Column: 12},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 13},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 14},
			{Kind: token.RBrace, Value: "}", Line: 3, Column: 15},
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
   Results
    LParent: "(" @3:9 (kind=39)
     Param
      Ident: "a" @3:10 (kind=3)
      Type
       NameType
        Name: "z" @3:12 (kind=3)
    RParent: ")" @3:13 (kind=40)
   Body
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("struct_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWType, Value: "type", Line: 4, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 4, Column: 6},
			{Kind: token.KWStruct, Value: "struct", Line: 4, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 18},
			{Kind: token.Ident, Value: "x", Line: 4, Column: 20},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 22},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 19},
			{Kind: token.RBrace, Value: "}", Line: 5, Column: 9},
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
     LBrace: "{" @3:9 (kind=41)
     Stmts
      Type: "type" @4:1 (kind=26)
      Name: "test" @4:6 (kind=3)
      Struct: "struct" @4:11 (kind=27)
      LBrace: "{" @4:18 (kind=41)
       Name: "x" @4:20 (kind=3)
       Type:
        NameType
         Name: "int" @4:22 (kind=12)
      RBrace: "}" @4:19 (kind=42)
     RBrace: "}" @5:9 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("interface_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWType, Value: "type", Line: 4, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 4, Column: 6},
			{Kind: token.KWInterface, Value: "interface", Line: 4, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 22},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 23},
			{Kind: token.RBrace, Value: "}", Line: 5, Column: 9},
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
     LBrace: "{" @3:9 (kind=41)
     Stmts
      Type: "type" @4:1 (kind=26)
      Name: "test" @4:6 (kind=3)
      Interface: "interface" @4:11 (kind=28)
      LBrace: "{" @4:22 (kind=41)
      RBrace: "}" @4:23 (kind=42)
     RBrace: "}" @5:9 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("enum_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 9},
			{Kind: token.KWType, Value: "type", Line: 4, Column: 1},
			{Kind: token.Ident, Value: "Color", Line: 4, Column: 6},
			{Kind: token.KWEnum, Value: "enum", Line: 4, Column: 11},
			{Kind: token.Assign, Value: "=", Line: 4, Column: 17},
			{Kind: token.Pipe, Value: "|", Line: 5, Column: 3},
			{Kind: token.Ident, Value: "Red", Line: 5, Column: 5},
			{Kind: token.RBrace, Value: "}", Line: 5, Column: 9},
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
     LBrace: "{" @3:9 (kind=41)
     Stmts
      Type: "type" @4:1 (kind=26)
       Name: "Color" @4:6 (kind=3)
       Public: true
       Enum: "enum" @4:11 (kind=74)
       Eq: "=" @4:17 (kind=49)
       Variants
        Ident: "Red" @5:5 (kind=3)
     RBrace: "}" @5:9 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("bad_return_types_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 9},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 10},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 14},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 15},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_return_types_x2", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 9},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 10},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 14},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 15},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 16},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 20},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 21},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 15},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_return_types_x3", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 9},
			{Kind: token.Ident, Value: "a", Line: 3, Column: 10},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 12},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 16},
			{Kind: token.Ident, Value: "b", Line: 3, Column: 17},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 20},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 21},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 15},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_return_types_x4", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 9},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 10},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 11},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 12},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_return_types_x5", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 9},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 10},
			{Kind: token.Ident, Value: "a", Line: 3, Column: 11},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 12},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 16},
			{Kind: token.Ident, Value: "b", Line: 3, Column: 17},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 18},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 19},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 20},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_return_types_x6", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 9},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 10},
			{Kind: token.Ident, Value: "a", Line: 3, Column: 11},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 12},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 16},
			{Kind: token.Ident, Value: "b", Line: 3, Column: 17},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 18},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 19},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 20},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_return_types_x7", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 9},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 10},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 11},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 12},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_return_types_x8", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 9},
			{Kind: token.Ident, Value: "a", Line: 3, Column: 10},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 11},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 16},
			{Kind: token.Ident, Value: "b", Line: 3, Column: 17},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 18},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 19},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 20},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 21},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_return_types_x9", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 9},
			{Kind: token.Ident, Value: "a", Line: 3, Column: 10},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 11},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 16},
			{Kind: token.Ident, Value: "b", Line: 3, Column: 17},
			{Kind: token.Ident, Value: "b", Line: 3, Column: 19},
			{Kind: token.Ident, Value: "b", Line: 3, Column: 20},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 21},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 22},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 23},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_return_types_x10", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 9},
			{Kind: token.Ident, Value: "a", Line: 3, Column: 10},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 11},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 16},
			{Kind: token.Ident, Value: "b", Line: 3, Column: 17},
			{Kind: token.Ident, Value: "b", Line: 3, Column: 19},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 20},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 21},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 22},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 23},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_return_types_x11", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 9},
			{Kind: token.Ident, Value: "a", Line: 3, Column: 10},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 12},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 16},
			{Kind: token.Ident, Value: "b", Line: 3, Column: 17},
			{Kind: token.KWReturn, Value: "return", Line: 3, Column: 19},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 20},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 21},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 15},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_return_types_x12", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 9},
			{Kind: token.Ident, Value: "a", Line: 3, Column: 10},
			{Kind: token.KWReturn, Value: "return", Line: 3, Column: 19},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 21},
			{Kind: token.Ident, Value: "b", Line: 3, Column: 22},
			{Kind: token.KWReturn, Value: "return", Line: 3, Column: 24},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 25},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 26},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 15},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_return_types_x13", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.KWStruct, Value: "struct", Line: 3, Column: 10},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 14},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 15},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_return_types_x14", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 9},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 10},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 14},
			{Kind: token.KWStruct, Value: "struct", Line: 3, Column: 16},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 24},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 25},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 15},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_return_types_x15", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 9},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 10},
			{Kind: token.KWStruct, Value: "struct", Line: 3, Column: 16},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 24},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 25},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 15},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_return_types_x16", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 9},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 10},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 14},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 16},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 17},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 24},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 25},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 15},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_return_types_x17", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 8},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 9},
			{Kind: token.Ident, Value: "a", Line: 3, Column: 10},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 16},
			{Kind: token.Ident, Value: "b", Line: 3, Column: 17},
			{Kind: token.Ident, Value: "z", Line: 3, Column: 19},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 20},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 21},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 15},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}
