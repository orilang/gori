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

// kindNext returns only the next token kind
func (p *Parser) kindNext(pos int) token.Kind {
	return p.peekNext(pos).Kind
}

// peekNext returns only the next token without
// advancing its position. Used as lookahead.
func (p *Parser) peekNext(pos int) token.Token {
	if pos >= len(p.Tokens) {
		return token.Token{Kind: token.EOF}
	}
	return p.Tokens[pos]
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

// consumeTo consumes all tokens until reaching
// the one passed or EOF
func (p *Parser) consumeTo(k token.Kind) {
	for p.kind() != k && p.kind() != token.EOF {
		_ = p.next()
	}
}

// newlineSincePrev is a boolean validating if we change line
func (p *Parser) newlineSincePrev() bool {
	if p.peek().Line == 0 {
		return false
	}
	return p.peek().Line > p.Tokens[p.position-1].Line
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

		case token.KWType:
			if token.IsValidTypeDecl(p.kindNext(p.position + 2)) {
				if p.kindNext(p.position+2) == token.KWStruct {
					f.Structs = append(f.Structs, p.parseStructType())
				} else if p.kindNext(p.position+2) == token.KWInterface {
					f.Interfaces = append(f.Interfaces, p.parseInterfaceType())
				} else if p.kindNext(p.position+2) == token.KWEnum {
					f.Enums = append(f.Enums, p.parseEnumDecl())
				} else {
					f.Sums = append(f.Sums, p.parseSumDecl())
				}
			} else {
				tok := p.peek()
				p.errors = append(p.errors, fmt.Errorf("%d:%d: unsupported file statement starting with %d %q", tok.Line, tok.Column, tok.Kind, tok.Value))
				p.consumeTo(token.RBrace)
			}

		default:
			if p.kindNext(p.position+1) == token.KWImplements {
				f.Implements = append(f.Implements, p.parseImplementsDecl())
			} else {
				tok := p.peek()
				p.errors = append(p.errors, fmt.Errorf("%d:%d: unsupported file statement starting with %d %q", tok.Line, tok.Column, tok.Kind, tok.Value))
				_ = p.next()
			}
		}
	}

	return f
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

	if p.kind() == token.KWType && token.IsValidTypeDecl(p.kindNext(p.position+2)) {
		{
			if p.kindNext(p.position+2) == token.KWStruct {
				return p.parseStructType()
			} else if p.kindNext(p.position+2) == token.KWInterface {
				return p.parseInterfaceType()
			} else if p.kindNext(p.position+2) == token.KWEnum {
				return p.parseEnumDecl()
			} else {
				return p.parseSumDecl()
			}
		}
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

// isPublic returns if the field is public or not
func isPublic(z token.Token) bool {
	ch := z.Value[0]
	return ch >= 'A' && ch <= 'Z'
}
