package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tohyung85/gophercises/exercise-10-blackjack/blackjack/game"
)

var startGameCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Game",
	Long:  "Starts a Blackjack Game",
	Run:   startGame,
}

func startGame(cmd *cobra.Command, args []string) {
	numDecks, _ := cmd.Flags().GetInt("decks")
	jokersToAdd, _ := cmd.Flags().GetInt("jokers")
	toShuffle, _ := cmd.Flags().GetBool("shuffle")
	suitsToOmit, _ := cmd.Flags().GetStringSlice("omit-suits")
	ranksToOmit, _ := cmd.Flags().GetStringSlice("omit-ranks")
	numPlayers, _ := cmd.Flags().GetInt("players")
	newGame := game.InitializeGame(numDecks, jokersToAdd, toShuffle, suitsToOmit, ranksToOmit, numPlayers)
	err := newGame.Start()
	if err != nil {
		fmt.Printf("%s", err)
	}
}
