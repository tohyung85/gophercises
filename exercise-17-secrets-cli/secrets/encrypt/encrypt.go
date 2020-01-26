package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"golang.org/x/crypto/scrypt"
)

const saltSize = 8

func Encrypt(passphrase string, data []byte) ([]byte, error) {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	dk, err := scrypt.Key([]byte(passphrase), salt, 1<<15, 8, 1, 32)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(dk)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	cipherText := gcm.Seal(nonce, nonce, data, nil)

	cipherTextWithSalt := append(salt, cipherText...)

	return cipherTextWithSalt, nil
}

func Decrypt(passphrase string, encryptedData []byte) ([]byte, error) {
	salt := encryptedData[:saltSize]
	dk, err := scrypt.Key([]byte(passphrase), salt, 1<<15, 8, 1, 32)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(dk)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceStart := saltSize + gcm.NonceSize()

	nonce, encryptedText := encryptedData[saltSize:nonceStart], encryptedData[nonceStart:]

	plainText, err := gcm.Open(nil, nonce, encryptedText, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}
