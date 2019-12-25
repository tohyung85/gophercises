package main

import (
	"fmt"
	"github.com/tohyung85/gophercises/exercise-9-cards-deck/deck"
)

func main() {
	removeCard := deck.Card{Rank: deck.Ace, Suit: deck.Club}
	myDeck := deck.NewDeck(
		deck.NumberDecks(2),
		deck.NumberJokers(1),
		deck.OmitRanks(deck.Two),
		deck.OmitSuits(deck.Heart, deck.Spade),
		deck.OmitCards(removeCard),
		deck.SortedBy(customSort),
	)
	fmt.Printf("There are %d cards in this deck\n", myDeck.CountDeck())
	fmt.Printf("Listing all cards in the deck:\n%s", myDeck)
}

func customSort(cards []deck.Card) func(i, j int) bool {
	return func(i, j int) bool {
		return cards[i].Suit.String() < cards[j].Suit.String()
	}
}
