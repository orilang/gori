package ir

type Program struct {
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
	Name string
}

type Const struct {
	Result string
	Type   string
	Value  string
}

type Call struct {
	Result string
	Name   string
	Args   []string
}

type Branch struct {
	Condition string
	True      string
	False     string
	Index     int
}

type Label struct {
	Name  string
	Index int
}

type Jump struct {
	Name  string
	Index int
}

type Assigment struct {
	Result string
	Value  string
}

type Unary struct {
	Result   string
	Operator string
	Value    string
}
