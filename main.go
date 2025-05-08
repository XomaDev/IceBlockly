package main

import (
	"IceBlockly/blocks"
	"encoding/xml"
	"fmt"
	"strings"
)

func main() {
	input := `
<xml xmlns="http://www.w3.org/1999/xhtml"><block xmlns="https://developers.google.com/blockly/xml" type="logic_operation"><mutation xmlns="http://www.w3.org/1999/xhtml" items="2"></mutation><field name="OP">OR</field><value name="A"><block type="logic_boolean"><field name="BOOL">TRUE</field></block></value><value name="B"><block type="logic_boolean"><field name="BOOL">FALSE</field></block></value></block></xml>
`

	decoder := xml.NewDecoder(strings.NewReader(input))
	decoder.Strict = false
	decoder.DefaultSpace = ""

	var root blocks.XmlRoot
	if err := decoder.Decode(&root); err != nil {
		panic(err)
	}

	parsedBlocks := allBlocks(root.Blocks)

	fmt.Println(parsedBlocks)
}
