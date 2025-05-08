package main

import (
	"IceBlockly/blocks"
	"encoding/xml"
	"fmt"
	"strings"
)

func main() {
	input := `
<xml xmlns="http://www.w3.org/1999/xhtml"><block xmlns="https://developers.google.com/blockly/xml" type="lists_slice"><value name="LIST"><block type="lists_create_with"><mutation xmlns="http://www.w3.org/1999/xhtml" items="3"></mutation><value name="ADD0"><block type="math_number"><field name="NUM">1</field></block></value><value name="ADD1"><block type="math_number"><field name="NUM">2</field></block></value><value name="ADD2"><block type="math_number"><field name="NUM">3</field></block></value></block></value><value name="INDEX1"><block type="math_number"><field name="NUM">2</field></block></value><value name="INDEX2"><block type="math_number"><field name="NUM">4</field></block></value></block></xml>
`

	decoder := xml.NewDecoder(strings.NewReader(input))
	decoder.Strict = false
	decoder.DefaultSpace = ""

	var root blocks.XmlRoot
	if err := decoder.Decode(&root); err != nil {
		panic(err)
	}

	parsedBlocks := allBlocks(root.Blocks)

	for i := range parsedBlocks {
		fmt.Println(parsedBlocks[i])
	}
}
