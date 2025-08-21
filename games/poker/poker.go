package poker

import (
	"fmt"

	"casino/entities"
	"casino/games"
	"casino/utils"
)

type WinLevel int

const (
	WinLevelDealer WinLevel = -1
	WinLevelPush   WinLevel = 0
	WinLevelUser   WinLevel = 1
)

type poker struct {
	dealer      utils.Dealer
	saveManager utils.SaveDataManager

	hands     map[entities.Role]entities.Hand
	userChips int
	ante      int
	pairPlus  int

	in  chan string
	out chan string

	quit func()
}

func NewPoker(
	dealer utils.Dealer,
	saveManager utils.SaveDataManager,
	in chan string,
	out chan string,
	quit func(),
) games.Game {
	return &poker{
		dealer:      dealer,
		saveManager: saveManager,
		in:          in,
		out:         out,
		quit:        quit,
	}
}

func (p *poker) Name() string {
	return "3-Card Poker"
}

func (p *poker) Play() {
	utils.Clear(p.out)
	save := p.saveManager.Read()
	p.userChips = save.RemainingChips

	utils.PrintBanner(p.Name(), p.out)
	p.wager()
}

func (p *poker) wager() {
	save := p.saveManager.Read()
	p.out <- utils.Dim(fmt.Sprintf("You have %d chips", save.RemainingChips))
	p.ante = utils.GetBet(p.in, p.out, fmt.Sprintf("How much do you want to ante? (max %d)", p.userChips/2), 1, p.userChips/2)
	p.userChips -= p.ante
	p.out <- "Pair plus payouts:"
	for hand, multiplier := range PokerHandToPairPlusMultiplier {
		p.out <- fmt.Sprintf("\t%s: %d to 1", PokerHandToString[hand], multiplier)
	}
	p.pairPlus = utils.GetBet(p.in, p.out, fmt.Sprintf("Pair plus? (max %d)", p.userChips/2), 0, p.userChips/2)
	p.userChips -= p.pairPlus

	save.RemainingChips = p.userChips
	p.saveManager.Save(save)

	p.deal()
}

func (p *poker) deal() {
	p.out <- utils.Dim("Shuffling the deck...")
	p.dealer.Shuffle()
	p.hands = map[entities.Role]entities.Hand{}

	// Create dealer initial hand and user initial hand
	for i := range 6 {
		role := entities.DealerRole
		hidden := true
		if i%2 == 0 {
			role = entities.UserRole
			hidden = false
		}

		hand := p.hands[role]
		card := p.dealer.Draw()
		card.Hidden = hidden
		hand.Cards = append(hand.Cards, card)
		p.hands[role] = hand
	}

	p.printUpdate()
	if p.pairPlus > 0 {
		userHandLevel := DeterminePokerHandLevel(p.hands[entities.UserRole])
		p.payoutPairPlus(userHandLevel)
	}

	p.lastChance()
}

func (p *poker) payoutPairPlus(level PokerHand) {
	multiplier := PokerHandToPairPlusMultiplier[level]
	bonus := p.pairPlus * multiplier

	if multiplier == 0 {
		p.out <- utils.Red("You lose your pair plus bet (-%d chips)", p.pairPlus)
		return
	}

	p.out <- utils.Dim("%s pays out %d to 1", PokerHandToString[level], multiplier)
	p.out <- utils.Green(utils.Bold("You win your pair plus bet! (+%d chips)", bonus))

	save := p.saveManager.Read()
	save.RemainingChips += p.pairPlus * multiplier
	p.saveManager.Save(save)
	p.userChips = save.RemainingChips
}

func (p *poker) lastChance() {
	userChoice := utils.GetInput(
		p.in,
		p.out,
		[]string{"play", "p", "fold", "f"},
		fmt.Sprintf("Play (p, %d chips) or fold (f)?", p.ante),
	)

	var folded bool
	switch userChoice {
	case "play", "p":
		save := p.saveManager.Read()
		save.RemainingChips -= p.ante
		p.saveManager.Save(save)
		p.userChips = save.RemainingChips
		folded = false
	case "fold", "f":
		folded = true
	}
	p.compareHands(folded)
}

func (p *poker) compareHands(folded bool) {
	for i := range p.hands[entities.DealerRole].Cards {
		p.hands[entities.DealerRole].Cards[i].Hidden = false
	}

	p.printUpdate()
	if folded {
		p.out <- utils.Red("You %s", utils.Bold("folded"))
	}

	dealerLevel := DeterminePokerHandLevel(p.hands[entities.DealerRole])
	userLevel := DeterminePokerHandLevel(p.hands[entities.UserRole])

	dealerStr := fmt.Sprintf("Dealer has %s", PokerHandToString[dealerLevel])
	if dealerLevel == HighCard {
		dealerHigh, _ := utils.MaxFunc(p.hands[entities.DealerRole].Cards, func(c entities.Card) int {
			return c.SortValue
		})
		dealerStr += fmt.Sprintf(" (%s high)", dealerHigh.Rank)
	}
	p.out <- utils.Dim(dealerStr)
	userStr := fmt.Sprintf("You have %s", PokerHandToString[userLevel])
	if userLevel == HighCard {
		userHigh, _ := utils.MaxFunc(p.hands[entities.UserRole].Cards, func(c entities.Card) int {
			return c.SortValue
		})
		userStr += fmt.Sprintf(" (%s high)", userHigh.Rank)
	}
	p.out <- utils.Dim(userStr)

	bonus := 0
	if userLevel-dealerLevel < 0 {
		loss := p.ante * 2
		if folded {
			loss = p.ante
		}
		p.out <- utils.Red(utils.Bold("Dealer wins (-%d chips)", loss))
	} else if userLevel-dealerLevel == 0 {
		winner, _ := ResolvePush(p.hands, userLevel)
		if folded || winner == entities.DealerRole {
			loss := p.ante * 2
			if folded {
				loss = p.ante
			}
			p.out <- utils.Red(utils.Bold("Dealer wins (-%d chips)", loss))
		} else if winner == entities.UserRole {
			bonus = p.ante * 3
			p.out <- utils.Green("You win! (+%d chips)", bonus)
		} else {
			p.out <- "Push, you get your chips back!"
			bonus = p.ante * 2
		}
	} else {
		bonus = p.ante * 3
		p.out <- utils.Green("You win! (+%d chips)", bonus)
	}

	save := p.saveManager.Read()
	save.RemainingChips += bonus
	p.saveManager.Save(save)
	p.userChips = save.RemainingChips

	p.endGame()
}

func (p *poker) endGame() {
	for _, hand := range p.hands {
		p.dealer.Discard(hand.Cards...)
	}
	p.hands = nil

	playAgainChoice := utils.GetInput(
		p.in,
		p.out,
		[]string{"yes", "y", "no", "n"},
		"Play again? Yes (y) or no (n)",
	)
	switch playAgainChoice {
	case "yes", "y":
		p.Play()
		return
	default:
		p.quit()
	}
}

func (p *poker) printUpdate() {
	utils.Clear(p.out)

	p.out <- utils.Dim("Your chips: %d", p.userChips)
	p.out <- utils.Dim("Anted: %d", p.ante)
	p.out <- utils.Dim("Pair Plus: %d", p.pairPlus)
	p.out <- utils.Divider()

	dealerHand := p.hands[entities.DealerRole]
	userHand := p.hands[entities.UserRole]

	for _, line := range utils.RenderHand(dealerHand) {
		p.out <- line
	}
	for _, line := range utils.RenderHand(userHand) {
		p.out <- line
	}

	p.out <- utils.Divider()
}
