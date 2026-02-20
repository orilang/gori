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
	case token.Token:
		d.line(indent, fmt.Sprintf("Token %s", fmtTok(v)))

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

		if len(v.Decls) > 0 {
			d.line(indent+1, "Decls")
			for _, decl := range v.Decls {
				d.decl(decl, indent+2)
			}
		}

		if len(v.Structs) > 0 {
			d.line(indent+1, "Structs")
			for _, v := range v.Structs {
				d.stmt(v, indent+2)
			}
		}

		if len(v.Interfaces) > 0 {
			d.line(indent+1, "Interfaces")
			for _, v := range v.Interfaces {
				d.stmt(v, indent+2)
			}
		}

		if len(v.Implements) > 0 {
			d.line(indent+1, "Implements")
			for _, v := range v.Implements {
				d.node(v, indent+2)
			}
		}

		if len(v.Enums) > 0 {
			d.line(indent+1, "Enums")
			for _, v := range v.Enums {
				d.node(v, indent+2)
			}
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
				d.kv(indent+3, "Ident", p.Name)
				d.line(indent+3, "Type")
				d.typ(p.Type, indent+4)
			}
		}

		if len(v.Results.List) > 0 {
			d.line(indent+1, "Results")
			if v.Results.LParen != (token.Token{}) {
				d.kv(indent+2, "LParent", v.Results.LParen)
			}

			for _, p := range v.Results.List {
				d.line(indent+3, "Param")
				if p.Name != (token.Token{}) {
					d.kv(indent+4, "Ident", p.Name)
				}
				d.line(indent+4, "Type")
				d.typ(p.Type, indent+5)
			}

			if v.Results.RParen != (token.Token{}) {
				d.kv(indent+2, "RParent", v.Results.RParen)
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

	case *BadExpr:
		d.line(indent+2, fmtBadExpr(v))

	case *BadStmt:
		d.line(indent+2, fmtBadStmt(v))

	case *BadDecl:
		d.line(indent+2, fmtBadDecl(v))

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

	case *ParenExpr:
		d.line(indent, "ParenExpr")
		if v.Inner == nil {
			return
		}
		d.expr(v.Inner, indent+1)

	case *BinaryExpr:
		d.line(indent, "BinaryExpr")
		d.expr(v.Left, indent+1)
		d.kv(indent+1, "Operator", v.Operator)
		d.expr(v.Right, indent+1)

	case *UnaryExpr:
		d.line(indent, "UnaryExpr")
		d.kv(indent+1, "Operator", v.Operator)
		d.expr(v.Right, indent+1)

	case *SelectorExpr:
		d.line(indent, "SelectorExpr")
		d.line(indent+1, "X:")
		d.expr(v.X, indent+2)
		d.kv(indent+1, "Dot", v.Dot)
		d.kv(indent+1, "Selector", v.Selector)

	case *IndexExpr:
		d.line(indent, "IndexExpr")
		d.line(indent+1, "X:")
		d.expr(v.X, indent+1)
		d.kv(indent+1, "LBracket", v.LBracket)
		d.expr(v.Index, indent+2)
		d.kv(indent+1, "RBracket", v.RBracket)

	case *CallExpr:
		d.line(indent, "CallExpr")
		d.line(indent+1, "Callee")
		d.expr(v.Callee, indent+2)
		d.kv(indent+1, "LParent", v.LParen)
		if len(v.Args) > 0 {
			d.line(indent+1, "Args:")
			for _, v := range v.Args {
				d.expr(v, indent+2)
			}
		}
		d.kv(indent+1, "RParent", v.RParen)

	case *AssignStmt:
		d.line(indent, "AssignStmt")
		d.line(indent+1, "Left")
		d.expr(v.Left, indent+2)
		d.kv(indent+1, "Operator", v.Operator)
		d.line(indent+1, "Right")
		d.expr(v.Right, indent+2)

	case *ExprStmt:
		d.expr(v.Expr, indent)

	case *ReturnStmt:
		d.line(indent, "ReturnStmt")
		if len(v.Values) > 0 {
			d.line(indent+1, "Values")
			for _, v := range v.Values {
				d.expr(v, indent+2)
			}
		}

	case *IfStmt:
		d.line(indent, "IfStmt")
		d.line(indent+1, "Condition")
		d.expr(v.Condition, indent+2)
		d.line(indent+1, "Then")
		d.stmt(v.Then, indent+2)

		if v.Else != nil {
			d.line(indent+1, "Else")
			d.stmt(v.Else, indent+2)
		}

	case *ForStmt:
		d.line(indent, "ForStmt")
		d.kv(indent+1, "For", v.For)
		if v.Init != nil {
			d.line(indent+1, "Init")
			d.stmt(v.Init, indent+2)
		}
		if v.Condition != nil {
			d.line(indent+1, "Condition")
			d.expr(v.Condition, indent+2)
		}
		if v.Post != nil {
			d.line(indent+1, "Post")
			d.stmt(v.Post, indent+2)
		}
		d.stmt(v.Body, indent+2)

	case *RangeStmt:
		d.line(indent, "RangeStmt")
		d.kv(indent+1, "For", v.For)
		if v.Key != nil {
			d.line(indent+1, "Key")
			d.expr(v.Key, indent+2)
		}
		if v.Value != nil {
			d.line(indent+1, "Condition")
			d.expr(v.Value, indent+2)
		}
		d.kv(indent+1, "Op", v.Op)
		d.kv(indent+1, "Range", v.Range)
		d.expr(v.X, indent+2)
		d.stmt(v.Body, indent+2)

	case *IncDecStmt:
		d.line(indent, "IncDecStmt")
		d.line(indent+1, "X:")
		d.expr(v.X, indent+2)
		d.kv(indent+1, "Operator", v.Operator)

	case *BreakStmt:
		d.kv(indent, "Break", v.Break)

	case *ContinueStmt:
		d.kv(indent, "Continue", v.Continue)

	case *SwitchStmt:
		d.line(indent, "SwitchStmt")
		d.kv(indent+1, "Switch", v.Switch)
		if v.Init != nil {
			d.line(indent+1, "Init:")
			d.stmt(v.Init, indent+2)
		}

		if v.Tag != nil {
			d.line(indent+1, "Init:")
			d.expr(v.Tag, indent+2)
		}

		d.kv(indent+1, "LBrace", v.LBrace)
		for _, vc := range v.Cases {
			d.kv(indent+1, "Case", vc.Case)
			if len(vc.Values) > 0 {
				d.line(indent+2, "Values:")
				for _, e := range vc.Values {
					d.expr(e, indent+3)
				}
			}
			d.kv(indent+1, "Colon", vc.Colon)
			if len(vc.Body) > 0 {
				d.line(indent+2, "Body:")
				for _, b := range vc.Body {
					d.stmt(b, indent+3)
				}
			}
		}
		d.kv(indent+1, "RBrace", v.RBrace)

	case *FallThroughStmt:
		d.kv(indent, "FallThrough", v.FallThroughStmt)

	case *StructType:
		d.kv(indent, "Type", v.TypeDecl)
		d.kv(indent, "Name", v.Name)
		d.kv(indent, "Struct", v.Struct)
		if v.Public {
			d.line(indent, "Public: true")
		}
		d.kv(indent, "LBrace", v.LBrace)
		if len(v.Fields) > 0 {
			for _, f := range v.Fields {
				d.kv(indent+1, "Name", f.Name)
				if f.Public {
					d.line(indent+1, "Public: true")
				}
				d.line(indent+1, "Type:")
				d.typ(f.Type, indent+2)
				if f.Eq != nil {
					d.kv(indent+1, "Eq", *f.Eq)
					d.expr(f.Default, indent+1)
				}
			}
		}
		d.kv(indent, "RBrace", v.RBrace)

	case *InterfaceType:
		d.kv(indent, "Type", v.TypeDecl)
		d.kv(indent, "Name", v.Name)
		d.kv(indent, "Interface", v.Interface)
		if v.Public {
			d.line(indent, "Public: true")
		}
		d.kv(indent, "LBrace", v.LBrace)
		if len(v.Embeds) > 0 {
			d.line(indent+2, "Embeds")
			for _, e := range v.Embeds {
				for _, p := range e.Parts {
					if p.Kind == token.Ident {
						d.kv(indent+3, "Ident", p)
					} else {
						d.kv(indent+3, "Dot", p)
					}
				}
			}
		}
		if len(v.Methods) > 0 {
			for _, f := range v.Methods {
				d.kv(indent+1, "Name", f.Name)
				d.line(indent+1, "Params")
				if len(f.Params) == 0 {
					d.line(indent+2, "(none)")
				} else {
					for _, p := range f.Params {
						d.line(indent+2, "Param")
						d.kv(indent+3, "Ident", p.Name)
						d.line(indent+3, "Type")
						d.typ(p.Type, indent+4)
					}
				}

				if len(f.Results.List) > 0 {
					d.line(indent+1, "Results")
					if f.Results.LParen != (token.Token{}) {
						d.kv(indent+2, "LParent", f.Results.LParen)
					}

					for _, p := range f.Results.List {
						d.line(indent+3, "Param")
						if p.Name != (token.Token{}) {
							d.kv(indent+4, "Ident", p.Name)
						}
						d.line(indent+4, "Type")
						d.typ(p.Type, indent+5)
					}

					if f.Results.RParen != (token.Token{}) {
						d.kv(indent+2, "RParent", f.Results.RParen)
					}
				}
			}
		}
		d.kv(indent, "RBrace", v.RBrace)

	case *ImplementsDecl:
		d.kv(indent, "Type", v.Type)
		d.kv(indent+1, "Implements", v.Implements)
		d.line(indent+1, "Interface")
		for _, p := range v.Interface.Parts {
			if p.Kind == token.Ident {
				d.kv(indent+2, "Ident", p)
			} else {
				d.kv(indent+2, "Dot", p)
			}
		}

	case *EnumType:
		d.kv(indent, "Type", v.TypeDecl)
		d.kv(indent+1, "Name", v.Name)
		if v.Public {
			d.line(indent+1, "Public: true")
		}
		d.kv(indent+1, "Enum", v.Enum)
		d.kv(indent+1, "LBrace", v.LBrace)
		d.line(indent+1, "Variants")
		for _, p := range v.Variants {
			d.kv(indent+2, "Ident", p)
		}
		d.kv(indent+1, "RBrace", v.RBrace)

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
		return fmt.Sprintf("BadType from @%d:%d to @%d:%d reason=%s from=%q to=%q", b.From.Line, b.From.Column, b.To.Line, b.To.Column, b.Reason, b.From.Value, b.To.Value)
	}
	return fmt.Sprintf("BadType at @%d:%d reason=%s value=%q", b.From.Line, b.From.Column, b.Reason, b.From.Value)
}

