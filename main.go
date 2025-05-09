package main

import (
	"fmt"
)

func main() {
	input := `
<xml xmlns="http://www.w3.org/1999/xhtml"><block xmlns="https://developers.google.com/blockly/xml" type="procedures_callreturn" inline="false"><mutation xmlns="http://www.w3.org/1999/xhtml" name="greet"><arg name="name"></arg></mutation><field name="PROCNAME">greet</field><value name="ARG0"><block type="text"><field name="TEXT">Eki</field></block></value></block></xml>
`

	parsedBlocks := ParseBlockly(input)

	for i := range parsedBlocks {
		fmt.Println(parsedBlocks[i])
	}
}
