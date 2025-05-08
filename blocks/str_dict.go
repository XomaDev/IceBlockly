package blocks

import "fmt"

func (d Pair) String() string {
	return fmt.Sprintf("%v:%v", d.Key, d.Value)
}

func (d MakeDict) String() string {
	return fmt.Sprintf("{%v}", JoinBlocks(d.Pairs, ", "))
}
