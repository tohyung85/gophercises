package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tohyung85/gophercises/exercise-17-secrets-cli/secrets/secret"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Key",
	Long:  "Gets key in encrypted file",
	Args:  cobra.MinimumNArgs(1),
	Run:   getSecret,
}

func init() {
	getCmd.Flags().StringP("key", "k", "", "your encoding key")
}

func getSecret(cmd *cobra.Command, args []string) {
	inputKey, _ := cmd.Flags().GetString("key")
	filePath := "/Users/joshuatan/Go/src/github.com/tohyung85/gophercises/exercise-17-secrets-cli/secrets/secretfile.txt"
	v, err := secret.FileVault(inputKey, filePath)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	val, err := v.Get(args[0])
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Printf("Key %s has value %s\n", inputKey, val)
}
