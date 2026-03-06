package parser

import (
	"fmt"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
)

// parseComptimeStmt returns comptime expression
func (p *Parser) parseComptimeStmt() ast.Stmt {
	x := p.expect(token.KWComptime, "expected 'comptime'")
	if p.kind() == token.KWConst {
		return &ast.ComptimeType{
			Comptime: x,
			Const:    p.parseConstDecl(),
		}
	}

	if p.kind() == token.KWFunc {
		return &ast.ComptimeType{
			Comptime: x,
			Func:     p.parseFuncDecl(),
		}
	}

	p.errors = append(p.errors, fmt.Errorf("%d:%d: expected 'const' or 'func', got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))

	p.consumeTo(token.EOF)
	return nil
}
