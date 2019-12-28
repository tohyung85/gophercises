package game

import (
	"bufio"
	"fmt"
	"github.com/tohyung85/gophercises/exercise-9-cards-deck/deck"
	"os"
	"strconv"
	"strings"
)

type Game struct {
	deck       *deck.Deck
	players    []*Participant
	dealer     *Participant
	inProgress bool
	rounds     []*Round
}

type Round struct {
	inProgress bool
	winners    map[int]struct{}
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
		playerHand := make([]deck.Card, 0)
		player := &Participant{i + 1, Player, playerHand}
		participants = append(participants, player)
	}
	dealerHand := make([]deck.Card, 0)
	dealer := &Participant{0, Dealer, dealerHand}

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
	var err error
	for game.inProgress {
		reader := bufio.NewReader(os.Stdin)
		round := &Round{true, make(map[int]struct{})}
		game.rounds = append(game.rounds, round)
		game.clearHands()
		err = game.dealCards()

		if err != nil {
			return err
		}
		for round.inProgress {
			game.printGameStatus(false)
			for _, p := range game.players {
				err = game.processMoveFor(p)
				if err != nil {
					return err
				}
			}
			err = game.processDealerMove(game.dealer)
			if err != nil {
				return err
			}
			game.processResults()
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
	dPoints := game.dealer.handPoints()
	currRound := game.rounds[len(game.rounds)-1]
	for _, p := range game.players {
		pPoints := p.handPoints()
		if pPoints > 21 {
			continue
		}
		if dPoints > 21 || pPoints > dPoints {
			currRound.winners[p.id] = struct{}{}
			continue
		}
	}
}

func (game *Game) processMoveFor(player *Participant) error {
	turnInProgress := true
	fmt.Println("-------------------------------------")
	fmt.Printf("Player %d: ", player.id)
	if len(player.hand) == 2 && player.handPoints() == 21 {
		game.rounds[len(game.rounds)-1].winners[player.id] = struct{}{}
		turnInProgress = false
	}
	reader := bufio.NewReader(os.Stdin)
	for turnInProgress {
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
				return err
			}
			player.hand = append(player.hand, card)
			fmt.Printf("Player drew: %v\n", card)
			if player.handPoints() > 21 {
				turnInProgress = false
			}
		} else {
			fmt.Printf("Player Stands\n")
			turnInProgress = false
		}
	}
	fmt.Printf("Turn over - Player %d's hand: %v (%d points)\n", player.id, player.hand, player.handPoints())
	return nil
}

func (game *Game) processDealerMove(player *Participant) error {
	fmt.Println("-------------------------------------")
	fmt.Printf("Dealer: ")
	turnInProgress := true
	for turnInProgress {
		points := player.handPoints()
		if points < 17 {
			card, err := game.deck.Draw()
			if err != nil {
				return err
			}
			fmt.Printf("Dealer drew: %v\n", card)
			player.hand = append(player.hand, card)
		} else {
			turnInProgress = false
			fmt.Printf("Dealer Stands\n")
		}
	}
	fmt.Printf("Turn over - Dealer's hand: %v (%d points)\n", player.hand, player.handPoints())
	return nil
}

func (game *Game) clearHands() {
	for _, p := range game.players {
		hand := make([]deck.Card, 0)
		p.hand = hand
	}
	game.dealer.hand = make([]deck.Card, 0)
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
		card, err := game.deck.Draw()
		if err != nil {
			return err
		}
		game.dealer.hand = append(game.dealer.hand, card)
	}
	return nil
}

func (game *Game) printGameStatus(finalStatus bool) {
	dealer := game.dealer
	currRound := game.rounds[len(game.rounds)-1]
	fmt.Println("-------------------------------------")
	if finalStatus {
		fmt.Printf("Dealer: %d Points\n", dealer.handPoints())
		for idx, c := range dealer.hand {
			fmt.Printf("Card %d: %s\n", idx+1, c.String())
		}
	} else {
		fmt.Println("-------------------------------------")
		fmt.Printf("Dealer:\n")
		for idx, c := range dealer.hand {
			if idx == 0 {
				fmt.Printf("Card %d: hidden (%s)\n", idx+1, c.String())
			} else {
				fmt.Printf("Card %d: %s\n", idx+1, c.String())
			}
		}
	}

	for _, p := range game.players {
		fmt.Println("-------------------------------------")
		wins := "Lost!"
		if finalStatus {
			_, won := currRound.winners[p.id]
			if won {
				wins = "Wins!"
			}
			fmt.Printf("Player %d %s: %d Points\n", p.id, wins, p.handPoints())
		} else {
			fmt.Printf("Player %d:\n", p.id)
		}
		for idx, c := range p.hand {
			fmt.Printf("Card %d: %s\n", idx+1, c.String())
		}
	}
}
