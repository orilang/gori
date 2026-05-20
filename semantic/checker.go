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

// declareNoShadow declares new symbol by enforcing there is no shadowing
// and return true when NOT exists
func (c *Checker) declareNoShadow(scope *Scope, sym *Symbol, kind string) bool {
	if scope == nil {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("cannot declare %q on nil scope", sym.Name)})
		return false
	}

	if sym.Name == "" || sym.Name == "_" {
		return true
	}

	if scope.Lookup(sym.Name) != nil {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("%s %q already declared", kind, sym.Name)})
		return false
	}

	if !scope.Declare(sym) {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("internal checker error: %s %q already declared in current scope", kind, sym.Name)})
		return false
	}

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

// Check performs the type checking step in order to validate
// all code definition structure and fill diagnostics when
// errros are found
func (c *Checker) Check(file *ast.File) []Diagnostics {
	c.collectTopLevelSymbols(file)
	c.createTypeObjects()
	c.resolveTypeDecls()
	c.resolveFuncSignatures()
	c.resolveMethodSignatures()
	c.checkImplementsDecls()
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
			if !c.declareNoShadow(c.pkgScope, &Symbol{
				Name: typeDeclName(decl),
				Kind: SymConst,
				Decl: decl,
			}, "symbol") {
				continue
			}
			c.constDecls = append(c.constDecls, d)

		case *ast.ImplementsDecl:
			c.implDecls = append(c.implDecls, d)
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
func (c *Checker) declareMethodSymbol(receiver *NamedType, fm *FuncMethod, decl *ast.FuncDecl) {
	seen := c.methods
	if seen == nil {
		c.methods = make(map[*NamedType]map[string]*FuncMethod)
	}

	rcv := c.methods[receiver]
	if rcv == nil {
		rcv = make(map[string]*FuncMethod)
		c.methods[receiver] = rcv
	}

	if _, exists := rcv[fm.Name]; exists {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("method %q already declared", fm.Name)})
		return
	}

	rcv[decl.Name.Value] = fm
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

		len, ok := c.evalArrayLen(v.Len)
		if !ok || len < 0 {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("invalid array length type")})
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
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("invalid map key type")})
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
					c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("division by 0 is forbidden")})
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

// resolveMethodSignatures resolves the ast FuncDecl declaration with receiver to return
// the semantic view of the FuncMethod. A diagnostic is emitted when duplicates found
func (c *Checker) resolveMethodSignatures() {
	for _, decl := range c.methodDecls {
		recvType := c.resolveType(decl.Receiver.Type)
		namedRcv, ok := recvType.(*NamedType)
		if !ok {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("receiver must be a named type got %#v", namedRcv)})
			continue
		}

		if _, ok := unwrapNamed(namedRcv).(*InterfaceType); ok {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("receiver cannot be an interface type")})
			continue
		}

		fm := &FuncMethod{
			Name: decl.Name.Value,
			FuncType: &FuncType{
				Params:  c.resolveParams("param", decl.Params),
				Results: c.resolveParams("result", decl.Results.List),
			},
		}

		c.declareMethodSymbol(namedRcv, fm, decl)
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
				c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("%s name %q already declared", kind, p.Name.Value)})
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
	sym.Decl = decl
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
			switch t.Operator.Kind {
			case token.Eq, token.Neq, token.And, token.Or, token.Lt, token.Lte, token.Gt, token.Gte:
				return TBool

			default:
				return left
			}
		}
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("invalid binary operation %q of type %s with type %s", t.Operator.Value, left.String(), right.String())})
		return TInvalid

	case *ast.CallExpr:
		if sel, ok := t.Callee.(*ast.SelectorExpr); ok {
			if named, ok := c.checkExpr(sel.X).(*NamedType); ok {
				if method, ok := c.lookupMethodType(named, sel.Selector.Value); ok {
					if len(t.Args) != len(method.FuncType.Params) {
						c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("too many arguments %s func params required %d got %d", method.Name, len(method.FuncType.Params), len(t.Args))})
						return TInvalid
					}
					for k, v := range method.FuncType.Params {
						x := c.checkExpr(t.Args[k])
						if !IsAssignableTo(v.Type, x) {
							c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("cannot assign value of type %T to const of type %T", v.Type, x)})
							return TInvalid
						}
					}
					if len(method.FuncType.Results) == 0 {
						return nil
					}
					if len(method.FuncType.Results) > 1 {
						c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("too many arguments returned by %s func", method.Name)})
						return TInvalid
					}
					return method.FuncType.Results[0].Type
				}
			}
		}

		calleeType := c.checkExpr(t.Callee)
		if named, ok := calleeType.(*NamedType); ok {
			if len(t.Args) != 1 {
				c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("too many arguments in %#v, expected 1 got %d", named.Name, len(t.Args))})
				return TInvalid
			}
			arg := c.checkExpr(t.Args[0])
			if IsConvertibleTo(arg, named) {
				return named
			}
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("cannot convert %#v to  %s", arg, named.Name)})
			return TInvalid
		}

		if builtin, ok := calleeType.(*BuiltinType); ok {
			if len(t.Args) != 1 {
				c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("too many arguments in %#v, expected 1 got %d", expr, len(t.Args))})
				return TInvalid
			}

			arg := c.checkExpr(t.Args[0])
			if IsConvertibleTo(arg, builtin) {
				return calleeType
			}
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("cannot convert %#v to %#v", arg, builtin)})
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

		case *HashMapType:
			if !IsIdentical(decl.Key, index) {
				c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("invalid hashmap index expression of type %#v", index)})
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

	for _, fn := range c.methodDecls {
		c.checkMethodBody(fn)
	}
}

