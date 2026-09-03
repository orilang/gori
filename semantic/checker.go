package semantic

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/token"
)

// NewScope allows us to create new scope base on provided config.
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

// declareNoShadow declares new symbol by enforcing there is no shadowing
// and return true when NOT exists
func (c *Checker) declareNoShadow(scope *Scope, sym *Symbol, kind string) bool {
	if scope == nil {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("cannot declare %q on nil scope", sym.Name)})
		return false
	}

	if sym.Name == "" || sym.Name == "_" {
		return true
	}

	if scope.Lookup(sym.Name) != nil {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("%s %q already declared", kind, sym.Name)})
		return false
	}

	if !scope.Declare(sym) {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("internal checker error: %s %q already declared in current scope", kind, sym.Name)})
		return false
	}

	return true
}

// Lookup allows us to loop over parent scope to find the provided symbol
// and returns its related Symbol if exists
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
	universalTypes := NewScope(nil)
	universalTypes.Declare(&Symbol{Name: "bool", Kind: SymType, Type: TBool})
	universalTypes.Declare(&Symbol{Name: "int", Kind: SymType, Type: TInt})
	universalTypes.Declare(&Symbol{Name: "int8", Kind: SymType, Type: TInt8})
	universalTypes.Declare(&Symbol{Name: "int32", Kind: SymType, Type: TInt32})
	universalTypes.Declare(&Symbol{Name: "int64", Kind: SymType, Type: TInt64})
	universalTypes.Declare(&Symbol{Name: "uint", Kind: SymType, Type: TUInt})
	universalTypes.Declare(&Symbol{Name: "uint8", Kind: SymType, Type: TUInt8})
	universalTypes.Declare(&Symbol{Name: "uint32", Kind: SymType, Type: TUInt32})
	universalTypes.Declare(&Symbol{Name: "uint64", Kind: SymType, Type: TUInt64})
	universalTypes.Declare(&Symbol{Name: "float", Kind: SymType, Type: TFloat})
	universalTypes.Declare(&Symbol{Name: "float32", Kind: SymType, Type: TFloat32})
	universalTypes.Declare(&Symbol{Name: "float64", Kind: SymType, Type: TFloat64})
	universalTypes.Declare(&Symbol{Name: "string", Kind: SymType, Type: TString})

	return &Checker{
		pkgScope: NewScope(universalTypes),
	}
}

// Check performs the type checking step to validate
// all code definition structure and fill diagnostics when
// errros are found
func (c *Checker) Check(file *ast.File) (Program, Diagnostics) {
	pf := &File{Package: file.PackageKW.Value}
	c.program.Files = append(c.program.Files, pf)
	c.programFileIndex = 0

	c.collectTopLevelSymbols(file)
	c.createTypeObjects()
	c.resolveTypeDecls()
	c.declareComptimeDecls()
	c.resolveFuncSignatures()
	c.resolveMethodSignatures()
	c.checkImplementsDecls()
	c.checkComptimeValues()
	c.checkTopLevelValues(file)
	c.checkFuncBodies()

	if len(c.errors) == 0 {
		return c.program, nil
	}
	return Program{}, c.errors
}

// HasErrors returns true when errors found
func (d Diagnostics) HasErrors() bool {
	return len(d) > 0
}

// collectTopLevelSymbols collects top levels symbols names first
// to make sure that all definitions exists before creating
// semantic objects and resolve remaining contents.
// This prevents having types that does not exists.
func (c *Checker) collectTopLevelSymbols(file *ast.File) {
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case ast.TypeDecl:
			if !c.declareNoShadow(c.pkgScope, &Symbol{
				Name: typeDeclName(decl),
				Kind: SymType,
				Decl: decl,
			}, "symbol") {
				continue
			}
			c.typeDecls = append(c.typeDecls, d)

		case *ast.FuncDecl:
			if d.Receiver != nil {
				c.methodDecls = append(c.methodDecls, d)
			} else {
				if !c.declareNoShadow(c.pkgScope, &Symbol{
					Name: d.Name.Value,
					Kind: SymFunc,
					Decl: d,
				}, "symbol") {
					continue
				}
				c.funcDecls = append(c.funcDecls, d)
			}

		case *ast.ConstDecl:
			sym := &Symbol{
				Name: typeDeclName(decl),
				Kind: SymConst,
				Decl: decl,
			}
			if !c.declareNoShadow(c.pkgScope, sym, "symbol") {
				continue
			}
			c.constDecls = append(c.constDecls, d)

		case *ast.ImplementsDecl:
			c.implDecls = append(c.implDecls, d)

		case *ast.ComptimeBlockDecl:
			c.comptimeDecls = append(c.comptimeDecls, d.Decl)
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

// declareMethodSymbol declares new type symbol with its name
// and append diagnostics errors when already exists
func (c *Checker) declareMethodSymbol(receiver *NamedType, fm *FuncMethod) {
	if c.methods == nil {
		c.methods = make(map[*NamedType]map[string]*FuncMethod)
	}

	rcv := c.methods[receiver]
	if rcv == nil {
		rcv = make(map[string]*FuncMethod)
		c.methods[receiver] = rcv
	}

	if _, exists := rcv[fm.Name]; exists {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("method %q already declared", fm.Name)})
		return
	}

	rcv[fm.Name] = fm
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
				Name: d.Name.Value,
				Decl: d,
			}

		case *ast.StructDecl:
			sym.Type = &NamedType{
				Name: d.Name.Value,
				Decl: d,
				UnderlyingType: &StructType{
					Decl: d,
				},
			}

		case *ast.InterfaceDecl:
			sym.Type = &NamedType{
				Name: d.Name.Value,
				Decl: d,
				UnderlyingType: &InterfaceType{
					Decl: d,
				},
			}

		case
			*ast.EnumDecl:
			sym.Type = &NamedType{
				Name: d.Name.Value,
				Decl: d,
				UnderlyingType: &EnumType{
					Decl: d,
				},
			}

		case *ast.SumDecl:
			sym.Type = &NamedType{
				Name: d.Name.Value,
				Decl: d,
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

// resolveType resolves the type passed as parameter to fetch and
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

		len, ok := c.evalArrayLen(v.Len)
		if !ok || len < 0 {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid array length type")})
			return TInvalid
		}
		return &ArrayType{Len: len, Elem: elem}

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
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid map key type")})
			return TInvalid
		}

		value := c.resolveType(v.ValueType)
		if value == nil {
			return TInvalid
		}

		if v.KindKW.Kind == token.KWHashMap {
			return &HashMapType{
				Key:   key,
				Value: value,
			}
		}

		return &MapType{
			Key:   key,
			Value: value,
		}
	}
	return nil
}

// evalArrayLen validates array length
func (c *Checker) evalArrayLen(expr ast.Expr) (int64, bool) {
	switch t := expr.(type) {
	case *ast.IntLitExpr:
		v, err := strconv.ParseInt(t.Name.Value, 10, 64)
		if err != nil {
			return 0, false
		}
		return v, true

	case *ast.ParenExpr:
		return c.evalArrayLen(t.Inner)

	case *ast.BinaryExpr:
		left, lok := c.evalArrayLen(t.Left)
		right, rok := c.evalArrayLen(t.Right)
		if lok && rok {
			switch t.Operator.Kind {
			case token.Plus:
				return left + right, true
			case token.Minus:
				return left - right, true
			case token.Slash:
				if right == 0 {
					c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("division by 0 is forbidden")})
					return 0, false
				}
				return left / right, true
			case token.Star:
				return left * right, true
			}
		}
		return 0, false

	case *ast.UnaryExpr:
		v, ok := c.evalArrayLen(t.Right)
		if !ok {
			return 0, false
		}

		switch t.Operator.Kind {
		case token.Plus:
			return v, true

		case token.Minus:
			return -v, true
		}
	}

	return 0, false
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
		if c.useScope || c.inComptimeFunc {
			sym = c.scope.Lookup(part.Value)
		} else {
			sym = c.pkgScope.LookupLocal(part.Value)
		}
		if sym == nil || sym.Kind != SymType {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("unknown %q type", part.Value)})
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
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("struct field %q already declared", field.Name.Value)})
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
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("method name %q already declared", fn.Name.Value)})
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

// resolveMethodSignatures resolves the ast FuncDecl declaration with receiver to return
// the semantic view of the FuncMethod. A diagnostic is emitted when duplicates found
func (c *Checker) resolveMethodSignatures() {
	for _, decl := range c.methodDecls {
		recvType := c.resolveType(decl.Receiver.Type)
		namedRcv, ok := recvType.(*NamedType)
		if !ok {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("receiver must be a named type got %#v", namedRcv)})
			continue
		}

		if _, ok := unwrapNamed(namedRcv).(*InterfaceType); ok {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("receiver cannot be an interface type")})
			continue
		}

		fm := &FuncMethod{
			Name: decl.Name.Value,
			FuncType: &FuncType{
				Params:  c.resolveParams("param", decl.Params),
				Results: c.resolveParams("result", decl.Results.List),
			},
		}

		c.declareMethodSymbol(namedRcv, fm)
	}
}

// checkImplementsDecls validates all "implements" declaration for interface
func (c *Checker) checkImplementsDecls() {
	for _, impl := range c.implDecls {
		c.checkImplementsDecl(impl)
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
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("%s name %q already declared", kind, p.Name.Value)})
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
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("variant %q already declared", v.Value)})
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
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("sum variant %q already declared", v.Name.Value)})
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
	valueType, expr := c.checkExpr(decl.Init)

	if !IsAssignableTo(targetType, valueType) {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("cannot assign value of type %T to const of type %T", valueType, targetType)})
		return
	}
	sym := c.pkgScope.Lookup(typeDeclName(decl))
	sym.Type = targetType
	sym.Decl = decl

	c.program.Files[c.programFileIndex].Decls = append(c.program.Files[c.programFileIndex].Decls, &ConstDecl{Name: sym.Name, Symbol: resolvedSymbol(sym), Eq: decl.Eq.Kind, Init: expr})
}

