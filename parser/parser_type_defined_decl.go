package parser

import (
	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
)

// parseDefinedDecl returns parsed defined types
func (p *Parser) parseDefinedDecl() ast.Decl {
	kwt := p.expect(token.KWType, "expected 'type'")
	kwi := p.expectValidIdent(token.Ident, true, "expected 'ident'")

	dt := &ast.DefinedTypeDecl{
		TypeDecl: kwt,
		Name:     kwi,
	}

	x := &ast.NamedType{}
	x.Parts = append(x.Parts, p.next())
	dt.Type = x

	if p.kind() == token.SemiComma {
		_ = p.expect(token.SemiComma, "expected ';'")
	}
	return dt
}
