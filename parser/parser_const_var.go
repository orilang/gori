package parser

import (
	"fmt"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
)

// parseConstDecl returns constant declaration
func (p *Parser) parseConstDecl() ast.Stmt {
	kw := p.expect(token.KWConst, "expected 'const'")
	name := p.expect(token.Ident, "expected constant name")

	typ, btyp, bad := p.parseVarConstType()
	if bad {
		p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", btyp.From.Line, btyp.From.Column, btyp.From.Kind, btyp.From.Value))
		return btyp
	}
	eq := p.expect(token.Assign, "expected '=")
	init := p.parseExpr(LOWEST)

	return &ast.ConstDeclStmt{
		ConstKW: kw,
		Name:    name,
		Type:    typ,
		Eq:      eq,
		Init:    init,
	}
}

// parseVarDecl returns variable declaration
func (p *Parser) parseVarDecl() ast.Stmt {
	kw := p.expect(token.KWVar, "expected 'var'")
	name := p.expect(token.Ident, "expected variable name")

	typ, btyp, bad := p.parseVarConstType()
	if bad {
		return btyp
	}
	eq := p.expect(token.Assign, "expected '=")
	init := p.parseExpr(LOWEST)

	return &ast.VarDeclStmt{
		VarKW: kw,
		Name:  name,
		Type:  typ,
		Eq:    eq,
		Init:  init,
	}
}

// parseVarConstType returns const/vars types
func (p *Parser) parseVarConstType() (*ast.NameType, *ast.BadType, bool) {
	typ := &ast.NameType{}
	btyp := &ast.BadType{}
	var bad bool
	tok := p.peek()

	if token.IsVarConstTypes(tok.Kind) {
		typ.Name = tok
	} else {
		p.errors = append(p.errors, fmt.Errorf("%d:%d: unsupported type with %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
		btyp.From = tok
		btyp.Reason = "unexpected type name"
		bad = true
	}

	_ = p.next()
	return typ, btyp, bad
}
