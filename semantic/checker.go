package semantic

import (
	"fmt"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
)

// NewScope allows use to create new scope base on provided config.
// Scope must be used for
// - package level type names
// - function names
// - const/var names
// - function parameters
// - local variables
// - block scoped names
func NewScope(parent *Scope) *Scope {
	return &Scope{
		Parent:  parent,
		Symbols: make(map[string]*Symbol),
	}
}

// Declare declares new symbol and return true when NOT exists
func (s *Scope) Declare(sym *Symbol) bool {
	if _, exists := s.Symbols[sym.Name]; exists {
		return false
	}
	s.Symbols[sym.Name] = sym
	return true
}

// Lookup allows us to loop over parent scope in order to find the provided one and returns its related Symbol if exists
func (s *Scope) Lookup(name string) *Symbol {
	for scope := s; scope != nil; scope = scope.Parent {
		if sym := scope.Symbols[name]; sym != nil {
			return sym
		}
	}
	return nil
}

// LookupLocal returns the symbol from map
func (s *Scope) LookupLocal(name string) *Symbol {
	return s.Symbols[name]
}

// NewChecker returns new checker pointer with pkgScope initialized to nil.
// pkgScope is set to nil because it's the parent scope.
// All new scopes will only be childrens
func NewChecker() *Checker {
	return &Checker{
		pkgScope: NewScope(nil),
	}
}

// Check performs the type checking step in order to validate
// all code definition structure and fill diagnostics when
// errros are found
func (c *Checker) Check(file *ast.File) []Diagnostics {
	c.collectTopLevelSymbols(file)
	c.createTypeObjects()
	c.resolveTypeDecls()
	c.resolveFuncSignatures()
	return c.errors
}

// collectTopLevelSymbols collects top levels symbols names first
// in order to make sure that all definitions exists before creating
// semantic objects and resolve remaining contents.
// This prevents having types that does not exists.
func (c *Checker) collectTopLevelSymbols(file *ast.File) {
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case ast.TypeDecl:
			c.declareTypeSymbol(d)
		case *ast.FuncDecl:
			c.declareFuncSymbol(d)
		}
	}
}

// typeDeclName returns the name of the type declaration
func typeDeclName(decl ast.Decl) string {
	switch d := decl.(type) {
	case *ast.DefinedTypeDecl:
		return d.Name.Value
	case *ast.StructDecl:
		return d.Name.Value
	case *ast.InterfaceDecl:
		return d.Name.Value
	case *ast.EnumDecl:
		return d.Name.Value
	case *ast.SumDecl:
		return d.Name.Value
	default:
		return ""
	}
}

// declareTypeSymbol declares new type symbol with its name
// and append diagnostics errors when already exists
func (c *Checker) declareTypeSymbol(decl ast.TypeDecl) {
	name := typeDeclName(decl)
	if name == "" {
		return
	}

	sym := &Symbol{
		Name: name,
		Kind: SymType,
		Decl: decl,
	}

	if !c.pkgScope.Declare(sym) {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("symbol %q already declared", name)})
		return
	}

	c.typeDecls = append(c.typeDecls, decl)
}

// declareFuncSymbol declares new type symbol with its name
// and append diagnostics errors when already exists
func (c *Checker) declareFuncSymbol(fn *ast.FuncDecl) {
	sym := &Symbol{
		Name: fn.Name.Value,
		Kind: SymFunc,
		Decl: fn,
	}

	if !c.pkgScope.Declare(sym) {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("symbol %q already declared", fn.Name.Value)})
		return
	}

	c.funcDecls = append(c.funcDecls, fn)
}

