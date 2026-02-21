package parser

import (
	"fmt"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
)

// parseSumDecl returns parsed sum
func (p *Parser) parseSumDecl() *ast.SumType {
	kwt := p.expect(token.KWType, "expected 'type'")
	kwi := p.expect(token.Ident, "expected 'ident'")
	kws := p.expect(token.KWSum, "expected 'sum'")
	lbrace := p.expect(token.LBrace, "expected '{'")

	st := &ast.SumType{
		TypeDecl: kwt,
		Name:     kwi,
		Public:   isPublic(kwi),
		Sum:      kws,
		LBrace:   lbrace,
	}

	for p.kind() != token.RBrace && p.kind() != token.EOF {
		if p.kind() != token.Ident {
			p.errors = append(p.errors, fmt.Errorf("%d:%d: expected 'ident', got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
			p.consumeTo(token.RBrace)
			return st
		}

		if p.kind() == token.Ident {
			if p.kindNext(p.position+1) == token.LParen {
				st.Methods = append(st.Methods, p.parseSumFuncSignature())
			} else {
				st.Variants = append(st.Variants, p.next())
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

		p.errors = append(p.errors, fmt.Errorf("%d:%d: expected ';' or newline after sum field, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))

		p.consumeTo(token.RBrace)
	}

	if len(st.Variants) == 0 && len(st.Methods) == 0 {
		p.errors = append(p.errors, fmt.Errorf("%d:%d: expected variant(s) or variant method(s) inside braces, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
		p.consumeTo(token.RBrace)
		return st
	}

	rbrace := p.expect(token.RBrace, "expected '}'")
	st.RBrace = rbrace

	return st
}

// parseFuncSignature returns function signature for interface
func (p *Parser) parseSumFuncSignature() ast.VariantsMethods {
	name := p.expect(token.Ident, "expected function name")
	_ = p.expect(token.LParen, "expected '(' after function name")

	f := ast.VariantsMethods{
		Name: name,
	}
	for p.kind() != token.RParen && p.kind() != token.EOF {
		if p.kind() == token.Comma {
			tok := p.next()
			p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
			p.consumeTo(token.Comma)
			return f
		}
		f.Params = append(f.Params, p.parseSumFuncSignatureParam())
		if p.kind() != token.Comma && p.kind() != token.RParen && p.kind() != token.EOF {
			tok := p.next()
			p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
			p.consumeTo(token.Comma)
			return f
		}

		if p.kind() == token.Comma {
			_ = p.next()
			if p.kind() == token.RParen || p.kind() == token.EOF {
				p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
				p.consumeTo(token.Comma)
				return f
			}
		}
	}

	if len(f.Params) == 0 {
		p.errors = append(p.errors, fmt.Errorf("%d:%d: expected param(s) inside parenthesis, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
		p.consumeTo(token.RBrace)
		return f
	}

	_ = p.expect(token.RParen, "expected ')' after function name")
	return f
}

// parseFuncSignatureParam returns function parameter
func (p *Parser) parseSumFuncSignatureParam() ast.Param {
	name := p.expect(token.Ident, "expected parameter identifier")
	return ast.Param{Name: name, Type: p.parseSumFuncSignatureParamType()}
}

// parseFuncSignatureParamType returns func parameter type
func (p *Parser) parseSumFuncSignatureParamType() *ast.NameType {
	typ := &ast.NameType{}
	tok := p.next()

	if token.IsFuncParamTypes(tok.Kind) {
		typ.Name = tok
	} else {
		p.errors = append(p.errors, fmt.Errorf("%d:%d: unsupported type with %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
		p.consumeTo(token.RParen)
	}
	return typ
}
