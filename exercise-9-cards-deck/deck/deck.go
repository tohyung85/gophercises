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
	cards []Card
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
		d.AddDecks(num - 1)
	}
}

func NumberJokers(num int) option {
	return func(d *Deck) {
		d.AddJokers(num)
	}
}

func SortedBy(customFunc func(c []Card) func(i, j int) bool) option {
	return func(d *Deck) {
		d.SortDeckCustom(customFunc)
	}
}

func OmitSuits(args ...Suit) option {
	return func(d *Deck) {
		for _, s := range args {
			d.RemoveCardsWithSuit(s)
		}
	}
}

func OmitRanks(args ...Rank) option {
	return func(d *Deck) {
		for _, r := range args {
			d.RemoveCardsWithNum(r)
		}
	}
}

func OmitCards(args ...Card) option {
	return func(d *Deck) {
		for _, c := range args {
			d.RemoveCard(c.Rank, c.Suit)
		}
	}
}

func ShuffledDeck(shuffle bool) option {
	return func(d *Deck) {
		d.ShuffleDeck()
	}
}

func NewDeck(opts ...option) *Deck { // Please include number of decks option first
	cards := make([]Card, 0)
	deck := &Deck{
		cards: cards,
	}

	deck.AddDecks(1) // You need at least one deck

	for _, opt := range opts {
		opt(deck)
	}

	return deck
}

func (d *Deck) AddDecks(num int) {
	for i := 0; i < num; i++ {
		for s := Spade; s <= Heart; s++ {
			for r := Ace; r <= King; r++ {
				card := Card{r, s}
				d.cards = append(d.cards, card)
			}
		}
	}
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

func (d *Deck) SortDeckDefault() {
	d.sortDeck(defSortFunc)
}

func (d *Deck) SortDeckCustom(less func(c []Card) func(i, j int) bool) {
	d.sortDeck(less)
}

func (d *Deck) ShuffleDeck() {
	d.sortDeck(shuffleSortFunc)
}

func (d *Deck) sortDeck(less func(c []Card) func(i, j int) bool) {
	sort.SliceStable(d.cards, less(d.cards))
}

func defSortFunc(cards []Card) func(i, j int) bool {
	return func(i, j int) bool {
		return absRank(cards[i]) < absRank(cards[j])
	}
}

func absRank(c Card) int {
	return int(c.Suit)*13 + int(c.Rank)
}

func shuffleSortFunc(cards []Card) func(i, j int) bool {
	return func(i, j int) bool {
		r := rand.New(rand.NewSource(99))
		if r.Intn(2) == 1 {
			return true
		} else {
			return false
		}
	}
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