// createTypeObjects create structured type objects
// after declareTypeSymbol func call by filling its Type field.
// Each type declaration creates only one semantic type object.
func (c *Checker) createTypeObjects() {
	for _, decl := range c.typeDecls {
		name := typeDeclName(decl)
		sym := c.pkgScope.Lookup(name)
		if sym == nil {
			continue
		}

		switch d := decl.(type) {
		case *ast.DefinedTypeDecl:
			sym.Type = &NamedType{
				Decl: d,
			}

		case *ast.StructDecl:
			sym.Type = &StructType{
				Decl: d,
			}

		case *ast.InterfaceDecl:
			sym.Type = &InterfaceType{
				Decl: d,
			}
		case
			*ast.EnumDecl:
			sym.Type = &EnumType{
				Decl: d,
			}

		case *ast.SumDecl:
			sym.Type = &SumType{
				Decl: d,
			}
		}
	}
}

// resolveTypeDecls resolves type declarations by binding the ast decl
// with semantic check declarations
func (c *Checker) resolveTypeDecls() {
	for _, decl := range c.typeDecls {
		name := typeDeclName(decl)
		sym := c.pkgScope.LookupLocal(name)
		if sym == nil || sym.Type == nil {
			continue
		}

		switch d := decl.(type) {
		case *ast.DefinedTypeDecl:
			t := sym.Type.(*NamedType)
			t.UnderlyingType = c.resolveType(d.Type)

		case *ast.StructDecl:
			t := sym.Type.(*StructType)
			t.Fields = c.resolveStructFields(d.Fields)

		case *ast.InterfaceDecl:
			t := sym.Type.(*InterfaceType)
			t.Methods = c.resolveInterfaceMethods(d.Methods)

		case *ast.EnumDecl:
			t := sym.Type.(*EnumType)
			t.Variants = c.resolveEnumVariants(d.Variants)

		case *ast.SumDecl:
			t := sym.Type.(*SumType)
			t.Variants = c.resolveSumVariants(d.Variants)
		}
	}
}

// resolveType resolves the type passed as parameter in order to fetch and
// return its semantic type
func (c *Checker) resolveType(t ast.Type) Type {
	switch v := t.(type) {
	case *ast.NamedType:
		return c.resolveNamedType(v)

	case *ast.ArrayType:
		elem := c.resolveType(v.Elem)
		// temporary keeping this
		if elem == nil {
			return TInvalid
		}
		return &ArrayType{Elem: elem}

	case *ast.SliceType:
		elem := c.resolveType(v.Elem)
		// temporary keeping this
		if elem == nil {
			return TInvalid
		}
		return &SliceType{Elem: elem}

	case *ast.MapType:
		key := c.resolveType(v.KeyType)
		if key == nil {
			return TInvalid
		}
		value := c.resolveType(v.ValueType)
		if value == nil {
			return TInvalid
		}
		m := &MapType{
			Key:   key,
			Value: value,
		}
		if v.KindKW.Kind == token.KWHashMap {
			m.Kind = MapHash
		}
		return m
	}
	return nil
}

// resolveNamedType resolves named type to return a semantic type
func (c *Checker) resolveNamedType(t *ast.NamedType) Type {
	// we take only the first id for now
	if len(t.Parts) != 1 {
		return TInvalid
	}

	part := t.Parts[0]
	switch part.Kind {
	case token.KWBool:
		return TBool
	case token.KWInt:
		return TInt
	case token.KWInt8:
		return TInt8
	case token.KWInt32:
		return TInt32
	case token.KWInt64:
		return TInt64
	case token.KWUint:
		return TUInt
	case token.KWUint8:
		return TUInt8
	case token.KWUint32:
		return TUInt32
	case token.KWUint64:
		return TUInt64
	case token.KWFloat:
		return TFloat
	case token.KWFloat32:
		return TFloat32
	case token.KWFloat64:
		return TFloat64
	case token.KWString:
		return TString
	case token.Ident:
		sym := c.pkgScope.LookupLocal(part.Value)
		if sym == nil || sym.Kind != SymType {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("unknown %q type", part.Value)})
			return TInvalid
		}
		return sym.Type
	}
	return nil
}

