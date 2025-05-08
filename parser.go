package main

import (
	"IceBlockly/blocks"
	"strconv"
	"strings"
)

func allBlocks(allBlocks []blocks.RawBlock) []blocks.Block {
	var parsedBlocks []blocks.Block
	for i := range allBlocks {
		parsedBlocks = append(parsedBlocks, parseBlock(allBlocks[i]))
	}
	return parsedBlocks
}

// TODO (future):
//  * when two Ints are TextJoined, we'd produce `123 + 123`
//    which is an arithmetic operation.
//    ==Solution==:
//    > While Blockly -> Mist
//      If any of the operands resolves to a String, keep as it is
//      If none of the operands resolve to a String, add a toString() wrapper call
//    > While Mist -> Blockly
//      We'd decide if the '+' means TextJoin or Arithmetic Add by looking
//      into the possible resolved types. If there's a string resolved, it's a TextJoin

func parseBlock(block blocks.RawBlock) blocks.Block {
	switch block.Type {
	case "logic_boolean":
		return blocks.LogicBoolean{RawBlock: block, Value: block.SingleField() == "TRUE"}
	case "logic_negate":
		return blocks.LogicNot{RawBlock: block, Value: parseBlock(block.SingleValue())}
	case "logic_compare", "logic_operation":
		return logicExpr(block)

	case "math_number":
		return blocks.MathNumber{RawBlock: block, Value: block.SingleField()}
	case "math_compare":
		return mathCompare(block)
	case "math_add":
		return blocks.MathExpr{Operator: "+", Operands: fromValues(block.Values)}
	case "math_subtract":
		return blocks.MathExpr{Operator: "-", Operands: fromValues(block.Values)}
	case "math_multiply":
		return blocks.MathExpr{Operator: "*", Operands: fromValues(block.Values)}
	case "math_division":
		return blocks.MathExpr{Operator: "/", Operands: fromValues(block.Values)}
	case "math_power":
		return blocks.MathExpr{Operator: "^", Operands: fromValues(block.Values)}
	case "math_bitwise":
		return mathBitwise(block, fromValues(block.Values))
	case "math_random_int":
		return mathRandom(block)
	case "math_random_float":
		return blocks.MathRandomFloat{}
	case "math_random_set_seed":
		return blocks.MathRandomSetSeed{RawBlock: block, Seed: parseBlock(block.SingleValue())}
	case "math_number_radix":
		return mathRadix(block)
	case "math_on_list", "math_trig", "math_sin", "math_cos", "math_tan":
		return blocks.MathFunc{RawBlock: block, Operation: strings.ToLower(block.SingleField()), Operands: fromValues(block.Values)}
	case "math_on_list2":
		return mathOnList2(block)
	case "math_mode_of_list":
		return blocks.MathFunc{RawBlock: block, Operation: "modeOfList", Operands: fromValues(block.Values)}
	case "math_atan2":
		return blocks.MathFunc{RawBlock: block, Operation: "atan2", Operands: fromValues(block.Values)}
	case "math_format_as_decimal":
		return blocks.MathFunc{RawBlock: block, Operation: "formatDecimals", Operands: fromValues(block.Values)}
	case "math_single":
		return mathSingle(block)
	case "math_divide":
		return mathDivide(block)
	case "math_convert_angles":
		return mathAngles(block)
	case "math_is_a_number":
		return mathIsNumber(block)
	case "math_convert_number":
		return mathConvertNumber(block)

	case "text":
		return blocks.TextString{RawBlock: block, Text: block.SingleField()}
	case "text_join":
		return blocks.TextExpr{RawBlock: block, Operation: "+", Operands: fromValues(block.Values)}
	case "text_length":
		return blocks.TextProperty{RawBlock: block, Property: "len", Text: parseBlock(block.SingleValue())}
	case "text_isEmpty":
		return blocks.TextProperty{RawBlock: block, Property: "isEmpty", Text: parseBlock(block.SingleValue())}
	case "text_trim":
		return blocks.TextProperty{RawBlock: block, Property: "trim", Text: parseBlock(block.SingleValue())}
	case "text_reverse":
		return blocks.TextProperty{RawBlock: block, Property: "reverse", Text: parseBlock(block.SingleValue())}
	case "text_split_at_spaces":
		return blocks.TextProperty{RawBlock: block, Property: "splitAtSpaces", Text: parseBlock(block.SingleValue())}
	case "text_compare":
		return textCompare(block)
	case "text_changeCase":
		return textChangeCase(block)
	case "text_starts_at":
		return textStartsWith(block)
	case "text_contains":
		return textContains(block)
	case "text_split":
		return textSplit(block)
	case "text_segment":
		return textSegment(block)
	case "text_replace_all":
		return textReplace(block)
	case "obfuscated_text":
		return blocks.TextObfuscate{RawBlock: block, Text: block.SingleField()}
	case "text_is_string":
		return blocks.TextIsString{RawBlock: block, Value: parseBlock(block.SingleValue())}
	case "text_replace_mappings":
		return textReplaceMap(block)

	case "lists_create_with":
		return blocks.MakeList{RawBlock: block, Elements: fromValues(block.Values)}
	case "lists_add_items":
		return listAddItem(block)
	case "lists_is_in":
		return listContainsItem(block)
	case "lists_length":
		return blocks.ListProperty{RawBlock: block, List: parseBlock(block.SingleValue()), Property: "size"}
	case "lists_is_empty":
		return blocks.ListProperty{RawBlock: block, List: parseBlock(block.SingleValue()), Property: "isEmpty"}
	case "lists_pick_random_item":
		return blocks.ListProperty{RawBlock: block, List: parseBlock(block.SingleValue()), Property: "random"}
	case "lists_position_in":
		return listIndexOf(block)
	case "lists_select_item":
		return listSelectItem(block)
	case "lists_insert_item":
		return listInsertItem(block)
	case "lists_replace_item":
		return listReplaceItem(block)
	case "lists_remove_item":
		return listRemoveItem(block)
	case "lists_append_list":
		return listAppendList(block)
	case "lists_copy":
		return blocks.ListProperty{RawBlock: block, List: parseBlock(block.SingleValue()), Property: "copy"}
	case "lists_reverse":
		return blocks.ListProperty{RawBlock: block, List: parseBlock(block.SingleValue()), Property: "reverse"}
	case "lists_to_csv_row":
		return blocks.ListProperty{RawBlock: block, List: parseBlock(block.SingleValue()), Property: "toCsvRow"}
	case "lists_to_csv_table":
		return blocks.ListProperty{RawBlock: block, List: parseBlock(block.SingleValue()), Property: "toCsvTable"}
	case "lists_sort":
		return blocks.ListProperty{RawBlock: block, List: parseBlock(block.SingleValue()), Property: "sort"}
	case "lists_is_list":
		return blocks.ListFunction{RawBlock: block, Function: "isList", Args: []blocks.Block{parseBlock(block.SingleValue())}}
	case "lists_from_csv_row":
		return blocks.ListFunction{RawBlock: block, Function: "listFromCsvRow", Args: []blocks.Block{parseBlock(block.SingleValue())}}
	case "lists_from_csv_table":
		return blocks.ListFunction{RawBlock: block, Function: "listFromCsvTable", Args: []blocks.Block{parseBlock(block.SingleValue())}}
	case "lists_but_first":
		return blocks.ListFunction{RawBlock: block, Function: "exceptFirst", Args: []blocks.Block{parseBlock(block.SingleValue())}}
	case "lists_but_last":
		return blocks.ListFunction{RawBlock: block, Function: "exceptLast", Args: []blocks.Block{parseBlock(block.SingleValue())}}
	case "lists_lookup_in_pairs":
		return listLookupPairs(block)
	case "lists_join_with_separator":
		return listJoin(block)
	case "lists_slice":
		return listSlice(block)
	case "lists_map":
		return listMap(block)
	case "lists_filter":
		return listFilter(block)
	case "lists_reduce":
		return listReduce(block)
	case "lists_sort_comparator":
		return listSortComparator(block)
	case "lists_sort_key":
		return listSortKeyComparator(block)
	case "lists_minimum_value":
		return listTransMin(block)
	case "lists_maximum_value":
		return listTransMax(block)

	case "pair":
		return dictPair(block)
	case "dictionaries_create_with":
		return blocks.MakeDict{RawBlock: block, Pairs: fromValues(block.Values)}
	default:
		panic("Unsupported block type: " + block.Type)
	}
}

