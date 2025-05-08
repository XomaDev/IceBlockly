package blocks

type LogicBoolean struct {
	RawBlock
	Value bool
}

type LogicNot struct {
	RawBlock
	Value Block
}

type LogicExpr struct {
	RawBlock
	Operator string
	Operands []Block
}

func (l LogicExpr) Continuous() bool {
	return false
}
