package semantic

import (
	"fmt"
	"strconv"

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
	c.checkTopLevelValues(file)
	c.checkFuncBodies()
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
		case *ast.ConstDecl:
			c.declareConstSymbol(d)
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
	case *ast.ConstDecl:
		return d.Name.Value
	default:
		return ""
	}
}

// exprName returns the name of the expression declaration
func exprName(decl ast.Expr) string {
	switch d := decl.(type) {
	case *ast.IdentExpr:
		return d.Name.Value
	case *ast.IndexExpr:
		return exprName(d.X)
	case *ast.SelectorExpr:
		return exprName(d.X)
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

// declareConstSymbol declares new type symbol with its name
// and append diagnostics errors when already exists
func (c *Checker) declareConstSymbol(decl *ast.ConstDecl) {
	name := typeDeclName(decl)
	if name == "" {
		return
	}

	sym := &Symbol{
		Name: name,
		Kind: SymConst,
		Decl: decl,
	}

	if !c.pkgScope.Declare(sym) {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("symbol %q already declared", name)})
		return
	}

	c.constDecls = append(c.constDecls, decl)
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
			sym.Type = &NamedType{
				UnderlyingType: &StructType{
					Decl: d,
				},
			}

		case *ast.InterfaceDecl:
			sym.Type = &NamedType{
				UnderlyingType: &InterfaceType{
					Decl: d,
				},
			}

		case
			*ast.EnumDecl:
			sym.Type = &NamedType{
				UnderlyingType: &EnumType{
					Decl: d,
				},
			}

		case *ast.SumDecl:
			sym.Type = &NamedType{
				UnderlyingType: &SumType{
					Decl: d,
				},
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
			t := sym.Type.(*NamedType)
			st := t.UnderlyingType.(*StructType)
			st.Fields = c.resolveStructFields(d.Fields)

		case *ast.InterfaceDecl:
			t := sym.Type.(*NamedType)
			it := t.UnderlyingType.(*InterfaceType)
			it.Methods = c.resolveInterfaceMethods(d.Methods)

		case *ast.EnumDecl:
			t := sym.Type.(*NamedType)
			et := t.UnderlyingType.(*EnumType)
			et.Variants = c.resolveEnumVariants(d.Variants)

		case *ast.SumDecl:
			t := sym.Type.(*NamedType)
			st := t.UnderlyingType.(*SumType)
			st.Variants = c.resolveSumVariants(d.Variants)
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
		// TODO: temporary keeping this
		if elem == nil {
			return TInvalid
		}
		len := v.Len.(*ast.IntLitExpr)
		l, _ := strconv.Atoi(len.Name.Value)
		return &ArrayType{Len: int64(l), Elem: elem}

	case *ast.SliceType:
		elem := c.resolveType(v.Elem)
		// TODO: temporary keeping this
		if elem == nil {
			return TInvalid
		}
		return &SliceType{Elem: elem}

	case *ast.MapType:
		key := c.resolveType(v.KeyType)
		if key == nil {
			return TInvalid
		}

		if !isMapKeyType(key) {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("invalid map key type")})
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
		var sym *Symbol
		if c.useScope {
			sym = c.scope.Lookup(part.Value)
		} else {
			sym = c.pkgScope.LookupLocal(part.Value)
		}
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
		if sym := c.pkgScope.Lookup(fn.Name.Value); sym != nil {
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

// checkTopLevelValues checks const values
func (c *Checker) checkTopLevelValues(file *ast.File) {
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.ConstDecl:
			c.checkConstDecl(d)
		}
	}
}

// checkConstDecl validates constant targetType and valueType.
// An error is emitted if any
func (c *Checker) checkConstDecl(decl *ast.ConstDecl) {
	targetType := c.resolveType(decl.Type)
	valueType := c.checkExpr(decl.Init)

	if !IsAssignableTo(targetType, valueType) {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("cannot assign value of type %T to const of type %T", valueType, targetType)})
		return
	}
	name := typeDeclName(decl)
	sym := c.pkgScope.Lookup(name)
	sym.Type = targetType
}