// checkFuncBody validates function body.
// An error is emitted if any
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
		if !c.declareNoShadow(c.scope, &Symbol{
			Name: p.Name,
			Kind: SymVar,
			Type: p.Type,
		}, "variable") {
			continue
		}
	}

	c.checkBlockStmt(fn.Body)
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
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("receiver must be a named type got %#v", namedRcv)})
		return
	}

	method, ok := c.lookupMethodType(namedRcv, fn.Name.Value)
	if !ok {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("method type undefined")})
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

	c.checkBlockStmt(fn.Body)
}

// checkBlockStmt loops over block statements in order to check/declare them
// within its dedicated local scope
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
		c.checkStmt(stmt)
	}
}

// checkStmt checks all possible statements kind.
// This is splitted here for reusability
func (c *Checker) checkStmt(stmt ast.Stmt) {
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
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("unsupported declaration %#v", decl)})
		}

	case *ast.AssignStmt:
		c.checkAssigmentStmt(t)

	case *ast.ReturnStmt:
		c.checkReturnStmt(t)

	case *ast.IncDecStmt:
		c.checkIncDecStmt(t)

	case *ast.ExprStmt:
		c.checkExprStmt(t)

	case *ast.IfStmt:
		c.checkIfStmt(t)

	case *ast.ForStmt:
		c.checkForStmt(t)

	case *ast.RangeStmt:
		c.checkRangeStmt(t)

	case *ast.SwitchStmt:
		c.checkSwitchStmt(t)

	case *ast.FallThroughStmt:
		c.checkFallThroughStmt(t)

	case *ast.BreakStmt:
		c.checkBreakStmt(t)

	case *ast.ContinueStmt:
		c.checkContinueStmt(t)

	default:
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("unsupported statement %#v", stmt)})
	}
}

// checkScopeConstDecl validates constant targetType and valueType.
// An error is emitted if any
func (c *Checker) checkScopeConstDecl(decl *ast.ConstDecl) {
	targetType := c.resolveType(decl.Type)
	valueType := c.checkExpr(decl.Init)

	if !IsAssignableTo(targetType, valueType) {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("cannot assign value of type %T to var of type %T", valueType, targetType)})
		return
	}

	if !c.declareNoShadow(c.scope, &Symbol{
		Name: decl.Name.Value,
		Kind: SymConst,
		Type: targetType,
		Decl: decl,
	}, "const") {
		return
	}
}

