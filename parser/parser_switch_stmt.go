package parser

import (
	"fmt"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
)

func (p *Parser) parseSwitchStmt() ast.Stmt {
	kw := p.expect(token.KWSwitch, "expected 'switch'")
	s := ast.SwitchStmt{
		Switch: kw,
	}

	if p.lookForInSwitchHeader(token.SemiComma) {
		s.Init = p.parseSimpleStmt()
		if p.kind() == token.SemiComma {
			_ = p.next()
		}

		if p.kind() != token.LBrace {
			s.Tag = p.parseExpr(LOWEST)
		}
		lb := p.expect(token.LBrace, "expected '{'")
		s.LBrace = lb

		return p.parseSwitchCasesStmt(s)
	}

	if p.kind() != token.LBrace {
		s.Tag = p.parseExpr(LOWEST)
	}
	lb := p.expect(token.LBrace, "expected '{'")
	s.LBrace = lb

	return p.parseSwitchCasesStmt(s)
}

func (p *Parser) parseSwitchCasesStmt(s ast.SwitchStmt) ast.Stmt {
	var dcount int
	for p.kind() != token.RBrace && p.kind() != token.EOF {
		switch p.kind() {
		case token.KWCase:
			ckw := p.expect(token.KWCase, "expected 'case'")
			scase := ast.CaseClause{
				Case: ckw,
			}

			if p.lookForInSwitchCaseHeader(token.Comma) {
				for p.kind() != token.Colon && p.kind() != token.RBrace && p.kind() != token.EOF {
					if p.kind() == token.Comma {
						p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
						return &ast.BadStmt{From: s.Switch, To: p.peek(), Reason: "expected expression before ','"}
					}

					scase.Values = append(scase.Values, p.parseExpr(LOWEST))
					if p.kind() == token.Colon {
						colon := p.expect(token.Colon, "expected ':'")
						scase.Colon = colon

						for p.kind() != token.KWCase && p.kind() != token.KWDefault && p.kind() != token.RBrace && p.kind() != token.EOF {
							scase.Body = append(scase.Body, p.parseStmt())
						}
						break
					}

					if p.kind() != token.Comma && p.kind() != token.Colon {
						p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
						return &ast.BadStmt{From: s.Switch, To: p.peek(), Reason: "expected ','"}
					}
					_ = p.expect(token.Comma, "expected ','")
				}
			} else {
				for p.kind() != token.Colon && p.kind() != token.RBrace && p.kind() != token.EOF {
					scase.Values = append(scase.Values, p.parseExpr(LOWEST))
					if p.kind() == token.Colon {
						colon := p.expect(token.Colon, "expected ':'")
						scase.Colon = colon

						for p.kind() != token.KWCase && p.kind() != token.KWDefault && p.kind() != token.RBrace && p.kind() != token.EOF {
							scase.Body = append(scase.Body, p.parseStmt())
						}
						break
					}

					if p.kind() != token.Colon {
						p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
						return &ast.BadStmt{From: s.Switch, To: p.peek(), Reason: "expected ','"}
					}
				}
			}
			s.Cases = append(s.Cases, scase)

		case token.KWDefault:
			dcount++
			if dcount > 1 {
				p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected 'default', got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
				return &ast.BadStmt{From: s.Switch, To: p.peek(), Reason: "expected only one 'default' case"}
			}

			dkw := p.expect(token.KWDefault, "expected 'default'")
			if p.kind() == token.Colon {
				colon := p.expect(token.Colon, "expected ':'")
				dcase := ast.CaseClause{
					Case:  dkw,
					Colon: colon,
				}

				// empty body
				if p.kind() == token.KWCase {
					s.Cases = append(s.Cases, dcase)
				} else {
					for p.kind() != token.KWCase && p.kind() != token.KWDefault && p.kind() != token.RBrace && p.kind() != token.EOF {
						dcase.Body = append(dcase.Body, p.parseStmt())
					}
					s.Cases = append(s.Cases, dcase)
				}
			}

			if p.kind() != token.KWDefault && p.kind() != token.KWCase && p.kind() != token.RBrace && p.kind() != token.EOF {
				p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected ':', got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
				return &ast.BadStmt{From: s.Switch, To: p.peek(), Reason: "expected ':' after 'default'"}
			}

		default:
			p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
			return &ast.BadStmt{From: s.Switch, To: p.peek(), Reason: "expected 'case' or 'default'"}
		}
	}

	rb := p.expect(token.RBrace, "expected '}'")
	s.RBrace = rb

	return &s
}

// parseFallThroughStmt returns expressions for parseStmt func
func (p *Parser) parseFallThroughStmt() ast.Stmt {
	kw := p.expect(token.KWFallThrough, "expected 'fallthrough'")
	if p.kind() == token.Comment {
		_ = p.next()
	}

	// any unauthorized statement is rejected after 'fallthrough'
	if p.kind() == token.RBrace || p.kind() == token.EOF || p.kind() == token.KWCase || p.kind() == token.KWDefault {
		return &ast.FallThroughStmt{
			FallThroughStmt: kw,
		}
	}

	p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected statement after 'fallthrough', got %v %q", kw.Line, kw.Column, p.peek().Kind, p.peek().Value))
	return &ast.BadStmt{From: p.peek(), Reason: "expected '}' or 'EOF' or new line"}
}
