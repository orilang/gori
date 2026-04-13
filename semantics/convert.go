package semantics

// IsIdentical returns whether both types are identical
func IsIdentical(a, b Type) bool {
	if a == nil || b == nil {
		return false
	}

	switch t1 := a.(type) {
	case *BuiltinType:
		if t2, ok := b.(*BuiltinType); ok {
			return t1.String() == t2.String()
		}

	case *NamedType:
		if t2, ok := b.(*NamedType); ok {
			return t1.Name == t2.Name && t1.UnderlyingType.String() == t2.UnderlyingType.String()
		}

	case *ArrayType:
		if t2, ok := b.(*ArrayType); ok {
			return t1.Elem.String() == t2.Elem.String()
		}

	case *SliceType:
		if t2, ok := b.(*SliceType); ok {
			return t1.Elem.String() == t2.Elem.String()
		}

	case *MapType:
		if t2, ok := b.(*MapType); ok {
			return t1.Kind == t2.Kind && t1.Key.String() == t2.Key.String() && t1.Value.String() == t2.Value.String()
		}

	case *StructType:
		if t2, ok := b.(*StructType); ok {
			if t1.String() == t2.String() && len(t1.Fields) == len(t2.Fields) {
				for k := range t1.Fields {
					if t1.Fields[k].Name == t2.Fields[k].Name && t1.Fields[k].Type.String() == t2.Fields[k].Type.String() {
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
			if t1.Name == t2.Name && t1.Type.String() == t2.Type.String() {
				return true
			}
		}

	case *FuncMethod:
		if t2, ok := b.(*FuncMethod); ok {
			if t1.String() == t2.String() && t1.Name == t2.Name {
				if t1.FuncType != nil && t2.FuncType != nil && len(t1.FuncType.Params) == len(t2.FuncType.Params) && len(t1.FuncType.Results) == len(t2.FuncType.Results) {
					for k := range t1.FuncType.Params {
						if IsIdentical(&t1.FuncType.Params[k], &t2.FuncType.Params[k]) {
							continue
						} else {
							return false
						}
					}
					for k := range t1.FuncType.Results {
						if IsIdentical(&t1.FuncType.Results[k], &t2.FuncType.Results[k]) {
							continue
						} else {
							return false
						}
					}
					return true
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
						for kv := range t1.Variants {
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
