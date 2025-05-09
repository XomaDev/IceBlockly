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

type DictMethod struct {
	RawBlock
	Dict   Block
	Method string
	Args   []Block
}

type DictProperty struct {
	RawBlock
	Dict     Block
	Property string
}

type DictWalkAll struct {
	RawBlock
}
