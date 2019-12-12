package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func main() {
	fmt.Println("Running Caesar Cipher Program")
	runCaesarCipherProgram()
}

func runCaesarCipherProgram() {
	shiftPtr := flag.Int("shift", 3, "Letters to shift each character by. Defaults to 3")
	flag.Parse()
	inp := ""
	reader := bufio.NewReader(os.Stdin)

	for {
		inp, _ = reader.ReadString('\n')
		inp = strings.TrimSpace(inp)

		if inp == "quit" {
			break
		}

		result := caesarCipher(inp, int32(*shiftPtr))
		fmt.Printf("Encrypted Message: %s, Decoded Message: %s\n", inp, result)
		fmt.Println(`Type "quit" to exit`)
	}

	return
}

func caesarCipher(s string, k int32) string {
	decodedStr := ""
	for _, char := range s {
		ascii := int32(char)
		newAscii := int32(-1)
		switch {
		case !unicode.IsLetter(char):
			newAscii = ascii
		case unicode.IsUpper(char):
			newAscii = 65 + (ascii+k-65)%26
		case unicode.IsLower(char):
			newAscii = 97 + (ascii+k-97)%26
		}
		decodedStr += string(newAscii)
	}
	return decodedStr
}
