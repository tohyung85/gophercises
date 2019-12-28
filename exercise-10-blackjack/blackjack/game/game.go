package game

import (
	"bufio"
	"fmt"
	"github.com/tohyung85/gophercises/exercise-9-cards-deck/deck"
	"os"
	"strconv"
	"strings"
)

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

type Game struct {
	deck       *deck.Deck
	players    []*Participant
	inProgress bool
}

type Round struct {
	inProgress bool
}

type Participant struct {
	role Role
	hand []deck.Card
}

func InitializeGame(decks int, jokers int, shuffle bool, removeSuits []string, removeRanks []string) *Game {
	gameDeck := initializeDeck(decks, jokers, shuffle, removeSuits, removeRanks)
	players := initializePlayers(1)
	game := &Game{gameDeck, players, true}
	fmt.Printf("%s", game)
	return game
}

func initializePlayers(numPlayers int) []*Participant {
	participants := make([]*Participant, 0)

	for i := 0; i < numPlayers; i++ {
		playerHand := make([]deck.Card, 0)
		player := &Participant{Player, playerHand}
		participants = append(participants, player)
	}
	dealerHand := make([]deck.Card, 0)
	dealer := &Participant{Dealer, dealerHand}
	participants = append(participants, dealer)

	return participants
}

func initializeDeck(decks int, jokers int, shuffle bool, removeSuits []string, removeRanks []string) *deck.Deck {
	fmt.Println("Initializing deck with following parameters:")
	fmt.Printf("%d decks\n", decks)
	fmt.Printf("%d jokers\n", jokers)
	fmt.Printf("Deck to be shuffled? %t\n", shuffle)
	fmt.Printf("Suits %v will be removed\n", removeSuits)
	fmt.Printf("Ranks %v will be removed\n", removeRanks)
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

func (game *Game) Start() error {
	var err error
	for game.inProgress {
		reader := bufio.NewReader(os.Stdin)
		round := &Round{true}
		game.clearHands()
		err = game.dealCards()

		if err != nil {
			return err
		}
		for round.inProgress {
			game.printGameStatus(false)
			for idx, player := range game.players {
				turnInProgress := true
				for turnInProgress {
					var move Moves
					if player.role == Player {
						fmt.Printf("Player %d: ", idx+1)
						move, err = game.processMoveFor(player)
						fmt.Printf("Player %d's hand: %v (%d points)\n", idx+1, player.hand, player.handPoints())
					} else {
						fmt.Printf("Dealer: ")
						move, err = game.processDealerMove(player)
						fmt.Printf("Dealer's hand: %v (%d points)\n", player.hand, player.handPoints())
					}
					if err != nil {
						return err
					}
					if move == Stand || player.handPoints() > 21 {
						turnInProgress = false
					}
				}

			}
			// Process game results
			// Show end result
			game.printGameStatus(true)
			round.inProgress = false
		}
		fmt.Printf("Round over. Do you want to go another round? Y to continue and N to quit.\n")
		userInp := ""
		for userInp != "Y" && userInp != "N" {
			userInp, _ = reader.ReadString('\n')
			userInp = strings.TrimSpace(userInp)
		}
		if userInp == "N" {
			game.inProgress = false
		}
	}
	err = nil
	return err
}

func (game *Game) processResults() {

}

func (game *Game) processMoveFor(player *Participant) (Moves, error) {
	reader := bufio.NewReader(os.Stdin)
	userInp := ""
	for userInp != "1" && userInp != "2" {
		fmt.Printf("Your Move?\n1. Hit\n2. Stand\n")
		userInp, _ = reader.ReadString('\n')
		userInp = strings.TrimSpace(userInp)
	}
	inpInt, _ := strconv.Atoi(userInp)
	if Moves(inpInt) == Hit {
		card, err := game.deck.Draw()
		if err != nil {
			return Error, err
		}
		player.hand = append(player.hand, card)
		fmt.Printf("Player drew: %v\n", card)
	} else {
		fmt.Printf("Player Stands\n")
	}
	return Moves(inpInt), nil
}

func (game *Game) processDealerMove(player *Participant) (Moves, error) {
	points := player.handPoints()
	var move Moves
	if points < 17 {
		move = Hit
		card, err := game.deck.Draw()
		if err != nil {
			return Error, err
		}
		fmt.Printf("Dealer drew: %v\n", card)
		player.hand = append(player.hand, card)
	} else {
		move = Stand
		fmt.Printf("Dealer Stands\n")
	}
	return move, nil
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

func (game *Game) clearHands() {
	for _, p := range game.players {
		hand := make([]deck.Card, 0)
		p.hand = hand
	}
}

func (game *Game) dealCards() error {
	for i := 0; i < 2; i++ {
		for _, p := range game.players {
			card, err := game.deck.Draw()
			if err != nil {
				return err
			}
			p.hand = append(p.hand, card)
		}
	}
	return nil
}

func (game *Game) printGameStatus(finalStatus bool) {
	dealer := game.players[len(game.players)-1]
	if finalStatus {
		fmt.Printf("Dealer: %d Points\n", dealer.handPoints())
		for idx, c := range dealer.hand {
			fmt.Printf("Card %d: %s\n", idx+1, c.String())
		}
	} else {
		fmt.Printf("Dealer:\n")
		for idx, c := range dealer.hand {
			if idx == 0 {
				fmt.Printf("Card %d: hidden (%s)\n", idx+1, c.String())
			} else {
				fmt.Printf("Card %d: %s\n", idx+1, c.String())
			}
		}
	}

	for idx, p := range game.players {
		if p.role == Dealer {
			continue
		}
		if finalStatus {
			wins := "Wins!"
			if p.handPoints() > 21 {
				wins = "Lost!"
			} else {
				if dealer.handPoints() > 21 {
					wins = "Wins"
				}
				if p.handPoints() < dealer.handPoints() && dealer.handPoints() < 22 {
					wins = "Lost!"
				}
				if p.handPoints() == dealer.handPoints() {
					wins = "Draw!"
				}
			}
			fmt.Printf("Player %d %s: %d Points\n", idx+1, wins, p.handPoints())
		} else {
			fmt.Printf("Player %d:\n", idx+1)
		}
		for idx, c := range p.hand {
			fmt.Printf("Card %d: %s\n", idx+1, c.String())
		}
	}
}
