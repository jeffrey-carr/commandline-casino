package mappers

import (
	"fmt"
	"strconv"

	"casino/entities"
)

func RankToValue(rank entities.StandardRank) int {
	switch rank {
	case entities.Ace:
		return 11
	case entities.Ten:
		fallthrough
	case entities.Jack:
		fallthrough
	case entities.Queen:
		fallthrough
	case entities.King:
		return 10
	default:
		v, _ := strconv.ParseInt(string(rank), 0, 64)
		return int(v)
	}
}

func CreateCardByIndex(i int, suit entities.StandardSuit) entities.Card {
	switch i {
	case 0:
		card := CreateCard(entities.Ace, suit)
		card.AltValue = 1
		return card
	case 1:
		return CreateCard(entities.Two, suit)
	case 2:
		return CreateCard(entities.Three, suit)
	case 3:
		return CreateCard(entities.Four, suit)
	case 4:
		return CreateCard(entities.Five, suit)
	case 5:
		return CreateCard(entities.Six, suit)
	case 6:
		return CreateCard(entities.Seven, suit)
	case 7:
		return CreateCard(entities.Eight, suit)
	case 8:
		return CreateCard(entities.Nine, suit)
	case 9:
		return CreateCard(entities.Ten, suit)
	case 10:
		return CreateCard(entities.Jack, suit)
	case 11:
		return CreateCard(entities.Queen, suit)
	case 12:
		return CreateCard(entities.King, suit)
	default:
		return entities.Card{}
	}
}

func CreateCard(
	rank entities.StandardRank,
	suit entities.StandardSuit,
) entities.Card {
	return entities.Card{
		Code:   fmt.Sprintf("%s%s", rank, suit),
		Rank:   rank,
		Suit:   suit,
		Value:  RankToValue(rank),
		Hidden: true,
	}
}
