package blocks

import (
	"fmt"
	"strings"
)

func (m MathNumber) String() string {
	return m.Value
}

func (m MathExpr) String() string {
	opStrings := make([]string, len(m.Operands))
	for i, op := range m.Operands {
		opStrings[i] = op.String()
	}
	return strings.Join(opStrings, " "+m.Operator+" ")
}

func (m MathRandomInt) String() string {
	return fmt.Sprintf("randint(%v, %v)", m.From, m.To)
}

func (m MathRandomFloat) String() string {
	return "randfloat()"
}
