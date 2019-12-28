package game

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tohyung85/gophercises/exercise-9-cards-deck/deck"
)

type Game struct {
	deck       *deck.Deck
	players    []*Participant
	dealer     *Participant
	inProgress bool
	rounds     []*Round
}

func InitializeGame(decks int, jokers int, shuffle bool, removeSuits []string, removeRanks []string, numPlayers int) *Game {
	gameDeck := initializeDeck(decks, jokers, shuffle, removeSuits, removeRanks)
	players, dealer := initializePlayers(numPlayers)
	game := &Game{gameDeck, players, dealer, true, make([]*Round, 0)}
	fmt.Printf("%s", game)
	return game
}

func initializePlayers(numPlayers int) ([]*Participant, *Participant) {
	participants := make([]*Participant, 0)

	for i := 0; i < numPlayers; i++ {
		player := &Participant{i + 1, Player, 100}
		participants = append(participants, player)
	}
	dealer := &Participant{0, Dealer, 1000000}

	return participants, dealer
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

func (game *Game) Start() error {
	for game.inProgress {
		round, err := game.startNewRound()
		if err != nil {
			return err
		}
		game.rounds = append(game.rounds, round)
		err = round.dealCards()
		if err != nil {
			return err
		}

		round.printStatus()
		for _, h := range round.hands {
			participant := h.belongTo
			if participant.role == Player {
				err := round.processPlayerHand(h)
				if err != nil {
					return err
				}
			} else {
				err := round.processDealerHand(h)
				if err != nil {
					return err
				}
			}
		}
		round.processResults()
		game.displayGameStatus()
		fmt.Printf("Round over. Do you want to go another round? Y to continue and N to quit.\n")
		reader := bufio.NewReader(os.Stdin)
		userInp := ""
		for userInp != "Y" && userInp != "N" {
			userInp, _ = reader.ReadString('\n')
			userInp = strings.TrimSpace(userInp)
		}
		if userInp == "N" {
			game.inProgress = false
		}
	}
	return nil
}

func (g *Game) startNewRound() (*Round, error) {
	hands := make([]*Hand, 0)
	for _, p := range g.players {
		hand := p.buyIn()
		hands = append(hands, hand)
	}
	dealerHand := &Hand{make([]deck.Card, 0), g.dealer, 0, false}
	hands = append(hands, dealerHand)

	return &Round{true, hands, g}, nil
}

func (game *Game) displayGameStatus() {
	fmt.Println("-----------------------------")
	fmt.Println("Current standing:")
	for _, p := range game.players {
		fmt.Printf("Player %d: %d\n", p.id, p.bankRoll)
	}
	fmt.Println("-----------------------------")
}
