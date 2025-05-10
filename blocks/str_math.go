package blocks

func (m MathNumber) String() string {
	return m.Value
}

func (m MathExpr) String() string {
	return JoinBlocks(m.Operands, " "+m.Operator+" ")
}

func (m MathRandomInt) String() string {
	return sprintf("randint(%v, %v)", m.From, m.To)
}

func (m MathRandomFloat) String() string {
	return "randfloat()"
}

func (m MathRandomSetSeed) String() string {
	return sprintf("randsetseed(%v)", m.Seed)
}

func (m MathRadix) String() string {
	return sprintf("radix(%v, \"%v\")", m.Radix, m.Number)
}

func (m MathFunc) String() string {
	return sprintf("%v(%v)", m.Operation, JoinBlocks(m.Operands, ", "))
}