// checkExpr returns the type of the expression
func (c *Checker) checkExpr(expr ast.Expr) (Type, Expr) {
	switch t := expr.(type) {
	case *ast.IntLitExpr:
		return TInt, &IntLitExpr{Type: TInt, Value: t.Name.Value}

	case *ast.FloatLitExpr:
		return TFloat, &FloatLitExpr{Type: TFloat, Value: t.Name.Value}

	case *ast.BoolLitExpr:
		return TBool, &BoolLitExpr{Type: TBool, Value: t.Name.Value}

	case *ast.StringLitExpr:
		return TString, &StringLitExpr{Type: TString, Value: t.Name.Value}

	case *ast.SliceLitExpr:
		return c.resolveType(t.Type), nil

	case *ast.IdentExpr:
		var sym *Symbol
		if c.useScope {
			sym = c.scope.Lookup(t.Name.Value)
		} else {
			sym = c.pkgScope.Lookup(t.Name.Value)
		}
		if sym == nil || sym.Type == nil {
			return TInvalid, nil
		}
		return sym.Type, &IdentExpr{Type: sym.Type, Symbol: resolvedSymbol(sym), Value: t.Name.Value}

	case *ast.UnaryExpr:
		right, ex := c.checkExpr(t.Right)
		if SupportsUnaryOp(right, t.Operator.Kind) {
			return right, &UnaryExpr{Type: right, Operator: t.Operator.Kind, Right: ex}
		}
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid unary operatation %s with type %s", t.Operator.Value, right.String())})
		return TInvalid, nil

	case *ast.BinaryExpr:
		left, lex := c.checkExpr(t.Left)
		right, rex := c.checkExpr(t.Right)

		if IsIdentical(left, right) && SupportsBinaryOp(left, t.Operator.Kind) {
			switch t.Operator.Kind {
			case token.Eq, token.Neq, token.And, token.Or, token.Lt, token.Lte, token.Gt, token.Gte:
				return TBool, &BinaryExpr{Type: TBool, Left: lex, Operator: t.Operator.Kind, Right: rex}

			default:
				return left, &BinaryExpr{Type: left, Left: lex, Operator: t.Operator.Kind, Right: rex}
			}
		}
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid binary operation %q of type %s with type %s", t.Operator.Value, left.String(), right.String())})
		return TInvalid, nil

	case *ast.CallExpr:
		if sel, ok := t.Callee.(*ast.SelectorExpr); ok {
			selx, _ := c.checkExpr(sel.X)
			if named, ok := selx.(*NamedType); ok {
				if method, ok := c.lookupMethodType(named, sel.Selector.Value); ok {
					if len(t.Args) != len(method.FuncType.Params) {
						c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("too many arguments for %s func params, expected %d got %d", method.Name, len(method.FuncType.Params), len(t.Args))})
						return TInvalid, nil
					}
					for k, v := range method.FuncType.Params {
						x, _ := c.checkExpr(t.Args[k])
						if !IsAssignableTo(v.Type, x) {
							c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("cannot assign value of type %T to const of type %T", v.Type, x)})
							return TInvalid, nil
						}
					}
					if len(method.FuncType.Results) == 0 {
						return nil, nil
					}
					if len(method.FuncType.Results) > 1 {
						c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("too many arguments returned by %q func", method.Name)})
						return TInvalid, nil
					}
					return method.FuncType.Results[0].Type, nil
				}
			}
		}

		var ce CallExpr
		calleeType, calleeExpr := c.checkExpr(t.Callee)
		if named, ok := calleeType.(*NamedType); ok {
			if len(t.Args) != 1 {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("too many arguments in %#v, expected 1 got %d", named.Name, len(t.Args))})
				return TInvalid, nil
			}
			arg, _ := c.checkExpr(t.Args[0])
			if IsConvertibleTo(arg, named) {
				return named, nil
			}
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("cannot convert %#v to %s", arg, named.Name)})
			return TInvalid, nil
		}

		if builtin, ok := calleeType.(*BuiltinType); ok {
			if len(t.Args) != 1 {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("too many arguments in %#v, expected 1 got %d", expr, len(t.Args))})
				return TInvalid, nil
			}

			arg, ex := c.checkExpr(t.Args[0])
			if IsConvertibleTo(arg, builtin) {
				return calleeType, &ConversionExpr{To: calleeType, Value: ex}
			}
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("cannot convert %#v to %#v", arg, builtin)})
			return TInvalid, nil
		}

		fn, ok := calleeType.(*FuncMethod)
		if !ok {
			return TInvalid, nil
		}
		if len(t.Args) != len(fn.FuncType.Params) {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("too many arguments for %s func params, expected %d got %d", fn.Name, len(fn.FuncType.Params), len(t.Args))})
			return TInvalid, nil
		}
		for k, v := range fn.FuncType.Params {
			x, argExpr := c.checkExpr(t.Args[k])
			if !IsAssignableTo(v.Type, x) {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("cannot assign value of type %T to const of type %T", v.Type, x)})
				return TInvalid, nil
			}
			ce.Args = append(ce.Args, argExpr)
		}
		if len(fn.FuncType.Results) == 0 {
			return nil, nil
		}
		if len(fn.FuncType.Results) > 1 {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("too many arguments returned by %q func", fn.Name)})
			return TInvalid, nil
		}
		ce.Callee = calleeExpr
		ce.CalleeType = calleeType
		return fn.FuncType.Results[0].Type, &ce

	case *ast.ParenExpr:
		return c.checkExpr(t.Inner)

	case *ast.IndexExpr:
		baseType, _ := c.checkExpr(t.X)
		underlying := unwrapNamed(baseType)
		index, _ := c.checkExpr(t.Index)

		switch decl := underlying.(type) {
		case *SliceType:
			if !IsIdentical(index, TInt) {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid index expression of type %#v", index)})
				return TInvalid, nil
			}
			return decl.Elem, nil

		case *ArrayType:
			if !IsIdentical(index, TInt) {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid index expression of type %#v", index)})
				return TInvalid, nil
			}
			return decl.Elem, nil

		case *MapType:
			if !IsIdentical(decl.Key, index) {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid map index expression of type %#v", index)})
				return TInvalid, nil
			}
			return decl.Value, nil

		case *HashMapType:
			if !IsIdentical(decl.Key, index) {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid hashmap index expression of type %#v", index)})
				return TInvalid, nil
			}
			return decl.Value, nil

		default:
			return TInvalid, nil
		}

	case *ast.SelectorExpr:
		return c.checkSelectorExpr(t), nil

	default:
		return TInvalid, nil
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

	for _, fn := range c.methodDecls {
		c.checkMethodBody(fn)
	}
}

// checkFuncBody validates function body.
// An error is emitted if any
func (c *Checker) checkFuncBody(fn *ast.FuncDecl) {
	oldScope := c.scope
	oldFunc := c.currentFunc
	oldUseScope := c.useScope
	oldInComptimeFunc := c.inComptimeFunc
	defer func() {
		c.scope = oldScope
		c.currentFunc = oldFunc
		c.useScope = oldUseScope
		c.inComptimeFunc = oldInComptimeFunc
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
	c.inComptimeFunc = sym.IsComptime

	for _, p := range fnType.FuncType.Params {
		if !c.declareNoShadow(c.scope, &Symbol{
			Name:       p.Name,
			Kind:       SymVar,
			Type:       p.Type,
			IsComptime: c.inComptimeFunc,
		}, "variable") {
			continue
		}
	}

	for _, p := range fnType.FuncType.Results {
		if !c.declareNoShadow(c.scope, &Symbol{
			Name:       p.Name,
			Kind:       SymVar,
			Type:       p.Type,
			IsComptime: c.inComptimeFunc,
		}, "variable") {
			continue
		}
	}

	blockStmt := c.checkBlockStmt(fn.Body, nil)
	if len(fnType.FuncType.Results) > 0 && blockStmt.returnFlowResult == flowFallsThrough {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("missing return statement")})
		return
	}

	fd := &FuncDecl{
		Name:    sym.Name,
		Symbol:  resolvedSymbol(sym),
		Params:  fnType.FuncType.Params,
		Results: fnType.FuncType.Results,
	}

	if blockStmt.blockStmt != nil {
		fd.Body = &BlockStmt{blockStmt.blockStmt}
	}
	c.program.Files[c.programFileIndex].Decls = append(c.program.Files[c.programFileIndex].Decls, fd)
}

// checkMethodBody validates method.
// An error is emitted if any
func (c *Checker) checkMethodBody(fn *ast.FuncDecl) {
	oldScope := c.scope
	oldFunc := c.currentFunc
	oldUseScope := c.useScope
	defer func() {
		c.scope = oldScope
		c.currentFunc = oldFunc
		c.useScope = oldUseScope
	}()

	recvType := c.resolveType(fn.Receiver.Type)
	namedRcv, ok := recvType.(*NamedType)
	if !ok {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("receiver must be a named type got %#v", namedRcv)})
		return
	}

	method, ok := c.lookupMethodType(namedRcv, fn.Name.Value)
	if !ok {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("method type undefined")})
		return
	}

	c.currentFunc = method.FuncType
	c.scope = NewScope(c.pkgScope)
	c.useScope = true

	if !c.declareNoShadow(c.scope, &Symbol{
		Name: fn.Receiver.Name.Value,
		Kind: SymVar,
		Type: recvType,
	}, "receiver",
	) {
		return
	}

	for _, p := range method.FuncType.Params {
		if !c.declareNoShadow(c.scope, &Symbol{
			Name: p.Name,
			Kind: SymVar,
			Type: p.Type,
		}, "variable") {
			return
		}
	}

	for _, p := range method.FuncType.Results {
		if !c.declareNoShadow(c.scope, &Symbol{
			Name: p.Name,
			Kind: SymVar,
			Type: p.Type,
		}, "variable") {
			continue
		}
	}

	blockStmt := c.checkBlockStmt(fn.Body, nil)
	if len(method.FuncType.Results) > 0 && blockStmt.returnFlowResult == flowFallsThrough {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("missing return statement")})
		return
	}
}

