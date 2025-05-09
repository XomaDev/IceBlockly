package blocks

import "fmt"

func (h HelperDropdown) String() string {
	return fmt.Sprintf("%v.%v", h.Key, h.Option)
}
