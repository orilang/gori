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
			p.errors = append(p.errors, fmt.Errorf("%d:%d: unsupported file statement starting with %d %q", tok.Line, tok.Column, tok.Kind, tok.Value))
			_ = p.next()
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
			_ = p.next()
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

	case tok.Kind == token.Ident:
		typ.Name = tok

	default:
		p.errors = append(p.errors, fmt.Errorf("%d:%d: unsupported type with %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
		btyp.From = tok
		btyp.Reason = "unexpected type name"
		bad = true
	}

	_ = p.next()
	return typ, btyp, bad
}

// parseBlock returns declaration within curly braces
func (p *Parser) parseBlock() *ast.BlockStmt {
	lb := p.expect(token.LBrace, "expected '{'")
	var stmts []ast.Stmt

	for p.kind() != token.RBrace && p.kind() != token.EOF {
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
	if p.kind() == token.KWConst {
		return p.parseConstDecl()
	}

	if p.kind() == token.KWVar {
		return p.parseVarDecl()
	}

	if p.kind() == token.KWReturn {
		return p.parseReturnStmtExpr()
	}

	if p.kind() == token.KWIf {
		return p.parseIfStmtExpr()
	}

	left := p.parseExpr(LOWEST)
	_, iok := left.(*ast.IdentExpr)
	_, sok := left.(*ast.SelectorExpr)
	_, xok := left.(*ast.IndexExpr)
	if (iok || sok || xok) && token.IsAssignment(p.kind()) {
		return p.parseStmtExpr(left)
	}

	_, cok := left.(*ast.CallExpr)
	_, bok := left.(*ast.BadExpr)
	if !cok && !bok {
		p.errors = append(p.errors, fmt.Errorf("%d:%d: unsupported statement starting with %v %q", left.Start().Line, left.Start().Column, left.Start().Kind, left.Start().Value))
		return &ast.BadStmt{From: left.Start(), Reason: "unsupported statement"}
	}
	return &ast.ExprStmt{Expr: left}
}

// parseConstDecl returns constant declaration
func (p *Parser) parseConstDecl() ast.Stmt {
	kw := p.expect(token.KWConst, "expected 'const'")
	name := p.expect(token.Ident, "expected constant name")

	typ, btyp, bad := p.parseType()
	if bad {
		return btyp
	}
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

	typ, btyp, bad := p.parseType()
	if bad {
		return btyp
	}
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
	if !token.IsPrefix(p.kind()) {
		tok := p.next()
		p.errors = append(p.errors, fmt.Errorf("%d:%d: expected prefix expression, got %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
		return &ast.BadExpr{From: tok, Reason: "unexpected prefix expression"}
	}

	var left ast.Expr
	left = p.parsePrefix()

	for p.kind() != token.EOF && p.peekPrecedence() >= minPrecedence {
		if token.IsPostfix(p.kind()) {
			left = p.parsePostfix(left)
		} else if token.IsInfix(p.kind()) {
			left = p.parseInfix(left)
		} else {
			break
		}
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

	case token.Minus, token.Not:
		expr = p.parseUnaryExpr()
	}

	return expr
}

// parseInfix returns expressions for parseExpr func
func (p *Parser) parseInfix(left ast.Expr) ast.Expr {
	expr := &ast.BinaryExpr{
		Left:     left,
		Operator: p.peek(),
	}

	l, lok := expr.Left.(*ast.BinaryExpr)
	if lok && token.IsChainingComparison(l.Operator.Kind) && token.IsChainingComparison(expr.Operator.Kind) {
		p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected chaining comparison expression, got %v %q", expr.Operator.Line, expr.Operator.Column, expr.Operator.Kind, expr.Operator.Value))
		return &ast.BadExpr{From: l.Operator, To: expr.Operator, Reason: "unexpected chaining comparison expression, use && (e.g. a < b && b < c)"}
	}

	precedence := p.peekPrecedence()
	_ = p.next()
	expr.Right = p.parseExpr(precedence + 1)

	return expr
}

// parsePostfix returns expressions for parseExpr func
func (p *Parser) parsePostfix(left ast.Expr) ast.Expr {
	var expr ast.Expr
	switch p.kind() {
	case token.Dot:
		expr = p.parseSelectorExpr(left)

	case token.LBracket:
		expr = p.parseIndexExpr(left)

	case token.LParen:
		expr = p.parseCallExpr(left)
	}

	return expr
}

// parseGroupExpr parses expressions within parenthesis
func (p *Parser) parseGroupExpr() *ast.ParenExpr {
	from := p.expect(token.LParen, "expected '('")
	g := &ast.ParenExpr{Left: from}

	if p.kind() == token.RParen {
		to := p.expect(token.RParen, "expected ')'")
		g.Right = to
		p.errors = append(p.errors, fmt.Errorf("%d:%d: expected expression inside parentheses, got %v %q", to.Line, to.Column, to.Kind, to.Value))
		g.Inner = &ast.BadExpr{From: from, To: to, Reason: "expected expression inside parentheses"}
		return g
	}

	g.Inner = p.parseExpr(LOWEST)
	to := p.expect(token.RParen, "expected ')'")
	g.Right = to

	return g
}

// parseUnaryExpr parses unary expression like - and !
func (p *Parser) parseUnaryExpr() *ast.UnaryExpr {
	u := &ast.UnaryExpr{Operator: p.next()}
	u.Right = p.parseExpr(PREFIX + 1)
	return u
}

// parseSelector returns expressions for parsePostfix func
func (p *Parser) parseSelectorExpr(left ast.Expr) ast.Expr {
	dot := p.expect(token.Dot, "expected '.'")
	selector := p.expect(token.Ident, "expected 'ident'")
	return &ast.SelectorExpr{
		X:        left,
		Dot:      dot,
		Selector: selector,
	}
}

// parseIndexSelector returns expressions for parsePostfix func
func (p *Parser) parseIndexExpr(left ast.Expr) ast.Expr {
	lb := p.expect(token.LBracket, "expected '['")
	index := p.parseExpr(LOWEST)
	rb := p.expect(token.RBracket, "expected ']'")
	return &ast.IndexExpr{
		X:        left,
		LBracket: lb,
		Index:    index,
		RBracket: rb,
	}
}

// parseCallExpr returns expressions for parsePostfix func
func (p *Parser) parseCallExpr(left ast.Expr) ast.Expr {
	lb := p.expect(token.LParen, "expected '('")
	var args []ast.Expr
	for p.kind() != token.RParen && p.kind() != token.EOF {
		if p.kind() == token.Comma {
			_ = p.next()
		}
		args = append(args, p.parseExpr(LOWEST))
		if p.kind() != token.Comma && p.kind() != token.RParen && p.kind() != token.EOF {
			p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
			return &ast.BadExpr{From: lb, To: p.peek(), Reason: "expected ','"}
		}
	}
	rb := p.expect(token.RParen, "expected ')'")
	if rb.Kind != token.RParen {
		return &ast.BadExpr{From: lb, To: rb, Reason: "expected ')'"}
	}
	return &ast.CallExpr{
		Callee: left,
		LParen: lb,
		Args:   args,
		RParen: rb,
	}
}

// parseStmtExpr returns expressions for parseStmt func
func (p *Parser) parseStmtExpr(left ast.Expr) *ast.AssignStmt {
	op := p.peek()
	_ = p.next()
	return &ast.AssignStmt{
		Left:     left,
		Operator: op,
		Right:    p.parseExpr(LOWEST),
	}
}

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
			_ = p.next()
		}
		args = append(args, p.parseExpr(LOWEST))
		if p.kind() == token.EOF || p.kind() == token.RBrace {
			break
		}

		if p.kind() != token.Comma && p.kind() != token.RBrace && p.kind() != token.EOF {
			p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
			return &ast.BadStmt{From: rn, To: p.peek(), Reason: "expected expression after ','"}
		}
	}

	return &ast.ReturnStmt{
		Return: rn,
		Values: args,
	}
}

