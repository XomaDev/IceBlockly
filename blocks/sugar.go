package blocks

import "strings"

func Pad(block Block) string {
	return " " + strings.Replace(block.String(), "\n", "\n  ", -1) + "\n"
}

func PadLine(line string) string {
	return " " + strings.Replace(line, "\n", "\n  ", -1) + "\n"
}

func PadBody(blocks []Block) string {
	var builder strings.Builder
	for _, block := range blocks {
		builder.WriteString(Pad(block))
	}
	return builder.String()
}

func JoinBlocks(blocks []Block, delimiter string) string {
	opStrings := make([]string, len(blocks))
	for i, op := range blocks {
		opStrings[i] = op.String()
	}
	return strings.Join(opStrings, delimiter)
}
