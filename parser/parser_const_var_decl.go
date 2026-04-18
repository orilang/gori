package parser

import (
	"fmt"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
)

// parseConstDecl returns constant declaration
func (p *Parser) parseConstDecl() ast.Decl {
	kw := p.expect(token.KWConst, "expected 'const'")
	name := p.expectValidIdent(token.Ident, true, "expected constant name")

	typ, btyp, bad := p.parseVarConstType()
	if bad {
		p.Errors = append(p.Errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", btyp.From.Line, btyp.From.Column, btyp.From.Kind, btyp.From.Value))
		return btyp
	}
	eq := p.expect(token.Assign, "expected '=")
	var init ast.Expr
	switch p.peek().Value {
	case "[":
		init = p.parseSliceElements()
	default:
		init = p.parseExpr(LOWEST)
	}

	return &ast.ConstDecl{
		ConstKW: kw,
		Name:    name,
		Type:    typ,
		Eq:      eq,
		Init:    init,
	}
}

// parseVarDecl returns variable declaration
func (p *Parser) parseVarDecl() ast.Decl {
	kw := p.expect(token.KWVar, "expected 'var'")
	name := p.expectValidIdent(token.Ident, true, "expected variable name")
	var view token.Token
	if p.kind() == token.KWView {
		view = p.expect(token.KWView, "expected 'view'")
	}

	typ, btyp, bad := p.parseVarConstType()
	if bad {
		return btyp
	}
	eq := p.expect(token.Assign, "expected '=")
	var init ast.Expr
	switch p.peek().Value {
	case "make":
		init = p.parseMakeExpr()
	case "[":
		// []string{}
		init = p.parseSliceElements()
	default:
		// x[1:]
		if p.lookForInSliceHeader(token.LBracket) {
			init = p.parseSliceExpr(p.parsePrefix())
		} else {
			init = p.parseExpr(LOWEST)
		}
	}

	return &ast.VarDecl{
		VarKW: kw,
		Name:  name,
		View:  view,
		Type:  typ,
		Eq:    eq,
		Init:  init,
	}
}

// parseVarConstType returns const/vars types
func (p *Parser) parseVarConstType() (ast.Type, *ast.BadType, bool) {
	typ := &ast.NamedType{}
	btyp := &ast.BadType{}
	var bad bool

	switch {
	case p.kind() == token.LBracket:
		return p.parseSliceOrArrayType(), nil, false

	case token.IsMapType(p.kind()):
		return p.parseMapsHashMapsDecl(), nil, false

	case token.IsVarConstTypes(p.kind()):
		typ.Parts = append(typ.Parts, p.next())
	default:
		tok := p.next()
		p.Errors = append(p.Errors, fmt.Errorf("%d:%d: unsupported type with %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
		btyp.From = tok
		btyp.Reason = "unexpected type name"
		bad = true
	}

	return typ, btyp, bad
}
