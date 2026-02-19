package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
	"github.com/stretchr/testify/assert"
)

func TestParser_interface_implements(t *testing.T) {
	assert := assert.New(t)

	t.Run("x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.Ident, Value: "X", Line: 3, Column: 1},
			{Kind: token.KWImplements, Value: "implements", Line: 3, Column: 3},
			{Kind: token.Ident, Value: "Y", Line: 3, Column: 15},
			{Kind: token.Ident, Value: "Z", Line: 3, Column: 1},
			{Kind: token.KWImplements, Value: "implements", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "A", Line: 4, Column: 15},
			{Kind: token.Dot, Value: ".", Line: 4, Column: 16},
			{Kind: token.Ident, Value: "B", Line: 4, Column: 17},
			{Kind: token.EOF, Value: "", Line: 5, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Implements
  Type: "X" @3:1 (kind=3)
   Implements: "implements" @3:3 (kind=72)
   Interface
    Ident: "Y" @3:15 (kind=3)
  Type: "Z" @3:1 (kind=3)
   Implements: "implements" @4:3 (kind=72)
   Interface
    Ident: "A" @4:15 (kind=3)
    Dot: "." @4:16 (kind=48)
    Ident: "B" @4:17 (kind=3)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("x2", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.Ident, Value: "X", Line: 3, Column: 1},
			{Kind: token.KWImplements, Value: "implements", Line: 3, Column: 3},
			{Kind: token.Ident, Value: "Y", Line: 3, Column: 15},
			{Kind: token.SemiComma, Value: ";", Line: 3, Column: 16},
			{Kind: token.Ident, Value: "Z", Line: 3, Column: 17},
			{Kind: token.KWImplements, Value: "implements", Line: 3, Column: 19},
			{Kind: token.Ident, Value: "A", Line: 3, Column: 30},
			{Kind: token.Dot, Value: ".", Line: 3, Column: 31},
			{Kind: token.Ident, Value: "B", Line: 3, Column: 32},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Implements
  Type: "X" @3:1 (kind=3)
   Implements: "implements" @3:3 (kind=72)
   Interface
    Ident: "Y" @3:15 (kind=3)
  Type: "Z" @3:17 (kind=3)
   Implements: "implements" @3:19 (kind=72)
   Interface
    Ident: "A" @3:30 (kind=3)
    Dot: "." @3:31 (kind=48)
    Ident: "B" @3:32 (kind=3)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("bad_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.Ident, Value: "X", Line: 3, Column: 1},
			{Kind: token.KWImplements, Value: "implements", Line: 3, Column: 3},
			{Kind: token.Ident, Value: "Y", Line: 3, Column: 15},
			{Kind: token.SemiComma, Value: ";", Line: 3, Column: 16},
			{Kind: token.Ident, Value: "Z", Line: 3, Column: 17},
			{Kind: token.KWImplements, Value: "implements", Line: 3, Column: 19},
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

			{Kind: token.Ident, Value: "X", Line: 3, Column: 1},
			{Kind: token.KWImplements, Value: "implements", Line: 3, Column: 3},
			{Kind: token.Ident, Value: "Y", Line: 3, Column: 15},
			{Kind: token.SemiComma, Value: ";", Line: 3, Column: 16},
			{Kind: token.Ident, Value: "Z", Line: 3, Column: 17},
			{Kind: token.KWImplements, Value: "implements", Line: 3, Column: 19},
			{Kind: token.Ident, Value: "A", Line: 3, Column: 30},
			{Kind: token.Dot, Value: ".", Line: 3, Column: 31},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}
