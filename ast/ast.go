package ast

import "github.com/orilang/gori/token"

func (*FuncDecl) declNode()          {}
func (*dumpType) declNode()          {}
func (*BadDecl) declNode()           {}
func (*ConstDecl) declNode()         {}
func (*VarDecl) declNode()           {}
func (*BadType) declNode()           {}
func (*StructDecl) declNode()        {}
func (*InterfaceDecl) declNode()     {}
func (*EnumDecl) declNode()          {}
func (*SumDecl) declNode()           {}
func (*ComptimeBlockDecl) declNode() {}
func (*ImplementsDecl) declNode()    {}
func (*DefinedTypeDecl) declNode()   {}
func (*NamedType) declNode()         {}

func (*StructDecl) typeDeclNode()      {}
func (*InterfaceDecl) typeDeclNode()   {}
func (*EnumDecl) typeDeclNode()        {}
func (*SumDecl) typeDeclNode()         {}
func (*DefinedTypeDecl) typeDeclNode() {}
func (*NamedType) typeDeclNode()       {}

func (*dumpType) typeNode()  {}
func (*BadType) typeNode()   {}
func (*NamedType) typeNode() {}
func (*SliceType) typeNode() {}
func (*ArrayType) typeNode() {}
func (*MapType) typeNode()   {}

func (*dumpType) exprNode()      {}
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
func (*SliceExpr) exprNode()     {}
func (*MakeExpr) exprNode()      {}
func (*SliceLitExpr) exprNode()  {}

func (*dumpType) stmtNode()        {}
func (*BlockStmt) stmtNode()       {}
func (*AssignStmt) stmtNode()      {}
func (*ExprStmt) stmtNode()        {}
func (*BadStmt) stmtNode()         {}
func (*ReturnStmt) stmtNode()      {}
func (*IfStmt) stmtNode()          {}
func (*ForStmt) stmtNode()         {}
func (*RangeStmt) stmtNode()       {}
func (*IncDecStmt) stmtNode()      {}
func (*BreakStmt) stmtNode()       {}
func (*ContinueStmt) stmtNode()    {}
func (*SwitchStmt) stmtNode()      {}
func (*FallThroughStmt) stmtNode() {}
func (*DeclStmt) stmtNode()        {}

func (x *FuncDecl) Start() token.Token {
	return x.FuncKW
}
func (x *FuncDecl) End() token.Token {
	if x.Body != nil {
		return x.Body.End()
	}
	return token.Token{}
}

func (x *IdentExpr) Start() token.Token { return x.Name }
func (x *IdentExpr) End() token.Token   { return x.Name }

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

func (x *ConstDecl) Start() token.Token { return x.ConstKW }
func (x *ConstDecl) End() token.Token   { return x.Init.End() }

func (x *VarDecl) Start() token.Token { return x.VarKW }
func (x *VarDecl) End() token.Token   { return x.Init.End() }

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

func (x *BadDecl) Start() token.Token { return x.From }
func (x *BadDecl) End() token.Token   { return x.To }

func (x *ForStmt) Start() token.Token { return x.ForKW }
func (x *ForStmt) End() token.Token {
	if x.Body != nil {
		return x.Body.End()
	}
	return token.Token{}
}

func (x *RangeStmt) Start() token.Token { return x.ForKW }
func (x *RangeStmt) End() token.Token {
	if x.Body != nil {
		return x.Body.End()
	}
	return token.Token{}
}

func (x *IncDecStmt) Start() token.Token { return x.X.Start() }
func (x *IncDecStmt) End() token.Token   { return x.Operator }

func (x *BreakStmt) Start() token.Token { return x.Break }
func (x *BreakStmt) End() token.Token   { return x.Break }

func (x *ContinueStmt) Start() token.Token { return x.Continue }
func (x *ContinueStmt) End() token.Token   { return x.Continue }

func (x *SwitchStmt) Start() token.Token { return x.Switch }
func (x *SwitchStmt) End() token.Token   { return x.RBrace }

func (x *CaseClause) Start() token.Token { return x.Case }
func (x *CaseClause) End() token.Token {
	if len(x.Body) > 0 {
		return x.Body[len(x.Body)-1].End()
	}
	return token.Token{}
}

func (x *FallThroughStmt) Start() token.Token { return x.FallThrough }
func (x *FallThroughStmt) End() token.Token   { return x.FallThrough }

func (x *StructDecl) Start() token.Token { return x.TypeDecl }
func (x *StructDecl) End() token.Token   { return x.RBrace }

func (x *InterfaceDecl) Start() token.Token { return x.TypeDecl }
func (x *InterfaceDecl) End() token.Token   { return x.RBrace }

func (x *EnumDecl) Start() token.Token { return x.TypeDecl }
func (x *EnumDecl) End() token.Token   { return x.RBrace }

func (x *SumDecl) Start() token.Token { return x.TypeDecl }
func (x *SumDecl) End() token.Token   { return x.RBrace }

func (x *SliceType) Start() token.Token { return x.LBracket }
func (x *SliceType) End() token.Token {
	if x.Elem != nil {
		return x.Elem.End()
	}
	return token.Token{}
}

func (x *ArrayType) Start() token.Token { return x.LBracket }
func (x *ArrayType) End() token.Token {
	if x.Elem != nil {
		return x.Elem.End()
	}
	return token.Token{}
}

func (x *SliceExpr) Start() token.Token { return x.X.Start() }
func (x *SliceExpr) End() token.Token   { return x.RBracket }

func (x *SliceLitExpr) Start() token.Token {
	if x.Type != nil {
		return x.Type.Start()
	}
	return token.Token{}
}
func (x *SliceLitExpr) End() token.Token { return x.RBrace }

func (x *NamedType) Start() token.Token {
	if len(x.Parts) > 0 {
		return x.Parts[0]
	}
	return token.Token{}
}
func (x *NamedType) End() token.Token {
	if len(x.Parts) > 0 {
		return x.Parts[len(x.Parts)-1]
	}
	return token.Token{}
}

func (x *ComptimeBlockDecl) Start() token.Token { return x.ComptimeKW }
func (x *ComptimeBlockDecl) End() token.Token   { return x.ComptimeKW }

func (x *MapType) Start() token.Token { return x.KindKW }
func (x *MapType) End() token.Token   { return x.ValueType.End() }

func (x *MakeExpr) Start() token.Token { return x.MakeKW }
func (x *MakeExpr) End() token.Token   { return x.RParen }

func (x *ImplementsDecl) Start() token.Token { return x.TypeName }
func (x *ImplementsDecl) End() token.Token {
	if x.Interface != nil {
		return x.Interface.End()
	}
	return token.Token{}
}

func (x *DeclStmt) Start() token.Token {
	if x.Decl != nil {
		return x.Decl.Start()
	}
	return token.Token{}
}

func (x *DeclStmt) End() token.Token {
	if x.Decl != nil {
		return x.Decl.End()
	}
	return token.Token{}
}

func (x *DefinedTypeDecl) Start() token.Token { return x.TypeDecl }
func (x *DefinedTypeDecl) End() token.Token {
	if x.Type != nil {
		return x.Type.End()
	}
	return token.Token{}
}
