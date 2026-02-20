package parser

import (
	"fmt"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
)

// parseImplementsDecl is in charge of parsing interface "implements" requirements
func (p *Parser) parseImplementsDecl() *ast.ImplementsDecl {
	kw := p.expect(token.Ident, "expected 'ident'")
	kwi := p.expect(token.KWImplements, "expected 'implements'")
	id := &ast.ImplementsDecl{
		Type:       kw,
		Implements: kwi,
	}

	if p.kind() != token.Ident {
		p.errors = append(p.errors, fmt.Errorf("%d:%d: expected ident after 'implements', got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
		p.consumeTo(token.RBrace)
		return id
	}

	id.Interface = p.parseInterfaceTypeEmbbed()
	if p.kind() == token.SemiComma {
		_ = p.next()
	}
	return id
}
