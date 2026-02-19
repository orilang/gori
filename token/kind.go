package token

type Kind uint16

const (
	Illegal Kind = iota
	EOF

	Comment // // or /* */

	// identifier and literals
	Ident
	IntLit
	FloatLit
	StringLit
	BoolLit

	// Keywords
	KWPackage
	KWImport
	KWFunc
	KWVar
	KWInt
	KWInt8
	KWInt32
	KWInt64
	KWUint
	KWUint8
	KWUint32
	KWUint64
	KWFloat
	KWFloat32
	KWFloat64
	KWConst
	KWString
	KWBool
	KWType
	KWStruct
	KWInterface

	KWIf
	KWElse
	KWFor
	KWBreak
	KWContinue
	KWSwitch
	KWCase
	KWDefault
	KWFallThrough
	KWReturn

	LParen   // (
	RParen   // )
	LBrace   // {
	RBrace   // }
	LBracket // [
	RBracket // ]

	Comma     // ,
	SemiComma // ;
	Colon     // :
	Dot       // .

	Assign // =
	Define // :=

	Plus    // +
	PlusEq  // +=
	PPlus   // ++
	Minus   // -
	MinusEq // -=
	MMinus  // --
	Star    // *
	StarEq  // *=
	Slash   // /
	SlashEq // /=
	Modulo  // %

	// Comparison operators
	Eq  // ==
	Neq // !=
	Lt  // <
	Lte // <=
	Gt  // >
	Gte // >=

	// Logical operators
	And // &&
	Or  // ||
	Not // !

	KWRange
	KWImplements

	Pipe // |
	KWEnum
)