func listSlice(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	return blocks.ListMethod{
		RawBlock:  block,
		List:      pVals["LIST"],
		Operation: "slice",
		Args:      []blocks.Block{pVals["INDEX1"], pVals["INDEX2"]},
	}
}

func listTransMax(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	pFields := makeFieldMap(block.Fields)
	return blocks.ListTransMax{
		RawBlock:       block,
		List:           pVals["LIST"],
		FirstItemName:  pFields["VAR1"],
		SecondItemName: pFields["VAR2"],
		TestExpr:       pVals["COMPARE"],
	}
}

func listTransMin(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	pFields := makeFieldMap(block.Fields)
	return blocks.ListTransMin{
		RawBlock:       block,
		List:           pVals["LIST"],
		FirstItemName:  pFields["VAR1"],
		SecondItemName: pFields["VAR2"],
		TestExpr:       pVals["COMPARE"],
	}
}

func listSortKeyComparator(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	return blocks.ListSortKey{
		RawBlock:  block,
		List:      pVals["LIST"],
		KeyName:   block.SingleField(),
		ApplyExpr: pVals["KEY"],
	}
}

func listSortComparator(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	pFields := makeFieldMap(block.Fields)
	return blocks.ListSort{
		RawBlock:       block,
		List:           pVals["LIST"],
		FirstItemName:  pFields["VAR1"],
		SecondItemName: pFields["VAR2"],
		TestExpr:       pVals["COMPARE"],
	}
}

