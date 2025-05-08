package main

import (
	"IceBlockly/blocks"
	"encoding/xml"
	"fmt"
	"strings"
)

func main() {
	input := `
<xml xmlns="http://www.w3.org/1999/xhtml"><block xmlns="https://developers.google.com/blockly/xml" type="text_replace_mappings"><field name="OP">LONGEST_STRING_FIRST</field><value name="MAPPINGS"><block type="dictionaries_create_with"><mutation xmlns="http://www.w3.org/1999/xhtml" items="0"></mutation></block></value><value name="TEXT"><block type="text"><field name="TEXT">Hello, World!</field></block></value></block></xml>
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
