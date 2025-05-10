package blocks

func (t TextString) String() string {
	return "\"" + t.Text + "\""
}

func (t TextExpr) String() string {
	return JoinBlocks(t.Operands, " "+t.Operation+" ")
}

func (t TextProperty) String() string {
	if t.Text.Continuous() {
		return sprintf("%v.%v", t.Text, t.Property)
	}
	return sprintf("(%v).%v", t.Text, t.Property)
}

func (t TextMethod) String() string {
	pFormat := "%v.%v(%v)"
	if !t.Text.Continuous() {
		pFormat = "(%v).%v.(%v)"
	}
	return sprintf(pFormat, t.Text, t.Method, JoinBlocks(t.Args, ", "))
}

func (t TextSegment) String() string {
	pFormat := "%v[%v:%v]"
	if !t.Text.Continuous() {
		pFormat = "(%v)[%v:%v]"
	}
	return sprintf(pFormat, t.Text, t.Start, t.Length)
}

func (t TextObfuscate) String() string {
	return sprintf("obfuscate(\"%v\")", t.Text)
}