// checkScopeVarDecl validates constant targetType and valueType.
// An error is emitted if any
func (c *Checker) checkScopeVarDecl(decl *ast.VarDecl) {
	targetType := c.resolveType(decl.Type)
	valueType := c.checkExpr(decl.Init)

	if !IsAssignableTo(targetType, valueType) {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("cannot assign value of type %T to var of type %T", valueType, targetType)})
		return
	}

	if !c.declareNoShadow(c.scope, &Symbol{
		Name: decl.Name.Value,
		Kind: SymVar,
		Type: targetType,
		Decl: decl,
	}, "variable") {
		return
	}
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
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("variable %#v not an identifier", decl)})
		return
	}

	if !c.declareNoShadow(c.scope, &Symbol{
		Name: x.Name.Value,
		Kind: SymVar,
		Type: valueType,
	}, "variable") {
		return
	}
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
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("type %q is undefined", decl.TypeName.Value)})
		return
	}

	ifaceType := c.resolveType(decl.Interface)
	if IsInvalid(ifaceType) {
		return
	}

	ifaceNamed, ok := ifaceType.(*NamedType)
	if !ok {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("implements target must be a named interface")})
		return
	}

	// iface, ok := unwrapNamed(ifaceNamed).(*InterfaceType)
	_, ok = unwrapNamed(ifaceNamed).(*InterfaceType)
	if !ok {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("type %q is not an interface", ifaceNamed.Name)})
		return
	}

	// TODO: when receivers/method sets exist:
	// if !c.implementsInterface(ifaceType, iface) {
	// 	return
	// }

	c.implInfos = append(c.implInfos, ImplInfo{
		TypeName:      decl.TypeName.Value,
		Type:          sym.Type,
		InterfaceName: ifaceNamed.Name,
		Interface:     ifaceNamed,
		Decl:          decl,
	})
}

// implementsInterface loops over interface methods to verify if all methods
// has been implemented with the right signatures
// func (c *Checker) implementsInterface(ifaceType Type, iface *InterfaceType) bool {
// 	for _, im := range iface.Methods {
// 		fn, ok := c.lookupMethod(ifaceType, im.Name)
// 		if !ok {
// 			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("method %q not implemented", im.Name)})
// 			return false
// 		}

// 		if !IsIdentical(im.FuncType, fn.FuncType) {
// 			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("wrong method %q signature", im.Name)})
// 			return false
// 		}
// 	}
// 	return true
// }

// lookupMethod loops over funcs methods to match provided name.
// It returns its semantic type and true when found
// func (c *Checker) lookupMethod(ifaceType Type, name string) (*FuncMethod, bool) {
// 	if _, ok := c.methods[ifaceType];ok {
// 	for k, fn := range c.methods {
// 		if fm, ok := fn[name]; ok  {
// 		}
// 	}
// 	}
// 	return nil, false
// }

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
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("call expression statement must be a function call, got %#v", call)})
		return
	}

	if _, ok := call.Callee.(*ast.SelectorExpr); ok {
		_ = c.checkExpr(stmt.Expr)
		return
	}

	calleType := c.checkExpr(call.Callee)
	if _, ok := calleType.(*FuncMethod); !ok {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("callee type expression statement must be a function call got %#v", calleType)})
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

// checkIfStmt validates if statement block
func (c *Checker) checkIfStmt(stmt *ast.IfStmt) {
	if stmt == nil {
		return
	}

	condType := c.checkExpr(stmt.Condition)
	if !IsBool(condType) {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("if condition must returned a boolean")})
		return
	}

	if stmt.Then != nil {
		c.checkBlockStmt(stmt.Then)
	}

	if stmt.Else != nil {
		switch t := stmt.Else.(type) {
		case *ast.BlockStmt:
			c.checkBlockStmt(t)
		default:
			c.checkStmt(t)
		}
	}
}

// checkForStmt validates plain for statement block
func (c *Checker) checkForStmt(stmt *ast.ForStmt) {
	if stmt == nil {
		return
	}

	oldLoopDepth := c.loopDepth
	c.loopDepth++
	defer func() {
		c.loopDepth = oldLoopDepth
	}()

	if stmt.Init != nil {
		c.checkStmt(stmt.Init)
	}

	if stmt.Condition != nil {
		condType := c.checkExpr(stmt.Condition)
		if !IsBool(condType) {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("if condition must return a boolean")})
			return
		}
	}

	if stmt.Post != nil {
		c.checkStmt(stmt.Post)
	}

	if stmt.Body != nil {
		c.checkBlockStmt(stmt.Body)
	}
}

