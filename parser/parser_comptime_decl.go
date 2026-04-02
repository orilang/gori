package parser

import (
	"fmt"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
)

// parseComptimeBlockDecl returns comptime expression
func (p *Parser) parseComptimeBlockDecl() ast.Decl {
	x := p.expect(token.KWComptime, "expected 'comptime'")
	if p.kind() == token.KWConst {
		c := &ast.ComptimeBlockDecl{
			ComptimeKW: x,
		}
		c.Decls = append(c.Decls, p.parseConstDecl())
		return c
	}

	if p.kind() == token.KWFunc {
		c := &ast.ComptimeBlockDecl{
			ComptimeKW: x,
		}
		c.Decls = append(c.Decls, p.parseFuncDecl())
		return c
	}

	p.errors = append(p.errors, fmt.Errorf("%d:%d: expected 'const' or 'func', got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))

	p.consumeTo(token.EOF)
	return nil
}
