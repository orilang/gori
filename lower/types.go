package lower

import (
	"bytes"

	"github.com/orilang/gori/commons"
	"github.com/orilang/gori/ir"
)

type Lower struct {
	// output when set to true outputs the HIR
	output bool

	files *commons.Files

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
