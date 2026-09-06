// ir holds intermediate representation (IR) of the AST
package ir

func (*Binary) instrNode()    {}
func (*Return) instrNode()    {}
func (*Const) instrNode()     {}
func (*Call) instrNode()      {}
func (*Branch) instrNode()    {}
func (*Label) instrNode()     {}
func (*Jump) instrNode()      {}
func (*Assigment) instrNode() {}
