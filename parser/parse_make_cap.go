package parser

import (
	"fmt"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
)

// parseMakeExpr returns make parsed expression
func (p *Parser) parseMakeExpr() *ast.MakeExpr {
	kw := p.expect(token.Ident, "expected 'make'")
	lp := p.expect(token.LParen, "expected '(")

	x := &ast.MakeExpr{
		MakeKW: kw,
		LParen: lp,
	}

	if !token.IsMakeTypes(p.kind()) {
		p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected map/hashmap or slice, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
		p.consumeTo(token.EOF)
		return x
	}

	var count int
	for p.kind() != token.RParen && p.kind() != token.EOF {
		if token.IsMapType(p.kind()) {
			x.Type = p.parseMapsHashMapsDecl()
		} else if p.lookForInSliceHeader(token.LBracket) {
			x.Type = new(p.parseSliceOrArrayType())
		}

		if p.kind() == token.Comma {
			count++
			_ = p.expect(token.Comma, "expected ','")
		}

		if count > 2 {
			p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected map/hashmap or slice, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
			p.consumeTo(token.RParen)
			return x
		}

		if p.kind() != token.RParen {
			k := p.parseExpr(LOWEST)
			if k.Start().Kind != token.IntLit {
				p.errors = append(p.errors, fmt.Errorf("%d:%d: expected intLit, got %v %q", k.Start().Line, k.Start().Column, k.Start().Kind, k.Start().Value))
				p.consumeTo(token.RParen)
				return x
			}
			x.Args = append(x.Args, k)
		}

		if p.kind() == token.RParen {
			break
		}
	}

	x.RParen = p.expect(token.RParen, "expected ')")

	if p.kind() == token.SemiComma {
		_ = p.expect(token.SemiComma, "expected ';'")
	}

	return x
}
