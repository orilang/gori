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

		if len(v.Const) > 0 {
			d.line(indent+1, "ConstDecls")
			for _, s := range v.Const {
				d.stmt(indent+2, s)
			}
		}

		if len(v.Decls) > 0 {
			d.line(indent+1, "Decls")
			for _, decl := range v.Decls {
				d.decl(indent+2, decl)
			}
		}

		if len(v.Structs) > 0 {
			d.line(indent+1, "Structs")
			for _, v := range v.Structs {
				d.stmt(indent+2, v)
			}
		}

		if len(v.Interfaces) > 0 {
			d.line(indent+1, "Interfaces")
			for _, v := range v.Interfaces {
				d.stmt(indent+2, v)
			}
		}

		if len(v.Implements) > 0 {
			d.line(indent+1, "Implements")
			for _, v := range v.Implements {
				d.node(indent+2, v)
			}
		}

		if len(v.Enums) > 0 {
			d.line(indent+1, "Enums")
			for _, v := range v.Enums {
				d.node(indent+2, v)
			}
		}

		if len(v.Sums) > 0 {
			d.line(indent+1, "Sums")
			for _, v := range v.Sums {
				d.node(indent+2, v)
			}
		}

		if len(v.Comptime) > 0 {
			d.line(indent+1, "ComptimeStmt")
			for _, v := range v.Comptime {
				d.node(indent+2, v)
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

	case *ConstDeclStmt:
		d.line(indent, "ConstDecl")
		d.kv(indent+1, "Const", v.ConstKW)
		d.kv(indent+1, "Name", v.Name)

		d.line(indent+1, "Type")
		if v.Type == nil {
			d.line(indent+2, "(none)")
		} else {
			d.typ(indent+2, v.Type)
		}

		d.kv(indent+1, "Eq", v.Eq)
		d.line(indent+1, "Init")
		if v.Init == nil {
			d.line(indent+2, "(none)")
		} else {
			d.expr(indent+2, v.Init)
		}

	case *VarDeclStmt:
		d.line(indent, "VarDeclStmt")
		d.kv(indent+1, "Var", v.VarKW)
		d.kv(indent+1, "Name", v.Name)

		d.line(indent+1, "Type")
		if v.Type == nil {
			d.line(indent+2, "(none)")
		} else {
			d.typ(indent+2, v.Type)
		}

		d.kv(indent+1, "Eq", v.Eq)
		d.line(indent+1, "Init")
		if v.Init == nil {
			d.line(indent+2, "(none)")
		} else {
			d.expr(indent+2, v.Init)
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
		d.kv(indent+1, "For", v.For)
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
		d.kv(indent+1, "For", v.For)
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
		d.kv(indent, "Break", v.Break)

	case *ContinueStmt:
		d.kv(indent, "Continue", v.Continue)

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
				d.typ(indent+2, f.Type)
				if f.Eq != nil {
					d.kv(indent+1, "Eq", *f.Eq)
					d.expr(indent+1, f.Default)
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

	case *SumType:
		d.kv(indent, "Type", v.TypeDecl)
		d.kv(indent+1, "Name", v.Name)
		d.kv(indent+1, "Sum", v.Sum)
		if v.Public {
			d.line(indent, "Public: true")
		}
		d.kv(indent+1, "LBrace", v.LBrace)
		if len(v.Variants) > 0 {
			d.line(indent+2, "Variants")
			for _, p := range v.Variants {
				d.kv(indent+3, "Ident", p)
			}
		}
		if len(v.Methods) > 0 {
			d.line(indent+2, "VariantMethods")
			for _, f := range v.Methods {
				d.kv(indent+3, "Methods", f.Name)
				d.line(indent+4, "Params")
				for _, p := range f.Params {
					d.line(indent+5, "Param")
					d.kv(indent+6, "Ident", p.Name)
					d.line(indent+6, "Type")
					d.typ(indent+7, p.Type)
				}
			}
		}
		d.kv(indent+1, "RBrace", v.RBrace)

	case *SliceExpr:
		d.expr(indent+1, v.X)
		d.kv(indent+2, "LBracket", v.LBracket)
		if v.Low != nil {
			d.expr(indent+1, v.Low)
		}
		if v.Colon != (token.Token{}) {
			d.kv(indent+2, "Colon", v.Colon)
		}
		if v.High != nil {
			d.expr(indent+1, v.High)
		}
		d.kv(indent+2, "RBracket", v.RBracket)

	case *TypeRef:
		d.sliceOrArrayType(indent+1, false, v.Parts)

	case *ComptimeType:
		d.kv(indent, "Comptime", v.Comptime)
		if v.Const != nil {
			d.node(indent, v.Const)
		}
		if v.Func != nil {
			d.node(indent, v.Func)
		}

	case *MapType:
		if v.KindKW.Kind == token.KWMap {
			d.kv(indent, "Map", v.KindKW)
		} else {
			d.kv(indent, "Hashmap", v.KindKW)
		}
		d.kv(indent, "LBracket", v.LBracket)
		d.line(indent, "KeyType:")
		for _, v := range v.KeyType.Parts {
			d.kv(indent+1, "Name", v)
		}
		d.kv(indent, "RBracket", v.RBracket)
		d.line(indent, "ValueType:")
		for _, v := range v.KeyType.Parts {
			d.kv(indent+1, "Name", v)
		}

	case *MakeExpr:
		d.kv(indent, "Make", v.MakeKW)
		d.kv(indent, "LParen", v.LParen)
		d.node(indent, v.Type)
		if len(v.Args) > 0 {
			d.line(indent, "Size:")
			d.node(indent+1, v.Args[0])
			if len(v.Args) == 2 {
				d.line(indent, "Cap:")
				d.node(indent+1, v.Args[1])
			}
		}
		d.kv(indent, "RParen", v.RParen)

	case *SliceElementsExpr:
		d.node(indent, new(v.Type))
		d.kv(indent, "LBrace", v.LBrace)
		if len(v.Elements) > 0 {
			d.line(indent+1, "Elements")
			for _, x := range v.Elements {
				d.node(indent+2, x)
			}
		}
		d.kv(indent, "RBrace", v.RBrace)

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

func (d *dumper) decl(indent int, n Decl) {
	switch v := n.(type) {
	case *FuncDecl, *BadDecl:
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
	case *NameType:
		d.node(indent, v)

	case *BadType:
		d.node(indent, v)

	case *TypeRef, *MapType:
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
	case *BlockStmt, *ConstDeclStmt, *VarDeclStmt, *AssignStmt, *ExprStmt, *BadStmt:
		d.node(indent, v)

	case *ReturnStmt, *IfStmt, *ForStmt, *RangeStmt, *IncDecStmt:
		d.node(indent, v)

	case *BreakStmt, *ContinueStmt, *SwitchStmt, *FallThroughStmt:
		d.node(indent, v)

	case *StructType, *InterfaceType, *EnumType, *SumType, *ComptimeType:
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

	case *IndexExpr, *CallExpr, *SliceExpr, *MakeExpr, *SliceElementsExpr:
		d.node(indent, v)

	default:
		if n == nil {
			d.line(indent, "(nil expr)")
			return
		}
		d.line(indent, fmt.Sprintf("<<unhandled expr %T>>", n))
	}
}

// sliceType returns slice type
func (d *dumper) sliceOrArrayType(indent int, printType bool, t []token.Token) {
	// if printType {
	// 	d.line(indent, "Type:")
	// }
	x := indent + 1
	for _, v := range t {
		switch v.Kind {
		case token.LBracket:
			d.kv(x, "LBracket", v)
		case token.IntLit:
			d.kv(x, "Size", v)
		case token.RBracket:
			d.kv(x, "RBracket", v)
		case token.Dot:
			d.kv(x, "Dot", v)
		default:
			d.kv(x, "Ident", v)
		}
	}
}
