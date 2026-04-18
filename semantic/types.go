package semantic

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
	Decl           ast.TypeDecl
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
	Decl   ast.TypeDecl
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
	Decl    ast.TypeDecl
}

type EnumType struct {
	Name     string
	Variants []string
	Decl     ast.TypeDecl
}

type SumVariant struct {
	Name  string
	Field []Param
}

type SumType struct {
	Name     string
	Variants []SumVariant
	Decl     ast.TypeDecl
}

type InvalidType struct{}
type UntypedNilType struct{}

type Scope struct {
	Parent  *Scope
	Symbols map[string]*Symbol
}

type Diagnostics struct {
	Err error
}
type Checker struct {
	pkgScope  *Scope
	errors    []Diagnostics
	typeDecls []ast.TypeDecl
	funcDecls []*ast.FuncDecl
}
