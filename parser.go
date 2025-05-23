package main

import (
	"IceBlockly/blocks"
	"encoding/xml"
	"sort"
	"strconv"
	"strings"
)

func ParseBlockly(xmlContent string) [][]blocks.Block {
	return sortAndGroup(allBlocks(parseXml(xmlContent)))
}

func sortAndGroup(pBlocks []blocks.Block) [][]blocks.Block {
	if len(pBlocks) == 0 {
		return [][]blocks.Block{}
	}
	sort.Slice(pBlocks, func(i, j int) bool {
		return pBlocks[i].Order() < pBlocks[j].Order()
	})
	var grouped [][]blocks.Block
	currGroup := []blocks.Block{pBlocks[0]}
	currOrder := pBlocks[0].Order()

	for i := 1; i < len(pBlocks); i++ {
		aBlock := pBlocks[i]
		if aBlock.Order() == currOrder {
			currGroup = append(currGroup, aBlock)
		} else {
			grouped = append(grouped, currGroup)
			currGroup = []blocks.Block{aBlock}
			currOrder = aBlock.Order()
		}
	}
	if len(currGroup) > 0 {
		grouped = append(grouped, currGroup)
	}
	return grouped
}

func parseXml(xmlContent string) []blocks.RawBlock {
	decoder := xml.NewDecoder(strings.NewReader(xmlContent))
	decoder.Strict = false
	decoder.DefaultSpace = ""

	var root blocks.XmlRoot
	if err := decoder.Decode(&root); err != nil {
		panic(err)
	}
	return root.Blocks
}

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
	case "controls_if":
		return blocks.CtrlIf{RawBlock: block, Conditions: fromValues(block.Values), Bodies: fromStatements(block.Statements)}
	case "controls_forRange":
		return ctrlForRange(block)
	case "controls_forEach":
		return blocks.CtrlForEach{
			RawBlock: block,
			VarName:  block.SingleField(),
			List:     parseBlock(block.SingleValue()),
			Body:     recursiveParse(*block.SingleStatement().Block)}
	case "controls_for_each_dict":
		return ctrlForEachDict(block)
	case "controls_while":
		return blocks.CtrlWhile{
			RawBlock:  block,
			Condition: parseBlock(block.SingleValue()),
			Body:      recursiveParse(*block.SingleStatement().Block)}
	case "controls_choose":
		return ctrlChoose(block)
	case "controls_do_then_return":
		return blocks.CtrlDo{
			RawBlock: block,
			Body:     recursiveParse(*block.SingleStatement().Block),
			Result:   parseBlock(block.SingleValue()),
		}
	case "controls_eval_but_ignore":
		return blocks.CtrlMethod{RawBlock: block, Operation: "println", Args: []blocks.Block{parseBlock(block.SingleValue())}}
	case "controls_openAnotherScreen":
		return blocks.CtrlMethod{RawBlock: block, Operation: "openScreen", Args: []blocks.Block{parseBlock(block.SingleValue())}}
	case "controls_openAnotherScreenWithStartValue":
		return ctrlOpenScreenValue(block)
	case "controls_getStartValue":
		return blocks.CtrlMethod{RawBlock: block, Operation: "getStartValue"}
	case "controls_closeScreen":
		return blocks.CtrlMethod{RawBlock: block, Operation: "closeScreen"}
	case "controls_closeScreenWithValue":
		return blocks.CtrlMethod{RawBlock: block, Operation: "closeScreenWithValue", Args: []blocks.Block{parseBlock(block.SingleValue())}}
	case "controls_closeApplication":
		return blocks.CtrlMethod{RawBlock: block, Operation: "closeApp"}
	case "controls_getPlainStartText":
		return blocks.CtrlMethod{RawBlock: block, Operation: "getPlainStartValue"}
	case "controls_closeScreenWithPlainText":
		return blocks.CtrlMethod{
			RawBlock:  block,
			Operation: "closeScreenWithPlainText",
			Args:      []blocks.Block{parseBlock(block.SingleValue())}}
	case "controls_break":
		return blocks.CtrlBreak{RawBlock: block}

	case "logic_boolean", "logic_true", "logic_false":
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
		return blocks.MathExpr{Operator: "+", Operands: fromMinVals(block.Values, 2)}
	case "math_subtract":
		return blocks.MathExpr{Operator: "-", Operands: fromMinVals(block.Values, 2)}
	case "math_multiply":
		return blocks.MathExpr{Operator: "*", Operands: fromMinVals(block.Values, 2)}
	case "math_division":
		return blocks.MathExpr{Operator: "/", Operands: fromMinVals(block.Values, 2)}
	case "math_power":
		return blocks.MathExpr{Operator: "^", Operands: fromMinVals(block.Values, 2)}
	case "math_bitwise":
		return mathBitwise(block, fromMinVals(block.Values, 2))
	case "math_random_int":
		return mathRandom(block)
	case "math_random_float":
		return blocks.MathRandomFloat{}
	case "math_random_set_seed":
		return blocks.MathRandomSetSeed{RawBlock: block, Seed: parseBlock(block.SingleValue())}
	case "math_number_radix":
		return mathRadix(block)
	case "math_on_list", "math_trig", "math_sin", "math_cos", "math_tan":
		return blocks.MathFunc{RawBlock: block, Operation: strings.ToLower(block.SingleField()), Operands: fromMinVals(block.Values, 1)}
	case "math_on_list2":
		return mathOnList2(block)
	case "math_mode_of_list":
		return blocks.MathFunc{RawBlock: block, Operation: "modeOfList", Operands: fromMinVals(block.Values, 1)}
	case "math_atan2":
		return blocks.MathFunc{RawBlock: block, Operation: "atan2", Operands: fromMinVals(block.Values, 1)}
	case "math_format_as_decimal":
		return blocks.MathFunc{RawBlock: block, Operation: "formatDecimals", Operands: fromMinVals(block.Values, 1)}
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
		return blocks.TextExpr{RawBlock: block, Operation: "+", Operands: fromMinVals(block.Values, 1)}
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
		return blocks.TextProperty{RawBlock: block, Property: "isString", Text: parseBlock(block.SingleValue())}
	case "text_replace_mappings":
		return textReplaceMap(block)

	case "lists_create_with":
		return blocks.MakeList{RawBlock: block, Elements: fromMinVals(block.Values, 1)}
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
		return blocks.ListProperty{RawBlock: block, List: parseBlock(block.SingleValue()), Property: "newList"}
	case "lists_reverse":
		return blocks.ListProperty{RawBlock: block, List: parseBlock(block.SingleValue()), Property: "reverse"}
	case "lists_to_csv_row":
		return blocks.ListProperty{RawBlock: block, List: parseBlock(block.SingleValue()), Property: "toCsvRow"}
	case "lists_to_csv_table":
		return blocks.ListProperty{RawBlock: block, List: parseBlock(block.SingleValue()), Property: "toCsvTable"}
	case "lists_sort":
		return blocks.ListProperty{RawBlock: block, List: parseBlock(block.SingleValue()), Property: "sort"}
	case "lists_is_list":
		return blocks.ListProperty{RawBlock: block, List: parseBlock(block.SingleValue()), Property: "isList"}
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
		return blocks.MakeDict{RawBlock: block, Pairs: fromMinVals(block.Values, 1)}
	case "dictionaries_lookup":
		return dictLookup(block)
	case "dictionaries_set_pair":
		return dictSetPair(block)
	case "dictionaries_delete_pair":
		return dictRemove(block)
	case "dictionaries_recursive_lookup":
		return dictLookupPath(block)
	case "dictionaries_recursive_set":
		return dictSetPath(block)
	case "dictionaries_getters":
		return dictGetters(block)
	case "dictionaries_is_key_in":
		return dictHasKey(block)
	case "dictionaries_length":
		return blocks.DictProperty{RawBlock: block, Dict: parseBlock(block.SingleValue()), Property: "numKeys"}
	case "dictionaries_alist_to_dict":
		return blocks.DictProperty{RawBlock: block, Dict: parseBlock(block.SingleValue()), Property: "toDict"}
	case "dictionaries_dict_to_alist":
		return blocks.DictProperty{RawBlock: block, Dict: parseBlock(block.SingleValue()), Property: "toList"}
	case "dictionaries_copy":
		return blocks.DictProperty{RawBlock: block, Dict: parseBlock(block.SingleValue()), Property: "newDict"}
	case "dictionaries_combine_dicts":
		return dictCombine(block)
	case "dictionaries_walk_tree":
		return dictWalkTree(block)
	case "dictionaries_walk_all":
		return blocks.DictWalkAll{RawBlock: block}
	case "dictionaries_is_dict":
		return blocks.DictProperty{RawBlock: block, Dict: parseBlock(block.SingleValue()), Property: "isDict"}

	case "color_black":
		return blocks.Color{RawBlock: block, Name: "Black"}
	case "color_white":
		return blocks.Color{RawBlock: block, Name: "White"}
	case "color_red":
		return blocks.Color{RawBlock: block, Name: "Red"}
	case "color_pink":
		return blocks.Color{RawBlock: block, Name: "Pink"}
	case "color_orange":
		return blocks.Color{RawBlock: block, Name: "Orange"}
	case "color_yellow":
		return blocks.Color{RawBlock: block, Name: "Yellow"}
	case "color_green":
		return blocks.Color{RawBlock: block, Name: "Green"}
	case "color_cyan":
		return blocks.Color{RawBlock: block, Name: "Cyan"}
	case "color_blue":
		return blocks.Color{RawBlock: block, Name: "Blue"}
	case "color_magenta":
		return blocks.Color{RawBlock: block, Name: "Magenta"}
	case "color_light_gray":
		return blocks.Color{RawBlock: block, Name: "LightGray"}
	case "color_dark_gray":
		return blocks.Color{RawBlock: block, Name: "DarkGray"}
	case "color_make_color":
		return blocks.MakeColor{RawBlock: block, List: parseBlock(block.SingleValue())}
	case "color_split_color":
		return blocks.SplitColor{RawBlock: block, Color: parseBlock(block.SingleValue())}

	case "global_declaration":
		return blocks.GlobalVar{RawBlock: block, Name: block.SingleField(), Value: parseBlock(block.SingleValue())}
	case "lexical_variable_get":
		return variableGet(block)
	case "lexical_variable_set":
		return variableSet(block)
	case "local_declaration_statement", "local_declaration_expression":
		return variableSmts(block)

	case "procedures_defnoreturn":
		return voidProcedure(block)
	case "procedures_defreturn":
		return returnProcedure(block)
	case "procedures_callnoreturn", "procedures_callreturn":
		return procedureCall(block)

	case "helpers_assets":
		return blocks.TextString{RawBlock: block, Text: block.SingleField()}
	case "helpers_dropdown":
		return blocks.HelperDropdown{RawBlock: block, Key: block.Mutation.Key, Option: block.SingleField()}

	case "component_component_block":
		return blocks.Component{RawBlock: block, Name: block.SingleField()}
	case "component_set_get":
		return propSetGet(block)
	case "component_event":
		return event(block)
	case "component_method":
		return method(block)
	case "component_all_component_block":
		return blocks.AllComponent{RawBlock: block, Type: block.Mutation.ComponentType}
	default:
		panic("Unsupported block type: " + block.Type)
	}
}

