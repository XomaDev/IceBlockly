package blocks

type MakeList struct {
	RawBlock
	Elements []Block
}

type ListMethod struct {
	RawBlock
	List      Block
	Operation string
	Args      []Block
}

type ListProperty struct {
	RawBlock
	List     Block
	Property string
}

type ListGet struct {
	RawBlock
	List  Block
	Index Block
}

type ListSet struct {
	RawBlock
	List    Block
	Index   Block
	Element Block
}

type ListFunction struct {
	RawBlock
	Function string
	Args     []Block
}

type ListMap struct {
	RawBlock
	List   Block
	AsName string
	To     Block
}

type ListFilter struct {
	RawBlock
	List   Block
	AsName string
	Test   Block
}

type ListReduce struct {
	RawBlock
	List         Block
	ItemName     string
	AnsSoFarName string
	InitExpr     Block
	ApplyExpr    Block
}

type ListSort struct {
	RawBlock
	List           Block
	FirstItemName  string
	SecondItemName string
	TestExpr       Block
}

type ListSortKey struct {
	RawBlock
	List      Block
	KeyName   string
	ApplyExpr Block
}

type ListTransMin struct {
	RawBlock
	List           Block
	FirstItemName  string
	SecondItemName string
	TestExpr       Block
}

type ListTransMax struct {
	RawBlock
	List           Block
	FirstItemName  string
	SecondItemName string
	TestExpr       Block
}
