package semantics

import (
	"github.com/orilang/gori/token"
)

// IsIdentical returns whether both types are identical
func IsIdentical(a, b Type) bool {
	if a == nil || b == nil {
		return false
	}

	switch t1 := a.(type) {
	case *BuiltinType:
		if t2, ok := b.(*BuiltinType); ok {
			return t1 == t2
		}

	case *NamedType:
		if t2, ok := b.(*NamedType); ok {
			return t1.Name == t2.Name && t1.UnderlyingType == t2.UnderlyingType
		}

	case *ArrayType:
		if t2, ok := b.(*ArrayType); ok {
			return t1.Elem == t2.Elem && t1.Len == t2.Len
		}

	case *SliceType:
		if t2, ok := b.(*SliceType); ok {
			return IsIdentical(t1.Elem, t2.Elem)
		}

	case *MapType:
		if t2, ok := b.(*MapType); ok {
			return t1.Kind == t2.Kind && t1.Key == t2.Key && t1.Value == t2.Value
		}

	case *StructType:
		if t2, ok := b.(*StructType); ok {
			if len(t1.Fields) == len(t2.Fields) {
				for k := range t1.Fields {
					if t1.Fields[k].Name == t2.Fields[k].Name && IsIdentical(t1.Fields[k].Type, t2.Fields[k].Type) {
						continue
					} else {
						return false
					}
				}
				return true
			}
		}

	case *FuncType:
		if t2, ok := b.(*FuncType); ok {
			if t1 != nil && t2 != nil && len(t1.Params) == len(t2.Params) && len(t1.Results) == len(t2.Results) {
				for k := range t1.Params {
					if t1.Params[k].Name == t2.Params[k].Name && IsIdentical(t1.Params[k].Type, t2.Params[k].Type) {
						continue
					} else {
						return false
					}
				}
				for k := range t1.Results {
					if t1.Results[k].Name == t2.Results[k].Name && IsIdentical(t1.Results[k].Type, t2.Results[k].Type) {
						continue
					} else {
						return false
					}
				}
				return true
			}
		}

	case *Param:
		if t2, ok := b.(*Param); ok {
			if t1.Name == t2.Name && t1.Type == t2.Type {
				return true
			}
		}

	case *FuncMethod:
		if t2, ok := b.(*FuncMethod); ok {
			if t1.Name == t2.Name {
				if IsIdentical(t1.FuncType, t2.FuncType) {
					return true
				} else {
					return false
				}
			}
		}

	case *InterfaceType:
		if t2, ok := b.(*InterfaceType); ok {
			if len(t1.Methods) == len(t2.Methods) {
				for k := range t1.Methods {
					if IsIdentical(&t1.Methods[k], &t2.Methods[k]) {
						continue
					} else {
						return false
					}
				}
				return true
			}
		}

	case *Enum:
		if t2, ok := b.(*Enum); ok {
			if t1.Name == t2.Name && len(t1.Variants) == len(t2.Variants) {
				for k := range t1.Variants {
					if t1.Variants[k] == t2.Variants[k] {
						continue
					} else {
						return false
					}
				}
				return true
			}
		}

	case *SumType:
		if t2, ok := b.(*SumType); ok {
			if t1.Name == t2.Name && len(t1.Variants) == len(t2.Variants) {
				for k := range t1.Variants {
					if t1.Variants[k].Name == t2.Variants[k].Name && len(t1.Variants[k].Field) == len(t2.Variants[k].Field) {
						for kv := range t1.Variants[k].Field {
							if IsIdentical(&t1.Variants[k].Field[kv], &t2.Variants[k].Field[kv]) {
								continue
							} else {
								return false
							}
						}
						return true
					}
				}
			}
		}
	}
	return false
}

// IsAssignableTo verifies if provided parameters are assignable
func IsAssignableTo(src, dst Type) bool {
	switch src.(type) {
	case *BuiltinType:
		if dst == nil {
			return false
		}

	case *ArrayType, *SliceType, *MapType, *FuncMethod, *InterfaceType, *InvalidType:
		if dst == nil {
			return true
		}
	}

	if dst == nil {
		return false
	}

	return IsIdentical(src, dst)
}

// IsNumeric verifies if provided parameters is numeric
func IsNumeric(t Type) bool {
	switch t1 := t.(type) {
	case *BuiltinType:
		switch t1 {
		case TInt, TInt8, TInt32, TInt64, TUInt, TUInt8, TUInt32, TUInt64, TFloat, TFloat32, TFloat64:
			return true
		}

	case *NamedType:
		return IsNumeric(t1.UnderlyingType)
	}
	return false
}

// IsInteger verifies if provided parameter is an integer
func IsInteger(t Type) bool {
	switch t1 := t.(type) {
	case *BuiltinType:
		switch t1 {
		case TInt, TInt8, TInt32, TInt64, TUInt, TUInt8, TUInt32, TUInt64:
			return true
		}

	case *NamedType:
		return IsInteger(t1.UnderlyingType)
	}
	return false
}

// IsBool verifies if provided parameter is a boolean
func IsBool(t Type) bool {
	switch t1 := t.(type) {
	case *BuiltinType:
		if t1 == TBool {
			return true
		}

	case *NamedType:
		return IsBool(t1.UnderlyingType)
	}
	return false
}

// IsString verifies if provided parameter is a string
func IsString(t Type) bool {
	switch t1 := t.(type) {
	case *BuiltinType:
		if t1 == TString {
			return true
		}

	case *NamedType:
		return IsString(t1.UnderlyingType)
	}
	return false
}

// IsInvalid verifies if provided parameter is invalid
func IsInvalid(t Type) bool {
	switch t.(type) {
	case *InvalidType:
		return true
	}
	return false
}

// IsConvertibleTo verifies if provided parameters are convertible
func IsConvertibleTo(src, dst Type) bool {
	if IsInvalid(src) || IsInvalid(dst) {
		return true
	}
	if IsIdentical(src, dst) {
		return true
	}
	if IsNumeric(src) && IsNumeric(dst) {
		return true
	}
	return false
}

// SupportsBinaryOp verifies if provided parameters supports binary operations
func SupportsBinaryOp(t Type, op token.Kind) bool {
	switch op {
	case token.Plus, token.Minus, token.Star, token.Slash:
		return IsNumeric(t)

	case token.Modulo:
		return IsInteger(t)

	case token.And, token.Or:
		return IsBool(t)

	case token.Eq, token.Neq:
		return IsComparable(t)

	case token.Lt, token.Lte, token.Gt, token.Gte:
		return IsOrdered(t)

	default:
		return false
	}
}

// SupportsUnaryOp verifies if provided parameters supports unary operations
func SupportsUnaryOp(t Type, op token.Kind) bool {
	switch op {
	case token.Not:
		return IsBool(t)

	case token.Plus, token.Minus:
		return IsNumeric(t)

	default:
		return false
	}
}

// IsComparable verifies if the provided parameter is comparable
// like == or !=
func IsComparable(t Type) bool {
	return IsBool(t) || IsNumeric(t) || IsString(t)
}

// IsOrdered verifies if the provided parameters is ordered
// like < <= > >=
func IsOrdered(t Type) bool {
	return IsNumeric(t) || IsString(t)
}
