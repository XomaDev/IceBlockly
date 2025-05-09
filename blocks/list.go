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

func (l ListMap) Continuous() bool {
	return false
}

type ListFilter struct {
	RawBlock
	List   Block
	AsName string
	Test   Block
}

func (l ListFilter) Continuous() bool {
	return false
}

type ListReduce struct {
	RawBlock
	List         Block
	ItemName     string
	AnsSoFarName string
	InitExpr     Block
	ApplyExpr    Block
}

func (l ListReduce) Continuous() bool {
	return false
}

type ListSort struct {
	RawBlock
	List           Block
	FirstItemName  string
	SecondItemName string
	TestExpr       Block
}

func (l ListSort) Continuous() bool {
	return false
}

type ListSortKey struct {
	RawBlock
	List      Block
	KeyName   string
	ApplyExpr Block
}

func (l ListSortKey) Continuous() bool {
	return false
}

type ListTransMin struct {
	RawBlock
	List           Block
	FirstItemName  string
	SecondItemName string
	TestExpr       Block
}

func (l ListTransMin) Continuous() bool {
	return false
}

type ListTransMax struct {
	RawBlock
	List           Block
	FirstItemName  string
	SecondItemName string
	TestExpr       Block
}

func (l ListTransMax) Continuous() bool {
	return false
}
