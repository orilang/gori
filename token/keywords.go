package token

var keywords = map[string]Kind{
	"package":     KWPackage,
	"import":      KWImport,
	"func":        KWFunc,
	"var":         KWVar,
	"int":         KWInt,
	"int8":        KWInt8,
	"int32":       KWInt32,
	"int64":       KWInt64,
	"uint":        KWUint,
	"uint8":       KWUint8,
	"uint32":      KWUint32,
	"uint64":      KWUint64,
	"float":       KWFloat,
	"float32":     KWFloat32,
	"float64":     KWFloat64,
	"const":       KWConst,
	"string":      KWString,
	"type":        KWType,
	"struct":      KWStruct,
	"interface":   KWInterface,
	"if":          KWIf,
	"else":        KWElse,
	"for":         KWFor,
	"break":       KWBreak,
	"continue":    KWContinue,
	"switch":      KWSwitch,
	"case":        KWCase,
	"default":     KWDefault,
	"fallthrough": KWFallThrough,
	"return":      KWReturn,
	"bool":        KWBool,
	"true":        BoolLit,
	"false":       BoolLit,
	"range":       KWRange,
}

var builtinTypes = map[Kind]bool{
	KWInt:       true,
	KWInt8:      true,
	KWInt32:     true,
	KWInt64:     true,
	KWUint:      true,
	KWUint8:     true,
	KWUint32:    true,
	KWUint64:    true,
	KWFloat:     true,
	KWFloat32:   true,
	KWFloat64:   true,
	KWConst:     true,
	KWString:    true,
	KWBool:      true,
	KWType:      true,
	KWStruct:    true,
	KWInterface: true,
}

var prefix = map[Kind]bool{
	LParen:    true,
	Ident:     true,
	IntLit:    true,
	FloatLit:  true,
	StringLit: true,
	BoolLit:   true,
	Plus:      true,
	Minus:     true,
	Not:       true,
}

var infix = map[Kind]bool{
	Plus:   true,
	Minus:  true,
	Star:   true,
	Slash:  true,
	Modulo: true,
	Eq:     true,
	Neq:    true,
	Lt:     true,
	Lte:    true,
	Gt:     true,
	Gte:    true,
	And:    true,
	Or:     true,
}

var postfix = map[Kind]bool{
	Dot:      true,
	LBracket: true,
	LParen:   true,
}

var comparison = map[Kind]bool{
	Eq:  true,
	Neq: true,
	Lt:  true,
	Lte: true,
	Gt:  true,
	Gte: true,
}

var chainingComparison = map[Kind]bool{
	Eq:  true,
	Neq: true,
	Lt:  true,
	Lte: true,
	Gt:  true,
	Gte: true,
}

var assignment = map[Kind]bool{
	Assign:  true,
	Define:  true,
	PlusEq:  true,
	MinusEq: true,
	StarEq:  true,
	SlashEq: true,
}

var rangeForAssigment = map[Kind]bool{
	Assign: true,
	Define: true,
}

var incDec = map[Kind]bool{
	PPlus:  true,
	MMinus: true,
}

var varConstTypes = map[Kind]bool{
	Ident:       true,
	KWInt:       true,
	KWInt8:      true,
	KWInt32:     true,
	KWInt64:     true,
	KWUint:      true,
	KWUint8:     true,
	KWUint32:    true,
	KWUint64:    true,
	KWFloat:     true,
	KWFloat32:   true,
	KWFloat64:   true,
	KWString:    true,
	KWBool:      true,
	KWInterface: true,
}

var funcParamTypes = map[Kind]bool{
	Ident:       true,
	KWInt:       true,
	KWInt8:      true,
	KWInt32:     true,
	KWInt64:     true,
	KWUint:      true,
	KWUint8:     true,
	KWUint32:    true,
	KWUint64:    true,
	KWFloat:     true,
	KWFloat32:   true,
	KWFloat64:   true,
	KWString:    true,
	KWBool:      true,
	KWInterface: true,
	KWFunc:      true,
}