// checkBlockStmt loops over block statements to check/declare them
// within its dedicated local scope
func (c *Checker) checkBlockStmt(block *ast.BlockStmt, returnInputVarsInitialized []string) (st stmtInfo) {
	if block == nil {
		st.returnFlowResult = flowFallsThrough
		st.returnedInputVarsInitialized = nil
		return
	}

	oldScope := c.scope
	defer func() {
		c.scope = oldScope
	}()

	blockScope := NewScope(c.scope)
	c.scope = blockScope

	var (
		cStmt stmtInfo
		bStmt []Stmt
	)
	flow := flowFallsThrough
	for i, stmt := range block.Stmts {
		if flow == flowReturns {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("unreachable code at %d:%d", stmt.Start().Line, stmt.End().Line)})
			return
		}

		if _, ok := stmt.(*ast.BreakStmt); ok {
			if i != len(block.Stmts)-1 {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("break must be the last statement of this block")})
				return
			}
		}

		if _, ok := stmt.(*ast.ContinueStmt); ok {
			if i != len(block.Stmts)-1 {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("continue must be the last statement of this block")})
				return
			}
		}

		// At first iteration, we need to hace block stmt return values
		// but after that we need to take result from statement iteration
		// and pass it down to the next one not to loose context
		if i == 0 {
			cStmt = c.checkStmt(stmt, returnInputVarsInitialized)
		} else {
			cStmt = c.checkStmt(stmt, cStmt.returnedInputVarsInitialized)
		}
		returnInputVarsInitialized = slices.Clone(cStmt.returnedInputVarsInitialized)
		bStmt = append(bStmt, cStmt.stmt)
		if flow == flowFallsThrough && cStmt.returnFlowResult == flowReturns {
			flow = flowReturns
		}
	}

	st.returnFlowResult = flow
	st.returnedInputVarsInitialized = slices.Clone(returnInputVarsInitialized)
	st.blockStmt = bStmt
	return
}

// checkStmt checks all possible statements kind.
// This is splitted here for reusability
func (c *Checker) checkStmt(stmt ast.Stmt, returnInputVarsInitialized []string) (st stmtInfo) {
	flow := flowFallsThrough

	switch t := stmt.(type) {
	case *ast.DeclStmt:
		switch decl := t.Decl.(type) {
		case *ast.ConstDecl:
			st.stmt = &DeclStmt{c.checkScopeConstDecl(decl)}
		case *ast.VarDecl:
			st.stmt = &DeclStmt{c.checkScopeVarDecl(decl)}

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
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("unsupported declaration %#v", decl)})
		}

	case *ast.AssignStmt:
		rvi, cStmt := c.checkAssigmentStmt(t, returnInputVarsInitialized)
		returnInputVarsInitialized = slices.Clone(rvi)
		st.stmt = cStmt

	case *ast.ReturnStmt:
		rflow, expr := c.checkReturnStmt(t, returnInputVarsInitialized)
		flow = rflow
		st.stmt = &ReturnStmt{expr}

	case *ast.IncDecStmt:
		c.checkIncDecStmt(t)

	case *ast.ExprStmt:
		c.checkExprStmt(t)

	case *ast.IfStmt:
		return c.checkIfStmt(t, returnInputVarsInitialized)

	case *ast.ForStmt:
		return c.checkForStmt(t, returnInputVarsInitialized)

	case *ast.RangeStmt:
		return c.checkRangeStmt(t, returnInputVarsInitialized)

	case *ast.SwitchStmt:
		return c.checkSwitchStmt(t, returnInputVarsInitialized)

	case *ast.FallThroughStmt:
		c.checkFallThroughStmt(t)

	case *ast.BreakStmt:
		returnInputVarsInitialized = c.checkBreakStmt(t, returnInputVarsInitialized)

	case *ast.ContinueStmt:
		c.checkContinueStmt(t)

	default:
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("unsupported statement %#v", stmt)})
	}

	st.returnFlowResult = flow
	st.returnedInputVarsInitialized = slices.Clone(returnInputVarsInitialized)
	return
}

// checkScopeConstDecl validates constant targetType and valueType.
// An error is emitted if any
func (c *Checker) checkScopeConstDecl(decl *ast.ConstDecl) Decl {
	targetType := c.checkTypeInCurrentMode(c.resolveType(decl.Type))
	valueType, expr := c.checkExprInCurrentMode(decl.Init)

	if !IsAssignableTo(targetType, valueType) {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("cannot assign value of type %T to var of type %T", valueType, targetType)})
		return nil
	}

	if !c.declareNoShadow(c.scope, &Symbol{
		Name:       decl.Name.Value,
		Kind:       SymConst,
		Type:       targetType,
		Decl:       decl,
		IsComptime: c.inComptimeFunc,
	}, "const") {
		return nil
	}

	sym := c.scope.Lookup(decl.Name.Value)
	return &ConstDecl{Name: sym.Name, Symbol: resolvedSymbol(sym), Eq: decl.Eq.Kind, Init: expr}
}

// checkScopeVarDecl validates constant targetType and valueType.
// An error is emitted if any
func (c *Checker) checkScopeVarDecl(decl *ast.VarDecl) Decl {
	targetType := c.checkTypeInCurrentMode(c.resolveType(decl.Type))
	valueType, expr := c.checkExprInCurrentMode(decl.Init)

	if !IsAssignableTo(targetType, valueType) {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("cannot assign value of type %T to var of type %T", valueType, targetType)})
		return nil
	}

	if !c.declareNoShadow(c.scope, &Symbol{
		Name:       decl.Name.Value,
		Kind:       SymVar,
		Type:       targetType,
		Decl:       decl,
		IsComptime: c.inComptimeFunc,
	}, "variable") {
		return nil
	}

	sym := c.scope.Lookup(decl.Name.Value)
	return &VarDecl{Name: sym.Name, Symbol: resolvedSymbol(sym), Eq: decl.Eq.Kind, Init: expr}
}

// checkAssignableExpr returns valid assignable expression.
// An error is emitted if any
func (c *Checker) checkAssignableExpr(expr ast.Expr) (Type, Expr) {
	switch t := expr.(type) {
	case *ast.IdentExpr:
		return c.checkExpr(t)

	case *ast.IndexExpr:
		return c.checkExpr(t)

	case *ast.SelectorExpr:
		return c.checkSelectorExpr(t), nil

	default:
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("unsupported expression %#v", expr)})
		return TInvalid, nil
	}
}

// checkSimpleAssignStmt validates simple assigment statements like x = 1 where x has already been defined.
// An error is emitted if any
func (c *Checker) checkSimpleAssignStmt(decl *ast.AssignStmt, returnInputVarsInitialized []string) []string {
	name := exprName(decl.Left)
	sym := c.scope.Lookup(name)
	if sym == nil {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("assigment %q is undefined", name)})
		return nil
	}

	if sym.Kind == SymConst {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("reassign const %q value is forbidden", name)})
		return nil
	}

	targetType, _ := c.checkAssignableExpr(decl.Left)
	valueType, _ := c.checkExprInCurrentMode(decl.Right)

	left, _ := c.checkExprInCurrentMode(decl.Left)
	if IsInvalid(left) {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid variable type %T", targetType)})
		return nil
	}

	right, _ := c.checkExprInCurrentMode(decl.Right)
	if IsInvalid(right) {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("cannot assign value of type %T to variable of type %T", valueType, targetType)})
		return nil
	}

	if !IsAssignableTo(targetType, valueType) {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("cannot assign value of type %T to variable of type %T", valueType, targetType)})
		return nil
	}

	sym.Type = targetType
	if c.currentFunc != nil && decl.Operator.Kind == token.Assign {
		for _, v := range c.currentFunc.Results {
			if v.Name != "" && v.Name == name {
				returnInputVarsInitialized = c.checkReturnVarsInitialized(returnInputVarsInitialized, name)
			}
		}
	}
	return returnInputVarsInitialized
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
// define assigment like x := 1 is forbidden as we cannot infer the value type. Is it an int? int32 etc?
// An error is emitted if any
func (c *Checker) checkDefineAssignStmt(decl *ast.AssignStmt) Stmt {
	valueType, expr := c.checkExprInCurrentMode(decl.Right)
	if IsInvalid(valueType) {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("expression %#v is invalid", decl.Right)})
		return nil
	}

	if isNumericExpr(decl.Right) {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("cannot use numeric only expression with define assigment declaration (:=)")})
		return nil
	}

	x, ok := decl.Left.(*ast.IdentExpr)
	if !ok {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("variable %#v not an identifier", decl)})
		return nil
	}

	sym := &Symbol{
		Name:       x.Name.Value,
		Kind:       SymVar,
		Type:       valueType,
		IsComptime: c.inComptimeFunc,
	}
	if !c.declareNoShadow(c.scope, sym, "variable") {
		return nil
	}

	return &DefineAssigmentStmt{
		Symbol: resolvedSymbol(sym),
		Right:  expr,
	}
}

// checkReturnStmt checks returned values statement types and length.
// An error is emitted if any
func (c *Checker) checkReturnStmt(decl *ast.ReturnStmt, returnInputVarsInitialized []string) (returnFlow, []Expr) {
	if decl == nil {
		return flowFallsThrough, nil
	}

	// This is needed when we have named returned values initialized and only use "return" keyword.
	// returnInputVarsInitialized have already been validated semanticly so
	// it's safe to compare lengths and return "flowReturns"
	if len(decl.Values) == 0 {
		if c.currentFunc != nil {
			for _, result := range c.currentFunc.Results {
				if result.Name == "" {
					c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("naked return requires named return values")})
					return flowFallsThrough, nil
				}

				if !slices.Contains(returnInputVarsInitialized, result.Name) {
					c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("returning uninitialized variable %q", result.Name)})
					return flowFallsThrough, nil
				}
			}
			return flowReturns, nil
		}
	}

	if len(c.currentFunc.Results) != len(decl.Values) {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("number of returned values is invalid, expected %d got %d", len(c.currentFunc.Results), len(decl.Values))})
		return flowFallsThrough, nil
	}

	for k, v := range decl.Values {
		expr, _ := c.checkExprInCurrentMode(v)
		if !IsIdentical(c.currentFunc.Results[k].Type, expr) {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("cannot use a value of type %T as %T in return statement", expr, c.currentFunc.Results[k].Type)})
			return flowFallsThrough, nil
		}
	}

	var expr []Expr
	for _, dv := range decl.Values {
		if x, ok := dv.(*ast.IdentExpr); ok {
			for _, r := range c.currentFunc.Results {
				if r.Name == x.Name.Value && !slices.Contains(returnInputVarsInitialized, x.Name.Value) {
					c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("returning uninitialized variable %q", x.Name.Value)})
				}
			}
		}

		_, e := c.checkExpr(dv)
		expr = append(expr, e)
	}

	return flowReturns, expr
}

