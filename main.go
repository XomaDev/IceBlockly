package main

import (
	"IceBlockly/blocks"
	"encoding/xml"
	"fmt"
	"strings"
)

func main() {
	input := `
<xml xmlns="http://www.w3.org/1999/xhtml"><block xmlns="https://developers.google.com/blockly/xml" type="dictionaries_is_dict"><value name="THING"><block type="dictionaries_create_with"><mutation xmlns="http://www.w3.org/1999/xhtml" items="0"></mutation></block></value></block></xml>
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
