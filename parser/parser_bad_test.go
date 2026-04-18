package parser

import (
	"testing"

	"github.com/orilang/gori/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParser_bad(t *testing.T) {
	assert := assert.New(t)

	t.Run("function_params_bad", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `func x(a!,b string){}
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseFuncDecl()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("stmt", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `!a=int=0
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseStmt()
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})

	t.Run("expr", func(t *testing.T) {
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		assert.Nil(err)
		data := `var a int=!
`
		parser := New(lex.FetchTokensFromString(data))
		pr := parser.parseExpr(LOWEST)
		assert.NotNil(pr)
		assert.Greater(len(parser.Errors), 0)
	})
}
