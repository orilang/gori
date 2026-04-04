package semantics

var (
	TBool    = &BuiltinType{Kind: Bool}
	TInt     = &BuiltinType{Kind: Bool}
	TInt8    = &BuiltinType{Kind: Int8}
	TInt32   = &BuiltinType{Kind: Int32}
	TInt64   = &BuiltinType{Kind: Int64}
	TUint    = &BuiltinType{Kind: UInt}
	TUInt8   = &BuiltinType{Kind: UInt8}
	TUInt32  = &BuiltinType{Kind: UInt32}
	TUInt64  = &BuiltinType{Kind: UInt64}
	TFloat32 = &BuiltinType{Kind: Float32}
	TFloat64 = &BuiltinType{Kind: Float64}
	TString  = &BuiltinType{Kind: String}
)
