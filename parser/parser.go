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

// lookForInForHeader loops against for statements to find
// provided token Kind and returns true when found
func (p *Parser) lookForInForHeader(k token.Kind) bool {
	pos := p.position
	if pos >= len(p.Tokens) {
		return false
	}
	for p.Tokens[pos].Kind != token.LBrace && p.Tokens[pos].Kind != token.EOF {
		if p.Tokens[pos].Kind == k {
			return true
		}
		pos++
	}

	return false
}

// lookForInSwitchHeader loops against for statements to find
// provided token Kind and returns true when found
func (p *Parser) lookForInSwitchHeader(k token.Kind) bool {
	pos := p.position
	if pos >= len(p.Tokens) {
		return false
	}
	for p.Tokens[pos].Kind != token.LBrace && p.Tokens[pos].Kind != token.EOF {
		if p.Tokens[pos].Kind == k {
			return true
		}
		pos++
	}

	return false
}

// lookForInSwitchCaseHeader loops against for statements to find
// provided token Kind and returns true when found
func (p *Parser) lookForInSwitchCaseHeader(k token.Kind) bool {
	pos := p.position
	if pos >= len(p.Tokens) {
		return false
	}
	for p.Tokens[pos].Kind != token.Colon && p.Tokens[pos].Kind != token.EOF {
		if p.Tokens[pos].Kind == k {
			return true
		}
		pos++
	}

	return false
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

	if p.kind() == token.KWFor {
		return p.parseForStmtExpr()
	}

	if p.kind() == token.KWBreak {
		return p.parseBreakStmt()
	}

	if p.kind() == token.KWContinue {
		return p.parseContinueStmt()
	}

	if p.kind() == token.KWSwitch {
		return p.parseSwitchStmt()
	}

	if p.kind() == token.KWFallThrough {
		return p.parseFallThroughStmt()
	}

	left := p.parseExpr(LOWEST)
	_, iok := left.(*ast.IdentExpr)
	_, sok := left.(*ast.SelectorExpr)
	_, xok := left.(*ast.IndexExpr)
	if iok || sok || xok {
		if token.IsAssignment(p.kind()) {
			return p.parseStmtExpr(left)
		}
		if token.IsIncDec(p.kind()) {
			return p.parseIncDecStmtExpr(left)
		}
	}
	if token.IsIncDec(p.kind()) {
		tok := p.next()
		p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected statement starting with %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
		return &ast.BadStmt{From: left.Start(), To: tok, Reason: "unexpected ++ or -- statement here"}
	}

	_, cok := left.(*ast.CallExpr)
	_, bok := left.(*ast.BadExpr)
	if !cok && !bok {
		p.errors = append(p.errors, fmt.Errorf("%d:%d: unsupported statement starting with %v %q", left.Start().Line, left.Start().Column, left.Start().Kind, left.Start().Value))
		return &ast.BadStmt{From: left.Start(), Reason: "unsupported statement"}
	}
	return &ast.ExprStmt{Expr: left}
}

// parseSimpleStmt returns declaration within parseBlock
func (p *Parser) parseSimpleStmt() ast.Stmt {
	left := p.parseExpr(LOWEST)
	_, iok := left.(*ast.IdentExpr)
	_, sok := left.(*ast.SelectorExpr)
	_, xok := left.(*ast.IndexExpr)
	if iok || sok || xok {
		if token.IsAssignment(p.kind()) {
			return p.parseStmtExpr(left)
		}
		if token.IsIncDec(p.kind()) {
			return p.parseIncDecStmtExpr(left)
		}
	}
	if token.IsIncDec(p.kind()) {
		tok := p.next()
		p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected statement starting with %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
		return &ast.BadStmt{From: left.Start(), To: tok, Reason: "unexpected ++ or -- statement here"}
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

func (p *Parser) parseCallExpr(left ast.Expr) ast.Expr {
	lb := p.expect(token.LParen, "expected '('")

	if p.kind() == token.RParen {
		rb := p.expect(token.RParen, "expected ')'")
		return &ast.CallExpr{
			Callee: left,
			LParen: lb,
			RParen: rb,
		}
	}

	var args []ast.Expr
	for p.kind() != token.RParen && p.kind() != token.EOF {
		if p.kind() == token.Comma {
			tok := p.next()
			p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
			return &ast.BadExpr{From: lb, To: tok, Reason: "expected expression not ','"}
		}

		args = append(args, p.parseExpr(LOWEST))
		if p.kind() != token.Comma && p.kind() != token.RParen && p.kind() != token.EOF {
			tok := p.next()
			p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
			return &ast.BadExpr{From: lb, To: tok, Reason: "expected ',' or ')'"}
		}

		if p.kind() == token.Comma {
			_ = p.next()
			if p.kind() == token.RParen || p.kind() == token.EOF {
				p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
				return &ast.BadExpr{From: lb, To: p.peek(), Reason: "expected expression after ','"}
			}
		}
	}

	rb := p.expect(token.RParen, "expected ')'")
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

// parseIncDecStmtExpr returns expressions for parseStmt func
func (p *Parser) parseIncDecStmtExpr(left ast.Expr) *ast.IncDecStmt {
	op := p.peek()
	_ = p.next()
	return &ast.IncDecStmt{
		X:        left,
		Operator: op,
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

// parseIfStmtExpr returns expressions for parseStmt func
func (p *Parser) parseIfStmtExpr() ast.Stmt {
	ifs := p.expect(token.KWIf, "expected 'if'")
	stmt := &ast.IfStmt{
		If: ifs,
	}

	if p.kind() == token.LBrace {
		p.errors = append(p.errors, fmt.Errorf("%d:%d: missing condition, got %v %q", p.peek().Line, p.peek().Column, p.peek().Kind, p.peek().Value))
		return &ast.BadStmt{From: ifs, To: p.peek(), Reason: "missing condition after 'if'"}
	}

	stmt.Condition = p.parseExpr(LOWEST)
	if token.IsAssignment(p.kind()) {
		tok := p.next()
		p.errors = append(p.errors, fmt.Errorf("%d:%d: assignment not allowed in if condition, got %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
		return &ast.BadStmt{From: ifs, To: tok, Reason: "assignment not allowed in if condition; use =="}
	}

	stmt.Then = p.parseBlock()

	if p.kind() == token.KWElse {
		_ = p.next()
		if p.kind() == token.KWIf {
			stmt.Else = p.parseIfStmtExpr()
		} else if p.kind() == token.LBrace {
			stmt.Else = p.parseBlock()
		} else {
			tok := p.next()
			p.errors = append(p.errors, fmt.Errorf("%d:%d: unexpected expression, got %v %q", tok.Line, tok.Column, tok.Kind, tok.Value))
			return &ast.BadStmt{From: ifs, To: tok, Reason: "expected expression '{' or 'if' after 'else'"}
		}
	}

	return stmt
}

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
