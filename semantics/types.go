package semantics

import "github.com/orilang/gori/ast"

type Type interface {
	typeNode()
	String() string
}

type BuiltInKind int

const (
	InvalidBuiltin BuiltInKind = iota
	Bool
	Int
	Int8
	Int32
	Int64
	UInt
	UInt8
	UInt32
	UInt64
	Float
	Float32
	Float64
	String
)

type BuiltinType struct {
	Kind BuiltInKind
}

type NamedType struct {
	Name           string
	UnderlyingType Type
	Decl           ast.Position
}

type ArrayType struct {
	Len  int64
	Elem Type
}

type SliceType struct {
	Elem Type
}

type MapKind int

const (
	MapRegular MapKind = iota
	MapHash
)

type MapType struct {
	Kind  MapKind
	Key   Type
	Value Type
}

type StructField struct {
	Name string
	Type Type
}

type StructType struct {
	Fields []StructField
	Decl   ast.Position
}

type Param struct {
	Name string
	Type Type
}

type FuncType struct {
	Params  []Param
	Results []Param
}

type FuncMethod struct {
	Name     string
	FuncType *FuncType
}

type InterfaceType struct {
	Methods []FuncMethod
	Decl    ast.Position
}

type Enum struct {
	Name     string
	Variants []string
	Decl     ast.Position
}

type SumVariant struct {
	Name  string
	Field []Param
}

type SumType struct {
	Name     string
	Variants []SumVariant
	Decl     ast.Position
}

type InvalidType struct{}
type UntypeNilType struct{}
