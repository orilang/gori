package parser

import (
	"fmt"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
)

// parseReturnStmtExpr returns expressions for parseStmt func
func (p *Parser) parseReturnStmtExpr() ast.Stmt {
	rn := p.expect(token.KWReturn, "expected 'return'")

	// return has no values
	if p.kind() == token.EOF || p.kind() == token.RBrace {
		return &ast.ReturnStmt{
			Return: rn,
		}
	}

	var args []ast.Expr
	for p.kind() != token.RBrace && p.kind() != token.EOF {
		if p.kind() == token.Comma {
			tok := p.next()
			p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
			return &ast.BadStmt{From: rn, To: tok, Reason: "expected expression not ','"}
		}

		args = append(args, p.parseExpr(LOWEST))
		if p.kind() == token.KWCase || p.kind() == token.KWDefault {
			break
		}

		if p.kind() != token.Comma && p.kind() != token.RBrace && p.kind() != token.EOF {
			tok := p.next()
			p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
			return &ast.BadStmt{From: rn, To: tok, Reason: "expected ',' or '}' after return value"}
		}

		if p.kind() == token.Comma {
			_ = p.next()
			if p.kind() == token.RBrace || p.kind() == token.EOF {
				p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
				return &ast.BadStmt{From: rn, To: p.peek(), Reason: "expected expression after ','"}
			}
		}
	}

	return &ast.ReturnStmt{
		Return: rn,
		Values: args,
	}
}
