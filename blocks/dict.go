package blocks

type Pair struct {
	RawBlock
	Key   Block
	Value Block
}

type MakeDict struct {
	RawBlock
	Pairs []Block
}