// fmtBadExpr returns bad expr content with line, column etc
func fmtBadExpr(b *BadExpr) string {
	if b.From != b.To && b.To != (token.Token{}) {
		return fmt.Sprintf("BadExpr at @%d:%d to @%d:%d reason=%s from=%q to=%q", b.From.Line, b.From.Column, b.To.Line, b.To.Column, b.Reason, b.From.Value, b.To.Value)
	}
	return fmt.Sprintf("BadExpr at @%d:%d reason=%s value=%q", b.From.Line, b.From.Column, b.Reason, b.From.Value)
}

// fmtBadStmt returns bad stmt content with line, column etc
func fmtBadStmt(b *BadStmt) string {
	if b.From != b.To && b.To != (token.Token{}) {
		return fmt.Sprintf("BadStmt at @%d:%d to @%d:%d reason=%s from=%q to=%q", b.From.Line, b.From.Column, b.To.Line, b.To.Column, b.Reason, b.From.Value, b.To.Value)
	}
	return fmt.Sprintf("BadStmt at @%d:%d reason=%s value=%q", b.From.Line, b.From.Column, b.Reason, b.From.Value)
}

// fmtBadDecl returns bad decl content with line, column etc
func fmtBadDecl(b *BadDecl) string {
	if b.From != b.To && b.To != (token.Token{}) {
		return fmt.Sprintf("BadDecl at @%d:%d to @%d:%d reason=%s from=%q to=%q", b.From.Line, b.From.Column, b.To.Line, b.To.Column, b.Reason, b.From.Value, b.To.Value)
	}
	return fmt.Sprintf("BadDecl at @%d:%d reason=%s value=%q", b.From.Line, b.From.Column, b.Reason, b.From.Value)
}

func (d *dumper) decl(n Decl, indent int) {
	switch v := n.(type) {
	case *FuncDecl, *BadDecl:
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
	case *BlockStmt, *ConstDeclStmt, *VarDeclStmt, *AssignStmt, *ExprStmt, *BadStmt:
		d.node(v, indent)

	case *ReturnStmt, *IfStmt, *ForStmt, *RangeStmt, *IncDecStmt:
		d.node(v, indent)

	case *BreakStmt, *ContinueStmt, *SwitchStmt, *FallThroughStmt:
		d.node(v, indent)

	case *StructType, *InterfaceType, *EnumType:
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

	case *BadExpr, *BinaryExpr, *ParenExpr, *UnaryExpr, *SelectorExpr:
		d.node(v, indent)

	case *IndexExpr, *CallExpr:
		d.node(v, indent)

	default:
		if n == nil {
			d.line(indent, "(nil expr)")
			return
		}
		d.line(indent, fmt.Sprintf("<<unhandled expr %T>>", n))
	}
}
