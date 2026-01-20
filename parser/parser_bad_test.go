package parser

import (
	"testing"

	"github.com/orilang/gori/token"
	"github.com/stretchr/testify/assert"
)

func TestParser_bad(t *testing.T) {
	assert := assert.New(t)

	t.Run("function_params_bad", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWFunc, Value: "func", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 1, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 1, Column: 7},

			{Kind: token.Ident, Value: "a", Line: 1, Column: 8},
			{Kind: token.Not, Value: "!", Line: 1, Column: 10},

			{Kind: token.Comma, Value: ",", Line: 1, Column: 11},

			{Kind: token.Ident, Value: "b", Line: 1, Column: 13},
			{Kind: token.KWString, Value: "string", Line: 1, Column: 15},

			{Kind: token.RParen, Value: ")", Line: 1, Column: 21},
			{Kind: token.LBrace, Value: "{", Line: 1, Column: 22},

			{Kind: token.RBrace, Value: "}", Line: 1, Column: 23},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.parseFuncDecl()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("stmt", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.Not, Value: "!"},
			{Kind: token.Ident, Value: "a"},
			{Kind: token.KWInt, Value: "int"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "0"},
			{Kind: token.EOF, Value: ""},
		}

		parser := New(input)
		pr := parser.parseStmt()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("expr", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "a"},
			{Kind: token.KWInt, Value: "int"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.Not, Value: "!"},
			{Kind: token.EOF, Value: ""},
		}

		parser := New(input)
		pr := parser.parseExpr(LOWEST)
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}
