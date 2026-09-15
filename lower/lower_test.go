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
	t.Run("x1", func(t *testing.T) {
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
    x = t2

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
    x = t2

`,
			},
			{
				data: `package main
func add(a int, b int) (c int) {
    return a + b
}

func main() {
    x := add(int(1), int(2))
}`,
				expected: `func add(a:int, b:int) -> (c:int)
entry:
    t0 = add_int a, b
    return t0

func main()
entry:
    t0 = const_int 1
    t1 = const_int 2
    t2 = add(t0, t1)
    x = t2

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

	t.Run("x2", func(t *testing.T) {
		tests := []struct {
			data     string
			expected string
		}{
			{
				data: `package main
func f(a int, b int) (c int) {
    if a > 0 {
      return a
		}
    return b
}

func main() {
    x := f(int(1), int(2))
}`,
				expected: `func f(a:int, b:int) -> (c:int)
entry:
    t0 = gt_bool a, 0
    branch t0, if_then_0, if_end_0

if_then_0:
    return a

if_end_0:
    return b

func main()
entry:
    t0 = const_int 1
    t1 = const_int 2
    t2 = f(t0, t1)
    x = t2

`,
			},
			{
				data: `package main
func f(a int, b int) (c int) {
    if a > 0 {
      return a
		} else {
      return b
    }
}

func main() {
    x := f(int(1), int(2))
}`,
				expected: `func f(a:int, b:int) -> (c:int)
entry:
    t0 = gt_bool a, 0
    branch t0, if_then_0, if_else_0

if_then_0:
    return a

if_else_0:
    return b

func main()
entry:
    t0 = const_int 1
    t1 = const_int 2
    t2 = f(t0, t1)
    x = t2

`,
			},
			{
				data: `package main
func add(a int, b int) int {
    return a + b
}

func main() {
    a := float64(1)
    x := add(int(a), int(2))
}`,
				expected: `func add(a:int, b:int) -> int
entry:
    t0 = add_int a, b
    return t0

func main()
entry:
    t0 = const_float64 1
    a = t0
    t1 = const_int a
    t2 = const_int 2
    t3 = add(t1, t2)
    x = t3

`,
			},
			{
				data: `package main
func f(a int, b int) (c int) {
    if a > int(0) {
      return a
    } else {
      x := int(1)
      return b
    }
}

func main() {
    x := f(int(1), int(2))
}`,
				expected: `func f(a:int, b:int) -> (c:int)
entry:
    t0 = const_int 0
    t1 = gt_bool a, t0
    branch t1, if_then_0, if_else_0

if_then_0:
    return a

if_else_0:
    t2 = const_int 1
    x = t2
    return b

func main()
entry:
    t0 = const_int 1
    t1 = const_int 2
    t2 = f(t0, t1)
    x = t2

`,
			},
			{
				data: `package main
func f(a int, b int) (c int) {
    if a > 0 {
      return a
    }
    return b
}

func main() {
    x := f(int(1), int(2))
}`,
				expected: `func f(a:int, b:int) -> (c:int)
entry:
    t0 = gt_bool a, 0
    branch t0, if_then_0, if_end_0

if_then_0:
    return a

if_end_0:
    return b

func main()
entry:
    t0 = const_int 1
    t1 = const_int 2
    t2 = f(t0, t1)
    x = t2

`,
			},
			{
				data: `package main
func f(a int, b int) (c int) {
    if a > int(0) {
      if a < int(2) {
        return int(1)
      }
      return int(0)
    } else {
      x := int(1)
      return b
    }
}

func main() {
    x := f(int(1), int(2))
}`,
				expected: `func f(a:int, b:int) -> (c:int)
entry:
    t0 = const_int 0
    t1 = gt_bool a, t0
    branch t1, if_then_0, if_else_0

if_then_0:
    t2 = const_int 2
    t3 = lt_bool a, t2
    branch t3, if_then_1, if_end_1

if_then_1:
    t4 = const_int 1
    return t4

if_end_1:
    t5 = const_int 0
    return t5

if_else_0:
    t6 = const_int 1
    x = t6
    return b

func main()
entry:
    t0 = const_int 1
    t1 = const_int 2
    t2 = f(t0, t1)
    x = t2

`,
			},
			{
				data: `package main
func f(a int, b int) (c int) {
    c = int(1)
    if a > int(0) {
      c = a
    } else {
      return b
    }
    return c
}

func main() {
    x := f(int(1), int(2))
}`,
				expected: `func f(a:int, b:int) -> (c:int)
entry:
    t0 = const_int 1
    c = t0
    t1 = const_int 0
    t2 = gt_bool a, t1
    branch t2, if_then_0, if_else_0

if_then_0:
    c = a
    jump if_end_0

if_else_0:
    return b

if_end_0:
    return c

func main()
entry:
    t0 = const_int 1
    t1 = const_int 2
    t2 = f(t0, t1)
    x = t2

`,
			},
			{
				data: `package main
func f(a int, b int) (c int) {
    if a > 0 {
      c = a
    } else {
      c = b
    }
    return c
}

func main() {
    x := f(int(1), int(2))
}`,
				expected: `func f(a:int, b:int) -> (c:int)
entry:
    t0 = gt_bool a, 0
    branch t0, if_then_0, if_else_0

if_then_0:
    c = a
    jump if_end_0

if_else_0:
    c = b
    jump if_end_0

if_end_0:
    return c

func main()
entry:
    t0 = const_int 1
    t1 = const_int 2
    t2 = f(t0, t1)
    x = t2

`,
			},
			{
				data: `package main
func f(a int, b int) (c int) {
    if a > 0 {
      c = a
    } else {
      c = b
    }
    return
}

func main() {
    x := f(int(1), int(2))
}`,
				expected: `func f(a:int, b:int) -> (c:int)
entry:
    t0 = gt_bool a, 0
    branch t0, if_then_0, if_else_0

if_then_0:
    c = a
    jump if_end_0

if_else_0:
    c = b
    jump if_end_0

if_end_0:
    return

func main()
entry:
    t0 = const_int 1
    t1 = const_int 2
    t2 = f(t0, t1)
    x = t2

`,
			},
			{
				data: `package main
func f(a int, b int) (c int) {
    if a > 0 {
      c = a
    } else if a < 0 {
      c = - int(1)
    } else {
      c = b
    }
    return
}

func main() {
    x := f(int(1), int(2))
}`,
				expected: `func f(a:int, b:int) -> (c:int)
entry:
    t0 = gt_bool a, 0
    branch t0, if_then_0, if_else_0

if_then_0:
    c = a
    jump if_end_0

if_else_0:
    t1 = lt_bool a, 0
    branch t1, if_then_1, if_else_1

if_then_1:
    t2 = const_int 1
    t3 = -t2
    c = t3
    jump if_end_1

if_else_1:
    c = b
    jump if_end_1

if_end_1:
    jump if_end_0

if_end_0:
    return

func main()
entry:
    t0 = const_int 1
    t1 = const_int 2
    t2 = f(t0, t1)
    x = t2

`,
			},
			{
				data: `package main
func f(a int, b int) (c int) {
    if a > 0 {
      c = a
    } else if a < 0 {
      c = int(-1)
    } else {
      c = b
    }
    return
}

func main() {
    x := f(int(1), int(2))
}`,
				expected: `func f(a:int, b:int) -> (c:int)
entry:
    t0 = gt_bool a, 0
    branch t0, if_then_0, if_else_0

if_then_0:
    c = a
    jump if_end_0

if_else_0:
    t1 = lt_bool a, 0
    branch t1, if_then_1, if_else_1

if_then_1:
    t2 = -1
    t3 = const_int t2
    c = t3
    jump if_end_1

if_else_1:
    c = b
    jump if_end_1

if_end_1:
    jump if_end_0

if_end_0:
    return

func main()
entry:
    t0 = const_int 1
    t1 = const_int 2
    t2 = f(t0, t1)
    x = t2

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

	t.Run("dummy_lower_tests", func(t *testing.T) {
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

	t.Run("dummy_Lower_error", func(t *testing.T) {
		l := NewLower(true)
		l.errors = append(l.errors, Diagnostic{Err: fmt.Errorf("dummy")})
		_, diag := l.Lower(semantic.Program{})
		require.Equal(t, true, diag.HasErrors())
	})

	t.Run("dummy_decl_error", func(t *testing.T) {
		l := NewLower(true)
		l.errors = append(l.errors, Diagnostic{Err: fmt.Errorf("dummy")})
		l.decl(&semantic.VarDecl{Name: "plop"})
		require.Equal(t, true, len(l.errors) > 0)
	})

	t.Run("dummy_expr_error", func(t *testing.T) {
		l := NewLower(true)
		l.lowerExpr(nil)
		require.Equal(t, true, len(l.errors) > 0)
	})

	t.Run("dummy_stmt_error", func(t *testing.T) {
		l := NewLower(true)
		l.lowerStmt(nil)
		require.Equal(t, true, len(l.errors) > 0)
	})

	t.Run("x3", func(t *testing.T) {
		tests := []struct {
			data     string
			expected string
		}{
			{
				data: `package main

func main() {
  a := int(0)

  switch a {
  case int(1), int(2):
    a = int(1)
  case int(3):
    a = int(1)
  }

  a = int(2)
}`,
				expected: `func main()
entry:
    t0 = const_int 0
    a = t0

switch_0_check_1:
    t2 = const_int 1
    t1 = eq_bool a, t2
    branch t1, switch_0_case_1, switch_0_check_2

switch_0_check_2:
    t3 = const_int 2
    t1 = eq_bool a, t3
    branch t1, switch_0_case_1, switch_0_check_3

switch_0_check_3:
    t4 = const_int 3
    t1 = eq_bool a, t4
    branch t1, switch_0_case_3, switch_0_end

switch_0_case_1:
    t5 = const_int 1
    a = t5
    jump switch_0_end

switch_0_case_3:
    t6 = const_int 1
    a = t6
    jump switch_0_end

switch_0_end:
    t7 = const_int 2
    a = t7

`,
			},
			{
				data: `package main
func foo() int {
  return int(0)
}

func bar() int {
  return int(1)
}

func baz() int {
  return int(2)
}

func main() {
  a := int(0)
  b := int(0)

  switch a {
    case int(1):
      b = foo()

    case int(2), int(3):
      b = bar()

    default:
      b = baz()
  }
}`,
				expected: `func foo() -> int
entry:
    t0 = const_int 0
    return t0

func bar() -> int
entry:
    t0 = const_int 1
    return t0

func baz() -> int
entry:
    t0 = const_int 2
    return t0

func main()
entry:
    t0 = const_int 0
    a = t0
    t1 = const_int 0
    b = t1

switch_0_check_1:
    t3 = const_int 1
    t2 = eq_bool a, t3
    branch t2, switch_0_case_1, switch_0_check_2

switch_0_check_2:
    t4 = const_int 2
    t2 = eq_bool a, t4
    branch t2, switch_0_case_2, switch_0_check_3

switch_0_check_3:
    t5 = const_int 3
    t2 = eq_bool a, t5
    branch t2, switch_0_case_2, switch_0_default

switch_0_case_1:
    t6 = foo()
    b = t6
    jump switch_0_end

switch_0_case_2:
    t7 = bar()
    b = t7
    jump switch_0_end

switch_0_default:
    t8 = baz()
    b = t8
    jump switch_0_end

switch_0_end:

`,
			},
			{
				data: `package main

func f(a int) int {
  return a * 2
}

func main() {
  a := int(0)

  switch z := f(a);z {
  case 1:
    fallthrough
  default:
    a = int(1)
  }
}`,
				expected: `func f(a:int) -> int
entry:
    t0 = mul_int a, 2
    return t0

func main()
entry:
    t0 = const_int 0
    a = t0
    t1 = f(a)
    z = t1

switch_0_check_1:
    t2 = eq_bool z, 1
    branch t2, switch_0_case_1, switch_0_default

switch_0_case_1:
    jump switch_0_default

switch_0_default:
    t3 = const_int 1
    a = t3
    jump switch_0_end

switch_0_end:

`,
			},
			{
				data: `package main

func main() {
  a := int(0)

  switch a {
  case 1:
    a = int(1)
  default:
    a = int(1)
  }

  a = int(2)
}`,
				expected: `func main()
entry:
    t0 = const_int 0
    a = t0

switch_0_check_1:
    t1 = eq_bool a, 1
    branch t1, switch_0_case_1, switch_0_default

switch_0_case_1:
    t2 = const_int 1
    a = t2
    jump switch_0_end

switch_0_default:
    t3 = const_int 1
    a = t3
    jump switch_0_end

switch_0_end:
    t4 = const_int 2
    a = t4

`,
			},
			{
				data: `package main

func f(a int) int {
  return a * 2
}

func main() {
  a := int(0)

  switch z := f(a);z {
  case 1:
    a = f(a)*2
  case 2, 3:
    fallthrough
  case 4, 5:
    a = f(a)*3
  default:
    a = int(1)
  }
}`,
				expected: `func f(a:int) -> int
entry:
    t0 = mul_int a, 2
    return t0

func main()
entry:
    t0 = const_int 0
    a = t0
    t1 = f(a)
    z = t1

switch_0_check_1:
    t2 = eq_bool z, 1
    branch t2, switch_0_case_1, switch_0_check_2

switch_0_check_2:
    t2 = eq_bool z, 2
    branch t2, switch_0_case_2, switch_0_check_3

switch_0_check_3:
    t2 = eq_bool z, 3
    branch t2, switch_0_case_2, switch_0_check_4

switch_0_check_4:
    t2 = eq_bool z, 4
    branch t2, switch_0_case_4, switch_0_check_5

switch_0_check_5:
    t2 = eq_bool z, 5
    branch t2, switch_0_case_4, switch_0_default

switch_0_case_1:
    t3 = f(a)
    t4 = mul_int t3, 2
    a = t4
    jump switch_0_end

switch_0_case_2:
    jump switch_0_case_4

switch_0_case_4:
    t5 = f(a)
    t6 = mul_int t5, 3
    a = t6
    jump switch_0_end

switch_0_default:
    t7 = const_int 1
    a = t7
    jump switch_0_end

switch_0_end:

`,
			},
			{
				data: `package main

func f(a int) int {
  return a * 2
}

func main() {
  a := int(0)

  switch z := f(a);z {
  case 2, 3:
    fallthrough
  case 1:
    a = f(a)*2
  default:
    a = int(1)
  }
}`,
				expected: `func f(a:int) -> int
entry:
    t0 = mul_int a, 2
    return t0

func main()
entry:
    t0 = const_int 0
    a = t0
    t1 = f(a)
    z = t1

switch_0_check_1:
    t2 = eq_bool z, 2
    branch t2, switch_0_case_1, switch_0_check_2

switch_0_check_2:
    t2 = eq_bool z, 3
    branch t2, switch_0_case_1, switch_0_check_3

switch_0_check_3:
    t2 = eq_bool z, 1
    branch t2, switch_0_case_3, switch_0_default

switch_0_case_1:
    jump switch_0_case_3

switch_0_case_3:
    t3 = f(a)
    t4 = mul_int t3, 2
    a = t4
    jump switch_0_end

switch_0_default:
    t5 = const_int 1
    a = t5
    jump switch_0_end

switch_0_end:

`,
			},
			{
				data: `package main

func f(a int) int {
  return a * 2
}

func main() {
  a := int(0)

  switch z := f(a);z {
  case 1:
  case 2, 3:
    a = 2
    fallthrough
  case 4, 5:
    a = f(2)
    a = a+1
    fallthrough
  case 6, 7:
    a = f(3)
    a = a+1
  default:
    a = int(1)
  }
}`,
				expected: `func f(a:int) -> int
entry:
    t0 = mul_int a, 2
    return t0

func main()
entry:
    t0 = const_int 0
    a = t0
    t1 = f(a)
    z = t1

switch_0_check_1:
    t2 = eq_bool z, 1
    branch t2, switch_0_case_1, switch_0_check_2

switch_0_check_2:
    t2 = eq_bool z, 2
    branch t2, switch_0_case_2, switch_0_check_3

switch_0_check_3:
    t2 = eq_bool z, 3
    branch t2, switch_0_case_2, switch_0_check_4

switch_0_check_4:
    t2 = eq_bool z, 4
    branch t2, switch_0_case_4, switch_0_check_5

switch_0_check_5:
    t2 = eq_bool z, 5
    branch t2, switch_0_case_4, switch_0_check_6

switch_0_check_6:
    t2 = eq_bool z, 6
    branch t2, switch_0_case_6, switch_0_check_7

switch_0_check_7:
    t2 = eq_bool z, 7
    branch t2, switch_0_case_6, switch_0_default

switch_0_case_1:
    jump switch_0_end

switch_0_case_2:
    a = 2
    jump switch_0_case_4

switch_0_case_4:
    t3 = f(2)
    a = t3
    t4 = add_int a, 1
    a = t4
    jump switch_0_case_6

switch_0_case_6:
    t5 = f(3)
    a = t5
    t6 = add_int a, 1
    a = t6
    jump switch_0_end

switch_0_default:
    t7 = const_int 1
    a = t7
    jump switch_0_end

switch_0_end:

`,
			},
			{
				data: `package main

func f(a int) int {
  return a * 2
}

func main() {
  a := int(0)

  switch {
  case a == 1:
  case a == 2:
    a = 2
    fallthrough
  case a == 3:
    a = f(2)
    a = a+1
    fallthrough
  case a == 4:
    a = f(3)
    a = a+1
  default:
    a = int(1)
  }
}`,
				expected: `func f(a:int) -> int
entry:
    t0 = mul_int a, 2
    return t0

func main()
entry:
    t0 = const_int 0
    a = t0

switch_0_check_1:
    t1 = eq_bool a, 1
    branch t1, switch_0_case_1, switch_0_check_2

switch_0_check_2:
    t2 = eq_bool a, 2
    branch t2, switch_0_case_2, switch_0_check_3

switch_0_check_3:
    t3 = eq_bool a, 3
    branch t3, switch_0_case_3, switch_0_check_4

switch_0_check_4:
    t4 = eq_bool a, 4
    branch t4, switch_0_case_4, switch_0_default

switch_0_case_1:
    jump switch_0_end

switch_0_case_2:
    a = 2
    jump switch_0_case_3

switch_0_case_3:
    t5 = f(2)
    a = t5
    t6 = add_int a, 1
    a = t6
    jump switch_0_case_4

switch_0_case_4:
    t7 = f(3)
    a = t7
    t8 = add_int a, 1
    a = t8
    jump switch_0_end

switch_0_default:
    t9 = const_int 1
    a = t9
    jump switch_0_end

switch_0_end:

`,
			},
			{
				data: `package main

func f(a int) int {
  return a * 2
}

func main() {
  a := int(0)

  switch {
  case a == 1:
  case a == 2:
    a = 2
    fallthrough
  case a == 3:
    a = f(2)
    a = a+1
    fallthrough
  default:
    a = int(1)
  case a == 4:
    a = f(3)
    a = a+1
  }
}`,
				expected: `func f(a:int) -> int
entry:
    t0 = mul_int a, 2
    return t0

func main()
entry:
    t0 = const_int 0
    a = t0

switch_0_check_1:
    t1 = eq_bool a, 1
    branch t1, switch_0_case_1, switch_0_check_2

switch_0_check_2:
    t2 = eq_bool a, 2
    branch t2, switch_0_case_2, switch_0_check_3

switch_0_check_3:
    t3 = eq_bool a, 3
    branch t3, switch_0_case_3, switch_0_check_5

switch_0_check_5:
    t4 = eq_bool a, 4
    branch t4, switch_0_case_5, switch_0_default

switch_0_case_1:
    jump switch_0_end

switch_0_case_2:
    a = 2
    jump switch_0_case_3

switch_0_case_3:
    t5 = f(2)
    a = t5
    t6 = add_int a, 1
    a = t6
    jump switch_0_default

switch_0_default:
    t7 = const_int 1
    a = t7
    jump switch_0_end

switch_0_case_5:
    t8 = f(3)
    a = t8
    t9 = add_int a, 1
    a = t9
    jump switch_0_end

switch_0_end:

`,
			},
			{
				data: `package main

func f(a int) int {
  return a * 2
}

func main() {
  a := int(0)

  switch {
  case a == 1:
  case a == 2:
    a = 2
    fallthrough
  case a == 3:
    a = f(2)
    a = a+1
  default:
    a = int(1)
  case a == 4:
    a = f(3)
    a = a+1
  }
}`,
				expected: `func f(a:int) -> int
entry:
    t0 = mul_int a, 2
    return t0

func main()
entry:
    t0 = const_int 0
    a = t0

switch_0_check_1:
    t1 = eq_bool a, 1
    branch t1, switch_0_case_1, switch_0_check_2

switch_0_check_2:
    t2 = eq_bool a, 2
    branch t2, switch_0_case_2, switch_0_check_3

switch_0_check_3:
    t3 = eq_bool a, 3
    branch t3, switch_0_case_3, switch_0_check_5

switch_0_check_5:
    t4 = eq_bool a, 4
    branch t4, switch_0_case_5, switch_0_default

switch_0_case_1:
    jump switch_0_end

switch_0_case_2:
    a = 2
    jump switch_0_case_3

switch_0_case_3:
    t5 = f(2)
    a = t5
    t6 = add_int a, 1
    a = t6
    jump switch_0_end

switch_0_default:
    t7 = const_int 1
    a = t7
    jump switch_0_end

switch_0_case_5:
    t8 = f(3)
    a = t8
    t9 = add_int a, 1
    a = t9
    jump switch_0_end

switch_0_end:

`,
			},
			{
				data: `package main

func f(a int) int {
  return a * 2
}

func main() {
  a := int(0)

  switch z := f(a);z {
  case 1:
  default:
    a = int(1)
  case 2, 3:
    a = 2
    fallthrough
  case 4, 5:
    a = f(2)
    a = a+1
    fallthrough
  case 6, 7:
    a = f(3)
    a = a+1
  }
}`,
				expected: `func f(a:int) -> int
entry:
    t0 = mul_int a, 2
    return t0

func main()
entry:
    t0 = const_int 0
    a = t0
    t1 = f(a)
    z = t1

switch_0_check_1:
    t2 = eq_bool z, 1
    branch t2, switch_0_case_1, switch_0_check_3

switch_0_check_3:
    t2 = eq_bool z, 2
    branch t2, switch_0_case_3, switch_0_check_4

switch_0_check_4:
    t2 = eq_bool z, 3
    branch t2, switch_0_case_3, switch_0_check_5

switch_0_check_5:
    t2 = eq_bool z, 4
    branch t2, switch_0_case_5, switch_0_check_6

switch_0_check_6:
    t2 = eq_bool z, 5
    branch t2, switch_0_case_5, switch_0_check_7

switch_0_check_7:
    t2 = eq_bool z, 6
    branch t2, switch_0_case_7, switch_0_check_8

switch_0_check_8:
    t2 = eq_bool z, 7
    branch t2, switch_0_case_7, switch_0_default

switch_0_case_1:
    jump switch_0_end

switch_0_default:
    t3 = const_int 1
    a = t3
    jump switch_0_end

switch_0_case_3:
    a = 2
    jump switch_0_case_5

switch_0_case_5:
    t4 = f(2)
    a = t4
    t5 = add_int a, 1
    a = t5
    jump switch_0_case_7

switch_0_case_7:
    t6 = f(3)
    a = t6
    t7 = add_int a, 1
    a = t7
    jump switch_0_end

switch_0_end:

`,
			},
			{
				data: `package main

func main() {
  a := int(0)

  switch a {
  default:
    a = int(1)
    fallthrough
  case 1:
    a = int(1)
  }

  a = int(2)
}`,
				expected: `func main()
entry:
    t0 = const_int 0
    a = t0

switch_0_check_2:
    t1 = eq_bool a, 1
    branch t1, switch_0_case_2, switch_0_default

switch_0_default:
    t2 = const_int 1
    a = t2
    jump switch_0_case_2

switch_0_case_2:
    t3 = const_int 1
    a = t3
    jump switch_0_end

switch_0_end:
    t4 = const_int 2
    a = t4

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
}
