package parser

import (
	"fmt"
	"os"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/lexer"
	"github.com/orilang/gori/token"
	"github.com/orilang/gori/walk"
)

// NewParser returns files config to StartParsing
func NewParser(config Config) (*Files, error) {
	w, err := walk.Walk(walk.Config{File: config.File, Directory: config.Directory})
	if err != nil {
		return nil, err
	}

	return &Files{
		Files:  w.Files,
		output: config.Output,
	}, nil
}

// StartParsing ranges over files to return the AST
func (f *Files) StartParsing() error {
	for _, file := range f.Files {
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		l := lexer.New(data)
		l.Tokenize()
		p := New(l.Tokens)

		if f.output {
			fmt.Printf("%s\n", ast.Dump(p.ParseFile()))
		}
	}
	return nil
}

func New(tokens []token.Token) *Parser {
	return &Parser{
		Tokens: tokens,
		size:   len(tokens),
	}
}

// peek returns only the current token
func (p *Parser) peek() token.Token {
	if p.position >= len(p.Tokens) {
		return token.Token{Kind: token.EOF}
	}
	return p.Tokens[p.position]
}

// kind returns only the current token kind
func (p *Parser) kind() token.Kind {
	return p.peek().Kind
}

// peekPrecedence returns the precedence level of the current token.Kind
func (p *Parser) peekPrecedence() int {
	if v, ok := precedence[p.kind()]; ok {
		return v
	}
	return LOWEST
}

// next returns the current token AND advance the position
func (p *Parser) next() token.Token {
	tok := p.peek()
	if tok.Kind != token.EOF {
		p.position++
	}
	return tok
}

// match compares the current token with the provided token Kind
func (p *Parser) match(k token.Kind) (token.Token, bool) {
	if p.kind() == k {
		return p.next(), true
	}
	return token.Token{}, false
}

// expect compares the current token with the provided token Kind.
// If not found, and error will be append to errors list.
// It will then  return the next token from the list
func (p *Parser) expect(k token.Kind, msg string) token.Token {
	tok := p.peek()
	if tok.Kind != k {
		p.errors = append(p.errors, fmt.Errorf("%d:%d %s (got %v %q)", tok.Line, tok.Column, msg, tok.Kind, tok.Value))
	}
	return p.next()
}

// ParseFile returns the content of the file being parsed
func (p *Parser) ParseFile() *ast.File {
	kw := p.expect(token.KWPackage, "expected 'package'")
	name := p.expect(token.Ident, "expected package name")
	f := &ast.File{
		PackageKW: kw,
		Name:      name,
	}

	for p.kind() != token.EOF {
		switch p.kind() {
		case token.KWConst:
			f.Const = append(f.Const, p.parseConstDecl())

		case token.KWFunc:
			f.Decls = append(f.Decls, p.parseFuncDecl())

		default:
			tok := p.peek()
			p.errors = append(p.errors, fmt.Errorf("%d:%d: unsupported statement starting with %d %q", tok.Line, tok.Column, tok.Kind, tok.Value))
			p.next()
		}
	}

	return f
}

// parseFuncDecl returns function declaration
func (p *Parser) parseFuncDecl() *ast.FuncDecl {
	kw := p.expect(token.KWFunc, "expected 'func'")
	name := p.expect(token.Ident, "expected function name")
	_ = p.expect(token.LParen, "expected '(' after function name")

	f := &ast.FuncDecl{
		FuncKW: kw,
		Name:   name,
	}
	for p.kind() != token.RParen && p.kind() != token.EOF {
		if p.kind() == token.Comma {
			p.next()
			continue
		}
		f.Params = append(f.Params, p.parseFuncParam())
	}

	_ = p.expect(token.RParen, "expected ')' after function name")

	body := p.parseBlock()
	f.Body = body

	return f
}

// parseFuncParam returns function parameter
func (p *Parser) parseFuncParam() ast.Param {
	name := p.expect(token.Ident, "expected parameter identifier")
	typ, btyp, bad := p.parseType()
	if bad {
		return ast.Param{Name: name, Type: btyp}
	}
	return ast.Param{Name: name, Type: typ}
}

