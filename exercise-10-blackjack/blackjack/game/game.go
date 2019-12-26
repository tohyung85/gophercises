package game

import (
	"fmt"
	"github.com/tohyung85/gophercises/exercise-9-cards-deck/deck"
)

type Role int

const (
	Dealer Role = iota
	Player
)

type Game struct {
	deck    *deck.Deck
	players []Participant
}

type Participant struct {
	role Role
}

func StartGame(decks int, jokers int, shuffle bool, removeSuits []string, removeRanks []string) {
	fmt.Println("Starting game with following parameters:")
	fmt.Printf("%d decks\n", decks)
	fmt.Printf("%d jokers\n", jokers)
	fmt.Printf("Deck to be shuffled? %t\n", shuffle)
	fmt.Printf("Suits %v will be removed\n", removeSuits)
	fmt.Printf("Ranks %v will be removed\n", removeRanks)
	gameDeck := initializeDeck(decks, jokers, shuffle, removeSuits, removeRanks)
	game := initializeGame(gameDeck)
	fmt.Printf("%s", game)
}

func initializeGame(gameDeck *deck.Deck) *Game {
	players := initializePlayers(1)
	return &Game{gameDeck, players}
}

func initializePlayers(numPlayers int) []Participant {
	dealer := Participant{Dealer}
	participants := []Participant{dealer}

	for i := 0; i < numPlayers; i++ {
		player := &Participant{Player}
		participants = append(participants, *player)
	}

	return participants
}

func initializeDeck(decks int, jokers int, shuffle bool, removeSuits []string, removeRanks []string) *deck.Deck {
	suitMap := deck.GetSuitMap()
	rankMap := deck.GetRankMap()

	omittedSuits := make([]deck.Suit, 0)
	for _, s := range removeSuits {
		suit, inMap := suitMap[s]
		if inMap {
			omittedSuits = append(omittedSuits, suit)
		}
	}
	omittedRanks := make([]deck.Rank, 0)
	for _, r := range removeRanks {
		rank, inMap := rankMap[r]
		if inMap {
			omittedRanks = append(omittedRanks, rank)
		}
	}
	return deck.NewDeck(deck.NumberDecks(decks), deck.NumberJokers(jokers), deck.ShuffledDeck(shuffle), deck.OmitRanks(omittedRanks...), deck.OmitSuits(omittedSuits...))
}

func (game *Game) String() string {
	return fmt.Sprintf("Game has %d players including the dealer\n", len(game.players))
}
