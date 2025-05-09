package main

import (
	"fmt"
)

func main() {
	input := `
<xml xmlns="http://www.w3.org/1999/xhtml"><block xmlns="https://developers.google.com/blockly/xml" type="component_all_component_block"><mutation xmlns="http://www.w3.org/1999/xhtml" component_type="Web"></mutation><field name="COMPONENT_TYPE_SELECTOR">Web</field></block></xml>
`

	parsedBlocks := ParseBlockly(input)

	for i := range parsedBlocks {
		fmt.Println(parsedBlocks[i])
	}
}
