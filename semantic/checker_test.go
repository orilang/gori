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
		require.Equal(t, 0, len(parser.Errors))
		check := NewChecker()

		assert.Equal(t, 0, len(check.Check(pr)))
		assert.Equal(t, scope.Symbols["UserID"].Name, check.pkgScope.Symbols["UserID"].Name)
		assert.Equal(t, scope.Symbols["UserID"].Kind, check.pkgScope.Symbols["UserID"].Kind)
		xx := check.pkgScope.Symbols["UserID"].Type.(*NamedType)
		assert.Equal(t, TInt, xx.UnderlyingType)

		assert.Equal(t, scope.Symbols["User"].Name, check.pkgScope.Symbols["User"].Name)
		assert.Equal(t, scope.Symbols["User"].Kind, check.pkgScope.Symbols["User"].Kind)
		stn := check.pkgScope.Symbols["User"].Type.(*NamedType)
		st := stn.UnderlyingType.(*StructType)
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
		require.Equal(t, 0, len(parser.Errors))
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
											Name:           "UserID",
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
		require.Equal(t, 0, len(parser.Errors))
		check := NewChecker()

		assert.Equal(t, 0, len(check.Check(pr)))
		assert.Equal(t, scope.Symbols["UserID"].Name, check.pkgScope.Symbols["UserID"].Name)
		assert.Equal(t, scope.Symbols["UserID"].Kind, check.pkgScope.Symbols["UserID"].Kind)
		xx := check.pkgScope.Symbols["UserID"].Type.(*NamedType)
		assert.Equal(t, TInt, xx.UnderlyingType)

		assert.Equal(t, scope.Symbols["test"].Name, check.pkgScope.Symbols["test"].Name)
		source := scope.Symbols["test"].Type.(*InterfaceType)
		destinationN := check.pkgScope.Symbols["test"].Type.(*NamedType)
		destination := destinationN.UnderlyingType.(*InterfaceType)
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
		require.Equal(t, 0, len(parser.Errors))
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
		require.Equal(t, 0, len(parser.Errors))
		check := NewChecker()

		assert.Equal(t, 0, len(check.Check(pr)))
		src := scope.Symbols["Color"].Type.(*EnumType)
		dstE := check.pkgScope.Symbols["Color"].Type.(*NamedType)
		dst := dstE.UnderlyingType.(*EnumType)
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
		require.Equal(t, 0, len(parser.Errors))
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
		require.Equal(t, 0, len(parser.Errors))
		check := NewChecker()

		assert.Equal(t, 0, len(check.Check(pr)))
		assert.Equal(t, scope.Symbols["test"].Name, check.pkgScope.Symbols["test"].Name)
		assert.Equal(t, scope.Symbols["test"].Kind, check.pkgScope.Symbols["test"].Kind)

		src := scope.Symbols["test"].Type.(*SumType)
		dstS := check.pkgScope.Symbols["test"].Type.(*NamedType)
		dst := dstS.UnderlyingType.(*SumType)
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
		require.Equal(t, 0, len(parser.Errors))
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
		require.Equal(t, 0, len(parser.Errors))
		check := NewChecker()

		assert.Equal(t, 0, len(check.Check(pr)))
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
		require.Equal(t, 0, len(parser.Errors))
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
								Name: "hmp", Type: &HashMapType{Key: TString, Value: TString},
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
		require.Equal(t, 0, len(parser.Errors))
		check := NewChecker()

		assert.Equal(t, 0, len(check.Check(pr)))
		assert.Equal(t, scope.Symbols["User"].Name, check.pkgScope.Symbols["User"].Name)
		assert.Equal(t, scope.Symbols["User"].Kind, check.pkgScope.Symbols["User"].Kind)
		src := scope.Symbols["User"].Type.(*StructType)
		dstS := check.pkgScope.Symbols["User"].Type.(*NamedType)
		dst := dstS.UnderlyingType.(*StructType)

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
		assert.Equal(t, mpsrc.Key, mpdst.Key)
		assert.Equal(t, mpsrc.Value, mpdst.Value)

		assert.Equal(t, src.Fields[3].Name, dst.Fields[3].Name)
		hmpsrc := src.Fields[3].Type.(*HashMapType)
		hmpdst := src.Fields[3].Type.(*HashMapType)
		assert.Equal(t, hmpsrc.Key, hmpdst.Key)
		assert.Equal(t, hmpsrc.Value, hmpdst.Value)
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
		require.Equal(t, 0, len(parser.Errors))
		check := NewChecker()

		assert.Equal(t, 0, len(check.Check(pr)))
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
		require.Equal(t, 0, len(parser.Errors))
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
		require.Equal(t, 0, len(parser.Errors))
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
		require.Equal(t, 0, len(parser.Errors))
		check := NewChecker()

		assert.Equal(t, 0, len(check.Check(pr)))
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
		require.Equal(t, 0, len(parser.Errors))
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
		require.Equal(t, 0, len(parser.Errors))
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

		check.checkAssigmentStmt(&ast.AssignStmt{
			Operator: token.Token{Kind: token.Slash},
		})

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
		check.checkAssignableExpr(&ast.IndexExpr{
			X:     &ast.IdentExpr{Name: token.Token{Value: "a"}},
			Index: &ast.BadExpr{},
		})
		check.checkAssignableExpr(&ast.BadExpr{})

		check.typeDecls = append(check.typeDecls, &ast.StructDecl{})
		check.createTypeObjects()
		check.resolveTypeDecls()
		check.checkIncDecStmt(&ast.IncDecStmt{X: &ast.BadExpr{}})

		check.resolveType(&ast.MapType{
			KeyType: &ast.NamedType{
				Parts: []token.Token{
					{
						Kind:  token.KWInt,
						Value: "int",
					},
				},
			},
			ValueType: &ast.NamedType{
				Parts: []token.Token{
					{
						Kind:  token.StringLit,
						Value: "zzzz",
					},
				},
			},
		})
	})

	t.Run("x9", func(t *testing.T) {
		tests := []struct {
			data string
			err  bool
		}{
			{
				data: `package main
type User struct {
	name string
	age  int
}
func f(u User) string {
	return u.name
}
`,
			},
			{
				data: `package main
type User struct {
	name string
	age  int
}
func age(u User) {
	u.age = 10
}
`,
			},
			{
				err: true,
				data: `package main
type User struct {
		name string
}
func bad(u User) {
		u.name = 1
}
`,
			},
			{
				err: true,
				data: `package main
type User struct {
		name string
}
func bad(u User) string {
		return u.age
}
`,
			},
			{
				err: true,
				data: `package main
type User struct {
		name string
}
func bad(u User) string {
		return u.unknown
}
`,
			},
			{
				data: `package main
type test interface {
	foo() string
}
func f(u test) string {
	return u.foo()
}
			`,
			},
			{
				data: `package main
type test interface {
	foo() string
}
func f(u test) {
	u.foo()
}
`,
			},
			{
				err: true,
				data: `package main
type test interface {
	foo() string
}
func f(u test) string {
	return u.unknown()
}
			`,
			},
			{
				err: true,
				data: `package main
type test interface {
	foo() (string,string)
}
func f(u test) string {
	return u.foo()
}
			`,
			},
			{
				err: true,
				data: `package main
type test interface {
	foo() string
}
func f(u test) {
	u.unknown()
}
`,
			},
			{
				data: `package main
func foo() string {
	return "foo"
}
func bar() string {
	return foo()
}
`,
			},
			{
				err: true,
				data: `package main
func bar() string {
	return foo()
}
`,
			},
			{
				data: `package main
func f(s []string) string {
	return s[0]
}
func x(s [5]string) string {
	return s[0]
}
			`,
			},
			{
				data: `package main
func f(m map[string]string) string {
	return m["x"]
}
`,
			},
			{
				data: `package main
func f(m map[string]string) string {
	m["k"] = "v"
}
`,
			},
			{
				data: `package main
func f(m map[string]string) map[string]string {
	return m
}
`,
			},
			{
				data: `package main
type UsersByID map[int]string

func f(m UsersByID) string {
	return m[1]
}
`,
			},
			{
				data: `package main
func f(m hashmap[string]string) string {
	m["k"] = "v"
}
`,
			},
			{
				data: `package main
func f(m hashmap[string]string) hashmap[string]string {
	return m
}
`,
			},
			{
				err: true,
				data: `package main
func f(m map[string]string) string {
	return m[0]
}
`,
			},
			{
				err: true,
				data: `package main
func f(m map[string]string) hashmap[string]string {
	return m
}
`,
			},
			{
				err: true,
				data: `package main
func f(m hashmap[string]string) string {
	return m[0]
}
`,
			},
		}

		for _, tc := range tests {
			lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
			require.NoError(t, err)
			parser := parser.New(lex.FetchTokensFromString(tc.data))
			pr := parser.ParseFile()
			require.Equal(t, 0, len(parser.Errors))
			check := NewChecker()

			result := check.Check(pr)
			if tc.err {
				assert.Greater(t, len(result), 0)
			} else {
				assert.Equal(t, 0, len(result))
			}
		}

		check := NewChecker()
		assert.Equal(t, TInvalid, check.checkSelectorExpr(&ast.SelectorExpr{X: &ast.BadExpr{}}))
	})

	t.Run("x10", func(t *testing.T) {
		tests := []struct {
			data string
			err  bool
		}{
			{
				data: `package main
type User [5]string
func f(u User) string {
	return u[1]
}
`,
			},
			{
				data: `package main
type User [1+2]string
`,
			},
			{
				data: `package main
type User [2-1]string
`,
			},
			{
				data: `package main
type User [2/1]string
`,
			},
			{
				data: `package main
type User [1*2]string
`,
			},
			{
				data: `package main
type User [+2]string
`,
			},
			{
				err: true,
				data: `package main
type User [1+2.4]string
`,
			},
			{
				data: `package main
type User [-(-2)]string
`,
			},
			{
				data: `package main
type User [- -1]string
`,
			},
			{
				err: true,
				data: `package main
type User [-1]string
`,
			},
			{
				err: true,
				data: `package main
type User [2/0]string
`,
			},
			{
				err: true,
				data: `package main
type User [-1.2]string
`,
			},
			{
				err: true,
				data: `package main
type User [1.2]string
`,
			},
			{
				err: true,
				data: `package main
type User [+1.2]string
`,
			},
			{
				err: true,
				data: `package main
type User [-(1+2)]string
`,
			},
			{
				err: true,
				data: `package main
type User [-a]string
`,
			},
			{
				err: true,
				data: `package main
type User [!true]string
`,
			},
		}

		for i, tc := range tests {
			lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
			require.NoError(t, err)
			parser := parser.New(lex.FetchTokensFromString(tc.data))
			pr := parser.ParseFile()
			require.Equal(t, 0, len(parser.Errors))
			check := NewChecker()

			result := check.Check(pr)
			if tc.err {
				assert.Greater(t, len(result), 0, i)
			} else {
				assert.Equal(t, 0, len(result), i)
			}
		}

		check := NewChecker()
		_, ok := check.evalArrayLen(&ast.IntLitExpr{Name: token.Token{Value: "a"}})
		assert.Equal(t, false, ok)
	})

	t.Run("x11", func(t *testing.T) {
		tests := []struct {
			data string
			err  bool
		}{
			{
				data: `package main
func x(a int, b int, c int) int {
 if a < b {
   if a < c {
	   return c
	 }
 } else {
   return b
 }
}
`,
			},
			{
				data: `package main
func x(a int, b int, c int) int {
 if a < b {
   return b
 } else if a < c {
	 return c
 } else {
   return a
 }
}
`,
			},
			{
				err: true,
				data: `package main
func x(a int, b int, c int) int {
 if a {
   return b
 }
}
`,
			},
		}

		for i, tc := range tests {
			lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
			require.NoError(t, err)
			parser := parser.New(lex.FetchTokensFromString(tc.data))
			pr := parser.ParseFile()
			require.Equal(t, 0, len(parser.Errors))
			check := NewChecker()

			result := check.Check(pr)
			if tc.err {
				assert.Greater(t, len(result), 0, i)
			} else {
				assert.Equal(t, 0, len(result), i)
			}
		}

		check := NewChecker()
		check.checkIfStmt(nil)
		check.scope = NewScope(check.pkgScope)
		check.useScope = true
		check.scope.Declare(&Symbol{
			Name: "a",
			Kind: SymVar,
			Type: TInt,
		})
		check.scope.Declare(&Symbol{
			Name: "b",
			Kind: SymVar,
			Type: TInt,
		})
		check.checkIfStmt(&ast.IfStmt{
			Condition: &ast.BinaryExpr{
				Left:     &ast.IdentExpr{Name: token.Token{Kind: token.Ident, Value: "a"}},
				Operator: token.Token{Kind: token.Lt, Value: "<"},
				Right:    &ast.IdentExpr{Name: token.Token{Kind: token.Ident, Value: "b"}},
			},
			Else: &ast.BadStmt{},
		})
		check.useScope = false
	})

	t.Run("x12", func(t *testing.T) {
		tests := []struct {
			data string
			err  bool
		}{
			{
				data: `package main
func x() bool {
	for {
		return false
	}
}
`,
			},
			{
				data: `package main
func x(a int, b int) bool {
	for a < b {
		return false
	}
}
`,
			},
			{
				data: `package main
func x(a int) int {
	for a = 0;a<5;a+=1 {
		return a
	}
}
`,
			},
			{
				data: `package main
func x() int {
 for a := int(0);a<5;a+=1 {
   return a
 }
}
`,
			},
			{
				err: true,
				data: `package main
func x() int {
	for a := int("a");a<5;a+=1 {
		return a
	}
}
`,
			},
			{
				err: true,
				data: `package main
func x() int {
	for a := int(1,2);a<5;a+=1 {
		return a
	}
}
`,
			},
		}

		for i, tc := range tests {
			lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
			require.NoError(t, err)
			parser := parser.New(lex.FetchTokensFromString(tc.data))
			pr := parser.ParseFile()
			require.Equal(t, 0, len(parser.Errors))
			check := NewChecker()

			result := check.Check(pr)
			if tc.err {
				assert.Greater(t, len(result), 0, i)
			} else {
				assert.Equal(t, 0, len(result), i)
			}
		}

		check := NewChecker()
		check.checkForStmt(nil)
		check.checkForStmt(&ast.ForStmt{
			Init: &ast.BadStmt{},
		})
		check.checkForStmt(&ast.ForStmt{
			Condition: &ast.BadExpr{},
		})
		check.checkForStmt(&ast.ForStmt{
			Post: &ast.BadStmt{},
		})
	})

	t.Run("x13", func(t *testing.T) {
		tests := []struct {
			data string
			err  bool
		}{
			{
				data: `package main
func x(z int){
	for i := range int(5) {
		z = i
	}
}
`,
			},
			{
				data: `package main
func x(){
	for k,v := range int(5) {}
}
`,
			},
			{
				data: `package main
func x(){
	for _,v := range int(5) {}
}
`,
			},
			{
				data: `package main
func x(){
	for k,_ := range int(5) {}
}
`,
			},
			{
				data: `package main
func x(){
	var k int = 0
	var v int = 0
	for kk,vv := range int(5) {}
}
`,
			},
			{
				err: true,
				data: `package main
func x(){
	for _,_ := range int(5) {}
}
`,
			},
			{
				err: true,
				data: `package main
func x(){
	for _ := range int(5) {}
}
`,
			},
			{
				err: true,
				data: `package main
func x(z string){
	for z = range 5 {
		z = i
	}
}
`,
			},
			{
				data: `package main
func x(z int){
	for z = range int(5) {}
}
`,
			},
			{
				err: true,
				data: `package main
func x(z string){
	for _ = range 5 {}
}
`,
			},
			{
				err: true,
				data: `package main
func x(z string){
	for z = range int(5) {
		z = i
	}
}
`,
			},
			{
				err: true,
				data: `package main
func x(z string){
	for z,_ = range int(5) {
		z = i
	}
}
`,
			},
			{
				err: true,
				data: `package main
func x(z string){
	for _,z = range int(5) {
		z = i
	}
}
`,
			},
			{
				err: true,
				data: `package main
func x(z int){
	for z,z := range int(5) {
		z = i
	}
}
`,
			},
			{
				err: true,
				data: `package main
func x(z int){
	for k,v := range float(5) {
		z = i
	}
}
`,
			},
			{
				data: `package main
func x(z []int){
	for k,v := range z {}
}
`,
			},
			{
				data: `package main
func x(z [5]int){
	for k,v := range z {}
}
`,
			},
			{
				data: `package main
func x(z map[string]int){
	for k,v := range z {}
}
`,
			},
			{
				data: `package main
func x(z hashmap[string]int){
	for k,v := range z {}
}
`,
			},
			{
				data: `package main
func x(z hashmap[string]int){
	for range z {}
}
`,
			},
			{
				data: `package main
func x(z hashmap[string]int){
	for range z {
	  break
	}
}
`,
			},
			{
				data: `package main
func x(z hashmap[string]int){
	for range z {
	  break
		// special case that will be handled by unreachable stmt
		a := int(0)
	}
}
`,
			},
			{
				data: `package main
func x(z hashmap[string]int){
	for k,v := range z {
	  break
	}
}
`,
			},
			{
				data: `package main
func x(z hashmap[string]int){
	for k,v := range z {
	  break
		// special case that will be handled by unreachable stmt
		a := int(0)
	}
}
`,
			},
			{
				err: true,
				data: `package main
func x(z hashmap[string]int){
	break
}
`,
			},
			{
				data: `package main
func x(){
 a := true
 for {
   if a {
	   break
	 }
 }
}
`,
			},
			{
				err: true,
				data: `package main
func f(x int) {
	switch x {
	case 1:
		break
	}
}
`,
			},
			{
				data: `package main
func x(z hashmap[string]int){
	for k,v := range z {
	  continue
	}
}
`,
			},
			{
				data: `package main
func x(z hashmap[string]int){
	for k,v := range z {
	  continue
		// special case that will be handled by unreachable stmt
		a := int(0)
	}
}
`,
			},
			{
				data: `package main
func x(){
 a := true
 for {
   if a {
	   continue
	 }
 }
}
`,
			},
			{
				err: true,
				data: `package main
func f(x int) {
	switch x {
	case 1:
		continue
	}
}
`,
			},
		}

		for i, tc := range tests {
			lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
			require.NoError(t, err)
			parser := parser.New(lex.FetchTokensFromString(tc.data))
			pr := parser.ParseFile()
			require.Equal(t, 0, len(parser.Errors))
			check := NewChecker()

			result := check.Check(pr)
			if tc.err {
				assert.Greater(t, len(result), 0, i)
			} else {
				assert.Equal(t, 0, len(result), i)
			}
		}

		check := NewChecker()
		check.checkRangeStmt(nil)
		check.checkRangeStmt(&ast.RangeStmt{X: &ast.BadExpr{}})
		_, _, ok := rangeVars(TInvalid)
		assert.Equal(t, false, ok)

		check.scope = NewScope(check.pkgScope)
		check.useScope = true
		check.checkRangeStmt(&ast.RangeStmt{
			X: &ast.CallExpr{
				Callee: &ast.IdentExpr{
					Name: token.Token{Kind: token.Ident, Value: "int"},
				},
				Args: []ast.Expr{
					&ast.IntLitExpr{Name: token.Token{Kind: token.IntLit, Value: "5"}},
				},
			},
			Op: token.Token{Kind: token.Slash},
		})
		check.useScope = false
	})

	t.Run("x14", func(t *testing.T) {
		tests := []struct {
			data string
			err  bool
		}{
			{
				err: true,
				data: `package main
const x int = int(1)

func f() {
	var x int = int(2)
}
`,
			},
			{
				err: true,
				data: `package main
const k int = int(1)

func x(zz hashmap[string]int){
	for k,v := range zz {}
}
`,
			},
			{
				err: true,
				data: `package main
func x(zz hashmap[string]int){
  const ca int = int(0)
  const ca int = int(0)
}
`,
			},
			{
				err: true,
				data: `package main
func x(zz hashmap[string]int){
  var va int = int(0)
  var va int = int(0)
}
`,
			},
			{
				err: true,
				data: `package main
func x(zz hashmap[string]int,zz hashmap[string]int){
  assign := int(0)
  assign := int(0)
}
`,
			},
			{
				err: true,
				data: `package main
const k int = int(1)

func x(zz hashmap[string]int){
	for k := range zz {}
}
`,
			},
			{
				err: true,
				data: `package main
const v int = int(1)

func x(zz hashmap[string]int){
	for k,v := range zz {}
}
`,
			},
		}

		for i, tc := range tests {
			lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
			require.NoError(t, err)
			parser := parser.New(lex.FetchTokensFromString(tc.data))
			pr := parser.ParseFile()
			require.Equal(t, 0, len(parser.Errors))
			check := NewChecker()

			result := check.Check(pr)
			if tc.err {
				assert.Greater(t, len(result), 0, i)
			} else {
				assert.Equal(t, 0, len(result), i)
			}
		}
	})

	t.Run("x15", func(t *testing.T) {
		tests := []struct {
			data string
			err  bool
		}{
			{
				data: `package main
func w() int {
  return 5
}
func b() {}
func c() {}
func x(a int) int {
  switch z:=w();z {
    case a:
      b()
    case 2:
      c()
	}
}
`,
			},
			{
				err: true,
				data: `package main
func w() int {
  return 5
}
func b() {}
func c() {}
func x(a int) int {
  switch z:=w();z {
    case z>a:
      b()
    case 2:
      c()
	}
}
`,
			},
			{
				err: true,
				data: `package main
func x(a int) int {
  switch z {
    case z>a:
      b()
    case 2:
      c()
	}
}
`,
			},
			{
				err: true,
				data: `package main
func w() int {
  return 5
}
func x(a int) int {
  switch z:=w();z {
    case z>a:
      b()
    case 2:
      c()
	}
}
`,
			},
			{
				err: true,
				data: `package main
func w() int {
  return 5
}
func x(a string) int {
  switch z:=w();z {
    case a:
      b()
    case 2:
      c()
	}
}
`,
			},
			{
				err: true,
				data: `package main
func w() int {
  return 5
}
func x(z map[string]string) int {
  switch z {
    case 1:
      b()
	}
}
`,
			},
			{
				err: true,
				data: `package main
func w() int {
  return 5
}
func b() {}
func c() {}
func x(a int) int {
  switch z:=w();z {}
}
`,
			},
			{
				data: `package main
func x() int {
  switch {
    case 1 == 2:
      return 0
	}
}
`,
			},
			{
				err: true,
				data: `package main
func x() int {
  switch {
    default:
    default:
    case 1 == 2:
      return 0
	}
}
`,
			},
			{
				err: true,
				data: `package main
func x() int {
  switch {
    case "a":
      return 0
	}
}
`,
			},
			{
				err: true,
				data: `package main
func x() int {
  switch {}
}
`,
			},
			{
				data: `package main
func doA() {}
func doB() {}
func x(a int) int {
  switch a {
	case 1:
		doA()
		fallthrough
	case 2:
		doB()
	}
}
`,
			},
			{
				err: true,
				data: `package main
func doA() {}
func doB() {}
func x(a int) int {
  switch a {
	case 1:
		doA()
		fallthrough
	case 2:
		doB()
	default:
	default:
	}
}
`,
			},
			{
				err: true,
				data: `package main
func doA() {}
func doB() {}
func x(a int) int {
  switch a {
	case 1:
		fallthrough
		doA()
	case 2:
		doB()
	}
}
`,
			},
			{
				err: true,
				data: `package main
func doA() {}
func x(a int) int {
  switch a {
	case 1:
		doA()
	case 2:
		fallthrough
	}
}
`,
			},
			{
				err: true,
				data: `package main
func doA() {}
func x(a int) int {
  switch a {
	case 1:
		doA()
	default:
		fallthrough
	}
}
`,
			},
			{
				err: true,
				data: `package main
func doA() {
  fallthrough
}
`,
			},
		}

		for i, tc := range tests {
			lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
			require.NoError(t, err)
			parser := parser.New(lex.FetchTokensFromString(tc.data))
			pr := parser.ParseFile()
			require.Equal(t, 0, len(parser.Errors))
			check := NewChecker()

			result := check.Check(pr)
			if tc.err {
				assert.Greater(t, len(result), 0, i)
			} else {
				assert.Equal(t, 0, len(result), i)
			}
		}

		check := NewChecker()
		check.checkSwitchStmt(nil)
	})

	t.Run("x16", func(t *testing.T) {
		tests := []struct {
			data string
			err  bool
		}{
			{
				data: `package main
type Shape sum {
  Circle(radius int)
}

func describe(s Shape) int {
	switch s {
			case Circle(r):
				return r
	}
}
`,
			},
			{
				err: true,
				data: `package main
type Shape sum {
	Circle(radius int)
	Rect(w int, h int)
}

func describe(s Shape) int {
	switch s {
			case Circle(r):
				return r
			case Rect(r):
				return r
	}
}
`,
			},
			{
				err: true,
				data: `package main
type Shape sum {
	Circle(radius int)
	Rect(w int, h int)
}

func describe(s Shape) int {
	switch s {
			case Circle(r):
				return r
	}
}
`,
			},
			{
				err: true,
				data: `package main
type Shape sum {
	Circle(radius int)
	Rect(w int, h int)
}

func describe(s Shape) int {
	switch s {
			case Circle(r):
				return r
			default:
				return int(1)
	}
}
`,
			},
			{
				err: true,
				data: `package main
type Shape sum {
	Circle(radius int)
	Rect(w int, h int)
}

func describe(s Shape) int {
	switch s {
			case Circle(r):
				return r
			case 1 == 2:
				return int(1)
	}
}
`,
			},
			{
				err: true,
				data: `package main
type Shape sum {
	Circle(radius int)
	Rect(w int, h int)
}

func describe(s Shape) int {
	switch s {
			case Circle(r):
				return r
			case RRect(w, h):
				return w*h
	}
}
`,
			},
			{
				err: true,
				data: `package main
type Shape sum {
	Circle(radius int)
	Rect(w int, h int)
}

func describe(s Shape) int {
	switch s {
			case Circle(r):
				return r
			case Circle(x):
				return x
	}
}
`,
			},
			{
				err: true,
				data: `package main
type Shape sum {
	Circle(radius int)
	Rect(w int, h int)
}

func describe(s Shape) int {
	switch s {
			case Circle(r):
				return r
			case Rect(w):
				return w
	}
}
`,
			},
			{
				data: `package main
type Shape sum {
	Circle(radius int)
	Rect(w int, h int)
}

func describe(s Shape) int {
	switch s {
			case Circle(r):
				return r
			case Rect(w, h):
				return w*h
	}
}
`,
			},
			{
				err: true,
				data: `package main
type Shape sum {
	Circle(radius int)
	Rect(w int, h int)
}

func describe(s Shape) int {
	switch s {
			case Circle(r):
				return r
			case Rect(h, h):
				return h*h
	}
}
`,
			},
			{
				err: true,
				data: `package main
type Shape sum {
	Circle(radius int)
	Rect(w int, h int)
}

func describe(s Shape) int {
	switch s {
			case Circle(r):
				return r
			case Rect(w, h):
				fallthrough
				return h
	}
}
`,
			},
			{
				err: true,
				data: `package main
type Shape sum {
	Circle(radius int)
	Rect(w int, h int)
}

func describe(s Shape) int {
  var h int = 0
	switch s {
			case Circle(r):
				return r
			case Rect(w, h):
				return h
	}
}
`,
			},
			{
				err: true,
				data: `package main
type Shape sum {
	Circle(radius int)
	Rect(w int, h int)
}

func describe(s Shape) int {
	switch s {
			case Circle(r):
				return r
			case Rect(r, h, h):
				return h*r
	}
}
`,
			},
			{
				err: true,
				data: `package main
type Shape sum {
	Circle(radius int)
	Rect(w int, h int)
}

func describe(s Shape) int {
	switch s {
			case Circle(0):
				return 0
			case Rect(w, h):
				return h*r
	}
}
`,
			},
			{
				err: true,
				data: `package main
type UserID int
type Shape sum {
	Circle(radius int)
	Rect(w int, h int)
}

func describe(s Shape) int {
	switch s {
			case UserID(0):
				return 0
			case Rect(w, h):
				return h*r
	}
}
`,
			},
			{
				err: true,
				data: `package main
type Shape sum {
	Circle(radius int)
	Rect(w int, h int)
}

func describe(s Shape) int {
	switch s {
			case a.Circle(0):
				return 0
			case Rect(w, h):
				return h*r
	}
}
`,
			},
		}

		for i, tc := range tests {
			lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
			require.NoError(t, err)
			parser := parser.New(lex.FetchTokensFromString(tc.data))
			pr := parser.ParseFile()
			require.Equal(t, 0, len(parser.Errors))
			check := NewChecker()

			result := check.Check(pr)
			if tc.err {
				assert.Greater(t, len(result), 0, i)
			} else {
				assert.Equal(t, 0, len(result), i)
			}
		}

		check := NewChecker()
		check.checkSwitchStmt(nil)
		body := []*Symbol{
			{Name: "test", Kind: SymVar, Type: TInt},
			{Name: "test", Kind: SymVar, Type: TInt},
		}
		check.checkSwitchSumBody(nil, body)
	})

	t.Run("x17", func(t *testing.T) {
		tests := []struct {
			data string
			err  bool
		}{
			{
				data: `package main
type Color enum {
  Red;Blue;Green;Yellow
}

func action(l Color) string {
	switch l {
		case Red:
			return "red"
		case Blue:
			return "blue"
		case Green:
			return "green"
		case Yellow:
			return "yellow"
	}
}
`,
			},
			{
				err: true,
				data: `package main
type Color enum {
  Red;Blue;Green;Yellow
}

func action(l Color) string {
	switch l {
		case Red:
			return "red"
		case Blue:
			return "blue"
		case Green:
			return "green"
		case Yellow:
			return "yellow"
		default:
			return "yellow"
	}
}
`,
			},
			{
				err: true,
				data: `package main
type Color enum {
  Red;Blue;Green;Yellow
}

func action(l Color) string {
	switch l {
		case int(0):
			return "red"
		case Blue:
			return "blue"
		case Green:
			return "green"
		case Yellow:
			return "yellow"
	}
}
`,
			},
			{
				err: true,
				data: `package main
type Color enum {
  Red;Blue;Green;Yellow
}

func action(l Color) string {
	switch l {
		case Reddy:
			return "red"
		case Blue:
			return "blue"
		case Green:
			return "green"
		case Yellow:
			return "yellow"
	}
}
`,
			},
			{
				err: true,
				data: `package main
type Color enum {
  Red;Blue;Green;Yellow
}

func action(l Color) string {
	switch l {
		case Red:
			return "red"
		case Red:
			return "red"
		case Blue:
			return "blue"
		case Green:
			return "green"
		case Yellow:
			return "yellow"
	}
}
`,
			},
			{
				err: true,
				data: `package main
type Color enum {
  Red;Blue;Green;Yellow
}

func action(l Color) string {
	switch l {
		case Red:
			return "red"
		case Blue:
			return "blue"
		case Green:
			return "green"
	}
}
`,
			},
			{
				err: true,
				data: `package main
type Color enum {
  Red;Blue;Green;Yellow
}

func action(l Color) string {
	switch l {
		case Red:
			return "red"
		case Blue:
			return "blue"
		case Green:
			fallthrough
		case Yellow:
			return "yellow"
	}
}
`,
			},
			{
				err: true,
				data: `package main
type Color enum {
  Red;Blue;Green;Yellow
}

func action(l Color) string {
	switch l {
		case Red, Red:
			return "red"
		case Blue:
			return "blue"
		case Green:
			return "green"
		case Yellow:
			return "yellow"
	}
}
`,
			},
		}

		for i, tc := range tests {
			lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
			require.NoError(t, err)
			parser := parser.New(lex.FetchTokensFromString(tc.data))
			pr := parser.ParseFile()
			require.Equal(t, 0, len(parser.Errors))
			check := NewChecker()

			result := check.Check(pr)
			if tc.err {
				assert.Greater(t, len(result), 0, i)
			} else {
				assert.Equal(t, 0, len(result), i)
			}
		}
	})

	t.Run("x18", func(t *testing.T) {
		tests := []struct {
			data string
			err  bool
		}{
			{
				err: true,
				data: `package main
func x(a int) int {
	switch a {
		case 2:
			x := int(0)
		case 2:
			x := int(0)
  }
}
`,
			},
			{
				err: true,
				data: `package main
func x(a float) int {
	switch a {
		case 2.0:
			x := int(0)
		case 2.0:
			x := int(0)
  }
}
`,
			},
			{
				err: true,
				data: `package main
func x(a string) int {
	switch a {
		case "x":
			x := int(0)
		case "x":
			x := int(0)
  }
}
`,
			},
			{
				err: true,
				data: `package main
func x(a bool) int {
	switch a {
		case true:
			x := int(0)
		case true:
			x := int(0)
  }
}
`,
			},
			{
				data: `package main
func xa() bool { return true}
func xb() bool { return false}
func x(a bool) int {
	switch a {
		case xa():
			x1 := int(0)
		case xb():
			x2 := int(0)
  }
}
`,
			},
			{
				data: `package main
func xa() bool { return true}
func xb() bool { return false}
func x(a bool) int {
	switch {
		case 1 == 2:
			x1 := int(0)
		case 1 == 2:
			x2 := int(0)
  }
}
`,
			},
		}

		for i, tc := range tests {
			lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
			require.NoError(t, err)
			parser := parser.New(lex.FetchTokensFromString(tc.data))
			pr := parser.ParseFile()
			for _, v := range parser.Errors {
				fmt.Println(v.Error())
			}
			require.Equal(t, 0, len(parser.Errors))
			check := NewChecker()

			result := check.Check(pr)
			for _, v := range result {
				fmt.Println("BBBB", v.Err.Error())
			}
			if tc.err {
				assert.Greater(t, len(result), 0, i)
			} else {
				assert.Equal(t, 0, len(result), i)
			}
		}
	})

	t.Run("x19", func(t *testing.T) {
		tests := []struct {
			data string
			err  bool
		}{
			{
				err: true,
				data: `package main
type A interface {
	X() int
	Y(a int)
	Z(b int) int
}

type UserID int
type XA struct{x int}

XA implements xxx
XB implements UserID
`,
			},
			{
				data: `package main
type A interface {
	X() int
	Y(a int)
	Z(b int) int
}

type XA struct{x int}
XA implements A

func X() int {
	return int(0)
}
func Y(a int) {}
func Z(a int) int {
	return a
}
			`,
			},
			{
				err: true,
				data: `package main
type A struct {}

type XA struct{x int}
XA implements A

func X() int {
	return int(0)
}
func Y(a int) {}
`,
			},
			// 			{
			// 				err: true,
			// 				data: `package main
			// type A interface {
			// 	X() int
			// 	Y(a int)
			// 	Z(b int) int
			// }

			// type XA struct{x int}
			// XA implements A

			// func X() int {
			// 	return int(0)
			// }
			// func Y(a int) {}
			// `,
			// 			},
			{
				err: true,
				data: `package main
type A interface {
	X() int
	Y(a int)
	Z(b int) int
}

type XA struct{x int}
XA implements A

func X() int {
	return int(0)
}
func Y(a int) {}
func Z(a int) (int,int) {
	return a
}
`,
			},
		}

		for i, tc := range tests {
			lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
			require.NoError(t, err)
			parser := parser.New(lex.FetchTokensFromString(tc.data))
			pr := parser.ParseFile()
			for _, v := range parser.Errors {
				fmt.Println(v.Error())
			}
			require.Equal(t, 0, len(parser.Errors))
			check := NewChecker()

			result := check.Check(pr)
			for _, v := range result {
				fmt.Println("BBBB", v.Err.Error())
			}
			if tc.err {
				assert.Greater(t, len(result), 0, i)
			} else {
				assert.Equal(t, 0, len(result), i)
			}
		}

		check := NewChecker()
		check.scope = NewScope(check.pkgScope)
		check.pkgScope.Declare(&Symbol{
			Name: "XA",
			Kind: SymType,
			Type: TInt,
		})
		check.checkImplementsDecl(&ast.ImplementsDecl{
			TypeName:  token.Token{Value: "XA"},
			Interface: &ast.BadType{},
		})
		assert.Greater(t, len(check.errors), 0)
	})
}
