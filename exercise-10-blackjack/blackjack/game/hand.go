package game

import (
	"fmt"

	"github.com/tohyung85/gophercises/exercise-9-cards-deck/deck"
)

type Hand struct {
	cards     []deck.Card
	belongTo  *Participant
	betAmount int
	won       bool
}

func (hand *Hand) Points() int {
	aces := make([]deck.Card, 0)
	totalPoints := 0
	for _, c := range hand.cards {
		switch c.Rank {
		case deck.Ace:
			aces = append(aces, c)
		case deck.King, deck.Queen, deck.Jack:
			totalPoints += 10
		default:
			totalPoints += int(c.Rank)
		}
	}
	for range aces {
		if totalPoints+11 > 21 {
			totalPoints += 1
		} else {
			totalPoints += 11
		}
	}
	return totalPoints
}

func (hand *Hand) isSplittable() bool {
	cards := hand.cards
	if len(cards) != 2 {
		return false
	}
	if cards[0].Rank == cards[1].Rank {
		return true
	}
	return false
}

func (hand *Hand) String() string {
	cardString := ""
	for idx, c := range hand.cards {
		cardString += fmt.Sprintf("Card %d: %s\n", idx+1, c)
	}
	return fmt.Sprintf("Bet Amount: %d\nHand: %d Points\n%s", hand.betAmount, hand.Points(), cardString)
}