func method(block blocks.RawBlock) blocks.Block {
	if block.Mutation.IsGeneric {
		pVals := makeValueMap(block.Values)
		var callArgs []blocks.Block

		for i := 0; ; i++ {
			aArg := pVals["ARG"+strconv.Itoa(i)]
			if aArg == nil {
				break
			}
			callArgs = append(callArgs, aArg)
		}

		return blocks.GenericMethodCall{
			RawBlock:  block,
			Component: pVals["COMPONENT"],
			Method:    block.Mutation.MethodName,
			Args:      callArgs,
		}
	}
	return blocks.MethodCall{
		RawBlock:  block,
		Component: block.Mutation.InstanceName,
		Method:    block.Mutation.MethodName,
		Args:      fromValues(block.Values),
	}
}

func event(block blocks.RawBlock) blocks.Block {
	var component string
	if block.Mutation.IsGeneric {
		component = block.Mutation.ComponentType
	} else {
		component = block.Mutation.InstanceName
	}
	return blocks.Event{
		IsGeneric:  block.Mutation.IsGeneric,
		RawBlock:   block,
		Component:  component,
		Event:      block.Mutation.EventName,
		Parameters: nil, // TODO ( fix 'em later )
		Body:       recursiveParse(*block.SingleStatement().Block),
	}
}

