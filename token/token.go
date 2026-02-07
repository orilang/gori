package token

// LookupKeyword lookup for the current keyword list
// and returns its kind if found otherwise identifier
func LookupKeyword(s string) Kind {
	if token, ok := keywords[s]; ok {
		return token
	}
	return Ident
}

// IsBuiltinType returns true when builtin types is found otherwise false
func IsBuiltinType(k Kind) bool {
	return builtinTypes[k]
}

// IsPrefix returns true when the provided kind is found is the list
func IsPrefix(k Kind) bool {
	return prefix[k]
}

// IsInfix returns true when the provided kind is found is the list
func IsInfix(k Kind) bool {
	return infix[k]
}

// IsPostfix returns true when the provided kind is found is the list
func IsPostfix(k Kind) bool {
	return postfix[k]
}

// IsComparison returns true when the provided kind is found is the list
func IsComparison(k Kind) bool {
	return comparison[k]
}

// IsChainingComparison returns true when the provided kind is found is the list.
// This must only be used to prevent stuffs like a < b < c
func IsChainingComparison(k Kind) bool {
	return chainingComparison[k]
}

// IsAssignment returns true when the provided kind is found is the list
func IsAssignment(k Kind) bool {
	return assignment[k]
}

// IsRangeForAssignment returns true when the provided kind is found is the list
func IsRangeForAssignment(k Kind) bool {
	return rangeForAssigment[k]
}

// IsRangeForAssignment returns true when the provided kind is found is the list
func IsIncDec(k Kind) bool {
	return incDec[k]
}

// IsVarConstTypes returns true when the provided kind is found is the list
func IsVarConstTypes(k Kind) bool {
	return varConstTypes[k]
}

// IsFuncParamTypes returns true when the provided kind is found is the list
func IsFuncParamTypes(k Kind) bool {
	return funcParamTypes[k]
}
