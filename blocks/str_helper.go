package blocks

func (h HelperDropdown) String() string {
	return sprintf("%v.%v", h.Key, h.Option)
}
