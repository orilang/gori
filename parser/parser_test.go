package parser

import (
	"syscall"
	"testing"

	"github.com/orilang/gori/ast"
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
}

func TestParser_parse_file(t *testing.T) {
	assert := assert.New(t)

	t.Run("error_global_var_forbidden", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "a"},
			{Kind: token.KWInt, Value: "int"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "0"},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 10},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 13},
			{Kind: token.RBrace, Value: "}", Line: 8, Column: 1},
			{Kind: token.EOF, Value: "", Line: 9, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("function", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},

			{Kind: token.KWConst, Value: "const"},
			{Kind: token.Ident, Value: "a"},
			{Kind: token.KWFloat, Value: "float"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.FloatLit, Value: "3.14"},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 10},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 13},

			{Kind: token.KWConst, Value: "const"},
			{Kind: token.Ident, Value: "ab"},
			{Kind: token.KWFloat, Value: "float"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.FloatLit, Value: "3.14"},

			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "a"},
			{Kind: token.KWInt, Value: "int"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "0"},

			{Kind: token.RBrace, Value: "}", Line: 8, Column: 1},
			{Kind: token.EOF, Value: "", Line: 9, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 ConstDecls
  ConstDecl
   Const: "const" @0:0 (kind=23)
   Name: "a" @0:0 (kind=3)
   Type
    NameType
     Name: "float" @0:0 (kind=20)
   Eq: "=" @0:0 (kind=49)
   Init
    FloatLitExpr
     Value: "3.14" @0:0 (kind=5)
 Decls
  FuncDecl
   Function: "func" @3:1 (kind=10)
   Name: "main" @3:6 (kind=3)
   Params
    (none)
   Body
    BlockStmt
     LBrace: "{" @3:13 (kind=41)
     Stmts
      ConstDecl
       Const: "const" @0:0 (kind=23)
       Name: "ab" @0:0 (kind=3)
       Type
        NameType
         Name: "float" @0:0 (kind=20)
       Eq: "=" @0:0 (kind=49)
       Init
        FloatLitExpr
         Value: "3.14" @0:0 (kind=5)
      VarDeclStmt
       Var: "var" @0:0 (kind=11)
       Name: "a" @0:0 (kind=3)
       Type
        NameType
         Name: "int" @0:0 (kind=12)
       Eq: "=" @0:0 (kind=49)
       Init
        IntLitExpr
         Value: "0" @0:0 (kind=4)
     RBrace: "}" @8:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("decls_none", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 9},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:9 (kind=3)
 Decls
  (none)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("function_params", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "dummy", Line: 1, Column: 9},

			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 1},
			{Kind: token.Ident, Value: "x", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 7},

			{Kind: token.Ident, Value: "a", Line: 3, Column: 8},
			{Kind: token.KWInt, Value: "int", Line: 3, Column: 10},

			{Kind: token.Comma, Value: ",", Line: 3, Column: 10},

			{Kind: token.Ident, Value: "b", Line: 3, Column: 15},
			{Kind: token.KWString, Value: "string", Line: 3, Column: 17},

			{Kind: token.RParen, Value: ")", Line: 3, Column: 23},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 24},

			{Kind: token.RBrace, Value: "}", Line: 3, Column: 25},
			{Kind: token.EOF, Value: "", Line: 4, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "dummy" @1:9 (kind=3)
 Decls
  FuncDecl
   Function: "func" @3:1 (kind=10)
   Name: "x" @3:6 (kind=3)
   Params
    Param
     Function: "a" @3:8 (kind=3)
     Type
      NameType
       Name: "int" @3:10 (kind=12)
    Param
     Function: "b" @3:15 (kind=3)
     Type
      NameType
       Name: "string" @3:17 (kind=24)
   Body
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})
}

