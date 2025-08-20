package utils

import (
	"casino/entities"
	"casino/mappers"
)

type Dealer struct {
	CurrentDeck entities.Deck
}

func NewDealer() Dealer {
	deck := GenerateStandardDeck()
	dealer := Dealer{
		CurrentDeck: deck,
	}
	dealer.Shuffle()

	return dealer
}

func (d *Dealer) Shuffle() {
	allCards := append(d.CurrentDeck.DrawPile, d.CurrentDeck.DiscardPile...)
	Shuffle(allCards)
	d.CurrentDeck.DrawPile = allCards
}

func (d *Dealer) Draw() entities.Card {
	if len(d.CurrentDeck.DrawPile) == 0 {
		return entities.Card{}
	}

	card := d.CurrentDeck.DrawPile[0]
	card.Hidden = false
	d.CurrentDeck.DrawPile = d.CurrentDeck.DrawPile[1:]

	return card
}

func (d *Dealer) Discard(cards ...entities.Card) {
	for _, card := range cards {
		card.Hidden = true
		d.CurrentDeck.DiscardPile = append(d.CurrentDeck.DiscardPile, card)
	}
}

func GenerateStandardDeck() entities.Deck {
	// A standard deck has 4 suits, A through K of each suite totalling 52 cards
	cardsBySuit := Map(entities.AllSuits, func(_ int, suit entities.StandardSuit) []entities.Card {
		var cards []entities.Card
		for i := range 13 {
			cards = append(cards, mappers.CreateCardByIndex(i, suit))
		}
		return cards
	})
	cards := Flatten(cardsBySuit)

	return entities.Deck{
		DrawPile: cards,
	}
}
