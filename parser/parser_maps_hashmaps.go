package parser

import (
	"fmt"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
)

// parseMapsHashMapsDecl parses maps or hashmaps declarations
func (p *Parser) parseMapsHashMapsDecl() *ast.MapType {
	var kw token.Token
	if p.kind() == token.KWMap {
		kw = p.expect(token.KWMap, "expected 'map'")
	} else {
		kw = p.expect(token.KWHashMap, "expected 'hashmap'")
	}

	lb := p.expect(token.LBracket, "expected '['")
	x := &ast.MapType{
		KindKW:   kw,
		LBracket: lb,
	}

	var keyType ast.TypeRef
	for p.kind() != token.RBracket && p.kind() != token.RParen && p.kind() != token.EOF {
		if !token.IsMapTypes(p.kind()) {
			p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected map/hashmap key type, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
			p.consumeTo(token.EOF)
			return x
		}
		keyType.Parts = append(keyType.Parts, p.next())

		if p.kind() == token.RBracket || p.kind() == token.RParen {
			break
		}
	}

	x.KeyType = keyType
	x.RBracket = p.expect(token.RBracket, "expected ']'")

	var valueType ast.TypeRef
	for p.kind() != token.Assign && p.kind() != token.EOF {
		if !token.IsMapTypes(p.kind()) {
			p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected map/hashmap key type, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
			p.consumeTo(token.EOF)
			return x
		}
		valueType.Parts = append(valueType.Parts, p.next())

		if p.newlineSincePrev() || p.kind() == token.Assign || p.kind() == token.RParen || p.kind() == token.Comma {
			break
		}
	}

	x.ValueType = valueType

	return x
}
