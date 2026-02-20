package parser

import (
	"fmt"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
)

// parseFuncDecl returns function declaration
func (p *Parser) parseFuncDecl() ast.Decl {
	kw := p.expect(token.KWFunc, "expected 'func'")
	name := p.expect(token.Ident, "expected function name")
	_ = p.expect(token.LParen, "expected '(' after function name")

	f := &ast.FuncDecl{
		FuncKW: kw,
		Name:   name,
	}
	for p.kind() != token.RParen && p.kind() != token.EOF {
		if p.kind() == token.Comma {
			tok := p.next()
			p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
			return &ast.BadDecl{From: kw, To: tok, Reason: "expected expression not ','"}
		}
		f.Params = append(f.Params, p.parseFuncParam())
		if p.kind() != token.Comma && p.kind() != token.RParen && p.kind() != token.EOF {
			tok := p.next()
			p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
			return &ast.BadDecl{From: kw, To: tok, Reason: "expected ',' or ')'"}
		}

		if p.kind() == token.Comma {
			_ = p.next()
			if p.kind() == token.RParen || p.kind() == token.EOF {
				p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
				return &ast.BadDecl{From: kw, To: p.peek(), Reason: "expected expression after ','"}
			}
		}
	}

	_ = p.expect(token.RParen, "expected ')' after function name")

	f.Results = p.parseFuncReturnTypes()
	body := p.parseBlock()
	f.Body = body

	return f
}

// parseFuncParam returns function parameter
func (p *Parser) parseFuncParam() ast.Param {
	name := p.expect(token.Ident, "expected parameter identifier")
	typ, btyp, bad := p.parseFuncParamType()
	if bad {
		return ast.Param{Name: name, Type: btyp}
	}
	return ast.Param{Name: name, Type: typ}
}

// parseFuncParamType returns func parameter type
func (p *Parser) parseFuncParamType() (*ast.NameType, *ast.BadType, bool) {
	typ := &ast.NameType{}
	btyp := &ast.BadType{}
	var bad bool
	tok := p.peek()

	if token.IsFuncParamTypes(tok.Kind) {
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

// parseFuncReturnTypes returns func return types
func (p *Parser) parseFuncReturnTypes() ast.ReturnTypes {
	if p.kind() == token.LBrace {
		return ast.ReturnTypes{}
	}

	var result ast.ReturnTypes
	if p.kind() == token.LParen {
		lp := p.expect(token.LParen, "expected '('")
		if p.kind() == token.RParen {
			p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
			btyp := &ast.BadType{
				From:   lp,
				To:     p.peek(),
				Reason: "expected parameter(s) before ')'",
			}
			result.List = append(result.List, ast.Param{Name: p.peek(), Type: btyp})
			return result
		}

		if p.kind() == token.Comma {
			p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
			btyp := &ast.BadType{
				From:   lp,
				To:     p.peek(),
				Reason: "expected expression before ','",
			}

			result.List = append(result.List, ast.Param{Name: p.peek(), Type: btyp})
			return result
		}

		result.LParen = lp
		// entering into kind: (indentA indentB, indentC indentD) or (indentA indentB)
		if p.kindNext(p.position+2) == token.Comma || p.kindNext(p.position+2) == token.RParen {
			for p.kind() != token.RParen && p.kind() != token.LBrace && p.kind() != token.EOF {
				param := p.parseFuncParam()
				_, bad := param.Type.(*ast.BadType)
				if bad {
					result.List = append(result.List, param)
					return result
				}
				result.List = append(result.List, param)

				if p.kind() != token.Comma && p.kind() != token.RParen {
					p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
					btyp := &ast.BadType{
						From:   lp,
						To:     p.peek(),
						Reason: "expected ',' after parameter(s)",
					}
					result.List = append(result.List, ast.Param{Name: p.peek(), Type: btyp})
					return result
				}

				if p.kind() == token.Comma {
					comma := p.expect(token.Comma, "expected ','")
					if p.kind() == token.RParen {
						p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
						btyp := &ast.BadType{
							From:   comma,
							To:     p.peek(),
							Reason: "expected parameter(s) after ','",
						}
						result.List = append(result.List, ast.Param{Name: p.peek(), Type: btyp})
						return result
					}
				}
			}
		} else {
			// entering into kind: (type, type)
			for p.kind() != token.RParen && p.kind() != token.LBrace && p.kind() != token.EOF {
				typ, btyp, bad := p.parseFuncParamType()
				if bad {
					result.List = append(result.List, ast.Param{Type: btyp})
					return result
				}
				result.List = append(result.List, ast.Param{Type: typ})

				if p.kind() != token.Comma && p.kind() != token.RParen {
					p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
					btyp := &ast.BadType{
						From:   lp,
						To:     p.peek(),
						Reason: "expected ',' after parameter(s)",
					}
					result.List = append(result.List, ast.Param{Name: p.peek(), Type: btyp})
					return result
				}

				if p.kind() == token.Comma {
					comma := p.expect(token.Comma, "expected ','")
					if p.kind() == token.RParen {
						p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
						btyp := &ast.BadType{
							From:   comma,
							To:     p.peek(),
							Reason: "expected parameter(s) after ','",
						}
						result.List = append(result.List, ast.Param{Name: p.peek(), Type: btyp})
						return result
					}
				}
			}
		}

		rp := p.expect(token.RParen, "expected ')'")
		result.RParen = rp
		return result
	}

	if p.kind() == token.Comma {
		p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
		btyp := &ast.BadType{
			From:   p.peek(),
			Reason: "expected builtin type or ident",
		}
		result.List = append(result.List, ast.Param{Name: p.peek(), Type: btyp})
		return result
	}

	next := p.peek()
	typ, btyp, bad := p.parseFuncParamType()
	if bad {
		p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", btyp.From.Line, btyp.From.Column, btyp.From.Kind, btyp.From.Value))
		result.List = append(result.List, ast.Param{Type: btyp})
		return result
	}
	result.List = append(result.List, ast.Param{Type: typ})

	if p.kind() != token.LBrace {
		p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
		btyp := &ast.BadType{
			From:   next,
			To:     p.peek(),
			Reason: "expected '{' after type",
		}
		result.List = append(result.List, ast.Param{Name: p.peek(), Type: btyp})
		return result
	}

	return result
}