// resolveStructFields resolves the ast struct field declaration to return
// the semantic view of the field list. A diagnostic is emitted when duplicates found
func (c *Checker) resolveStructFields(fields []*ast.FieldDecl) []StructField {
	var out []StructField
	seen := make(map[string]*ast.FieldDecl)

	for _, field := range fields {
		if prev := seen[field.Name.Value]; prev != nil {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("struct field %q already declared", field.Name.Value)})
			continue
		}

		seen[field.Name.Value] = field
		out = append(out, StructField{
			Name: field.Name.Value,
			Type: c.resolveType(field.Type),
		})
	}
	return out
}

// resolveInterfaceMethods resolves the ast interface declaration to return
// the semantic view of the method list. A diagnostic is emitted when duplicates found
func (c *Checker) resolveInterfaceMethods(methods []ast.InterfaceMethod) []FuncMethod {
	var out []FuncMethod
	seen := make(map[string]*ast.InterfaceMethod)

	for _, fn := range methods {
		if prev := seen[fn.Name.Value]; prev != nil {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("method name %q already declared", fn.Name.Value)})
			continue
		}

		seen[fn.Name.Value] = &fn
		out = append(out, c.resolveInterfaceMethod(fn))
	}
	return out
}

// resolveInterfaceMethod resolves the ast InterfaceMethod declaration to return
// the semantic view of the FuncMethod. A diagnostic is emitted when duplicates found
func (c *Checker) resolveInterfaceMethod(m ast.InterfaceMethod) FuncMethod {
	return FuncMethod{
		Name: m.Name.Value,
		FuncType: &FuncType{
			Params:  c.resolveParams("param", m.Params),
			Results: c.resolveParams("result", m.Results.List),
		},
	}
}

// resolveFuncSignatures resolves the ast FuncDecl declaration to return
// the semantic view of the FuncMethod. A diagnostic is emitted when duplicates found
func (c *Checker) resolveFuncSignatures() {
	for _, fn := range c.funcDecls {
		sym := c.pkgScope.Lookup(fn.Name.Value)
		if sym != nil {
			sym.Type = &FuncMethod{
				Name: fn.Name.Value,
				FuncType: &FuncType{
					Params:  c.resolveParams("param", fn.Params),
					Results: c.resolveParams("result", fn.Results.List),
				},
			}
		}
	}
}

// resolveParams resolves the ast FuncDecl declaration to return
// the semantic view of the Param list. A diagnostic is emitted when duplicates found
func (c *Checker) resolveParams(kind string, pr []ast.Param) []Param {
	var out []Param
	seen := make(map[string]*ast.Param)

	for _, p := range pr {
		if p.Name.Value != "" {
			if prev := seen[p.Name.Value]; prev != nil {
				c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("method %s name %q already declared", kind, p.Name.Value)})
				continue
			}

			seen[p.Name.Value] = &p
		}
		out = append(out, Param{
			Name: p.Name.Value,
			Type: c.resolveType(p.Type),
		})
	}
	return out
}

// resolveEnumVariants resolves the ast FuncDecl declaration to return
// the semantic view of the string list. A diagnostic is emitted when duplicates found
func (c *Checker) resolveEnumVariants(variants []token.Token) []string {
	var out []string
	seen := make(map[string]*token.Token)

	for _, v := range variants {
		if prev := seen[v.Value]; prev != nil {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("variant %q already declared", v.Value)})
			continue
		}

		seen[v.Value] = &v
		out = append(out, v.Value)
	}
	return out
}

// resolveSumVariants resolves the ast FuncDecl declaration to return
// the semantic view of the SumVariant list. A diagnostic is emitted when duplicates found
func (c *Checker) resolveSumVariants(variants []ast.SumVariant) []SumVariant {
	var out []SumVariant
	seen := make(map[string]*ast.SumVariant)

	for _, v := range variants {
		if prev := seen[v.Name.Value]; prev != nil {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("sum variant %q already declared", v.Name.Value)})
			continue
		}

		seen[v.Name.Value] = &v
		out = append(out, SumVariant{
			Name:  v.Name.Value,
			Field: c.resolveParams("variant", v.Params),
		})
	}
	return out
}
