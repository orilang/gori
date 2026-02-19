package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
	"github.com/stretchr/testify/assert"
)

func TestParser_interface_type(t *testing.T) {
	assert := assert.New(t)

	t.Run("empty_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWInterface, Value: "interface", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 22},
			{Kind: token.RBrace, Value: "}", Line: 3, Column: 23},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Interfaces
  Type: "type" @3:1 (kind=26)
  Name: "test" @3:6 (kind=3)
  Interface: "interface" @3:11 (kind=28)
  LBrace: "{" @3:22 (kind=41)
  RBrace: "}" @3:23 (kind=42)
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
			{Kind: token.KWInterface, Value: "interface", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 22},
			{Kind: token.RBrace, Value: "}", Line: 3, Column: 23},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Interfaces
  Type: "type" @3:1 (kind=26)
  Name: "Test" @3:6 (kind=3)
  Interface: "interface" @3:11 (kind=28)
  Public: true
  LBrace: "{" @3:22 (kind=41)
  RBrace: "}" @3:23 (kind=42)
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
			{Kind: token.KWInterface, Value: "interface", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 18},
			{Kind: token.Ident, Value: "X", Line: 4, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 4},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 5},
			{Kind: token.Ident, Value: "error", Line: 4, Column: 7},
			{Kind: token.RBrace, Value: "}", Line: 5, Column: 1},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Interfaces
  Type: "type" @3:1 (kind=26)
  Name: "test" @3:6 (kind=3)
  Interface: "interface" @3:11 (kind=28)
  LBrace: "{" @3:18 (kind=41)
   Name: "X" @4:3 (kind=3)
   Params
    (none)
   Results
     Param
      Type
       NameType
        Name: "error" @4:7 (kind=3)
  RBrace: "}" @5:1 (kind=42)
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
			{Kind: token.KWInterface, Value: "interface", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 21},
			{Kind: token.Ident, Value: "X", Line: 3, Column: 22},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 23},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 24},
			{Kind: token.Ident, Value: "error", Line: 3, Column: 26},
			{Kind: token.RBrace, Value: "}", Line: 3, Column: 31},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Interfaces
  Type: "type" @3:1 (kind=26)
  Name: "test" @3:6 (kind=3)
  Interface: "interface" @3:11 (kind=28)
  LBrace: "{" @3:21 (kind=41)
   Name: "X" @3:22 (kind=3)
   Params
    (none)
   Results
     Param
      Type
       NameType
        Name: "error" @3:26 (kind=3)
  RBrace: "}" @3:31 (kind=42)
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
			{Kind: token.KWInterface, Value: "interface", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 21},
			{Kind: token.Ident, Value: "X", Line: 3, Column: 22},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 23},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 24},
			{Kind: token.Ident, Value: "error", Line: 3, Column: 26},
			{Kind: token.SemiComma, Value: ";", Line: 3, Column: 31},
			{Kind: token.Ident, Value: "Y", Line: 3, Column: 32},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 33},
			{Kind: token.Ident, Value: "a", Line: 3, Column: 34},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 36},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 39},
			{Kind: token.Ident, Value: "error", Line: 3, Column: 41},
			{Kind: token.RBrace, Value: "}", Line: 3, Column: 46},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Interfaces
  Type: "type" @3:1 (kind=26)
  Name: "test" @3:6 (kind=3)
  Interface: "interface" @3:11 (kind=28)
  LBrace: "{" @3:21 (kind=41)
   Name: "X" @3:22 (kind=3)
   Params
    (none)
   Results
     Param
      Type
       NameType
        Name: "error" @3:26 (kind=3)
   Name: "Y" @3:32 (kind=3)
   Params
    Param
     Ident: "a" @3:34 (kind=3)
     Type
      NameType
       Name: "int" @3:36 (kind=12)
   Results
     Param
      Type
       NameType
        Name: "error" @3:41 (kind=3)
  RBrace: "}" @3:46 (kind=42)
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
			{Kind: token.KWInterface, Value: "interface", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 21},
			{Kind: token.Ident, Value: "X", Line: 4, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 4},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 5},
			{Kind: token.Ident, Value: "error", Line: 4, Column: 6},
			{Kind: token.Ident, Value: "Y", Line: 5, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 5, Column: 4},
			{Kind: token.Ident, Value: "a", Line: 5, Column: 5},
			{Kind: token.KWInt, Value: "int", Line: 5, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 5, Column: 10},
			{Kind: token.Ident, Value: "error", Line: 5, Column: 12},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 1},
			{Kind: token.EOF, Value: "", Line: 7, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Interfaces
  Type: "type" @3:1 (kind=26)
  Name: "test" @3:6 (kind=3)
  Interface: "interface" @3:11 (kind=28)
  LBrace: "{" @3:21 (kind=41)
   Name: "X" @4:3 (kind=3)
   Params
    (none)
   Results
     Param
      Type
       NameType
        Name: "error" @4:6 (kind=3)
   Name: "Y" @5:3 (kind=3)
   Params
    Param
     Ident: "a" @5:5 (kind=3)
     Type
      NameType
       Name: "int" @5:7 (kind=12)
   Results
     Param
      Type
       NameType
        Name: "error" @5:12 (kind=3)
  RBrace: "}" @6:1 (kind=42)
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
			{Kind: token.KWInterface, Value: "interface", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 21},
			{Kind: token.Ident, Value: "X", Line: 4, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 4},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 5},
			{Kind: token.Ident, Value: "error", Line: 4, Column: 6},
			{Kind: token.Ident, Value: "Y", Line: 5, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 5, Column: 4},
			{Kind: token.Ident, Value: "a", Line: 5, Column: 5},
			{Kind: token.KWInt, Value: "int", Line: 5, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 5, Column: 10},
			{Kind: token.Ident, Value: "error", Line: 5, Column: 12},
			{Kind: token.Comment, Value: "// comment", Line: 5, Column: 19},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 1},
			{Kind: token.EOF, Value: "", Line: 7, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Interfaces
  Type: "type" @3:1 (kind=26)
  Name: "test" @3:6 (kind=3)
  Interface: "interface" @3:11 (kind=28)
  LBrace: "{" @3:21 (kind=41)
   Name: "X" @4:3 (kind=3)
   Params
    (none)
   Results
     Param
      Type
       NameType
        Name: "error" @4:6 (kind=3)
   Name: "Y" @5:3 (kind=3)
   Params
    Param
     Ident: "a" @5:5 (kind=3)
     Type
      NameType
       Name: "int" @5:7 (kind=12)
   Results
     Param
      Type
       NameType
        Name: "error" @5:12 (kind=3)
  RBrace: "}" @6:1 (kind=42)
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
			{Kind: token.KWInterface, Value: "interface", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 21},
			{Kind: token.Ident, Value: "X", Line: 4, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 4},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 5},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 6},
			{Kind: token.Ident, Value: "y", Line: 4, Column: 7},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 9},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 12},
			{Kind: token.Ident, Value: "z", Line: 4, Column: 16},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 16},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 19},
			{Kind: token.RBrace, Value: "}", Line: 5, Column: 1},
			{Kind: token.EOF, Value: "", Line: 6, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Interfaces
  Type: "type" @3:1 (kind=26)
  Name: "test" @3:6 (kind=3)
  Interface: "interface" @3:11 (kind=28)
  LBrace: "{" @3:21 (kind=41)
   Name: "X" @4:3 (kind=3)
   Params
    (none)
   Results
    LParent: "(" @4:6 (kind=39)
     Param
      Ident: "y" @4:7 (kind=3)
      Type
       NameType
        Name: "int" @4:9 (kind=12)
     Param
      Ident: "z" @4:16 (kind=3)
      Type
       NameType
        Name: "int" @4:16 (kind=12)
    RParent: ")" @4:19 (kind=40)
  RBrace: "}" @5:1 (kind=42)
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
			{Kind: token.KWInterface, Value: "interface", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 21},
			{Kind: token.Ident, Value: "X", Line: 4, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 4},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 5},
			{Kind: token.Ident, Value: "error", Line: 4, Column: 6},
			{Kind: token.SemiComma, Value: ";", Line: 4, Column: 11},
			{Kind: token.Ident, Value: "Y", Line: 4, Column: 12},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 13},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 14},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 16},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 19},
			{Kind: token.Ident, Value: "error", Line: 4, Column: 21},
			{Kind: token.Comment, Value: "// comment", Line: 4, Column: 27},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Interfaces
  Type: "type" @3:1 (kind=26)
  Name: "test" @3:6 (kind=3)
  Interface: "interface" @3:11 (kind=28)
  LBrace: "{" @3:21 (kind=41)
   Name: "X" @4:3 (kind=3)
   Params
    (none)
   Results
     Param
      Type
       NameType
        Name: "error" @4:6 (kind=3)
   Name: "Y" @4:12 (kind=3)
   Params
    Param
     Ident: "a" @4:14 (kind=3)
     Type
      NameType
       Name: "int" @4:16 (kind=12)
   Results
     Param
      Type
       NameType
        Name: "error" @4:21 (kind=3)
  RBrace: "}" @4:1 (kind=42)
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
			{Kind: token.KWInterface, Value: "interface", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 21},
			{Kind: token.Ident, Value: "X", Line: 4, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 4},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 5},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 13},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 14},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 15},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 16},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 17},
			{Kind: token.Comment, Value: "// comment", Line: 4, Column: 19},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Interfaces
  Type: "type" @3:1 (kind=26)
  Name: "test" @3:6 (kind=3)
  Interface: "interface" @3:11 (kind=28)
  LBrace: "{" @3:21 (kind=41)
   Name: "X" @4:3 (kind=3)
   Params
    (none)
   Results
    LParent: "(" @4:13 (kind=39)
     Param
      Type
       NameType
        Name: "a" @4:14 (kind=3)
     Param
      Type
       NameType
        Name: "int" @4:16 (kind=12)
    RParent: ")" @4:17 (kind=40)
  RBrace: "}" @4:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("x9", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWInterface, Value: "interface", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 21},
			{Kind: token.Ident, Value: "X", Line: 4, Column: 3},
			{Kind: token.Dot, Value: ".", Line: 4, Column: 4},
			{Kind: token.Ident, Value: "Y", Line: 4, Column: 5},
			{Kind: token.Comment, Value: "// comment", Line: 4, Column: 7},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Interfaces
  Type: "type" @3:1 (kind=26)
  Name: "test" @3:6 (kind=3)
  Interface: "interface" @3:11 (kind=28)
  LBrace: "{" @3:21 (kind=41)
    Embeds
     Ident: "X" @4:3 (kind=3)
     Dot: "." @4:4 (kind=48)
     Ident: "Y" @4:5 (kind=3)
  RBrace: "}" @4:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("bad_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWInterface, Value: "interface", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 18},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 20},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 22},
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
			{Kind: token.KWInterface, Value: "interface", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 18},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 20},
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
			{Kind: token.KWInterface, Value: "interface", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 21},
			{Kind: token.Ident, Value: "X", Line: 4, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 4},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 5},
			{Kind: token.Ident, Value: "error", Line: 4, Column: 6},
			{Kind: token.Ident, Value: "Y", Line: 5, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 5, Column: 4},
			{Kind: token.Ident, Value: "a", Line: 5, Column: 5},
			{Kind: token.KWInt, Value: "int", Line: 5, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 5, Column: 10},
			{Kind: token.Ident, Value: "error", Line: 5, Column: 12},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 19},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 1},
			{Kind: token.EOF, Value: "", Line: 7, Column: 1},
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
			{Kind: token.KWInterface, Value: "interface", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 21},
			{Kind: token.Ident, Value: "X", Line: 4, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 4},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 5},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 6},
			{Kind: token.Ident, Value: "error", Line: 4, Column: 7},
			{Kind: token.Ident, Value: "Y", Line: 5, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 5, Column: 4},
			{Kind: token.Ident, Value: "a", Line: 5, Column: 5},
			{Kind: token.KWInt, Value: "int", Line: 5, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 5, Column: 10},
			{Kind: token.Ident, Value: "error", Line: 5, Column: 12},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 1},
			{Kind: token.EOF, Value: "", Line: 7, Column: 1},
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

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWInterface, Value: "interface", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 21},
			{Kind: token.Ident, Value: "X", Line: 4, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 5, Column: 4},
			{Kind: token.Ident, Value: "a", Line: 5, Column: 5},
			{Kind: token.KWInt, Value: "int", Line: 5, Column: 7},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 11},
			{Kind: token.RParen, Value: ")", Line: 5, Column: 12},
			{Kind: token.Ident, Value: "error", Line: 5, Column: 14},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 1},
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

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWInterface, Value: "interface", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 21},
			{Kind: token.Ident, Value: "X", Line: 5, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 5, Column: 4},
			{Kind: token.Ident, Value: "a", Line: 5, Column: 5},
			{Kind: token.KWInt, Value: "int", Line: 5, Column: 7},
			{Kind: token.Ident, Value: "x", Line: 4, Column: 11},
			{Kind: token.RParen, Value: ")", Line: 5, Column: 12},
			{Kind: token.Ident, Value: "error", Line: 5, Column: 14},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 1},
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

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWInterface, Value: "interface", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 21},
			{Kind: token.Ident, Value: "X", Line: 4, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 4},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 5},
			{Kind: token.KWStruct, Value: "struct", Line: 4, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 15},
			{Kind: token.Ident, Value: "error", Line: 4, Column: 17},
			{Kind: token.RBrace, Value: "}", Line: 5, Column: 1},
			{Kind: token.EOF, Value: "", Line: 6, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_x8", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWInterface, Value: "interface", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 21},
			{Kind: token.Ident, Value: "X", Line: 4, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 4},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 5},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 12},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 13},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 14},
			{Kind: token.RBrace, Value: "}", Line: 5, Column: 1},
			{Kind: token.EOF, Value: "", Line: 6, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_x9", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWInterface, Value: "interface", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 21},
			{Kind: token.Ident, Value: "X", Line: 4, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 4},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 5},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 12},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 13},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 14},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 15},
			{Kind: token.RBrace, Value: "}", Line: 5, Column: 1},
			{Kind: token.EOF, Value: "", Line: 6, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_x10", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWInterface, Value: "interface", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 21},
			{Kind: token.Ident, Value: "X", Line: 4, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 4},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 5},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 12},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 13},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 14},
			{Kind: token.Ident, Value: "b", Line: 4, Column: 16},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 17},
			{Kind: token.Ident, Value: "c", Line: 4, Column: 18},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 19},
			{Kind: token.RBrace, Value: "}", Line: 5, Column: 1},
			{Kind: token.EOF, Value: "", Line: 6, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_x11", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWInterface, Value: "interface", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 21},
			{Kind: token.Ident, Value: "X", Line: 4, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 4},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 5},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 12},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 13},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 14},
			{Kind: token.Ident, Value: "b", Line: 4, Column: 16},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 17},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 18},
			{Kind: token.RBrace, Value: "}", Line: 5, Column: 1},
			{Kind: token.EOF, Value: "", Line: 6, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_x12", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWInterface, Value: "interface", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 21},
			{Kind: token.Ident, Value: "X", Line: 4, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 4},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 5},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 13},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 14},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 15},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 16},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 21},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 22},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_x13", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWInterface, Value: "interface", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 21},
			{Kind: token.Ident, Value: "X", Line: 4, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 4},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 5},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 13},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 14},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 15},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 16},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 21},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 22},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_x14", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWInterface, Value: "interface", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 21},
			{Kind: token.Ident, Value: "X", Line: 4, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 4},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 5},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 6},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_x15", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWInterface, Value: "interface", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 21},
			{Kind: token.Ident, Value: "X", Line: 4, Column: 3},
			{Kind: token.Dot, Value: ".", Line: 4, Column: 4},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}
