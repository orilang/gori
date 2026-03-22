package parser

import (
	"fmt"
	"slices"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
)

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

		if p.newlineSincePrev() {
			break
		}
	}

	return tp
}

// parseSliceElements returns slice elements
func (p *Parser) parseSliceElements() ast.SliceElementsExpr {
	se := ast.SliceElementsExpr{
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
