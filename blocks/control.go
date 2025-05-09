package blocks

type CtrlIf struct {
	RawBlock
	Conditions []Block
	Bodies     [][]Block
}

type CtrlForRange struct {
	RawBlock
	VarName string
	Start   Block
	End     Block
	Step    Block
	Body    []Block
}

type CtrlForEach struct {
	RawBlock
	VarName string
	List    Block
	Body    []Block
}

type CtrlForEachDict struct {
	RawBlock
	KeyName   string
	ValueName string
	Dict      Block
	Body      []Block
}

type CtrlWhile struct {
	RawBlock
	Condition Block
	Body      []Block
}

type CtrlChoose struct {
	RawBlock
	Condition Block
	Then      Block
	Else      Block
}

func (c CtrlChoose) Continuous() bool {
	return false
}

type CtrlDo struct {
	RawBlock
	Body   []Block
	Result Block
}

func (c CtrlDo) Continuous() bool {
	return false
}

type CtrlMethod struct {
	RawBlock
	Operation string
	Args      []Block
}

type CtrlBreak struct {
	RawBlock
}
