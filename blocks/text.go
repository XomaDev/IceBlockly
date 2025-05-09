package blocks

type TextString struct {
	RawBlock
	Text string
}

type TextExpr struct {
	RawBlock
	Operation string
	Operands  []Block
}

func (t TextExpr) Continuous() bool {
	return false
}

type TextProperty struct {
	RawBlock
	Property string
	Text     Block
}

type TextMethod struct {
	RawBlock
	Method string
	Text   Block
	Args   []Block
}

type TextSegment struct {
	RawBlock
	Text   Block
	Start  Block
	Length Block
}

type TextObfuscate struct {
	RawBlock
	Text string
}