func propSetGet(block blocks.RawBlock) blocks.Block {
	pFields := makeFieldMap(block.Fields)

	property := pFields["PROP"]
	isSet := block.Mutation.SetOrGet == "set"

	if block.Mutation.IsGeneric {
		pVals := makeValueMap(block.Values)
		if isSet {
			return blocks.GenericPropertySet{
				RawBlock:  block,
				Component: pVals["COMPONENT"],
				Property:  property,
				Value:     pVals["VALUE"],
			}
		}
		return blocks.GenericPropertyGet{
			RawBlock:  block,
			Component: pVals["COMPONENT"],
			Property:  property,
		}
	}
	if isSet {
		return blocks.PropertySet{
			RawBlock:  block,
			Component: pFields["COMPONENT_SELECTOR"],
			Property:  property,
			Value:     parseBlock(block.SingleValue()),
		}
	}
	return blocks.PropertyGet{
		RawBlock:  block,
		Component: pFields["COMPONENT_SELECTOR"],
		Property:  property,
	}
}

func ctrlOpenScreenValue(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	return blocks.CtrlMethod{
		RawBlock:  block,
		Operation: "openScreenWithValue",
		Args:      []blocks.Block{pVals["SCREENNAME"], pVals["STARTVALUE"]},
	}
}

