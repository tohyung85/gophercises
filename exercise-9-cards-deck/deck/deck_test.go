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

	peekCardNumber(1, "A", "S", t, deck)
	peekCardNumber(20, "7", "D", t, deck)
	peekCardNumber(52, "K", "H", t, deck)

	deck = NewDeck(NumberDecks(2), NumberJokers(2), ShuffledDeck(true))

	if deck.CountDeck() == 106 {
		t.Logf("Success: Deck created has 106 cards!\n")
	} else {
		t.Errorf("Failed: Deck has %d cards instead of 106 cards\n", deck.CountDeck())
	}

	jokerPos := deck.FindCardPosition("Joker", "Joker")
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
		return deck.cards[i].Suit < deck.cards[j].Suit
	})

	peekCardNumber(1, "A", "C", t, deck)
	peekCardNumber(52, "K", "S", t, deck)
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

	peekCardNumber(54, "Joker", "Joker", t, deck)
}

func TestRemoveCards(t *testing.T) {
	deck := NewDeck()
	deck.RemoveCard("A", "S")
	checkCardExist("A", "S", t, deck)

	deck.RemoveCardsWithNum("2")
	checkCardExist("2", "S", t, deck)
	checkCardExist("2", "H", t, deck)
	checkCardExist("2", "C", t, deck)
	checkCardExist("2", "D", t, deck)

	deck.RemoveCardsWithSuit("H")
	checkCardExist("5", "H", t, deck)
	checkCardExist("K", "H", t, deck)
	checkCardExist("A", "H", t, deck)
}

func TestDrawDeck(t *testing.T) {
	deck := NewDeck()
	c, err := deck.Draw()
	if err != nil {
		t.Errorf("Failed: error encountered: %s", err)
		return
	}
	if c.Number == "K" && c.Suit == "H" {
		t.Logf("Success: Got last card K of H")
	} else {
		t.Errorf("Failed: Got %s of %s instead of K of H", c.Number, c.Suit)
	}
	if deck.CountDeck() == 51 {
		t.Logf("Success: last card got removed")
	} else {
		t.Errorf("Failed: Deck still has 52 cards")
	}
}

func checkCardExist(num string, suit string, t *testing.T, deck *Deck) {
	cardPos := deck.FindCardPosition(num, suit)
	if cardPos == -1 {
		t.Logf("Success: Removed Card %s of %s no longer in deck", num, suit)
	} else {
		t.Errorf("Failed: %s of %s still in deck", num, suit)
	}
}

func peekCardNumber(toDraw int, expectedNum string, expectedSuit string, t *testing.T, deck *Deck) {
	card, err := deck.Peek(toDraw)
	checkError(err, t)
	if card.Number == expectedNum && card.Suit == expectedSuit {
		t.Logf("Success: card no. %d on the deck is %s of %s", toDraw, card.Number, card.Suit)
	} else {
		t.Errorf("Failed: card no. %d on the deck is %s of %s instead of %s of %s", toDraw, card.Number, card.Suit, expectedNum, expectedSuit)
	}
}

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Errorf("Failed: error encountered: %s", err)
	}
}
