package parser

import (
	"fmt"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
)

// parseStructType returns parsed struct
func (p *Parser) parseStructType() *ast.StructType {
	kwt := p.expect(token.KWType, "expected 'type'")
	kwi := p.expect(token.Ident, "expected 'ident'")
	kws := p.expect(token.KWStruct, "expected 'struct'")
	lbrace := p.expect(token.LBrace, "expected '{'")

	st := &ast.StructType{
		TypeDecl: kwt,
		Name:     kwi,
		Public:   isPublic(kwi),
		Struct:   kws,
		LBrace:   lbrace,
	}

	for p.kind() != token.RBrace && p.kind() != token.EOF {
		st.Fields = append(st.Fields, p.parseStructTypeField())

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

		p.errors = append(p.errors, fmt.Errorf("%d:%d: expected ';' or newline after struct field, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))

		p.consumeTo(token.RBrace)
	}
	rbrace := p.expect(token.RBrace, "expected '}'")
	st.RBrace = rbrace

	return st
}

// parseStructTypeField is in charge of parsing struct field
func (p *Parser) parseStructTypeField() *ast.FieldDecl {
	kw := p.expect(token.Ident, "expected 'ident'")
	fd := &ast.FieldDecl{
		Name:   kw,
		Public: isPublic(kw),
	}
	tok := p.peek()

	if !token.IsStructFieldTypes(tok.Kind) {
		fd.Type = &ast.BadType{From: tok, Reason: "unexpected type name"}
		p.errors = append(p.errors, fmt.Errorf("%d:%d: unsupported type with %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
		p.consumeTo(token.EOF)
		return fd
	}
	fd.Type = &ast.NameType{Name: tok}
	_ = p.next()

	if p.kind() == token.Assign {
		kwa := p.expect(token.Assign, "expected '='")
		fd.Eq = &kwa
		fd.Default = p.parseExpr(LOWEST)
	}
	return fd
}
