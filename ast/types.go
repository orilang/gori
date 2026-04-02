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
	// Const      []Stmt
	Decls []Decl
	// Structs    []*StructType
	// Interfaces []*InterfaceType
	// Implements []*ImplementsDecl
	// Enums      []*EnumType
	// Sums       []*SumType
	// Comptime   []Stmt
}

// FuncDecl holds function parsed content
type FuncDecl struct {
	FuncKW  token.Token
	Name    token.Token
	Params  []Param
	Results ReturnTypes
	Body    *BlockStmt
}

// Params holds func parameter
type Param struct {
	Name token.Token
	Type Type
}

type ReturnTypes struct {
	LParen token.Token
	List   []Param
	RParen token.Token
}

// BlockStmt holds content between curly braces
type BlockStmt struct {
	LBrace token.Token
	Stmts  []Stmt
	RBrace token.Token
}

// ConstDecl holds constant content
type ConstDecl struct {
	ConstKW token.Token // KWConst
	Name    token.Token // Ident
	Type    Type        // Optional
	Eq      token.Token // Assign or Define
	Init    Expr
}

// VarDec holds variable content
type VarDecl struct {
	VarKW token.Token // KWVar
	Name  token.Token // Ident
	View  token.Token // Optional
	Type  Type
	Eq    token.Token // Assign or Define
	Init  Expr
}

// IdentExpr holds identifier content
type IdentExpr struct {
	Name token.Token // Ident
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

type DeclStmt struct {
	Decl Decl
}

type Decl interface {
	Position
	declNode()
}

type Type interface {
	Position
	typeNode()
}

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
	ForKW     token.Token
	Init      Stmt
	Condition Expr
	Post      Stmt
	Body      *BlockStmt
}

type RangeStmt struct {
	ForKW token.Token
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

type SwitchStmt struct {
	Switch token.Token
	Init   Stmt
	Tag    Expr
	LBrace token.Token
	Cases  []CaseClause
	RBrace token.Token
}

type CaseClause struct {
	Case   token.Token
	Values []Expr
	Colon  token.Token
	Body   []Stmt
}

type FallThroughStmt struct {
	FallThrough token.Token
}

type StructDecl struct {
	TypeDecl token.Token
	Name     token.Token
	Struct   token.Token
	Public   bool
	LBrace   token.Token
	Fields   []*FieldDecl
	RBrace   token.Token
}

type FieldDecl struct {
	Name    token.Token
	Public  bool
	Type    Type
	Eq      *token.Token // nil if no default
	Default Expr         // nil if no default
}

type NamedType struct {
	// e.g "pkg.Type" or "Type"
	Parts []token.Token // identifiers around dots
}

type InterfaceDecl struct {
	TypeDecl  token.Token
	Name      token.Token
	Public    bool
	Interface token.Token
	LBrace    token.Token
	Embeds    []Type
	Methods   []InterfaceMethod
	RBrace    token.Token
}

type InterfaceMethod struct {
	Name    token.Token
	Params  []Param
	Results ReturnTypes
}

type ImplementsDecl struct {
	TypeName   token.Token
	Implements token.Token
	Interface  Type
}

type EnumDecl struct {
	TypeDecl token.Token
	Name     token.Token
	Public   bool
	Enum     token.Token
	LBrace   token.Token
	Variants []token.Token
	RBrace   token.Token
}

type SumDecl struct {
	TypeDecl token.Token
	Name     token.Token
	Public   bool
	Sum      token.Token
	LBrace   token.Token
	Variants []SumVariant
	RBrace   token.Token
}

type SumVariant struct {
	Name   token.Token
	Params []Param
}

type SliceType struct {
	LBracket token.Token
	RBracket token.Token
	Elem     Type
}

type ArrayType struct {
	LBracket token.Token
	Len      Expr
	RBracket token.Token
	Elem     Type
}

type SliceLitExpr struct {
	Type     Type
	LBrace   token.Token
	Elements []Expr
	RBrace   token.Token
}

type SliceExpr struct {
	X        Expr
	LBracket token.Token
	Low      Expr
	Colon    token.Token
	High     Expr
	RBracket token.Token
}

type ComptimeBlockDecl struct {
	ComptimeKW token.Token
	Decls      []Decl
}

type MapType struct {
	KindKW    token.Token // map or hashmap
	LBracket  token.Token
	KeyType   Type
	RBracket  token.Token
	ValueType Type
}

type MakeExpr struct {
	MakeKW token.Token
	LParen token.Token
	Type   Type
	Args   []Expr // optional
	RParen token.Token
}
