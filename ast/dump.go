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
	d.node(0, input)
	return b.String()
}

func (d *dumper) node(indent int, n any) {
	switch v := n.(type) {
	case token.Token:
		d.line(indent, fmt.Sprintf("Token: %s", fmtTok(v)))

	case *File:
		d.line(indent, "File")
		d.kv(indent+1, "Package", v.PackageKW)
		d.kv(indent+1, "Name", v.Name)

		if len(v.Decls) > 0 {
			d.line(indent+1, "Decls")
			for _, decl := range v.Decls {
				d.decl(indent+2, decl)
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
				d.typ(indent+4, p.Type)
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
				d.typ(indent+5, p.Type)
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
		d.stmt(indent+2, v.Body)

	case *BlockStmt:
		if v.Stmts == nil {
			return
		}
		d.line(indent, "BlockStmt")
		d.kv(indent+1, "LBrace", v.LBrace)
		d.line(indent+1, "Stmts")
		for _, s := range v.Stmts {
			d.stmt(indent+2, s)
		}

		d.kv(indent+1, "RBrace", v.RBrace)

	case *ConstDecl:
		d.line(indent, "ConstDecl")
		d.kv(indent+1, "Const", v.ConstKW)
		d.kv(indent+1, "Name", v.Name)

		d.line(indent+1, "Type")
		d.typ(indent+2, v.Type)

		d.kv(indent+1, "Eq", v.Eq)
		d.line(indent+1, "Init")
		d.expr(indent+2, v.Init)

	case *VarDecl:
		d.line(indent, "VarDecl")
		d.kv(indent+1, "Var", v.VarKW)
		d.kv(indent+1, "Name", v.Name)
		if v.View != (token.Token{}) {
			d.kv(indent+1, "View", v.View)
		}

		d.line(indent+1, "Type")
		d.typ(indent+2, v.Type)

		d.kv(indent+1, "Eq", v.Eq)
		d.line(indent+1, "Init")
		d.expr(indent+2, v.Init)

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

	case *NamedType:
		d.line(indent, "NamedType")
		for k, p := range v.Parts {
			if k%2 == 0 {
				d.kv(indent+1, "Ident", p)
			} else {
				d.kv(indent+1, "Dot", p)
			}
		}

	case *ParenExpr:
		d.line(indent, "ParenExpr")
		if v.Inner == nil {
			return
		}
		d.expr(indent+1, v.Inner)

	case *BinaryExpr:
		d.line(indent, "BinaryExpr")
		d.expr(indent+1, v.Left)
		d.kv(indent+1, "Operator", v.Operator)
		d.expr(indent+1, v.Right)

	case *UnaryExpr:
		d.line(indent, "UnaryExpr")
		d.kv(indent+1, "Operator", v.Operator)
		d.expr(indent+1, v.Right)

	case *SelectorExpr:
		d.line(indent, "SelectorExpr")
		d.line(indent+1, "X:")
		d.expr(indent+2, v.X)
		d.kv(indent+1, "Dot", v.Dot)
		d.kv(indent+1, "Selector", v.Selector)

	case *IndexExpr:
		d.line(indent, "IndexExpr")
		d.line(indent+1, "X:")
		d.expr(indent+1, v.X)
		d.kv(indent+1, "LBracket", v.LBracket)
		d.expr(indent+2, v.Index)
		d.kv(indent+1, "RBracket", v.RBracket)

	case *CallExpr:
		d.line(indent, "CallExpr")
		d.line(indent+1, "Callee")
		d.expr(indent+2, v.Callee)
		d.kv(indent+1, "LParent", v.LParen)
		if len(v.Args) > 0 {
			d.line(indent+1, "Args:")
			for _, v := range v.Args {
				d.expr(indent+2, v)
			}
		}
		d.kv(indent+1, "RParent", v.RParen)

	case *AssignStmt:
		d.line(indent, "AssignStmt")
		d.line(indent+1, "Left")
		d.expr(indent+2, v.Left)
		d.kv(indent+1, "Operator", v.Operator)
		d.line(indent+1, "Right")
		d.expr(indent+2, v.Right)

	case *ExprStmt:
		d.expr(indent, v.Expr)

	case *ReturnStmt:
		d.line(indent, "ReturnStmt")
		if len(v.Values) > 0 {
			d.line(indent+1, "Values")
			for _, v := range v.Values {
				d.expr(indent+2, v)
			}
		}

	case *IfStmt:
		d.line(indent, "IfStmt")
		d.line(indent+1, "Condition")
		d.expr(indent+2, v.Condition)
		d.line(indent+1, "Then")
		d.stmt(indent+2, v.Then)

		if v.Else != nil {
			d.line(indent+1, "Else")
			d.stmt(indent+2, v.Else)
		}

	case *ForStmt:
		d.line(indent, "ForStmt")
		d.kv(indent+1, "For", v.ForKW)
		if v.Init != nil {
			d.line(indent+1, "Init")
			d.stmt(indent+2, v.Init)
		}
		if v.Condition != nil {
			d.line(indent+1, "Condition")
			d.expr(indent+2, v.Condition)
		}
		if v.Post != nil {
			d.line(indent+1, "Post")
			d.stmt(indent+2, v.Post)
		}
		d.stmt(indent+2, v.Body)

	case *RangeStmt:
		d.line(indent, "RangeStmt")
		d.kv(indent+1, "For", v.ForKW)
		if v.Key != nil {
			d.line(indent+1, "Key")
			d.expr(indent+2, v.Key)
		}
		if v.Value != nil {
			d.line(indent+1, "Condition")
			d.expr(indent+2, v.Value)
		}
		d.kv(indent+1, "Op", v.Op)
		d.kv(indent+1, "Range", v.Range)
		d.expr(indent+2, v.X)
		d.stmt(indent+2, v.Body)

	case *IncDecStmt:
		d.line(indent, "IncDecStmt")
		d.line(indent+1, "X:")
		d.expr(indent+2, v.X)
		d.kv(indent+1, "Operator", v.Operator)

	case *BreakStmt:
		d.line(indent, "BreakStmt")
		d.kv(indent+1, "Break", v.Break)

	case *ContinueStmt:
		d.line(indent, "ContinueStmt")
		d.kv(indent+1, "Continue", v.Continue)

	case *SwitchStmt:
		d.line(indent, "SwitchStmt")
		d.kv(indent+1, "Switch", v.Switch)
		if v.Init != nil {
			d.line(indent+1, "Init:")
			d.stmt(indent+2, v.Init)
		}

		if v.Tag != nil {
			d.line(indent+1, "Init:")
			d.expr(indent+2, v.Tag)
		}

		d.kv(indent+1, "LBrace", v.LBrace)
		for _, vc := range v.Cases {
			d.kv(indent+1, "Case", vc.Case)
			if len(vc.Values) > 0 {
				d.line(indent+2, "Values:")
				for _, e := range vc.Values {
					d.expr(indent+3, e)
				}
			}
			d.kv(indent+1, "Colon", vc.Colon)
			if len(vc.Body) > 0 {
				d.line(indent+2, "Body:")
				for _, b := range vc.Body {
					d.stmt(indent+3, b)
				}
			}
		}
		d.kv(indent+1, "RBrace", v.RBrace)

	case *FallThroughStmt:
		d.line(indent, "FallThroughStmt")
		d.kv(indent+1, "FallThrough", v.FallThrough)

	case *StructDecl:
		d.line(indent, "StructDecl:")
		d.kv(indent+1, "Type", v.TypeDecl)
		d.kv(indent+1, "Name", v.Name)
		d.kv(indent+1, "Struct", v.Struct)
		if v.Public {
			d.line(indent+1, "Public: true")
		}
		d.kv(indent+1, "LBrace", v.LBrace)
		if len(v.Fields) > 0 {
			for _, f := range v.Fields {
				d.kv(indent+2, "Name", f.Name)
				if f.Public {
					d.line(indent+2, "Public: true")
				}
				d.line(indent+2, "Type:")
				d.typ(indent+3, f.Type)
				if f.Eq != nil {
					d.kv(indent+2, "Eq", *f.Eq)
					d.expr(indent+2, f.Default)
				}
			}
		}
		d.kv(indent+1, "RBrace", v.RBrace)

	case *InterfaceDecl:
		d.line(indent, "InterfaceDecl:")
		d.kv(indent+1, "Type", v.TypeDecl)
		d.kv(indent+1, "Name", v.Name)
		d.kv(indent+1, "Interface", v.Interface)
		if v.Public {
			d.line(indent+1, "Public: true")
		}
		d.kv(indent+1, "LBrace", v.LBrace)
		if len(v.Embeds) > 0 {
			d.line(indent+2, "Embeds")
			for _, e := range v.Embeds {
				d.node(indent+3, e)
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
						d.typ(indent+4, p.Type)
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
						d.typ(indent+5, p.Type)
					}

					if f.Results.RParen != (token.Token{}) {
						d.kv(indent+2, "RParent", f.Results.RParen)
					}
				}
			}
		}
		d.kv(indent+1, "RBrace", v.RBrace)

	case *ImplementsDecl:
		d.line(indent, "ImplementsDecl:")
		d.kv(indent+1, "TypeName", v.TypeName)
		d.kv(indent+2, "Implements", v.Implements)
		d.line(indent+2, "Interface")
		d.node(indent+3, v.Interface)

	case *EnumDecl:
		d.line(indent, "EnumDecl:")
		d.kv(indent+1, "Type", v.TypeDecl)
		d.kv(indent+2, "Name", v.Name)
		if v.Public {
			d.line(indent+2, "Public: true")
		}
		d.kv(indent+2, "Enum", v.Enum)
		d.kv(indent+2, "LBrace", v.LBrace)
		d.line(indent+3, "Variants")
		for _, p := range v.Variants {
			d.kv(indent+4, "Ident", p)
		}
		d.kv(indent+2, "RBrace", v.RBrace)

	case *SumDecl:
		d.line(indent, "SumDecl:")
		d.kv(indent+1, "Type", v.TypeDecl)
		d.kv(indent+2, "Name", v.Name)
		d.kv(indent+2, "Sum", v.Sum)
		if v.Public {
			d.line(indent+1, "Public: true")
		}
		d.kv(indent+1, "LBrace", v.LBrace)
		if len(v.Variants) > 0 {
			d.line(indent+2, "Variants")
			for _, f := range v.Variants {
				d.kv(indent+3, "SumVariant", f.Name)
				if len(f.Params) > 0 {
					d.line(indent+4, "Params")
					for _, p := range f.Params {
						d.line(indent+5, "Param")
						d.kv(indent+6, "Ident", p.Name)
						d.line(indent+6, "Type")
						d.typ(indent+7, p.Type)
					}
				}
			}
		}
		d.kv(indent+1, "RBrace", v.RBrace)

	case *ComptimeBlockDecl:
		d.line(indent, "CompTimeBlockDecl:")
		d.kv(indent+1, "Comptime", v.ComptimeKW)
		for _, dec := range v.Decls {
			d.node(indent+2, dec)
		}

	case *SliceExpr:
		d.line(indent, "SliceExpr")
		d.expr(indent+1, v.X)
		d.kv(indent+1, "LBracket", v.LBracket)
		if v.Low != nil {
			d.expr(indent+1, v.Low)
		}
		if v.Colon != (token.Token{}) {
			d.kv(indent+2, "Colon", v.Colon)
		}
		if v.High != nil {
			d.expr(indent+1, v.High)
		}
		d.kv(indent+1, "RBracket", v.RBracket)

	case *SliceType:
		d.line(indent, "SliceType:")
		d.kv(indent+1, "LBracket", v.LBracket)
		d.kv(indent+1, "RBracket", v.RBracket)
		d.node(indent+1, v.Elem)

	case *ArrayType:
		d.line(indent, "ArrayType:")
		d.kv(indent+1, "LBracket", v.LBracket)
		d.expr(indent+1, v.Len)
		d.kv(indent+1, "RBracket", v.RBracket)
		d.node(indent+1, v.Elem)

	case *MapType:
		d.line(indent, "MapType:")
		if v.KindKW.Kind == token.KWMap {
			d.kv(indent+1, "Map", v.KindKW)
		} else {
			d.kv(indent+1, "Hashmap", v.KindKW)
		}
		d.kv(indent+1, "LBracket", v.LBracket)
		d.line(indent+1, "KeyType:")
		d.node(indent+2, v.KeyType)
		d.kv(indent+1, "RBracket", v.RBracket)
		d.line(indent+1, "ValueType:")
		d.node(indent+2, v.ValueType)

	case *MakeExpr:
		d.line(indent, "MakeExpr:")
		d.kv(indent+1, "Make", v.MakeKW)
		d.kv(indent+1, "LParen", v.LParen)
		d.node(indent+1, v.Type)
		if len(v.Args) > 0 {
			d.line(indent+1, "Size:")
			d.node(indent+2, v.Args[0])
			if len(v.Args) == 2 {
				d.line(indent+1, "Cap:")
				d.node(indent+2, v.Args[1])
			}
		}
		d.kv(indent+1, "RParen", v.RParen)

	case *SliceLitExpr:
		d.node(indent, v.Type)
		d.kv(indent+1, "LBrace", v.LBrace)
		if len(v.Elements) > 0 {
			d.line(indent+2, "Elements")
			for _, x := range v.Elements {
				d.node(indent+3, x)
			}
		}
		d.kv(indent+2, "RBrace", v.RBrace)

	case *DeclStmt:
		d.node(indent, v.Decl)

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
	if strings.Contains(t.Value, `"`) {
		return fmt.Sprintf("%s @%d:%d (kind=%d)", t.Value, t.Line, t.Column, t.Kind)
	}
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

func (d *dumper) decl(indent int, n Decl) {
	switch v := n.(type) {
	case *FuncDecl, *BadDecl, *ConstDecl, *VarDecl, *InterfaceDecl, *EnumDecl, *SumDecl:
		d.node(indent, v)

	case *ComptimeBlockDecl, *StructDecl, *ImplementsDecl:
		d.node(indent, v)

	default:
		if n == nil {
			d.line(indent, "(nil decl)")
			return
		}
		d.line(indent, fmt.Sprintf("<<unhandled decl %T>>", n))
	}
}

func (d *dumper) typ(indent int, n Type) {
	switch v := n.(type) {
	case *NamedType, *BadType, *MapType, *ArrayType, *SliceType:
		d.node(indent, v)

	default:
		if n == nil {
			d.line(indent, "(nil type)")
			return
		}
		d.line(indent, fmt.Sprintf("<<unhandled type %T>>", n))
	}
}

func (d *dumper) stmt(indent int, n Stmt) {
	switch v := n.(type) {
	case *BlockStmt, *AssignStmt, *ExprStmt, *BadStmt:
		d.node(indent, v)

	case *ReturnStmt, *IfStmt, *ForStmt, *RangeStmt, *IncDecStmt:
		d.node(indent, v)

	case *BreakStmt, *ContinueStmt, *SwitchStmt, *FallThroughStmt, *DeclStmt:
		d.node(indent, v)

	default:
		if n == nil {
			d.line(indent, "(nil stmt)")
			return
		}
		d.line(indent, fmt.Sprintf("<<unhandled stmt %T>>", n))
	}
}

func (d *dumper) expr(indent int, n Expr) {
	switch v := n.(type) {
	case *IdentExpr, *IntLitExpr, *FloatLitExpr, *BoolLitExpr, *StringLitExpr:
		d.node(indent, v)

	case *BadExpr, *BinaryExpr, *ParenExpr, *UnaryExpr, *SelectorExpr:
		d.node(indent, v)

	case *IndexExpr, *CallExpr, *SliceExpr, *MakeExpr, *SliceLitExpr:
		d.node(indent, v)

	default:
		if n == nil {
			d.line(indent, "(nil expr)")
			return
		}
		d.line(indent, fmt.Sprintf("<<unhandled expr %T>>", n))
	}
}