func listReduce(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	pFields := makeFieldMap(block.Fields)
	return blocks.ListReduce{
		RawBlock:     block,
		List:         pVals["LIST"],
		ItemName:     pFields["VAR1"],
		AnsSoFarName: pFields["VAR2"],
		InitExpr:     pVals["INITANSWER"],
		ApplyExpr:    pVals["COMBINE"],
	}
}

func listFilter(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	return blocks.ListFilter{
		RawBlock: block,
		List:     pVals["LIST"],
		AsName:   block.SingleField(),
		Test:     pVals["TEST"],
	}
}

func listMap(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	return blocks.ListMap{
		RawBlock: block,
		List:     pVals["LIST"],
		AsName:   block.SingleField(),
		To:       pVals["TO"],
	}
}

func listJoin(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	return blocks.ListMethod{
		RawBlock:  block,
		List:      pVals["LIST"],
		Operation: "join",
		Args:      []blocks.Block{pVals["SEPARATOR"]},
	}
}

func listLookupPairs(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	return blocks.ListMethod{
		RawBlock:  block,
		List:      pVals["LIST"],
		Operation: "lookup",
		Args:      []blocks.Block{pVals["KEY"], pVals["NOTFOUND"]},
	}
}

func listAppendList(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	return blocks.ListMethod{
		RawBlock:  block,
		List:      pVals["LIST0"],
		Operation: "append",
		Args:      []blocks.Block{pVals["LIST1"]},
	}
}

func listRemoveItem(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	return blocks.ListMethod{
		RawBlock:  block,
		List:      pVals["LIST"],
		Operation: "remove",
		Args:      []blocks.Block{pVals["INDEX"]},
	}
}

func listReplaceItem(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	return blocks.ListSet{
		RawBlock: block,
		List:     pVals["LIST"],
		Index:    pVals["NUM"],
		Element:  pVals["ITEM"],
	}
}

func listInsertItem(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	return blocks.ListMethod{
		RawBlock:  block,
		List:      pVals["LIST"],
		Operation: "insert",
		Args:      []blocks.Block{pVals["INDEX"], pVals["ITEM"]},
	}
}

func listSelectItem(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	return blocks.ListGet{
		RawBlock: block,
		List:     pVals["LIST"],
		Index:    pVals["NUM"],
	}
}

func listIndexOf(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	return blocks.ListMethod{
		RawBlock:  block,
		List:      pVals["LIST"],
		Operation: "indexOf",
		Args:      []blocks.Block{pVals["ITEM"]},
	}
}

func listContainsItem(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	return blocks.ListMethod{
		RawBlock:  block,
		List:      pVals["LIST"],
		Operation: "contains",
		Args:      []blocks.Block{pVals["ITEM"]},
	}
}

func listAddItem(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	numElements := block.Mutation.ItemCount
	arrElements := make([]blocks.Block, numElements)

	for i := 0; i < numElements; i++ {
		arrElements[i] = pVals["ITEM"+strconv.Itoa(i)]
	}
	return blocks.ListMethod{
		RawBlock:  block,
		List:      pVals["LIST"],
		Operation: "add",
		Args:      arrElements,
	}
}

func dictPair(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	return blocks.Pair{
		RawBlock: block,
		Key:      pVals["KEY"],
		Value:    pVals["VALUE"],
	}
}

func textReplaceMap(block blocks.RawBlock) blocks.Block {
	var pOperation string
	switch block.SingleField() {
	case "LONGEST_STRING_FIRST":
		pOperation = "replaceMapLongestFirst"
	case "DICTIONARY_ORDER":
		pOperation = "replaceMap"
	}
	pVals := makeValueMap(block.Values)
	return blocks.TextMethod{
		RawBlock: block,
		Method:   pOperation,
		Text:     pVals["TEXT"],
		Args:     []blocks.Block{pVals["MAPPINGS"]},
	}
}

