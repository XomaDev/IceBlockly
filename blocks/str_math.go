package blocks

import (
	"fmt"
)

func (m MathNumber) String() string {
	return m.Value
}

func (m MathExpr) String() string {
	return JoinBlocks(m.Operands, " "+m.Operator+" ")
}

func (m MathRandomInt) String() string {
	return fmt.Sprintf("randint(%v, %v)", m.From, m.To)
}

func (m MathRandomFloat) String() string {
	return "randfloat()"
}

func (m MathRandomSetSeed) String() string {
	return fmt.Sprintf("randsetseed(%v)", m.Seed)
}

func (m MathRadix) String() string {
	return fmt.Sprintf("radix(%v, \"%v\")", m.Radix, m.Number)
}

func (m MathFunc) String() string {
	return fmt.Sprintf("%v(%v)", m.Operation, JoinBlocks(m.Operands, ", "))
}
