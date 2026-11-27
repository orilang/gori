package lower

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/orilang/gori/ir"
)

// dump allows us to print LIR in a human readable format
func dump(input any) string {
	var b bytes.Buffer
	d := dumper{w: &b}
	d.node(0, input)
	return b.String()
}

func (d *dumper) node(indent int, n any) {
	switch t := n.(type) {
	case ir.Program:
		for _, f := range t.Types {
			d.node(indent, f)
		}

		for _, f := range t.Funcs {
			d.node(indent, f)
		}

	case *ir.Func:
		d.line(indent, d.funcDecl(true, t.Receiver, t.Name, t.Params, t.Results))
		if t.Blocks != nil {
			_, _ = fmt.Fprintf(d.w, "%s:\n", t.Label)
			for _, block := range t.Blocks {
				d.node(indent+2, block)
			}
		}
		d.w.WriteString("\n")

	case *ir.Block:
		for _, v := range t.Instructions {
			d.node(indent+2, v)
		}

	case *ir.Binary:
		d.line(indent, fmt.Sprintf("%s = %s %s, %s", t.Result, t.Op, t.Left, t.Right))

	case *ir.Return:
		if len(t.Values) == 0 {
			d.line(indent, "return")
		} else {
			d.line(indent, fmt.Sprintf("return %s", strings.Join(t.Values, ", ")))
		}

	case *ir.Const:
		d.line(indent, fmt.Sprintf("%s = const_%s %s", t.Result, t.Type, t.Value))

	case *ir.Call:
		if t.FromFunc {
			d.line(indent, fmt.Sprintf("%s = call %s(%s)", t.Result, t.Name, strings.Join(t.Args, ", ")))
		} else {
			d.line(indent, fmt.Sprintf("%s = %s(%s)", t.Result, t.Name, strings.Join(t.Args, ", ")))
		}

	case *ir.CallMethod:
		d.line(indent, fmt.Sprintf("%s = call_method %s, %s(%s)", t.Result, t.Name, t.Selector, strings.Join(t.Args, ", ")))

	case *ir.Branch:
		var branches []string
		if t.Condition != "" {
			branches = append(branches, t.Condition)
		}
		for _, v := range t.List {
			if v.NoSuffix {
				branches = append(branches, v.Name)
			} else {
				branches = append(branches, fmt.Sprintf("%s_%d", v.Name, v.Index))
			}
		}

		d.line(indent, fmt.Sprintf("branch %s", strings.Join(branches, ", ")))

	case *ir.Label:
		if t.NoSuffix {
			d.line(0, fmt.Sprintf("\n%s:", t.Name))
		} else {
			d.line(0, fmt.Sprintf("\n%s_%d:", t.Name, t.Index))
		}

	case *ir.Jump:
		if t.NoSuffix {
			d.line(indent, fmt.Sprintf("jump %s", t.Name))
		} else {
			d.line(indent, fmt.Sprintf("jump %s_%d", t.Name, t.Index))
		}

	case *ir.Assigment:
		d.line(indent, fmt.Sprintf("%s %s %s", t.Result, t.Operator, t.Value))

	case *ir.Extract:
		d.line(indent, fmt.Sprintf("%s = extract %s, %d", t.Result, t.Value, t.Index))

	case *ir.Unary:
		d.line(indent, fmt.Sprintf("%s = %s%s", t.Result, t.Operator, t.Value))

	case *ir.IncDec:
		d.line(indent, fmt.Sprintf("%s%s", t.Value, t.Operator))

	case *ir.Index:
		d.line(indent, fmt.Sprintf("%s = index %s, %s", t.Result, t.Name, t.Value))

	case *ir.Len:
		d.line(indent, fmt.Sprintf("%s = len %s", t.Result, t.Value))

	case *ir.MapIter:
		d.line(indent, fmt.Sprintf("%s = map_iter %s", t.Result, t.Value))

	case *ir.IterNext:
		d.line(indent, fmt.Sprintf("%s = iter_next %s", t.Result, t.Value))

	case *ir.IterKey:
		d.line(indent, fmt.Sprintf("%s = iter_key %s", t.Result, t.Value))

	case *ir.IterValue:
		d.line(indent, fmt.Sprintf("%s = iter_value %s", t.Result, t.Value))

	case *ir.DefinedTypeDecl:
		d.line(indent, fmt.Sprintf("type %s %s", t.Name, t.Type))

	case *ir.StructDecl:
		d.line(indent, fmt.Sprintf("type %s struct {", t.Name))
		for _, f := range t.Fields {
			if f.Default != "" {
				d.line(indent+2, fmt.Sprintf("%s %s = %s", f.Name, f.Type, f.Default))
			} else {
				d.line(indent+2, fmt.Sprintf("%s %s", f.Name, f.Type))
			}
		}
		d.line(indent, "}\n")

	case *ir.FieldGet:
		d.line(indent, fmt.Sprintf("%s = field_get %s, %s", t.Result, t.X, t.Selector))

	default:
		if n == nil {
			d.line(indent, "(nil intermediate representation)")
			return
		}
		d.line(indent, fmt.Sprintf("<<unhandled intermediate representation %T>>", n))
	}
}

// line writes the content with indentation
func (d *dumper) line(indent int, s string) {
	d.w.WriteString(strings.Repeat(" ", indent))
	d.w.WriteString(s)
	d.w.WriteString("\n")
}

func (d *dumper) funcDecl(printFn bool, receiver *ir.Receiver, name string, params, results []ir.Param) string {
	var l []string

	final := name
	if printFn {
		if receiver != nil {
			final = fmt.Sprintf("func (%s:%s) %s", receiver.Name, receiver.Type, name)
		} else {
			final = "func " + name
		}
	}

	for _, v := range params {
		l = append(l, v.Name+":"+v.Type)
	}
	final = final + "(" + strings.Join(l, ", ") + ")"

	switch len(results) {
	case 0:
		return final

	case 1:
		if results[0].Name == "" {
			return final + " -> " + results[0].Type
		}
		return final + " -> (" + results[0].Name + ":" + results[0].Type + ")"

	default:
		l = nil
		for _, v := range results {
			if v.Name == "" {
				l = append(l, v.Type)
			} else {
				l = append(l, v.Name+":"+v.Type)
			}
		}
		return final + " -> (" + strings.Join(l, ", ") + ")"
	}
}