func textReplace(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	return blocks.TextMethod{
		RawBlock: block,
		Method:   "replace",
		Text:     pVals["TEXT"],
		Args:     []blocks.Block{pVals["SEGMENT"], pVals["REPLACEMENT"]},
	}
}

func textSegment(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	return blocks.TextSegment{
		RawBlock: block,
		Text:     pVals["TEXT"],
		Start:    pVals["START"],
		Length:   pVals["LENGTH"],
	}
}

func textSplit(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	var pOperation string
	switch block.SingleField() {
	case "SPLIT":
		pOperation = "split"
	case "SPLITATFIRST":
		pOperation = "splitFirst"
	case "SPLITATANY":
		pOperation = "splitAny"
	case "SPLITATFIRSTOFANY":
		pOperation = "splitFirstOfAny"
	default:
		panic("Unsupported TextSplit block operation: " + block.SingleField())
	}
	return blocks.TextMethod{
		RawBlock: block,
		Method:   pOperation,
		Text:     pVals["TEXT"],
		Args:     []blocks.Block{pVals["AT"]},
	}
}

func textContains(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	var pOperation string
	switch block.SingleField() {
	case "CONTAINS":
		pOperation = "contains"
	case "CONTAINS_ANY":
		pOperation = "containsAny"
	case "CONTAINS_ALL":
		pOperation = "containsAll"
	default:
		panic("Unsupported TextContains operation: " + block.SingleField())
	}
	return blocks.TextMethod{
		RawBlock: block,
		Method:   pOperation,
		Text:     pVals["TEXT"],
		Args:     []blocks.Block{pVals["PIECE"]},
	}
}

func textStartsWith(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	return blocks.TextMethod{
		RawBlock: block,
		Method:   "startsWith",
		Text:     pVals["TEXT"],
		Args:     []blocks.Block{pVals["PIECE"]},
	}
}

func textChangeCase(block blocks.RawBlock) blocks.Block {
	var pOperation string
	switch block.SingleField() {
	case "UPCASE":
		pOperation = "upper"
	case "DOWNCASE":
		pOperation = "lower"
	default:
		panic("Unsupported TextChangeCase operation type: " + block.SingleField())
	}
	return blocks.TextProperty{RawBlock: block, Property: pOperation, Text: parseBlock(block.SingleValue())}
}

func textCompare(block blocks.RawBlock) blocks.Block {
	var pOperation string
	switch block.SingleField() {
	case "EQ":
		pOperation = "==="
	case "NEQ":
		pOperation = "!=="
	case "LT":
		pOperation = "<<"
	case "RT":
		pOperation = ">>"
	default:
		panic("Unknown Text Compare operation: " + block.SingleField())
	}
	return blocks.TextExpr{RawBlock: block, Operation: pOperation, Operands: fromValues(block.Values)}
}

func logicExpr(block blocks.RawBlock) blocks.Block {
	var pOperation string
	switch block.SingleField() {
	case "EQ":
		pOperation = "=="
	case "NEQ":
		pOperation = "!="
	case "AND":
		pOperation = "&&"
	case "OR":
		pOperation = "||"
	default:
		panic("Unknown Logic Compare operation: " + block.SingleField())
	}
	return blocks.LogicExpr{
		RawBlock: blocks.RawBlock{},
		Operator: pOperation,
		Operands: fromValues(block.Values),
	}
}

func mathConvertNumber(block blocks.RawBlock) blocks.Block {
	var pOperation string
	switch block.SingleField() {
	case "DEC_TO_HEX":
		pOperation = "toHex"
	case "HEX_TO_DEC":
		pOperation = "fromHex"
	case "DEC_TO_BIN":
		pOperation = "toBin"
	case "BIN_TO_DEC":
		pOperation = "fromBin"
	default:
		panic("Unknown MathConvertNumber type: " + block.SingleField())
	}
	return blocks.MathFunc{RawBlock: block, Operation: pOperation, Operands: fromValues(block.Values)}
}

func mathIsNumber(block blocks.RawBlock) blocks.Block {
	var pOperation string
	switch block.SingleField() {
	case "NUMBER":
		pOperation = "isNum"
	case "BINARY":
		pOperation = "isBin"
	case "HEXADECIMAL":
		pOperation = "isHexa"
	case "BASE10":
		pOperation = "isDecimal"
	default:
		panic("Unknown MathIsNumber type: " + block.SingleField())
	}
	return blocks.MathFunc{RawBlock: block, Operation: pOperation, Operands: fromValues(block.Values)}
}

