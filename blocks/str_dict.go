package blocks

import "fmt"

func (d Pair) String() string {
	return fmt.Sprintf("%v:%v", d.Key, d.Value)
}

func (d MakeDict) String() string {
	return fmt.Sprintf("{%v}", JoinBlocks(d.Pairs, ", "))
}

func (d DictMethod) String() string {
	pFormat := "%v.%v(%v)"
	if !d.Dict.Continuous() {
		pFormat = "(%v).%v(%v)"
	}
	return fmt.Sprintf(pFormat, d.Dict, d.Method, JoinBlocks(d.Args, ", "))
}

func (d DictProperty) String() string {
	pFormat := "%v.%v"
	if !d.Dict.Continuous() {
		pFormat = "(%v).%v"
	}
	return fmt.Sprintf(pFormat, d.Dict, d.Property)
}

func (d DictWalkAll) String() string {
	return "WalkAll"
}
