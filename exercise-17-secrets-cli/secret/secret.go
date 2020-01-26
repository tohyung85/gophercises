package secret

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tohyung85/gophercises/exercise-17-secrets-cli/encrypt"
)

type vault struct {
	passphrase string
	filePath   string
}

func FileVault(passphrase string, path string) (*vault, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	encryptedContents, err := encrypt.Encrypt(passphrase, contents)
	if err != nil {
		return nil, err
	}

	encryptedContentsStr := string(encryptedContents)
	err = overwriteFile(encryptedContentsStr, path)
	if err != nil {
		return nil, err
	}

	return &vault{passphrase, path}, nil
}

// func (v *vault) Get(keyName string) {

// }

func DecryptFile(passphrase string, path string) error {
	encryptedContents, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("error: %s", err)
		return err
	}

	decryptedContents, err := encrypt.Decrypt(passphrase, encryptedContents)

	err = overwriteFile(string(decryptedContents), path)
	if err != nil {
		fmt.Printf("%s", err)
		return err
	}
	return nil
}

func readFileLines(file *os.File) ([]string, error) {
	fileLines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileLines = append(fileLines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return fileLines, nil
}

func overwriteFile(lines string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	fmt.Fprintf(w, "%s", lines)
	return w.Flush()
}

func parseLine(line string) (map[string]string, error) {
	keyVal := strings.Split(line, "=")
	if len(keyVal) != 2 {
		return nil, fmt.Errorf("Unable to parse line. Invalid format")
	}
	keyValMap := make(map[string]string)
	keyValMap[keyVal[0]] = keyVal[1]
	return keyValMap, nil
}

func parseMapToStringSlice(keyVal map[string]string) ([]string, error) {
	if len(keyVal) != 2 {
		return nil, fmt.Errorf("Unable to parse Data. Invalid format")
	}

	lineArr := make([]string, 0)

	for key, val := range keyVal {
		line := fmt.Sprintf("%s=%s", key, val)
		lineArr = append(lineArr, line)
	}

	return lineArr, nil
}

func encryptLine(line string, passphrase string) ([]byte, error) {
	encryptedText, err := encrypt.Encrypt(passphrase, []byte(line))
	if err != nil {
		return nil, err
	}
	return encryptedText, nil
}

func decryptLine(line string, passphrase string) ([]byte, error) {
	decryptedText, err := encrypt.Decrypt(passphrase, []byte(line))
	if err != nil {
		return nil, err
	}
	return decryptedText, nil

}
