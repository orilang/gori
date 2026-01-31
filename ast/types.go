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

// BadStmt holds returned bad stmt with reason
type BadStmt struct {
	From, To token.Token
	Reason   string
}

// BadDecl holds returned bad declaration with reason
type BadDecl struct {
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

// BinaryExpr handles complex binary expressions like addition
// multiplication etc
type BinaryExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

// UnaryExpr handles expressions like - or !
type UnaryExpr struct {
	Operator token.Token
	Right    Expr
}

// SelectorExpr handles field access expressions like a.b, a.b.c etc
type SelectorExpr struct {
	X        Expr
	Dot      token.Token
	Selector token.Token
}

// IndexExpr handles index/slicing like a[b] etc
type IndexExpr struct {
	X        Expr
	LBracket token.Token
	Index    Expr
	RBracket token.Token
}

// CallExpr handles function call
type CallExpr struct {
	Callee Expr
	LParen token.Token
	Args   []Expr
	RParen token.Token
}

// AssignStmt handles assignement expressions
type AssignStmt struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

// ExprStmt is used by Stmt
type ExprStmt struct {
	Expr Expr
}

type Decl interface{ declNode() }
type Type interface{ typeNode() }
type Stmt interface {
	Position
	stmtNode()
}
type Expr interface {
	Position
	exprNode()
}

type Position interface {
	Start() token.Token
	End() token.Token
}

type ReturnStmt struct {
	Return token.Token
	Values []Expr
}

type IfStmt struct {
	If        token.Token
	Condition Expr
	Then      *BlockStmt
	Else      Stmt
}

type ForStmt struct {
	For       token.Token
	Init      Stmt
	Condition Expr
	Post      Stmt
	Body      *BlockStmt
}

type RangeStmt struct {
	For   token.Token
	Key   *IdentExpr
	Value *IdentExpr
	Op    token.Token
	Range token.Token
	X     Expr
	Body  *BlockStmt
}

type IncDecStmt struct {
	X        Expr        // must be assignable: Ident/Selector/Index
	Operator token.Token // ++ or --
}

type BreakStmt struct {
	Break token.Token
}

type ContinueStmt struct {
	Continue token.Token
}
