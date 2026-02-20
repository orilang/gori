package parser

import (
	"fmt"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
)

// parseIfStmtExpr returns expressions for parseStmt func
func (p *Parser) parseIfStmtExpr() ast.Stmt {
	ifs := p.expect(token.KWIf, "expected 'if'")
	stmt := &ast.IfStmt{
		If: ifs,
	}

	if p.kind() == token.LBrace {
		p.errors = append(p.errors, fmt.Errorf("%d:%d: missing condition, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
		return &ast.BadStmt{From: ifs, To: p.peek(), Reason: "missing condition after 'if'"}
	}

	stmt.Condition = p.parseExpr(LOWEST)
	if token.IsAssignment(p.kind()) {
		tok := p.next()
		p.errors = append(p.errors, fmt.Errorf("%d:%d: assignment not allowed in if condition, got %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
		return &ast.BadStmt{From: ifs, To: tok, Reason: "assignment not allowed in if condition; use =="}
	}

	stmt.Then = p.parseBlock()

	if p.kind() == token.KWElse {
		_ = p.next()
		if p.kind() == token.KWIf {
			stmt.Else = p.parseIfStmtExpr()
		} else if p.kind() == token.LBrace {
			stmt.Else = p.parseBlock()
		} else {
			tok := p.next()
			p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
			return &ast.BadStmt{From: ifs, To: tok, Reason: "expected expression '{' or 'if' after 'else'"}
		}
	}

	return stmt
}
