package semantics

var (
	TBool    = &BuiltinType{Kind: Bool}
	TInt     = &BuiltinType{Kind: Bool}
	TInt8    = &BuiltinType{Kind: Int8}
	TInt32   = &BuiltinType{Kind: Int32}
	TInt64   = &BuiltinType{Kind: Int64}
	TUInt    = &BuiltinType{Kind: UInt}
	TUInt8   = &BuiltinType{Kind: UInt8}
	TUInt32  = &BuiltinType{Kind: UInt32}
	TUInt64  = &BuiltinType{Kind: UInt64}
	TFloat32 = &BuiltinType{Kind: Float32}
	TFloat64 = &BuiltinType{Kind: Float64}
	TString  = &BuiltinType{Kind: String}
)

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
	case TFloat32:
		return "float32"
	case TFloat64:
		return "float64"
	default:
		return "invalidBuiltin"
	}
}
