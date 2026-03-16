package parser

import (
	"fmt"
	"slices"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
)

// parseSliceDecl returns parsed slice
func (p *Parser) parseSliceDecl() *ast.SliceType {
	st := &ast.SliceType{}
	if p.kind() == token.KWVar {
		kwt := p.expect(token.KWVar, "expected 'var'")
		kwi := p.expectValidIdent(token.Ident, true, "expected 'ident'")

		st.VarConstKW = kwt
		st.Name = kwi
		st.Type = p.parseSliceOrArrayType()

		if p.kind() == token.Assign {
			st.Eq = p.next()
			st.Elements = p.parseSliceElements()
		}
	} else if p.kind() == token.KWConst {
		kwt := p.expect(token.KWConst, "expected 'const'")
		kwi := p.expectValidIdent(token.Ident, true, "expected 'ident'")

		st.VarConstKW = kwt
		st.Name = kwi
		st.Type = p.parseSliceOrArrayType()

		if p.kind() != token.Assign {
			p.errors = append(p.errors, fmt.Errorf("%d:%d: expected '=', got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
			p.consumeTo(token.EOF)
			return st
		}

		st.Eq = p.next()
		st.Elements = p.parseSliceElements()
	}
	if p.kind() == token.SemiComma {
		_ = p.next()
	}

	return st
}

// parseSliceOrArrayType returns slice or array type
func (p *Parser) parseSliceOrArrayType() ast.TypeRef {
	var tp ast.TypeRef
	lb := p.expect(token.LBracket, "expected '['")
	if p.kind() == token.IntLit {
		size := p.expect(token.IntLit, "expected 'intLit'")
		rb := p.expect(token.RBracket, "expected ']'")
		tp.Parts = append(tp.Parts, lb, size, rb)
	} else {
		rb := p.expect(token.RBracket, "expected ']'")
		tp.Parts = append(tp.Parts, lb, rb)
	}
	kindList := []token.Kind{
		token.Comma,
		token.SemiComma,
		token.RParen,
		token.LBrace,
		token.Assign,
		token.EOF,
	}

	for !slices.Contains(kindList, p.kind()) {
		if !token.IsSliceType(p.kind()) {
			p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected slice/array type, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
			p.consumeTo(token.EOF)
			return tp
		}
		tp.Parts = append(tp.Parts, p.next())

		if p.kind() == token.SemiComma {
			_ = p.next()
			break
		}

		if p.newlineSincePrev() {
			break
		}
	}

	return tp
}

// parseSliceElements returns slice elements
func (p *Parser) parseSliceElements() ast.SliceElements {
	se := ast.SliceElements{
		Type: p.parseSliceOrArrayType(),
	}
	lb := p.expect(token.LBrace, "expected '{'")
	se.LBrace = lb

	for p.kind() != token.RBrace && p.kind() != token.EOF {
		se.Elements = append(se.Elements, p.parseExpr(LOWEST))

		if p.kind() == token.Comma {
			_ = p.next()
			continue
		}
		if p.kind() == token.RBrace {
			break
		}
	}

	rb := p.expect(token.RBrace, "expected '}'")
	se.RBrace = rb
	return se
}

// parseSliceViewDecl returns parsed slice view
func (p *Parser) parseSliceViewDecl() *ast.SliceViewType {
	kwt := p.expect(token.KWVar, "expected 'var'")
	kwi := p.expectValidIdent(token.Ident, true, "expected 'ident'")
	kwv := p.expect(token.KWView, "expected 'view'")

	st := &ast.SliceViewType{
		VarKW: kwt,
		Name:  kwi,
		View:  kwv,
		Type:  p.parseSliceOrArrayType(),
	}

	st.Eq = p.expect(token.Assign, "expected '='")
	st.Elements = p.parseSliceExpr(p.parsePrefix())

	if p.kind() == token.SemiComma {
		_ = p.next()
	}

	return st
}

// parseSliceExpr returns expressions for parsePostfix func
func (p *Parser) parseSliceExpr(left ast.Expr) *ast.SliceExpr {
	lb := p.expect(token.LBracket, "LBracket expected '['")
	x := &ast.SliceExpr{
		X:        left,
		LBracket: lb,
	}

	if p.kind() == token.Colon {
		x.Colon = p.next()
		if p.kind() != token.RBracket {
			x.High = p.parseExpr(LOWEST)
		}
		x.RBracket = p.expect(token.RBracket, "RBracket expected ']'")
		return x
	}

	x.Low = p.parseExpr(LOWEST)
	if p.kind() != token.Colon {
		x.RBracket = p.expect(token.RBracket, "RBracket expected ']'")
		return x
	}

	x.Colon = p.next()
	if p.kind() != token.RBracket {
		x.High = p.parseExpr(LOWEST)
	}
	x.RBracket = p.expect(token.RBracket, "RBracket expected ']'")
	return x
}

// parseArrayDecl returns parsed array
func (p *Parser) parseArrayDecl() *ast.ArrayType {
	st := &ast.ArrayType{}
	if p.kind() == token.KWVar {
		kwt := p.expect(token.KWVar, "expected 'var'")
		kwi := p.expectValidIdent(token.Ident, true, "expected 'ident'")

		st.VarConstKW = kwt
		st.Name = kwi
		st.Type = p.parseSliceOrArrayType()

		if p.kind() == token.Assign {
			st.Eq = p.next()
			st.Elements = p.parseSliceElements()
		}
	} else if p.kind() == token.KWConst {
		kwt := p.expect(token.KWConst, "expected 'const'")
		kwi := p.expectValidIdent(token.Ident, true, "expected 'ident'")

		st.VarConstKW = kwt
		st.Name = kwi
		st.Type = p.parseSliceOrArrayType()

		if p.kind() != token.Assign {
			p.errors = append(p.errors, fmt.Errorf("%d:%d: expected '=', got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
			p.consumeTo(token.EOF)
			return st
		}

		st.Eq = p.next()
		st.Elements = p.parseSliceElements()
	}
	if p.kind() == token.SemiComma {
		_ = p.next()
	}

	return st
}

// parseArrayViewDecl returns parsed array view
func (p *Parser) parseArrayViewDecl() *ast.ArrayViewType {
	kwt := p.expect(token.KWVar, "expected 'var'")
	kwi := p.expectValidIdent(token.Ident, true, "expected 'ident'")
	kwv := p.expect(token.KWView, "expected 'view'")

	st := &ast.ArrayViewType{
		VarKW: kwt,
		Name:  kwi,
		View:  kwv,
		Type:  p.parseSliceOrArrayType(),
	}

	st.Eq = p.expect(token.Assign, "expected '='")
	st.Elements = p.parseSliceExpr(p.parsePrefix())

	if p.kind() == token.SemiComma {
		_ = p.next()
	}

	return st
}
