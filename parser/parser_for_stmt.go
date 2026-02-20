package parser

import (
	"fmt"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
)

// parseForStmtExpr returns expressions for parseStmt func
func (p *Parser) parseForStmtExpr() ast.Stmt {
	ftok := p.expect(token.KWFor, "expected 'for'")
	fstmt := &ast.ForStmt{
		For: ftok,
	}
	rstmt := &ast.RangeStmt{
		For: ftok,
	}

	// infinite loop
	if p.kind() == token.LBrace {
		fstmt.Body = p.parseForBlockStmt()
		return fstmt
	}

	if p.lookForInForHeader(token.KWRange) {
		// for range x {}
		if p.kind() == token.KWRange {
			rg := p.expect(token.KWRange, "expected 'range'")
			rstmt.Range = rg
			if p.kind() == token.LBrace {
				p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
				return &ast.BadStmt{From: ftok, To: p.peek(), Reason: "expected expression before '{'"}
			}
			rstmt.X = p.parseExpr(LOWEST)

			if p.kind() != token.LBrace {
				p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
				return &ast.BadStmt{From: ftok, To: p.peek(), Reason: "expected '{' after expression"}
			}

			rstmt.Body = p.parseForBlockStmt()
			return rstmt
		}

		// var k1, k2 ast.Expr
		if p.kind() == token.Comma {
			tok := p.next()
			p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
			return &ast.BadStmt{From: ftok, To: tok, Reason: "expected expression not ','"}
		}

		k1 := p.expect(token.Ident, "expected 'identifier'")
		if k1.Kind != token.Ident {
			return &ast.BadStmt{From: ftok, To: k1, Reason: "expected identifier"}
		}
		xk1 := &ast.IdentExpr{Name: k1}

		// for k,v := range x {}
		if p.kind() == token.Comma {
			_ = p.next()
			rstmt.Key = xk1
			k2 := p.expect(token.Ident, "expected 'identifier'")
			if k2.Kind != token.Ident {
				return &ast.BadStmt{From: ftok, To: k2, Reason: "expected identifier"}
			}

			rstmt.Value = &ast.IdentExpr{Name: k2}
			if !token.IsRangeForAssignment(p.kind()) {
				p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
				return &ast.BadStmt{From: ftok, To: p.peek(), Reason: "expected '=' or ':=' after expression"}
			}
			op := p.next()
			rstmt.Op = op
			rg := p.expect(token.KWRange, "expected 'range'")
			rstmt.Range = rg
			if p.kind() == token.LBrace {
				p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
				return &ast.BadStmt{From: ftok, To: p.peek(), Reason: "expected expression before '{'"}
			}
			rstmt.X = p.parseExpr(LOWEST)

			if p.kind() != token.LBrace {
				p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
				return &ast.BadStmt{From: ftok, To: p.peek(), Reason: "expected '{' after expression"}
			}

			rstmt.Body = p.parseForBlockStmt()
			return rstmt
		}

		// for v := range x {}
		if !token.IsRangeForAssignment(p.kind()) {
			p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
			return &ast.BadStmt{From: ftok, To: p.peek(), Reason: "expected '=' or ':=' after expression"}
		}

		rstmt.Key = xk1
		op := p.next()
		rstmt.Op = op
		rg := p.expect(token.KWRange, "expected 'range'")
		rstmt.Range = rg
		if p.kind() == token.LBrace {
			p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
			return &ast.BadStmt{From: ftok, To: p.peek(), Reason: "expected expression before '{'"}
		}
		rstmt.X = p.parseExpr(LOWEST)

		if p.kind() != token.LBrace {
			p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
			return &ast.BadStmt{From: ftok, To: p.peek(), Reason: "expected '{' after expression"}
		}

		rstmt.Body = p.parseForBlockStmt()
		return rstmt
	}

	// for condition {}
	if !p.lookForInForHeader(token.SemiComma) {
		fstmt.Condition = p.parseExpr(LOWEST)
		if p.kind() != token.LBrace {
			p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
			return &ast.BadStmt{From: ftok, To: p.peek(), Reason: "expected '{' after condition"}
		}

		fstmt.Body = p.parseForBlockStmt()
		return fstmt
	}

	// for init; condition; post {}
	fstmt.Init = p.parseSimpleStmt()
	if !isValidForInit(fstmt.Init) {
		p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected init statement in for header, got %v %q", fstmt.Init.Start().Line, fstmt.Init.Start().Column, fstmt.Init.Start().Kind, fstmt.Init.Start().Value))
		return &ast.BadStmt{From: fstmt.Init.Start(), To: fstmt.Init.End(), Reason: "invalid init statement in for header"}
	}
	sm1 := p.expect(token.SemiComma, "expected ';'")
	if sm1.Kind != token.SemiComma {
		return &ast.BadStmt{From: ftok, To: sm1, Reason: "expected ';'"}
	}

	fstmt.Condition = p.parseExpr(LOWEST)
	sm2 := p.expect(token.SemiComma, "expected ';'")
	if sm2.Kind != token.SemiComma {
		return &ast.BadStmt{From: ftok, To: sm2, Reason: "expected ';'"}
	}

	if p.kind() == token.LBrace {
		p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
		return &ast.BadStmt{From: ftok, To: p.peek(), Reason: "expected expression before '{'"}
	}

	fstmt.Post = p.parseSimpleStmt()
	if !isValidForPost(fstmt.Post) {
		p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected post statement in for header, got %v %q", fstmt.Post.Start().Line, fstmt.Post.Start().Column, fstmt.Post.Start().Kind, fstmt.Post.Start().Value))
		return &ast.BadStmt{From: fstmt.Post.Start(), To: fstmt.Post.End(), Reason: "invalid post statement in for header"}
	}
	if p.kind() != token.LBrace {
		tok := p.next()
		p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
		return &ast.BadStmt{From: fstmt.Post.Start(), To: tok, Reason: "expected '{'"}
	}

	fstmt.Body = p.parseForBlockStmt()
	return fstmt
}

