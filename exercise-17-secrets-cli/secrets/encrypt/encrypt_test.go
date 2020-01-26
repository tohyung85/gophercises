package encrypt

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	plaintext := "Super Secret Text"
	encryptedData, err := Encrypt("password", []byte(plaintext))
	if err != nil {
		t.Errorf("Error: %s\n", err)
	}

	if string(encryptedData) != plaintext {
		t.Logf("Success: Converted %s to %s\n", plaintext, string(encryptedData))
	} else {
		t.Errorf("Failed text was not encrypted: %s\n", string(encryptedData))
	}

	decipheredText, err := Decrypt("password", encryptedData)
	if err != nil {
		t.Errorf("Error: %s\n", err)
	}

	if string(decipheredText) == plaintext {
		t.Logf("Success: Got back text of %s\n", plaintext)
	} else {
		t.Errorf("Failed: Incorrect decryption\n")
	}
}

func TestEncryptDecryptFail(t *testing.T) {
	plaintext := "Super Secret Text"
	encryptedData, err := Encrypt("password", []byte(plaintext))
	if err != nil {
		t.Errorf("Error: %s\n", err)
	}

	decipheredText, err := Decrypt("passwords", encryptedData)
	if err != nil {
		t.Logf("Success: Encryption in place %s\n", err)
	} else {
		t.Errorf("Failed: Deciphered Text to %s, with incorrect passphrase\n", decipheredText)
	}
}