// checkExpr returns the type of the expression
func (c *Checker) checkExpr(expr ast.Expr) Type {
	switch t := expr.(type) {
	case *ast.IntLitExpr:
		return TInt

	case *ast.FloatLitExpr:
		return TFloat

	case *ast.BoolLitExpr:
		return TBool

	case *ast.StringLitExpr:
		return TString

	case *ast.IdentExpr:
		var sym *Symbol
		if c.useScope {
			sym = c.scope.Lookup(t.Name.Value)
		} else {
			sym = c.pkgScope.LookupLocal(t.Name.Value)
		}
		if sym == nil || sym.Type == nil {
			return TInvalid
		}
		return sym.Type

	case *ast.UnaryExpr:
		right := c.checkExpr(t.Right)
		if SupportsUnaryOp(right, t.Operator.Kind) {
			return right
		}
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("invalid unary operatation %s with type %s", t.Operator.Value, right.String())})
		return TInvalid

	case *ast.BinaryExpr:
		left := c.checkExpr(t.Left)
		right := c.checkExpr(t.Right)

		if IsIdentical(left, right) && SupportsBinaryOp(left, t.Operator.Kind) {
			return left
		}
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("invalid binary operation %q of type %s with type %s", t.Operator.Value, left.String(), right.String())})
		return TInvalid

	case *ast.CallExpr:
		calleeType := c.checkExpr(t.Callee)
		if named, ok := calleeType.(*NamedType); ok {
			if len(t.Args) != 1 {
				c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("too many arguments in %s", named.Name)})
				return TInvalid
			}
			if IsConvertibleTo(named, c.checkExpr(t.Args[0])) {
				return named
			}
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("cannot convert in %s", named.Name)})
			return TInvalid
		}

		fn, ok := calleeType.(*FuncMethod)
		if !ok {
			return TInvalid
		}
		if len(t.Args) != len(fn.FuncType.Params) {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("too many arguments %s func params required %d got %d", fn.Name, len(fn.FuncType.Params), len(t.Args))})
			return TInvalid
		}
		for k, v := range fn.FuncType.Params {
			x := c.checkExpr(t.Args[k])
			if !IsAssignableTo(v.Type, x) {
				c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("cannot assign value of type %T to const of type %T", v.Type, x)})
				return TInvalid
			}
		}
		if len(fn.FuncType.Results) == 0 {
			return nil
		}
		if len(fn.FuncType.Results) > 1 {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("too many arguments returned by %s func", fn.Name)})
			return TInvalid
		}
		return fn.FuncType.Results[0].Type

	case *ast.ParenExpr:
		return c.checkExpr(t.Inner)

	case *ast.IndexExpr:
		baseType := c.checkExpr(t.X)
		underlying := unwrapNamed(baseType)
		index := c.checkExpr(t.Index)

		switch decl := underlying.(type) {
		case *SliceType:
			if !IsIdentical(index, TInt) {
				c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("invalid index expression of type %#v", index)})
				return TInvalid
			}
			return decl.Elem

		case *ArrayType:
			if !IsIdentical(index, TInt) {
				c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("invalid index expression of type %#v", index)})
				return TInvalid
			}
			return decl.Elem

		case *MapType:
			if !IsIdentical(decl.Key, index) {
				c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("invalid map index expression of type %#v", index)})
				return TInvalid
			}
			return decl.Value

		default:
			return TInvalid
		}

	case *ast.SelectorExpr:
		return c.checkSelectorExpr(t)

	default:
		return TInvalid
	}
}

// isMapKeyType validates provided type as key type
func isMapKeyType(t Type) bool {
	if IsNumeric(t) || IsBool(t) || IsString(t) {
		return true
	}
	return false
}

// checkFuncBodies validates functions bodies.
// An error is emitted if any
func (c *Checker) checkFuncBodies() {
	for _, fn := range c.funcDecls {
		c.checkFuncBody(fn)
	}
}

