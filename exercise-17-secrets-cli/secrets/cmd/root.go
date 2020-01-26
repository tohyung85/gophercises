package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Secrets Manager",
	Long:  "Cli for managing encryption",
}

func init() {
}

func Execute() error {
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(decryptCmd)
	return rootCmd.Execute()
}
