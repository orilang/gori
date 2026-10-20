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
		require.Equal(t, false, target.Symbol.FromFunc)
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
		require.Equal(t, false, target.Symbol.FromFunc)
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
		require.Equal(t, false, target.Symbol.FromFunc)
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
		require.Equal(t, false, target.Symbol.FromFunc)

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
		require.Equal(t, false, target1.Symbol.FromFunc)

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
		require.Equal(t, false, cs.Symbol.FromFunc)

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
		require.Equal(t, 1, len(da.Symbol))
		require.Equal(t, "a", da.Symbol[0].Name)
		require.Equal(t, TInt, da.Symbol[0].Type)
		require.Equal(t, true, da.Symbol[0].FromFunc)

		require.Equal(t, 1, len(da.Right))
		dar, ok := da.Right[0].(*CallExpr)
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
		require.Equal(t, 1, len(da1.Symbol))
		require.Equal(t, "b", da1.Symbol[0].Name)
		require.Equal(t, TInt64, da1.Symbol[0].Type)
		require.Equal(t, false, da1.Symbol[0].FromFunc)

		require.Equal(t, 1, len(da1.Right))
		ceDar1, ok := da1.Right[0].(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt64, ceDar1.To)

		dav1, ok := ceDar1.Value.(*IdentExpr)
		require.Equal(t, true, ok)
		require.NotNil(t, dav1.Symbol)
		require.Equal(t, "a", dav1.Value)
		require.Equal(t, false, dav1.Symbol.FromFunc)

		da2, ok := fn1.Body.Stmts[2].(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(da2.Symbol))
		require.Equal(t, "b", da2.Symbol[0].Name)
		require.Equal(t, TInt64, da2.Symbol[0].Type)
		require.Equal(t, false, da2.Symbol[0].FromFunc)

		require.Equal(t, 1, len(da2.Right))
		ceDar2, ok := da2.Right[0].(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt64, ceDar2.To)

		dav2, ok := ceDar2.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.NotNil(t, dav2.Type)
		require.Equal(t, "42", dav2.Value)
	})

	t.Run("x9", func(t *testing.T) {
		data := `package main

func main() {
  a := int(0)
  b := int(1)

	switch a {
	case 1:
    b = int(1)
  case 2, 3:
    b = int(2)
	default:
    b = a
	}
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

		fn1, ok := pf.Decls[0].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, fn1.Name, "main")
		require.NotNil(t, fn1.Symbol)
		require.Equal(t, SymFunc, fn1.Symbol.Kind)
		require.Equal(t, 0, len(fn1.Params))
		require.Equal(t, 0, len(fn1.Results))

		require.NotNil(t, fn1.Body)
		require.Equal(t, 3, len(fn1.Body.Stmts))

		da1, ok := fn1.Body.Stmts[0].(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(da1.Right))
		require.Equal(t, "a", da1.Symbol[0].Name)
		require.Equal(t, TInt, da1.Symbol[0].Type)

		require.Equal(t, 1, len(da1.Right))
		dar1, ok := da1.Right[0].(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, dar1.To)
		arg1, ok := dar1.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, arg1.Type)
		require.Equal(t, "0", arg1.Value)

		da2, ok := fn1.Body.Stmts[1].(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(da2.Symbol))
		require.Equal(t, "b", da2.Symbol[0].Name)
		require.Equal(t, TInt, da2.Symbol[0].Type)

		require.Equal(t, 1, len(da2.Right))
		dar2, ok := da2.Right[0].(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, dar2.To)
		arg2, ok := dar2.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, arg2.Type)
		require.Equal(t, "1", arg2.Value)

		da3, ok := fn1.Body.Stmts[2].(*SwitchStmt)
		require.Equal(t, true, ok)
		require.Nil(t, da3.Init)
		require.NotNil(t, da3.Tag)
		swt, ok := da3.Tag.(*IdentExpr)
		require.Equal(t, true, ok)
		require.NotNil(t, swt.Symbol)
		require.Equal(t, TInt, swt.Type)
		require.Equal(t, "a", swt.Value)
		require.NotNil(t, da3.Cases)
		for ci, cc := range da3.Cases {
			switch ci {
			case 0:
				require.Equal(t, token.KWCase, cc.Case)
				require.Equal(t, 1, len(cc.Values))
				swc, ok := cc.Values[0].(*IntLitExpr)
				require.Equal(t, true, ok)
				require.Equal(t, TInt, swc.Type)
				require.Equal(t, "1", swc.Value)
			case 1:
				require.Equal(t, token.KWCase, cc.Case)
				require.Equal(t, 2, len(cc.Values))
				for _, ccv := range cc.Values {
					swcc, ok := ccv.(*IntLitExpr)
					require.Equal(t, true, ok)
					require.Equal(t, TInt, swcc.Type)
				}
			case 2:
				require.Equal(t, token.KWDefault, cc.Case)
				require.Equal(t, 0, len(cc.Values))
			}
			require.NotNil(t, cc.Body)
		}
	})

	t.Run("x10", func(t *testing.T) {
		data := `package main

func main() {
  a := int(0)
  b := int(1)
  c := int(1)

	switch a {
	case 1:
    b = int(1)
    c = int(1)
  case 2, 3:
    b = int(2)
    c = int(2)
	default:
    b = a
    c = int(2)
	}
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
		require.Equal(t, 1, len(pf.Decls))

		fn1, ok := pf.Decls[0].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, fn1.Name, "main")
		require.NotNil(t, fn1.Symbol)
		require.Equal(t, SymFunc, fn1.Symbol.Kind)
		require.Equal(t, 0, len(fn1.Params))
		require.Equal(t, 0, len(fn1.Results))

		require.NotNil(t, fn1.Body)
		require.Equal(t, 4, len(fn1.Body.Stmts))

		da1, ok := fn1.Body.Stmts[0].(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(da1.Symbol))
		require.Equal(t, "a", da1.Symbol[0].Name)
		require.Equal(t, TInt, da1.Symbol[0].Type)

		require.Equal(t, 1, len(da1.Right))
		dar1, ok := da1.Right[0].(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, dar1.To)
		arg1, ok := dar1.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, arg1.Type)
		require.Equal(t, "0", arg1.Value)

		da2, ok := fn1.Body.Stmts[1].(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(da2.Symbol))
		require.Equal(t, "b", da2.Symbol[0].Name)
		require.Equal(t, TInt, da2.Symbol[0].Type)

		require.Equal(t, 1, len(da2.Right))
		dar2, ok := da2.Right[0].(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, dar2.To)
		arg2, ok := dar2.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, arg2.Type)
		require.Equal(t, "1", arg2.Value)

		da3, ok := fn1.Body.Stmts[2].(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(da3.Symbol))
		require.Equal(t, "c", da3.Symbol[0].Name)
		require.Equal(t, TInt, da3.Symbol[0].Type)

		require.Equal(t, 1, len(da3.Right))
		dar3, ok := da3.Right[0].(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, dar2.To)
		arg3, ok := dar3.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, arg3.Type)
		require.Equal(t, "1", arg3.Value)

		da4, ok := fn1.Body.Stmts[3].(*SwitchStmt)
		require.Equal(t, true, ok)
		require.Nil(t, da4.Init)
		require.NotNil(t, da4.Tag)
		swt, ok := da4.Tag.(*IdentExpr)
		require.Equal(t, true, ok)
		require.NotNil(t, swt.Symbol)
		require.Equal(t, TInt, swt.Type)
		require.Equal(t, "a", swt.Value)
		require.NotNil(t, da4.Cases)
		for ci, cc := range da4.Cases {
			switch ci {
			case 0:
				require.Equal(t, token.KWCase, cc.Case)
				require.Equal(t, 1, len(cc.Values))
				swc, ok := cc.Values[0].(*IntLitExpr)
				require.Equal(t, true, ok)
				require.Equal(t, TInt, swc.Type)
				require.Equal(t, "1", swc.Value)
			case 1:
				require.Equal(t, token.KWCase, cc.Case)
				require.Equal(t, 2, len(cc.Values))
				for _, ccv := range cc.Values {
					swcc, ok := ccv.(*IntLitExpr)
					require.Equal(t, true, ok)
					require.Equal(t, TInt, swcc.Type)
				}
			case 2:
				require.Equal(t, token.KWDefault, cc.Case)
				require.Equal(t, 0, len(cc.Values))
			}
			require.NotNil(t, cc.Body)
		}
	})

	t.Run("x11", func(t *testing.T) {
		data := `package main

func main() {
  a := int(0)

  switch a {
  case 1:
    fallthrough
  default:
    a = int(1)
    a = a+1
  }
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
		require.Equal(t, 1, len(pf.Decls))

		fn1, ok := pf.Decls[0].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, fn1.Name, "main")
		require.NotNil(t, fn1.Symbol)
		require.Equal(t, SymFunc, fn1.Symbol.Kind)
		require.Equal(t, 0, len(fn1.Params))
		require.Equal(t, 0, len(fn1.Results))

		require.NotNil(t, fn1.Body)
		require.Equal(t, 2, len(fn1.Body.Stmts))

		da1, ok := fn1.Body.Stmts[0].(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(da1.Symbol))
		require.Equal(t, "a", da1.Symbol[0].Name)
		require.Equal(t, TInt, da1.Symbol[0].Type)

		require.Equal(t, 1, len(da1.Right))
		dar1, ok := da1.Right[0].(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, dar1.To)
		arg1, ok := dar1.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, arg1.Type)
		require.Equal(t, "0", arg1.Value)

		da2, ok := fn1.Body.Stmts[1].(*SwitchStmt)
		require.Equal(t, true, ok)
		require.Nil(t, da2.Init)
		require.NotNil(t, da2.Tag)
		swt, ok := da2.Tag.(*IdentExpr)
		require.Equal(t, true, ok)
		require.NotNil(t, swt.Symbol)
		require.Equal(t, TInt, swt.Type)
		require.Equal(t, "a", swt.Value)
		require.NotNil(t, da2.Cases)
		for ci, cc := range da2.Cases {
			switch ci {
			case 0:
				require.Equal(t, token.KWCase, cc.Case)
				require.Equal(t, 1, len(cc.Values))
				swc, ok := cc.Values[0].(*IntLitExpr)
				require.Equal(t, true, ok)
				require.Equal(t, TInt, swc.Type)
				require.Equal(t, "1", swc.Value)

				require.NotNil(t, cc.Body)
				require.Equal(t, 1, len(cc.Body))

				_, ok = cc.Body[0].(*FallThroughStmt)
				require.Equal(t, true, ok)

			case 1:
				require.Equal(t, token.KWDefault, cc.Case)
				require.Equal(t, 0, len(cc.Values))
				require.NotNil(t, cc.Body)
				require.Equal(t, 2, len(cc.Body))
			}
		}
	})

	t.Run("x12", func(t *testing.T) {
		data := `package main
func test() (string, int) {
  return "yes", int(1)
}

func f() {
  _,_ := test()
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
		require.Equal(t, 2, len(pf.Decls))

		fn, ok := pf.Decls[0].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, "test", fn.Name)
		require.NotNil(t, fn.Symbol)
		require.Equal(t, SymFunc, fn.Symbol.Kind)
		require.Equal(t, 0, len(fn.Params))
		require.Equal(t, 2, len(fn.Results))
		require.Equal(t, TString, fn.Results[0].Type)
		require.Equal(t, TInt, fn.Results[1].Type)

		require.NotNil(t, fn.Body)
		rs, ok := fn.Body.Stmts[0].(*ReturnStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 2, len(rs.Values))

		rv1, ok := rs.Values[0].(*StringLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TString, rv1.Type)
		require.Equal(t, "\"yes\"", rv1.Value)

		rv2, ok := rs.Values[1].(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, rv2.To)
		v2, ok := rv2.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, v2.Type)
		require.Equal(t, "1", v2.Value)

		fn1, ok := pf.Decls[1].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, fn1.Name, "f")
		require.NotNil(t, fn1.Symbol)
		require.Equal(t, SymFunc, fn1.Symbol.Kind)
		require.Equal(t, 0, len(fn1.Params))
		require.Equal(t, 0, len(fn1.Results))

		require.NotNil(t, fn1.Body)
		require.Equal(t, 1, len(fn1.Body.Stmts))

		da, ok := fn1.Body.Stmts[0].(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 2, len(da.Symbol))
		require.Equal(t, "_", da.Symbol[0].Name)
		require.Equal(t, TString, da.Symbol[0].Type)

		require.Equal(t, "_", da.Symbol[1].Name)
		require.Equal(t, TInt, da.Symbol[1].Type)

		require.Equal(t, 1, len(da.Right))
		dar, ok := da.Right[0].(*CallExpr)
		require.Equal(t, true, ok)
		fnx, ok := dar.CalleeType.(*FuncMethod)
		require.Equal(t, true, ok)
		require.Equal(t, "test", fnx.Name)
	})

	t.Run("x13", func(t *testing.T) {
		data := `package main
func test() (string, int) {
  return "yes", int(1)
}

func f() {
  a,_ := test()
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
		require.Equal(t, 2, len(pf.Decls))

		fn, ok := pf.Decls[0].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, "test", fn.Name)
		require.NotNil(t, fn.Symbol)
		require.Equal(t, SymFunc, fn.Symbol.Kind)
		require.Equal(t, 0, len(fn.Params))
		require.Equal(t, 2, len(fn.Results))
		require.Equal(t, TString, fn.Results[0].Type)
		require.Equal(t, TInt, fn.Results[1].Type)

		require.NotNil(t, fn.Body)
		rs, ok := fn.Body.Stmts[0].(*ReturnStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 2, len(rs.Values))

		rv1, ok := rs.Values[0].(*StringLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TString, rv1.Type)
		require.Equal(t, "\"yes\"", rv1.Value)

		rv2, ok := rs.Values[1].(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, rv2.To)
		v2, ok := rv2.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, v2.Type)
		require.Equal(t, "1", v2.Value)

		fn1, ok := pf.Decls[1].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, fn1.Name, "f")
		require.NotNil(t, fn1.Symbol)
		require.Equal(t, SymFunc, fn1.Symbol.Kind)
		require.Equal(t, 0, len(fn1.Params))
		require.Equal(t, 0, len(fn1.Results))

		require.NotNil(t, fn1.Body)
		require.Equal(t, 1, len(fn1.Body.Stmts))

		da, ok := fn1.Body.Stmts[0].(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 2, len(da.Symbol))
		require.Equal(t, "a", da.Symbol[0].Name)
		require.Equal(t, TString, da.Symbol[0].Type)
		require.Equal(t, true, da.Symbol[0].FromFunc)

		require.Equal(t, "_", da.Symbol[1].Name)
		require.Equal(t, TInt, da.Symbol[1].Type)
		require.Equal(t, true, da.Symbol[1].FromFunc)

		require.Equal(t, 1, len(da.Right))
		dar, ok := da.Right[0].(*CallExpr)
		require.Equal(t, true, ok)
		fnx, ok := dar.CalleeType.(*FuncMethod)
		require.Equal(t, true, ok)
		require.Equal(t, "test", fnx.Name)
	})

	t.Run("x14", func(t *testing.T) {
		data := `package main
func test() (string, int) {
  return "yes", int(1)
}

func f() {
  a, b := test()
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
		require.Equal(t, 2, len(pf.Decls))

		fn, ok := pf.Decls[0].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, "test", fn.Name)
		require.NotNil(t, fn.Symbol)
		require.Equal(t, SymFunc, fn.Symbol.Kind)
		require.Equal(t, 0, len(fn.Params))
		require.Equal(t, 2, len(fn.Results))
		require.Equal(t, TString, fn.Results[0].Type)
		require.Equal(t, TInt, fn.Results[1].Type)

		require.NotNil(t, fn.Body)
		rs, ok := fn.Body.Stmts[0].(*ReturnStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 2, len(rs.Values))

		rv1, ok := rs.Values[0].(*StringLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TString, rv1.Type)
		require.Equal(t, "\"yes\"", rv1.Value)

		rv2, ok := rs.Values[1].(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, rv2.To)
		v2, ok := rv2.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, v2.Type)
		require.Equal(t, "1", v2.Value)

		fn1, ok := pf.Decls[1].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, fn1.Name, "f")
		require.NotNil(t, fn1.Symbol)
		require.Equal(t, SymFunc, fn1.Symbol.Kind)
		require.Equal(t, 0, len(fn1.Params))
		require.Equal(t, 0, len(fn1.Results))

		require.NotNil(t, fn1.Body)
		require.Equal(t, 1, len(fn1.Body.Stmts))

		da, ok := fn1.Body.Stmts[0].(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 2, len(da.Symbol))
		require.Equal(t, "a", da.Symbol[0].Name)
		require.Equal(t, TString, da.Symbol[0].Type)
		require.Equal(t, true, da.Symbol[0].FromFunc)

		require.Equal(t, "b", da.Symbol[1].Name)
		require.Equal(t, TInt, da.Symbol[1].Type)
		require.Equal(t, true, da.Symbol[1].FromFunc)

		require.Equal(t, 1, len(da.Right))
		dar, ok := da.Right[0].(*CallExpr)
		require.Equal(t, true, ok)
		fnx, ok := dar.CalleeType.(*FuncMethod)
		require.Equal(t, true, ok)
		require.Equal(t, "test", fnx.Name)
	})

	t.Run("x15", func(t *testing.T) {
		data := `package main
func test() (string, int) {
  return "yes", int(1)
}

func f() {
  a, _ := test()
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
		require.Equal(t, 2, len(pf.Decls))

		fn, ok := pf.Decls[0].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, "test", fn.Name)
		require.NotNil(t, fn.Symbol)
		require.Equal(t, SymFunc, fn.Symbol.Kind)
		require.Equal(t, 0, len(fn.Params))
		require.Equal(t, 2, len(fn.Results))
		require.Equal(t, TString, fn.Results[0].Type)
		require.Equal(t, TInt, fn.Results[1].Type)

		require.NotNil(t, fn.Body)
		rs, ok := fn.Body.Stmts[0].(*ReturnStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 2, len(rs.Values))

		rv1, ok := rs.Values[0].(*StringLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TString, rv1.Type)
		require.Equal(t, "\"yes\"", rv1.Value)

		rv2, ok := rs.Values[1].(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, rv2.To)
		v2, ok := rv2.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, v2.Type)
		require.Equal(t, "1", v2.Value)

		fn1, ok := pf.Decls[1].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, fn1.Name, "f")
		require.NotNil(t, fn1.Symbol)
		require.Equal(t, SymFunc, fn1.Symbol.Kind)
		require.Equal(t, 0, len(fn1.Params))
		require.Equal(t, 0, len(fn1.Results))

		require.NotNil(t, fn1.Body)
		require.Equal(t, 1, len(fn1.Body.Stmts))

		da, ok := fn1.Body.Stmts[0].(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 2, len(da.Symbol))
		require.Equal(t, "a", da.Symbol[0].Name)
		require.Equal(t, TString, da.Symbol[0].Type)
		require.Equal(t, true, da.Symbol[0].FromFunc)

		require.Equal(t, "_", da.Symbol[1].Name)
		require.NotNil(t, da.Symbol[1].Type)

		require.Equal(t, 1, len(da.Right))
		dar, ok := da.Right[0].(*CallExpr)
		require.Equal(t, true, ok)
		fnx, ok := dar.CalleeType.(*FuncMethod)
		require.Equal(t, true, ok)
		require.Equal(t, "test", fnx.Name)
	})

	t.Run("x17", func(t *testing.T) {
		data := `package main
func test1() string {
  return "yes"
}

func test2() int {
  return int(0)
}

func f() {
  a,b := test1(),test2()
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
		require.Equal(t, 3, len(pf.Decls))

		fn1, ok := pf.Decls[0].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, "test1", fn1.Name)
		require.NotNil(t, fn1.Symbol)
		require.Equal(t, SymFunc, fn1.Symbol.Kind)
		require.Equal(t, 0, len(fn1.Params))
		require.Equal(t, 1, len(fn1.Results))
		require.Equal(t, TString, fn1.Results[0].Type)

		require.NotNil(t, fn1.Body)
		rs1, ok := fn1.Body.Stmts[0].(*ReturnStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(rs1.Values))

		rv1, ok := rs1.Values[0].(*StringLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TString, rv1.Type)
		require.Equal(t, "\"yes\"", rv1.Value)

		fn2, ok := pf.Decls[1].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, "test2", fn2.Name)
		require.NotNil(t, fn2.Symbol)
		require.Equal(t, SymFunc, fn2.Symbol.Kind)
		require.Equal(t, 0, len(fn2.Params))
		require.Equal(t, 1, len(fn2.Results))
		require.Equal(t, TInt, fn2.Results[0].Type)

		require.NotNil(t, fn2.Body)
		rs2, ok := fn2.Body.Stmts[0].(*ReturnStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(rs2.Values))

		ce1, ok := rs2.Values[0].(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, ce1.To)
		rv2, ok := ce1.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, rv2.Type)
		require.Equal(t, "0", rv2.Value)

		fn3, ok := pf.Decls[2].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, fn3.Name, "f")
		require.NotNil(t, fn3.Symbol)
		require.Equal(t, SymFunc, fn3.Symbol.Kind)
		require.Equal(t, 0, len(fn3.Params))
		require.Equal(t, 0, len(fn3.Results))

		require.NotNil(t, fn3.Body)
		require.Equal(t, 1, len(fn3.Body.Stmts))

		da, ok := fn3.Body.Stmts[0].(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 2, len(da.Symbol))
		require.Equal(t, "a", da.Symbol[0].Name)
		require.Equal(t, TString, da.Symbol[0].Type)
		require.Equal(t, true, da.Symbol[0].FromFunc)

		require.Equal(t, "b", da.Symbol[1].Name)
		require.Equal(t, TInt, da.Symbol[1].Type)
		require.Equal(t, true, da.Symbol[1].FromFunc)

		require.Equal(t, 2, len(da.Right))
		ca1, ok := da.Right[0].(*CallExpr)
		require.Equal(t, true, ok)
		require.Equal(t, 0, len(ca1.Args))

		callee1, ok := ca1.Callee.(*IdentExpr)
		require.Equal(t, true, ok)
		require.Equal(t, "test1", callee1.Value)
		calleeType1, ok := ca1.CalleeType.(*FuncMethod)
		require.Equal(t, true, ok)
		require.Equal(t, "test1", calleeType1.Name)

		ca2, ok := da.Right[1].(*CallExpr)
		require.Equal(t, true, ok)
		require.Equal(t, 0, len(ca2.Args))

		callee2, ok := ca2.Callee.(*IdentExpr)
		require.Equal(t, true, ok)
		require.Equal(t, "test2", callee2.Value)
		calleeType2, ok := ca2.CalleeType.(*FuncMethod)
		require.Equal(t, true, ok)
		require.Equal(t, "test2", calleeType2.Name)
	})

	t.Run("x18", func(t *testing.T) {
		data := `package main
func test1() string {
  return "yes"
}

func test2() int {
  return int(0)
}

func f() {
  a, b := "no",int(1)
  a,b = test1(),test2()
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
		require.Equal(t, 3, len(pf.Decls))

		fn1, ok := pf.Decls[0].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, "test1", fn1.Name)
		require.NotNil(t, fn1.Symbol)
		require.Equal(t, SymFunc, fn1.Symbol.Kind)
		require.Equal(t, 0, len(fn1.Params))
		require.Equal(t, 1, len(fn1.Results))
		require.Equal(t, TString, fn1.Results[0].Type)

		require.NotNil(t, fn1.Body)
		rs1, ok := fn1.Body.Stmts[0].(*ReturnStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(rs1.Values))

		rv1, ok := rs1.Values[0].(*StringLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TString, rv1.Type)
		require.Equal(t, "\"yes\"", rv1.Value)

		fn2, ok := pf.Decls[1].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, "test2", fn2.Name)
		require.NotNil(t, fn2.Symbol)
		require.Equal(t, SymFunc, fn2.Symbol.Kind)
		require.Equal(t, 0, len(fn2.Params))
		require.Equal(t, 1, len(fn2.Results))
		require.Equal(t, TInt, fn2.Results[0].Type)

		require.NotNil(t, fn2.Body)
		rs2, ok := fn2.Body.Stmts[0].(*ReturnStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(rs2.Values))

		ce1, ok := rs2.Values[0].(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, ce1.To)
		rv2, ok := ce1.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, rv2.Type)
		require.Equal(t, "0", rv2.Value)

		fn3, ok := pf.Decls[2].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, fn3.Name, "f")
		require.NotNil(t, fn3.Symbol)
		require.Equal(t, SymFunc, fn3.Symbol.Kind)
		require.Equal(t, 0, len(fn3.Params))
		require.Equal(t, 0, len(fn3.Results))

		require.NotNil(t, fn3.Body)
		require.Equal(t, 2, len(fn3.Body.Stmts))

		da1, ok := fn3.Body.Stmts[0].(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 2, len(da1.Symbol))
		require.Equal(t, "a", da1.Symbol[0].Name)
		require.Equal(t, TString, da1.Symbol[0].Type)
		require.Equal(t, false, da1.Symbol[0].FromFunc)

		require.Equal(t, "b", da1.Symbol[1].Name)
		require.Equal(t, TInt, da1.Symbol[1].Type)
		require.Equal(t, false, da1.Symbol[1].FromFunc)

		require.Equal(t, 2, len(da1.Right))
		darv1, ok := da1.Right[0].(*StringLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TString, darv1.Type)
		require.Equal(t, "\"no\"", darv1.Value)

		ce2, ok := da1.Right[1].(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, ce2.To)

		darv2, ok := ce2.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, darv2.Type)
		require.Equal(t, "1", darv2.Value)

		da2, ok := fn3.Body.Stmts[1].(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 2, len(da2.Symbol))
		require.Equal(t, "a", da2.Symbol[0].Name)
		require.Equal(t, TString, da2.Symbol[0].Type)
		require.Equal(t, false, da2.Symbol[0].FromFunc)

		require.Equal(t, "b", da2.Symbol[1].Name)
		require.Equal(t, TInt, da2.Symbol[1].Type)
		require.Equal(t, false, da2.Symbol[1].FromFunc)

		require.Equal(t, 2, len(da2.Right))
		ca1, ok := da2.Right[0].(*CallExpr)
		require.Equal(t, true, ok)
		require.Equal(t, 0, len(ca1.Args))

		callee1, ok := ca1.Callee.(*IdentExpr)
		require.Equal(t, true, ok)
		require.Equal(t, "test1", callee1.Value)
		calleeType1, ok := ca1.CalleeType.(*FuncMethod)
		require.Equal(t, true, ok)
		require.Equal(t, "test1", calleeType1.Name)

		ca2, ok := da2.Right[1].(*CallExpr)
		require.Equal(t, true, ok)
		require.Equal(t, 0, len(ca2.Args))

		callee2, ok := ca2.Callee.(*IdentExpr)
		require.Equal(t, true, ok)
		require.Equal(t, "test2", callee2.Value)
		calleeType2, ok := ca2.CalleeType.(*FuncMethod)
		require.Equal(t, true, ok)
		require.Equal(t, "test2", calleeType2.Name)
	})

	t.Run("x19", func(t *testing.T) {
		data := `package main
func test1() string {
  return "yes"
}

func test2() int {
  return int(0)
}

func f() {
  a, b := "no",int(1)
  a,_ = test1(),test2()
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
		require.Equal(t, 3, len(pf.Decls))

		fn1, ok := pf.Decls[0].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, "test1", fn1.Name)
		require.NotNil(t, fn1.Symbol)
		require.Equal(t, SymFunc, fn1.Symbol.Kind)
		require.Equal(t, 0, len(fn1.Params))
		require.Equal(t, 1, len(fn1.Results))
		require.Equal(t, TString, fn1.Results[0].Type)

		require.NotNil(t, fn1.Body)
		rs1, ok := fn1.Body.Stmts[0].(*ReturnStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(rs1.Values))

		rv1, ok := rs1.Values[0].(*StringLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TString, rv1.Type)
		require.Equal(t, "\"yes\"", rv1.Value)

		fn2, ok := pf.Decls[1].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, "test2", fn2.Name)
		require.NotNil(t, fn2.Symbol)
		require.Equal(t, SymFunc, fn2.Symbol.Kind)
		require.Equal(t, 0, len(fn2.Params))
		require.Equal(t, 1, len(fn2.Results))
		require.Equal(t, TInt, fn2.Results[0].Type)

		require.NotNil(t, fn2.Body)
		rs2, ok := fn2.Body.Stmts[0].(*ReturnStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(rs2.Values))

		ce1, ok := rs2.Values[0].(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, ce1.To)
		rv2, ok := ce1.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, rv2.Type)
		require.Equal(t, "0", rv2.Value)

		fn3, ok := pf.Decls[2].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, fn3.Name, "f")
		require.NotNil(t, fn3.Symbol)
		require.Equal(t, SymFunc, fn3.Symbol.Kind)
		require.Equal(t, 0, len(fn3.Params))
		require.Equal(t, 0, len(fn3.Results))

		require.NotNil(t, fn3.Body)
		require.Equal(t, 2, len(fn3.Body.Stmts))

		da1, ok := fn3.Body.Stmts[0].(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 2, len(da1.Symbol))
		require.Equal(t, "a", da1.Symbol[0].Name)
		require.Equal(t, TString, da1.Symbol[0].Type)

		require.Equal(t, "b", da1.Symbol[1].Name)
		require.Equal(t, TInt, da1.Symbol[1].Type)

		require.Equal(t, 2, len(da1.Right))
		darv1, ok := da1.Right[0].(*StringLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TString, darv1.Type)
		require.Equal(t, "\"no\"", darv1.Value)

		ce2, ok := da1.Right[1].(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, ce2.To)

		darv2, ok := ce2.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, darv2.Type)
		require.Equal(t, "1", darv2.Value)

		da2, ok := fn3.Body.Stmts[1].(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 2, len(da2.Symbol))
		require.Equal(t, "a", da2.Symbol[0].Name)
		require.Equal(t, TString, da2.Symbol[0].Type)

		require.Equal(t, "", da2.Symbol[1].Name)
		require.Equal(t, true, da2.Symbol[1].IsBlank)
		require.NotNil(t, da2.Symbol[1].Type)

		require.Equal(t, 2, len(da2.Right))
		ca1, ok := da2.Right[0].(*CallExpr)
		require.Equal(t, true, ok)
		require.Equal(t, 0, len(ca1.Args))

		callee1, ok := ca1.Callee.(*IdentExpr)
		require.Equal(t, true, ok)
		require.Equal(t, "test1", callee1.Value)
		calleeType1, ok := ca1.CalleeType.(*FuncMethod)
		require.Equal(t, true, ok)
		require.Equal(t, "test1", calleeType1.Name)

		ca2, ok := da2.Right[1].(*CallExpr)
		require.Equal(t, true, ok)
		require.Equal(t, 0, len(ca2.Args))

		callee2, ok := ca2.Callee.(*IdentExpr)
		require.Equal(t, true, ok)
		require.Equal(t, "test2", callee2.Value)
		calleeType2, ok := ca2.CalleeType.(*FuncMethod)
		require.Equal(t, true, ok)
		require.Equal(t, "test2", calleeType2.Name)
	})

	t.Run("x20", func(t *testing.T) {
		data := `package main
func multi(a int) int {
  return a*2
}

func f() {
  a := int(1)
  a = multi(int(4))
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
		require.Equal(t, 2, len(pf.Decls))

		fn1, ok := pf.Decls[0].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, "multi", fn1.Name)
		require.NotNil(t, fn1.Symbol)
		require.Equal(t, SymFunc, fn1.Symbol.Kind)
		require.Equal(t, 1, len(fn1.Params))
		require.Equal(t, 1, len(fn1.Results))
		require.Equal(t, "a", fn1.Params[0].Name)
		require.Equal(t, TInt, fn1.Params[0].Type)
		require.Equal(t, TInt, fn1.Results[0].Type)

		require.NotNil(t, fn1.Body)
		rs1, ok := fn1.Body.Stmts[0].(*ReturnStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(rs1.Values))

		rv1, ok := rs1.Values[0].(*BinaryExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, rv1.Type)
		rv1i, ok := rv1.Left.(*IdentExpr)
		require.Equal(t, true, ok)
		require.Equal(t, "a", rv1i.Value)
		require.Equal(t, TInt, rv1i.Type)
		require.Equal(t, token.Star, rv1.Operator)

		fn3, ok := pf.Decls[1].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, fn3.Name, "f")
		require.NotNil(t, fn3.Symbol)
		require.Equal(t, SymFunc, fn3.Symbol.Kind)
		require.Equal(t, 0, len(fn3.Params))
		require.Equal(t, 0, len(fn3.Results))

		require.NotNil(t, fn3.Body)
		require.Equal(t, 2, len(fn3.Body.Stmts))

		da1, ok := fn3.Body.Stmts[0].(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(da1.Symbol))
		require.Equal(t, "a", da1.Symbol[0].Name)
		require.Equal(t, TInt, da1.Symbol[0].Type)

		require.Equal(t, 1, len(da1.Right))
		ce1, ok := da1.Right[0].(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, ce1.To)

		darv2, ok := ce1.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, darv2.Type)
		require.Equal(t, "1", darv2.Value)

		da2, ok := fn3.Body.Stmts[1].(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(da2.Symbol))
		require.Equal(t, "a", da2.Symbol[0].Name)
		require.Equal(t, false, da2.Symbol[0].FromFunc)
		require.Equal(t, TInt, da2.Symbol[0].Type)

		require.Equal(t, 1, len(da2.Right))
		ca1, ok := da2.Right[0].(*CallExpr)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(ca1.Args))

		callee1, ok := ca1.Callee.(*IdentExpr)
		require.Equal(t, true, ok)
		require.Equal(t, "multi", callee1.Value)
		calleeType1, ok := ca1.CalleeType.(*FuncMethod)
		require.Equal(t, true, ok)
		require.Equal(t, "multi", calleeType1.Name)
	})

	t.Run("x21", func(t *testing.T) {
		data := `package main
func f(a int, b int) bool {
  for a < b {
    return false
  }
  return false
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

		fn1, ok := pf.Decls[0].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, "f", fn1.Name)
		require.NotNil(t, fn1.Symbol)
		require.Equal(t, SymFunc, fn1.Symbol.Kind)
		require.Equal(t, 2, len(fn1.Params))
		require.Equal(t, "a", fn1.Params[0].Name)
		require.Equal(t, TInt, fn1.Params[0].Type)
		require.Equal(t, "b", fn1.Params[1].Name)
		require.Equal(t, TInt, fn1.Params[1].Type)

		require.Equal(t, 1, len(fn1.Results))
		require.Equal(t, "", fn1.Results[0].Name)
		require.Equal(t, TBool, fn1.Results[0].Type)
		require.NotNil(t, fn1.Body)
		require.Equal(t, 2, len(fn1.Body.Stmts))

		fr1, ok := fn1.Body.Stmts[0].(*ForStmt)
		require.Equal(t, true, ok)
		require.Nil(t, fr1.Init)
		require.NotNil(t, fr1.Condition)
		require.Nil(t, fr1.Post)

		bn1, ok := fr1.Condition.(*BinaryExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TBool, bn1.Type)

		require.NotNil(t, bn1.Left)
		l1, ok := bn1.Left.(*IdentExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, l1.Type)
		require.Equal(t, "a", l1.Value)

		r1, ok := bn1.Right.(*IdentExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, r1.Type)
		require.Equal(t, "b", r1.Value)

		require.Equal(t, token.Lt, bn1.Operator)
		require.NotNil(t, fr1.Body)
		require.Equal(t, 1, len(fr1.Body.Stmts))
		rst1, ok := fr1.Body.Stmts[0].(*ReturnStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(rst1.Values))

		bl1, ok := rst1.Values[0].(*BoolLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TBool, bl1.Type)
		require.Equal(t, "false", bl1.Value)

		rst2, ok := fn1.Body.Stmts[1].(*ReturnStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(rst2.Values))

		bl2, ok := rst2.Values[0].(*BoolLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TBool, bl2.Type)
		require.Equal(t, "false", bl2.Value)
	})

	t.Run("x22", func(t *testing.T) {
		data := `package main
func f() bool {
  for {
    break
  }
  return false
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

		fn1, ok := pf.Decls[0].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, "f", fn1.Name)
		require.NotNil(t, fn1.Symbol)
		require.Equal(t, SymFunc, fn1.Symbol.Kind)
		require.Equal(t, 0, len(fn1.Params))

		require.Equal(t, 1, len(fn1.Results))
		require.Equal(t, "", fn1.Results[0].Name)
		require.Equal(t, TBool, fn1.Results[0].Type)
		require.NotNil(t, fn1.Body)
		require.Equal(t, 2, len(fn1.Body.Stmts))

		fr1, ok := fn1.Body.Stmts[0].(*ForStmt)
		require.Equal(t, true, ok)
		require.Nil(t, fr1.Init)
		require.Nil(t, fr1.Condition)
		require.Nil(t, fr1.Post)

		require.NotNil(t, fr1.Body)
		require.Equal(t, 1, len(fr1.Body.Stmts))
		_, ok = fr1.Body.Stmts[0].(*BreakStmt)
		require.Equal(t, true, ok)

		rst1, ok := fn1.Body.Stmts[1].(*ReturnStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(rst1.Values))

		bl1, ok := rst1.Values[0].(*BoolLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TBool, bl1.Type)
		require.Equal(t, "false", bl1.Value)
	})

	t.Run("x23", func(t *testing.T) {
		data := `package main
func f(run bool) (result int) {
  result = int(1)

  for i := int(0); run; i += int(1) {
    continue
  }

  return result
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

		fn1, ok := pf.Decls[0].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, "f", fn1.Name)
		require.NotNil(t, fn1.Symbol)
		require.Equal(t, SymFunc, fn1.Symbol.Kind)
		require.Equal(t, 1, len(fn1.Params))
		require.Equal(t, "run", fn1.Params[0].Name)
		require.Equal(t, TBool, fn1.Params[0].Type)

		require.Equal(t, 1, len(fn1.Results))
		require.Equal(t, "result", fn1.Results[0].Name)
		require.Equal(t, TInt, fn1.Results[0].Type)
		require.NotNil(t, fn1.Body)
		require.Equal(t, 3, len(fn1.Body.Stmts))

		as1, ok := fn1.Body.Stmts[0].(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(as1.Right))

		require.Equal(t, 1, len(as1.Symbol))
		require.Equal(t, "result", as1.Symbol[0].Name)
		require.Equal(t, TInt, as1.Symbol[0].Type)

		ce1, ok := as1.Right[0].(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, ce1.To)

		av1, ok := ce1.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, av1.Type)
		require.Equal(t, "1", av1.Value)

		fr1, ok := fn1.Body.Stmts[1].(*ForStmt)
		require.Equal(t, true, ok)
		require.NotNil(t, fr1.Init)
		require.NotNil(t, fr1.Condition)
		require.NotNil(t, fr1.Post)

		init, ok := fr1.Init.(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(init.Symbol))
		require.Equal(t, "i", init.Symbol[0].Name)

		require.Equal(t, token.Define, init.Operator)

		require.Equal(t, 1, len(init.Right))
		initr, ok := init.Right[0].(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, initr.To)

		iav1, ok := initr.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, iav1.Type)
		require.Equal(t, "0", iav1.Value)

		cond, ok := fr1.Condition.(*IdentExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TBool, cond.Type)
		require.Equal(t, "run", cond.Value)

		post, ok := fr1.Post.(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(post.Symbol))
		require.Equal(t, "i", post.Symbol[0].Name)

		require.Equal(t, token.PlusEq, post.Operator)

		require.Equal(t, 1, len(post.Right))
		postr, ok := post.Right[0].(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, postr.To)

		pav1, ok := postr.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, pav1.Type)
		require.Equal(t, "1", pav1.Value)

		require.NotNil(t, fr1.Body)
		require.Equal(t, 1, len(fr1.Body.Stmts))
		_, ok = fr1.Body.Stmts[0].(*ContinueStmt)
		require.Equal(t, true, ok)

		rst1, ok := fn1.Body.Stmts[2].(*ReturnStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(rst1.Values))

		id1, ok := rst1.Values[0].(*IdentExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, id1.Type)
		require.Equal(t, "result", id1.Value)
	})

	t.Run("x24", func(t *testing.T) {
		data := `package main
func f() (result int) {
  result = int(1)
  for k := range int(5) {
    result = int(2)
  }
  return
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

		fn1, ok := pf.Decls[0].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, "f", fn1.Name)
		require.NotNil(t, fn1.Symbol)
		require.Equal(t, SymFunc, fn1.Symbol.Kind)
		require.Equal(t, 0, len(fn1.Params))
		require.Equal(t, 1, len(fn1.Results))
		require.Equal(t, "result", fn1.Results[0].Name)
		require.Equal(t, TInt, fn1.Results[0].Type)

		require.NotNil(t, fn1.Body)
		as1, ok := fn1.Body.Stmts[0].(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(as1.Right))

		require.Equal(t, 1, len(as1.Symbol))
		require.Equal(t, "result", as1.Symbol[0].Name)
		require.Equal(t, TInt, as1.Symbol[0].Type)

		ce1, ok := as1.Right[0].(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, ce1.To)

		av1, ok := ce1.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, av1.Type)
		require.Equal(t, "1", av1.Value)

		fr1, ok := fn1.Body.Stmts[1].(*RangeStmt)
		require.Equal(t, true, ok)
		require.NotNil(t, fr1.Key)

		key, ok := fr1.Key.(*IdentExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, key.Type)
		require.Equal(t, "k", key.Value)
		require.Equal(t, token.Define, fr1.Op)
		require.Equal(t, token.KWRange, fr1.Range)

		rg1, ok := fr1.X.(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, rg1.To)

		arv1, ok := rg1.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, arv1.Type)
		require.Equal(t, "5", arv1.Value)

		require.Equal(t, 1, len(fr1.Body.Stmts))
		br1, ok := fr1.Body.Stmts[0].(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(br1.Right))

		require.Equal(t, 1, len(br1.Symbol))
		require.Equal(t, "result", br1.Symbol[0].Name)
		require.Equal(t, TInt, br1.Symbol[0].Type)

		rst1, ok := fn1.Body.Stmts[2].(*ReturnStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 0, len(rst1.Values))
	})

	t.Run("x25", func(t *testing.T) {
		data := `package main
func f() (result int) {
  result = int(1)
  for k,v := range int(5) {
    result = int(2)
  }
  return
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

		fn1, ok := pf.Decls[0].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, "f", fn1.Name)
		require.NotNil(t, fn1.Symbol)
		require.Equal(t, SymFunc, fn1.Symbol.Kind)
		require.Equal(t, 0, len(fn1.Params))
		require.Equal(t, 1, len(fn1.Results))
		require.Equal(t, "result", fn1.Results[0].Name)
		require.Equal(t, TInt, fn1.Results[0].Type)

		require.NotNil(t, fn1.Body)
		as1, ok := fn1.Body.Stmts[0].(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(as1.Right))

		require.Equal(t, 1, len(as1.Symbol))
		require.Equal(t, "result", as1.Symbol[0].Name)
		require.Equal(t, TInt, as1.Symbol[0].Type)

		ce1, ok := as1.Right[0].(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, ce1.To)

		av1, ok := ce1.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, av1.Type)
		require.Equal(t, "1", av1.Value)

		fr1, ok := fn1.Body.Stmts[1].(*RangeStmt)
		require.Equal(t, true, ok)
		require.NotNil(t, fr1.Key)
		require.NotNil(t, fr1.Value)

		key, ok := fr1.Key.(*IdentExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, key.Type)
		require.Equal(t, "k", key.Value)

		value, ok := fr1.Value.(*IdentExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, value.Type)
		require.Equal(t, "v", value.Value)

		require.Equal(t, token.Define, fr1.Op)
		require.Equal(t, token.KWRange, fr1.Range)

		rg1, ok := fr1.X.(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, rg1.To)

		arv1, ok := rg1.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, arv1.Type)
		require.Equal(t, "5", arv1.Value)

		require.Equal(t, 1, len(fr1.Body.Stmts))
		br1, ok := fr1.Body.Stmts[0].(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(br1.Right))

		require.Equal(t, 1, len(br1.Symbol))
		require.Equal(t, "result", br1.Symbol[0].Name)
		require.Equal(t, TInt, br1.Symbol[0].Type)

		rst1, ok := fn1.Body.Stmts[2].(*ReturnStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 0, len(rst1.Values))
	})

	t.Run("x26", func(t *testing.T) {
		data := `package main
func f(z int){
	for z = range int(5) {}
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

		fn1, ok := pf.Decls[0].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, "f", fn1.Name)
		require.NotNil(t, fn1.Symbol)
		require.Equal(t, SymFunc, fn1.Symbol.Kind)
		require.Equal(t, 1, len(fn1.Params))
		require.Equal(t, "z", fn1.Params[0].Name)
		require.Equal(t, TInt, fn1.Params[0].Type)
		require.Equal(t, 0, len(fn1.Results))

		require.NotNil(t, fn1.Body)
		fr1, ok := fn1.Body.Stmts[0].(*RangeStmt)
		require.Equal(t, true, ok)
		require.NotNil(t, fr1.Key)
		require.Nil(t, fr1.Value)
		require.NotNil(t, fr1.Body)

		key, ok := fr1.Key.(*IdentExpr)
		require.Equal(t, true, ok)
		require.Equal(t, "z", key.Symbol.Name)
		require.Equal(t, TInt, key.Symbol.Type)

		ce1, ok := fr1.X.(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, ce1.To)

		av1, ok := ce1.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, av1.Type)
		require.Equal(t, "5", av1.Value)

		require.Equal(t, 0, len(fr1.Body.Stmts))
	})

	t.Run("x27", func(t *testing.T) {
		data := `package main
func f(run bool) (result int) {
  result = int(1)

  for i := int(0); run; i++ {
    continue
  }

  return result
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

		fn1, ok := pf.Decls[0].(*FuncDecl)
		require.Equal(t, true, ok)
		require.Equal(t, "f", fn1.Name)
		require.NotNil(t, fn1.Symbol)
		require.Equal(t, SymFunc, fn1.Symbol.Kind)
		require.Equal(t, 1, len(fn1.Params))
		require.Equal(t, "run", fn1.Params[0].Name)
		require.Equal(t, TBool, fn1.Params[0].Type)

		require.Equal(t, 1, len(fn1.Results))
		require.Equal(t, "result", fn1.Results[0].Name)
		require.Equal(t, TInt, fn1.Results[0].Type)
		require.NotNil(t, fn1.Body)
		require.Equal(t, 3, len(fn1.Body.Stmts))

		as1, ok := fn1.Body.Stmts[0].(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(as1.Right))

		require.Equal(t, 1, len(as1.Symbol))
		require.Equal(t, "result", as1.Symbol[0].Name)
		require.Equal(t, TInt, as1.Symbol[0].Type)

		ce1, ok := as1.Right[0].(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, ce1.To)

		av1, ok := ce1.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, av1.Type)
		require.Equal(t, "1", av1.Value)

		fr1, ok := fn1.Body.Stmts[1].(*ForStmt)
		require.Equal(t, true, ok)
		require.NotNil(t, fr1.Init)
		require.NotNil(t, fr1.Condition)
		require.NotNil(t, fr1.Post)

		init, ok := fr1.Init.(*AssigmentStmt)
		require.Equal(t, true, ok)
		require.Equal(t, 1, len(init.Symbol))
		require.Equal(t, "i", init.Symbol[0].Name)

		require.Equal(t, token.Define, init.Operator)

		require.Equal(t, 1, len(init.Right))
		initr, ok := init.Right[0].(*ConversionExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, initr.To)

		iav1, ok := initr.Value.(*IntLitExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, iav1.Type)
		require.Equal(t, "0", iav1.Value)

		cond, ok := fr1.Condition.(*IdentExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TBool, cond.Type)
		require.Equal(t, "run", cond.Value)

		post, ok := fr1.Post.(*IncDecStmt)
		require.Equal(t, true, ok)
		id1, ok := post.X.(*IdentExpr)
		require.Equal(t, true, ok)
		require.Equal(t, TInt, id1.Type)
		require.Equal(t, "i", id1.Symbol.Name)
		require.Equal(t, token.PPlus, post.Operator)
	})
}
