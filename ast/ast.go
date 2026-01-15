package ast

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
func (*BadType) exprNode()       {}
func (*ParenExpr) exprNode()     {}
func (*BadExpr) exprNode()       {}
func (*BinaryExpr) exprNode()    {}
func (*UnaryExpr) exprNode()     {}
func (*SelectorExpr) exprNode()  {}
func (*IndexExpr) exprNode()     {}
func (*CallExpr) exprNode()      {}
func (*dumpType) exprNode()      {}
