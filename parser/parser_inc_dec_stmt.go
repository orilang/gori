package parser

import "github.com/orilang/gori/ast"

// parseIncDecStmtExpr returns expressions for parseStmt func
func (p *Parser) parseIncDecStmtExpr(left ast.Expr) *ast.IncDecStmt {
	op := p.peek()
	_ = p.next()
	return &ast.IncDecStmt{
		X:        left,
		Operator: op,
	}
}
