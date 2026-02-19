package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
	"github.com/stretchr/testify/assert"
)

func TestParser_type_struct(t *testing.T) {
	assert := assert.New(t)

	t.Run("empty_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWStruct, Value: "struct", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 18},
			{Kind: token.RBrace, Value: "}", Line: 3, Column: 19},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Structs
  Type: "type" @3:1 (kind=26)
  Name: "test" @3:6 (kind=3)
  Struct: "struct" @3:11 (kind=27)
  LBrace: "{" @3:18 (kind=41)
  RBrace: "}" @3:19 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("empty_x2", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "Test", Line: 3, Column: 6},
			{Kind: token.KWStruct, Value: "struct", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 18},
			{Kind: token.RBrace, Value: "}", Line: 3, Column: 19},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Structs
  Type: "type" @3:1 (kind=26)
  Name: "Test" @3:6 (kind=3)
  Struct: "struct" @3:11 (kind=27)
  Public: true
  LBrace: "{" @3:18 (kind=41)
  RBrace: "}" @3:19 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWStruct, Value: "struct", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 18},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 20},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 22},
			{Kind: token.RBrace, Value: "}", Line: 3, Column: 19},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Structs
  Type: "type" @3:1 (kind=26)
  Name: "test" @3:6 (kind=3)
  Struct: "struct" @3:11 (kind=27)
  LBrace: "{" @3:18 (kind=41)
   Name: "x" @3:20 (kind=3)
   Type:
    NameType
     Name: "int" @3:22 (kind=12)
  RBrace: "}" @3:19 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("x2", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWStruct, Value: "struct", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 18},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 20},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 22},
			{Kind: token.RBrace, Value: "}", Line: 3, Column: 19},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Structs
  Type: "type" @3:1 (kind=26)
  Name: "test" @3:6 (kind=3)
  Struct: "struct" @3:11 (kind=27)
  LBrace: "{" @3:18 (kind=41)
   Name: "x" @3:20 (kind=3)
   Type:
    NameType
     Name: "int" @3:22 (kind=12)
  RBrace: "}" @3:19 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("x3", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWStruct, Value: "struct", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 18},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 20},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 22},
			{Kind: token.SemiComma, Value: ";", Line: 3, Column: 25},
			{Kind: token.Ident, Value: "Y", Line: 3, Column: 27},
			{Kind: token.KWInt, Value: "string", Line: 3, Column: 29},
			{Kind: token.RBrace, Value: "}", Line: 3, Column: 19},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Structs
  Type: "type" @3:1 (kind=26)
  Name: "test" @3:6 (kind=3)
  Struct: "struct" @3:11 (kind=27)
  LBrace: "{" @3:18 (kind=41)
   Name: "x" @3:20 (kind=3)
   Type:
    NameType
     Name: "int" @3:22 (kind=12)
   Name: "Y" @3:27 (kind=3)
   Public: true
   Type:
    NameType
     Name: "string" @3:29 (kind=12)
  RBrace: "}" @3:19 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("x4", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWStruct, Value: "struct", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 18},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 20},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 22},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 1},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Structs
  Type: "type" @3:1 (kind=26)
  Name: "test" @3:6 (kind=3)
  Struct: "struct" @3:11 (kind=27)
  LBrace: "{" @3:18 (kind=41)
   Name: "x" @3:20 (kind=3)
   Type:
    NameType
     Name: "int" @3:22 (kind=12)
  RBrace: "}" @4:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("x5", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWStruct, Value: "struct", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 18},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 20},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 22},
			{Kind: token.Comment, Value: "// comment", Line: 3, Column: 26},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 1},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Structs
  Type: "type" @3:1 (kind=26)
  Name: "test" @3:6 (kind=3)
  Struct: "struct" @3:11 (kind=27)
  LBrace: "{" @3:18 (kind=41)
   Name: "x" @3:20 (kind=3)
   Type:
    NameType
     Name: "int" @3:22 (kind=12)
  RBrace: "}" @4:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("x6", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWStruct, Value: "struct", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 18},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 20},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 22},
			{Kind: token.Assign, Value: "=", Line: 3, Column: 26},
			{Kind: token.IntLit, Value: "5", Line: 3, Column: 28},
			{Kind: token.RBrace, Value: "}", Line: 3, Column: 30},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Structs
  Type: "type" @3:1 (kind=26)
  Name: "test" @3:6 (kind=3)
  Struct: "struct" @3:11 (kind=27)
  LBrace: "{" @3:18 (kind=41)
   Name: "x" @3:20 (kind=3)
   Type:
    NameType
     Name: "int" @3:22 (kind=12)
   Eq: "=" @3:26 (kind=49)
   IntLitExpr
    Value: "5" @3:28 (kind=4)
  RBrace: "}" @3:30 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("x7", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWStruct, Value: "struct", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 18},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 20},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 22},
			{Kind: token.Assign, Value: "=", Line: 3, Column: 26},
			{Kind: token.IntLit, Value: "5", Line: 3, Column: 28},
			{Kind: token.Ident, Value: "y", Line: 4, Column: 20},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 22},
			{Kind: token.Assign, Value: "=", Line: 4, Column: 26},
			{Kind: token.IntLit, Value: "5", Line: 4, Column: 28},
			{Kind: token.RBrace, Value: "}", Line: 5, Column: 30},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Structs
  Type: "type" @3:1 (kind=26)
  Name: "test" @3:6 (kind=3)
  Struct: "struct" @3:11 (kind=27)
  LBrace: "{" @3:18 (kind=41)
   Name: "x" @3:20 (kind=3)
   Type:
    NameType
     Name: "int" @3:22 (kind=12)
   Eq: "=" @3:26 (kind=49)
   IntLitExpr
    Value: "5" @3:28 (kind=4)
   Name: "y" @4:20 (kind=3)
   Type:
    NameType
     Name: "int" @4:22 (kind=12)
   Eq: "=" @4:26 (kind=49)
   IntLitExpr
    Value: "5" @4:28 (kind=4)
  RBrace: "}" @5:30 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("x8", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWStruct, Value: "struct", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 18},
			{Kind: token.Ident, Value: "x", Line: 4, Column: 20},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 22},
			{Kind: token.Assign, Value: "=", Line: 4, Column: 27},
			{Kind: token.Ident, Value: "join", Line: 4, Column: 11},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 15},
			{Kind: token.StringLit, Value: "a", Line: 4, Column: 16},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 19},
			{Kind: token.StringLit, Value: "b", Line: 4, Column: 20},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 23},
			{Kind: token.RBrace, Value: "}", Line: 5, Column: 1},
			{Kind: token.EOF, Value: "", Line: 6, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Structs
  Type: "type" @3:1 (kind=26)
  Name: "test" @3:6 (kind=3)
  Struct: "struct" @3:11 (kind=27)
  LBrace: "{" @3:18 (kind=41)
   Name: "x" @4:20 (kind=3)
   Type:
    NameType
     Name: "int" @4:22 (kind=12)
   Eq: "=" @4:27 (kind=49)
   CallExpr
    Callee
     IdentExpr
      Name: "join" @4:11 (kind=3)
    LParent: "(" @4:15 (kind=39)
    Args:
     StringLitExpr
      Value: "a" @4:16 (kind=6)
     StringLitExpr
      Value: "b" @4:20 (kind=6)
    RParent: ")" @4:23 (kind=40)
  RBrace: "}" @5:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("bad_type", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.Ident, Value: "structt", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 18},
			{Kind: token.RBrace, Value: "}", Line: 3, Column: 19},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWStruct, Value: "struct", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 18},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 20},
			{Kind: token.RBrace, Value: "}", Line: 3, Column: 21},
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

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWStruct, Value: "struct", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 18},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 20},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 22},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 25},
			{Kind: token.RBrace, Value: "}", Line: 3, Column: 19},
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

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWStruct, Value: "struct", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 18},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 20},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 22},
			{Kind: token.SemiComma, Value: ";", Line: 3, Column: 25},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 26},
			{Kind: token.RBrace, Value: "}", Line: 3, Column: 19},
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

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWStruct, Value: "struct", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 18},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 20},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 22},
			{Kind: token.Assign, Value: "=", Line: 3, Column: 25},
			{Kind: token.RBrace, Value: "}", Line: 3, Column: 19},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}
