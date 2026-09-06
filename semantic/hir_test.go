package semantic

import (
	"fmt"
	"testing"

	"github.com/orilang/gori/lexer"
	"github.com/orilang/gori/parser"
	"github.com/orilang/gori/token"
	"github.com/stretchr/testify/require"
)

func TestSemantic_hir(t *testing.T) {
	t.Run("x1", func(t *testing.T) {
		data := `package main
const a int = int(0)
`

		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		require.Equal(t, 0, len(parser.Errors))
		check := NewChecker()

		program, diagnostics := check.Check(pr)
		require.Equal(t, 0, len(diagnostics))
		require.Equal(t, 1, len(program.Files))

		pf := program.Files[0]
		require.Equal(t, 1, len(pf.Decls))
		target, ok := pf.Decls[0].(*ConstDecl)
		require.Equal(t, true, ok)
		require.Equal(t, "a", target.Name)
		require.Equal(t, TInt, target.Symbol.Type)
		ce, ok := target.Init.(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, ce.To)
		lit, ok := ce.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, lit.Type)
		require.Equal(t, "0", lit.Value)
	})

	t.Run("x2", func(t *testing.T) {
		data := `package main
const a string = "a"
`

		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		require.Equal(t, 0, len(parser.Errors))
		check := NewChecker()

		program, diagnostics := check.Check(pr)
		require.Equal(t, 0, len(diagnostics))
		require.Equal(t, 1, len(program.Files))

		pf := program.Files[0]
		require.Equal(t, 1, len(pf.Decls))
		target, ok := pf.Decls[0].(*ConstDecl)
		require.Equal(t, true, ok)
		require.Equal(t, "a", target.Name)
		require.NotNil(t, target.Symbol)
		require.Equal(t, target.Symbol.Type, TString)
		lit, ok := target.Init.(*StringLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TString, lit.Type)
		require.Equal(t, fmt.Sprintf("%q", "a"), lit.Value)
	})

	t.Run("x3", func(t *testing.T) {
		data := `package main
const a bool = true
`

		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		for _, v := range parser.Errors {
			fmt.Println(v.Error())
		}
		require.Equal(t, 0, len(parser.Errors))
		check := NewChecker()

		program, diagnostics := check.Check(pr)
		require.Equal(t, 0, len(diagnostics))
		require.Equal(t, 1, len(program.Files))

		pf := program.Files[0]
		require.Equal(t, 1, len(pf.Decls))
		target, ok := pf.Decls[0].(*ConstDecl)
		require.Equal(t, true, ok)
		require.Equal(t, "a", target.Name)
		require.NotNil(t, target.Symbol)
		require.Equal(t, TBool, target.Symbol.Type)
		lit, ok := target.Init.(*BoolLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TBool, lit.Type)
		require.Equal(t, "true", lit.Value)
	})

	t.Run("x4", func(t *testing.T) {
		data := `package main
const a float64 = float64(0)
`

		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		require.Equal(t, 0, len(parser.Errors))
		check := NewChecker()

		program, diagnostics := check.Check(pr)
		require.Equal(t, 0, len(diagnostics))
		require.Equal(t, 1, len(program.Files))

		pf := program.Files[0]
		require.Equal(t, 1, len(pf.Decls))
		target, ok := pf.Decls[0].(*ConstDecl)
		require.Equal(t, true, ok)
		require.Equal(t, "a", target.Name)
		require.NotNil(t, target.Symbol)
		require.Equal(t, TFloat64, target.Symbol.Type)
		ce, ok := target.Init.(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TFloat64, ce.To)
		lit, ok := ce.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, lit.Type)
		require.Equal(t, "0", lit.Value)
	})

	t.Run("x5", func(t *testing.T) {
		data := `package main
const a int = int(0)
`

		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		require.Equal(t, 0, len(parser.Errors))
		check := NewChecker()

		program, diagnostics := check.Check(pr)
		require.Equal(t, 0, len(diagnostics))
		require.Equal(t, 1, len(program.Files))

		pf := program.Files[0]
		require.Equal(t, 1, len(pf.Decls))
		target, ok := pf.Decls[0].(*ConstDecl)
		require.Equal(t, true, ok)
		require.Equal(t, "a", target.Name)
		require.NotNil(t, target.Symbol)
		require.Equal(t, target.Symbol.Type, TInt)
		ce, ok := target.Init.(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, ce.To)
		lit, ok := ce.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, lit.Type)
		require.Equal(t, "0", lit.Value)
	})

	t.Run("x6", func(t *testing.T) {
		data := `package main
const a int = int(0)
const b int = a + int(1)
`

		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		require.Equal(t, 0, len(parser.Errors))
		check := NewChecker()

		program, diagnostics := check.Check(pr)
		require.Equal(t, 0, len(diagnostics))
		require.Equal(t, 1, len(program.Files))

		pf := program.Files[0]
		require.Equal(t, 2, len(pf.Decls))

		target, ok := pf.Decls[0].(*ConstDecl)
		require.Equal(t, true, ok)
		require.Equal(t, "a", target.Name)
		require.NotNil(t, target.Symbol)
		require.Equal(t, TInt, target.Symbol.Type)

		ce, ok := target.Init.(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, ce.To)
		lit, ok := ce.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, lit.Type)
		require.Equal(t, "0", lit.Value)

		target1, ok := pf.Decls[1].(*ConstDecl)
		require.Equal(t, true, ok)
		require.Equal(t, "b", target1.Name)
		require.NotNil(t, target.Symbol)
		require.Equal(t, TInt, target1.Symbol.Type)

		init1, ok := target1.Init.(*BinaryExpr)
		require.Equal(t, true, ok)
		require.Equal(t, init1.Type, TInt)

		left, ok := init1.Left.(*IdentExpr)
		require.Equal(t, true, ok)
		require.NotNil(t, left.Symbol)
		require.Equal(t, TInt, left.Type)
		require.Equal(t, "a", left.Value)
		require.Equal(t, token.Plus, init1.Operator)
		ce1, ok := init1.Right.(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, ce1.To)
		right, ok := ce1.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, right.Type)
		require.Equal(t, "1", right.Value)
	})

	t.Run("x7", func(t *testing.T) {
		data := `package main
func main() {
	const a int = int(0)
	const b int = -a
}
`

		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		require.Equal(t, 0, len(parser.Errors))
		check := NewChecker()

		program, diagnostics := check.Check(pr)
		require.Equal(t, 0, len(diagnostics))
		require.Equal(t, 1, len(program.Files))

		pf := program.Files[0]
		require.Equal(t, 1, len(pf.Decls))

		fn, ok := pf.Decls[0].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, fn.Name, "main")
		require.NotNil(t, fn.Symbol)
		require.Equal(t, SymFunc, fn.Symbol.Kind)
		require.Equal(t, 0, len(fn.Params))
		require.Equal(t, 0, len(fn.Results))

		require.NotNil(t, fn.Body)
		require.Equal(t, 2, len(fn.Body.Stmts))
		ds, ok := fn.Body.Stmts[0].(*DeclStmt)
		require.Equal(t, true, ok)
		cs, ok := ds.Decl.(*ConstDecl)
		require.Equal(t, true, ok)
		require.Equal(t, "a", cs.Name)
		require.NotNil(t, cs.Symbol)
		require.Equal(t, TInt, cs.Symbol.Type)

		ce, ok := cs.Init.(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, ce.To)
		lit, ok := ce.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, lit.Type)
		require.Equal(t, "0", lit.Value)

		ds1, ok := fn.Body.Stmts[1].(*DeclStmt)
		require.Equal(t, true, ok)
		cs1, ok := ds1.Decl.(*ConstDecl)
		require.Equal(t, true, ok)
		require.Equal(t, "b", cs1.Name)
		require.NotNil(t, cs1.Symbol)
		require.Equal(t, TInt, cs1.Symbol.Type)

		ce1, ok := cs1.Init.(*UnaryExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, ce1.Type)
		require.Equal(t, token.Minus, ce1.Operator)
		iden, ok := ce1.Right.(*IdentExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, iden.Type)
		require.NotNil(t, iden.Symbol)
		require.Equal(t, "a", iden.Value)
	})

	t.Run("x8", func(t *testing.T) {
		data := `package main
func add(a int, b int) int {
    return a + b
}

func main() {
    a := add(int(1), int(2))
		b := int64(a)
		b = int64(42)
}
`

		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		require.Equal(t, 0, len(parser.Errors))
		check := NewChecker()

		program, diagnostics := check.Check(pr)
		for _, d := range diagnostics {
			fmt.Println(d.Err.Error())
		}
		require.Equal(t, 0, len(diagnostics))
		require.Equal(t, 1, len(program.Files))

		pf := program.Files[0]
		require.Equal(t, 2, len(pf.Decls))

		fn, ok := pf.Decls[0].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, "add", fn.Name)
		require.NotNil(t, fn.Symbol)
		require.Equal(t, SymFunc, fn.Symbol.Kind)
		require.Equal(t, 2, len(fn.Params))
		require.Equal(t, "a", fn.Params[0].Name)
		require.Equal(t, TInt, fn.Params[0].Type)
		require.Equal(t, "b", fn.Params[1].Name)
		require.Equal(t, TInt, fn.Params[1].Type)
		require.Equal(t, 1, len(fn.Results))
		require.Equal(t, TInt, fn.Results[0].Type)

		require.NotNil(t, fn.Body)
		rs, ok := fn.Body.Stmts[0].(*ReturnStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(rs.Values))
		be, ok := rs.Values[0].(*BinaryExpr)
		require.Equal(t, true, ok)

		left, ok := be.Left.(*IdentExpr)
		require.Equal(t, true, ok)
		require.NotNil(t, left.Symbol)
		require.Equal(t, TInt, left.Type)
		require.Equal(t, "a", left.Value)

		require.Equal(t, token.Plus, be.Operator)

		right, ok := be.Right.(*IdentExpr)
		require.Equal(t, true, ok)
		require.NotNil(t, right.Symbol)
		require.Equal(t, TInt, right.Type)
		require.Equal(t, "b", right.Value)

		fn1, ok := pf.Decls[1].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, fn1.Name, "main")
		require.NotNil(t, fn1.Symbol)
		require.Equal(t, SymFunc, fn1.Symbol.Kind)
		require.Equal(t, 0, len(fn1.Params))
		require.Equal(t, 0, len(fn1.Results))

		require.NotNil(t, fn1.Body)
		require.Equal(t, 3, len(fn1.Body.Stmts))

		da, ok := fn1.Body.Stmts[0].(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, "a", da.Symbol.Name)
		require.Equal(t, TInt, da.Symbol.Type)

		dar, ok := da.Right.(*CallExpr)
		require.Equal(t, true, ok)

		dav, ok := dar.Callee.(*IdentExpr)
		require.Equal(t, true, ok)
		require.NotNil(t, dav.Symbol)
		require.Equal(t, "add", dav.Value)

		require.Equal(t, 2, len(dar.Args))

		ceArg1, ok := dar.Args[0].(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, ceArg1.To)
		arg1, ok := ceArg1.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, arg1.Type)
		require.Equal(t, "1", arg1.Value)

		ceArg2, ok := dar.Args[1].(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, ceArg2.To)
		arg2, ok := ceArg2.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, arg2.Type)
		require.Equal(t, "2", arg2.Value)

		da1, ok := fn1.Body.Stmts[1].(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, "b", da1.Symbol.Name)
		require.Equal(t, TInt64, da1.Symbol.Type)

		ceDar1, ok := da1.Right.(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt64, ceDar1.To)

		dav1, ok := ceDar1.Value.(*IdentExpr)
		require.Equal(t, true, ok)
		require.NotNil(t, dav1.Symbol)
		require.Equal(t, "a", dav1.Value)

		da2, ok := fn1.Body.Stmts[2].(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, "b", da2.Symbol.Name)
		require.Equal(t, TInt64, da2.Symbol.Type)

		ceDar2, ok := da2.Right.(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt64, ceDar2.To)

		dav2, ok := ceDar2.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.NotNil(t, dav2.Type)
		require.Equal(t, "42", dav2.Value)
	})
}
