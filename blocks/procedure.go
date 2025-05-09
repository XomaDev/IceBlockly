package blocks

type VoidProcedure struct {
	RawBlock
	Name       string
	Parameters []string
	Body       []Block
}

func (p VoidProcedure) Order() int {
	return 1
}

type ReturnProcedure struct {
	RawBlock
	Name       string
	Parameters []string
	Return     Block
}

func (p ReturnProcedure) Order() int {
	return 2
}

type ProcedureCall struct {
	RawBlock
	Name       string
	Parameters []string
	Args       []Block
}
