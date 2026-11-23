package lower

import (
	"bytes"

	"github.com/orilang/gori/ir"
	"github.com/orilang/gori/semantic"
)

type Lower struct {
	// output when set to true outputs the HIR
	output bool

	types  []ir.Instruction
	funcs  []*ir.Func
	errors []Diagnostic

	tIndex       int
	labelIndex   int
	instructions []ir.Instruction
	blocks       []*ir.Block
	blockName    string
}

// Files holds all files to use for tokenization
type Files struct {
	// Files holds the list of files to parse
	Files []string

	// output when set to true outputs the AST
	output bool
}

type Diagnostics []Diagnostic
type Diagnostic struct {
	Err error
}

type dumper struct {
	w *bytes.Buffer
}

type exprKind int

const (
	exprKindInvalid exprKind = iota
	exprKindInterger
	exprKindArray
	exprKindSlice
	exprKindMap
	exprKindHashMap
	exprKindFuncMethod
	exprKindNamedType
)

type infos struct {
	values         []ir.Value
	flowTerminates bool
	isBlank        bool
	breaks         bool
	exprKind       exprKind
	exprType       semantic.Type
	// Used only when exprKind == exprKindNamedType
	exprUnderlyingKind exprKind
}

type forStmtKind int

const (
	forStmtInfinite forStmtKind = iota
	forStmtConditional
	forStmt3Clauses
)

type inputConfig struct {
	forStmt forStmtKind

	// continues must be filled in for loop
	// with continue statement
	continues *ir.Jump

	// breaks must be filled in for loop
	// with break statement
	breaks *ir.Jump

	// isForSwitch must be used when we are in a for loop
	// and found a switch. We need to connect to it with a Jump
	isForSwitch bool
}
