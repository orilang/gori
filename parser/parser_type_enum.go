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
	kweq := p.expect(token.Assign, "expected '='")

	ed := &ast.EnumType{
		TypeDecl: kwp,
		Name:     kwi,
		Public:   isPublic(kwi),
		Eq:       kweq,
		Enum:     kwe,
	}

	if p.kind() != token.Pipe {
		p.errors = append(p.errors, fmt.Errorf("%d:%d: expected '| after '=', got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
		p.consumeTo(token.EOF)
		return ed
	}

	for p.kind() != token.EOF {
		if p.kind() == token.Pipe {
			_ = p.next()
		}

		if p.kind() != token.Ident {
			p.errors = append(p.errors, fmt.Errorf("%d:%d: expected ident after '|', got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
			p.consumeTo(token.EOF)
		}

		ed.Variants = append(ed.Variants, p.next())

		if p.kind() == token.Comment {
			_ = p.next()
		}

		if p.kind() != token.Pipe {
			break
		}
	}

	return ed
}
