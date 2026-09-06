// lower package simplifies semantic symbols into n Lower Intermediate representation (LIR).
// In other terms, it defines how an AST becomes an IR.
package lower

import (
	"fmt"
	"slices"

	"github.com/orilang/gori/ir"
	"github.com/orilang/gori/semantic"
	"github.com/orilang/gori/token"
)

// NewLower initialize Lower-level Intermediate Representation requirements
func NewLower(output bool) *Lower {
	return &Lower{
		output: output,
	}
}

func (l *Lower) Lower(p semantic.Program) (ir.Program, Diagnostics) {
	for _, v := range p.Files {
		l.decls(v.Decls)
	}

	if len(l.errors) == 0 {
		pgm := ir.Program{}
		pgm.Funcs = append(pgm.Funcs, l.funcs...)

		if l.output {
			fmt.Printf("%s\n", dump(pgm))
		}

		return pgm, l.errors
	}
	return ir.Program{}, l.errors
}

// HasErrors returns true when errors found
func (d Diagnostics) HasErrors() bool {
	return len(d) > 0
}

func (l *Lower) decls(decl []semantic.Decl) {
	for _, v := range decl {
		l.decl(v)
	}
}

func (l *Lower) decl(decl semantic.Decl) {
	switch t := decl.(type) {
	case *semantic.FuncDecl:
		l.funcs = append(l.funcs, l.fn(t))

	default:
		l.errors = append(l.errors, Diagnostic{Err: fmt.Errorf("unsupported declaration %T", t)})
	}
}

func (l *Lower) fn(decl *semantic.FuncDecl) *ir.Func {
	oldBlockName := l.blockName
	oldblocks := l.blocks
	oldInstructions := slices.Clone(l.instructions)
	oldTIndex := l.tIndex
	oldLabelIndex := l.labelIndex
	defer func() {
		l.blockName = oldBlockName
		l.blocks = oldblocks
		l.instructions = oldInstructions
		l.tIndex = oldTIndex
		l.labelIndex = oldLabelIndex
	}()
	l.blockName = "entry"

	f := &ir.Func{Name: decl.Name}
	var param, result []ir.Param

	for _, p := range decl.Params {
		param = append(param, ir.Param{Name: p.Name, Type: p.Type.String()})
	}
	f.Params = param

	for _, p := range decl.Results {
		result = append(result, ir.Param{Name: p.Name, Type: p.Type.String()})
	}
	f.Results = result

	if decl.Body != nil {
		f.Label = "entry"
		_ = l.lower(decl.Body.Stmts)
		f.Blocks = l.blocks
	}
	return f
}

// lower lowers any input stmt/expr to later create an instruction
func (l *Lower) lower(input any) ir.Value {
	switch in := input.(type) {
	case []semantic.Stmt:
		for _, st := range in {
			_ = l.lower(st)
			l.blocks = append(l.blocks, &ir.Block{
				Instructions: slices.Clone(l.instructions),
			})
			l.instructions = nil
		}

	case semantic.Stmt:
		return l.lowerStmt(in)

	case semantic.Expr:
		return l.lowerExpr(in)

	default:
		l.errors = append(l.errors, Diagnostic{Err: fmt.Errorf("unsupported input %T", in)})
	}

	return ir.Value("")
}