func ctrlChoose(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	return blocks.CtrlChoose{
		RawBlock:  block,
		Condition: pVals["TEST"],
		Then:      pVals["THENRETURN"],
		Else:      pVals["ELSERETURN"],
	}
}

func ctrlForEachDict(block blocks.RawBlock) blocks.Block {
	pFields := makeFieldMap(block.Fields)
	return blocks.CtrlForEachDict{
		RawBlock:  block,
		KeyName:   pFields["KEY"],
		ValueName: pFields["VALUE"],
		Dict:      parseBlock(block.SingleValue()),
		Body:      recursiveParse(*block.SingleStatement().Block),
	}
}

func ctrlForRange(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	return blocks.CtrlForRange{
		RawBlock: block,
		VarName:  block.SingleField(),
		Start:    pVals["START"],
		End:      pVals["END"],
		Step:     pVals["STEP"],
		Body:     recursiveParse(*block.SingleStatement().Block),
	}
}

func procedureCall(block blocks.RawBlock) blocks.Block {
	mutArgsNames := block.Mutation.Args
	paramNames := make([]string, len(mutArgsNames))
	for i := range mutArgsNames {
		paramNames[i] = mutArgsNames[i].Name
	}
	procedureName := block.SingleField()
	args := fromValues(block.Values)
	return blocks.ProcedureCall{
		RawBlock:   block,
		Name:       procedureName,
		Parameters: paramNames,
		Args:       args,
	}
}

func returnProcedure(block blocks.RawBlock) blocks.Block {
	procedureName := makeFieldMap(block.Fields)["NAME"]
	mutArgs := block.Mutation.Args
	paramNames := make([]string, len(mutArgs))
	for i := range mutArgs {
		paramNames[i] = mutArgs[i].Name
	}
	return blocks.ReturnProcedure{
		RawBlock:   block,
		Name:       procedureName,
		Parameters: paramNames,
		Return:     parseBlock(block.SingleValue()),
	}
}

func voidProcedure(block blocks.RawBlock) blocks.Block {
	procedureName := makeFieldMap(block.Fields)["NAME"]
	mutArgs := block.Mutation.Args
	paramNames := make([]string, len(mutArgs))
	for i := range mutArgs {
		paramNames[i] = mutArgs[i].Name
	}
	return blocks.VoidProcedure{
		RawBlock:   block,
		Name:       procedureName,
		Parameters: paramNames,
		Body:       recursiveParse(*block.SingleStatement().Block),
	}
}

func variableSmts(block blocks.RawBlock) blocks.Block {
	numOfVars := len(block.Mutation.LocalNames)
	fieldMap := makeFieldMap(block.Fields)
	valueMap := makeValueMap(block.Values)

	varNames := make([]string, numOfVars)
	varValues := make([]blocks.Block, numOfVars)

	for i := 0; i < numOfVars; i++ {
		varNames[i] = fieldMap["VAR"+strconv.Itoa(i)]
		varValues[i] = valueMap["DECL"+strconv.Itoa(i)]
	}
	if block.GetType() == "local_declaration_statement" {
		return blocks.VarBody{
			RawBlock:  block,
			VarNames:  varNames,
			VarValues: varValues,
			Body:      recursiveParse(*block.SingleStatement().Block),
		}
	}
	return blocks.VarResult{
		RawBlock:  block,
		VarNames:  varNames,
		VarValues: varValues,
		Result:    valueMap["RETURN"],
	}
}

func variableSet(block blocks.RawBlock) blocks.Block {
	varName := block.SingleField()
	isGlobal := strings.HasPrefix(varName, "global ")
	if isGlobal {
		varName = varName[len("global "):]
	}
	return blocks.VarSet{RawBlock: block, Global: isGlobal, Name: varName, Value: parseBlock(block.SingleValue())}
}

func variableGet(block blocks.RawBlock) blocks.Block {
	varName := block.Fields[0].Name
	if varName == "VAR" {
		varName = block.SingleField()
	}
	isGlobal := strings.HasPrefix(varName, "global ")
	if isGlobal {
		varName = varName[len("global "):]
	}
	return blocks.VarGet{RawBlock: block, Global: isGlobal, Name: varName}
}

func dictWalkTree(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	return blocks.DictMethod{
		RawBlock: block,
		Dict:     pVals["DICT"],
		Method:   "walkTree",
		Args:     []blocks.Block{pVals["PATH"]},
	}
}

