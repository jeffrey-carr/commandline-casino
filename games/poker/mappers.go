package poker

import (
	"casino/entities"
	"casino/utils"
	"slices"
)

func DeterminePokerHandLevel(hand entities.Hand) PokerHand {
	cards := hand.Cards
	slices.SortFunc(cards, func(card1, card2 entities.Card) int {
		return card1.SortValue - card2.SortValue
	})
	isStraight := true
	isFlush := true
	for i := range cards[:len(cards)-1] {
		card := cards[i]
		nextCard := cards[i+1]
		if card.SortValue+1 != nextCard.SortValue {
			isStraight = false
		}
		if card.Suit != nextCard.Suit {
			isFlush = false
		}
	}

	if isFlush && isStraight {
		return StraightFlush
	}

	uniqueCardValues := make(map[int]struct{}, len(hand.Cards))
	for _, card := range cards {
		uniqueCardValues[card.SortValue] = struct{}{}
	}
	if len(uniqueCardValues) == 1 {
		return ThreeOfAKind
	}

	if isStraight {
		return Straight
	}

	if isFlush {
		return Flush
	}

	if len(uniqueCardValues) == 2 {
		return Pair
	}

	return HighCard
}

func ResolvePush(hands map[entities.Role]entities.Hand, handLevel PokerHand) (entities.Role, bool) {
	dealerHand := hands[entities.DealerRole]
	userHand := hands[entities.UserRole]

	if handLevel == Pair {
		// Figure out what the pair is
		dealerPairValue := GetPairValues(dealerHand)
		userPairValue := GetPairValues(userHand)
		if dealerPairValue > userPairValue {
			return entities.DealerRole, true
		} else if dealerPairValue < userPairValue {
			return entities.UserRole, true
		} else {
			return entities.Role("none"), false
		}
	}

	// Otherwise, we can just figure out what the highest card is
	dealerHighest, _ := utils.MaxFunc(dealerHand.Cards, func(card entities.Card) int {
		return card.SortValue
	})
	userHighest, _ := utils.MaxFunc(userHand.Cards, func(card entities.Card) int {
		return card.SortValue
	})

	if dealerHighest.SortValue > userHighest.SortValue {
		return entities.DealerRole, true
	}
	if dealerHighest.SortValue < userHighest.SortValue {
		return entities.UserRole, true
	}

	return entities.Role("none"), false
}

func GetPairValues(hand entities.Hand) int {
	if hand.Cards[0].SortValue == hand.Cards[1].SortValue {
		return hand.Cards[0].SortValue
	}

	return hand.Cards[2].SortValue
}