func (c *Checker) checkFuncBody(fn *ast.FuncDecl) {
	oldScope := c.scope
	oldFunc := c.currentFunc
	defer func() {
		c.scope = oldScope
		c.currentFunc = oldFunc
		c.useScope = false
	}()

	sym := c.pkgScope.Lookup(fn.Name.Value)
	if sym == nil {
		return
	}

	fnType, ok := sym.Type.(*FuncMethod)
	if !ok {
		return
	}

	c.currentFunc = fnType.FuncType
	c.scope = NewScope(c.pkgScope)
	c.useScope = true

	for _, p := range fnType.FuncType.Params {
		c.scope.Declare(&Symbol{
			Name: p.Name,
			Kind: SymVar,
			Type: p.Type,
		})
	}

	c.checkBlockStmt(fn.Body)
}

func (c *Checker) checkBlockStmt(block *ast.BlockStmt) {
	if block == nil {
		return
	}

	oldScope := c.scope
	defer func() {
		c.scope = oldScope
	}()

	blockScope := NewScope(c.scope)
	c.scope = blockScope

	for _, stmt := range block.Stmts {
		switch t := stmt.(type) {
		case *ast.DeclStmt:
			switch decl := t.Decl.(type) {
			case *ast.ConstDecl:
				c.checkScopeConstDecl(decl)

			case *ast.VarDecl:
				c.checkScopeVarDecl(decl)

			case *ast.DefinedTypeDecl:
				c.checkDefinedTypeDecl(decl)

			case *ast.StructDecl:
				c.checkStructDecl(decl)

			case *ast.EnumDecl:
				c.checkEnumDecl(decl)

			case *ast.SumDecl:
				c.checkSumDecl(decl)

			case *ast.InterfaceDecl:
				c.checkInterfaceDecl(decl)

			default:
				c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("declaration %#v not managed", decl)})
			}

		case *ast.AssignStmt:
			switch t.Operator.Kind {
			case token.Assign:
				c.checkSimpleAssignStmt(t)

			case token.Define:
				c.checkDefineAssignStmt(t)
			}

		case *ast.ReturnStmt:
			c.checkReturnStmt(t)

		case *ast.IncDecStmt:
			c.checkIncDecStmt(t)

		case *ast.ExprStmt:
			c.checkExprStmt(t)

		default:
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("statement %#v not managed", stmt)})
		}
	}
}

// checkScopeConstDecl validates constant targetType and valueType.
// An error is emitted if any
func (c *Checker) checkScopeConstDecl(decl *ast.ConstDecl) {
	sym := c.scope.Lookup(decl.Name.Value)
	if sym != nil {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("const %s already declared", decl.Name.Value)})
		return
	}

	targetType := c.resolveType(decl.Type)
	valueType := c.checkExpr(decl.Init)

	if !IsAssignableTo(targetType, valueType) {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("cannot assign value of type %T to var of type %T", valueType, targetType)})
		return
	}

	c.scope.Declare(&Symbol{
		Name: decl.Name.Value,
		Kind: SymConst,
		Type: targetType,
		Decl: decl,
	})
}

// checkScopeVarDecl validates constant targetType and valueType.
// An error is emitted if any
func (c *Checker) checkScopeVarDecl(decl *ast.VarDecl) {
	sym := c.scope.Lookup(decl.Name.Value)
	if sym != nil {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("variable %s already declared", decl.Name.Value)})
		return
	}

	targetType := c.resolveType(decl.Type)
	valueType := c.checkExpr(decl.Init)

	if !IsAssignableTo(targetType, valueType) {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("cannot assign value of type %T to var of type %T", valueType, targetType)})
		return
	}

	c.scope.Declare(&Symbol{
		Name: decl.Name.Value,
		Kind: SymVar,
		Type: targetType,
		Decl: decl,
	})
}

// checkAssignableExpr returns valid assignable expression.
// An error is emitted if any
func (c *Checker) checkAssignableExpr(expr ast.Expr) Type {
	switch t := expr.(type) {
	case *ast.IdentExpr:
		return c.checkExpr(t)

	case *ast.IndexExpr:
		return c.checkExpr(t)

	case *ast.SelectorExpr:
		return c.checkSelectorExpr(t)

	default:
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("unsupported expression %#v", expr)})
		return TInvalid
	}
}

