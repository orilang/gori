package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/orilang/gori/token"
)

// Dump allows us to print AST in a human readable format
// for debugging
func Dump(input any) string {
	var b bytes.Buffer
	d := dumper{w: &b}
	d.node(input, 0)
	return b.String()
}

func (d *dumper) node(n any, indent int) {
	switch v := n.(type) {
	case *File:
		d.line(indent, "File")
		d.kv(indent+1, "Package", v.PackageKW)
		d.kv(indent+1, "Name", v.Name)

		if len(v.Const) > 0 {
			d.line(indent+1, "ConstDecls")
			for _, s := range v.Const {
				d.stmt(s, indent+2)
			}
		}

		d.line(indent+1, "Decls")
		if len(v.Decls) == 0 {
			d.line(indent+2, "(none)")
			return
		}
		for _, decl := range v.Decls {
			d.decl(decl, indent+2)
		}

	case *FuncDecl:
		d.line(indent, "FuncDecl")
		d.kv(indent+1, "Function", v.FuncKW)
		d.kv(indent+1, "Name", v.Name)

		d.line(indent+1, "Params")
		if len(v.Params) == 0 {
			d.line(indent+2, "(none)")
		} else {
			for _, p := range v.Params {
				d.line(indent+2, "Param")
				d.kv(indent+3, "Function", p.Name)
				d.line(indent+3, "Type")
				d.typ(p.Type, indent+4)
			}
		}

		d.line(indent+1, "Body")
		if v.Body == nil {
			d.line(indent+2, "(none)")
			return
		}
		d.stmt(v.Body, indent+2)

	case *BlockStmt:
		if v.Stmts == nil {
			return
		}
		d.line(indent, "BlockStmt")
		d.kv(indent+1, "LBrace", v.LBrace)
		d.line(indent+1, "Stmts")
		for _, s := range v.Stmts {
			d.stmt(s, indent+2)
		}

		d.kv(indent+1, "RBrace", v.RBrace)

	case *ConstDeclStmt:
		d.line(indent, "ConstDecl")
		d.kv(indent+1, "Const", v.ConstKW)
		d.kv(indent+1, "Name", v.Name)

		d.line(indent+1, "Type")
		if v.Type == nil {
			d.line(indent+2, "(none)")
		} else {
			d.typ(v.Type, indent+2)
		}

		d.kv(indent+1, "Eq", v.Eq)
		d.line(indent+1, "Init")
		if v.Init == nil {
			d.line(indent+2, "(none)")
		} else {
			d.expr(v.Init, indent+2)
		}

	case *VarDeclStmt:
		d.line(indent, "VarDeclStmt")
		d.kv(indent+1, "Var", v.VarKW)
		d.kv(indent+1, "Name", v.Name)

		d.line(indent+1, "Type")
		if v.Type == nil {
			d.line(indent+2, "(none)")
		} else {
			d.typ(v.Type, indent+2)
		}

		d.kv(indent+1, "Eq", v.Eq)
		d.line(indent+1, "Init")
		if v.Init == nil {
			d.line(indent+2, "(none)")
		} else {
			d.expr(v.Init, indent+2)
		}

	case *BadType:
		d.line(indent+2, fmtBadType(v))

	case *IdentExpr:
		d.line(indent, "IdentExpr")
		d.kv(indent+1, "Name", v.Name)

	case *IntLitExpr:
		d.line(indent, "IntLitExpr")
		d.kv(indent+1, "Value", v.Name)

	case *FloatLitExpr:
		d.line(indent, "FloatLitExpr")
		d.kv(indent+1, "Value", v.Name)

	case *BoolLitExpr:
		d.line(indent, "BoolLitExpr")
		d.kv(indent+1, "Value", v.Name)

	case *StringLitExpr:
		d.line(indent, "StringLitExpr")
		d.kv(indent+1, "Value", v.Name)

	case *NameType:
		d.line(indent, "NameType")
		d.kv(indent+1, "Name", v.Name)

	case token.Token:
		d.line(indent, fmt.Sprintf("Token %s", fmtTok(v)))

	default:
		if n == nil {
			d.line(indent, "(nil)")
			return
		}
		d.line(indent, fmt.Sprintf("<<unhandled %T>>", n))
	}
}

// line writes the content with indentation
func (d *dumper) line(indent int, s string) {
	d.w.WriteString(strings.Repeat(" ", indent))
	d.w.WriteString(s)
	d.w.WriteString("\n")
}

func (d *dumper) kv(indent int, key string, t token.Token) {
	d.line(indent, fmt.Sprintf("%s: %s", key, fmtTok(t)))
}

// fmtTok returns token content with line, column etc
func fmtTok(t token.Token) string {
	return fmt.Sprintf("%q @%d:%d (kind=%d)", t.Value, t.Line, t.Column, t.Kind)
}

// fmtBadType returns bad type content with line, column etc
func fmtBadType(b *BadType) string {
	if b.From != b.To && b.To != (token.Token{}) {
		return fmt.Sprintf("BadType from @%d:%d to @%d:%d reason=%s value=%s", b.From.Line, b.From.Column, b.To.Line, b.To.Column, b.Reason, b.From.Value)
	}
	return fmt.Sprintf("BadType at @%d:%d reason=%s value=%s", b.From.Line, b.From.Column, b.Reason, b.From.Value)
}

func (d *dumper) decl(n Decl, indent int) {
	switch v := n.(type) {
	case *FuncDecl:
		d.node(v, indent)

	default:
		if n == nil {
			d.line(indent, "(nil decl)")
			return
		}
		d.line(indent, fmt.Sprintf("<<unhandled decl %T>>", n))
	}
}

func (d *dumper) typ(n Type, indent int) {
	switch v := n.(type) {
	case *NameType:
		d.node(v, indent)

	case *BadType:
		d.node(v, indent)

	default:
		if n == nil {
			d.line(indent, "(nil type)")
			return
		}
		d.line(indent, fmt.Sprintf("<<unhandled type %T>>", n))
	}
}

func (d *dumper) stmt(n Stmt, indent int) {
	switch v := n.(type) {
	case *BlockStmt, *ConstDeclStmt, *VarDeclStmt:
		d.node(v, indent)

	default:
		if n == nil {
			d.line(indent, "(nil stmt)")
			return
		}
		d.line(indent, fmt.Sprintf("<<unhandled stmt %T>>", n))
	}
}

func (d *dumper) expr(n Expr, indent int) {
	switch v := n.(type) {
	case *IdentExpr, *IntLitExpr, *FloatLitExpr, *BoolLitExpr, *StringLitExpr:
		d.node(v, indent)

	default:
		if n == nil {
			d.line(indent, "(nil expr)")
			return
		}
		d.line(indent, fmt.Sprintf("<<unhandled expr %T>>", n))
	}
}
