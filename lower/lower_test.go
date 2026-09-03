package lower

import (
	"fmt"
	"testing"

	"github.com/orilang/gori/ir"
	"github.com/orilang/gori/lexer"
	"github.com/orilang/gori/parser"
	"github.com/orilang/gori/semantic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLower_Lower(t *testing.T) {
	t.Run("lower", func(t *testing.T) {
		tests := []struct {
			data     string
			expected string
		}{
			{
				data: `package main

func multiply(a int, b int) int {
    return a * b + b * a
}

func main() {
    x := multiply(int(1), int(2))
}`,
				expected: `func multiply(a:int, b:int) -> int
entry:
    t0 = mul_int a, b
    t1 = mul_int b, a
    t2 = add_int t0, t1
    return t2

func main()
entry:
    t0 = const_int 1
    t1 = const_int 2
    t2 = multiply(t0, t1)

`,
			},
			{
				data: `package main

func add(a int, b int) int {
    return a + b
}

func main() {
    x := add(int(1), int(2))
}`,
				expected: `func add(a:int, b:int) -> int
entry:
    t0 = add_int a, b
    return t0

func main()
entry:
    t0 = const_int 1
    t1 = const_int 2
    t2 = add(t0, t1)

`,
			},
		}

		for i, tc := range tests {
			lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
			require.NoError(t, err)
			parser := parser.New(lex.FetchTokensFromString(tc.data))
			pr := parser.ParseFile()
			require.Equal(t, 0, len(parser.Errors))
			check := semantic.NewChecker()

			spr, diagnostics := check.Check(pr)
			require.Equal(t, false, diagnostics.HasErrors())

			l := NewLower(true)
			lir, ldiagnostics := l.Lower(spr)
			require.Equal(t, false, ldiagnostics.HasErrors())

			assert.Equal(t, tc.expected, dump(lir), i)
		}
	})

	t.Run("lower_dummy_tests", func(t *testing.T) {
		tests := []struct {
			data     any
			expected ir.Value
			err      bool
		}{
			{
				data:     &semantic.IntLitExpr{Type: semantic.TInt, Value: "1"},
				expected: ir.Value("1"),
			},
			{
				data:     &semantic.BoolLitExpr{Type: semantic.TBool, Value: "true"},
				expected: ir.Value("true"),
			},
			{
				data:     &semantic.StringLitExpr{Type: semantic.TInt, Value: "yes"},
				expected: ir.Value("yes"),
			},
			{
				data:     &semantic.StringLitExpr{Type: semantic.TString, Value: "yes"},
				expected: ir.Value("yes"),
			},
			{
				data:     &semantic.FloatLitExpr{Type: semantic.TFloat, Value: "1.0"},
				expected: ir.Value("1.0"),
			},
			{
				data:     semantic.StringLitExpr{Type: semantic.TInt, Value: "yes"},
				expected: ir.Value(""),
				err:      true,
			},
		}

		for i, tc := range tests {
			l := NewLower(true)
			require.Equal(t, tc.expected, l.lower(tc.data), i)
			require.Equal(t, tc.err, len(l.errors) > 0, i)
		}
	})

	t.Run("Lower_error", func(t *testing.T) {
		l := NewLower(true)
		l.errors = append(l.errors, Diagnostic{Err: fmt.Errorf("dummy")})
		_, diag := l.Lower(semantic.Program{})
		require.Equal(t, true, diag.HasErrors())
	})

	t.Run("decl_error", func(t *testing.T) {
		l := NewLower(true)
		l.errors = append(l.errors, Diagnostic{Err: fmt.Errorf("dummy")})
		l.decl(&semantic.VarDecl{Name: "plop"})
		require.Equal(t, true, len(l.errors) > 0)
	})
}
