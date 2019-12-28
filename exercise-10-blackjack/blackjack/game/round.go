package game

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Round struct {
	inProgress   bool
	hands        []*Hand
	belongToGame *Game
}

func (r *Round) dealCards() error {
	for i := 0; i < 2; i++ {
		for _, h := range r.hands {
			err := r.drawCardToHand(h)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Round) drawCardToHand(hand *Hand) error {
	card, err := r.belongToGame.deck.Draw()
	hand.cards = append(hand.cards, card)
	return err
}

func (r *Round) processDealerHand(hand *Hand) error {
	fmt.Println("-------------------------------------")
	fmt.Printf("Dealer: ")
	turnInProgress := true
	for turnInProgress {
		points := hand.Points()
		if points < 17 {
			card, err := r.belongToGame.deck.Draw()
			if err != nil {
				return err
			}
			fmt.Printf("Dealer drew: %v\n", card)
			hand.cards = append(hand.cards, card)
		} else {
			turnInProgress = false
			fmt.Printf("Dealer Stands\n")
		}
	}
	fmt.Printf("Turn over - Dealer's hand: %v (%d points)\n", hand.cards, hand.Points())
	return nil
}

func (r *Round) processPlayerHand(hand *Hand) error {
	turnInProgress := true
	fmt.Println("-------------------------------------")
	fmt.Printf("Player %d: ", hand.belongTo.id)
	if len(hand.cards) == 2 && hand.Points() == 21 {
		hand.won = true
		turnInProgress = false
	}
	for turnInProgress {
		switch getUserMove() {
		case Hit:
			card, err := r.belongToGame.deck.Draw()
			if err != nil {
				return err
			}
			hand.cards = append(hand.cards, card)
			fmt.Printf("Player drew: %v\n", card)
			if hand.Points() > 21 {
				turnInProgress = false
			}
		case Stand:
			fmt.Printf("Player Stands\n")
			turnInProgress = false
		default:
			fmt.Printf("Player Stands\n")
			turnInProgress = false
		}
	}
	fmt.Printf("Turn over - Player %d's hand: %s (%d points)\n", hand.belongTo.id, hand.cards, hand.Points())
	return nil
}

func getUserMove() Moves {
	reader := bufio.NewReader(os.Stdin)
	userInp := ""
	for userInp != "1" && userInp != "2" {
		fmt.Printf("Your Move?\n1. Hit\n2. Stand\n")
		userInp, _ = reader.ReadString('\n')
		userInp = strings.TrimSpace(userInp)
	}
	inpInt, _ := strconv.Atoi(userInp)

	return Moves(inpInt)
}

func (round *Round) printStatus() {
	fmt.Println("Printing status...")
	for _, h := range round.hands {
		fmt.Println("-------------------------------------")
		participant := h.belongTo
		if participant.role == Dealer {
			fmt.Printf("Dealer:\n")
			for idx, c := range h.cards {
				if idx == 0 {
					fmt.Printf("Card %d: hidden (%s)\n", idx+1, c)
				} else {
					fmt.Printf("Card %d: %s\n", idx+1, c)
				}
			}
			continue
		}

		fmt.Printf("Player %d:\n", participant.id)
		for idx, c := range h.cards {
			fmt.Printf("Card %d: %s\n", idx+1, c)
		}
	}
}

func (round *Round) processResults() {
	fmt.Println("-----------------------------")
	fmt.Println("Round over!")
	dealerHand := round.hands[len(round.hands)-1]
	dPoints := dealerHand.Points()
	for _, h := range round.hands {
		if h.belongTo.role == Dealer {
			fmt.Printf("Dealer has a hand of %s (%d points)\n", h.cards, h.Points())
			continue
		}
		if h.won {
			fmt.Printf("Player %d with hand %s (%d points) got a blackjack! Get paid $%d\n", h.belongTo.id, h.cards, h.Points(), h.betAmount*2)
			h.belongTo.getPaid(h.betAmount * 2)
			continue
		}
		if h.Points() > 21 {
			fmt.Printf("Player %d with hand %s (%d points) went burst! Lost $%d\n", h.belongTo.id, h.cards, h.Points(), h.betAmount)
			continue
		}
		if dPoints > 21 && h.Points() < 21 {
			fmt.Printf("Player %d with hand %s (%d points) won! Get paid $%d\n", h.belongTo.id, h.cards, h.Points(), h.betAmount*2)
			h.belongTo.getPaid(h.betAmount * 2)
			h.won = true
			continue
		}
		if dPoints == h.Points() {
			fmt.Printf("Player %d with hand %s (%d points) drew! Got back $%d\n", h.belongTo.id, h.cards, h.Points(), h.betAmount)
			h.belongTo.getPaid(h.betAmount)
			continue
		}
		if h.Points() > dPoints {
			fmt.Printf("Player %d with hand %s (%d points) won! Won $%d\n", h.belongTo.id, h.cards, h.Points(), h.betAmount*2)
			h.belongTo.getPaid(h.betAmount * 2)
			h.won = true
		}
	}
}
