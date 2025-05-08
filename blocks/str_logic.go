package blocks

func (l LogicBoolean) String() string {
	if l.Value {
		return "true"
	}
	return "false"
}

func (l LogicNot) String() string {
	return "!" + l.Value.String()
}

func (l LogicExpr) String() string {
	return JoinBlocks(l.Operands, " "+l.Operator+" ")
}
