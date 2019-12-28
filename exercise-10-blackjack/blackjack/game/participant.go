package game

import "github.com/tohyung85/gophercises/exercise-9-cards-deck/deck"

type Moves int

const (
	Hit   Moves = 1
	Stand Moves = 2
	Error Moves = 3
)

type Role int

const (
	Dealer Role = iota
	Player
)

type Participant struct {
	id   int
	role Role
	hand []deck.Card
}

func (participant *Participant) handPoints() int {
	aces := make([]deck.Card, 0)
	totalPoints := 0
	for _, c := range participant.hand {
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
