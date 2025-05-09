package main

import (
	"fmt"
)

func main() {
	input := `
<xml xmlns="http://www.w3.org/1999/xhtml"><block xmlns="https://developers.google.com/blockly/xml" type="local_declaration_expression"><mutation xmlns="http://www.w3.org/1999/xhtml"><localname name="name"></localname><localname name="x"></localname></mutation><field name="VAR0">name</field><field name="VAR1">x</field><value name="DECL0"><block type="text"><field name="TEXT">Car</field></block></value><value name="DECL1"><block type="text"><field name="TEXT">Bekku</field></block></value><value name="RETURN"><block type="text"><field name="TEXT">Hola World</field></block></value></block></xml>
`

	parsedBlocks := ParseBlockly(input)

	for i := range parsedBlocks {
		fmt.Println(parsedBlocks[i])
	}
}
