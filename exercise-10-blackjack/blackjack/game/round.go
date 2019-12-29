package game

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tohyung85/gophercises/exercise-9-cards-deck/deck"
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
	fmt.Printf("Turn over - Dealer:\n%s", hand)
	return nil
}

func (r *Round) processPlayerHand(hand *Hand, idx int) error {
	turnInProgress := true
	fmt.Println("-------------------------------------")
	fmt.Printf("Player %d: \n", hand.belongTo.id)
	fmt.Printf("%s", hand)
	if len(hand.cards) == 2 && hand.Points() == 21 {
		hand.won = true
		turnInProgress = false
	}
	for turnInProgress {
		switch getUserMove(hand) {
		case Hit:
			card, err := r.belongToGame.deck.Draw()
			if err != nil {
				return err
			}
			hand.cards = append(hand.cards, card)
			fmt.Printf("Player drew: %v\n%s", card, hand)
			if hand.Points() > 21 {
				turnInProgress = false
			}
		case Stand:
			fmt.Printf("Player Stands\n")
			turnInProgress = false
		case Split:
			fmt.Printf("Player splits\n")
			r.splitPlayerHand(hand, idx)
		default:
			fmt.Printf("Player Stands\n")
			turnInProgress = false
		}
	}
	fmt.Printf("Turn over - Player:\n%s", hand)
	return nil
}

func (r *Round) splitPlayerHand(h *Hand, idx int) {
	newCards := append(make([]deck.Card, 0), h.cards[1])
	h.cards = append(make([]deck.Card, 0), h.cards[0])
	newHand := &Hand{newCards, h.belongTo, h.betAmount, false}
	h.belongTo.bankRoll -= h.betAmount
	r.drawCardToHand(h)
	r.drawCardToHand(newHand)

	fmt.Printf("Current Hand is now:\n%s", h)
	fmt.Printf("New Hand is now:\n%s", newHand)

	r.hands = append(r.hands, &Hand{})
	copy(r.hands[idx+2:], r.hands[idx+1:])
	r.hands[idx+1] = newHand
}

func getUserMove(hand *Hand) Moves {
	reader := bufio.NewReader(os.Stdin)
	userInp := ""
	for !(userInp == "1" || userInp == "2" || (userInp == "3" && hand.isSplittable())) {
		queryString := ""
		if hand.isSplittable() {
			queryString = "3. Split\n"
		}
		fmt.Printf("Your Move?\n1. Hit\n2. Stand\n%s", queryString)
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
	fmt.Printf("Number of hands: %d", len(round.hands))
	dealerHand := round.hands[len(round.hands)-1]
	dPoints := dealerHand.Points()
	for _, h := range round.hands {
		if h.belongTo.role == Dealer {
			fmt.Printf("Dealer: %d points\n%s", h.Points(), h)
			fmt.Println("-----------------------------")
			continue
		}
		if h.won {
			fmt.Printf("Player %d: %d points\n%sBlackjack!!Get paid $%d\n", h.belongTo.id, h.Points(), h, h.betAmount*2)
			fmt.Println("-----------------------------")
			h.belongTo.getPaid(h.betAmount * 2)
			continue
		}
		if h.Points() > 21 {
			fmt.Printf("Player %d: %d points\n%sYou went burst! Lost $%d\n", h.belongTo.id, h.Points(), h, h.betAmount)
			fmt.Println("-----------------------------")
			continue
		}
		if dPoints > 21 && h.Points() <= 21 {
			fmt.Printf("Player %d: %d points\n%sDealer went burst! You won! Get paid $%d\n", h.belongTo.id, h.Points(), h, h.betAmount*2)
			fmt.Println("-----------------------------")
			h.belongTo.getPaid(h.betAmount * 2)
			h.won = true
			continue
		}
		if dPoints == h.Points() {
			fmt.Printf("Player %d: %d points\n%sIt's a draw! Got back $%d\n", h.belongTo.id, h.Points(), h, h.betAmount*2)
			fmt.Println("-----------------------------")
			h.belongTo.getPaid(h.betAmount)
			continue
		}
		if h.Points() > dPoints {
			fmt.Printf("Player %d: %d points\n%sYou won! Get paid $%d\n", h.belongTo.id, h.Points(), h, h.betAmount*2)
			fmt.Println("-----------------------------")
			h.belongTo.getPaid(h.betAmount * 2)
			h.won = true
			continue
		}
		fmt.Printf("Player %d: %d points\n%sYou lost! Lost $%d\n", h.belongTo.id, h.Points(), h, h.betAmount)
		fmt.Println("-----------------------------")
	}
}
