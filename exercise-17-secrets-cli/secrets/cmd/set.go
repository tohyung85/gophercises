package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tohyung85/gophercises/exercise-17-secrets-cli/secrets/secret"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set Key",
	Long:  "Sets key in encrypted file",
	Args:  cobra.MinimumNArgs(2),
	Run:   setKey,
}

func init() {
	setCmd.Flags().StringP("key", "k", "", "your encoding key")
}

func setKey(cmd *cobra.Command, args []string) {
	encodeKey, _ := cmd.Flags().GetString("key")
	inputKey := args[0]
	inputVal := args[1]

	filePath := "/Users/joshuatan/Go/src/github.com/tohyung85/gophercises/exercise-17-secrets-cli/secrets/secretfile.txt"
	v, err := secret.FileVault(encodeKey, filePath)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	err = v.Set(inputKey, inputVal)
	if err != nil {
		fmt.Printf("Error Setting Key Val: %s\n", err)
		return
	}

	fmt.Printf("Successfully set key %s with value %s\n", inputKey, inputVal)
}