// checkForStmt validates for assigment statement
func (c *Checker) checkAssigmentStmt(stmt *ast.AssignStmt) {
	switch stmt.Operator.Kind {
	case token.Assign, token.PlusEq, token.MinusEq:
		c.checkSimpleAssignStmt(stmt)
	case token.Define:
		c.checkDefineAssignStmt(stmt)
	default:
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("unsupported assigment in for statement %#v", stmt)})
		return
	}
}

// checkRangeStmt validates range statement block
func (c *Checker) checkRangeStmt(stmt *ast.RangeStmt) {
	if stmt == nil {
		return
	}

	oldLoopDepth := c.loopDepth
	c.loopDepth++
	defer func() {
		c.loopDepth = oldLoopDepth
	}()

	iteratorType := c.checkExpr(stmt.X)
	if IsInvalid(iteratorType) {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("range expression is invalid")})
		return
	}

	rangekeyType, rangeValueType, ok := rangeVars(iteratorType)
	if !ok {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("unsupported range var type")})
		return
	}

	// for range x {}
	if stmt.Op == (token.Token{}) && stmt.Body != nil {
		c.checkBlockStmt(stmt.Body)
		return
	}

	switch stmt.Op.Kind {
	case token.Assign:
		if stmt.Key != nil && stmt.Value != nil {
			key := c.checkExpr(stmt.Key)
			if stmt.Key.Name.Value != "_" && !IsAssignableTo(rangekeyType, key) {
				c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("invalid range key type, expected %#v, got %#v", rangekeyType, key)})
				return
			}

			value := c.checkExpr(stmt.Value)
			if stmt.Value.Name.Value != "_" && !IsAssignableTo(rangeValueType, value) {
				c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("invalid range value type, expected %#v, got %#v", rangeValueType, value)})
				return
			}
		} else if stmt.Key != nil {
			if stmt.Key.Name.Value == "_" {
				c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("blank identifier for this range key is forbidden")})
				return
			}

			key := c.checkExpr(stmt.Key)
			if !IsAssignableTo(rangekeyType, key) {
				c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("invalid range key type, expected %#v, got %#v", rangekeyType, key)})
				return
			}
		}

	case token.Define:
		if stmt.Key != nil && stmt.Value != nil {
			if stmt.Key.Name.Value == "_" && stmt.Value.Name.Value == "_" {
				c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("range key and value cannot be both blank identifiers")})
				return
			}

			oldScope := c.scope
			defer func() {
				c.scope = oldScope
			}()

			blockScope := NewScope(c.scope)
			c.scope = blockScope

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
				c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("blank identifier for this range key is forbidden")})
				return
			}

			oldScope := c.scope
			defer func() {
				c.scope = oldScope
			}()

			blockScope := NewScope(c.scope)
			c.scope = blockScope

			if !c.declareNoShadow(c.scope, &Symbol{
				Name: stmt.Key.Name.Value,
				Kind: SymVar,
				Type: rangekeyType,
			}, "variable") {
				return
			}
		}

	default:
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("forbidden range token %s", stmt.Op.Value)})
		return
	}

	if stmt.Body != nil {
		c.checkBlockStmt(stmt.Body)
	}
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
func (c *Checker) checkSwitchStmt(stmt *ast.SwitchStmt) {
	if stmt == nil {
		return
	}

	oldScope := c.scope
	defer func() {
		c.scope = oldScope
	}()

	switchScope := NewScope(c.scope)
	c.scope = switchScope

	if stmt.Init != nil {
		c.checkStmt(stmt.Init)
	}

	var dcount int
	if stmt.Tag != nil {
		tagType := c.checkExpr(stmt.Tag)
		if IsInvalid(tagType) {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("tag expression is invalid")})
			return
		}

		underlying := unwrapNamed(tagType)
		switch underlying.(type) {
		case *SumType:
			c.checkSwitchStmtSumType(tagType, stmt)
			return
		case *EnumType:
			c.checkSwitchStmtEnumType(tagType, stmt)
			return
		}

		if !IsComparable(tagType) {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("tag type is not comparable got %#v", tagType)})
			return
		}

		if len(stmt.Cases) == 0 {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("switch statement has 0 cases")})
			return
		}

		seen := make(map[constKey]bool, len(stmt.Cases))
		for i, cc := range stmt.Cases {
			if cc.Case.Kind == token.KWDefault {
				dcount++

				if dcount > 1 {
					c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("multiple default clause is forbidden, got %d", dcount)})
					return
				}
			}

			for _, v := range cc.Values {
				vExpr := c.checkExpr(v)
				if !IsIdentical(tagType, vExpr) {
					c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("tag and case are not identical expected %#v got %#v", tagType, vExpr)})
					return
				}

				if ck, ok := c.constKey(v, vExpr); ok {
					if seen[ck] {
						c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("duplicate case expression %s at %d:%d", ck.value, v.Start().Line, v.End().Line)})
						return
					}
					seen[ck] = true
				}
			}

			if cc.Body != nil {
				c.checkSwitchBody(cc.Body, i == len(stmt.Cases)-1)
			}
		}
	} else {
		if len(stmt.Cases) == 0 {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("switch statement has 0 cases")})
			return
		}

		for i, cc := range stmt.Cases {
			if cc.Case.Kind == token.KWDefault {
				dcount++

				if dcount > 1 {
					c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("multiple default clause is forbidden, got %d", dcount)})
					return
				}
			}

			for _, v := range cc.Values {
				vExpr := c.checkExpr(v)
				// TODO: Tagless switch duplicates is a job for the linter
				// as it's difficult for the checker to properly handle every cases
				// without any burden

				if !IsBool(vExpr) {
					c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("switch case expected boolean got %#v", vExpr)})
					return
				}
			}

			if cc.Body != nil {
				c.checkSwitchBody(cc.Body, i == len(stmt.Cases)-1)
			}
		}
	}
}