func TestParser_parse_const(t *testing.T) {
	assert := assert.New(t)
	t.Run("float_lit", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWConst, Value: "const"},
			{Kind: token.Ident, Value: "a"},
			{Kind: token.KWFloat, Value: "float"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.FloatLit, Value: "3.14"},
		}

		parser := New(input)
		pr := parser.parseConstDecl()
		result := `ConstDecl
 Const: "const" @0:0 (kind=23)
 Name: "a" @0:0 (kind=3)
 Type
  NameType
   Name: "float" @0:0 (kind=20)
 Eq: "=" @0:0 (kind=49)
 Init
  FloatLitExpr
   Value: "3.14" @0:0 (kind=5)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})
}
func TestParser_parse_var(t *testing.T) {
	assert := assert.New(t)

	t.Run("int_lit", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "a"},
			{Kind: token.KWInt, Value: "int"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.IntLit, Value: "0"},
		}

		parser := New(input)
		pr := parser.parseVarDecl()
		result := `VarDeclStmt
 Var: "var" @0:0 (kind=11)
 Name: "a" @0:0 (kind=3)
 Type
  NameType
   Name: "int" @0:0 (kind=12)
 Eq: "=" @0:0 (kind=49)
 Init
  IntLitExpr
   Value: "0" @0:0 (kind=4)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("float_lit", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "a"},
			{Kind: token.KWFloat, Value: "float"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.FloatLit, Value: "3.14"},
		}

		parser := New(input)
		pr := parser.parseVarDecl()
		result := `VarDeclStmt
 Var: "var" @0:0 (kind=11)
 Name: "a" @0:0 (kind=3)
 Type
  NameType
   Name: "float" @0:0 (kind=20)
 Eq: "=" @0:0 (kind=49)
 Init
  FloatLitExpr
   Value: "3.14" @0:0 (kind=5)
`
		assert.Equal(result, ast.Dump(pr))
	})

	t.Run("bool_lit", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "a"},
			{Kind: token.KWBool, Value: "bool"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.BoolLit, Value: "true"},
		}

		parser := New(input)
		pr := parser.parseVarDecl()
		result := `VarDeclStmt
 Var: "var" @0:0 (kind=11)
 Name: "a" @0:0 (kind=3)
 Type
  NameType
   Name: "bool" @0:0 (kind=25)
 Eq: "=" @0:0 (kind=49)
 Init
  BoolLitExpr
   Value: "true" @0:0 (kind=7)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("string_lit", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "a"},
			{Kind: token.KWString, Value: "string"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.StringLit, Value: "ok"},
		}

		parser := New(input)
		pr := parser.parseVarDecl()
		result := `VarDeclStmt
 Var: "var" @0:0 (kind=11)
 Name: "a" @0:0 (kind=3)
 Type
  NameType
   Name: "string" @0:0 (kind=24)
 Eq: "=" @0:0 (kind=49)
 Init
  StringLitExpr
   Value: "ok" @0:0 (kind=6)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("indent_lit", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWVar, Value: "var"},
			{Kind: token.Ident, Value: "a"},
			{Kind: token.KWFloat, Value: "float"},
			{Kind: token.Assign, Value: "="},
			{Kind: token.Ident, Value: "x"},
		}

		parser := New(input)
		pr := parser.parseVarDecl()
		result := `VarDeclStmt
 Var: "var" @0:0 (kind=11)
 Name: "a" @0:0 (kind=3)
 Type
  NameType
   Name: "float" @0:0 (kind=20)
 Eq: "=" @0:0 (kind=49)
 Init
  IdentExpr
   Name: "x" @0:0 (kind=3)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})
}

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
			{Kind: token.IntLit, Value: "!"},
		}

		parser := New(input)
		pr := parser.parseExpr(LOWEST)
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}

func TestParser_expr(t *testing.T) {
	assert := assert.New(t)
	t.Run("grouping", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.LParen, Value: "(", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "a", Line: 1, Column: 2},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 3},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
		pr := parser.parseExpr(LOWEST)
		result := `ParenExpr
 IdentExpr
  Name: "a" @1:2 (kind=3)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("error_unclosed", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.LParen, Value: "(", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "a", Line: 1, Column: 2},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
		pr := parser.parseExpr(LOWEST)
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("error_empty", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.LParen, Value: "(", Line: 1, Column: 1},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 2},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
		pr := parser.parseExpr(LOWEST)
		assert.Contains(ast.Dump(pr), "BadExpr")
		assert.Contains(ast.Dump(pr), "expected expression inside parentheses")
		assert.Equal(0, len(parser.errors))
	})

	t.Run("grouping", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.LParen, Value: "(", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "a", Line: 1, Column: 2},
			{Kind: token.Ident, Value: "+", Line: 1, Column: 3},
			{Kind: token.Ident, Value: "b", Line: 1, Column: 4},
			{Kind: token.RParen, Value: ")", Line: 1, Column: 5},
			{Kind: token.EOF, Value: "", Line: 2, Column: 1},
		}

		parser := New(input)
		pr := parser.parseExpr(LOWEST)
		result := `ParenExpr
 IdentExpr
  Name: "b" @1:4 (kind=3)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})
}
