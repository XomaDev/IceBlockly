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

type MathRandomSetSeed struct {
	RawBlock
	Seed Block
}

type MathRadix struct {
	RawBlock
	Radix  int
	Number string
}

type MathFunc struct {
	RawBlock
	Operation string
	Operands  []Block
}
