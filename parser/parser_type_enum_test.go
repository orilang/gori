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
			{Kind: token.Assign, Value: "=", Line: 3, Column: 17},
			{Kind: token.Pipe, Value: "|", Line: 3, Column: 3},
			{Kind: token.Ident, Value: "Red", Line: 3, Column: 5},
			{Kind: token.Pipe, Value: "|", Line: 3, Column: 3},
			{Kind: token.Ident, Value: "Blue", Line: 3, Column: 5},
			{Kind: token.Pipe, Value: "|", Line: 3, Column: 3},
			{Kind: token.Ident, Value: "Green", Line: 3, Column: 5},
			{Kind: token.Pipe, Value: "|", Line: 3, Column: 3},
			{Kind: token.Ident, Value: "Yellow", Line: 3, Column: 5},
			{Kind: token.Comment, Value: "// comment", Line: 3, Column: 12},
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
   Eq: "=" @3:17 (kind=49)
   Variants
    Ident: "Red" @3:5 (kind=3)
    Ident: "Blue" @3:5 (kind=3)
    Ident: "Green" @3:5 (kind=3)
    Ident: "Yellow" @3:5 (kind=3)
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
			{Kind: token.Assign, Value: "=", Line: 3, Column: 17},
			{Kind: token.Pipe, Value: "|", Line: 4, Column: 19},
			{Kind: token.Ident, Value: "Red", Line: 4, Column: 21},
			{Kind: token.Pipe, Value: "|", Line: 5, Column: 25},
			{Kind: token.Ident, Value: "Blue", Line: 5, Column: 27},
			{Kind: token.Pipe, Value: "|", Line: 6, Column: 32},
			{Kind: token.Ident, Value: "Green", Line: 6, Column: 34},
			{Kind: token.Pipe, Value: "|", Line: 7, Column: 40},
			{Kind: token.Ident, Value: "Yellow", Line: 7, Column: 42},
			{Kind: token.Comment, Value: "// comment", Line: 3, Column: 49},
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
   Eq: "=" @3:17 (kind=49)
   Variants
    Ident: "Red" @4:21 (kind=3)
    Ident: "Blue" @5:27 (kind=3)
    Ident: "Green" @6:34 (kind=3)
    Ident: "Yellow" @7:42 (kind=3)
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
			{Kind: token.Assign, Value: "=", Line: 3, Column: 17},
			{Kind: token.Assign, Value: "=", Line: 3, Column: 18},
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
			{Kind: token.Assign, Value: "=", Line: 3, Column: 17},
			{Kind: token.Pipe, Value: "|", Line: 4, Column: 19},
			{Kind: token.Pipe, Value: "|", Line: 4, Column: 20},
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
			{Kind: token.Assign, Value: "=", Line: 3, Column: 17},
			{Kind: token.Pipe, Value: "|", Line: 4, Column: 19},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 20},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}