// constKey returns true if expression is a literal one. It will be used by basic switch case
// in order to find duplicates
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
func (c *Checker) checkSwitchBody(body []ast.Stmt, isLastCaseClause bool) {
	oldScope := c.scope
	oldInSwitchCase := c.inSwitchCase
	defer func() {
		c.scope = oldScope
		c.inSwitchCase = oldInSwitchCase
	}()
	c.inSwitchCase = true

	bodyScope := NewScope(c.scope)
	c.scope = bodyScope
	for i, b := range body {
		if _, ok := b.(*ast.FallThroughStmt); ok {
			if i != len(body)-1 {
				c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("fallthrough must be the last statement of the switch case body")})
				return
			}

			if isLastCaseClause {
				c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("fallthrough is forbidden inside last switch case")})
				return
			}
			continue
		}

		c.checkStmt(b)
	}
}

// checkFallThroughStmt produces an error when not into switch case
func (c *Checker) checkFallThroughStmt(_ *ast.FallThroughStmt) {
	if !c.inSwitchCase {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("fallthrough is forbidden outside of switch case")})
	}
}

// checkBreakStmt produces an error when not into for loop statement
func (c *Checker) checkBreakStmt(_ *ast.BreakStmt) {
	if c.loopDepth == 0 {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("break is forbidden outside of loop")})
	}
}

// checkContinueStmt produces an error when not into for loop statement
func (c *Checker) checkContinueStmt(_ *ast.ContinueStmt) {
	if c.loopDepth == 0 {
		c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("continue is forbidden outside of loop")})
	}
}

