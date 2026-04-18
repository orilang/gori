package parser

import (
	"fmt"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
)

// parseInterfaceType returns parsed interface
func (p *Parser) parseInterfaceDecl() ast.Decl {
	kwt := p.expect(token.KWType, "expected 'type'")
	kwi := p.expectValidIdent(token.Ident, true, "expected 'ident'")
	kws := p.expect(token.KWInterface, "expected 'interface'")
	lbrace := p.expect(token.LBrace, "expected '{'")

	it := &ast.InterfaceDecl{
		TypeDecl:  kwt,
		Name:      kwi,
		Public:    isPublic(kwi),
		Interface: kws,
		LBrace:    lbrace,
	}

	for p.kind() != token.RBrace && p.kind() != token.EOF {
		if p.kind() == token.Ident {
			if p.kindNext(p.position+1) == token.LParen {
				it.Methods = append(it.Methods, p.parseFuncSignature())
			} else {
				it.Embeds = append(it.Embeds, p.parseInterfaceTypeEmbbed())
			}
		}

		if p.kind() == token.Comment {
			_ = p.next()
		}

		if p.kind() == token.SemiComma {
			_ = p.next()
			continue
		}

		if p.kind() == token.RBrace {
			break
		}

		if p.newlineSincePrev() {
			continue
		}

		p.Errors = append(p.Errors, fmt.Errorf("%d:%d: expected ';' or newline after interface field, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))

		p.consumeTo(token.RBrace)
	}

	rbrace := p.expect(token.RBrace, "expected '}'")
	it.RBrace = rbrace

	return it
}

// parseFuncSignature returns function signature for interface
func (p *Parser) parseFuncSignature() ast.InterfaceMethod {
	name := p.expectValidIdent(token.Ident, true, "expected function name")
	_ = p.expect(token.LParen, "expected '(' after function name")

	f := ast.InterfaceMethod{
		Name: name,
	}
	for p.kind() != token.RParen && p.kind() != token.EOF {
		if p.kind() == token.Comma {
			tok := p.next()
			p.Errors = append(p.Errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
			p.consumeTo(token.Comma)
			return f
		}
		f.Params = append(f.Params, p.parseFuncSignatureParam())
		if p.kind() != token.Comma && p.kind() != token.RParen && p.kind() != token.EOF {
			tok := p.next()
			p.Errors = append(p.Errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
			p.consumeTo(token.Comma)
			return f
		}

		if p.kind() == token.Comma {
			_ = p.next()
			if p.kind() == token.RParen || p.kind() == token.EOF {
				p.Errors = append(p.Errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
				p.consumeTo(token.Comma)
				return f
			}
		}
	}

	_ = p.expect(token.RParen, "expected ')' after function name")
	f.Results = p.parseFuncSignatureReturnTypes()

	return f
}

// parseFuncSignatureParam returns function parameter
func (p *Parser) parseFuncSignatureParam() ast.Param {
	name := p.expectValidIdent(token.Ident, true, "expected parameter identifier")
	return ast.Param{Name: name, Type: p.parseFuncSignatureParamType()}
}

// parseFuncSignatureParamType returns func parameter type
func (p *Parser) parseFuncSignatureParamType() ast.Type {
	typ := &ast.NamedType{}
	tok := p.next()

	if token.IsFuncParamTypes(tok.Kind) {
		typ.Parts = append(typ.Parts, tok)
	} else {
		p.Errors = append(p.Errors, fmt.Errorf("%d:%d: unsupported type with %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
		p.consumeTo(token.RParen)
	}
	return typ
}

// parseFuncSignatureReturnTypes returns func return types
func (p *Parser) parseFuncSignatureReturnTypes() ast.ReturnTypes {
	var result ast.ReturnTypes
	if p.kind() == token.LParen {
		lp := p.expect(token.LParen, "expected '('")
		if p.kind() == token.RParen {
			p.Errors = append(p.Errors, fmt.Errorf("%d:%d: expected parameter(s) before ')', got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
			result.List = append(result.List, ast.Param{Type: p.parseFuncSignatureParamType()})
			return result
		}

		if p.kind() == token.Comma {
			p.Errors = append(p.Errors, fmt.Errorf("%d:%d: expected expression before ',', got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))

			result.List = append(result.List, ast.Param{Type: p.parseFuncSignatureParamType()})
			return result
		}

		result.LParen = lp
		// entering into kind: (indentA indentB, indentC indentD) or (indentA indentB)
		if p.kindNext(p.position+2) == token.Comma || p.kindNext(p.position+2) == token.RParen {
			for p.kind() != token.RParen && p.kind() != token.LBrace && p.kind() != token.EOF {
				result.List = append(result.List, p.parseFuncSignatureParam())

				if p.kind() != token.Comma && p.kind() != token.RParen {
					p.Errors = append(p.Errors, fmt.Errorf("%d:%d: expected ',' after parameter(s), got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
					return result
				}

				if p.kind() == token.Comma {
					_ = p.expect(token.Comma, "expected ','")
					if p.kind() == token.RParen {
						p.Errors = append(p.Errors, fmt.Errorf("%d:%d: expected parameter(s) after ',', got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
						return result
					}
				}
			}
		} else {
			// entering into kind: (type, type)
			for p.kind() != token.RParen && p.kind() != token.LBrace && p.kind() != token.EOF {
				result.List = append(result.List, ast.Param{Type: p.parseFuncSignatureParamType()})

				if p.kind() != token.Comma && p.kind() != token.RParen {
					p.Errors = append(p.Errors, fmt.Errorf("%d:%d: expected ',' after parameter(s), got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
					return result
				}

				if p.kind() == token.Comma {
					_ = p.expect(token.Comma, "expected ','")
					if p.kind() == token.RParen {
						p.Errors = append(p.Errors, fmt.Errorf("%d:%d: expected parameter(s) after ',', got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
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
		p.Errors = append(p.Errors, fmt.Errorf("%d:%d: expected builtin type or ident, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
		return result
	}

	result.List = append(result.List, ast.Param{Type: p.parseFuncSignatureParamType()})
	return result
}

// parseInterfaceTypeEmbbed returns embbed signature for interface
func (p *Parser) parseInterfaceTypeEmbbed() ast.Type {
	kw := p.expectValidIdent(token.Ident, true, "expected 'ident'")

	nt := &ast.NamedType{}
	nt.Parts = append(nt.Parts, kw)
	if p.kind() == token.Dot {
		nt.Parts = append(nt.Parts, p.next())
		if p.kind() != token.Ident {
			p.Errors = append(p.Errors, fmt.Errorf("%d:%d: expected ident after '.', got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
			p.consumeTo(token.RBrace)
			return nt
		}
		nt.Parts = append(nt.Parts, p.next())
	}
	return nt
}