// isValidForInit is only used in for init statement
func isValidForInit(s ast.Stmt) bool {
	switch st := s.(type) {
	case *ast.AssignStmt:
		return true
	case *ast.ExprStmt:
		_, ok := st.Expr.(*ast.CallExpr)
		return ok
	default:
		return false
	}
}

// isValidForPost is only used in for post statement
func isValidForPost(s ast.Stmt) bool {
	switch st := s.(type) {
	case *ast.AssignStmt, *ast.IncDecStmt:
		return true
	case *ast.ExprStmt:
		_, ok := st.Expr.(*ast.CallExpr)
		return ok
	default:
		return false
	}
}

// parseForBlockStmt is an helper to make sure we inc/dec
// loopDepth counter
func (p *Parser) parseForBlockStmt() *ast.BlockStmt {
	p.loopDepth++
	defer func() {
		p.loopDepth--
	}()
	return p.parseBlock()
}

// parseBreakStmt returns expressions for parseStmt func
func (p *Parser) parseBreakStmt() ast.Stmt {
	if p.loopDepth == 0 {
		tok := p.next()
		p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected break expression outside for loop, got %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
		return &ast.BadStmt{From: tok, Reason: "expected 'break' inside 'for' loop"}
	}

	kw := p.expect(token.KWBreak, "expected 'break'")
	if p.kind() == token.Comment {
		_ = p.next()
	}
	// we do not accept labels for now so anything unauthorized is rejected
	// maybe labels will be supported later
	if p.kind() == token.RBrace || p.kind() == token.EOF || p.kind() == token.SemiComma || p.peek().Line > kw.Line {
		if p.kind() == token.SemiComma {
			_ = p.next()
		}
		return &ast.BreakStmt{
			Break: kw,
		}
	}

	p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected statement after 'break', got %v %q", kw.Line, kw.Column, p.peek().Kind, p.peek().Value))
	return &ast.BadStmt{From: p.peek(), Reason: "expected '}' or 'EOF' or new line"}
}

// parseContinueStmt returns expressions for parseStmt func
func (p *Parser) parseContinueStmt() ast.Stmt {
	if p.loopDepth == 0 {
		tok := p.next()
		p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected continue expression outside for loop, got %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
		return &ast.BadStmt{From: tok, Reason: "expected 'continue' inside 'for' loop"}
	}

	kw := p.expect(token.KWContinue, "expected 'continue'")
	if p.kind() == token.Comment {
		_ = p.next()
	}
	// any unauthorized statement is rejected after 'continue'
	if p.kind() == token.RBrace || p.kind() == token.EOF || p.kind() == token.SemiComma || p.peek().Line > kw.Line {
		if p.kind() == token.SemiComma {
			_ = p.next()
		}
		return &ast.ContinueStmt{
			Continue: kw,
		}
	}

	p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected statement after 'continue', got %v %q", kw.Line, kw.Column, p.peek().Kind, p.peek().Value))
	return &ast.BadStmt{From: p.peek(), Reason: "expected '}' or 'EOF' or new line"}
}
