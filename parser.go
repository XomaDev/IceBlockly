package main

import (
	"IceBlockly/blocks"
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
	case "math_number":
		return blocks.MathNumber{RawBlock: block, Value: block.Field}
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
	default:
		panic("Unsupported block type: " + block.Type)
	}
}

func mathBitwise(block blocks.RawBlock, operands []blocks.Block) blocks.Block {
	var mathOp string
	switch block.Field {
	case "BITAND":
		mathOp = "&"
	case "BITOR":
		mathOp = "|"
	case "BITXOR":
		mathOp = "^^"
	default:
		panic("Unsupported bitwise operator: " + block.Field)
	}
	return blocks.MathExpr{Operator: mathOp, Operands: operands}
}

func mathRandom(block blocks.RawBlock) blocks.Block {
	valMap := makeValueMap(block.Values)
	return blocks.MathRandomInt{RawBlock: block, From: valMap["FROM"], To: valMap["TO"]}
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
