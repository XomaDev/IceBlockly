package blocks

type Component struct {
	RawBlock
	Name string
}

type PropertySet struct {
	RawBlock
	Component string
	Property  string
	Value     Block
}

type PropertyGet struct {
	RawBlock
	Component string
	Property  string
}

type GenericPropertySet struct {
	RawBlock
	Component Block
	Property  string
	Value     Block
}

type GenericPropertyGet struct {
	RawBlock
	Component Block
	Property  string
}

type Event struct {
	RawBlock
	IsGeneric  bool
	Component  string
	Event      string
	Parameters []string
	Body       []Block
}

func (c Event) Order() int {
	if c.IsGeneric {
		return 3
	}
	return 4
}

type MethodCall struct {
	RawBlock
	Component string
	Method    string
	Args      []Block
}

type GenericMethodCall struct {
	RawBlock
	Component Block
	Method    string
	Args      []Block
}

type AllComponent struct {
	RawBlock
	Type string
}
