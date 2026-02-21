package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
	"github.com/stretchr/testify/assert"
)

func TestParser_type_sum(t *testing.T) {
	assert := assert.New(t)

	t.Run("x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWSum, Value: "sum", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 16},
			{Kind: token.Ident, Value: "Circle", Line: 3, Column: 18},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 24},
			{Kind: token.Ident, Value: "radius", Line: 3, Column: 25},
			{Kind: token.KWInt, Value: "float", Line: 3, Column: 32},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 37},
			{Kind: token.SemiComma, Value: ";", Line: 3, Column: 38},
			{Kind: token.Ident, Value: "Rect", Line: 3, Column: 40},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 44},
			{Kind: token.Ident, Value: "w", Line: 3, Column: 45},
			{Kind: token.KWInt, Value: "float", Line: 3, Column: 47},
			{Kind: token.Comma, Value: ",", Line: 3, Column: 52},
			{Kind: token.Ident, Value: "h", Line: 3, Column: 54},
			{Kind: token.KWInt, Value: "float", Line: 3, Column: 56},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 61},
			{Kind: token.SemiComma, Value: ";", Line: 3, Column: 62},
			{Kind: token.Ident, Value: "None", Line: 3, Column: 63},
			{Kind: token.RBrace, Value: "}", Line: 3, Column: 67},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Sums
  Type: "type" @3:1 (kind=26)
   Name: "test" @3:6 (kind=3)
   Sum: "sum" @3:11 (kind=75)
   LBrace: "{" @3:16 (kind=41)
    Variants
     Ident: "None" @3:63 (kind=3)
    VariantMethods
     Methods: "Circle" @3:18 (kind=3)
      Params
       Param
        Ident: "radius" @3:25 (kind=3)
        Type
         NameType
          Name: "float" @3:32 (kind=12)
     Methods: "Rect" @3:40 (kind=3)
      Params
       Param
        Ident: "w" @3:45 (kind=3)
        Type
         NameType
          Name: "float" @3:47 (kind=12)
       Param
        Ident: "h" @3:54 (kind=3)
        Type
         NameType
          Name: "float" @3:56 (kind=12)
   RBrace: "}" @3:67 (kind=42)
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
			{Kind: token.KWSum, Value: "sum", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 16},
			{Kind: token.Ident, Value: "Circle", Line: 4, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 9},
			{Kind: token.Ident, Value: "radius", Line: 4, Column: 10},
			{Kind: token.KWInt, Value: "float", Line: 4, Column: 17},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 22},
			{Kind: token.Ident, Value: "Rect", Line: 5, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 5, Column: 7},
			{Kind: token.Ident, Value: "w", Line: 5, Column: 8},
			{Kind: token.KWInt, Value: "float", Line: 5, Column: 10},
			{Kind: token.Comma, Value: ",", Line: 5, Column: 15},
			{Kind: token.Ident, Value: "h", Line: 5, Column: 17},
			{Kind: token.KWInt, Value: "float", Line: 5, Column: 19},
			{Kind: token.RParen, Value: ")", Line: 5, Column: 24},
			{Kind: token.Ident, Value: "None", Line: 6, Column: 3},
			{Kind: token.Comment, Value: "// comment", Line: 6, Column: 8},
			{Kind: token.RBrace, Value: "}", Line: 7, Column: 1},
			{Kind: token.EOF, Value: "", Line: 8, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Sums
  Type: "type" @3:1 (kind=26)
   Name: "test" @3:6 (kind=3)
   Sum: "sum" @3:11 (kind=75)
   LBrace: "{" @3:16 (kind=41)
    Variants
     Ident: "None" @6:3 (kind=3)
    VariantMethods
     Methods: "Circle" @4:3 (kind=3)
      Params
       Param
        Ident: "radius" @4:10 (kind=3)
        Type
         NameType
          Name: "float" @4:17 (kind=12)
     Methods: "Rect" @5:3 (kind=3)
      Params
       Param
        Ident: "w" @5:8 (kind=3)
        Type
         NameType
          Name: "float" @5:10 (kind=12)
       Param
        Ident: "h" @5:17 (kind=3)
        Type
         NameType
          Name: "float" @5:19 (kind=12)
   RBrace: "}" @7:1 (kind=42)
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
			{Kind: token.KWSum, Value: "sum", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 13},
			{Kind: token.RBrace, Value: "}", Line: 3, Column: 24},
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
			{Kind: token.KWSum, Value: "sum", Line: 3, Column: 11},
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
			{Kind: token.KWSum, Value: "sum", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 21},
			{Kind: token.Ident, Value: "X", Line: 4, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 4},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 5},
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
			{Kind: token.KWSum, Value: "sum", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 21},
			{Kind: token.Ident, Value: "X", Line: 4, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 4},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 5},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 6},
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
			{Kind: token.KWSum, Value: "sum", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 21},
			{Kind: token.Ident, Value: "X", Line: 4, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 5, Column: 4},
			{Kind: token.Ident, Value: "a", Line: 5, Column: 5},
			{Kind: token.KWInt, Value: "int", Line: 5, Column: 7},
			{Kind: token.Comma, Value: ",", Line: 4, Column: 11},
			{Kind: token.RParen, Value: ")", Line: 5, Column: 12},
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
			{Kind: token.KWSum, Value: "sum", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 4, Column: 21},
			{Kind: token.Ident, Value: "X", Line: 5, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 5, Column: 4},
			{Kind: token.Ident, Value: "a", Line: 5, Column: 5},
			{Kind: token.KWInt, Value: "int", Line: 5, Column: 7},
			{Kind: token.Ident, Value: "x", Line: 4, Column: 11},
			{Kind: token.RParen, Value: ")", Line: 5, Column: 12},
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
			{Kind: token.KWSum, Value: "sum", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 21},
			{Kind: token.Ident, Value: "X", Line: 4, Column: 3},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 4},
			{Kind: token.Ident, Value: "a", Line: 4, Column: 5},
			{Kind: token.KWStruct, Value: "struct", Line: 4, Column: 7},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 15},
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
			{Kind: token.KWSum, Value: "sum", Line: 3, Column: 11},
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

	t.Run("bad_x9", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWSum, Value: "sum", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 21},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 3},
			{Kind: token.RBrace, Value: "}", Line: 4, Column: 1},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}
