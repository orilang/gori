package parser

import (
	"syscall"
	"testing"

	"github.com/orilang/gori/lexer"
	"github.com/orilang/gori/token"
	"github.com/stretchr/testify/assert"
)

func TestParser_parse_common(t *testing.T) {
	assert := assert.New(t)

	t.Run("err_no_such_file", func(t *testing.T) {
		_, err := NewParser(Config{File: "xxxx.ori"})
		assert.ErrorIs(err, syscall.Errno(2))
	})

	t.Run("err_start_lexing", func(t *testing.T) {
		parse := &Files{Files: []string{"xxxx.ori"}}
		assert.ErrorIs(parse.StartParsing(), syscall.Errno(2))
	})

	t.Run("peek_eof", func(t *testing.T) {
		input := "main"

		lex := lexer.New([]byte(input))
		parse := New(lex.Tokens)
		parse.position = len(input)
		result := parse.peek()

		assert.Equal(token.EOF, result.Kind)
	})

	t.Run("match_true", func(t *testing.T) {
		input := "package"

		lex := lexer.New([]byte(input))
		lex.Tokenize()
		parse := New(lex.Tokens)
		_, ok := parse.match(token.KWPackage)

		assert.Equal(true, ok)
	})

	t.Run("match_false", func(t *testing.T) {
		input := "main"

		lex := lexer.New([]byte(input))
		parse := New(lex.Tokens)
		parse.position = len(input)
		_, ok := parse.match(token.Ident)

		assert.Equal(false, ok)
	})

	t.Run("expect_ok", func(t *testing.T) {
		input := "package"

		lex := lexer.New([]byte(input))
		lex.Tokenize()
		parse := New(lex.Tokens)
		_ = parse.expect(token.KWPackage, "ok")

		assert.Nil(parse.errors)
	})

	t.Run("expect_errors", func(t *testing.T) {
		input := "package"

		lex := lexer.New([]byte(input))
		lex.Tokenize()
		parse := New(lex.Tokens)
		tok := parse.expect(token.Illegal, "nok")
		assert.NotNil(parse.errors)
		assert.Equal(token.KWPackage, tok.Kind)
	})

	t.Run("peekPrecedence_lowest", func(t *testing.T) {
		input := "package"

		lex := lexer.New([]byte(input))
		lex.Tokenize()
		parse := New(lex.Tokens)
		assert.Equal(LOWEST, parse.peekPrecedence())
	})

	t.Run("peekPrecedence_multiplicative", func(t *testing.T) {
		input := "*"

		lex := lexer.New([]byte(input))
		lex.Tokenize()
		parse := New(lex.Tokens)
		assert.Equal(MULTIPLICATIVE, parse.peekPrecedence())
	})
}
