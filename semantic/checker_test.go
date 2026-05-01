package semantic

import (
	"fmt"
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/lexer"
	"github.com/orilang/gori/parser"
	"github.com/orilang/gori/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSemantic_checker(t *testing.T) {
	t.Run("x1", func(t *testing.T) {
		data :=
			`package main
type UserID int

type User struct {
	id UserID
}
`
		scope := &Scope{
			Parent: nil,
			Symbols: map[string]*Symbol{
				"UserID": {
					Name: "UserID",
					Kind: SymType,
					Type: &NamedType{
						UnderlyingType: TInt,
					},
				},
				"User": {
					Name: "User",
					Kind: SymType,
					Type: &StructType{
						Fields: []StructField{
							{
								Name: "id",
								Type: &NamedType{
									Name:           "UserID",
									UnderlyingType: TInt,
								},
							},
						},
					},
				},
			},
		}

		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.Equal(t, 0, len(parser.Errors))
		check := NewChecker()

		assert.Equal(t, 0, len(check.Check(pr)))
		assert.Equal(t, scope.Parent, check.pkgScope.Parent)
		assert.Equal(t, scope.Symbols["UserID"].Name, check.pkgScope.Symbols["UserID"].Name)
		assert.Equal(t, scope.Symbols["UserID"].Kind, check.pkgScope.Symbols["UserID"].Kind)
		xx := check.pkgScope.Symbols["UserID"].Type.(*NamedType)
		assert.Equal(t, TInt, xx.UnderlyingType)

		assert.Equal(t, scope.Symbols["User"].Name, check.pkgScope.Symbols["User"].Name)
		assert.Equal(t, scope.Symbols["User"].Kind, check.pkgScope.Symbols["User"].Kind)
		st := check.pkgScope.Symbols["User"].Type.(*StructType)
		xy := st.Fields[0].Type.(*NamedType)
		assert.Equal(t, TInt, TInt, xy.UnderlyingType)
	})

	t.Run("x1_duplicate", func(t *testing.T) {
		data :=
			`package main
type UserID int
type UserID int

type User struct {
	id UserID
	id UserID
}
`
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.Equal(t, 0, len(parser.Errors))
		check := NewChecker()
		assert.Greater(t, len(check.Check(pr)), 0)
	})

	t.Run("x2", func(t *testing.T) {
		data :=
			`package main
type UserID int
type test interface{
	x(a string, b bool) (c float64, d float32, e uint, f uint32, g uint64 )
	y(a string, b bool) (int8, int32, int64)
	z(a string, b bool) (c uint8, d float, e UserID)
}
`
		scope := &Scope{
			Parent: nil,
			Symbols: map[string]*Symbol{
				"UserID": {
					Name: "UserID",
					Kind: SymType,
					Type: &NamedType{
						UnderlyingType: TInt,
					},
				},
				"test": {
					Name: "test",
					Kind: SymType,
					Type: &InterfaceType{
						Methods: []FuncMethod{
							{
								Name: "x",
								FuncType: &FuncType{
									Params: []Param{
										{Name: "a", Type: TString},
										{Name: "b", Type: TBool},
									},
									Results: []Param{
										{Name: "c", Type: TFloat64},
										{Name: "d", Type: TFloat32},
										{Name: "e", Type: TUInt},
										{Name: "f", Type: TUInt32},
										{Name: "g", Type: TUInt64},
									},
								},
							},
							{
								Name: "y",
								FuncType: &FuncType{
									Params: []Param{
										{Name: "a", Type: TString},
										{Name: "b", Type: TBool},
									},
									Results: []Param{
										{Type: TInt8},
										{Type: TInt32},
										{Type: TInt64},
									},
								},
							},
							{
								Name: "z",
								FuncType: &FuncType{
									Params: []Param{
										{Name: "a", Type: TString},
										{Name: "b", Type: TBool},
									},
									Results: []Param{
										{Name: "c", Type: TUInt8},
										{Name: "d", Type: TFloat},
										{Name: "e", Type: &NamedType{
											UnderlyingType: TInt,
										}},
									},
								},
							},
						},
					},
				},
			},
		}

		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.Equal(t, 0, len(parser.Errors))
		check := NewChecker()

		assert.Equal(t, 0, len(check.Check(pr)))
		assert.Equal(t, scope.Parent, check.pkgScope.Parent)
		assert.Equal(t, scope.Symbols["UserID"].Name, check.pkgScope.Symbols["UserID"].Name)
		assert.Equal(t, scope.Symbols["UserID"].Kind, check.pkgScope.Symbols["UserID"].Kind)
		xx := check.pkgScope.Symbols["UserID"].Type.(*NamedType)
		assert.Equal(t, TInt, xx.UnderlyingType)

		assert.Equal(t, scope.Symbols["test"].Name, check.pkgScope.Symbols["test"].Name)
		source := scope.Symbols["test"].Type.(*InterfaceType)
		destination := check.pkgScope.Symbols["test"].Type.(*InterfaceType)
		assert.Equal(t, len(source.Methods), len(destination.Methods))

		for k := range len(source.Methods) {
			assert.Equal(t, source.Methods[k].Name, destination.Methods[k].Name)
			assert.Equal(t, source.Methods[k].FuncType.Params, destination.Methods[k].FuncType.Params)
			assert.Equal(t, len(source.Methods[k].FuncType.Results), len(destination.Methods[k].FuncType.Results))

			for r := range len(source.Methods[k].FuncType.Results) {
				if k == len(source.Methods)-1 && r == len(source.Methods[k].FuncType.Results)-1 {
					src := source.Methods[k].FuncType.Results[r].Type.(*NamedType)
					dst := destination.Methods[k].FuncType.Results[r].Type.(*NamedType)
					assert.Equal(t, src.Name, dst.Name)
					assert.Equal(t, src.UnderlyingType, dst.UnderlyingType)
				} else {
					assert.Equal(t, source.Methods[k].FuncType.Results[r].Name, destination.Methods[k].FuncType.Results[r].Name)
					assert.Equal(t, source.Methods[k].FuncType.Results[r].Type, destination.Methods[k].FuncType.Results[r].Type)
				}
			}
		}
	})

	t.Run("x2_duplicate", func(t *testing.T) {
		data :=
			`package main
type test interface{
	x(a string, b bool) (c float64, d float32, e uint, f uint32, g uint64 )
	x(a string, b bool) (c float64, d float32, e uint, f uint32, g uint64, g uint64 )
	y(a string, a string, b bool) (int8, int32, zzz)
}
`
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.Equal(t, 0, len(parser.Errors))
		check := NewChecker()
		assert.Greater(t, len(check.Check(pr)), 0)
	})

	t.Run("x3", func(t *testing.T) {
		data :=
			`package main
type Color enum {
  Red;Blue;Green;Yellow
}
`
		scope := &Scope{
			Parent: nil,
			Symbols: map[string]*Symbol{
				"Color": {
					Name: "Color",
					Kind: SymType,
					Type: &EnumType{
						Variants: []string{"Red", "Blue", "Green", "Yellow"},
					},
				},
			},
		}

		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.Equal(t, 0, len(parser.Errors))
		check := NewChecker()

		assert.Equal(t, 0, len(check.Check(pr)))
		assert.Equal(t, scope.Parent, check.pkgScope.Parent)
		src := scope.Symbols["Color"].Type.(*EnumType)
		dst := check.pkgScope.Symbols["Color"].Type.(*EnumType)
		assert.Equal(t, src.Variants, dst.Variants)
	})

	t.Run("x3_duplicate", func(t *testing.T) {
		data :=
			`package main
type Color enum {
  Red;Blue;Green;Yellow;Yellow
}
`
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.Equal(t, 0, len(parser.Errors))
		check := NewChecker()
		assert.Greater(t, len(check.Check(pr)), 0)
	})

	t.Run("x4", func(t *testing.T) {
		data :=
			`package main
type test sum {
  Circle(radius float);Rect(w float, h float);None
}
`
		scope := &Scope{
			Parent: nil,
			Symbols: map[string]*Symbol{
				"test": {
					Name: "test",
					Kind: SymType,
					Type: &SumType{
						Name: "test",
						Variants: []SumVariant{
							{
								Name: "Circle",
								Field: []Param{
									{Name: "radius", Type: TFloat},
								},
							},
							{
								Name: "Rect",
								Field: []Param{
									{Name: "w", Type: TFloat},
									{Name: "h", Type: TFloat},
								},
							},
							{Name: "None"},
						},
					},
				},
			},
		}

		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.Equal(t, 0, len(parser.Errors))
		check := NewChecker()

		assert.Equal(t, 0, len(check.Check(pr)))
		assert.Equal(t, scope.Parent, check.pkgScope.Parent)
		assert.Equal(t, scope.Symbols["test"].Name, check.pkgScope.Symbols["test"].Name)
		assert.Equal(t, scope.Symbols["test"].Kind, check.pkgScope.Symbols["test"].Kind)

		src := scope.Symbols["test"].Type.(*SumType)
		dst := check.pkgScope.Symbols["test"].Type.(*SumType)
		assert.Equal(t, len(src.Variants), len(dst.Variants))
		for k := range len(src.Variants) {
			assert.Equal(t, src.Variants[k], dst.Variants[k])
		}
	})

	t.Run("x4_duplicate", func(t *testing.T) {
		data :=
			`package main
type test sum {
  Circle(radius float);Rect(w float, h float);None;None
}
`
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.Equal(t, 0, len(parser.Errors))
		check := NewChecker()
		assert.Greater(t, len(check.Check(pr)), 0)
	})

	t.Run("type_decl_name_lookup", func(t *testing.T) {
		assert.Equal(t, "", typeDeclName(nil))
	})

	t.Run("declare_type", func(t *testing.T) {
		check := NewChecker()
		check.declareTypeSymbol(nil)
	})

	t.Run("lookup", func(t *testing.T) {
		check := NewChecker()
		assert.Nil(t, check.pkgScope.Lookup(""))
	})

	t.Run("resolve_type", func(t *testing.T) {
		check := NewChecker()
		assert.Nil(t, check.resolveType(nil))
		assert.Equal(t, TInvalid, check.resolveType(&ast.ArrayType{}))
		assert.Equal(t, TInvalid, check.resolveType(&ast.SliceType{}))
		assert.Equal(t, TInvalid, check.resolveType(&ast.MapType{}))
		assert.Equal(t, TInvalid, check.resolveType(&ast.MapType{KeyType: &ast.NamedType{Parts: []token.Token{{Kind: token.Ident, Value: "string"}}}}))
	})

	t.Run("resolve_named_type", func(t *testing.T) {
		check := NewChecker()
		assert.Nil(t, check.resolveNamedType(&ast.NamedType{Parts: []token.Token{{Kind: token.KWIf}}}))
		p := []token.Token{
			{
				Kind:  token.Ident,
				Value: "x",
			},
			{
				Kind:  token.Ident,
				Value: "y",
			},
		}
		assert.Equal(t, TInvalid, check.resolveNamedType(&ast.NamedType{Parts: p}))
	})

	t.Run("x5", func(t *testing.T) {
		data :=
			`package main
type UserID int

func ok(a UserID, b UserID) UserID {
	return a
}
`
		scope := &Scope{
			Parent: nil,
			Symbols: map[string]*Symbol{
				"UserID": {
					Name: "UserID",
					Kind: SymType,
					Type: &NamedType{
						UnderlyingType: TInt,
					},
				},
				"ok": {
					Name: "ok",
					Kind: SymFunc,
					Type: &FuncMethod{
						Name: "ok",
						FuncType: &FuncType{
							Params: []Param{
								{Name: "a", Type: &NamedType{Name: "UserID", UnderlyingType: TInt}},
								{Name: "b", Type: &NamedType{Name: "UserID", UnderlyingType: TInt}},
							},
							Results: []Param{
								{Type: &NamedType{Name: "UserID", UnderlyingType: TInt}},
							},
						},
					},
				},
			},
		}

		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.Equal(t, 0, len(parser.Errors))
		check := NewChecker()

		assert.Equal(t, 0, len(check.Check(pr)))
		assert.Equal(t, scope.Parent, check.pkgScope.Parent)
		assert.Equal(t, scope.Symbols["UserID"].Name, check.pkgScope.Symbols["UserID"].Name)
		assert.Equal(t, scope.Symbols["UserID"].Kind, check.pkgScope.Symbols["UserID"].Kind)
		xx := check.pkgScope.Symbols["UserID"].Type.(*NamedType)
		assert.Equal(t, TInt, xx.UnderlyingType)

		assert.Equal(t, scope.Symbols["ok"].Name, check.pkgScope.Symbols["ok"].Name)
		assert.Equal(t, scope.Symbols["ok"].Kind, check.pkgScope.Symbols["ok"].Kind)
		src := scope.Symbols["ok"].Type.(*FuncMethod)
		dst := check.pkgScope.Symbols["ok"].Type.(*FuncMethod)
		assert.Equal(t, src.Name, dst.Name)

		src2 := src.FuncType.Params[0].Type.(*NamedType)
		dst2 := dst.FuncType.Params[0].Type.(*NamedType)
		assert.Equal(t, src2.UnderlyingType, dst2.UnderlyingType)

		src3 := src.FuncType.Params[1].Type.(*NamedType)
		dst3 := dst.FuncType.Params[1].Type.(*NamedType)
		assert.Equal(t, src3.UnderlyingType, dst3.UnderlyingType)

		rsrc1 := src.FuncType.Results[0].Type.(*NamedType)
		rdst1 := dst.FuncType.Results[0].Type.(*NamedType)
		assert.Equal(t, rsrc1.UnderlyingType, rdst1.UnderlyingType)
	})

	t.Run("x5_duplicate", func(t *testing.T) {
		data :=
			`package main
type UserID int

func ok(a UserID, b UserID, b UserID) UserID {
	return a + b
}
func ok(a UserID, b UserID) UserID {
	return a + b
}
func okk(a UserID, b UserID) UserID {
	return c
}
`
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.Equal(t, 0, len(parser.Errors))
		check := NewChecker()
		assert.Greater(t, len(check.Check(pr)), 0)
	})

	t.Run("x6", func(t *testing.T) {
		data :=
			`package main
type User struct {
	ids []int
	ar  [5]int
  mp  map[string]string
  hmp hashmap[string]string
}
`
		scope := &Scope{
			Parent: nil,
			Symbols: map[string]*Symbol{
				"User": {
					Name: "User",
					Kind: SymType,
					Type: &StructType{
						Fields: []StructField{
							{
								Name: "ids", Type: &SliceType{Elem: TInt},
							},
							{
								Name: "ar", Type: &ArrayType{Len: 5, Elem: TInt},
							},
							{
								Name: "mp", Type: &MapType{Key: TString, Value: TString},
							},
							{
								Name: "hmp", Type: &MapType{Kind: MapHash, Key: TString, Value: TString},
							},
						},
					},
				},
			},
		}

		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.Equal(t, 0, len(parser.Errors))
		check := NewChecker()

		assert.Equal(t, 0, len(check.Check(pr)))
		assert.Equal(t, scope.Parent, check.pkgScope.Parent)
		assert.Equal(t, scope.Symbols["User"].Name, check.pkgScope.Symbols["User"].Name)
		assert.Equal(t, scope.Symbols["User"].Kind, check.pkgScope.Symbols["User"].Kind)
		src := scope.Symbols["User"].Type.(*StructType)
		dst := check.pkgScope.Symbols["User"].Type.(*StructType)

		assert.Equal(t, src.Fields[0].Name, dst.Fields[0].Name)
		ssrc := src.Fields[0].Type.(*SliceType)
		sdst := src.Fields[0].Type.(*SliceType)
		assert.Equal(t, ssrc.Elem, sdst.Elem)

		assert.Equal(t, src.Fields[1].Name, dst.Fields[1].Name)
		asrc := src.Fields[1].Type.(*ArrayType)
		adst := src.Fields[1].Type.(*ArrayType)
		assert.Equal(t, asrc.Len, adst.Len)
		assert.Equal(t, asrc.Elem, adst.Elem)

		assert.Equal(t, src.Fields[2].Name, dst.Fields[2].Name)
		mpsrc := src.Fields[2].Type.(*MapType)
		mpdst := src.Fields[2].Type.(*MapType)
		assert.Equal(t, mpsrc.Kind, mpdst.Kind)
		assert.Equal(t, mpsrc.Key, mpdst.Key)
		assert.Equal(t, mpsrc.Value, mpdst.Value)
	})

	t.Run("declare_const_symbol", func(t *testing.T) {
		check := NewChecker()
		x := &ast.ConstDecl{}
		check.declareConstSymbol(x)
	})

	t.Run("check_expr", func(t *testing.T) {
		check := NewChecker()
		check.checkExpr(nil)
	})

	t.Run("x7", func(t *testing.T) {
		data :=
			`package main
const a int = 1
const b float = 1.0
const c bool = true
const d string = "test"
const e int = 1+1
const f int = -1*2
const g int = 1
const h int = 1+g
func x() int {
  return 1
}
const i int = 1+x()
type UserID int
const j UserID = UserID(1)
`
		scope := &Scope{
			Parent: nil,
			Symbols: map[string]*Symbol{
				"a": {Name: "a", Kind: SymConst, Type: TInt},
				"b": {Name: "b", Kind: SymConst, Type: TFloat},
				"c": {Name: "c", Kind: SymConst, Type: TBool},
				"d": {Name: "d", Kind: SymConst, Type: TString},
				"e": {Name: "e", Kind: SymConst, Type: TInt},
				"f": {Name: "f", Kind: SymConst, Type: TInt},
				"g": {Name: "g", Kind: SymConst, Type: TInt},
				"h": {Name: "h", Kind: SymConst, Type: TInt},
				"x": {Name: "x", Kind: SymFunc, Type: &FuncMethod{Name: "x", FuncType: &FuncType{Results: []Param{{Type: TInt}}}}},
				"i": {Name: "i", Kind: SymConst, Type: TInt},
				"UserID": {
					Name: "UserID",
					Kind: SymType,
					Type: &NamedType{UnderlyingType: TInt},
				},
				"j": {Name: "j", Kind: SymConst, Type: &NamedType{UnderlyingType: TInt}},
			},
		}

		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.Equal(t, 0, len(parser.Errors))
		check := NewChecker()

		assert.Equal(t, 0, len(check.Check(pr)))
		assert.Equal(t, scope.Parent, check.pkgScope.Parent)
		assert.Equal(t, scope.Symbols["a"].Name, check.pkgScope.Symbols["a"].Name)
		assert.Equal(t, scope.Symbols["a"].Kind, check.pkgScope.Symbols["a"].Kind)
		assert.Equal(t, scope.Symbols["a"].Type, check.pkgScope.Symbols["a"].Type)

		assert.Equal(t, scope.Symbols["b"].Name, check.pkgScope.Symbols["b"].Name)
		assert.Equal(t, scope.Symbols["b"].Kind, check.pkgScope.Symbols["b"].Kind)
		assert.Equal(t, scope.Symbols["b"].Type, check.pkgScope.Symbols["b"].Type)

		assert.Equal(t, scope.Symbols["c"].Name, check.pkgScope.Symbols["c"].Name)
		assert.Equal(t, scope.Symbols["c"].Kind, check.pkgScope.Symbols["c"].Kind)
		assert.Equal(t, scope.Symbols["c"].Type, check.pkgScope.Symbols["c"].Type)

		assert.Equal(t, scope.Symbols["d"].Name, check.pkgScope.Symbols["d"].Name)
		assert.Equal(t, scope.Symbols["d"].Kind, check.pkgScope.Symbols["d"].Kind)
		assert.Equal(t, scope.Symbols["d"].Type, check.pkgScope.Symbols["d"].Type)

		assert.Equal(t, scope.Symbols["e"].Name, check.pkgScope.Symbols["e"].Name)
		assert.Equal(t, scope.Symbols["e"].Kind, check.pkgScope.Symbols["e"].Kind)
		assert.Equal(t, scope.Symbols["e"].Type, check.pkgScope.Symbols["e"].Type)

		assert.Equal(t, scope.Symbols["f"].Name, check.pkgScope.Symbols["f"].Name)
		assert.Equal(t, scope.Symbols["f"].Kind, check.pkgScope.Symbols["f"].Kind)
		assert.Equal(t, scope.Symbols["f"].Type, check.pkgScope.Symbols["f"].Type)

		assert.Equal(t, scope.Symbols["g"].Name, check.pkgScope.Symbols["g"].Name)
		assert.Equal(t, scope.Symbols["g"].Kind, check.pkgScope.Symbols["g"].Kind)
		assert.Equal(t, scope.Symbols["g"].Type, check.pkgScope.Symbols["g"].Type)

		assert.Equal(t, scope.Symbols["h"].Name, check.pkgScope.Symbols["h"].Name)
		assert.Equal(t, scope.Symbols["h"].Kind, check.pkgScope.Symbols["h"].Kind)
		assert.Equal(t, scope.Symbols["h"].Type, check.pkgScope.Symbols["h"].Type)

		assert.Equal(t, scope.Symbols["x"].Name, check.pkgScope.Symbols["x"].Name)
		assert.Equal(t, scope.Symbols["x"].Kind, check.pkgScope.Symbols["x"].Kind)
		fsrc := scope.Symbols["x"].Type.(*FuncMethod)
		fdst := check.pkgScope.Symbols["x"].Type.(*FuncMethod)
		assert.Equal(t, fsrc.Name, fdst.Name)
		assert.Equal(t, fsrc.FuncType.Params, fdst.FuncType.Params)
		assert.Equal(t, fsrc.FuncType.Results[0].Type, fdst.FuncType.Results[0].Type)

		assert.Equal(t, scope.Symbols["i"].Name, check.pkgScope.Symbols["i"].Name)
		assert.Equal(t, scope.Symbols["i"].Kind, check.pkgScope.Symbols["i"].Kind)
		assert.Equal(t, scope.Symbols["i"].Type, check.pkgScope.Symbols["i"].Type)

		assert.Equal(t, scope.Symbols["UserID"].Name, check.pkgScope.Symbols["UserID"].Name)
		assert.Equal(t, scope.Symbols["UserID"].Kind, check.pkgScope.Symbols["UserID"].Kind)
		xx := check.pkgScope.Symbols["UserID"].Type.(*NamedType)
		assert.Equal(t, TInt, xx.UnderlyingType)

		assert.Equal(t, scope.Symbols["j"].Name, check.pkgScope.Symbols["j"].Name)
		assert.Equal(t, scope.Symbols["j"].Kind, check.pkgScope.Symbols["j"].Kind)

		src := scope.Symbols["j"].Type.(*NamedType)
		dst := check.pkgScope.Symbols["j"].Type.(*NamedType)
		assert.Equal(t, src.UnderlyingType, dst.UnderlyingType)
	})

	t.Run("x7_error", func(t *testing.T) {
		data :=
			`package main
const a float = "test"
const b int = 1+g
func x() string {
	return "string"
}
func y() (int,string) {
	return 1
}
func z(a int, b int) int {
	return a+b
}
const c int = 1+x()
const cc int = x()+1
const d int = 1+y()
const dd int = y()+1
const e int = z(1,"a"+1)
const e1 int = z(1,"a"+1,5)
const e2 int = z(1,1.1)
const e3 int = zz()
type UserID int
type User int
const l UserID = User(1)
const l1 UserID = User(1,1)
const l2 UserID = User(UserID(1))
const l3 UserID = UserID(1) + User(1)
const l4 UserID = UserID(true)
const m int = !1
`
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.Equal(t, 0, len(parser.Errors))
		check := NewChecker()
		assert.Greater(t, len(check.Check(pr)), 0)
	})

	t.Run("x7_duplicate", func(t *testing.T) {
		data :=
			`package main
const a int = 1
const a float = 1.0
`
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.Equal(t, 0, len(parser.Errors))
		check := NewChecker()
		assert.Greater(t, len(check.Check(pr)), 0)
	})

	t.Run("x8", func(t *testing.T) {
		data :=
			`package main
func fa() {
	const a int = 0
	var b int = 0
	b = 1
	x := b+1
	x++
	x = x*(b+10)
	var s string = ""
	const sl []string = []string{"a","b","c"}
	s = sl[0]
	s = sl[b]
	ss := sl[0]
	const ar [5]int = [5]int{1,2,3}
	b = ar[0]
	type User int
	const l User = User(1)
	ll := User(1)
	type test sum {
		Circle(radius float);Rect(w float, h float);None
	}
	type Color enum {
		Red;Blue;Green;Yellow
	}
	type testi interface{
		x(a string, b bool) (c float64, d float32, e uint, f uint32, g uint64 )
	}
	type UserID struct {
		id string
	}
	sstring := "ori"
	ssbool := true
	fempty()
}
func fempty() {}
`
		scope := &Scope{
			Parent: nil,
			Symbols: map[string]*Symbol{
				"fa": {
					Name: "fa",
					Kind: SymFunc,
					Type: &FuncMethod{
						Name:     "fa",
						FuncType: &FuncType{},
					},
				},
				"a": {Name: "a", Kind: SymConst, Type: TInt},
				"b": {Name: "b", Kind: SymVar, Type: TInt},
				"x": {Name: "x", Kind: SymVar, Type: TInt},
				"s": {Name: "s", Kind: SymVar, Type: TString},
				"sl": {
					Name: "sl",
					Kind: SymConst,
					Type: &SliceType{Elem: TString},
				},
				"ar": {
					Name: "ar",
					Kind: SymConst,
					Type: &SliceType{Elem: TString},
				},
				"j": {Name: "j", Kind: SymConst, Type: &NamedType{UnderlyingType: TInt}},
				"fempty": {
					Name: "fempty",
					Kind: SymFunc,
					Type: &FuncMethod{
						Name:     "fempty",
						FuncType: &FuncType{},
					},
				},
			},
		}

		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.Equal(t, 0, len(parser.Errors))
		check := NewChecker()

		x := check.Check(pr)
		for _, v := range x {
			fmt.Println("BBBB", v.Err.Error())
		}
		assert.Equal(t, 0, len(x))
		// assert.Equal(t, 0, len(check.Check(pr)))
		assert.Equal(t, scope.Parent, check.pkgScope.Parent)
		fmt.Printf("check.pkgScope.Symbols %#v\n", check.pkgScope.Symbols)

		assert.Equal(t, scope.Symbols["fa"].Name, check.pkgScope.Symbols["fa"].Name)
		assert.Equal(t, scope.Symbols["fa"].Kind, check.pkgScope.Symbols["fa"].Kind)
		fsrc := scope.Symbols["fa"].Type.(*FuncMethod)
		fdst := check.pkgScope.Symbols["fa"].Type.(*FuncMethod)
		assert.Equal(t, fsrc.Name, fdst.Name)
		assert.Equal(t, len(fsrc.FuncType.Params), len(fdst.FuncType.Params))
		assert.Equal(t, len(fdst.FuncType.Results), len(fdst.FuncType.Results))

		assert.Equal(t, scope.Symbols["fempty"].Name, check.pkgScope.Symbols["fempty"].Name)
		assert.Equal(t, scope.Symbols["fempty"].Kind, check.pkgScope.Symbols["fempty"].Kind)
		f2src := scope.Symbols["fempty"].Type.(*FuncMethod)
		f2dst := check.pkgScope.Symbols["fempty"].Type.(*FuncMethod)
		assert.Equal(t, f2src.Name, f2dst.Name)
		assert.Equal(t, len(f2src.FuncType.Params), len(f2dst.FuncType.Params))
		assert.Equal(t, len(f2dst.FuncType.Results), len(f2dst.FuncType.Results))
	})

	t.Run("x8_error", func(t *testing.T) {
		// 		tests := []struct {
		// 			data string
		// 			err  bool
		// 		}{
		// 			{
		// 				data: `package main
		// func fa() {
		// 	const a int = 0
		// 	var b int = 0
		// 	b = 1
		// 	x := b+1
		// 	x++
		// }`,
		// 			},
		// 			{
		// 				data: `package main
		// func fa() {
		// 	var s string = ""
		// 	const sl []string = []string{"a","b","c"}
		// 	s = sl[0]
		// 	s = sl[b]
		// 	ss := sl[0]
		// 	const ar [5]int = [5]int{1,2,3}
		// 	b = ar[0]
		// }`,
		// 			},
		// 		}

		// 		for _, tc := range tests {
		// 			lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		// 			require.NoError(t, err)
		// 			parser := parser.New(lex.FetchTokensFromString(tc.data))
		// 			pr := parser.ParseFile()
		// 			assert.Equal(t, 0, len(parser.Errors))
		// 			check := NewChecker()
		// 			assert.Greater(t, len(check.Check(pr)), 0)
		// 		}

		data :=
			`package main
const zz int = 0
func yy() (int, int){
  return 2,3
}
func y() {
	const a int = 0
	a++
	a = 5
  var b int = 0
	b = true
	bb := b
	bb := b
	z = 1
	x := 1
	x := xx
	xx++
	var s string = "a"
	s++
	const sl []string = []string{"a","b","c"}
	s = sl[s]
	const ar [5]int = [5]int{1,2,3}
	b = ar[s]
	y3 := yy()
	y4 := 1+2
	y5 := (1 + 2) + 3
	y6 := 1 + (2 + 3)
	y7 := -1
	y8 := (1)
	type User int
	y9 := User(1) + 1
  b := y9 + User(1)
	User(1)
	z1.b := "c"
}
`
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.Equal(t, 0, len(parser.Errors))
		check := NewChecker()
		assert.Greater(t, len(check.Check(pr)), 0)
		check.checkFuncBody(&ast.FuncDecl{Name: token.Token{Value: "zz"}})
	})

	t.Run("x8_duplicate", func(t *testing.T) {
		data :=
			`package main
func y() {
  const a int = 0
  const a int = 0
  const aa int = true
  var b int = 0
  var b int = 0
  var c int = true
	type User int
	type User int
	type Color enum {
    Red;Blue;Green;Yellow
  }
	type Color enum {
		Red;Blue;Green;Yellow
	}
	type test interface{
		x(a string, b bool) (c float64, d float32, e uint, f uint32, g uint64 )
	}
	type test interface{
		x(a string, b bool) (c float64, d float32, e uint, f uint32, g uint64 )
	}
	type test sum {
		Circle(radius float);Rect(w float, h float);None
	}
	type test sum {
		Circle(radius float);Rect(w float, h float);None
	}
	type UserID struct {
		id string
	}
	type UserID struct {
		id string
	}
}
`
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		assert.Equal(t, 0, len(parser.Errors))
		check := NewChecker()
		assert.Greater(t, len(check.Check(pr)), 0)
	})

	t.Run("check_block_stmt", func(t *testing.T) {
		check := NewChecker()
		check.checkBlockStmt(nil)
	})

	t.Run("check_func_body", func(t *testing.T) {
		check := NewChecker()
		check.checkFuncBody(&ast.FuncDecl{})
		check.checkReturnStmt(nil)
		check.checkExprStmt(&ast.ExprStmt{
			Expr: &ast.IntLitExpr{
				Name: token.Token{
					Value: "0",
				},
			},
		})

		check.checkSimpleAssignStmt(
			&ast.AssignStmt{
				Left: &ast.IntLitExpr{
					Name: token.Token{
						Value: "0",
					},
				},
				Right: &ast.StringLitExpr{
					Name: token.Token{
						Value: "test",
					},
				},
			},
		)

		check.checkBlockStmt(&ast.BlockStmt{
			Stmts: []ast.Stmt{
				&ast.BadStmt{},
				&ast.DeclStmt{
					Decl: &ast.BadDecl{},
				},
			},
		})

		check.checkExpr(&ast.IndexExpr{
			X:     &ast.IdentExpr{Name: token.Token{Value: "a"}},
			Index: &ast.BadExpr{},
		})

		check.typeDecls = append(check.typeDecls, &ast.StructDecl{})
		check.createTypeObjects()
		check.resolveTypeDecls()

		check.checkIncDecStmt(&ast.IncDecStmt{X: &ast.BadExpr{}})
	})
}
