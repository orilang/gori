package ast

import "github.com/orilang/gori/token"

func (*FuncDecl) declNode()      {}
func (*dumpType) declNode()      {}
func (*BlockStmt) stmtNode()     {}
func (*ConstDeclStmt) stmtNode() {}
func (*VarDeclStmt) stmtNode()   {}
func (*IdentExpr) stmtNode()     {}
func (*dumpType) stmtNode()      {}
func (*NameType) typeNode()      {}
func (*BadType) typeNode()       {}
func (*dumpType) typeNode()      {}
func (*IdentExpr) exprNode()     {}
func (*IntLitExpr) exprNode()    {}
func (*FloatLitExpr) exprNode()  {}
func (*BoolLitExpr) exprNode()   {}
func (*StringLitExpr) exprNode() {}
func (*ParenExpr) exprNode()     {}
func (*BadExpr) exprNode()       {}
func (*BinaryExpr) exprNode()    {}
func (*UnaryExpr) exprNode()     {}
func (*SelectorExpr) exprNode()  {}
func (*IndexExpr) exprNode()     {}
func (*CallExpr) exprNode()      {}
func (*dumpType) exprNode()      {}
func (*AssignStmt) stmtNode()    {}
func (*ExprStmt) stmtNode()      {}
func (*BadStmt) stmtNode()       {}
func (*BadType) stmtNode()       {}
func (*ReturnStmt) stmtNode()    {}
func (*IfStmt) stmtNode()        {}

func (x *IdentExpr) Start() token.Token { return x.Name }
func (x *IdentExpr) End() token.Token   { return x.Name }

func (x *NameType) Start() token.Token { return x.Name }
func (x *NameType) End() token.Token   { return x.Name }

func (x *IntLitExpr) Start() token.Token { return x.Name }
func (x *IntLitExpr) End() token.Token   { return x.Name }

func (x *FloatLitExpr) Start() token.Token { return x.Name }
func (x *FloatLitExpr) End() token.Token   { return x.Name }

func (x *BoolLitExpr) Start() token.Token { return x.Name }
func (x *BoolLitExpr) End() token.Token   { return x.Name }

func (x *StringLitExpr) Start() token.Token { return x.Name }
func (x *StringLitExpr) End() token.Token   { return x.Name }

func (x *ParenExpr) Start() token.Token { return x.Left }
func (x *ParenExpr) End() token.Token   { return x.Right }

func (x *BinaryExpr) Start() token.Token { return x.Left.Start() }
func (x *BinaryExpr) End() token.Token   { return x.Right.End() }

func (x *UnaryExpr) Start() token.Token { return x.Operator }
func (x *UnaryExpr) End() token.Token   { return x.Right.End() }

func (x *SelectorExpr) Start() token.Token { return x.X.Start() }
func (x *SelectorExpr) End() token.Token   { return x.Selector }

func (x *IndexExpr) Start() token.Token { return x.X.Start() }
func (x *IndexExpr) End() token.Token   { return x.RBracket }

func (x *CallExpr) Start() token.Token { return x.Callee.Start() }
func (x *CallExpr) End() token.Token   { return x.RParen }

func (x *ExprStmt) Start() token.Token { return x.Expr.Start() }
func (x *ExprStmt) End() token.Token   { return x.Expr.End() }

func (x *BadExpr) Start() token.Token { return x.From }
func (x *BadExpr) End() token.Token   { return x.To }

func (x *dumpType) Start() token.Token { return token.Token{} }
func (x *dumpType) End() token.Token   { return token.Token{} }

func (x *ReturnStmt) Start() token.Token { return x.Return }
func (x *ReturnStmt) End() token.Token {
	if len(x.Values) > 0 {
		return x.Values[len(x.Values)-1].End()
	}
	return token.Token{}
}

func (x *BlockStmt) Start() token.Token { return x.LBrace }
func (x *BlockStmt) End() token.Token   { return x.RBrace }

func (x *ConstDeclStmt) Start() token.Token { return x.ConstKW }
func (x *ConstDeclStmt) End() token.Token   { return x.Init.End() }

func (x *VarDeclStmt) Start() token.Token { return x.VarKW }
func (x *VarDeclStmt) End() token.Token   { return x.Init.End() }

func (x *AssignStmt) Start() token.Token { return x.Left.Start() }
func (x *AssignStmt) End() token.Token   { return x.Right.End() }

func (x *BadStmt) Start() token.Token { return x.From }
func (x *BadStmt) End() token.Token   { return x.To }

func (x *BadType) Start() token.Token { return x.From }
func (x *BadType) End() token.Token   { return x.To }

func (x *IfStmt) Start() token.Token { return x.If }
func (x *IfStmt) End() token.Token {
	if x.Else != nil {
		return x.Else.End()
	}
	return token.Token{}
}
