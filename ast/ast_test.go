package ast

import (
	"testing"

	"github.com/orilang/gori/token"
	"github.com/stretchr/testify/assert"
)

func TestAst_position(t *testing.T) {
	assert := assert.New(t)
	// The following are mostly dump tests as they are only to validate
	// the structed but not its values

	t.Run("ident_expr", func(t *testing.T) {
		z := token.Token{
			Kind:   token.Ident,
			Value:  "a",
			Line:   1,
			Column: 1,
		}
		x := &IdentExpr{z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("name_type", func(t *testing.T) {
		z := token.Token{
			Kind:   token.KWInt,
			Value:  "int",
			Line:   1,
			Column: 1,
		}
		x := &NameType{z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("int_lit_expr", func(t *testing.T) {
		z := token.Token{
			Kind:   token.IntLit,
			Value:  "1",
			Line:   1,
			Column: 1,
		}
		x := &IntLitExpr{z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("float_lit_expr", func(t *testing.T) {
		z := token.Token{
			Kind:   token.FloatLit,
			Value:  "3.14",
			Line:   1,
			Column: 1,
		}
		x := &FloatLitExpr{z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("bool_lit_expr", func(t *testing.T) {
		z := token.Token{
			Kind:   token.BoolLit,
			Value:  "true",
			Line:   1,
			Column: 1,
		}
		x := &BoolLitExpr{z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("string_lit_expr", func(t *testing.T) {
		z := token.Token{
			Kind:   token.StringLit,
			Value:  "string",
			Line:   1,
			Column: 1,
		}
		x := &StringLitExpr{z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("parent_expr", func(t *testing.T) {
		z := token.Token{
			Kind:   token.Ident,
			Value:  "a",
			Line:   1,
			Column: 1,
		}
		x := &ParenExpr{Left: z, Right: z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("binary_expr", func(t *testing.T) {
		z := token.Token{
			Kind:   token.Ident,
			Value:  "a",
			Line:   1,
			Column: 1,
		}
		x := &BinaryExpr{Left: &StringLitExpr{z}, Right: &StringLitExpr{z}}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("unary_expr", func(t *testing.T) {
		z := token.Token{
			Kind:   token.Plus,
			Value:  "+",
			Line:   1,
			Column: 1,
		}
		x := &UnaryExpr{Operator: z, Right: &StringLitExpr{z}}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("selector_expr", func(t *testing.T) {
		z := token.Token{
			Kind:   token.Ident,
			Value:  "a",
			Line:   1,
			Column: 1,
		}
		x := &SelectorExpr{X: &IdentExpr{z}, Selector: z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("index_expr", func(t *testing.T) {
		z := token.Token{
			Kind:   token.Ident,
			Value:  "a",
			Line:   1,
			Column: 1,
		}
		x := &IndexExpr{X: &IdentExpr{z}, RBracket: z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("call_expr", func(t *testing.T) {
		z := token.Token{
			Kind:   token.Ident,
			Value:  "a",
			Line:   1,
			Column: 1,
		}
		x := &CallExpr{Callee: &IdentExpr{z}, RParen: z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("expr_stmt", func(t *testing.T) {
		z := token.Token{
			Kind:   token.Ident,
			Value:  "a",
			Line:   1,
			Column: 1,
		}
		x := &ExprStmt{&IdentExpr{z}}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("bad_expr", func(t *testing.T) {
		z := token.Token{
			Kind:   token.Ident,
			Value:  "a",
			Line:   1,
			Column: 1,
		}
		x := &BadExpr{From: z, To: z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("dump_type", func(t *testing.T) {
		x := &dumpType{S: ""}
		assert.Equal(token.Token{}, x.Start())
		assert.Equal(token.Token{}, x.End())
	})

	t.Run("return_stmt", func(t *testing.T) {
		z := token.Token{
			Kind:   token.Ident,
			Value:  "a",
			Line:   1,
			Column: 1,
		}
		x := &ReturnStmt{Return: z, Values: []Expr{&IdentExpr{z}}}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("return_stmt_no_end", func(t *testing.T) {
		z := token.Token{
			Kind:   token.Ident,
			Value:  "a",
			Line:   1,
			Column: 1,
		}
		x := &ReturnStmt{Return: z}
		assert.Equal(z, x.Start())
		assert.Equal(token.Token{}, x.End())
	})

	t.Run("block_stmt", func(t *testing.T) {
		z := token.Token{
			Kind:   token.LBrace,
			Value:  "a",
			Line:   1,
			Column: 1,
		}
		x := &BlockStmt{LBrace: z, RBrace: z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("const_decl_stmt", func(t *testing.T) {
		z := token.Token{
			Kind:   token.KWConst,
			Value:  "a",
			Line:   1,
			Column: 1,
		}
		x := &ConstDeclStmt{ConstKW: z, Init: &IdentExpr{z}}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("var_decl_stmt", func(t *testing.T) {
		z := token.Token{
			Kind:   token.KWVar,
			Value:  "a",
			Line:   1,
			Column: 1,
		}
		x := &VarDeclStmt{VarKW: z, Init: &IdentExpr{z}}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("assign_stmt", func(t *testing.T) {
		z := token.Token{
			Kind:   token.KWVar,
			Value:  "a",
			Line:   1,
			Column: 1,
		}
		x := &AssignStmt{Left: &IdentExpr{z}, Right: &IdentExpr{z}}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("bad_stmt_from_to", func(t *testing.T) {
		z := token.Token{
			Kind:   token.KWVar,
			Value:  "a",
			Line:   1,
			Column: 1,
		}
		x := &BadStmt{From: z, To: z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("bad_stmt_from", func(t *testing.T) {
		z := token.Token{
			Kind:   token.KWVar,
			Value:  "a",
			Line:   1,
			Column: 1,
		}
		x := &BadStmt{From: z}
		assert.Equal(z, x.Start())
		assert.Equal(token.Token{}, x.End())
	})

	t.Run("bad_type_from_to", func(t *testing.T) {
		z := token.Token{
			Kind:   token.KWVar,
			Value:  "a",
			Line:   1,
			Column: 1,
		}
		x := &BadType{From: z, To: z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("bad_type_from", func(t *testing.T) {
		z := token.Token{
			Kind:   token.KWVar,
			Value:  "a",
			Line:   1,
			Column: 1,
		}
		x := &BadType{From: z}
		assert.Equal(z, x.Start())
		assert.Equal(token.Token{}, x.End())
	})

	t.Run("if_stmt_else_nil", func(t *testing.T) {
		z := token.Token{
			Kind:   token.KWIf,
			Value:  "if",
			Line:   1,
			Column: 1,
		}
		x := &IfStmt{If: z}
		assert.Equal(z, x.Start())
		assert.Equal(token.Token{}, x.End())
	})

	t.Run("if_stmt_else_not_nil", func(t *testing.T) {
		z := token.Token{
			Kind:   token.KWIf,
			Value:  "if",
			Line:   1,
			Column: 1,
		}
		x := &IfStmt{If: z, Else: &IdentExpr{z}}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("bad_decl_from_to", func(t *testing.T) {
		z := token.Token{
			Kind:   token.KWVar,
			Value:  "a",
			Line:   1,
			Column: 1,
		}
		x := &BadDecl{From: z, To: z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("bad_decl_from", func(t *testing.T) {
		z := token.Token{
			Kind:   token.KWVar,
			Value:  "a",
			Line:   1,
			Column: 1,
		}
		x := &BadDecl{From: z}
		assert.Equal(z, x.Start())
		assert.Equal(token.Token{}, x.End())
	})

	t.Run("for_stmt_from", func(t *testing.T) {
		z := token.Token{
			Kind:   token.KWFor,
			Value:  "for",
			Line:   1,
			Column: 1,
		}
		x := &ForStmt{For: z}
		assert.Equal(z, x.Start())
		assert.Equal(token.Token{}, x.End())
	})

	t.Run("for_stmt_from_to", func(t *testing.T) {
		z := token.Token{
			Kind:   token.KWFor,
			Value:  "for",
			Line:   1,
			Column: 1,
		}
		l := token.Token{
			Kind:   token.LBrace,
			Value:  "for",
			Line:   1,
			Column: 1,
		}
		r := token.Token{
			Kind:   token.RBrace,
			Value:  "for",
			Line:   1,
			Column: 1,
		}
		x := &ForStmt{For: z, Body: &BlockStmt{LBrace: l, RBrace: r}}
		assert.Equal(z, x.Start())
		assert.Equal(r, x.End())
	})

	t.Run("for_range_stmt_from", func(t *testing.T) {
		z := token.Token{
			Kind:   token.KWFor,
			Value:  "for",
			Line:   1,
			Column: 1,
		}
		x := &RangeStmt{For: z}
		assert.Equal(z, x.Start())
		assert.Equal(token.Token{}, x.End())
	})

	t.Run("for_range_stmt_from_to", func(t *testing.T) {
		z := token.Token{
			Kind:   token.KWFor,
			Value:  "for",
			Line:   1,
			Column: 1,
		}
		l := token.Token{
			Kind:   token.LBrace,
			Value:  "for",
			Line:   1,
			Column: 1,
		}
		r := token.Token{
			Kind:   token.RBrace,
			Value:  "for",
			Line:   1,
			Column: 1,
		}
		x := &RangeStmt{For: z, Body: &BlockStmt{LBrace: l, RBrace: r}}
		assert.Equal(z, x.Start())
		assert.Equal(r, x.End())
	})

	t.Run("for_range_stmt_from_to", func(t *testing.T) {
		z := token.Token{
			Kind:   token.Ident,
			Value:  "i",
			Line:   1,
			Column: 1,
		}
		op := token.Token{
			Kind:   token.PPlus,
			Value:  "++",
			Line:   1,
			Column: 3,
		}
		x := &IncDecStmt{X: &IdentExpr{z}, Operator: op}
		assert.Equal(z, x.Start())
		assert.Equal(op, x.End())
	})

	t.Run("break_stmt", func(t *testing.T) {
		z := token.Token{
			Kind:   token.KWBreak,
			Value:  "break",
			Line:   1,
			Column: 1,
		}
		x := &BreakStmt{Break: z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("continue_stmt", func(t *testing.T) {
		z := token.Token{
			Kind:   token.KWContinue,
			Value:  "continue",
			Line:   1,
			Column: 1,
		}
		x := &ContinueStmt{Continue: z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})

	t.Run("switch_stmt", func(t *testing.T) {
		z := token.Token{
			Kind:   token.KWSwitch,
			Value:  "switch",
			Line:   1,
			Column: 1,
		}
		r := token.Token{
			Kind:   token.RBrace,
			Value:  "}",
			Line:   1,
			Column: 3,
		}
		x := &SwitchStmt{Switch: z, RBrace: r}
		assert.Equal(z, x.Start())
		assert.Equal(r, x.End())
	})

	t.Run("case_clause_stmt_x1", func(t *testing.T) {
		z := token.Token{
			Kind:   token.KWCase,
			Value:  "case",
			Line:   1,
			Column: 1,
		}
		x := &CaseClause{Case: z}
		assert.Equal(z, x.Start())
		assert.Equal(token.Token{}, x.End())
	})

	t.Run("case_clause_stmt_x2", func(t *testing.T) {
		z := token.Token{
			Kind:   token.KWCase,
			Value:  "case",
			Line:   1,
			Column: 1,
		}
		rt := token.Token{
			Kind:   token.KWReturn,
			Value:  "return",
			Line:   2,
			Column: 1,
		}

		rtv := token.Token{
			Kind:   token.Ident,
			Value:  "a",
			Line:   2,
			Column: 8,
		}
		x := &CaseClause{Case: z}
		rtstmt := &ReturnStmt{Return: rt}
		rtstmt.Values = append(rtstmt.Values, &IdentExpr{rtv})
		x.Body = append(x.Body, rtstmt)
		assert.Equal(z, x.Start())
		assert.Equal(rtv, x.End())
	})

	t.Run("fallthrough_stmt", func(t *testing.T) {
		z := token.Token{
			Kind:   token.KWFallThrough,
			Value:  "continue",
			Line:   1,
			Column: 1,
		}
		x := &FallThroughStmt{FallThroughStmt: z}
		assert.Equal(z, x.Start())
		assert.Equal(z, x.End())
	})
}
