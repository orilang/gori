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
	case *semantic.FallThroughStmt:
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
		br := &ir.Branch{Condition: string(cond)}
		br.List = append(br.List, ir.BranchSub{Name: "if_then", Index: labelIndex})

		if len(stmt.Else) > 0 {
			br.List = append(br.List, ir.BranchSub{Name: "if_else", Index: labelIndex})
		} else {
			br.List = append(br.List, ir.BranchSub{Name: "if_end", Index: labelIndex})
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

	case *semantic.SwitchStmt:
		labelIndex := l.labelIndex
		l.labelIndex++

		type branch struct {
			name  string
			index int
			dft   bool
		}

		type swLabel struct {
			label      string
			branches   []branch
			dft        bool
			index      int
			isMulti    bool
			multiIndex int
			value      semantic.Expr
		}

		type next struct {
			jump  string
			index int
			dft   bool
		}

		type swCase struct {
			dft   bool
			casee string
			index int
			body  []semantic.Stmt
			next  next
		}

		var (
			swLabels   []swLabel
			swCases    []swCase
			nextJump   bool
			hasDefault bool
		)

		idx := 1
		for _, sc := range stmt.Cases {
			swc := swLabel{}

			switch len(sc.Values) {
			case 0:
				hasDefault = true
				swc.dft = true
				swc.label = fmt.Sprintf("switch_%d_default", labelIndex)
				br := branch{name: swc.label, dft: true}
				swc.branches = append(swc.branches, br)
				lgl := len(swLabels)
				if lgl > 0 {
					swLabels[lgl-1].branches = append(swLabels[lgl-1].branches, br)
				}
				swLabels = append(swLabels, swc)

				lgh := len(swCases)
				if nextJump && lgh > 0 {
					swCases[lgh-1].next.jump = swc.label
					swCases[lgh-1].next.index = idx
					swCases[lgh-1].next.dft = true
				}

				swCases = append(swCases, swCase{
					dft:   true,
					casee: swc.label,
					body:  sc.Body,
				})
				idx++

			case 1:
				swc.label = fmt.Sprintf("switch_%d_check", labelIndex)
				swc.index = idx
				swc.value = sc.Values[0]
				casee := fmt.Sprintf("switch_%d_case", labelIndex)
				br := branch{name: casee, index: idx}
				swc.branches = append(swc.branches, br)
				lgl := len(swLabels)
				if lgl > 0 {
					br.name = swc.label
					swLabels[lgl-1].branches = append(swLabels[lgl-1].branches, br)
				}
				swLabels = append(swLabels, swc)

				lgh := len(swCases)
				if nextJump && lgh > 0 {
					swCases[lgh-1].next.jump = casee
					swCases[lgh-1].next.index = idx
				}

				swCases = append(swCases, swCase{
					casee: casee,
					index: idx,
					body:  sc.Body,
				})

				if isFallingThroughStmt(sc.Body) {
					nextJump = true
				} else {
					nextJump = false
				}

				idx++

			default:
				swc.isMulti = true
				swc.multiIndex = idx

				casee := fmt.Sprintf("switch_%d_case", labelIndex)
				lgh := len(swCases)
				if nextJump && lgh > 0 {
					swCases[lgh-1].next.jump = casee
					swCases[lgh-1].next.index = idx
				}

				swCases = append(swCases, swCase{
					casee: casee,
					index: idx,
					body:  sc.Body,
				})

				if isFallingThroughStmt(sc.Body) {
					nextJump = true
				} else {
					nextJump = false
				}

				br := branch{name: casee, index: idx}
				swc.branches = append(swc.branches, br)
				for _, sv := range sc.Values {
					swc.label = fmt.Sprintf("switch_%d_check", labelIndex)
					swc.index = idx
					swc.value = sv
					lgl := len(swLabels)
					if lgl > 0 {
						br := branch{name: swc.label, index: idx}
						swLabels[lgl-1].branches = append(swLabels[lgl-1].branches, br)
					}
					swLabels = append(swLabels, swc)
					idx++
				}
			}
		}

		lgl := len(swLabels)
		if lgl > 0 {
			if hasDefault {
				swLabels[lgl-1].branches = append(swLabels[lgl-1].branches, branch{
					name: fmt.Sprintf("switch_%d_default", labelIndex),
					dft:  true,
				})
			} else {
				swLabels[lgl-1].branches = append(swLabels[lgl-1].branches, branch{
					name: fmt.Sprintf("switch_%d_end", labelIndex),
					dft:  true,
				})
			}
		}

		var (
			tagValue ir.Value
			result   string
		)
		if stmt.Init != nil {
			tagValue = l.lower(stmt.Init)
		}

		if stmt.Tag != nil {
			tagValue = l.lower(stmt.Tag)
			result = fmt.Sprintf("t%d", l.tIndex)
			l.tIndex++
		}

		for _, sc := range swLabels {
			if !sc.dft {
				lb := &ir.Label{
					Name:     sc.label,
					Index:    sc.index,
					NoSuffix: sc.dft,
				}
				l.instructions = append(l.instructions, lb)

				br := &ir.Branch{}
				if stmt.Tag == nil {
					br.Condition = string(l.lower(sc.value))
				} else {
					l.instructions = append(l.instructions, &ir.Binary{
						Result: result,
						Op:     token.BinaryOpString(token.Eq) + "_" + semantic.TBool.String(),
						Left:   string(tagValue),
						Right:  string(l.lower(sc.value)),
					})
					br.Condition = string(result)
				}

				if sc.isMulti {
					for _, v := range sc.branches {
						br.List = append(br.List, ir.BranchSub{
							Name:     v.name,
							Index:    v.index,
							NoSuffix: v.dft,
						})
					}
				} else {
					for _, v := range sc.branches {
						br.List = append(br.List, ir.BranchSub{
							Name:     v.name,
							Index:    v.index,
							NoSuffix: v.dft,
						})
					}
				}
				l.instructions = append(l.instructions, br)
			}
		}

		post := func(body []semantic.Stmt, next next) {
			var isTerminating bool
			for _, st := range body {
				isTerminating = isTerminatingStmt(st)
				if !isFallingThroughStmt([]semantic.Stmt{st}) {
					l.lowerStmt(st)
				}
			}

			if next.jump != "" {
				l.instructions = append(l.instructions, &ir.Jump{
					Name:     next.jump,
					Index:    next.index,
					NoSuffix: next.dft,
				})
			} else {
				if !isTerminating {
					l.instructions = append(l.instructions, &ir.Jump{
						Name:     fmt.Sprintf("switch_%d_end", labelIndex),
						NoSuffix: true,
					})
				}
			}
		}

		for _, sc := range swCases {
			lb := &ir.Label{
				Name:  sc.casee,
				Index: sc.index,
			}
			if sc.dft {
				lb.NoSuffix = true
			}
			l.instructions = append(l.instructions, lb)

			post(sc.body, sc.next)
		}

		l.instructions = append(l.instructions, &ir.Label{
			Name:     fmt.Sprintf("switch_%d_end", labelIndex),
			NoSuffix: true,
		})

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

	case *semantic.UnaryExpr:
		ex := l.lower(expr.Right)
		t := fmt.Sprintf("t%d", l.tIndex)
		l.instructions = append(l.instructions, &ir.Unary{
			Result:   t,
			Operator: token.UnaryOpString(expr.Operator),
			Value:    string(ex),
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

// isFallingThroughStmt returns true if statement is fallthrough
func isFallingThroughStmt(stmt []semantic.Stmt) bool {
	switch len(stmt) {
	case 0:
		return false
	case 1:
		if _, ok := stmt[0].(*semantic.FallThroughStmt); ok {
			return true
		}
		return false
	default:
		if _, ok := stmt[len(stmt)-1].(*semantic.FallThroughStmt); ok {
			return true
		}
		return false
	}
}
