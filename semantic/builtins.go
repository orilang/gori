package semantic

import "strconv"

var (
	TInvalid = &InvalidType{}
	TBool    = &BuiltinType{Kind: Bool}
	TInt     = &BuiltinType{Kind: Bool}
	TInt8    = &BuiltinType{Kind: Int8}
	TInt32   = &BuiltinType{Kind: Int32}
	TInt64   = &BuiltinType{Kind: Int64}
	TUInt    = &BuiltinType{Kind: UInt}
	TUInt8   = &BuiltinType{Kind: UInt8}
	TUInt32  = &BuiltinType{Kind: UInt32}
	TUInt64  = &BuiltinType{Kind: UInt64}
	TFloat   = &BuiltinType{Kind: Float}
	TFloat32 = &BuiltinType{Kind: Float32}
	TFloat64 = &BuiltinType{Kind: Float64}
	TString  = &BuiltinType{Kind: String}
)

func (t *InvalidType) typeNode()      {}
func (t *InvalidType) String() string { return "<invalid>" }

func (t *UntypedNilType) typeNode()      {}
func (t *UntypedNilType) String() string { return "nil" }

func (t *BuiltinType) typeNode() {}
func (t *BuiltinType) String() string {
	switch t {
	case TBool:
		return "bool"
	case TInt:
		return "int"
	case TInt8:
		return "int8"
	case TInt32:
		return "int32"
	case TInt64:
		return "int64"
	case TUInt:
		return "uint"
	case TUInt8:
		return "uint8"
	case TUInt32:
		return "uint32"
	case TUInt64:
		return "uint64"
	case TFloat:
		return "float"
	case TFloat32:
		return "float32"
	case TFloat64:
		return "float64"
	case TString:
		return "string"
	default:
		return "invalidBuiltin"
	}
}

func (t *NamedType) typeNode() {}
func (t *NamedType) String() string {
	return t.Name
}

func (t *ArrayType) typeNode() {}
func (t *ArrayType) String() string {
	return "[" + strconv.FormatInt(t.Len, 10) + "]" + t.Elem.String()
}

func (t *SliceType) typeNode() {}
func (t *SliceType) String() string {
	return "[]" + t.Elem.String()
}

func (t *MapType) typeNode() {}
func (t *MapType) String() string {
	prefix := "map"
	if t.Kind == MapHash {
		prefix = "hashmap"
	}
	return prefix + "[" + t.Key.String() + "]" + t.Value.String()
}

func (t *StructType) typeNode() {}
func (t *StructType) String() string {
	return "struct"
}

func (t *FuncType) typeNode() {}
func (t *FuncType) String() string {
	return "funcType"
}

func (t *Param) typeNode() {}
func (t *Param) String() string {
	return "param"
}

func (t *FuncMethod) typeNode() {}
func (t *FuncMethod) String() string {
	return "funcMethod"
}

func (t *InterfaceType) typeNode() {}
func (t *InterfaceType) String() string {
	return "interface"
}

func (t *EnumType) typeNode() {}
func (t *EnumType) String() string {
	return "enum"
}

func (t *SumType) typeNode() {}
func (t *SumType) String() string {
	return "sum"
}
