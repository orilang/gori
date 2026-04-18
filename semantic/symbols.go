// Symbol collection is when you walk top-level declarations
// to register their names before checking their bodies
package semantic

import "github.com/orilang/gori/ast"

type SymbolKind int

const (
	SymInvalid SymbolKind = iota
	SymType
	SymFunc
	SymVar
	SymConst
)

type Symbol struct {
	Name string
	Kind SymbolKind
	Type Type
	Decl ast.Decl
}
