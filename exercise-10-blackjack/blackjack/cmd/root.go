package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "blackjack",
	Short: "Blackjack Game",
	Long:  "Starts a Blackjack Game",
}

func init() {
	startGameCmd.Flags().IntP("decks", "d", 1, "Number of decks to initialize with")
	startGameCmd.Flags().IntP("jokers", "j", 0, "Number of Jokers deck will be initialized with")
	startGameCmd.Flags().IntP("players", "p", 1, "Number of players in this game excluding the dealer")
	startGameCmd.Flags().BoolP("shuffle", "s", false, "Specify if deck should be shuffled")
	startGameCmd.Flags().StringSliceP("omit-suits", "o", []string{}, "List of suits (H,D,S,C) to be omitted from the deck")
	startGameCmd.Flags().StringSliceP("omit-ranks", "r", []string{}, "List of ranks (A,1,2...,J,Q,K) to be omitted from the deck")
}

func Execute() error {
	rootCmd.AddCommand(startGameCmd)
	return rootCmd.Execute()
}
