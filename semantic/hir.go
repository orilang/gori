package semantic

func (*ConstDecl) declNode()      {}
func (*VarDecl) declNode()        {}
func (*FuncDecl) declNode()       {}
func (*StructDecl) declNode()     {}
func (*InterfaceDecl) declNode()  {}
func (*EnumDecl) declNode()       {}
func (*SumDecl) declNode()        {}
func (*ImplementsDecl) declNode() {}

func (*ConversionExpr) exprNode() {}
func (*IdentExpr) exprNode()      {}
func (*IntLitExpr) exprNode()     {}
func (*FloatLitExpr) exprNode()   {}
func (*BoolLitExpr) exprNode()    {}
func (*StringLitExpr) exprNode()  {}
func (*UnaryExpr) exprNode()      {}
func (*BinaryExpr) exprNode()     {}
func (*CallExpr) exprNode()       {}

func (*DeclStmt) stmtNode()            {}
func (*ReturnStmt) stmtNode()          {}
func (*DefineAssigmentStmt) stmtNode() {}
func (*BlockStmt) stmtNode()           {}
