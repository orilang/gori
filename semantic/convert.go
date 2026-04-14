package semantic

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
			return t1.Decl == t2.Decl
		}

	case *ArrayType:
		if t2, ok := b.(*ArrayType); ok {
			return IsIdentical(t1.Elem, t2.Elem) && t1.Len == t2.Len
		}

	case *SliceType:
		if t2, ok := b.(*SliceType); ok {
			return IsIdentical(t1.Elem, t2.Elem)
		}

	case *MapType:
		if t2, ok := b.(*MapType); ok {
			return t1.Kind == t2.Kind && IsIdentical(t1.Key, t2.Key) && IsIdentical(t1.Value, t2.Value)
		}

	case *StructType:
		if t2, ok := b.(*StructType); ok {
			return t1.Decl == t2.Decl
		}

	case *FuncType:
		if t2, ok := b.(*FuncType); ok {
			if t1 != nil && t2 != nil && len(t1.Params) == len(t2.Params) && len(t1.Results) == len(t2.Results) {
				for k := range t1.Params {
					if !IsIdentical(t1.Params[k].Type, t2.Params[k].Type) {
						return false
					}
				}
				for k := range t1.Results {
					if !IsIdentical(t1.Results[k].Type, t2.Results[k].Type) {
						return false
					}
				}
				return true
			}
		}

	case *FuncMethod:
		if t2, ok := b.(*FuncMethod); ok {
			if t1.Name == t2.Name {
				if !IsIdentical(t1.FuncType, t2.FuncType) {
					return false
				}
				return true
			}
		}

	case *InterfaceType:
		if t2, ok := b.(*InterfaceType); ok {
			return t1.Decl == t2.Decl
		}

	case *EnumType:
		if t2, ok := b.(*EnumType); ok {
			return t1.Decl == t2.Decl
		}

	case *SumType:
		if t2, ok := b.(*SumType); ok {
			return t1.Decl == t2.Decl
		}
	}
	return false
}

// IsAssignableTo verifies if provided parameters are assignable.
// example targetType = SliceType{Elem: TInt} and valueType = nil
func IsAssignableTo(targetType, valueType Type) bool {
	if IsInvalid(targetType) || IsInvalid(valueType) {
		return true
	}
	if IsUntypedNilType(valueType) {
		return IsNilAssignable(targetType)
	}
	return IsIdentical(targetType, valueType)
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
func IsConvertibleTo(fromType, toType Type) bool {
	if IsInvalid(fromType) || IsInvalid(toType) {
		return true
	}
	if IsIdentical(fromType, toType) {
		return true
	}
	if IsNumeric(fromType) && IsNumeric(toType) {
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

// IsOrdered verifies if the provided parameter is ordered
// like < <= > >=
func IsOrdered(t Type) bool {
	return IsNumeric(t) || IsString(t)
}

// IsUntypedNilType verifies if the provided parameter is a nil type
func IsUntypedNilType(t Type) bool {
	if _, ok := t.(*UntypedNilType); ok {
		return true
	}
	return false
}

func IsNilAssignable(t Type) bool {
	switch t.(type) {
	case *SliceType, *MapType:
		return true
	}
	return false
}
