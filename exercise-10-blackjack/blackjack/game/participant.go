package game

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tohyung85/gophercises/exercise-9-cards-deck/deck"
)

type Moves int

const (
	Hit   Moves = 1
	Stand Moves = 2
	Error Moves = 3
)

type Role string

const (
	Dealer Role = "Dealer"
	Player Role = "Player"
)

type Participant struct {
	id       int
	role     Role
	bankRoll int
}

func (player *Participant) buyIn() *Hand {
	reader := bufio.NewReader(os.Stdin)
	inputOk := false
	var hand *Hand
	for !inputOk {
		fmt.Printf("Player %d please place bet amount: ", player.id)
		userInp, _ := reader.ReadString('\n')
		userInp = strings.TrimSpace(userInp)
		betAmount, err := strconv.Atoi(userInp)
		if err == nil {
			player.bankRoll -= betAmount
			hand = &Hand{make([]deck.Card, 0), player, betAmount, false}
			inputOk = true
		}
	}
	return hand
}

func (player *Participant) getPaid(winnings int) {
	player.bankRoll += winnings
}
