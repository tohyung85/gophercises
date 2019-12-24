package main

import (
	"fmt"
	"github.com/tohyung85/gophercises/exercise-9-cards-deck/deck"
)

func main() {
	myDeck := deck.NewDeck(deck.NumberDecks(1), deck.NumberJokers(1), deck.ShuffledDeck(true))
	fmt.Printf("There are %d cards in this deck\n", myDeck.CountDeck())
	myDeck.ListCards()
}
