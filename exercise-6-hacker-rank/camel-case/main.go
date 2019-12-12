package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

func main() {
	fmt.Println("Starting camel case program...")
	runCamelCaseProgram()

}

func runCamelCaseProgram() {
	inp := ""
	reader := bufio.NewReader(os.Stdin)

	for {
		inp, _ = reader.ReadString('\n')
		inp = strings.TrimSpace(inp)

		if inp == "quit" {
			break
		}

		result := camelcase(inp)
		fmt.Printf("There are %d words in %s\n", result, inp)
		fmt.Println(`Type "quit" to exit`)
	}

	return
}

func camelcase(s string) int32 {
	fmt.Printf("Analysing %s...\n", s)
	wordCount := 0

	if s == "" {
		return int32(wordCount)
	}

	wordCount++

	for i := 0; i < len(s)-1; i++ {
		nextChar, _ := utf8.DecodeRuneInString(s[i+1:])
		if unicode.IsUpper(nextChar) {
			wordCount++
		}
	}

	return int32(wordCount)
}