// parseIfStmtExpr returns expressions for parseStmt func
func (p *Parser) parseIfStmtExpr() ast.Stmt {
	ifs := p.expect(token.KWIf, "expected 'if'")
	stmt := &ast.IfStmt{
		If: ifs,
	}

	if p.peek().Kind == token.LBrace {
		p.errors = append(p.errors, fmt.Errorf("%d:%d: missing condition, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
		return &ast.BadStmt{From: ifs, To: p.peek(), Reason: "missing condition after 'if'"}
	}

	stmt.Condition = p.parseExpr(LOWEST)
	if token.IsAssignment(p.kind()) {
		p.errors = append(p.errors, fmt.Errorf("%d:%d: assignement not allowed in if condition, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
		return &ast.BadStmt{From: ifs, To: p.peek(), Reason: "assignement not allowed in if condition; use =="}
	}

	stmt.Then = p.parseBlock()

	if p.peek().Kind == token.KWElse {
		_ = p.next()
		if p.peek().Kind == token.KWIf {
			stmt.Else = p.parseIfStmtExpr()
		} else if p.kind() == token.LBrace {
			stmt.Else = p.parseBlock()
		} else {
			p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
			return &ast.BadStmt{From: ifs, To: p.peek(), Reason: "expected expression '{' or 'if' after 'else'"}
		}
	}

	return stmt
}
