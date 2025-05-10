package blocks

import (
	"fmt"
	"strings"
)

func Pad(block Block) string {
	return " " + strings.Replace(sprintf("%v", block), "\n", "\n  ", -1) + "\n"
}

func PadDirect(code string) string {
	return " " + strings.Replace(code, "\n", "\n  ", -1)
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
		opStrings[i] = sprintf("%v", op)
	}
	return strings.Join(opStrings, delimiter)
}

func sprintf(f string, args ...interface{}) string {
	safeArgs := make([]interface{}, len(args))
	for i, a := range args {
		if a == nil {
			safeArgs[i] = "(empty)"
		} else {
			safeArgs[i] = a
		}
	}
	return fmt.Sprintf(f, safeArgs...)
}
