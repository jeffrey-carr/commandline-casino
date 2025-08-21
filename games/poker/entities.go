package poker

type PokerHand int

const (
	HighCard      PokerHand = 0
	Pair          PokerHand = 1
	Flush         PokerHand = 2
	Straight      PokerHand = 3
	ThreeOfAKind  PokerHand = 4
	StraightFlush PokerHand = 5
)

var PokerHandToString = map[PokerHand]string{
	HighCard:      "High Card",
	Pair:          "Pair",
	Flush:         "Flush",
	Straight:      "Straight",
	ThreeOfAKind:  "Three of a kind",
	StraightFlush: "Straight flush",
}

var PokerHandToPairPlusMultiplier = map[PokerHand]int{
	Pair:          1,
	Flush:         3,
	Straight:      6,
	ThreeOfAKind:  30,
	StraightFlush: 40,
}
