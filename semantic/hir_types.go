package semantic

import (
	"github.com/orilang/gori/token"
)

// High-Level Intermediate Representation (HIR) holds all HIR types.
// This is a safe (type checked) semantic representation of the AST
// that will be used by the Low Intermediate Representation (LIR)

type Program struct {
	Files []*File
}

type File struct {
	Package string
	Decls   []Decl
}

type Decl interface {
	declNode()
}

type Expr interface {
	exprNode()
}

type Stmt interface {
	stmtNode()
}

type ResolvedSymbol struct {
	Name string
	Kind SymbolKind
	Type Type
}

type ConstDecl struct {
	Name   string
	Symbol ResolvedSymbol
	Eq     token.Kind
	Init   Expr
}

type VarDecl struct {
	Name   string
	Symbol ResolvedSymbol
	Eq     token.Kind
	Init   Expr
}

type StructDecl struct {
	Name   string
	Symbol ResolvedSymbol
	Fields []*FieldDecl
}

type FieldDecl struct {
	Name   string
	Symbol ResolvedSymbol
	// Type    Type
	Eq      *string // nil if no default
	Default Expr    // nil if no default
}

type InterfaceDecl struct {
	Name    string
	Symbol  ResolvedSymbol
	Embeds  []Type
	Methods []InterfaceMethod
}

type InterfaceMethod struct {
	Name    string
	Symbol  ResolvedSymbol
	Params  []Param
	Results []Param
}

type ImplementsDecl struct {
	Implements string
	Symbol     ResolvedSymbol
	Interface  Type
}

type EnumDecl struct {
	Name     string
	Symbol   ResolvedSymbol
	Variants []string
}

type SumDecl struct {
	Name     string
	Symbol   ResolvedSymbol
	Variants []SumVariant
}

type FuncDecl struct {
	Name     string
	Symbol   ResolvedSymbol
	Receiver *Receiver
	Params   []Param
	Results  []Param
	Body     *BlockStmt
}

type Receiver struct {
	Name   string
	Symbol ResolvedSymbol
	Shared bool
	Type   Type
}

type BlockStmt struct {
	Stmts []Stmt
}

type DeclStmt struct {
	Decl Decl
}

type ExprStmt struct {
	Expr Expr
}

type ReturnStmt struct {
	Values []Expr
}

type AssigmentStmt struct {
	Symbol ResolvedSymbol
	Right  Expr
}

type IfStmt struct {
	Condition Expr
	Then      []Stmt
	Else      []Stmt
}

type ConversionExpr struct {
	To    Type
	Value Expr
}

type IntLitExpr struct {
	Type  Type
	Value string
}

type FloatLitExpr struct {
	Type  Type
	Value string
}

type BoolLitExpr struct {
	Type  Type
	Value string
}

type StringLitExpr struct {
	Type  Type
	Value string
}

type IdentExpr struct {
	Type   Type
	Symbol ResolvedSymbol
	Value  string
}

type UnaryExpr struct {
	Type     Type
	Operator token.Kind
	Right    Expr
}

type BinaryExpr struct {
	Type     Type
	Left     Expr
	Operator token.Kind
	Right    Expr
}

type CallExpr struct {
	Callee     Expr
	CalleeType Type
	Args       []Expr
}
