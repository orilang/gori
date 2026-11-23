package ir

type Program struct {
	Types []Instruction
	Funcs []*Func
}

type Func struct {
	Name    string
	Params  []Param
	Results []Param
	Blocks  []*Block
	Label   string
}

type Param struct {
	Name string
	Type string
}

type Block struct {
	Instructions []Instruction
}

type Instruction interface {
	instrNode()
}

type Value string

type Binary struct {
	Result string
	Op     string
	Left   string
	Right  string
}

type Return struct {
	Values []string
}

type Const struct {
	Result string
	Type   string
	Value  string
}

type Call struct {
	Result   string
	Name     string
	Args     []string
	FromFunc bool
}

type BranchSub struct {
	Name     string
	Index    int
	NoSuffix bool
}
type Branch struct {
	Condition string
	List      []BranchSub
}

type Label struct {
	Name     string
	Index    int
	NoSuffix bool
}

type Jump struct {
	Name     string
	Index    int
	NoSuffix bool
}

type Assigment struct {
	Result   string
	Operator string
	Value    string
}

type Extract struct {
	Result string
	Value  string
	Index  int
}

type Unary struct {
	Result   string
	Operator string
	Value    string
}

type IncDec struct {
	Value    string
	Operator string
}

type Index struct {
	Result string
	Name   string
	Value  string
}

type Len struct {
	Result string
	Value  string
}

type MapIter struct {
	Result string
	Value  string
}

type IterNext struct {
	Result string
	Value  string
}

type IterKey struct {
	Result string
	Value  string
}

type IterValue struct {
	Result string
	Value  string
}

type DefinedTypeDecl struct {
	Name string
	Type string
}
