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

	if token.IsMapType(p.kind()) {
		dt.Type = p.parseMapsHashMapsDecl()
	} else if p.kind() == token.LBracket {
		dt.Type = p.parseSliceOrArrayType()
	} else {
		dt.Type = &ast.NamedType{
			Parts: []token.Token{p.next()},
		}
	}

	if p.kind() == token.SemiComma {
		_ = p.expect(token.SemiComma, "expected ';'")
	}
	return dt
}
