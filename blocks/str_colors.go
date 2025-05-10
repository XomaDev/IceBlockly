package blocks

func (c Color) String() string {
	return "Color." + c.Name
}

func (c MakeColor) String() string {
	return sprintf("MakeColor(%v)", c.List)
}

func (c SplitColor) String() string {
	pFormat := "%v.split"
	if !c.Color.Continuous() {
		pFormat = "(%v).split"
	}
	return sprintf(pFormat, c.Color)
}
