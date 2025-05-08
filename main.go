package main

import (
	"IceBlockly/blocks"
	"encoding/xml"
	"fmt"
	"strings"
)

func main() {
	input := `
<xml
    xmlns="http://www.w3.org/1999/xhtml">
    <block
        xmlns="https://developers.google.com/blockly/xml" type="math_random_int">
        <value name="FROM">
            <block type="math_number">
                <field name="NUM">1</field>
            </block>
        </value>
        <value name="TO">
            <block type="math_number">
                <field name="NUM">100</field>
            </block>
        </value>
    </block>
</xml>

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
