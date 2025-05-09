package blocks

import (
	"fmt"
	"strings"
)

func (p VoidProcedure) String() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("fun %v(%v) {\n", p.Name, strings.Join(p.Parameters, ", ")))
	for i := range p.Body {
		builder.WriteString(Pad(p.Body[i]))
	}
	builder.WriteString("}")
	return builder.String()
}

func (p ReturnProcedure) String() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("fun %v(%v) {\n", p.Name, strings.Join(p.Parameters, ", ")))
	builder.WriteString(Pad(p.Return))
	builder.WriteString("}")
	return builder.String()
}

func (p ProcedureCall) String() string {
	prepArgs := make([]string, len(p.Args))
	for i := range p.Args {
		prepArgs[i] = fmt.Sprintf("%v: %v", p.Parameters[i], p.Args[i])
	}
	return fmt.Sprintf("%v(%v)", p.Name, strings.Join(prepArgs, ", "))
}