// checkSimpleAssignStmt validates simple assigment statements like x = 1 where x has already been defined.
// An error is emitted if any
func (c *Checker) checkSimpleAssignStmt(decl *ast.AssignStmt) {
	name := exprName(decl.Left)
	sym := c.scope.Lookup(name)
	if sym == nil {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("assigment %s is undefined", name)})
		return
	}

	if sym.Kind == SymConst {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("reassign const %s value is forbidden", name)})
		return
	}

	targetType := c.checkAssignableExpr(decl.Left)
	valueType := c.checkExpr(decl.Right)

	if !IsAssignableTo(targetType, valueType) {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("cannot assign value of type %T to variable of type %T", valueType, targetType)})
		return
	}

	sym.Type = targetType
}

// isNumericExpr detects if expression is a numeric only expression
func isNumericExpr(expr ast.Expr) bool {
	switch t := expr.(type) {
	case *ast.IntLitExpr, *ast.FloatLitExpr:
		return true

	case *ast.ParenExpr:
		return isNumericExpr(t.Inner)

	case *ast.BinaryExpr:
		return isNumericExpr(t.Left) && isNumericExpr(t.Right)

	case *ast.UnaryExpr:
		return isNumericExpr(t.Right)

	default:
		return false
	}
}

// checkDefineAssignStmt defines new assigment statement where x := y and y has already been defined.
// define assigment like x = 1 is forbidden as we cannot infer the value type.
// An error is emitted if any
func (c *Checker) checkDefineAssignStmt(decl *ast.AssignStmt) {
	valueType := c.checkExpr(decl.Right)
	if IsInvalid(valueType) {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("expression %#v is invalid", decl.Right)})
		return
	}

	if isNumericExpr(decl.Right) {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("cannot use numeric only expression in := declaration")})
		return
	}

	x, ok := decl.Left.(*ast.IdentExpr)
	if !ok {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("variable %#v is not an identifier", decl)})
		return
	}

	sym := c.scope.Lookup(x.Name.Value)
	if sym != nil {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("variable %s is already defined", x.Name.Value)})
		return
	}

	c.scope.Declare(&Symbol{
		Name: x.Name.Value,
		Kind: SymVar,
		Type: valueType,
	})
}

// checkReturnStmt checks returned values statement types and length.
// An error is emitted if any
func (c *Checker) checkReturnStmt(decl *ast.ReturnStmt) {
	if decl == nil {
		return
	}

	if len(c.currentFunc.Results) != len(decl.Values) {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("number of returned values is invalid, wanted %d got %d", len(c.currentFunc.Results), len(decl.Values))})
		return
	}

	for k, v := range decl.Values {
		expr := c.checkExpr(v)
		if !IsIdentical(c.currentFunc.Results[k].Type, expr) {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("cannot use a value of type %T as %T in return statement", expr, c.currentFunc.Results[k].Type)})
			return
		}
	}
}

// checkIncDecStmt validates increment/decrement statement.
// An error is emitted if any
func (c *Checker) checkIncDecStmt(decl *ast.IncDecStmt) {
	x, ok := decl.X.(*ast.IdentExpr)
	if !ok {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("variable %#v is not an identifier", decl)})
		return
	}

	sym := c.scope.Lookup(x.Name.Value)
	if sym == nil {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("assigment %s is undefined", x.Name.Value)})
		return
	}

	if sym.Kind == SymConst {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("const %s cannot be modified", x.Name.Value)})
		return
	}

	if !IsNumeric(sym.Type) {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("variable %s is a non-numeric type", x.Name.Value)})
		return
	}
}

// checkDefinedTypeDecl validates type declaration statement.
// An error is emitted if any
func (c *Checker) checkDefinedTypeDecl(decl *ast.DefinedTypeDecl) {
	sym := c.scope.Lookup(decl.Name.Value)
	if sym != nil {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("type %s already declared", decl.Name.Value)})
		return
	}

	c.scope.Declare(&Symbol{
		Name: decl.Name.Value,
		Kind: SymType,
		Type: &NamedType{
			Decl:           decl,
			UnderlyingType: c.resolveType(decl.Type),
		},
	})
}

