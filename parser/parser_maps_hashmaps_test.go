package parser

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
	"github.com/stretchr/testify/assert"
)

func TestParser_maps_hashmaps(t *testing.T) {
	assert := assert.New(t)

	t.Run("map", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 1},
			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 9},
			{Kind: token.Ident, Value: "main", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 10},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 12},
			{Kind: token.KWVar, Value: "var", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "x", Line: 4, Column: 7},
			{Kind: token.KWMap, Value: "map", Line: 4, Column: 9},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 12},
			{Kind: token.KWString, Value: "string", Line: 4, Column: 13},
			{Kind: token.RBracket, Value: "]", Line: 4, Column: 19},
			{Kind: token.KWString, Value: "string", Line: 4, Column: 20},
			{Kind: token.Assign, Value: "=", Line: 4, Column: 27},
			{Kind: token.Ident, Value: "make", Line: 4, Column: 29},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 33},
			{Kind: token.KWMap, Value: "map", Line: 4, Column: 34},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 37},
			{Kind: token.KWString, Value: "string", Line: 4, Column: 38},
			{Kind: token.RBracket, Value: "]", Line: 4, Column: 44},
			{Kind: token.KWString, Value: "string", Line: 4, Column: 45},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 51},

			{Kind: token.KWVar, Value: "var", Line: 5, Column: 3},
			{Kind: token.Ident, Value: "y", Line: 5, Column: 7},
			{Kind: token.KWHashMap, Value: "hashmap", Line: 5, Column: 9},
			{Kind: token.LBracket, Value: "[", Line: 5, Column: 16},
			{Kind: token.KWString, Value: "string", Line: 5, Column: 17},
			{Kind: token.RBracket, Value: "]", Line: 5, Column: 23},
			{Kind: token.KWString, Value: "string", Line: 5, Column: 24},
			{Kind: token.Assign, Value: "=", Line: 5, Column: 31},
			{Kind: token.Ident, Value: "make", Line: 5, Column: 29},
			{Kind: token.LParen, Value: "(", Line: 5, Column: 33},
			{Kind: token.KWHashMap, Value: "hashmap", Line: 5, Column: 38},
			{Kind: token.LBracket, Value: "[", Line: 5, Column: 37},
			{Kind: token.KWString, Value: "string", Line: 5, Column: 38},
			{Kind: token.RBracket, Value: "]", Line: 5, Column: 44},
			{Kind: token.KWString, Value: "string", Line: 5, Column: 45},
			{Kind: token.RParen, Value: ")", Line: 5, Column: 51},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 1},
			{Kind: token.EOF, Value: "", Line: 7, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		result := `File
 Package: "package" @1:1 (kind=8)
 Name: "main" @1:1 (kind=3)
 Decls
  FuncDecl
   Function: "func" @3:9 (kind=10)
   Name: "main" @3:6 (kind=3)
   Params
    (none)
   Body
    BlockStmt
     LBrace: "{" @3:12 (kind=41)
     Stmts
      VarDeclStmt
       Var: "var" @4:3 (kind=11)
       Name: "x" @4:7 (kind=3)
       Type
        Map: "map" @4:9 (kind=79)
        LBracket: "[" @4:12 (kind=43)
        KeyType:
         Name: "string" @4:13 (kind=24)
        RBracket: "]" @4:19 (kind=44)
        ValueType:
         Name: "string" @4:13 (kind=24)
       Eq: "=" @4:27 (kind=49)
       Init
        Make: "make" @4:29 (kind=3)
        LParen: "(" @4:33 (kind=39)
        Map: "map" @4:34 (kind=79)
        LBracket: "[" @4:37 (kind=43)
        KeyType:
         Name: "string" @4:38 (kind=24)
        RBracket: "]" @4:44 (kind=44)
        ValueType:
         Name: "string" @4:38 (kind=24)
        RParen: ")" @4:51 (kind=40)
      VarDeclStmt
       Var: "var" @5:3 (kind=11)
       Name: "y" @5:7 (kind=3)
       Type
        Hashmap: "hashmap" @5:9 (kind=80)
        LBracket: "[" @5:16 (kind=43)
        KeyType:
         Name: "string" @5:17 (kind=24)
        RBracket: "]" @5:23 (kind=44)
        ValueType:
         Name: "string" @5:17 (kind=24)
       Eq: "=" @5:31 (kind=49)
       Init
        Make: "make" @5:29 (kind=3)
        LParen: "(" @5:33 (kind=39)
        Hashmap: "hashmap" @5:38 (kind=80)
        LBracket: "[" @5:37 (kind=43)
        KeyType:
         Name: "string" @5:38 (kind=24)
        RBracket: "]" @5:44 (kind=44)
        ValueType:
         Name: "string" @5:38 (kind=24)
        RParen: ")" @5:51 (kind=40)
     RBrace: "}" @6:1 (kind=42)
`
		assert.Equal(result, ast.Dump(pr))
		assert.Equal(0, len(parser.errors))
	})

	t.Run("bad_x1", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 1},
			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 9},
			{Kind: token.Ident, Value: "main", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 10},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 12},
			{Kind: token.KWVar, Value: "var", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "x", Line: 4, Column: 7},
			{Kind: token.KWMap, Value: "map", Line: 4, Column: 9},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 12},
			{Kind: token.KWStruct, Value: "struct", Line: 4, Column: 13},
			{Kind: token.RBracket, Value: "]", Line: 4, Column: 19},
			{Kind: token.KWString, Value: "string", Line: 4, Column: 20},
			{Kind: token.Assign, Value: "=", Line: 4, Column: 27},
			{Kind: token.Ident, Value: "make", Line: 4, Column: 29},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 33},
			{Kind: token.KWMap, Value: "map", Line: 4, Column: 34},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 37},
			{Kind: token.KWStruct, Value: "struct", Line: 4, Column: 38},
			{Kind: token.RBracket, Value: "]", Line: 4, Column: 44},
			{Kind: token.KWString, Value: "string", Line: 4, Column: 45},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 51},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 1},
			{Kind: token.EOF, Value: "", Line: 7, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})

	t.Run("bad_x2", func(t *testing.T) {
		input := []token.Token{
			{Kind: token.KWPackage, Value: "package", Line: 1, Column: 1},
			{Kind: token.Ident, Value: "main", Line: 1, Column: 1},
			{Kind: token.KWFunc, Value: "func", Line: 3, Column: 9},
			{Kind: token.Ident, Value: "main", Line: 3, Column: 6},
			{Kind: token.LParen, Value: "(", Line: 3, Column: 10},
			{Kind: token.RParen, Value: ")", Line: 3, Column: 11},
			{Kind: token.LBrace, Value: "{", Line: 3, Column: 12},
			{Kind: token.KWVar, Value: "var", Line: 4, Column: 3},
			{Kind: token.Ident, Value: "x", Line: 4, Column: 7},
			{Kind: token.KWMap, Value: "map", Line: 4, Column: 9},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 12},
			{Kind: token.KWString, Value: "string", Line: 4, Column: 13},
			{Kind: token.RBracket, Value: "]", Line: 4, Column: 19},
			{Kind: token.KWStruct, Value: "struct", Line: 4, Column: 20},
			{Kind: token.Assign, Value: "=", Line: 4, Column: 27},
			{Kind: token.Ident, Value: "make", Line: 4, Column: 29},
			{Kind: token.LParen, Value: "(", Line: 4, Column: 33},
			{Kind: token.KWMap, Value: "map", Line: 4, Column: 34},
			{Kind: token.LBracket, Value: "[", Line: 4, Column: 37},
			{Kind: token.KWStruct, Value: "struct", Line: 4, Column: 38},
			{Kind: token.RBracket, Value: "]", Line: 4, Column: 44},
			{Kind: token.KWString, Value: "string", Line: 4, Column: 45},
			{Kind: token.RParen, Value: ")", Line: 4, Column: 51},
			{Kind: token.RBrace, Value: "}", Line: 6, Column: 1},
			{Kind: token.EOF, Value: "", Line: 7, Column: 1},
		}

		parser := New(input)
		pr := parser.ParseFile()
		assert.NotNil(pr)
		assert.Greater(len(parser.errors), 0)
	})
}
