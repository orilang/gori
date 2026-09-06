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
		if len(t.Funcs) > 0 {
			for _, f := range t.Funcs {
				d.node(indent, f)
			}
		}

	case *ir.Func:
		d.line(indent, "func "+t.Name+d.funcDecl("arg", t.Params)+d.funcDecl("", t.Results))
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
		if t.Name == "" {
			d.line(indent, "return")
		} else {
			d.line(indent, fmt.Sprintf("return %s", t.Name))
		}

	case *ir.Const:
		d.line(indent, fmt.Sprintf("%s = const_%s %s", t.Result, t.Type, t.Value))

	case *ir.Call:
		d.line(indent, fmt.Sprintf("%s = %s(%s)", t.Result, t.Name, strings.Join(t.Args, ", ")))

	case *ir.Branch:
		d.line(indent, fmt.Sprintf("branch %s, %s_%d, %s_%d", t.Condition, t.True, t.Index, t.False, t.Index))

	case *ir.Label:
		d.line(0, fmt.Sprintf("\n%s_%d:", t.Name, t.Index))

	case *ir.Jump:
		d.line(indent, fmt.Sprintf("jump %s_%d", t.Name, t.Index))

	case *ir.Assigment:
		d.line(indent, fmt.Sprintf("%s = %s", t.Result, t.Value))

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

func (d *dumper) funcDecl(kind string, param []ir.Param) string {
	var l []string
	if kind == "arg" {
		for _, v := range param {
			l = append(l, v.Name+":"+v.Type)
		}
		return "(" + strings.Join(l, ", ") + ")"
	}

	switch len(param) {
	case 0:
		return ""

	case 1:
		if param[0].Name == "" {
			return " -> " + param[0].Type
		}
		return " -> (" + param[0].Name + ":" + param[0].Type + ")"

	default:
		for _, v := range param {
			if v.Name == "" {
				l = append(l, v.Type)
			} else {
				l = append(l, v.Name+":"+v.Type)
			}
		}
		return "(" + strings.Join(l, ", ") + ")"
	}
}
