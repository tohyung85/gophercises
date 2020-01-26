package main

import (
	"github.com/tohyung85/gophercises/exercise-17-secrets-cli/secret"
)

func main() {
	filePath := "/Users/joshuatan/Go/src/github.com/tohyung85/gophercises/exercise-17-secrets-cli/secretfile.txt"
	// decryptFilePath := "/Users/joshuatan/Go/src/github.com/tohyung85/gophercises/exercise-17-secrets-cli/secretfile_encrypt.txt"
	// secret.FileVault("password", filePath)
	secret.DecryptFile("password", filePath)
}
