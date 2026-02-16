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

	t.Run("look_for_x1", func(t *testing.T) {
		input := "*"

		lex := lexer.New([]byte(input))
		lex.Tokenize()
		parse := New(lex.Tokens)
		assert.Equal(false, parse.lookForInForHeader(token.Comma))
	})

	t.Run("look_for_x2", func(t *testing.T) {
		input := "a 1"

		lex := lexer.New([]byte(input))
		lex.Tokenize()
		parse := New(lex.Tokens)
		assert.Equal(true, parse.lookForInForHeader(token.IntLit))
	})

	t.Run("look_for_x3", func(t *testing.T) {
		input := "*"

		lex := lexer.New([]byte(input))
		lex.Tokenize()
		parse := New(lex.Tokens)
		parse.position = len(lex.Tokens)
		assert.Equal(false, parse.lookForInForHeader(token.Comma))
	})

	t.Run("look_for_in_switch_header_x1", func(t *testing.T) {
		input := "*"

		lex := lexer.New([]byte(input))
		lex.Tokenize()
		parse := New(lex.Tokens)
		assert.Equal(false, parse.lookForInSwitchHeader(token.SemiComma))
	})

	t.Run("look_for_in_switch_header_x2", func(t *testing.T) {
		input := "switch x:=f();x"

		lex := lexer.New([]byte(input))
		lex.Tokenize()
		parse := New(lex.Tokens)
		assert.Equal(true, parse.lookForInSwitchHeader(token.SemiComma))
	})

	t.Run("look_for_in_switch_header_x3", func(t *testing.T) {
		input := "*"

		lex := lexer.New([]byte(input))
		lex.Tokenize()
		parse := New(lex.Tokens)
		parse.position = len(lex.Tokens)
		assert.Equal(false, parse.lookForInSwitchHeader(token.SemiComma))
	})

	t.Run("look_for_in_switch_case_header_x1", func(t *testing.T) {
		input := "*"

		lex := lexer.New([]byte(input))
		lex.Tokenize()
		parse := New(lex.Tokens)
		assert.Equal(false, parse.lookForInSwitchCaseHeader(token.Comma))
	})

	t.Run("look_for_in_switch_case_header_x2", func(t *testing.T) {
		input := "case 1,2"

		lex := lexer.New([]byte(input))
		lex.Tokenize()
		parse := New(lex.Tokens)
		assert.Equal(true, parse.lookForInSwitchCaseHeader(token.Comma))
	})

	t.Run("look_for_in_switch_case_header_x3", func(t *testing.T) {
		input := "*"

		lex := lexer.New([]byte(input))
		lex.Tokenize()
		parse := New(lex.Tokens)
		parse.position = len(lex.Tokens)
		assert.Equal(false, parse.lookForInSwitchCaseHeader(token.Comma))
	})

	t.Run("peek_next_eof", func(t *testing.T) {
		input := "main"

		lex := lexer.New([]byte(input))
		parse := New(lex.Tokens)
		parse.position = len(input)
		result := parse.peekNext(parse.position + 1)

		assert.Equal(token.EOF, result.Kind)
	})

	t.Run("is_public_true", func(t *testing.T) {
		input := token.Token{Kind: token.Ident, Value: "A"}
		assert.Equal(true, isPublic(input))
	})

	t.Run("is_public_false", func(t *testing.T) {
		input := token.Token{Kind: token.Ident, Value: "a"}
		assert.Equal(false, isPublic(input))
	})

	t.Run("new_line_since_prev_true", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWType, Value: "type", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "test", Line: 3, Column: 6},
			{Kind: token.KWStruct, Value: "struct", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 18},
			{Kind: token.Ident, Value: "x", Line: 4, Column: 18},
			{Kind: token.KWInt, Value: "int", Line: 4, Column: 18},
			{Kind: token.RBrace, Value: "}", Line: 5, Column: 19},
			{Kind: token.EOF, Value: "", Line: 6, Column: 1},
		}
		p := New(input)
		p.position = 6
		assert.Equal(true, p.newlineSincePrev())
	})

	t.Run("new_line_since_prev_false", func(t *testing.T) {
		input := []token.Token{}
		p := New(input)
		assert.Equal(false, p.newlineSincePrev())
	})
}
