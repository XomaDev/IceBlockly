package blocks

func (d Pair) String() string {
	return sprintf("%v:%v", d.Key, d.Value)
}

func (d MakeDict) String() string {
	return sprintf("{%v}", JoinBlocks(d.Pairs, ", "))
}

func (d DictMethod) String() string {
	pFormat := "%v.%v(%v)"
	if !d.Dict.Continuous() {
		pFormat = "(%v).%v(%v)"
	}
	return sprintf(pFormat, d.Dict, d.Method, JoinBlocks(d.Args, ", "))
}

func (d DictProperty) String() string {
	pFormat := "%v.%v"
	if !d.Dict.Continuous() {
		pFormat = "(%v).%v"
	}
	return sprintf(pFormat, d.Dict, d.Property)
}

func (d DictWalkAll) String() string {
	return "WalkAll"
}