func mathAngles(block blocks.RawBlock) blocks.Block {
	var pOperation string
	switch block.SingleField() {
	case "RADIANS_TO_DEGREES":
		pOperation = "toDegrees"
	case "DEGREES_TO_RADIANS":
		pOperation = "toRadians"
	default:
		panic("Unsupported math angle type: " + block.SingleField())
	}
	return blocks.MathFunc{RawBlock: block, Operation: pOperation, Operands: fromValues(block.Values)}
}

func mathDivide(block blocks.RawBlock) blocks.Block {
	var pOperation string
	switch block.SingleField() {
	case "MODULO":
		pOperation = "mod"
	case "REMAINDER":
		pOperation = "rem"
	case "QUOTIENT":
		pOperation = "quot"
	default:
		panic("Unsupported math divide type: " + block.SingleField())
	}
	return blocks.MathFunc{RawBlock: block, Operation: pOperation, Operands: fromValues(block.Values)}
}

func mathSingle(block blocks.RawBlock) blocks.Block {
	pOperation := strings.ToLower(block.SingleField())
	switch pOperation {
	case "ln":
		pOperation = "log"
	case "ceiling":
		pOperation = "ceil"
	}
	return blocks.MathFunc{RawBlock: block, Operation: pOperation, Operands: fromValues(block.Values)}
}

func mathOnList2(block blocks.RawBlock) blocks.Block {
	var pOperation string
	switch block.SingleField() {
	case "AVG":
		pOperation = "avg"
	case "MIN":
		pOperation = "minList"
	case "MAX":
		pOperation = "maxList"
	case "GM":
		pOperation = "geoMean"
	case "SD":
		pOperation = "stdDev"
	case "SE":
		pOperation = "stdErr"
	default:
		panic("Unsupported math_on_list2 operation: " + block.SingleField())
	}
	return blocks.MathFunc{RawBlock: block, Operation: pOperation, Operands: fromValues(block.Values)}
}

func mathRadix(block blocks.RawBlock) blocks.Block {
	pFields := makeFieldMap(block.Fields)
	var radix int
	switch pFields["OP"] {
	case "DEC":
		radix = 10
	case "BIN":
		radix = 2
	case "OCT":
		radix = 8
	case "HEX":
		radix = 16
	}
	return blocks.MathRadix{RawBlock: block, Radix: radix, Number: pFields["NUM"]}
}

func mathBitwise(block blocks.RawBlock, operands []blocks.Block) blocks.Block {
	var mathOp string
	switch block.SingleField() {
	case "BITAND":
		mathOp = "&"
	case "BITOR":
		mathOp = "|"
	case "BITXOR":
		mathOp = "^^"
	default:
		panic("Unsupported bitwise operator: " + block.SingleField())
	}
	return blocks.MathExpr{Operator: mathOp, Operands: operands}
}

func mathRandom(block blocks.RawBlock) blocks.Block {
	valMap := makeValueMap(block.Values)
	return blocks.MathRandomInt{RawBlock: block, From: valMap["FROM"], To: valMap["TO"]}
}

func mathCompare(block blocks.RawBlock) blocks.Block {
	var pOperation string
	switch block.SingleField() {
	case "EQ":
		pOperation = "=="
	case "NEQ":
		pOperation = "!="
	case "LT":
		pOperation = "<"
	case "LTE":
		pOperation = "<="
	case "GT":
		pOperation = ">"
	case "GTE":
		pOperation = ">="
	default:
		panic("Unsupported MathCompare operation: " + block.SingleField())
	}
	return blocks.MathExpr{RawBlock: block, Operator: pOperation, Operands: fromValues(block.Values)}
}

func makeFieldMap(allFields []blocks.Field) map[string]string {
	valueMap := make(map[string]string, len(allFields))
	for _, fil := range allFields {
		valueMap[fil.Name] = fil.Value
	}
	return valueMap
}

func makeValueMap(allValues []blocks.Value) map[string]blocks.Block {
	valueMap := make(map[string]blocks.Block, len(allValues))
	for _, val := range allValues {
		valueMap[val.Name] = parseBlock(val.Block)
	}
	return valueMap
}

func fromValues(allValues []blocks.Value) []blocks.Block {
	arrBlocks := make([]blocks.Block, len(allValues))
	for i := range allValues {
		arrBlocks[i] = parseBlock(allValues[i].Block)
	}
	return arrBlocks
}
