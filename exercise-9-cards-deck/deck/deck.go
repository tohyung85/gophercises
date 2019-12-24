package deck

import (
	"fmt"
	"math/rand"
	"sort"
)

type Deck struct {
	cards        []Card
	NumberDecks  int
	NumberJokers int
	ShuffledDeck bool
}

type Card struct {
	Number string
	Suit   string
}

type option func(*Deck)

var deckNumbers [13]string = [13]string{
	"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K",
}

var deckSuits [4]string = [4]string{
	"S", "D", "C", "H",
}

func NumberDecks(num int) option {
	return func(d *Deck) {
		d.NumberDecks = num
	}
}

func NumberJokers(num int) option {
	return func(d *Deck) {
		d.NumberJokers = num
	}
}

func ShuffledDeck(shuffle bool) option {
	return func(d *Deck) {
		d.ShuffledDeck = shuffle
	}
}

func NewDeck(opts ...option) *Deck {
	cards := make([]Card, 0)
	deck := &Deck{
		cards:        cards,
		NumberDecks:  1,
		NumberJokers: 0,
		ShuffledDeck: false,
	}

	for _, opt := range opts {
		opt(deck)
	}

	for i := 0; i < deck.NumberDecks; i++ {
		for _, suit := range deckSuits {
			for _, num := range deckNumbers {
				card := Card{num, suit}
				deck.cards = append(deck.cards, card)
			}
		}
	}

	deck.AddJokers(deck.NumberJokers)
	if deck.ShuffledDeck {
		deck.ShuffleDeck()
	}

	return deck
}

func (d *Deck) CountDeck() int {
	return len(d.cards)
}

func (d *Deck) Peek(num int) (Card, error) {
	if num < 1 || num > d.CountDeck() {
		return Card{}, fmt.Errorf("Number must be more than 1 and less than size of deck (%d)", d.CountDeck())
	}
	return d.cards[num-1], nil
}

func (d *Deck) FindCardPosition(num string, suit string) int {
	for idx, c := range d.cards {
		if c.Number == num && c.Suit == suit {
			return idx + 1
		}
	}
	return -1
}

func SortDeck(d *Deck, compareFunc func(i, j int) bool) {
	sort.SliceStable(d.cards, compareFunc)
}

func (d *Deck) ShuffleDeck() {
	sort.Slice(d.cards, func(i, j int) bool {
		r := rand.New(rand.NewSource(99))
		if r.Intn(2) == 1 {
			return true
		} else {
			return false
		}
	})
}

func (d *Deck) Draw() (Card, error) {
	if d.CountDeck() < 1 {
		return Card{}, fmt.Errorf("Deck is empty!")
	}

	c, err := d.Peek(d.CountDeck())
	if err != nil {
		return Card{}, err
	}
	d.RemoveCard(c.Number, c.Suit)
	return c, nil
}

func (d *Deck) RemoveCard(num string, suit string) {
	newDeckOfCards := make([]Card, 0)
	for _, c := range d.cards {
		if c.Number == num && c.Suit == suit {
			continue
		}
		newDeckOfCards = append(newDeckOfCards, c)
	}
	d.cards = newDeckOfCards
}

func (d *Deck) RemoveCardsWithSuit(suit string) {
	newDeckOfCards := make([]Card, 0)
	for _, c := range d.cards {
		if c.Suit == suit {
			continue
		}
		newDeckOfCards = append(newDeckOfCards, c)
	}
	d.cards = newDeckOfCards
}

func (d *Deck) RemoveCardsWithNum(num string) {
	newDeckOfCards := make([]Card, 0)
	for _, c := range d.cards {
		if c.Number == num {
			continue
		}
		newDeckOfCards = append(newDeckOfCards, c)
	}
	d.cards = newDeckOfCards
}

func (d *Deck) AddJokers(numToAdd int) {
	for i := 0; i < numToAdd; i++ {
		joker := Card{"Joker", "Joker"}
		d.cards = append(d.cards, joker)
	}
}