// checkIncDecStmt validates increment/decrement statement.
// An error is emitted if any
func (c *Checker) checkIncDecStmt(decl *ast.IncDecStmt) {
	x, ok := decl.X.(*ast.IdentExpr)
	if !ok {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("variable %#v is not an identifier", decl)})
		return
	}

	sym := c.scope.Lookup(x.Name.Value)
	if sym == nil {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("assigment %q is undefined", x.Name.Value)})
		return
	}

	if sym.Kind == SymConst {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("const %q cannot be modified", x.Name.Value)})
		return
	}

	if !IsNumeric(sym.Type) {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("variable %q is a non-numeric type", x.Name.Value)})
		return
	}
}

// checkDefinedTypeDecl validates type declaration statement.
// An error is emitted if any
func (c *Checker) checkDefinedTypeDecl(decl *ast.DefinedTypeDecl) {
	if !c.declareNoShadow(c.scope, &Symbol{
		Name: decl.Name.Value,
		Kind: SymType,
		Type: &NamedType{
			Name:           decl.Name.Value,
			Decl:           decl,
			UnderlyingType: c.resolveType(decl.Type),
		},
	}, "type") {
		return
	}
}

// checkStructDecl validates struct statement.
// An error is emitted if any
func (c *Checker) checkStructDecl(decl *ast.StructDecl) {
	if !c.declareNoShadow(c.scope, &Symbol{
		Name: decl.Name.Value,
		Kind: SymType,
		Type: &NamedType{
			Name: decl.Name.Value,
			Decl: decl,
			UnderlyingType: &StructType{
				Decl:   decl,
				Fields: c.resolveStructFields(decl.Fields),
			},
		},
	}, "type") {
		return
	}
}

// checkEnumDecl validates struct statement.
// An error is emitted if any
func (c *Checker) checkEnumDecl(decl *ast.EnumDecl) {
	if !c.declareNoShadow(c.scope, &Symbol{
		Name: decl.Name.Value,
		Kind: SymType,
		Type: &NamedType{
			Name: decl.Name.Value,
			Decl: decl,
			UnderlyingType: &EnumType{
				Decl:     decl,
				Variants: c.resolveEnumVariants(decl.Variants),
			},
		},
	}, "type") {
		return
	}
}

// checkSumDecl validates struct statement.
// An error is emitted if any
func (c *Checker) checkSumDecl(decl *ast.SumDecl) {
	if !c.declareNoShadow(c.scope, &Symbol{
		Name: decl.Name.Value,
		Kind: SymType,
		Type: &NamedType{
			Name: decl.Name.Value,
			Decl: decl,
			UnderlyingType: &SumType{
				Decl:     decl,
				Variants: c.resolveSumVariants(decl.Variants),
			},
		},
	}, "type") {
		return
	}
}

// checkInterfaceDecl validates interface declaration statement.
// An error is emitted if any
func (c *Checker) checkInterfaceDecl(decl *ast.InterfaceDecl) {
	if !c.declareNoShadow(c.scope, &Symbol{
		Name: decl.Name.Value,
		Kind: SymType,
		Type: &NamedType{
			Name: decl.Name.Value,
			Decl: decl,
			UnderlyingType: &InterfaceType{
				Decl:    decl,
				Methods: c.resolveInterfaceMethods(decl.Methods),
			},
		},
	}, "type") {
		return
	}
}

// checkImplementsDecl validates "implements" declaration for interface.
// An error is emitted if any
func (c *Checker) checkImplementsDecl(decl *ast.ImplementsDecl) {
	sym := c.pkgScope.Lookup(decl.TypeName.Value)
	if sym == nil || sym.Kind != SymType {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("type %q is undefined", decl.TypeName.Value)})
		return
	}

	implementer, ok := sym.Type.(*NamedType)
	if !ok {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("%q is not a named type", decl.TypeName.Value)})
		return
	}

	if _, ok := unwrapNamed(implementer).(*InterfaceType); ok {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("%q cannot implement another interface", decl.TypeName.Value)})
		return
	}

	ifaceType := c.resolveType(decl.Interface)
	if IsInvalid(ifaceType) {
		return
	}

	ifaceNamed, ok := ifaceType.(*NamedType)
	if !ok {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("implements target must be a named interface")})
		return
	}

	iface, ok := unwrapNamed(ifaceNamed).(*InterfaceType)
	if !ok {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("type %q is not an interface", ifaceNamed.Name)})
		return
	}

	if !c.implementsInterface(implementer, iface) {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("%q interface implementation is invalid", decl.TypeName.Value)})
		return
	}

	c.implInfos = append(c.implInfos, ImplInfo{
		TypeName:      decl.TypeName.Value,
		Type:          implementer,
		InterfaceName: ifaceNamed.Name,
		Interface:     ifaceNamed,
		Decl:          decl,
	})
}

// implementsInterface loops over interface methods to verify if all methods
// have been implemented with the right signatures
func (c *Checker) implementsInterface(named *NamedType, iface *InterfaceType) bool {
	for _, im := range iface.Methods {
		fn, ok := c.lookupMethodType(named, im.Name)
		if !ok {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("method %q not implemented", im.Name)})
			return false
		}

		if !IsIdentical(im.FuncType, fn.FuncType) {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("wrong method %q signature", im.Name)})
			return false
		}
	}
	return true
}

// lookupMethodType loops over methods to match provided named type and func name.
// It returns its semantic type and true when found
func (c *Checker) lookupMethodType(named *NamedType, name string) (*FuncMethod, bool) {
	if c.methods == nil {
		return nil, false
	}

	if method, ok := c.methods[named]; ok {
		if m, ok := method[name]; ok {
			return m, true
		}
	}
	return nil, false
}

// checkExprStmt validates expression statement.
// An error is emitted if any
func (c *Checker) checkExprStmt(stmt *ast.ExprStmt) {
	call, ok := stmt.Expr.(*ast.CallExpr)
	if !ok {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("call expression statement must be a function call, got %#v", call)})
		return
	}

	if _, ok := call.Callee.(*ast.SelectorExpr); ok {
		_, _ = c.checkExprInCurrentMode(stmt.Expr)
		return
	}

	calleType, _ := c.checkExprInCurrentMode(call.Callee)
	fn, ok := calleType.(*FuncMethod)
	if !ok {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("callee type expression statement must be a function call got %#v", calleType)})
		return
	}

	if len(fn.FuncType.Results) > 0 {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("calling function with returned values are forbidden without assignment, expected 0, got %d", len(fn.FuncType.Results))})
		return
	}

	_, _ = c.checkExprInCurrentMode(stmt.Expr)
}

// checkSelectorExpr validates selector expression and return its type.
// An error is emitted if any
func (c *Checker) checkSelectorExpr(expr *ast.SelectorExpr) Type {
	baseType, _ := c.checkExpr(expr.X)
	underlying := unwrapNamed(baseType)

	switch t := underlying.(type) {
	case *StructType:
		tp, ok := lookupStructField(t, expr.Selector.Value)
		if !ok {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("unknown field %q", expr.Selector.Value)})
			return TInvalid
		}
		return tp

	case *InterfaceType:
		tp, ok := lookupInterfaceMethods(t, expr.Selector.Value)
		if !ok {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("unknown method %q", expr.Selector.Value)})
			return TInvalid
		}
		return tp

	default:
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid type")})
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

// checkIfStmt validates if statement block
func (c *Checker) checkIfStmt(stmt *ast.IfStmt, returnInputVarsInitialized []string) (st stmtInfo) {
	if stmt == nil {
		st.returnFlowResult = flowFallsThrough
		st.returnedInputVarsInitialized = slices.Clone(returnInputVarsInitialized)
		return
	}

	condType, _ := c.checkExprInCurrentMode(stmt.Condition)
	if !IsBool(condType) {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("if condition must returned a boolean")})
		st.returnFlowResult = flowFallsThrough
		st.returnedInputVarsInitialized = slices.Clone(returnInputVarsInitialized)
		return
	}

	var (
		initializedIfThen, initializedIfElse bool
		rvi                                  []string
		thenStmt, elseStmt                   stmtInfo
	)

	copyIf := slices.Clone(returnInputVarsInitialized)
	copyElse := slices.Clone(returnInputVarsInitialized)
	if stmt.Then != nil {
		initializedIfThen = true
		thenStmt = c.checkBlockStmt(stmt.Then, returnInputVarsInitialized)
	}

	if stmt.Else != nil {
		initializedIfElse = true
		switch t := stmt.Else.(type) {
		case *ast.BlockStmt:
			elseStmt = c.checkBlockStmt(t, returnInputVarsInitialized)
		default:
			elseStmt = c.checkStmt(t, returnInputVarsInitialized)
		}
	}

	if initializedIfThen && initializedIfElse {
		if thenStmt.returnFlowResult == flowFallsThrough && elseStmt.returnFlowResult == flowFallsThrough {
			var x, y []string
			for _, v := range thenStmt.returnedInputVarsInitialized {
				if !slices.Contains(copyIf, v) {
					x = append(x, v)
				}
			}

			for _, v := range elseStmt.returnedInputVarsInitialized {
				if !slices.Contains(copyElse, v) {
					y = append(y, v)
				}
			}

			for _, v := range x {
				if slices.Contains(y, v) {
					rvi = append(rvi, v)
				}
			}

			for _, v := range copyIf {
				if !slices.Contains(rvi, v) {
					rvi = append(rvi, v)
				}
			}

			if len(x) == 0 {
				rvi = slices.Clone(y)
			}

			if len(x) == 0 && len(y) == 0 {
				rvi = slices.Clone(copyIf)
			}

			if len(x) == 0 && len(y) > 0 {
				rvi = nil
			}

			st.returnFlowResult = flowFallsThrough
			st.returnedInputVarsInitialized = rvi
			return
		}

		if thenStmt.returnFlowResult == flowReturns && elseStmt.returnFlowResult == flowFallsThrough {
			st.returnFlowResult = flowFallsThrough
			st.returnedInputVarsInitialized = slices.Clone(elseStmt.returnedInputVarsInitialized)
			return
		}

		if thenStmt.returnFlowResult == flowFallsThrough && elseStmt.returnFlowResult == flowReturns {
			st.returnFlowResult = flowFallsThrough
			st.returnedInputVarsInitialized = slices.Clone(thenStmt.returnedInputVarsInitialized)
			return
		}

		if thenStmt.returnFlowResult == flowReturns && elseStmt.returnFlowResult == flowReturns {
			st.returnFlowResult = flowReturns
			st.returnedInputVarsInitialized = slices.Clone(returnInputVarsInitialized)
			return
		}
	}

	if initializedIfThen && !initializedIfElse {
		rvi = slices.Clone(copyIf)
	}

	st.returnFlowResult = flowFallsThrough
	st.returnedInputVarsInitialized = slices.Clone(rvi)
	return
}

