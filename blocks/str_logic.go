package blocks

func (l LogicBoolean) String() string {
	return sprintf("%t", l.Value)
}

func (l LogicNot) String() string {
	pFormat := "!%v"
	if !l.Continuous() {
		pFormat = "!(%v)"
	}
	return sprintf(pFormat, l.Value)
}

func (l LogicExpr) String() string {
	return JoinBlocks(l.Operands, " "+l.Operator+" ")
}
