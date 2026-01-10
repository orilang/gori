package ast

import (
	"bytes"

	"github.com/orilang/gori/token"
)

// dumper holds requirements to print AST content in
// human readable form
type dumper struct {
	w *bytes.Buffer
}

// File holds requirements from parsed file
type File struct {
	PackageKW token.Token
	Name      token.Token
	Const     []Stmt
	Decls     []Decl
}

// FuncDecl holds function parsed content
type FuncDecl struct {
	FuncKW token.Token
	Name   token.Token
	Params []Param
	Body   *BlockStmt
}

// Params holds func parameter
type Param struct {
	Name token.Token
	Type Type
}

// BlockStmt holds content between curly braces
type BlockStmt struct {
	LBrace token.Token
	Stmts  []Stmt
	RBrace token.Token
}

// ConstDeclStmt holds constant content
type ConstDeclStmt struct {
	ConstKW token.Token // KWConst
	Name    token.Token // Ident
	Type    Type        // Optional
	Eq      token.Token // Assign or Define
	Init    Expr
}

// VarDeclStmt holds variable content
type VarDeclStmt struct {
	VarKW token.Token // KWVar
	Name  token.Token // Ident
	Type  Type        // Optional
	Eq    token.Token // Assign or Define
	Init  Expr
}

// IdentExpr holds identifier content
type IdentExpr struct {
	Name token.Token // Ident
}

// NameType holds identifier type
type NameType struct {
	Name token.Token
}

// BadType holds returned bad type with reason
type BadType struct {
	From, To token.Token
	Reason   string
}

// BadExpr holds returned bad expr with reason
type BadExpr struct {
	From, To token.Token
	Reason   string
}

// dumpType is only used for testing purpose
type dumpType struct {
	S string
}

// IntLitExpr holds integer literal content
type IntLitExpr struct {
	Name token.Token
}

// FloatLitExpr holds integer literal content
type FloatLitExpr struct {
	Name token.Token
}

// BoolLitExpr holds bool content
type BoolLitExpr struct {
	Name token.Token
}

// StringLitExpr holds string content
type StringLitExpr struct {
	Name token.Token
}

// ParenExpr handles contents between open and closing parenthesis
type ParenExpr struct {
	Left  token.Token
	Inner Expr
	Right token.Token
}

type Decl interface{ declNode() }
type Type interface{ typeNode() }
type Stmt interface{ stmtNode() }
type Expr interface{ exprNode() }
