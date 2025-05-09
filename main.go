package main

import (
	"fmt"
)

func main() {
	input := `
<xml xmlns="http://www.w3.org/1999/xhtml"><block xmlns="https://developers.google.com/blockly/xml" type="helpers_dropdown"><mutation xmlns="http://www.w3.org/1999/xhtml" key="ScreenAnimation"></mutation><field name="OPTION">Default</field></block></xml>
`

	parsedBlocks := ParseBlockly(input)

	for i := range parsedBlocks {
		fmt.Println(parsedBlocks[i])
	}
}
