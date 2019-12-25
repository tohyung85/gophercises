//go:generate stringer -type=Rank,Suit

package deck

import (
	"fmt"
	"math/rand"
	"sort"
)

type Suit int

type Rank int

const (
	Spade Suit = iota
	Diamond
	Club
	Heart
	Joker
)

const (
	_ Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

type Deck struct {
	cards        []Card
	NumberDecks  int
	NumberJokers int
	ShuffledDeck bool
}

type Card struct {
	Rank
	Suit
}

type option func(*Deck)

func (c *Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}

	return fmt.Sprintf("%s of %ss", c.Rank, c.Suit)
}

func (d *Deck) String() string {
	returnStr := ""
	if len(d.cards) < 1 {
		return "There are no cards in the deck!"
	}
	for idx, c := range d.cards {
		returnStr += fmt.Sprintf("%d. %s\n", idx+1, c)
	}
	return returnStr
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
		for s := Spade; s <= Heart; s++ {
			for r := Ace; r <= King; r++ {
				card := Card{r, s}
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

func (d *Deck) FindCardPosition(r Rank, s Suit) int {
	for idx, c := range d.cards {
		if c.Rank == r && c.Suit == s {
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
	d.RemoveCard(c.Rank, c.Suit)
	return c, nil
}

func (d *Deck) RemoveCard(r Rank, s Suit) {
	newDeckOfCards := make([]Card, 0)
	for _, c := range d.cards {
		if c.Rank == r && c.Suit == s {
			continue
		}
		newDeckOfCards = append(newDeckOfCards, c)
	}
	d.cards = newDeckOfCards
}

func (d *Deck) RemoveCardsWithSuit(s Suit) {
	newDeckOfCards := make([]Card, 0)
	for _, c := range d.cards {
		if c.Suit == s {
			continue
		}
		newDeckOfCards = append(newDeckOfCards, c)
	}
	d.cards = newDeckOfCards
}

func (d *Deck) RemoveCardsWithNum(r Rank) {
	newDeckOfCards := make([]Card, 0)
	for _, c := range d.cards {
		if c.Rank == r {
			continue
		}
		newDeckOfCards = append(newDeckOfCards, c)
	}
	d.cards = newDeckOfCards
}

func (d *Deck) AddJokers(numToAdd int) {
	for i := 0; i < numToAdd; i++ {
		joker := Card{Ace, Joker}
		d.cards = append(d.cards, joker)
	}
}
