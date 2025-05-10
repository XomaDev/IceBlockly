package blocks

func (l MakeList) String() string {
	return "list(" + JoinBlocks(l.Elements, ", ") + ")"
}

func (l ListMethod) String() string {
	pFormat := "%v.%v(%v)"
	if !l.List.Continuous() {
		pFormat = "(%v).%v(%v)"
	}
	return sprintf(pFormat, l.List, l.Operation, JoinBlocks(l.Args, ", "))
}

func (l ListFunction) String() string {
	return sprintf("%v(%v)", l.Function, JoinBlocks(l.Args, ", "))
}

func (l ListProperty) String() string {
	pFormat := "%v.%v"
	if !l.List.Continuous() {
		pFormat = "(%v).%v"
	}
	return sprintf(pFormat, l.List, l.Property)
}

func (l ListGet) String() string {
	pFormat := "%v[%v]"
	if !l.List.Continuous() {
		pFormat = "(%v).[%v]"
	}
	return sprintf(pFormat, l.List, l.Index)
}

func (l ListSet) String() string {
	pFormat := "%v[%v] = %v"
	if !l.List.Continuous() {
		pFormat = "(%v)[%v] = %v"
	}
	return sprintf(pFormat, l.List, l.Index, l.Element)
}

func (l ListMap) String() string {
	pFormat := "%v"
	if !l.List.Continuous() {
		pFormat = "(%v)"
	}
	pFormat += ".map { %v ->\n%v}"
	return sprintf(pFormat, l.List, l.AsName, Pad(l.To))
}

func (l ListFilter) String() string {
	pFormat := "%v"
	if !l.List.Continuous() {
		pFormat = "(%v)"
	}
	pFormat += ".filter { %v ->\n%v}"
	return sprintf(pFormat, l.List, l.AsName, Pad(l.Test))
}

func (l ListReduce) String() string {
	pFormat := "%v"
	if !l.List.Continuous() {
		pFormat = "(%v)"
	}
	pFormat += ".reduce(%v) { %v, %v ->\n%v}"
	return sprintf(pFormat, l.List, l.InitExpr, l.ItemName, l.AnsSoFarName, Pad(l.ApplyExpr))
}

func (l ListSort) String() string {
	pFormat := "%v"
	if !l.List.Continuous() {
		pFormat = "(%v)"
	}
	pFormat += ".sort { %v, %v ->\n%v}"
	return sprintf(pFormat, l.List, l.FirstItemName, l.SecondItemName, Pad(l.TestExpr))
}

func (l ListSortKey) String() string {
	pFormat := "%v"
	if !l.List.Continuous() {
		pFormat = "(%v)"
	}
	pFormat += ".sortByKey { %v ->\n%v}"
	return sprintf(pFormat, l.List, l.KeyName, Pad(l.ApplyExpr))
}

func (l ListTransMin) String() string {
	pFormat := "%v"
	if !l.List.Continuous() {
		pFormat = "(%v)"
	}
	pFormat += ".min { %v, %v ->\n%v}"
	return sprintf(pFormat, l.List, l.FirstItemName, l.SecondItemName, Pad(l.TestExpr))
}

func (l ListTransMax) String() string {
	pFormat := "%v"
	if !l.List.Continuous() {
		pFormat = "(%v)"
	}
	pFormat += ".max { %v, %v ->\n%v}"
	return sprintf(pFormat, l.List, l.FirstItemName, l.SecondItemName, Pad(l.TestExpr))
}
