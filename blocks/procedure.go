package blocks

type VoidProcedure struct {
	RawBlock
	Name       string
	Parameters []string
	Body       []Block
}

type ReturnProcedure struct {
	RawBlock
	Name       string
	Parameters []string
	Return     Block
}

type ProcedureCall struct {
	RawBlock
	Name       string
	Parameters []string
	Args       []Block
}
