package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tohyung85/gophercises/exercise-17-secrets-cli/secrets/secret"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "encrypt file with key",
	Long:  "Encrypt file",
	Args:  cobra.MinimumNArgs(1),
	Run:   encryptFile,
}

func init() {
}

func encryptFile(cmd *cobra.Command, args []string) {
	filePath := "/Users/joshuatan/Go/src/github.com/tohyung85/gophercises/exercise-17-secrets-cli/secrets/secretfile.txt"
	inputKey := args[0]
	v, err := secret.FileVault(inputKey, filePath)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	err = v.EncryptFile()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Printf("Successfully encrypted file at %s\n", filePath)
}