// checkForStmt validates plain for statement block
func (c *Checker) checkForStmt(stmt *ast.ForStmt, returnInputVarsInitialized []string) (st stmtInfo) {
	if stmt == nil {
		st.returnFlowResult = flowFallsThrough
		st.returnedInputVarsInitialized = slices.Clone(returnInputVarsInitialized)
		return
	}

	oldLoopDepth := c.loopDepth
	c.loopDepth++
	oldscope := c.scope
	oldBreakFound := c.breakFound
	oldBreakInputVarsInitialized := slices.Clone(c.breakInputVarsInitialized)
	defer func() {
		c.loopDepth = oldLoopDepth
		c.scope = oldscope
		c.breakFound = oldBreakFound
		c.breakInputVarsInitialized = oldBreakInputVarsInitialized
	}()

	c.scope = NewScope(c.scope)
	c.breakFound = false

	var initStmt, bodyStmt stmtInfo
	afterInit := slices.Clone(returnInputVarsInitialized)
	if stmt.Init != nil {
		initStmt = c.checkStmt(stmt.Init, returnInputVarsInitialized)
		afterInit = slices.Clone(initStmt.returnedInputVarsInitialized)
	}

	if stmt.Condition != nil {
		condType, _ := c.checkExprInCurrentMode(stmt.Condition)
		if !IsBool(condType) {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("for condition must return a boolean")})
			st.returnFlowResult = flowFallsThrough
			st.returnedInputVarsInitialized = slices.Clone(returnInputVarsInitialized)
			return
		}
	}

	if stmt.Post != nil {
		c.checkStmt(stmt.Post, nil)
	}

	if stmt.Body != nil {
		if stmt.Condition == nil && len(stmt.Body.Stmts) == 0 {
			st.returnFlowResult = flowFallsThrough
			st.returnedInputVarsInitialized = slices.Clone(returnInputVarsInitialized)
			return
		}

		bodyStmt = c.checkBlockStmt(stmt.Body, afterInit)
		if stmt.Condition == nil {
			if c.breakFound {
				st.returnFlowResult = flowFallsThrough
				st.returnedInputVarsInitialized = slices.Clone(c.breakInputVarsInitialized)
				return
			}

			st.returnFlowResult = bodyStmt.returnFlowResult
			st.returnedInputVarsInitialized = slices.Clone(returnInputVarsInitialized)
			return
		}
	}

	if stmt.Condition != nil || c.breakFound {
		st.returnFlowResult = flowFallsThrough
		st.returnedInputVarsInitialized = slices.Clone(afterInit)
		return
	}

	st.returnFlowResult = flowReturns
	st.returnedInputVarsInitialized = slices.Clone(afterInit)
	return
}

// checkForStmt validates for assigment statement
func (c *Checker) checkAssigmentStmt(stmt *ast.AssignStmt, returnInputVarsInitialized []string) ([]string, Stmt) {
	var s Stmt
	switch stmt.Operator.Kind {
	case token.Assign, token.PlusEq, token.MinusEq:
		returnInputVarsInitialized = c.checkSimpleAssignStmt(stmt, returnInputVarsInitialized)
	case token.Define:
		s = c.checkDefineAssignStmt(stmt)
	default:
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("unsupported assigment in for statement %#v", stmt)})
	}
	return returnInputVarsInitialized, s
}

// checkRangeStmt validates range statement block
func (c *Checker) checkRangeStmt(stmt *ast.RangeStmt, returnInputVarsInitialized []string) (st stmtInfo) {
	st.returnFlowResult = flowFallsThrough
	st.returnedInputVarsInitialized = slices.Clone(returnInputVarsInitialized)
	if stmt == nil {
		return
	}

	oldLoopDepth := c.loopDepth
	c.loopDepth++
	oldscope := c.scope
	oldBreakFound := c.breakFound
	oldBreakInputVarsInitialized := slices.Clone(c.breakInputVarsInitialized)
	defer func() {
		c.loopDepth = oldLoopDepth
		c.scope = oldscope
		c.breakFound = oldBreakFound
		c.breakInputVarsInitialized = oldBreakInputVarsInitialized
	}()

	c.scope = NewScope(c.scope)
	c.breakFound = false

	iteratorType, _ := c.checkExprInCurrentMode(stmt.X)
	if IsInvalid(iteratorType) {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("range expression is invalid")})
		return
	}

	rangekeyType, rangeValueType, ok := rangeVars(iteratorType)
	if !ok {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("unsupported range var type")})
		return
	}

	// for range x {}
	if stmt.Op == (token.Token{}) && stmt.Body != nil {
		// Range block must always be a flowFallsThrough as it might not be executed
		c.checkBlockStmt(stmt.Body, returnInputVarsInitialized)
		return
	}

	switch stmt.Op.Kind {
	case token.Assign:
		if stmt.Key != nil && stmt.Value != nil {
			name := exprName(stmt.Key)
			sym := c.scope.Lookup(name)
			if sym == nil {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("assigment %q is undefined", name)})
				return
			}

			if sym.Kind == SymConst {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("reassign const %q value is forbidden", name)})
				return
			}

			key, _ := c.checkExpr(stmt.Key)
			if stmt.Key.Name.Value != "_" && !IsAssignableTo(rangekeyType, key) {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid range key type, expected %#v, got %#v", rangekeyType, key)})
				return
			}

			name = exprName(stmt.Value)
			sym = c.scope.Lookup(name)
			if sym == nil {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("assigment %q is undefined", name)})
				return
			}

			if sym.Kind == SymConst {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("reassign const %q value is forbidden", name)})
				return
			}

			value, _ := c.checkExpr(stmt.Value)
			if stmt.Value.Name.Value != "_" && !IsAssignableTo(rangeValueType, value) {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid range value type, expected %#v, got %#v", rangeValueType, value)})
				return
			}

			if c.currentFunc != nil {
				for _, v := range c.currentFunc.Results {
					if v.Name != "" && v.Name == stmt.Key.Name.Value {
						returnInputVarsInitialized = c.checkReturnVarsInitialized(returnInputVarsInitialized, stmt.Key.Name.Value)
					}
					if v.Name != "" && v.Name == stmt.Value.Name.Value {
						returnInputVarsInitialized = c.checkReturnVarsInitialized(returnInputVarsInitialized, stmt.Value.Name.Value)
					}
				}
			}
		} else if stmt.Key != nil {
			if stmt.Key.Name.Value == "_" {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("blank identifier for this range key is forbidden")})
				return
			}

			name := exprName(stmt.Key)
			sym := c.scope.Lookup(name)
			if sym == nil {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("assigment %q is undefined", name)})
				return
			}

			if sym.Kind == SymConst {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("reassign const %q value is forbidden", name)})
				return
			}

			key, _ := c.checkExpr(stmt.Key)
			if !IsAssignableTo(rangekeyType, key) {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid range key type, expected %#v, got %#v", rangekeyType, key)})
				return
			}

			if c.currentFunc != nil {
				for _, v := range c.currentFunc.Results {
					if v.Name != "" && v.Name == stmt.Key.Name.Value {
						returnInputVarsInitialized = c.checkReturnVarsInitialized(returnInputVarsInitialized, stmt.Key.Name.Value)
					}
				}
			}
		}

	case token.Define:
		if stmt.Key != nil && stmt.Value != nil {
			if stmt.Key.Name.Value == "_" && stmt.Value.Name.Value == "_" {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("range key and value cannot be both blank identifiers")})
				return
			}

			if !c.declareNoShadow(c.scope, &Symbol{
				Name: stmt.Key.Name.Value,
				Kind: SymVar,
				Type: rangekeyType,
			}, "variable") {
				return
			}

			if !c.declareNoShadow(c.scope, &Symbol{
				Name: stmt.Value.Name.Value,
				Kind: SymVar,
				Type: rangeValueType,
			}, "variable") {
				return
			}
		} else if stmt.Key != nil {
			if stmt.Key.Name.Value == "_" {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("blank identifier for this range key is forbidden")})
				return
			}

			if !c.declareNoShadow(c.scope, &Symbol{
				Name: stmt.Key.Name.Value,
				Kind: SymVar,
				Type: rangekeyType,
			}, "variable") {
				return
			}
		}

	default:
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("forbidden range token %q", stmt.Op.Value)})
		return
	}

	if stmt.Body != nil {
		// Range block must always be a flowFallsThrough as it might not be executed
		c.checkBlockStmt(stmt.Body, returnInputVarsInitialized)
	}
	return
}

// rangeVars returns underlying type var
func rangeVars(t Type) (Type, Type, bool) {
	un := unwrapNamed(t)
	switch x := un.(type) {
	case *ArrayType:
		return TInt, x.Elem, true
	case *SliceType:
		return TInt, x.Elem, true
	case *MapType:
		return x.Key, x.Value, true
	case *HashMapType:
		return x.Key, x.Value, true
	case *BuiltinType:
		switch x {
		case TInt, TInt8, TInt32, TInt64, TUInt, TUInt8, TUInt32, TUInt64:
			return t, t, true
		default:
			return nil, nil, false
		}
	default:
		return nil, nil, false
	}
}