// checkSwitchStmtSumType validates switch statement block
func (c *Checker) checkSwitchStmtSumType(tagType Type, stmt *ast.SwitchStmt) {
	underlying := unwrapNamed(tagType)
	sm := underlying.(*SumType)

	seen := make(map[string]bool, len(stmt.Cases))
	for _, cc := range stmt.Cases {
		if cc.Case.Kind == token.KWDefault {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("default is forbidden inside sum type switch")})
			return
		}

		call, ok := cc.Values[0].(*ast.CallExpr)
		if !ok {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("sum case must be a variant, got %#v", call)})
			return
		}

		callee, ok := call.Callee.(*ast.IdentExpr)
		if !ok {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("sum variant case must use a variant identifier, got %#v", callee)})
			return
		}

		variantName, ok := fetchSumVariant(callee.Name.Value, sm)
		if !ok {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("unknown variant name %q", callee.Name.Value)})
			return
		}

		if seen[variantName.Name] {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("duplicate variant name %s", variantName.Name)})
			return
		}
		seen[variantName.Name] = true

		if len(variantName.Field) != len(call.Args) {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("variant arguments length invalid, expected %d got %d", len(variantName.Field), len(call.Args))})
			return
		}

		var bindings []*Symbol
		seenBindings := make(map[string]bool)
		for k, v := range call.Args {
			ident, ok := v.(*ast.IdentExpr)
			if !ok {
				c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("variant expected identifier got %#v", v)})
				return
			}

			// TODO: if we need _ we will have to update the parser
			// if ident.Name.Value == "_" {
			// 	continue
			// }

			if seenBindings[ident.Name.Value] {
				c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("variable %q already declared", ident.Name.Value)})
				return
			}
			seenBindings[ident.Name.Value] = true

			if c.scope.Lookup(ident.Name.Value) != nil {
				c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("variable %q already declared", ident.Name.Value)})
				return
			}

			bindings = append(bindings, &Symbol{
				Name: ident.Name.Value,
				Kind: SymVar,
				Type: variantName.Field[k].Type,
			})
		}

		if cc.Body != nil {
			c.checkSwitchSumBody(cc.Body, bindings)
		}
	}

	for _, v := range sm.Variants {
		if !seen[v.Name] {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("missing variant case %q", v.Name)})
			return
		}
	}
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
func (c *Checker) checkSwitchSumBody(body []ast.Stmt, bindings []*Symbol) {
	oldScope := c.scope
	oldInSwitchCase := c.inSwitchCase
	defer func() {
		c.scope = oldScope
		c.inSwitchCase = oldInSwitchCase
	}()
	c.inSwitchCase = true

	bodyScope := NewScope(c.scope)
	c.scope = bodyScope

	for _, b := range bindings {
		if !c.declareNoShadow(c.scope, b, "variable") {
			return
		}
	}

	for _, b := range body {
		if _, ok := b.(*ast.FallThroughStmt); ok {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("fallthrough is forbidden in sum switch type")})
			return
		}

		c.checkStmt(b)
	}
}

// checkSwitchStmtEnumType validates switch statement block
func (c *Checker) checkSwitchStmtEnumType(tagType Type, stmt *ast.SwitchStmt) {
	underlying := unwrapNamed(tagType)
	en := underlying.(*EnumType)

	seen := make(map[string]bool, len(stmt.Cases))
	for _, cc := range stmt.Cases {
		if cc.Case.Kind == token.KWDefault {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("default is forbidden inside enum type switch")})
			return
		}

		if len(cc.Values) != 1 {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("enum case must have exactly one variant")})
			return
		}

		ident, ok := cc.Values[0].(*ast.IdentExpr)
		if !ok {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("enum case must be a variant, got %#v", ident)})
			return
		}

		variantName, ok := fetchEnumVariant(ident.Name.Value, en)
		if !ok {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("unknown variant name %q", ident.Name.Value)})
			return
		}

		if seen[variantName] {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("duplicate variant name %s", variantName)})
			return
		}
		seen[variantName] = true

		if cc.Body != nil {
			c.checkSwitchEnumBody(cc.Body)
		}
	}

	for _, v := range en.Variants {
		if !seen[v] {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("missing variant case %q", v)})
			return
		}
	}
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
func (c *Checker) checkSwitchEnumBody(body []ast.Stmt) {
	oldScope := c.scope
	oldInSwitchCase := c.inSwitchCase
	defer func() {
		c.scope = oldScope
		c.inSwitchCase = oldInSwitchCase
	}()
	c.inSwitchCase = true

	bodyScope := NewScope(c.scope)
	c.scope = bodyScope

	for _, b := range body {
		if _, ok := b.(*ast.FallThroughStmt); ok {
			c.errors = append(c.errors, Diagnostics{Err: fmt.Errorf("fallthrough is forbidden in enum switch type")})
			return
		}
		c.checkStmt(b)
	}
}
