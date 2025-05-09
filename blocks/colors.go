package blocks

type Color struct {
	RawBlock
	Name string
}

type MakeColor struct {
	RawBlock
	List Block
}

type SplitColor struct {
	RawBlock
	Color Block
}
