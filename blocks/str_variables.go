package blocks

import (
	"fmt"
	"strings"
)

func (v GlobalVar) String() string {
	return fmt.Sprintf("glob %v = %v", v.Name, v.Value)
}

func (v VarGet) String() string {
	pFormat := "%v"
	if v.Global {
		pFormat = "glob.%v"
	}
	return fmt.Sprintf(pFormat, v.Name)
}

func (v VarSet) String() string {
	pFormat := "%v = %v"
	if v.Global {
		pFormat = "glob.%v = %v"
	}
	return fmt.Sprintf(pFormat, v.Name, v.Value)
}

func (v VarBody) String() string {
	var builder strings.Builder
	for i := range v.VarValues {
		builder.WriteString(fmt.Sprintf("val %v = %v\n", v.VarNames[i], v.VarValues[i]))
	}
	var pBody []string
	for _, expr := range v.Body {
		pBody = append(pBody, expr.String())
	}
	return builder.String() + strings.Join(pBody, "\n")
}

func (v VarResult) String() string {
	var builder strings.Builder
	for i := range v.VarValues {
		builder.WriteString(fmt.Sprintf("val %v = %v\n", v.VarNames[i], v.VarValues[i]))
	}
	return builder.String() + v.Result.String()
}
