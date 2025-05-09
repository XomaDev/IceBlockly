package blocks

import "fmt"

func (l LogicBoolean) String() string {
	return fmt.Sprintf("%t", l.Value)
}

func (l LogicNot) String() string {
	pFormat := "!%v"
	if !l.Continuous() {
		pFormat = "!(%v)"
	}
	return fmt.Sprintf(pFormat, l.Value)
}

func (l LogicExpr) String() string {
	return JoinBlocks(l.Operands, " "+l.Operator+" ")
}
