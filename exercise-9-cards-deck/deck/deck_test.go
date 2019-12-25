package deck

import (
	"testing"
)

func TestNewDeck(t *testing.T) {
	deck := NewDeck()

	if deck.CountDeck() == 52 {
		t.Logf("Success: Deck created has 52 cards!\n")
	} else {
		t.Errorf("Failed: Deck has %d cards instead of 52 cards\n", deck.CountDeck())
	}

	peekCardNumber(1, Ace, Spade, t, deck)
	peekCardNumber(20, Seven, Diamond, t, deck)
	peekCardNumber(52, King, Heart, t, deck)

	deck = NewDeck(NumberDecks(2), NumberJokers(2), ShuffledDeck(true))

	if deck.CountDeck() == 106 {
		t.Logf("Success: Deck created has 106 cards!\n")
	} else {
		t.Errorf("Failed: Deck has %d cards instead of 106 cards\n", deck.CountDeck())
	}

	jokerPos := deck.FindCardPosition(Ace, Joker)
	if jokerPos > -1 {
		t.Logf("Success: There is Joker!\n")
	} else {
		t.Errorf("Failed: There are no Jokers!\n")
	}

	if jokerPos != 105 && jokerPos != 106 {
		t.Logf("Success: Deck was shuffled!\n")
	} else {
		t.Errorf("Failed: Deck not shuffled!\n")
	}
}

func TestSort(t *testing.T) {
	deck := NewDeck()

	SortDeck(deck, func(i, j int) bool {
		return deck.cards[i].Suit.String() < deck.cards[j].Suit.String()
	})

	peekCardNumber(1, Ace, Club, t, deck)
	peekCardNumber(52, King, Spade, t, deck)
}

func TestShuffle(t *testing.T) {
	deck := NewDeck()

	deck.ShuffleDeck()

	// deck.ListCards()

}

func TestAddJokers(t *testing.T) {
	deck := NewDeck()
	numberToAdd := 3
	initNumCards := deck.CountDeck()
	deck.AddJokers(numberToAdd)

	if deck.CountDeck() == initNumCards+numberToAdd {
		t.Logf("Success: 3 additional cards added!")
	} else {
		t.Errorf("Failed: %d cards in deck vs expected of %d", deck.CountDeck(), initNumCards+numberToAdd)
	}

	peekCardNumber(54, Ace, Joker, t, deck)
}

func TestRemoveCards(t *testing.T) {
	deck := NewDeck()
	deck.RemoveCard(Ace, Spade)
	checkCardExist(Ace, Spade, t, deck)

	deck.RemoveCardsWithNum(Two)
	checkCardExist(Two, Spade, t, deck)
	checkCardExist(Two, Heart, t, deck)
	checkCardExist(Two, Club, t, deck)
	checkCardExist(Two, Diamond, t, deck)

	deck.RemoveCardsWithSuit(Heart)
	checkCardExist(Five, Heart, t, deck)
	checkCardExist(King, Heart, t, deck)
	checkCardExist(Ace, Heart, t, deck)
}

func TestDrawDeck(t *testing.T) {
	deck := NewDeck()
	c, err := deck.Draw()
	if err != nil {
		t.Errorf("Failed: error encountered: %s", err)
		return
	}
	if c.Rank == King && c.Suit == Heart {
		t.Logf("Success: Got last card %s", c)
	} else {
		t.Errorf("Failed: Got %s instead of King of Hearts", c)
	}
	if deck.CountDeck() == 51 {
		t.Logf("Success: last card got removed")
	} else {
		t.Errorf("Failed: Deck still has 52 cards")
	}
}

func checkCardExist(r Rank, s Suit, t *testing.T, deck *Deck) {
	cardPos := deck.FindCardPosition(r, s)
	if cardPos == -1 {
		t.Logf("Success: Removed Card %s of %ss no longer in deck", r, s)
	} else {
		t.Errorf("Failed: %s of %s still in deck", r, s)
	}
}

func peekCardNumber(toDraw int, expectedRank Rank, expectedSuit Suit, t *testing.T, deck *Deck) {
	card, err := deck.Peek(toDraw)
	checkError(err, t)
	if card.Rank == expectedRank && card.Suit == expectedSuit {
		t.Logf("Success: card no. %d on the deck is %s of %s", toDraw, card.Rank, card.Suit)
	} else {
		t.Errorf("Failed: card no. %d on the deck is %s instead of %s of %s\n", toDraw, card, expectedRank, expectedSuit)
	}
}

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Errorf("Failed: error encountered: %s", err)
	}
}
