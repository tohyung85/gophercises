package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tohyung85/gophercises/exercise-17-secrets-cli/secrets/secret"
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "decrypt file with key",
	Long:  "Decrypt file",
	Args:  cobra.MinimumNArgs(1),
	Run:   decryptFile,
}

func init() {
}

func decryptFile(cmd *cobra.Command, args []string) {
	filePath := "/Users/joshuatan/Go/src/github.com/tohyung85/gophercises/exercise-17-secrets-cli/secrets/secretfile.txt"
	inputKey := args[0]
	v, err := secret.FileVault(inputKey, filePath)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	err = v.DecryptFile()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Printf("Successfully decrypted file at %s\n", filePath)
}