func dictCombine(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	return blocks.DictMethod{
		RawBlock: block,
		Dict:     pVals["DICT2"],
		Method:   "mergeInto",
		Args:     []blocks.Block{pVals["DICT1"]},
	}
}

func dictHasKey(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	return blocks.DictMethod{
		RawBlock: block,
		Dict:     pVals["DICT"],
		Method:   "hasKey",
		Args:     []blocks.Block{pVals["KEY"]},
	}
}

func dictGetters(block blocks.RawBlock) blocks.Block {
	var pOperation string
	switch block.SingleField() {
	case "KEYS":
		pOperation = "keys"
	case "VALUES":
		pOperation = "values"
	default:
		panic("Unknown DictGetters operation: " + block.SingleField())
	}
	return blocks.DictProperty{
		RawBlock: block,
		Dict:     parseBlock(block.SingleValue()),
		Property: pOperation,
	}
}

func dictSetPath(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	return blocks.DictMethod{
		RawBlock: block,
		Dict:     pVals["DICT"],
		Method:   "setAtPath",
		Args:     []blocks.Block{pVals["KEYS"], pVals["VALUE"]},
	}
}

func dictLookupPath(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	return blocks.DictMethod{
		RawBlock: block,
		Dict:     pVals["DICT"],
		Method:   "getAtPath",
		Args:     []blocks.Block{pVals["KEYS"], pVals["NOTFOUND"]},
	}
}

func dictRemove(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	return blocks.DictMethod{
		RawBlock: block,
		Dict:     pVals["DICT"],
		Method:   "remove",
		Args:     []blocks.Block{pVals["KEY"]},
	}
}

func dictSetPair(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	return blocks.DictMethod{
		RawBlock: block,
		Dict:     pVals["DICT"],
		Method:   "set",
		Args:     []blocks.Block{pVals["KEY"], pVals["VALUE"]},
	}
}

func dictLookup(block blocks.RawBlock) blocks.Block {
	pVals := makeValueMap(block.Values)
	return blocks.DictMethod{
		RawBlock: block,
		Dict:     pVals["DICT"],
		Method:   "get",
		Args:     []blocks.Block{pVals["KEY"], pVals["NOTFOUND"]},
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
	return blocks.TextExpr{RawBlock: block, Operation: pOperation, Operands: fromMinVals(block.Values, 2)}
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
		Operands: fromMinVals(block.Values, 2),
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
	return blocks.MathFunc{RawBlock: block, Operation: pOperation, Operands: fromMinVals(block.Values, 1)}
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
	return blocks.MathFunc{RawBlock: block, Operation: pOperation, Operands: fromMinVals(block.Values, 1)}
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
	return blocks.MathFunc{RawBlock: block, Operation: pOperation, Operands: fromMinVals(block.Values, 1)}
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
	return blocks.MathFunc{RawBlock: block, Operation: pOperation, Operands: fromMinVals(block.Values, 2)}
}

func mathSingle(block blocks.RawBlock) blocks.Block {
	pOperation := strings.ToLower(block.SingleField())
	switch pOperation {
	case "ln":
		pOperation = "log"
	case "ceiling":
		pOperation = "ceil"
	}
	return blocks.MathFunc{RawBlock: block, Operation: pOperation, Operands: fromMinVals(block.Values, 1)}
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
	return blocks.MathFunc{RawBlock: block, Operation: pOperation, Operands: fromMinVals(block.Values, 1)}
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
	return blocks.MathExpr{RawBlock: block, Operator: pOperation, Operands: fromMinVals(block.Values, 2)}
}

func recursiveParse(currBlock blocks.RawBlock) []blocks.Block {
	var pParsed []blocks.Block
	for {
		pParsed = append(pParsed, parseBlock(currBlock))
		if currBlock.Next == nil {
			break
		}
		currBlock = *currBlock.Next.Block
	}
	return pParsed
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

func fromStatements(allSmts []blocks.Statement) [][]blocks.Block {
	arrBlocks := make([][]blocks.Block, len(allSmts))
	for i := range allSmts {
		arrBlocks[i] = recursiveParse(*allSmts[i].Block)
	}
	return arrBlocks
}

func fromValues(allValues []blocks.Value) []blocks.Block {
	arrBlocks := make([]blocks.Block, len(allValues))
	for i := range allValues {
		arrBlocks[i] = parseBlock(allValues[i].Block)
	}
	return arrBlocks
}

func fromMinVals(allValues []blocks.Value, minCount int) []blocks.Block {
	arrBlocks := make([]blocks.Block, max(minCount, len(allValues)))
	for i := range allValues {
		arrBlocks[i] = parseBlock(allValues[i].Block)
	}
	return arrBlocks
}
