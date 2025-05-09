package blocks

import (
	"fmt"
	"strings"
)

func (c CtrlIf) String() string {
	var builder strings.Builder

	numConditions := len(c.Conditions)
	currI := 0

	builder.WriteString("if ")
	for {
		builder.WriteString(fmt.Sprintf("(%v) {\n", c.Conditions[currI]))
		builder.WriteString(PadBody(c.Bodies[currI]))
		builder.WriteString("}")
		currI += 1
		if currI < numConditions {
			builder.WriteString(" elif ")
		} else {
			break
		}
	}
	if currI < len(c.Bodies) {
		// it's an else clause!
		builder.WriteString(" else {\n")
		for _, expr := range c.Bodies[currI] {
			builder.WriteString(Pad(expr))
		}
		builder.WriteString("}")
	}
	return builder.String()
}

func (c CtrlForRange) String() string {
	return fmt.Sprintf("for (%v: %v to %v by %v) {\n%v}", c.VarName, c.Start, c.End, c.Step, PadBody(c.Body))
}

func (c CtrlForEach) String() string {
	return fmt.Sprintf("each (%v -> %v) {\n%v}", c.VarName, c.List, PadBody(c.Body))
}

func (c CtrlForEachDict) String() string {
	return fmt.Sprintf("each (%v::%v -> %v) {\n%v}", c.KeyName, c.ValueName, c.Dict, PadBody(c.Body))
}

func (c CtrlWhile) String() string {
	return fmt.Sprintf("while (%v) {\n%v}", c.Condition, PadBody(c.Body))
}

func (c CtrlChoose) String() string {
	return fmt.Sprintf("%v\n%v%v", c.Condition, PadLine("? "+c.Then.String()), PadLine(": "+c.Else.String()))
}

func (c CtrlDo) String() string {
	return fmt.Sprintf("do {\n%v} -> %v", PadBody(c.Body), c.Result)
}

func (c CtrlMethod) String() string {
	return fmt.Sprintf("%v(%v)", c.Operation, JoinBlocks(c.Args, ", "))
}

func (c CtrlBreak) String() string {
	return "break"
}