// checkSwitchStmt validates switch statement block
func (c *Checker) checkSwitchStmt(stmt *ast.SwitchStmt, returnInputVarsInitialized []string) (st stmtInfo) {
	st.returnFlowResult = flowFallsThrough
	st.returnedInputVarsInitialized = slices.Clone(returnInputVarsInitialized)
	if stmt == nil {
		return
	}

	oldScope := c.scope
	defer func() {
		c.scope = oldScope
	}()

	c.scope = NewScope(c.scope)

	if stmt.Init != nil {
		// switch examples:
		// switch z:=w();z {
		// switch z=w();z {
		// the init is the first z
		cStmt := c.checkStmt(stmt.Init, returnInputVarsInitialized)

		if val, ok := stmt.Init.(*ast.AssignStmt); ok && val.Operator.Kind == token.Assign {
			name := exprName(val.Left)
			sym := c.scope.Lookup(name)
			if sym != nil && sym.Kind == SymConst {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("reassign const %q value is forbidden", name)})
				return
			}

			for _, v := range c.currentFunc.Results {
				if v.Name != "" && v.Name == name {
					returnInputVarsInitialized = c.checkReturnVarsInitialized(cStmt.returnedInputVarsInitialized, name)
				}
			}
		}
	}

	var (
		dcount                  int
		hasFallThrough          bool
		caseIsEmpty             bool
		caseHasNoReturn         bool
		caseHasExitedNormally   bool
		defaultReturnFlowResult returnFlow
	)

	intersection := slices.Clone(returnInputVarsInitialized)
	intersectionFilter := func(sStmt stmtInfo) {
		if sStmt.returnFlowResult == flowFallsThrough && !sStmt.switchCaseHasFallThrough {
			if !caseHasExitedNormally {
				caseHasExitedNormally = true
				intersection = slices.Clone(sStmt.returnedInputVarsInitialized)
			} else {
				var tmp []string
				for _, vi := range intersection {
					for _, rvi := range sStmt.returnedInputVarsInitialized {
						if vi == rvi {
							tmp = append(tmp, rvi)
						}
					}
				}
				intersection = slices.Clone(tmp)
			}
		}
	}

	if stmt.Tag != nil {
		// switch examples:
		// switch a {
		// switch z:=w();z {
		// switch z=w();z {
		// the tag is "a" or the last z
		tagType, _ := c.checkExprInCurrentMode(stmt.Tag)
		if IsInvalid(tagType) {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("tag expression is invalid")})
			return
		}

		underlying := unwrapNamed(tagType)
		switch underlying.(type) {
		case *SumType:
			return c.checkSwitchStmtSumType(tagType, stmt, returnInputVarsInitialized)
		case *EnumType:
			return c.checkSwitchStmtEnumType(tagType, stmt, returnInputVarsInitialized)
		}

		if !IsComparable(tagType) {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("tag type is not comparable got %#v", tagType)})
			return
		}

		if len(stmt.Cases) == 0 {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("switch statement has 0 cases")})
			return
		}

		seen := make(map[constKey]bool, len(stmt.Cases))
		for i, cc := range stmt.Cases {
			if cc.Case.Kind == token.KWDefault {
				dcount++

				if dcount > 1 {
					c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("multiple default clauses are forbidden, got %d", dcount)})
					return
				}
			}

			for _, v := range cc.Values {
				vExpr, _ := c.checkExprInCurrentMode(v)
				if !IsIdentical(tagType, vExpr) {
					c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("tag and case are not identical expected %#v got %#v", tagType, vExpr)})
					return
				}

				if ck, ok := c.constKey(v, vExpr); ok {
					if seen[ck] {
						c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("duplicate case expression %s at %d:%d", ck.value, v.Start().Line, v.End().Line)})
						return
					}
					seen[ck] = true
				}
			}

			if cc.Body != nil {
				sStmt := c.checkSwitchBody(cc.Body, i == len(stmt.Cases)-1, returnInputVarsInitialized)

				if sStmt.switchCaseHasFallThrough {
					hasFallThrough = true
					continue
				}

				intersectionFilter(sStmt)

				if dcount == 1 {
					defaultReturnFlowResult = sStmt.returnFlowResult

					if hasFallThrough {
						hasFallThrough = false
						if !caseHasNoReturn && sStmt.returnFlowResult == flowReturns {
							continue
						}
					}

					if sStmt.returnFlowResult == flowReturns {
						continue
					}

					caseHasNoReturn = true
					if len(sStmt.returnedInputVarsInitialized) == 0 {
						caseIsEmpty = true
					}
					continue
				}

				if hasFallThrough {
					hasFallThrough = false
					if !caseHasNoReturn && sStmt.returnFlowResult == flowReturns {
						continue
					}
				}

				if !caseHasNoReturn && sStmt.returnFlowResult == flowReturns {
					continue
				}

				caseHasNoReturn = true
				if len(sStmt.returnedInputVarsInitialized) == 0 {
					caseIsEmpty = true
				}
			} else {
				caseIsEmpty = true
				caseHasExitedNormally = true
				intersection = slices.Clone(returnInputVarsInitialized)
			}
		}
	} else {
		// example switch
		/*
			switch {
			case a == 1:
				x := int(1)
			}
		*/
		if len(stmt.Cases) == 0 {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("switch statement has 0 cases")})
			return
		}

		for i, cc := range stmt.Cases {
			if cc.Case.Kind == token.KWDefault {
				dcount++

				if dcount > 1 {
					c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("multiple default clauses are forbidden, got %d", dcount)})
					return
				}
			}

			for _, v := range cc.Values {
				vExpr, _ := c.checkExprInCurrentMode(v)
				// TODO: Tagless switch duplicates is a job for the linter
				// as it's difficult for the checker to properly handle every cases
				// without any burden

				if !IsBool(vExpr) {
					c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("switch case expected boolean got %#v", vExpr)})
					return
				}
			}

			if cc.Body != nil {
				sStmt := c.checkSwitchBody(cc.Body, i == len(stmt.Cases)-1, returnInputVarsInitialized)

				if sStmt.switchCaseHasFallThrough {
					hasFallThrough = true
					continue
				}

				intersectionFilter(sStmt)

				if dcount == 1 {
					defaultReturnFlowResult = sStmt.returnFlowResult

					if hasFallThrough {
						hasFallThrough = false
						if !caseHasNoReturn && sStmt.returnFlowResult == flowReturns {
							continue
						}
					}

					if sStmt.returnFlowResult == flowReturns {
						continue
					}

					caseHasNoReturn = true
					if len(sStmt.returnedInputVarsInitialized) == 0 {
						caseIsEmpty = true
					}
					continue
				}

				if hasFallThrough {
					hasFallThrough = false
					if !caseHasNoReturn && sStmt.returnFlowResult == flowReturns {
						continue
					}
				}

				if !caseHasNoReturn && sStmt.returnFlowResult == flowReturns {
					continue
				}

				caseHasNoReturn = true
				if len(sStmt.returnedInputVarsInitialized) == 0 {
					caseIsEmpty = true
				}
			} else {
				caseIsEmpty = true
				caseHasExitedNormally = true
				intersection = slices.Clone(returnInputVarsInitialized)
			}
		}
	}

	st.returnedInputVarsInitialized = slices.Clone(intersection)
	if dcount == 0 {
		st.returnedInputVarsInitialized = slices.Clone(returnInputVarsInitialized)
	}

	if caseHasNoReturn || caseIsEmpty {
		st.returnFlowResult = flowFallsThrough
		return
	}

	if defaultReturnFlowResult == flowReturns {
		st.returnFlowResult = flowReturns
		return
	}

	st.returnFlowResult = flowFallsThrough
	return
}

// constKey returns true if expression is a literal one.
// It will be used by basic switch case to find duplicates
func (c *Checker) constKey(expr ast.Expr, typ Type) (constKey, bool) {
	switch e := expr.(type) {
	case *ast.IntLitExpr:
		return constKey{
			typeID: typ.String(),
			kind:   constInt,
			value:  e.Name.Value,
		}, true

	case *ast.FloatLitExpr:
		return constKey{
			typeID: typ.String(),
			kind:   constFloat,
			value:  e.Name.Value,
		}, true

	case *ast.BoolLitExpr:
		return constKey{
			typeID: typ.String(),
			kind:   constBool,
			value:  e.Name.Value,
		}, true

	case *ast.StringLitExpr:
		return constKey{
			typeID: typ.String(),
			kind:   constString,
			value:  e.Name.Value,
		}, true
	}
	return constKey{}, false
}

// checkSwitchBody loops over switch base body for validation
func (c *Checker) checkSwitchBody(body []ast.Stmt, isLastCaseClause bool, returnInputVarsInitialized []string) (st stmtInfo) {
	oldScope := c.scope
	oldInSwitchCase := c.inSwitchCase
	defer func() {
		c.scope = oldScope
		c.inSwitchCase = oldInSwitchCase
	}()
	c.inSwitchCase = true

	bodyScope := NewScope(c.scope)
	c.scope = bodyScope

	var (
		cStmt            stmtInfo
		returnFlowResult returnFlow
	)

	cStmt.returnedInputVarsInitialized = slices.Clone(returnInputVarsInitialized)
	for i, b := range body {
		if returnFlowResult == flowReturns {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("unreachable code at %d:%d", b.Start().Line, b.End().Line)})
			return
		}

		if _, ok := b.(*ast.FallThroughStmt); ok {
			if i != len(body)-1 {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("fallthrough must be the last statement of the switch case body")})
				return
			}

			if isLastCaseClause {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("fallthrough is forbidden inside last switch case")})
				return
			}
			cStmt.switchCaseHasFallThrough = true
			continue
		}

		cStmt = c.checkStmt(b, cStmt.returnedInputVarsInitialized)
		returnFlowResult = cStmt.returnFlowResult
	}

	return cStmt
}

// checkFallThroughStmt produces an error when not into switch case
func (c *Checker) checkFallThroughStmt(_ *ast.FallThroughStmt) {
	if !c.inSwitchCase {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("fallthrough is forbidden outside of switch case")})
	}
}

