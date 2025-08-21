package entities

import "fmt"

type StandardSuit string

const (
	Spade   StandardSuit = "♠"
	Club    StandardSuit = "♣"
	Heart   StandardSuit = "♥"
	Diamond StandardSuit = "♦"
)

var AllSuits = []StandardSuit{Spade, Club, Heart, Diamond}

type StandardRank string

const (
	Ace   StandardRank = "A"
	Two   StandardRank = "2"
	Three StandardRank = "3"
	Four  StandardRank = "4"
	Five  StandardRank = "5"
	Six   StandardRank = "6"
	Seven StandardRank = "7"
	Eight StandardRank = "8"
	Nine  StandardRank = "9"
	Ten   StandardRank = "10"
	Jack  StandardRank = "J"
	Queen StandardRank = "Q"
	King  StandardRank = "K"
)

type Card struct {
	Code      string
	Rank      StandardRank
	Suit      StandardSuit
	Value     int
	AltValue  int
	IsAlt     bool
	Hidden    bool
	SortValue int
}

func (c Card) GetValue() int {
	if c.IsAlt {
		return c.AltValue
	}

	return c.Value
}

type Hand struct {
	Cards []Card
}

func (h Hand) HasHidden() bool {
	for _, card := range h.Cards {
		if card.Hidden {
			return true
		}
	}

	return false
}

func (h Hand) String() string {
	str := ""
	for _, card := range h.Cards {
		if card.Hidden {
			str += " [?]"
		} else {
			str = fmt.Sprintf("%s %s", str, card.Code)
		}
	}

	return str
}

type Deck struct {
	DrawPile    []Card
	DiscardPile []Card
}

type Role string

const (
	UserRole   Role = "user"
	DealerRole Role = "dealer"
)

func (h *Hand) TotalValue() int {
	total := 0
	for _, card := range h.Cards {
		if card.Hidden {
			continue
		}

		total += card.Value
	}

	return total
}

type ShuffleOpts struct {
	NumDecks int
}

type DrawOpts struct {
	Count int
}
