package parser

import (
	"fmt"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
)

// parseEnumDecl is in charge of parsing enum type
func (p *Parser) parseEnumDecl() *ast.EnumType {
	kwp := p.expect(token.KWType, "expected 'type'")
	kwi := p.expect(token.Ident, "expected 'ident'")
	kwe := p.expect(token.KWEnum, "expected 'enum'")
	lbrace := p.expect(token.LBrace, "expected '{'")

	ed := &ast.EnumType{
		TypeDecl: kwp,
		Name:     kwi,
		Public:   isPublic(kwi),
		Enum:     kwe,
		LBrace:   lbrace,
	}

	for p.kind() != token.RBrace && p.kind() != token.EOF {
		if p.kind() != token.Ident {
			p.errors = append(p.errors, fmt.Errorf("%d:%d: expected 'ident', got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
			p.consumeTo(token.RBrace)
			return ed
		}

		if p.kind() == token.Ident {
			ed.Variants = append(ed.Variants, p.next())
		}

		if p.kind() == token.Comment {
			_ = p.next()
		}

		if p.kind() == token.SemiComma {
			_ = p.next()
			continue
		}

		if p.kind() == token.RBrace {
			break
		}

		if p.newlineSincePrev() {
			continue
		}

		p.errors = append(p.errors, fmt.Errorf("%d:%d: expected ';' or newline after ident, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
		p.consumeTo(token.RBrace)
	}

	if len(ed.Variants) == 0 {
		p.errors = append(p.errors, fmt.Errorf("%d:%d: expected 'ident' inside braces, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
		p.consumeTo(token.RBrace)
		return ed
	}

	rbrace := p.expect(token.RBrace, "expected '}'")
	ed.RBrace = rbrace

	return ed
}
