package lexer

import (
	"syscall"
	"testing"

	"github.com/orilang/gori/token"
	"github.com/stretchr/testify/assert"
)

func TestLexer_lexer(t *testing.T) {
	assert := assert.New(t)

	t.Run("err_no_such_file", func(t *testing.T) {
		_, err := NewLexer(Config{File: "xxxx.ori"})
		assert.ErrorIs(err, syscall.Errno(2))
	})

	t.Run("err_start_lexing", func(t *testing.T) {
		lex := &Files{Files: []string{"xxxx.ori"}}
		assert.ErrorIs(lex.StartLexing(), syscall.Errno(2))
	})

	t.Run("success", func(t *testing.T) {
		lex, err := NewLexer(Config{Directory: "testdata/success"})
		assert.Nil(err)
		assert.Nil(lex.StartLexing())
	})

	t.Run("illegal", func(t *testing.T) {
		lex, err := NewLexer(Config{Directory: "testdata/illegal"})
		assert.Nil(err)
		assert.Nil(lex.StartLexing())
	})

	t.Run("basic", func(t *testing.T) {
		input := `package main

func main() {}
`
		result := []token.Token{
			{Kind: token.KWPackage, Value: "package"},
			{Kind: token.Ident, Value: "main"},
			{Kind: token.KWFunc, Value: "func"},
			{Kind: token.Ident, Value: "main"},
			{Kind: token.LParen, Value: "("},
			{Kind: token.RParen, Value: ")"},
			{Kind: token.LBrace, Value: "{"},
			{Kind: token.RBrace, Value: "}"},
		}
		lex := New(input)
		lex.Tokenize()
		for i, r := range result {
			assert.Equal(r.Kind, lex.Tokens[i].Kind)
			assert.Equal(r.Value, lex.Tokens[i].Value)
		}
		assert.Equal(len(result), len(lex.Tokens))
	})

	t.Run("vars", func(t *testing.T) {
		input := `package main

func main() {
  var a int = 0
	var b uint = 0
  var c int8 = 0
  var d int32 = 0
  var e int64 = 120

  var f float = 0
  var g float32 = 0
  var h float64 = 0
  var pi float64 = 3.14
  var pi2 float64 = 3.14 // comment
  var pi3 float64 = 3.141,592,653,59

  var i bool = false
  // new test
  var j bool = true // test
  var bt string = "true"
  var bf string = "false"
}
`
		result := []token.Token{
			{Kind: token.KWPackage, Value: "package"},
			{Kind: token.Ident, Value: "main"},
			{Kind: token.KWFunc, Value: "func"},
			{Kind: token.Ident, Value: "main"},
			{Kind: token.LParen, Value: "("},
			{Kind: token.RParen, Value: ")"},
			{Kind: token.LBrace, Value: "{"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "a"},
			{Kind: token.KWInt, Value: "int"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "0"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "b"},
			{Kind: token.KWUint, Value: "uint"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "0"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "c"},
			{Kind: token.KWInt8, Value: "int8"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "0"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "d"},
			{Kind: token.KWInt32, Value: "int32"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "0"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "e"},
			{Kind: token.KWInt64, Value: "int64"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "120"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "f"},
			{Kind: token.KWFloat, Value: "float"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "0"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "g"},
			{Kind: token.KWFloat32, Value: "float32"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "0"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "h"},
			{Kind: token.KWFloat64, Value: "float64"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "0"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "pi"},
			{Kind: token.KWFloat64, Value: "float64"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.FloatLit, Value: "3.14"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "pi2"},
			{Kind: token.KWFloat64, Value: "float64"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.FloatLit, Value: "3.14"},
			{Kind: token.Comment, Value: "// comment"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "pi3"},
			{Kind: token.KWFloat64, Value: "float64"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.FloatLit, Value: "3.141,592,653,59"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "i"},
			{Kind: token.KWBool, Value: "bool"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.BoolLit, Value: "false"},

			{Kind: token.Comment, Value: "// new test"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "j"},
			{Kind: token.KWBool, Value: "bool"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.BoolLit, Value: "true"},
			{Kind: token.Comment, Value: "// test"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "bt"},
			{Kind: token.KWString, Value: "string"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.StringLit, Value: `"true"`},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "bf"},
			{Kind: token.KWString, Value: "string"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.StringLit, Value: `"false"`},

			{Kind: token.RBrace, Value: "}"},
		}
		lex := New(input)
		lex.Tokenize()
		for i, r := range result {
			assert.Equal(r.Kind, lex.Tokens[i].Kind, r.Value)
			assert.Equal(r.Value, lex.Tokens[i].Value)
		}
		assert.Equal(len(result), len(lex.Tokens))
	})

	t.Run("math", func(t *testing.T) {
		input := `package main

func main() {
  var add int = 2 + 2
  var sub int -= 1
  var plus int += 1
  var multiply int = 2 * 2
  var substract int = 2 - 2
  var divide int = 2 / 2
  var modulo int = 2 % 2
}
`
		result := []token.Token{
			{Kind: token.KWPackage, Value: "package"},
			{Kind: token.Ident, Value: "main"},
			{Kind: token.KWFunc, Value: "func"},
			{Kind: token.Ident, Value: "main"},
			{Kind: token.LParen, Value: "("},
			{Kind: token.RParen, Value: ")"},
			{Kind: token.LBrace, Value: "{"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "add"},
			{Kind: token.KWInt, Value: "int"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "2"},
			{Kind: token.Plus, Value: "+"},
			{Kind: token.IntLit, Value: "2"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "sub"},
			{Kind: token.KWInt, Value: "int"},
			{Kind: token.MinusEq, Value: "-="},
			{Kind: token.IntLit, Value: "1"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "plus"},
			{Kind: token.KWInt, Value: "int"},
			{Kind: token.PlusEq, Value: "+="},
			{Kind: token.IntLit, Value: "1"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "multiply"},
			{Kind: token.KWInt, Value: "int"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "2"},
			{Kind: token.Star, Value: "*"},
			{Kind: token.IntLit, Value: "2"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "substract"},
			{Kind: token.KWInt, Value: "int"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "2"},
			{Kind: token.Minus, Value: "-"},
			{Kind: token.IntLit, Value: "2"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "divide"},
			{Kind: token.KWInt, Value: "int"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "2"},
			{Kind: token.Slash, Value: "/"},
			{Kind: token.IntLit, Value: "2"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "modulo"},
			{Kind: token.KWInt, Value: "int"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "2"},
			{Kind: token.Modulo, Value: "%"},
			{Kind: token.IntLit, Value: "2"},

			{Kind: token.RBrace, Value: "}"},
		}
		lex := New(input)
		lex.Tokenize()
		for i, r := range result {
			assert.Equal(r.Kind, lex.Tokens[i].Kind)
			assert.Equal(r.Value, lex.Tokens[i].Value)
		}
		assert.Equal(len(result), len(lex.Tokens))
	})

	t.Run("bool", func(t *testing.T) {
		input := `package main

func main() {
  var i bool = false
  // new test
  var j bool = true // test
  var a bool = 1 > 1
  var b bool = 1 >= 1
  var c bool = 1 < 1
  var d bool = 1 <= 1
  var e bool = 1 == 1
  var f bool = 1 != 1
}
`
		result := []token.Token{
			{Kind: token.KWPackage, Value: "package"},
			{Kind: token.Ident, Value: "main"},
			{Kind: token.KWFunc, Value: "func"},
			{Kind: token.Ident, Value: "main"},
			{Kind: token.LParen, Value: "("},
			{Kind: token.RParen, Value: ")"},
			{Kind: token.LBrace, Value: "{"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "i"},
			{Kind: token.KWBool, Value: "bool"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.BoolLit, Value: "false"},

			{Kind: token.Comment, Value: "// new test"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "j"},
			{Kind: token.KWBool, Value: "bool"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.BoolLit, Value: "true"},
			{Kind: token.Comment, Value: "// test"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "a"},
			{Kind: token.KWBool, Value: "bool"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "1"},
			{Kind: token.Gt, Value: ">"},
			{Kind: token.IntLit, Value: "1"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "b"},
			{Kind: token.KWBool, Value: "bool"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "1"},
			{Kind: token.Gte, Value: ">="},
			{Kind: token.IntLit, Value: "1"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "c"},
			{Kind: token.KWBool, Value: "bool"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "1"},
			{Kind: token.Lt, Value: "<"},
			{Kind: token.IntLit, Value: "1"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "d"},
			{Kind: token.KWBool, Value: "bool"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "1"},
			{Kind: token.Lte, Value: "<="},
			{Kind: token.IntLit, Value: "1"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "e"},
			{Kind: token.KWBool, Value: "bool"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "1"},
			{Kind: token.Eq, Value: "=="},
			{Kind: token.IntLit, Value: "1"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "f"},
			{Kind: token.KWBool, Value: "bool"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "1"},
			{Kind: token.Neq, Value: "!="},
			{Kind: token.IntLit, Value: "1"},

			{Kind: token.RBrace, Value: "}"},
		}
		lex := New(input)
		lex.Tokenize()
		for i, r := range result {
			assert.Equal(r.Kind, lex.Tokens[i].Kind)
			assert.Equal(r.Value, lex.Tokens[i].Value)
		}
		assert.Equal(len(result), len(lex.Tokens))
	})

	t.Run("illegal_vars", func(t *testing.T) {
		input := `package main

func main() {
  var pi1 float64 = 3.141.592,653,59
  var pi2 float64 = 3.141.
  var pi3 float64 = 3.141,
  var x string = "test
}

`
		result := []token.Token{
			{Kind: token.KWPackage, Value: "package"},
			{Kind: token.Ident, Value: "main"},
			{Kind: token.KWFunc, Value: "func"},
			{Kind: token.Ident, Value: "main"},
			{Kind: token.LParen, Value: "("},
			{Kind: token.RParen, Value: ")"},
			{Kind: token.LBrace, Value: "{"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "pi1"},
			{Kind: token.KWFloat64, Value: "float64"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.Illegal, Value: "3.141.592,653,59"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "pi2"},
			{Kind: token.KWFloat64, Value: "float64"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.Illegal, Value: "3.141."},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "pi3"},
			{Kind: token.KWFloat64, Value: "float64"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.Illegal, Value: "3.141,"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "x"},
			{Kind: token.KWString, Value: "string"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.Illegal, Value: `"test`},

			{Kind: token.RBrace, Value: "}"},
		}
		lex := New(input)
		lex.Tokenize()
		for i, r := range result {
			assert.Equal(r.Kind, lex.Tokens[i].Kind)
			assert.Equal(r.Value, lex.Tokens[i].Value)
		}
		assert.Equal(len(result), len(lex.Tokens))
	})

	t.Run("comment", func(t *testing.T) {
		input := `package main

func main() {
  var pi2 float64 = 3.14 // comment
  // new test
  var j bool = true // test
}
`
		result := []token.Token{
			{Kind: token.KWPackage, Value: "package"},
			{Kind: token.Ident, Value: "main"},
			{Kind: token.KWFunc, Value: "func"},
			{Kind: token.Ident, Value: "main"},
			{Kind: token.LParen, Value: "("},
			{Kind: token.RParen, Value: ")"},
			{Kind: token.LBrace, Value: "{"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "pi2"},
			{Kind: token.KWFloat64, Value: "float64"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.FloatLit, Value: "3.14"},
			{Kind: token.Comment, Value: "// comment"},

			{Kind: token.Comment, Value: "// new test"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "j"},
			{Kind: token.KWBool, Value: "bool"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.BoolLit, Value: "true"},
			{Kind: token.Comment, Value: "// test"},

			{Kind: token.RBrace, Value: "}"},
		}
		lex := New(input)
		lex.Tokenize()
		for i, r := range result {
			assert.Equal(r.Kind, lex.Tokens[i].Kind)
			assert.Equal(r.Value, lex.Tokens[i].Value)
		}
		assert.Equal(len(result), len(lex.Tokens))
	})
}
