package blocks

import "fmt"

func (l LogicBoolean) String() string {
	return fmt.Sprintf("%t", l.Value)
}

func (l LogicNot) String() string {
	return "!" + l.Value.String()
}

func (l LogicExpr) String() string {
	return JoinBlocks(l.Operands, " "+l.Operator+" ")
}
