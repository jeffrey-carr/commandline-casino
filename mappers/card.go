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
		card := CreateCard(entities.Ace, suit, i)
		card.AltValue = 1
		return card
	case 1:
		return CreateCard(entities.Two, suit, i)
	case 2:
		return CreateCard(entities.Three, suit, i)
	case 3:
		return CreateCard(entities.Four, suit, i)
	case 4:
		return CreateCard(entities.Five, suit, i)
	case 5:
		return CreateCard(entities.Six, suit, i)
	case 6:
		return CreateCard(entities.Seven, suit, i)
	case 7:
		return CreateCard(entities.Eight, suit, i)
	case 8:
		return CreateCard(entities.Nine, suit, i)
	case 9:
		return CreateCard(entities.Ten, suit, i)
	case 10:
		return CreateCard(entities.Jack, suit, i)
	case 11:
		return CreateCard(entities.Queen, suit, i)
	case 12:
		return CreateCard(entities.King, suit, i)
	default:
		return entities.Card{}
	}
}

func CreateCard(
	rank entities.StandardRank,
	suit entities.StandardSuit,
	sortValue int,
) entities.Card {
	return entities.Card{
		Code:      fmt.Sprintf("%s%s", rank, suit),
		Rank:      rank,
		Suit:      suit,
		Value:     RankToValue(rank),
		Hidden:    true,
		SortValue: sortValue,
	}
}