// checkBreakStmt produces an error when not into for loop statement
func (c *Checker) checkBreakStmt(_ *ast.BreakStmt, returnInputVarsInitialized []string) []string {
	if c.loopDepth == 0 {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("break is forbidden outside of a loop")})
		return nil
	}

	if !c.breakFound {
		c.breakFound = true
		c.breakInputVarsInitialized = slices.Clone(returnInputVarsInitialized)
		return returnInputVarsInitialized
	}

	var intersection []string
	for _, v := range returnInputVarsInitialized {
		if slices.Contains(c.breakInputVarsInitialized, v) {
			intersection = append(intersection, v)
		}
	}

	c.breakInputVarsInitialized = slices.Clone(intersection)
	return returnInputVarsInitialized
}

// checkContinueStmt produces an error when not into for loop statement
func (c *Checker) checkContinueStmt(_ *ast.ContinueStmt) {
	if c.loopDepth == 0 {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("continue is forbidden outside of a loop")})
	}
}

// checkSwitchStmtSumType validates switch statement block
func (c *Checker) checkSwitchStmtSumType(tagType Type, stmt *ast.SwitchStmt, returnInputVarsInitialized []string) (st stmtInfo) {
	underlying := unwrapNamed(tagType)
	sm := underlying.(*SumType)

	var (
		cStmt                 stmtInfo
		caseHasExitedNormally bool
		caseHasNoReturn       bool
		caseIsEmpty           bool
	)

	seen := make(map[string]bool, len(stmt.Cases))

	intersection := slices.Clone(returnInputVarsInitialized)
	intersectionFilter := func(sStmt stmtInfo) {
		if sStmt.returnFlowResult == flowFallsThrough {
			if !caseHasExitedNormally {
				caseHasExitedNormally = true
				intersection = slices.Clone(sStmt.returnedInputVarsInitialized)
			} else {
				var tmp []string
				for _, vi := range intersection {
					for _, rvi := range sStmt.returnedInputVarsInitialized {
						if vi == rvi {
							tmp = append(tmp, rvi)
						}
					}
				}
				intersection = slices.Clone(tmp)
			}
		}
	}

	for _, cc := range stmt.Cases {
		if cc.Case.Kind == token.KWDefault {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("default is forbidden inside sum type switch")})
			return
		}

		if len(cc.Values) != 1 {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("sum case must have exactly one variant")})
			return
		}

		call, ok := cc.Values[0].(*ast.CallExpr)
		if !ok {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("sum case must be a variant, got %#v", call)})
			return
		}

		callee, ok := call.Callee.(*ast.IdentExpr)
		if !ok {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("sum variant case must use a variant identifier, got %#v", callee)})
			return
		}

		variantName, ok := fetchSumVariant(callee.Name.Value, sm)
		if !ok {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("unknown variant name %q", callee.Name.Value)})
			return
		}

		if seen[variantName.Name] {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("duplicate variant name %q", variantName.Name)})
			return
		}
		seen[variantName.Name] = true

		if len(variantName.Field) != len(call.Args) {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("variant arguments length invalid, expected %d got %d", len(variantName.Field), len(call.Args))})
			return
		}

		var bindings []*Symbol
		seenBindings := make(map[string]bool)
		for k, v := range call.Args {
			ident, ok := v.(*ast.IdentExpr)
			if !ok {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("variant expected identifier got %#v", v)})
				return
			}

			// TODO: if we need _ we will have to update the parser
			// if ident.Name.Value == "_" {
			// 	continue
			// }

			if seenBindings[ident.Name.Value] {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("variable %q already declared", ident.Name.Value)})
				return
			}
			seenBindings[ident.Name.Value] = true

			if c.scope.Lookup(ident.Name.Value) != nil {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("variable %q already declared", ident.Name.Value)})
				return
			}

			bindings = append(bindings, &Symbol{
				Name: ident.Name.Value,
				Kind: SymVar,
				Type: variantName.Field[k].Type,
			})
		}

		if cc.Body != nil {
			cStmt = c.checkSwitchSumBody(cc.Body, bindings, returnInputVarsInitialized)

			intersectionFilter(cStmt)
			if cStmt.returnFlowResult == flowReturns {
				continue
			}

			caseHasNoReturn = true
			if len(cStmt.returnedInputVarsInitialized) == 0 {
				caseIsEmpty = true
			}
		} else {
			caseIsEmpty = true
			caseHasExitedNormally = true
			intersection = slices.Clone(returnInputVarsInitialized)
		}
	}

	for _, v := range sm.Variants {
		if !seen[v.Name] {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("missing variant case %q", v.Name)})
			return
		}
	}

	st.returnedInputVarsInitialized = slices.Clone(intersection)
	if caseHasNoReturn || caseIsEmpty {
		st.returnFlowResult = flowFallsThrough
		return
	}

	st.returnFlowResult = flowReturns
	st.returnedInputVarsInitialized = slices.Clone(returnInputVarsInitialized)
	return
}

// fetchSumVariant fetches variant from sum type when found
func fetchSumVariant(name string, sm *SumType) (SumVariant, bool) {
	for _, v := range sm.Variants {
		if v.Name == name {
			return v, true
		}
	}

	return SumVariant{}, false
}

// checkSwitchSumBody loops over sum type switch base body for validation
// and creates related statements
func (c *Checker) checkSwitchSumBody(body []ast.Stmt, bindings []*Symbol, returnInputVarsInitialized []string) (st stmtInfo) {
	oldScope := c.scope
	oldInSwitchCase := c.inSwitchCase
	defer func() {
		c.scope = oldScope
		c.inSwitchCase = oldInSwitchCase
	}()
	c.inSwitchCase = true

	c.scope = NewScope(c.scope)

	for _, b := range bindings {
		if !c.declareNoShadow(c.scope, b, "variable") {
			return
		}
	}

	var (
		cStmt            stmtInfo
		returnFlowResult returnFlow
	)

	cStmt.returnedInputVarsInitialized = slices.Clone(returnInputVarsInitialized)
	for _, b := range body {
		if returnFlowResult == flowReturns {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("unreachable code at %d:%d", b.Start().Line, b.End().Line)})
			return
		}

		if _, ok := b.(*ast.FallThroughStmt); ok {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("fallthrough is forbidden in sum switch type")})
			return
		}

		cStmt = c.checkStmt(b, cStmt.returnedInputVarsInitialized)
		returnFlowResult = cStmt.returnFlowResult
	}

	return cStmt
}

// checkSwitchStmtEnumType validates switch statement block
func (c *Checker) checkSwitchStmtEnumType(tagType Type, stmt *ast.SwitchStmt, returnInputVarsInitialized []string) (st stmtInfo) {
	underlying := unwrapNamed(tagType)
	en := underlying.(*EnumType)

	var (
		cStmt                 stmtInfo
		caseHasExitedNormally bool
		caseHasNoReturn       bool
		caseIsEmpty           bool
	)
	seen := make(map[string]bool, len(stmt.Cases))

	intersection := slices.Clone(returnInputVarsInitialized)
	intersectionFilter := func(sStmt stmtInfo) {
		if sStmt.returnFlowResult == flowFallsThrough {
			if !caseHasExitedNormally {
				caseHasExitedNormally = true
				intersection = slices.Clone(sStmt.returnedInputVarsInitialized)
			} else {
				var tmp []string
				for _, vi := range intersection {
					for _, rvi := range sStmt.returnedInputVarsInitialized {
						if vi == rvi {
							tmp = append(tmp, rvi)
						}
					}
				}
				intersection = slices.Clone(tmp)
			}
		}
	}

	for _, cc := range stmt.Cases {
		if cc.Case.Kind == token.KWDefault {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("default is forbidden inside enum type switch")})
			return
		}

		if len(cc.Values) != 1 {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("enum case must have exactly one variant")})
			return
		}

		ident, ok := cc.Values[0].(*ast.IdentExpr)
		if !ok {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("enum case must be a variant, got %#v", ident)})
			return
		}

		variantName, ok := fetchEnumVariant(ident.Name.Value, en)
		if !ok {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("unknown variant name %q", ident.Name.Value)})
			return
		}

		if seen[variantName] {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("duplicate variant name %q", variantName)})
			return
		}
		seen[variantName] = true

		if cc.Body != nil {
			cStmt = c.checkSwitchEnumBody(cc.Body, returnInputVarsInitialized)

			intersectionFilter(cStmt)
			if cStmt.returnFlowResult == flowReturns {
				continue
			}

			caseHasNoReturn = true
			if len(cStmt.returnedInputVarsInitialized) == 0 {
				caseIsEmpty = true
			}
		} else {
			caseIsEmpty = true
			caseHasExitedNormally = true
			intersection = slices.Clone(returnInputVarsInitialized)
		}
	}

	for _, v := range en.Variants {
		if !seen[v] {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("missing variant case %q", v)})
			return
		}
	}

	st.returnedInputVarsInitialized = slices.Clone(intersection)
	if caseHasNoReturn || caseIsEmpty {
		st.returnFlowResult = flowFallsThrough
		return
	}

	st.returnFlowResult = flowReturns
	st.returnedInputVarsInitialized = slices.Clone(returnInputVarsInitialized)
	return
}

// fetchEnumVariant fetches variant from enum type when found
func fetchEnumVariant(name string, en *EnumType) (string, bool) {
	for _, v := range en.Variants {
		if v == name {
			return v, true
		}
	}

	return "", false
}

// checkSwitchEnumBody loops over enum type switch base body for validation
// and creates related statements
func (c *Checker) checkSwitchEnumBody(body []ast.Stmt, returnInputVarsInitialized []string) (st stmtInfo) {
	oldScope := c.scope
	oldInSwitchCase := c.inSwitchCase
	defer func() {
		c.scope = oldScope
		c.inSwitchCase = oldInSwitchCase
	}()
	c.inSwitchCase = true

	c.scope = NewScope(c.scope)

	var (
		cStmt            stmtInfo
		returnFlowResult returnFlow
	)

	cStmt.returnedInputVarsInitialized = slices.Clone(returnInputVarsInitialized)
	for _, b := range body {
		if returnFlowResult == flowReturns {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("unreachable code at %d:%d", b.Start().Line, b.End().Line)})
			return
		}

		if _, ok := b.(*ast.FallThroughStmt); ok {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("fallthrough is forbidden in enum switch type")})
			return
		}
		cStmt = c.checkStmt(b, cStmt.returnedInputVarsInitialized)
		returnFlowResult = cStmt.returnFlowResult
	}

	return cStmt
}