// lowerStmt lowers expr to later create an instruction
func (l *Lower) lowerStmt(t semantic.Stmt) ir.Value {
	switch stmt := t.(type) {
	case *semantic.ReturnStmt:
		if len(stmt.Values) == 0 {
			l.instructions = append(l.instructions, &ir.Return{})
			return ir.Value("")
		}

		var rt []ir.Value
		for _, v := range stmt.Values {
			st := l.lower(v)
			rt = append(rt, st)
			l.instructions = append(l.instructions, &ir.Return{
				Name: string(st),
			})
		}

		if len(rt) == 1 {
			return rt[0]
		}

	case *semantic.AssigmentStmt:
		t := l.lower(stmt.Right)
		l.instructions = append(l.instructions, &ir.Assigment{Result: stmt.Symbol.Name, Value: string(t)})
		return ir.Value(stmt.Symbol.Name)

	case *semantic.IfStmt:
		labelIndex := l.labelIndex
		l.labelIndex++

		cond := l.lower(stmt.Condition)
		br := &ir.Branch{Condition: string(cond), True: "if_then", False: "if_else", Index: labelIndex}
		if len(stmt.Else) > 0 {
			br.False = "if_else"
		} else {
			br.False = "if_end"
		}
		l.instructions = append(l.instructions, br)

		l.instructions = append(l.instructions, &ir.Label{Name: "if_then", Index: labelIndex})
		jump := &ir.Jump{Name: "if_end", Index: labelIndex}
		end := &ir.Label{Name: "if_end", Index: labelIndex}

		var thenTerminating bool
		for _, v := range stmt.Then {
			thenTerminating = isTerminatingStmt(v)
			_ = l.lower(v)
		}

		if !thenTerminating {
			l.instructions = append(l.instructions, jump)
		}

		var elseTerminating bool
		if len(stmt.Else) > 0 {
			l.instructions = append(l.instructions, &ir.Label{Name: "if_else", Index: labelIndex})
			for _, v := range stmt.Else {
				elseTerminating = isTerminatingStmt(v)
				_ = l.lower(v)
			}

			if !elseTerminating {
				l.instructions = append(l.instructions, jump)
			}
		}

		if !(thenTerminating && elseTerminating) {
			l.instructions = append(l.instructions, end)
		}

	default:
		l.errors = append(l.errors, Diagnostic{Err: fmt.Errorf("unsupported statement %T", stmt)})
	}

	return ir.Value("")
}

// lowerExpr lowers expr to later create an instruction
func (l *Lower) lowerExpr(t semantic.Expr) ir.Value {
	switch expr := t.(type) {
	case *semantic.IntLitExpr:
		return ir.Value(expr.Value)

	case *semantic.FloatLitExpr:
		return ir.Value(expr.Value)

	case *semantic.BoolLitExpr:
		return ir.Value(expr.Value)

	case *semantic.StringLitExpr:
		return ir.Value(expr.Value)

	case *semantic.IdentExpr:
		return ir.Value(expr.Value)

	case *semantic.BinaryExpr:
		left := l.lower(expr.Left)
		right := l.lower(expr.Right)
		t := fmt.Sprintf("t%d", l.tIndex)

		l.instructions = append(l.instructions, &ir.Binary{
			Result: t,
			Op:     token.BinaryOpString(expr.Operator) + "_" + expr.Type.String(),
			Left:   string(left),
			Right:  string(right),
		})

		l.tIndex++
		return ir.Value(t)

	case *semantic.ConversionExpr:
		value := l.lower(expr.Value)
		t := fmt.Sprintf("t%d", l.tIndex)
		l.instructions = append(l.instructions, &ir.Const{
			Result: t,
			Type:   expr.To.String(),
			Value:  string(value),
		})

		l.tIndex++
		return ir.Value(t)

	case *semantic.CallExpr:
		var args []string
		for _, v := range expr.Args {
			args = append(args, string(l.lower(v)))
		}

		callee := l.lower(expr.Callee)
		t := fmt.Sprintf("t%d", l.tIndex)
		l.instructions = append(l.instructions, &ir.Call{
			Result: t,
			Name:   string(callee),
			Args:   args,
		})

		l.tIndex++
		return ir.Value(t)

	default:
		l.errors = append(l.errors, Diagnostic{Err: fmt.Errorf("unsupported expression %T", expr)})
	}

	return ir.Value("")
}

// isTerminatingStmt returns true if statement is return/break/continue
func isTerminatingStmt(stmt semantic.Stmt) bool {
	if _, ok := stmt.(*semantic.ReturnStmt); ok {
		return true
	}
	return false
}
