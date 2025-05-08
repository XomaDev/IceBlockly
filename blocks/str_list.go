package blocks

func (l MakeList) String() string {
	return "[" + JoinBlocks(l.Elements, ", ") + "]"
}
