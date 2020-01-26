package secret

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tohyung85/gophercises/exercise-17-secrets-cli/secrets/encrypt"
)

type vault struct {
	passphrase string
	filePath   string
}

func FileVault(passphrase string, path string) (*vault, error) {
	return &vault{passphrase, path}, nil
}

func (v *vault) Set(keyName string, keyVal string) error {
	keyValMap, err := getKeyValMap(v.passphrase, v.filePath)
	if err != nil {
		return err
	}

	keyValMap[keyName] = keyVal

	lineArr, err := parseMapToStringSlice(keyValMap)
	if err != nil {
		return err
	}

	fullString := strings.Join(lineArr, "\n")

	err = overwriteFile(fullString, v.filePath)
	if err != nil {
		return err
	}

	err = v.EncryptFile()
	if err != nil {
		return err
	}

	return nil
}

func (v *vault) Get(keyName string) (string, error) {
	keyValMap, err := getKeyValMap(v.passphrase, v.filePath)
	if err != nil {
		return "", err
	}

	val, inMap := keyValMap[keyName]
	if !inMap {
		return "", fmt.Errorf("Error: Key of %s does not exist!\n", keyName)
	}
	return val, nil
}

func getKeyValMap(passphrase string, path string) (map[string]string, error) {
	decryptedString, err := getDecryptedContents(passphrase, path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(decryptedString, "\n")

	keyValMap := make(map[string]string)

	for _, l := range lines {
		key, val, err := parseLine(l)
		if err != nil {
			fmt.Printf("Error Parsing line: %s\n %s", l, err)
			continue
		}
		keyValMap[key] = val
	}

	return keyValMap, nil
}

func getDecryptedContents(passphrase string, path string) (string, error) {
	encryptedContents, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("error: %s", err)
		return "", err
	}

	decryptedContents, err := encrypt.Decrypt(passphrase, encryptedContents)
	if err != nil {
		fmt.Printf("error: %s", err)
		return "", err
	}

	return string(decryptedContents), nil
}

func (v *vault) EncryptFile() error {
	contents, err := ioutil.ReadFile(v.filePath)
	if err != nil {
		return err
	}

	encryptedContents, err := encrypt.Encrypt(v.passphrase, contents)
	if err != nil {
		return err
	}

	encryptedContentsStr := string(encryptedContents)
	err = overwriteFile(encryptedContentsStr, v.filePath)
	if err != nil {
		return err
	}

	return nil
}

func (v *vault) DecryptFile() error {
	encryptedContents, err := ioutil.ReadFile(v.filePath)
	if err != nil {
		fmt.Printf("error: %s", err)
		return err
	}

	decryptedContents, err := encrypt.Decrypt(v.passphrase, encryptedContents)

	err = overwriteFile(string(decryptedContents), v.filePath)
	if err != nil {
		fmt.Printf("%s", err)
		return err
	}
	return nil
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

func parseLine(line string) (string, string, error) {
	keyVal := strings.Split(line, "=")
	if len(keyVal) != 2 {
		return "", "", fmt.Errorf("Unable to parse line. Invalid format")
	}
	return keyVal[0], keyVal[1], nil
}

func parseMapToStringSlice(keyVal map[string]string) ([]string, error) {
	lineArr := make([]string, 0)

	for key, val := range keyVal {
		line := fmt.Sprintf("%s=%s", key, val)
		lineArr = append(lineArr, line)
	}

	return lineArr, nil
}
