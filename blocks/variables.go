package blocks

type GlobalVar struct {
	RawBlock
	Name  string
	Value Block
}

type VarGet struct {
	RawBlock
	Global bool
	Name   string
}

type VarSet struct {
	RawBlock
	Global bool
	Name   string
	Value  Block
}

type VarBody struct {
	RawBlock
	VarNames  []string
	VarValues []Block
	Body      []Block
}

type VarResult struct {
	RawBlock
	VarNames  []string
	VarValues []Block
	Result    Block
}
