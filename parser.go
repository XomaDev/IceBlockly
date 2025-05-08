package main

import (
	"IceBlockly/blocks"
	"strings"
)

func allBlocks(allBlocks []blocks.RawBlock) []blocks.Block {
	var parsedBlocks []blocks.Block
	for i := range allBlocks {
		parsedBlocks = append(parsedBlocks, parseBlock(allBlocks[i]))
	}
	return parsedBlocks
}

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
	default:
		panic("Unsupported block type: " + block.Type)
	}
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
