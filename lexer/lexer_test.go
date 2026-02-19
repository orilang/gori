package lexer

import (
	"os"
	"path/filepath"
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
		workingDir, err := os.Getwd()
		assert.Nil(err)

		lex, err := NewLexer(Config{Directory: filepath.Join(workingDir, "..", "testdata/success")})
		assert.Nil(err)
		assert.Nil(lex.StartLexing())
	})

	t.Run("illegal", func(t *testing.T) {
		workingDir, err := os.Getwd()
		assert.Nil(err)

		lex, err := NewLexer(Config{Directory: filepath.Join(workingDir, "..", "testdata/illegal")})
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
			{Kind: token.EOF, Value: ""},
		}
		lex := New([]byte(input))
		lex.Tokenize()
		for i, r := range result {
			assert.Equal(r.Kind, lex.Tokens[i].Kind)
			assert.Equal(r.Value, lex.Tokens[i].Value)
		}
		assert.Equal(len(result), len(lex.Tokens))
	})

	t.Run("line_column", func(t *testing.T) {
		input := `package main

func main() {
// comment
/*
multi line
*/
}
`
		result := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},
			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 10},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 13},
			{Kind: token.Comment, Value: "// comment", Line: 4, Column: 1},
			{Kind: token.Comment, Value: `/*
multi line
*/`, Line: 5, Column: 1},
			{Kind: token.RBrace, Value: "}", Line: 8, Column: 1},
			{Kind: token.EOF, Value: "", Line: 9, Column: 1},
		}
		lex := New([]byte(input))
		lex.Tokenize()
		assert.Equal(result, lex.Tokens)
	})

	t.Run("vars", func(t *testing.T) {
		input := `package main

func main() {
  var a int = 0
	a++
	a--
	var b uint = 0
  var c int8 = 0
  var d int32 = 0
  var e int64 = 120

  var f float = 0
  var g float32 = 0
  var h float64 = 0
  var pi float64 = 3.14
  var pi2 float64 = 3.14 // comment
  var pi3 float64 = 3.141_592_653_59
  var pix float64 = .14
  var piy float64 = 3_14

  var i bool = false
  // new test
  var j bool = true // test
  var bt string = "true"
  var bf string = "false"
	x:=[]int{1}
	xx:=[x:]
	x1 := 1;x2:=1
	x3.z=1
	backslash:="a\"b"
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

			{Kind: token.Ident, Value: "a"},
			{Kind: token.PPlus, Value: "++"},

			{Kind: token.Ident, Value: "a"},
			{Kind: token.MMinus, Value: "--"},

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
			{Kind: token.FloatLit, Value: "3.141_592_653_59"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "pix"},
			{Kind: token.KWFloat64, Value: "float64"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.FloatLit, Value: ".14"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "piy"},
			{Kind: token.KWFloat64, Value: "float64"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "3_14"},

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

			{Kind: token.Ident, Value: "x"},
			{Kind: token.Define, Value: ":="},
			{Kind: token.LBracket, Value: "["},
			{Kind: token.RBracket, Value: "]"},
			{Kind: token.KWInt, Value: "int"},
			{Kind: token.LBrace, Value: "{"},
			{Kind: token.IntLit, Value: "1"},
			{Kind: token.RBrace, Value: "}"},

			{Kind: token.Ident, Value: "xx"},
			{Kind: token.Define, Value: ":="},
			{Kind: token.LBracket, Value: "["},
			{Kind: token.Ident, Value: "x"},
			{Kind: token.Colon, Value: ":"},
			{Kind: token.RBracket, Value: "]"},

			{Kind: token.Ident, Value: "x1"},
			{Kind: token.Define, Value: ":="},
			{Kind: token.IntLit, Value: "1"},
			{Kind: token.SemiComma, Value: ";"},
			{Kind: token.Ident, Value: "x2"},
			{Kind: token.Define, Value: ":="},
			{Kind: token.IntLit, Value: "1"},

			{Kind: token.Ident, Value: "x3"},
			{Kind: token.Dot, Value: "."},
			{Kind: token.Ident, Value: "z"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "1"},

			{Kind: token.Ident, Value: "backslash"},
			{Kind: token.Define, Value: ":="},
			{Kind: token.StringLit, Value: `"a\"b"`},

			{Kind: token.RBrace, Value: "}"},
			{Kind: token.EOF, Value: ""},
		}
		lex := New([]byte(input))
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
  var starEq int *=2
  var multiply int = 2 * 2
  var substract int = 2 - 2
  var divide int = 2 / 2
  var divide2 int /= 2
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
			{Kind: token.Ident, Value: "starEq"},
			{Kind: token.KWInt, Value: "int"},
			{Kind: token.StarEq, Value: "*="},
			{Kind: token.IntLit, Value: "2"},

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
			{Kind: token.Ident, Value: "divide2"},
			{Kind: token.KWInt, Value: "int"},
			{Kind: token.SlashEq, Value: "/="},
			{Kind: token.IntLit, Value: "2"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "modulo"},
			{Kind: token.KWInt, Value: "int"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "2"},
			{Kind: token.Modulo, Value: "%"},
			{Kind: token.IntLit, Value: "2"},

			{Kind: token.RBrace, Value: "}"},
			{Kind: token.EOF, Value: ""},
		}
		lex := New([]byte(input))
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
  var g bool = e || f
  var h bool = e && f
  var i bool=!e
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

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "g"},
			{Kind: token.KWBool, Value: "bool"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.Ident, Value: "e"},
			{Kind: token.Or, Value: "||"},
			{Kind: token.Ident, Value: "f"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "h"},
			{Kind: token.KWBool, Value: "bool"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.Ident, Value: "e"},
			{Kind: token.And, Value: "&&"},
			{Kind: token.Ident, Value: "f"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "i"},
			{Kind: token.KWBool, Value: "bool"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.Not, Value: "!"},
			{Kind: token.Ident, Value: "e"},

			{Kind: token.RBrace, Value: "}"},
			{Kind: token.EOF, Value: ""},
		}
		lex := New([]byte(input))
		lex.Tokenize()
		for i, r := range result {
			assert.Equal(r.Kind, lex.Tokens[i].Kind)
			assert.Equal(r.Value, lex.Tokens[i].Value)
		}
		assert.Equal(len(result), len(lex.Tokens))
	})

	t.Run("illegal_numbers", func(t *testing.T) {
		input := `package main

func main() {
  var pi1 float64 = 3.141.592_653_59
  var pi2 float64 = 3.141.
  var pi3 float64 = 3.141_
  var pi4 float64 = _.14
  var pi5 float64=3._14
  var pi6 float64 = 3_.14
  var pi7 float64 = 3__141
  var pi8 float64 = _3_141
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
			{Kind: token.Illegal, Value: "3.141.592_653_59"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "pi2"},
			{Kind: token.KWFloat64, Value: "float64"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.Illegal, Value: "3.141."},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "pi3"},
			{Kind: token.KWFloat64, Value: "float64"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.Illegal, Value: "3.141_"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "pi4"},
			{Kind: token.KWFloat64, Value: "float64"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.Illegal, Value: "_.14"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "pi5"},
			{Kind: token.KWFloat64, Value: "float64"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.Illegal, Value: "3._14"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "pi6"},
			{Kind: token.KWFloat64, Value: "float64"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.Illegal, Value: "3_.14"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "pi7"},
			{Kind: token.KWFloat64, Value: "float64"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.Illegal, Value: "3__141"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "pi8"},
			{Kind: token.KWFloat64, Value: "float64"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.Illegal, Value: "_3_141"},

			{Kind: token.RBrace, Value: "}"},
			{Kind: token.EOF, Value: ""},
		}
		lex := New([]byte(input))
		lex.Tokenize()
		for i, r := range result {
			assert.Equal(r.Kind, lex.Tokens[i].Kind, i)
			assert.Equal(r.Value, lex.Tokens[i].Value, i)
		}
		assert.Equal(len(result), len(lex.Tokens))
	})

	t.Run("illegal_string", func(t *testing.T) {
		input := `package main

func main() {
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
			{Kind: token.Ident, Value: "x"},
			{Kind: token.KWString, Value: "string"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.Illegal, Value: "\"test\n}\n\n"},

			{Kind: token.EOF, Value: ""},
		}
		lex := New([]byte(input))
		lex.Tokenize()
		for i, r := range result {
			assert.Equal(r.Kind, lex.Tokens[i].Kind, i)
			assert.Equal(r.Value, lex.Tokens[i].Value)
		}
		assert.Equal(len(result), len(lex.Tokens))
	})

	t.Run("illegal_characters", func(t *testing.T) {
		input := `package main

func main() {
  var a int = 1 & 1
  var b int = 1 | 1
	_ := 1
	_c := 1
	#
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
			{Kind: token.IntLit, Value: "1"},
			{Kind: token.Illegal, Value: "&"},
			{Kind: token.IntLit, Value: "1"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "b"},
			{Kind: token.KWInt, Value: "int"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "1"},
			{Kind: token.Pipe, Value: "|"},
			{Kind: token.IntLit, Value: "1"},

			// _ here is only temporary
			// in future tests it won't be
			{Kind: token.Illegal, Value: "_"},
			{Kind: token.Define, Value: ":="},
			{Kind: token.IntLit, Value: "1"},

			{Kind: token.Illegal, Value: "_c"},
			{Kind: token.Define, Value: ":="},
			{Kind: token.IntLit, Value: "1"},

			{Kind: token.Illegal, Value: "#"},

			{Kind: token.RBrace, Value: "}"},
			{Kind: token.EOF, Value: ""},
		}
		lex := New([]byte(input))
		lex.Tokenize()
		for i, r := range result {
			assert.Equal(r.Kind, lex.Tokens[i].Kind, i)
			assert.Equal(r.Value, lex.Tokens[i].Value)
		}
		assert.Equal(len(result), len(lex.Tokens))
	})

	t.Run("illegal_multiline_comment", func(t *testing.T) {
		input := `package main

func main() {
  /*
  var pi1 float64 = 3.141.592_653_59
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

			{Kind: token.Illegal, Value: `/*
  var pi1 float64 = 3.141.592_653_59
}

`},

			{Kind: token.EOF, Value: ""},
		}
		lex := New([]byte(input))
		lex.Tokenize()
		for i, r := range result {
			assert.Equal(r.Kind, lex.Tokens[i].Kind, i)
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
  var k bool = true /* k bool */
	/* a * b */
	/* a /* b */ c */
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

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "k"},
			{Kind: token.KWBool, Value: "bool"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.BoolLit, Value: "true"},
			{Kind: token.Comment, Value: "/* k bool */"},

			{Kind: token.Comment, Value: "/* a * b */"},
			{Kind: token.Comment, Value: "/* a /* b */"},

			{Kind: token.Ident, Value: "c"},
			{Kind: token.Star, Value: "*"},
			{Kind: token.Slash, Value: "/"},

			{Kind: token.RBrace, Value: "}"},
			{Kind: token.EOF, Value: ""},
		}
		lex := New([]byte(input))
		lex.Tokenize()
		for i, r := range result {
			assert.Equal(r.Kind, lex.Tokens[i].Kind, i)
			assert.Equal(r.Value, lex.Tokens[i].Value, i)
		}
		assert.Equal(len(result), len(lex.Tokens))
	})

	t.Run("fetch_next_token", func(t *testing.T) {
		input := "main"

		lex := New([]byte(input))
		lex.position = 4
		_, ok := lex.fetchNextToken()

		assert.Equal(false, ok)
	})
}
