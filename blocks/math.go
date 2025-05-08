package blocks

type MathNumber struct {
	RawBlock
	Value string
}

type MathExpr struct {
	RawBlock
	Operator string
	Operands []Block
}

type MathRandomInt struct {
	RawBlock
	From Block
	To   Block
}

type MathRandomFloat struct {
	RawBlock
}
