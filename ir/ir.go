// ir holds intermediate representation (IR) of the AST
package ir

func (*Binary) instrNode()          {}
func (*Return) instrNode()          {}
func (*Const) instrNode()           {}
func (*Call) instrNode()            {}
func (*Branch) instrNode()          {}
func (*Label) instrNode()           {}
func (*Jump) instrNode()            {}
func (*Assigment) instrNode()       {}
func (*Extract) instrNode()         {}
func (*Unary) instrNode()           {}
func (*IncDec) instrNode()          {}
func (*Index) instrNode()           {}
func (*Len) instrNode()             {}
func (*MapIter) instrNode()         {}
func (*IterNext) instrNode()        {}
func (*IterKey) instrNode()         {}
func (*IterValue) instrNode()       {}
func (*DefinedTypeDecl) instrNode() {}
func (*StructDecl) instrNode()      {}
func (*FieldGet) instrNode()        {}
