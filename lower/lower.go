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

	case *semantic.VarDecl:
		init := l.lower(t.Init)
		l.instructions = append(l.instructions, &ir.Const{
			Result: t.Symbol.Name,
			Type:   t.Symbol.Type.String(),
			Value:  singleValue(init.values),
		})

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
func (l *Lower) lower(input any) (data infos) {
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

	return
}

// lowerStmt lowers expr to later create an instruction
func (l *Lower) lowerStmt(t semantic.Stmt) (data infos) {
	switch stmt := t.(type) {
	case *semantic.FallThroughStmt:
	case *semantic.ReturnStmt:
		data.returns = true
		if len(stmt.Values) == 0 {
			l.instructions = append(l.instructions, &ir.Return{})
			data.values = append(data.values, ir.Value(""))
			return
		}

		var rt []ir.Value
		irr := &ir.Return{}
		for _, v := range stmt.Values {
			st := l.lower(v)
			rt = append(rt, st.values...)
			for _, rv := range st.values {
				irr.Values = append(irr.Values, string(rv))
			}
		}

		l.instructions = append(l.instructions, irr)
		data.values = rt
		return

	case *semantic.AssigmentStmt:
		var results []infos
		for _, sr := range stmt.Right {
			results = append(results, l.lower(sr))
		}

		for k := range stmt.Right {
			if len(stmt.Right) > 1 {
				if !stmt.Symbol[k].IsBlank {
					l.instructions = append(l.instructions, &ir.Assigment{Result: stmt.Symbol[k].Name, Value: singleValue(results[k].values)})
					data.values = append(data.values, ir.Value(stmt.Symbol[k].Name))
				}
			} else {
				for sk, ss := range stmt.Symbol {
					if !ss.IsBlank {
						if ss.FromFunc && ss.IsMultiValue {
							l.instructions = append(l.instructions, &ir.Extract{Result: ss.Name, Value: singleValue(results[k].values), Index: sk})
						} else {
							l.instructions = append(l.instructions, &ir.Assigment{Result: ss.Name, Value: singleValue(results[k].values)})
						}
						data.values = append(data.values, ir.Value(ss.Name))
					}
				}
			}
		}
		return

	case *semantic.IfStmt:
		labelIndex := l.labelIndex
		l.labelIndex++

		cond := l.lower(stmt.Condition)
		br := &ir.Branch{Condition: singleValue(cond.values)}
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

		var thenr infos
		for _, v := range stmt.Then {
			thenr = l.lower(v)
		}

		if !thenr.returns {
			l.instructions = append(l.instructions, jump)
		}

		var elser infos
		if len(stmt.Else) > 0 {
			l.instructions = append(l.instructions, &ir.Label{Name: "if_else", Index: labelIndex})
			for _, v := range stmt.Else {
				elser = l.lower(v)
			}

			if !elser.returns {
				l.instructions = append(l.instructions, jump)
			}
		}

		x := thenr.returns && elser.returns
		if !(x) {
			l.instructions = append(l.instructions, end)
		}
		data.returns = thenr.returns && elser.returns

	case *semantic.SwitchStmt:
		labelIndex := l.labelIndex
		l.labelIndex++

		type branch struct {
			name  string
			index int
			dft   bool
		}

		type swLabel struct {
			label    string
			branches []branch
			dft      bool
			index    int
			isMulti  bool
			value    semantic.Expr
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
			swLabels         []swLabel
			swCases          []swCase
			isFallingThrough bool
			hasDefault       bool
		)

		idx := 1
		fillPrevBranch := -1
		for k, sc := range stmt.Cases {
			swc := swLabel{}

			switch len(sc.Values) {
			case 0:
				hasDefault = true
				swc.dft = true
				swc.label = fmt.Sprintf("switch_%d_default", labelIndex)
				br := branch{name: swc.label, dft: true}
				swc.branches = append(swc.branches, br)

				lgl := len(swLabels)
				if lgl > 0 && len(stmt.Cases)-1 == k {
					// if last, fill previous label
					swLabels[lgl-1].branches = append(swLabels[lgl-1].branches, br)
				}
				if lgl > 0 {
					// default ONLY participates at case level when the previous one
					// falltrough in here.
					// It DOES NOT participates at label level and
					// its placements remain irrelevant regarding the control flow graph (CFG)
					fillPrevBranch = lgl - 1
				}
				swLabels = append(swLabels, swc)

				lgh := len(swCases)
				if isFallingThrough && lgh > 0 {
					swCases[lgh-1].next.jump = swc.label
					swCases[lgh-1].next.index = idx
					swCases[lgh-1].next.dft = true
				}

				if isFallingThroughStmt(sc.Body) {
					isFallingThrough = true
				} else {
					isFallingThrough = false
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

				if fillPrevBranch != -1 {
					br := branch{name: swc.label, index: idx}
					swLabels[fillPrevBranch].branches = append(swLabels[fillPrevBranch].branches, br)
					fillPrevBranch = -1
				}
				br := branch{name: casee, index: idx}
				swc.branches = append(swc.branches, br)
				lgl := len(swLabels)
				if lgl > 0 {
					br.name = swc.label
					swLabels[lgl-1].branches = append(swLabels[lgl-1].branches, br)
				}
				swLabels = append(swLabels, swc)

				lgh := len(swCases)
				if isFallingThrough && lgh > 0 {
					swCases[lgh-1].next.jump = casee
					swCases[lgh-1].next.index = idx
				}

				if isFallingThroughStmt(sc.Body) {
					isFallingThrough = true
				} else {
					isFallingThrough = false
				}

				swCases = append(swCases, swCase{
					casee: casee,
					index: idx,
					body:  sc.Body,
				})

				idx++

			default:
				swc.isMulti = true

				casee := fmt.Sprintf("switch_%d_case", labelIndex)
				lgh := len(swCases)
				if isFallingThrough && lgh > 0 {
					swCases[lgh-1].next.jump = casee
					swCases[lgh-1].next.index = idx
				}

				if isFallingThroughStmt(sc.Body) {
					isFallingThrough = true
				} else {
					isFallingThrough = false
				}

				swCases = append(swCases, swCase{
					casee: casee,
					index: idx,
					body:  sc.Body,
				})

				br := branch{name: casee, index: idx}
				swc.branches = append(swc.branches, br)
				for _, sv := range sc.Values {
					swc.label = fmt.Sprintf("switch_%d_check", labelIndex)
					swc.index = idx
					swc.value = sv

					if fillPrevBranch != -1 {
						br := branch{name: swc.label, index: idx}
						swLabels[fillPrevBranch].branches = append(swLabels[fillPrevBranch].branches, br)
						fillPrevBranch = -1
					}

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

		var tagValue, result string
		if stmt.Init != nil {
			tagValue = singleValue(l.lower(stmt.Init).values)
		}

		if stmt.Tag != nil {
			tagValue = singleValue(l.lower(stmt.Tag).values)
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
					br.Condition = singleValue(l.lower(sc.value).values)
				} else {
					l.instructions = append(l.instructions, &ir.Binary{
						Result: result,
						Op:     token.BinaryOpString(token.Eq) + "_" + semantic.TBool.String(),
						Left:   tagValue,
						Right:  singleValue(l.lower(sc.value).values),
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

		post := func(body []semantic.Stmt, next next) bool {
			var returns bool
			for _, st := range body {
				if !isFallingThroughStmt([]semantic.Stmt{st}) {
					result := l.lowerStmt(st)
					returns = result.returns
				}
			}

			if next.jump != "" {
				l.instructions = append(l.instructions, &ir.Jump{
					Name:     next.jump,
					Index:    next.index,
					NoSuffix: next.dft,
				})
			} else {
				if !returns {
					l.instructions = append(l.instructions, &ir.Jump{
						Name:     fmt.Sprintf("switch_%d_end", labelIndex),
						NoSuffix: true,
					})
				}
			}
			return returns
		}

		var returns bool
		for k, sc := range swCases {
			lb := &ir.Label{
				Name:  sc.casee,
				Index: sc.index,
			}
			if sc.dft {
				lb.NoSuffix = true
			}
			l.instructions = append(l.instructions, lb)

			result := post(sc.body, sc.next)
			if k == 0 {
				returns = result
			} else {
				returns = returns && result
			}
		}

		data.returns = hasDefault && returns
		if !data.returns {
			l.instructions = append(l.instructions, &ir.Label{
				Name:     fmt.Sprintf("switch_%d_end", labelIndex),
				NoSuffix: true,
			})
		}

	case *semantic.DeclStmt:
		l.decl(stmt.Decl)

	default:
		l.errors = append(l.errors, Diagnostic{Err: fmt.Errorf("unsupported statement %T", stmt)})
	}

	return
}

// lowerExpr lowers expr to later create an instruction
func (l *Lower) lowerExpr(t semantic.Expr) (data infos) {
	switch expr := t.(type) {
	case *semantic.IntLitExpr:
		data.values = append(data.values, ir.Value(expr.Value))
		return

	case *semantic.FloatLitExpr:
		data.values = append(data.values, ir.Value(expr.Value))
		return

	case *semantic.BoolLitExpr:
		data.values = append(data.values, ir.Value(expr.Value))
		return

	case *semantic.StringLitExpr:
		data.values = append(data.values, ir.Value(expr.Value))
		return

	case *semantic.IdentExpr:
		data.values = append(data.values, ir.Value(expr.Value))
		return

	case *semantic.BinaryExpr:
		left := l.lower(expr.Left)
		right := l.lower(expr.Right)
		t := fmt.Sprintf("t%d", l.tIndex)

		l.instructions = append(l.instructions, &ir.Binary{
			Result: t,
			Op:     token.BinaryOpString(expr.Operator) + "_" + expr.Type.String(),
			Left:   singleValue(left.values),
			Right:  singleValue(right.values),
		})

		l.tIndex++
		data.values = append(data.values, ir.Value(t))
		return

	case *semantic.ConversionExpr:
		value := l.lower(expr.Value)
		t := fmt.Sprintf("t%d", l.tIndex)
		l.instructions = append(l.instructions, &ir.Const{
			Result: t,
			Type:   expr.To.String(),
			Value:  singleValue(value.values),
		})

		l.tIndex++
		data.values = append(data.values, ir.Value(t))
		return

	case *semantic.CallExpr:
		var args []string
		for _, v := range expr.Args {
			args = append(args, singleValue(l.lower(v).values))
		}

		callee := l.lower(expr.Callee)
		t := fmt.Sprintf("t%d", l.tIndex)
		l.instructions = append(l.instructions, &ir.Call{
			Result:   t,
			Name:     singleValue(callee.values),
			Args:     args,
			FromFunc: expr.FromFunc,
		})

		l.tIndex++
		data.values = append(data.values, ir.Value(t))
		return

	case *semantic.UnaryExpr:
		ex := l.lower(expr.Right)
		t := fmt.Sprintf("t%d", l.tIndex)
		l.instructions = append(l.instructions, &ir.Unary{
			Result:   t,
			Operator: token.UnaryOpString(expr.Operator),
			Value:    singleValue(ex.values),
		})
		l.tIndex++
		data.values = append(data.values, ir.Value(t))
		return

	default:
		l.errors = append(l.errors, Diagnostic{Err: fmt.Errorf("unsupported expression %T", expr)})
	}

	return
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

// singleValue returns "" if the slice is empty
// otherwise the related string
func singleValue(s []ir.Value) string {
	if len(s) == 0 {
		return ""
	}
	return string(s[0])
}