// checkStructDecl validates struct statement.
// An error is emitted if any
func (c *Checker) checkStructDecl(decl *ast.StructDecl) {
	sym := c.scope.Lookup(decl.Name.Value)
	if sym != nil {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("type %s already declared", decl.Name.Value)})
		return
	}

	c.scope.Declare(&Symbol{
		Name: decl.Name.Value,
		Kind: SymType,
		Type: &NamedType{
			UnderlyingType: &StructType{
				Decl:   decl,
				Fields: c.resolveStructFields(decl.Fields),
			},
		},
	})
}

// checkEnumDecl validates struct statement.
// An error is emitted if any
func (c *Checker) checkEnumDecl(decl *ast.EnumDecl) {
	sym := c.scope.Lookup(decl.Name.Value)
	if sym != nil {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("type %s already declared", decl.Name.Value)})
		return
	}

	c.scope.Declare(&Symbol{
		Name: decl.Name.Value,
		Kind: SymType,
		Type: &NamedType{
			UnderlyingType: &EnumType{
				Decl:     decl,
				Variants: c.resolveEnumVariants(decl.Variants),
			},
		},
	})
}

// checkSumDecl validates struct statement.
// An error is emitted if any
func (c *Checker) checkSumDecl(decl *ast.SumDecl) {
	sym := c.scope.Lookup(decl.Name.Value)
	if sym != nil {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("type %s already declared", decl.Name.Value)})
		return
	}

	c.scope.Declare(&Symbol{
		Name: decl.Name.Value,
		Kind: SymType,
		Type: &NamedType{
			UnderlyingType: &SumType{
				Decl:     decl,
				Variants: c.resolveSumVariants(decl.Variants),
			},
		},
	})
}

// checkInterfaceDecl validates interface declaration statement.
// An error is emitted if any
func (c *Checker) checkInterfaceDecl(decl *ast.InterfaceDecl) {
	sym := c.scope.Lookup(decl.Name.Value)
	if sym != nil {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("type %s already declared", decl.Name.Value)})
		return
	}

	c.scope.Declare(&Symbol{
		Name: decl.Name.Value,
		Kind: SymType,
		Type: &NamedType{
			UnderlyingType: &InterfaceType{
				Decl:    decl,
				Methods: c.resolveInterfaceMethods(decl.Methods),
			},
		},
	})
}

// checkExprStmt validates expression statement.
// An error is emitted if any
func (c *Checker) checkExprStmt(stmt *ast.ExprStmt) {
	call, ok := stmt.Expr.(*ast.CallExpr)
	if !ok {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("expression statement must be a function call")})
		return
	}

	calleType := c.checkExpr(call.Callee)
	if _, ok := calleType.(*FuncMethod); !ok {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("expression statement must be a function call")})
		return
	}

	_ = c.checkExpr(stmt.Expr)
}

// checkSelectorExpr validates selector expression and return its type.
// An error is emitted if any
func (c *Checker) checkSelectorExpr(expr *ast.SelectorExpr) Type {
	baseType := c.checkExpr(expr.X)
	underlying := unwrapNamed(baseType)

	switch t := underlying.(type) {
	case *StructType:
		tp, ok := lookupStructField(t, expr.Selector.Value)
		if !ok {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("unknown field %q", expr.Selector.Value)})
			return TInvalid
		}
		return tp

	case *InterfaceType:
		tp, ok := lookupInterfaceMethods(t, expr.Selector.Value)
		if !ok {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("unknown method %q", expr.Selector.Value)})
			return TInvalid
		}
		return tp

	default:
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("invalid type")})
		return TInvalid
	}
}

// unwrapNamed returns underlying type if it's a named type otherwise the initial type
func unwrapNamed(t Type) Type {
	for {
		named, ok := t.(*NamedType)
		if !ok {
			return t
		}
		t = named.UnderlyingType
	}
}

// lookupStructField loops into struct field to match provided name.
// It returns its type and true when found
func lookupStructField(st *StructType, name string) (Type, bool) {
	for _, f := range st.Fields {
		if f.Name == name {
			return f.Type, true
		}
	}
	return nil, false
}

// lookupInterfaceMethods loops into interface methods to match provided name.
// It returns its type and true when found
func lookupInterfaceMethods(it *InterfaceType, name string) (Type, bool) {
	for i := range it.Methods {
		if it.Methods[i].Name == name {
			return &it.Methods[i], true
		}
	}
	return nil, false
}
