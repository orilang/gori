package parser

import (
	"fmt"
	"slices"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
)

// parseSliceOrArrayType returns slice or array type
func (p *Parser) parseSliceOrArrayType() ast.Type {
	var (
		slice ast.SliceType
		array ast.ArrayType
		nt    ast.NamedType
	)
	kindList := []token.Kind{
		token.Comma,
		token.SemiComma,
		token.RParen,
		token.LBrace,
		token.Assign,
		token.EOF,
	}

	lb := p.expect(token.LBracket, "expected '['")
	if p.kind() == token.IntLit {
		array.LBracket = lb
		size := p.parseExpr(LOWEST)
		rb := p.expect(token.RBracket, "expected ']'")
		array.Len = size
		array.RBracket = rb

		for !slices.Contains(kindList, p.kind()) {
			if !token.IsSliceType(p.kind()) {
				p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected array type, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
				p.consumeTo(token.EOF)
				return &array
			}
			nt.Parts = append(nt.Parts, p.next())

			if p.newlineSincePrev() {
				break
			}
		}

		array.Elem = &nt
		return &array
	}

	rb := p.expect(token.RBracket, "expected ']'")
	slice.LBracket = lb
	slice.RBracket = rb

	for !slices.Contains(kindList, p.kind()) {
		if !token.IsSliceType(p.kind()) {
			p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected slice type, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
			p.consumeTo(token.EOF)
			return &slice
		}
		nt.Parts = append(nt.Parts, p.next())

		if p.newlineSincePrev() {
			break
		}
	}

	slice.Elem = &nt
	return &slice
}

// parseSliceElements returns slice elements
func (p *Parser) parseSliceElements() *ast.SliceLitExpr {
	se := &ast.SliceLitExpr{
		Type: p.parseSliceOrArrayType(),
	}
	lb := p.expect(token.LBrace, "expected '{'")
	se.LBrace = lb

	for p.kind() != token.RBrace && p.kind() != token.EOF {
		se.Elements = append(se.Elements, p.parseExpr(LOWEST))

		if p.kind() == token.Comma {
			_ = p.expect(token.Comma, "expected ','")
			continue
		}
		if p.kind() == token.RBrace {
			break
		}
	}

	rb := p.expect(token.RBrace, "expected '}'")
	se.RBrace = rb

	if p.kind() == token.SemiComma {
		_ = p.expect(token.SemiComma, "expected ';'")
	}

	return se
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

	if p.kind() == token.SemiComma {
		_ = p.expect(token.SemiComma, "expected ';'")
	}

	return x
}
