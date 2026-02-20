package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
	"github.com/stretchr/testify/assert"
)

func TestParser_type_enum(t *testing.T) {
	assert := assert.New(t)

	t.Run("x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "Color", Line: 3, Column: 6},
			{Kind: token.KWEnum, Value: "enum", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 17},
			{Kind: token.Ident, Value: "Red", Line: 3, Column: 19},
			{Kind: token.SemiComma, Value: ";", Line: 3, Column: 22},
			{Kind: token.Ident, Value: "Blue", Line: 3, Column: 23},
			{Kind: token.SemiComma, Value: ";", Line: 3, Column: 27},
			{Kind: token.Ident, Value: "Green", Line: 3, Column: 28},
			{Kind: token.SemiComma, Value: ";", Line: 3, Column: 33},
			{Kind: token.Ident, Value: "Yellow", Line: 3, Column: 34},
			{Kind: token.RBrace, Value: "}", Line: 3, Column: 40},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Enums
  Type: "type" @3:1 (kind=26)
   Name: "Color" @3:6 (kind=3)
   Public: true
   Enum: "enum" @3:11 (kind=74)
   LBrace: "{" @3:17 (kind=41)
   Variants
    Ident: "Red" @3:19 (kind=3)
    Ident: "Blue" @3:23 (kind=3)
    Ident: "Green" @3:28 (kind=3)
    Ident: "Yellow" @3:34 (kind=3)
   RBrace: "}" @3:40 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("x2", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "Color", Line: 3, Column: 6},
			{Kind: token.KWEnum, Value: "enum", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 17},
			{Kind: token.Ident, Value: "Red", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "Blue", Line: 5, Column: 3},
			{Kind: token.Ident, Value: "Green", Line: 6, Column: 3},
			{Kind: token.Ident, Value: "Yellow", Line: 7, Column: 3},
			{Kind: token.Comment, Value: "// comment", Line: 7, Column: 10},
			{Kind: token.RBrace, Value: "}", Line: 8, Column: 1},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Enums
  Type: "type" @3:1 (kind=26)
   Name: "Color" @3:6 (kind=3)
   Public: true
   Enum: "enum" @3:11 (kind=74)
   LBrace: "{" @3:17 (kind=41)
   Variants
    Ident: "Red" @4:3 (kind=3)
    Ident: "Blue" @5:3 (kind=3)
    Ident: "Green" @6:3 (kind=3)
    Ident: "Yellow" @7:3 (kind=3)
   RBrace: "}" @8:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("bad_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "Color", Line: 3, Column: 6},
			{Kind: token.KWEnum, Value: "enum", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 17},
			{Kind: token.RBrace, Value: "}", Line: 3, Column: 18},
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
			{Kind: token.Ident, Value: "Color", Line: 3, Column: 6},
			{Kind: token.KWEnum, Value: "enum", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 17},
			{Kind: token.KWString, Value: "string", Line: 3, Column: 19},
			{Kind: token.RBrace, Value: "}", Line: 3, Column: 26},
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
			{Kind: token.Ident, Value: "Color", Line: 3, Column: 6},
			{Kind: token.KWEnum, Value: "enum", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 17},
			{Kind: token.Ident, Value: "Red", Line: 3, Column: 19},
			{Kind: token.Ident, Value: "Green", Line: 3, Column: 24},
			{Kind: token.RBrace, Value: "}", Line: 3, Column: 30},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}