// declareComptimeDecls declares all "comptime" declarations
func (c *Checker) declareComptimeDecls() {
	for _, decl := range c.comptimeDecls {
		switch t := decl.(type) {
		case *ast.ConstDecl:
			// *ast.ConstDecl IS EMPTY ON PURPOSE because WE MUST to validate target and value type
			// before declaring symbol
		case *ast.FuncDecl:
			c.declareComptimeFuncSymbol(t)
		default:
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid comptime")})
		}
	}
}

// declareComptimeFuncSymbol declares "comptime" func declarations.
// An error is emitted if any
func (c *Checker) declareComptimeFuncSymbol(decl *ast.FuncDecl) {
	if decl.Receiver != nil {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("comptime methods are forbidden")})
		return
	}

	if !c.declareNoShadow(c.pkgScope, &Symbol{
		Name:       decl.Name.Value,
		Kind:       SymFunc,
		Decl:       decl,
		IsComptime: true,
	}, "symbol") {
		return
	}
	c.funcDecls = append(c.funcDecls, decl)
}

// checkComptimeValues validates all "comptime" declarations.
// An error is emitted if any
func (c *Checker) checkComptimeValues() {
	for _, decl := range c.comptimeDecls {
		switch t := decl.(type) {
		case *ast.ConstDecl:
			c.checkComptimeConstDecl(t)
		case *ast.FuncDecl:
			c.checkComptimeFuncDecl(t)
		}
	}
}

// checkComptimeConstDecl validates "comptime" const declarations.
// An error is emitted if any
func (c *Checker) checkComptimeConstDecl(decl *ast.ConstDecl) {
	targetType := c.resolveType(decl.Type)
	if !c.isValidComptimeType(targetType) {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid comptime type")})
		return
	}

	valueType := c.checkComptimeExpr(decl.Init)
	if !IsAssignableTo(targetType, valueType) {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("cannot assign value of type %T to const of type %T", valueType, targetType)})
		return
	}

	name := typeDeclName(decl)
	if !c.declareNoShadow(c.pkgScope, &Symbol{
		Name:       name,
		Kind:       SymConst,
		Type:       targetType,
		Decl:       decl,
		IsComptime: true,
	}, "symbol") {
		return
	}

	c.comptimeInfos = append(c.comptimeInfos, ComptimeInfo{
		Name: name,
		Kind: SymConst,
		Decl: decl,
	})
	c.constDecls = append(c.constDecls, decl)
}

// isValidComptimeType validates comptime type
func (c *Checker) isValidComptimeType(t Type) bool {
	if t == nil || IsInvalid(t) {
		return false
	}

	u := unwrapNamed(t)
	switch u.(type) {
	case *MapType, *HashMapType, *SliceType, *StructType, *SumType, *EnumType, *InterfaceType:
		return false
	}
	return true
}

// checkComptimeExpr validates allowed comptime expression
func (c *Checker) checkComptimeExpr(expr ast.Expr) Type {
	switch t := expr.(type) {
	case *ast.IntLitExpr, *ast.FloatLitExpr, *ast.BoolLitExpr, *ast.StringLitExpr:
		ce, _ := c.checkExpr(t)
		return ce

	case *ast.ParenExpr:
		return c.checkComptimeExpr(t.Inner)

	case *ast.UnaryExpr:
		if IsInvalid(c.checkComptimeExpr(t.Right)) {
			return TInvalid
		}
		ce, _ := c.checkExpr(expr)
		return ce

	case *ast.BinaryExpr:
		left := c.checkComptimeExpr(t.Left)
		right := c.checkComptimeExpr(t.Right)
		if IsInvalid(left) || IsInvalid(right) {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid comptime binary expression")})
			return TInvalid
		}
		ce, _ := c.checkExpr(expr)
		return ce

	case *ast.CallExpr:
		return c.checkComptimeCallExpr(t)

	case *ast.IdentExpr:
		var sym *Symbol
		if c.useScope && c.scope != nil {
			sym = c.scope.Lookup(t.Name.Value)
		} else {
			sym = c.pkgScope.Lookup(t.Name.Value)
		}
		if sym == nil || sym.Type == nil {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid comptime symbol")})
			return TInvalid
		}

		// this must stay as is. vars are only allowed in functions
		if sym.Kind == SymVar && !c.inComptimeFunc {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("variables are forbidden with comptime outside of function")})
			return TInvalid
		}

		return sym.Type

	case *ast.IndexExpr, *ast.SliceExpr, *ast.SliceLitExpr:
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("forbidden comptime expression")})
		return TInvalid

	default:
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("unsupported comptime expression")})
		return TInvalid
	}
}

// checkComptimeCallExpr validates allowed comptime call expression
func (c *Checker) checkComptimeCallExpr(expr *ast.CallExpr) Type {
	calleeTypeCheck := c.checkComptimeExpr(expr.Callee)
	if IsInvalid(calleeTypeCheck) {
		return TInvalid
	}

	calleeType := c.checkComptimeExpr(expr.Callee)
	if named, ok := calleeType.(*NamedType); ok {
		if len(expr.Args) != 1 {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("too many arguments in %#v, expected 1 got %d", named.Name, len(expr.Args))})
			return TInvalid
		}

		arg := c.checkComptimeExpr(expr.Args[0])
		if IsInvalid(arg) {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid comptime type")})
			return TInvalid
		}

		if IsConvertibleTo(arg, named) {
			return named
		}
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("cannot convert %#v to %s", arg, named.Name)})
		return TInvalid
	}

	if builtin, ok := calleeType.(*BuiltinType); ok {
		if len(expr.Args) != 1 {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("too many arguments in %#v, expected 1 got %d", expr, len(expr.Args))})
			return TInvalid
		}

		arg := c.checkComptimeExpr(expr.Args[0])
		if IsInvalid(arg) {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid comptime type")})
			return TInvalid
		}

		if IsConvertibleTo(arg, builtin) {
			return calleeType
		}
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("cannot convert %#v to %#v", arg, builtin)})
		return TInvalid
	}

	if fn, ok := calleeType.(*FuncMethod); ok {
		sym := c.pkgScope.Lookup(fn.Name)
		if sym == nil || sym.Kind != SymFunc || !sym.IsComptime {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid comptime function")})
			return TInvalid
		}

		if len(expr.Args) != len(fn.FuncType.Params) {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid comptime func length argument, expected %d got %d", len(expr.Args), len(fn.FuncType.Params))})
			return TInvalid
		}

		for i, p := range fn.FuncType.Params {
			argType := c.checkComptimeExpr(expr.Args[i])
			if !c.isValidComptimeType(p.Type) {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid comptime type")})
				return TInvalid
			}

			if !IsAssignableTo(p.Type, argType) {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid comptime type")})
				return TInvalid
			}
		}

		if len(fn.FuncType.Results) != 1 {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("comptime func %q must return exactly one value", fn.Name)})
			return TInvalid
		}

		for _, p := range fn.FuncType.Results {
			if !c.isValidComptimeType(p.Type) {
				c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid comptime type")})
				return TInvalid
			}
		}

		return fn.FuncType.Results[0].Type
	}

	c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid comptime expression")})
	return TInvalid
}

// checkComptimeFuncDecl validates "comptime" func declarations.
// An error is emitted if any
func (c *Checker) checkComptimeFuncDecl(decl *ast.FuncDecl) {
	// we enforce the validation of the declaration here because we can falsly
	// append comptimeInfos with wrong declarations when it's not called.
	// So we extra valides all funcs

	sym := c.pkgScope.Lookup(decl.Name.Value)
	if sym == nil || sym.Type == nil {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid comptime func")})
		return
	}

	if sym.Kind != SymFunc || !sym.IsComptime {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid comptime declaration")})
		return
	}

	fn, ok := sym.Type.(*FuncMethod)
	if !ok {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid comptime symbol")})
		return
	}

	for _, p := range fn.FuncType.Params {
		if !c.isValidComptimeType(p.Type) {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid comptime type")})
			return
		}
	}

	if len(fn.FuncType.Results) != 1 {
		c.errors = append(c.errors, Diagnostic{
			Err: fmt.Errorf("comptime func %q must return exactly one value", decl.Name.Value),
		})
		return
	}

	for _, p := range fn.FuncType.Results {
		if !c.isValidComptimeType(p.Type) {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid comptime type")})
			return
		}
	}

	c.comptimeInfos = append(c.comptimeInfos, ComptimeInfo{
		Name: decl.Name.Value,
		Kind: SymFunc,
		Decl: decl,
	})
}

// checkExprInCurrentMode checks if we are in comptime func or not and returns expression Type.
// When in comptime func we validate authorized expressions
func (c *Checker) checkExprInCurrentMode(expr ast.Expr) (Type, Expr) {
	if c.inComptimeFunc {
		t := c.checkComptimeExpr(expr)
		if !c.isValidComptimeType(t) {
			c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid comptime type")})
			return TInvalid, nil
		}
		return t, nil
	}
	return c.checkExpr(expr)
}

// checkTypeInCurrentMode checks if we are in comptime func or not and returns expression Type.
// When in comptime func we validate authorized expressions
func (c *Checker) checkTypeInCurrentMode(t Type) Type {
	if c.inComptimeFunc && !c.isValidComptimeType(t) {
		c.errors = append(c.errors, Diagnostic{Err: fmt.Errorf("invalid comptime type")})
		return TInvalid
	}
	return t
}

// checkReturnVarsInitialized verifies if new vars are already present in current list of vars.
// When not, it will append it to the current list
func (c *Checker) checkReturnVarsInitialized(a []string, s string) []string {
	if s != "" && !slices.Contains(a, s) {
		a = append(a, s)
	}
	return a
}

// resolvedSymbol resolves type checked Symbol for
// the High-Level Intermediate Representation
func resolvedSymbol(s *Symbol) ResolvedSymbol {
	if s == nil {
		return ResolvedSymbol{}
	}

	return ResolvedSymbol{
		Name: s.Name,
		Kind: s.Kind,
		Type: s.Type,
	}
}
