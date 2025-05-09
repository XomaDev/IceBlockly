package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Enter XML:")

	// Read multiline input from terminal
	var inputLines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		inputLines = append(inputLines, line)
	}
	input := strings.Join(inputLines, "\n")

	parsedBlocks := ParseBlockly(input)

	for i := range parsedBlocks {
		fmt.Println(parsedBlocks[i])
	}
}
