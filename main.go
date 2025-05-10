package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
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

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic occurred in ParseBlockly:", r)
			debug.PrintStack()
		}
	}()
	groups := ParseBlockly(input)
	for _, group := range groups {
		for _, block := range group {
			fmt.Println(block)
			if block.Order() > 0 {
				fmt.Println()
			}
		}
		fmt.Println()
	}
}
