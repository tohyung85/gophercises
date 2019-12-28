package game

import (
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
