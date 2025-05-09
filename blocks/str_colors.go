package blocks

import "fmt"

func (c Color) String() string {
	return "Color." + c.Name
}

func (c MakeColor) String() string {
	return "MakeColor(" + c.List.String() + ")"
}

func (c SplitColor) String() string {
	pFormat := "%v.split"
	if !c.Color.Continuous() {
		pFormat = "(%v).split"
	}
	return fmt.Sprintf(pFormat, c.Color)
}