// parseType returns parameter type
func (p *Parser) parseType() (*ast.NameType, *ast.BadType, bool) {
	typ := &ast.NameType{}
	btyp := &ast.BadType{}
	var bad bool
	tok := p.peek()
	switch {
	case token.IsBuiltinType(tok.Kind):
		typ.Name = tok

	default:
		p.errors = append(p.errors, fmt.Errorf("%d:%d: unsupported type with %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
		btyp.From = tok
		btyp.Reason = "unexpected type name"
		bad = true
	}

	p.next()
	return typ, btyp, bad
}

// parseBlock returns declaration within curly braces
func (p *Parser) parseBlock() *ast.BlockStmt {
	lb := p.expect(token.LBrace, "expected '{'")
	var stmts []ast.Stmt

	for p.kind() != token.RBrace {
		stmts = append(stmts, p.parseStmt())
	}
	rb := p.expect(token.RBrace, "expected '}'")

	return &ast.BlockStmt{
		LBrace: lb,
		Stmts:  stmts,
		RBrace: rb,
	}
}

// parseStmt returns declaration within parseBlock
func (p *Parser) parseStmt() ast.Stmt {
	switch p.kind() {
	case token.KWConst:
		return p.parseConstDecl()

	case token.KWVar:
		return p.parseVarDecl()

	default:
		tok := p.peek()
		p.errors = append(p.errors, fmt.Errorf("%d:%d: unsupported statement starting with %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
		p.next()
	}

	return &ast.BlockStmt{}
}

// parseConstDecl returns constant declaration
func (p *Parser) parseConstDecl() ast.Stmt {
	kw := p.expect(token.KWConst, "expected 'const'")
	name := p.expect(token.Ident, "expected constant name")

	typeTok := p.peek()
	p.next()
	typ := &ast.NameType{Name: typeTok}

	eq := p.expect(token.Assign, "expected '=")
	init := p.parseExpr(LOWEST)

	return &ast.ConstDeclStmt{
		ConstKW: kw,
		Name:    name,
		Type:    typ,
		Eq:      eq,
		Init:    init,
	}
}

// parseVarDecl returns variable declaration
func (p *Parser) parseVarDecl() ast.Stmt {
	kw := p.expect(token.KWVar, "expected 'var'")
	name := p.expect(token.Ident, "expected variable name")

	typeTok := p.peek()
	p.next()
	typ := &ast.NameType{Name: typeTok}

	eq := p.expect(token.Assign, "expected '=")
	init := p.parseExpr(LOWEST)

	return &ast.VarDeclStmt{
		VarKW: kw,
		Name:  name,
		Type:  typ,
		Eq:    eq,
		Init:  init,
	}
}

// parseExpr returns the type of declaration being parsed
func (p *Parser) parseExpr(minPrecedence int) ast.Expr {
	if !token.IsPrefix(p.kind()) && !token.IsInfix(p.kind()) {
		tok := p.peek()
		p.errors = append(p.errors, fmt.Errorf("%d:%d: expected expression, got %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
		return &ast.BadExpr{From: tok, Reason: "unexpected expression"}
	}

	var left ast.Expr
	if token.IsPrefix(p.kind()) {
		left = p.parsePrefix()
	}

	for p.kind() != token.EOF && token.IsInfix(p.kind()) && p.peekPrecedence() >= minPrecedence {
		left = p.parseInfix(left)
	}
	return left
}

// parsePrefix returns expressions for parseExpr func
func (p *Parser) parsePrefix() ast.Expr {
	var expr ast.Expr
	switch p.kind() {
	case token.IntLit:
		expr = &ast.IntLitExpr{Name: p.next()}

	case token.FloatLit:
		expr = &ast.FloatLitExpr{Name: p.next()}

	case token.BoolLit:
		expr = &ast.BoolLitExpr{Name: p.next()}

	case token.StringLit:
		expr = &ast.StringLitExpr{Name: p.next()}

	case token.Ident:
		expr = &ast.IdentExpr{Name: p.next()}

	case token.LParen:
		expr = p.parseGroupExpr()
	}

	return expr
}

// parseInfix returns expressions for parseExpr func
func (p *Parser) parseInfix(left ast.Expr) ast.Expr {
	expr := &ast.BinaryExpr{
		Left:     left,
		Operator: p.peek(),
	}

	precedence := p.peekPrecedence()
	_ = p.next()
	expr.Right = p.parseExpr(precedence)

	return expr
}

func (p *Parser) parseGroupExpr() *ast.ParenExpr {
	from := p.peek()
	g := &ast.ParenExpr{Left: from}
	p.next()

	for p.kind() != token.RParen && p.kind() != token.EOF {
		g.Inner = p.parseExpr(LOWEST)
	}

	to := p.peek()
	_ = p.expect(token.RParen, "expected ')")

	g.Right = to
	if g.Inner == nil {
		g.Inner = &ast.BadExpr{From: from, To: to, Reason: "expected expression inside parentheses"}
		return g
	}

	return g
}
