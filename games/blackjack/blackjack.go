package blackjack

import (
	"errors"
	"fmt"
	"strings"

	"casino/entities"
	"casino/games"
	"casino/utils"
)

const DealerStandValue = 17

type blackjack struct {
	dealer      utils.Dealer
	saveManager utils.SaveDataManager

	hands     map[entities.Role]entities.Hand
	wager     int
	userChips int

	in  chan string
	out chan string

	quit func()
}

func NewBlackjack(
	dealer utils.Dealer,
	saveManager utils.SaveDataManager,
	in chan string,
	out chan string,
	quit func(),
) games.Game {
	return &blackjack{
		dealer:      dealer,
		saveManager: saveManager,
		in:          in,
		out:         out,
		quit:        quit,
	}
}

func (b *blackjack) Name() string {
	return "Blackjack"
}

func (b *blackjack) Play() {
	utils.Clear(b.out)
	save := b.saveManager.Read()
	b.userChips = save.RemainingChips

	utils.PrintBanner(b.Name(), b.out)
	b.out <- fmt.Sprintf("%sDealer stands on %d%s", utils.AnsiDim, DealerStandValue, utils.AnsiReset)
	b.out <- utils.Dim(fmt.Sprintf("You have %d chips remaining", save.RemainingChips))
	b.wager = utils.GetBet(b.in, b.out, "How much would you like to wager?", 1, save.RemainingChips)

	b.start()
}

func (b *blackjack) start() {
	b.out <- utils.Yellow(utils.Bold(fmt.Sprintf("Wagered %d chips", b.wager)))
	save := b.saveManager.Read()
	save.RemainingChips -= b.wager
	b.userChips = save.RemainingChips
	b.saveManager.Save(save)

	b.out <- utils.Dim("Shuffling the deck...")
	b.dealer.Shuffle()
	b.hands = map[entities.Role]entities.Hand{}

	// Create dealer initial hand and user initial hand
	for i := range 4 {
		role := entities.DealerRole
		if i%2 == 0 {
			role = entities.UserRole
		}

		hand := b.hands[role]
		hand.Cards = append(hand.Cards, b.dealer.Draw())
		b.hands[role] = hand
	}

	b.hands[entities.DealerRole].Cards[1].Hidden = true

	b.printUpdate()
	b.run()
}

func (b *blackjack) run() {
	b.out <- ""

	if b.sumHand(b.hands[entities.UserRole]) == 21 {
		b.out <- utils.Green(utils.Bold("BLACKJACK"))
		b.endGame()
		return
	}

	b.out <- utils.Cyan("Your move → ") + utils.Bold("Hit (h)") + " / " + utils.Bold("Stay (s)")

	for line := range b.in {
		switch strings.ToLower(strings.TrimSpace(line)) {
		case "hit":
			fallthrough
		case "h":
			bust, err := b.hit(entities.UserRole)
			if err != nil {
				b.out <- utils.Red(err.Error())
				continue
			}

			b.printUpdate()

			if !bust {
				b.out <- utils.Cyan("Your move → ") + utils.Bold("Hit (h)") + " / " + utils.Bold("Stay (s)")
				continue
			}

			b.out <- utils.Red("BUST!")
		case "stay":
		case "s":
		default:
			b.out <- utils.Yellow(fmt.Sprintf("Unknown command: %s", line))
			b.out <- utils.Cyan("Your move → ") + utils.Bold("Hit (h)") + " / " + utils.Bold("Stay (s)")
			continue
		}

		break
	}

	b.runDealer()
}

func (b *blackjack) runDealer() {
	if b.sumHand(b.hands[entities.UserRole]) > 21 {
		b.endGame()
		return
	}

	// Reveal all dealer cards
	dealerHand := b.hands[entities.DealerRole]
	dealerHand.Cards = utils.Map(dealerHand.Cards, func(_ int, card entities.Card) entities.Card {
		card.Hidden = false
		return card
	})
	b.hands[entities.DealerRole] = dealerHand

	if b.sumHand(dealerHand) == 21 {
		b.out <- utils.Red(utils.Bold("DEALER BLACKJACK"))
		b.printUpdate()
		b.endGame()
	}

	for b.sumHand(b.hands[entities.DealerRole]) < DealerStandValue {
		b.hit(entities.DealerRole)
	}

	b.printUpdate()

	if b.sumHand(b.hands[entities.DealerRole]) > 21 {
		b.out <- utils.Green("DEALER BUST!")
	}

	b.endGame()
}

func (b *blackjack) endGame() {
	dealerHand := b.hands[entities.DealerRole]
	userHand := b.hands[entities.UserRole]
	dealerShowing := b.sumHand(dealerHand)
	userShowing := b.sumHand(userHand)

	userWins := userShowing <= 21 && (userShowing > dealerShowing || dealerShowing > 21)
	tie := dealerShowing == userShowing

	winnings := 0
	if userWins {
		winnings = b.wager * 2
		b.out <- utils.Green(utils.Bold(fmt.Sprintf("YOU WIN! +%d chips", winnings)))
	} else if tie {
		winnings = b.wager
		b.out <- utils.Yellow(utils.Bold(fmt.Sprintf("Tie - win your chips back (+%d chips)", winnings)))
	} else {
		b.out <- utils.Red(utils.Bold(fmt.Sprintf("Dealer wins (-%d chips)", b.wager)))
	}

	stats := b.saveManager.Read()
	stats.RemainingChips += winnings
	b.saveManager.Save(stats)

	b.out <- fmt.Sprintf("New total: %d", stats.RemainingChips)
	b.out <- fmt.Sprintf(
		"%sYou%s: %d\t%sDealer%s: %d",
		utils.AnsiBold, utils.AnsiReset, userShowing,
		utils.AnsiBold, utils.AnsiReset, dealerShowing,
	)

	b.wager = 0
	// Discard all cards
	for _, hand := range b.hands {
		b.dealer.Discard(hand.Cards...)
	}
	b.hands = nil

	b.out <- "Play again? " + utils.Bold("Yes (y)") + " or " + utils.Bold("no (n)")

	for line := range b.in {
		lowerLine := strings.ToLower(line)
		switch lowerLine {
		case "yes":
			fallthrough
		case "y":
			b.Play()
		case "no":
			fallthrough
		case "n":
			b.quit()
			return
		default:
			b.out <- utils.Yellow(fmt.Sprintf("Unknown command: %s", line))
			b.out <- "Play again?" + utils.Bold("Yes (y)") + "or " + utils.Bold("no (n)")
		}
	}
}

// hit allows the player to draw another card.
// Returns true if the hand busts
func (b *blackjack) hit(role entities.Role) (bool, error) {
	if b.hands == nil {
		return false, errors.New("No hand!")
	}

	hand := b.hands[role]
	hand.Cards = append(hand.Cards, b.dealer.Draw())
	b.hands[role] = hand

	return b.sumHand(hand) > 21, nil
}

func (b blackjack) printUpdate() {
	utils.Clear(b.out)

	b.out <- utils.Dim(fmt.Sprintf("Your chips: %d", b.userChips))
	b.out <- utils.Dim(fmt.Sprintf("Wagered: %d", b.wager))
	b.out <- utils.Divider()

	dealerHand := b.hands[entities.DealerRole]
	userHand := b.hands[entities.UserRole]

	dealterShowing := b.sumHand(dealerHand)
	userShowing := b.sumHand(userHand)

	if dealerHand.HasHidden() {
		b.out <- utils.Bold("Dealer") + fmt.Sprintf("\t(showing %d)", dealterShowing)
	} else {
		b.out <- utils.Bold("Dealer") + fmt.Sprintf("\t(total %d)", dealterShowing)
	}
	for _, line := range utils.RenderHand(dealerHand) {
		b.out <- line
	}

	b.out <- ""
	b.out <- utils.Bold("You") + fmt.Sprintf("\t(total %d)", userShowing)
	for _, line := range utils.RenderHand(userHand) {
		b.out <- line
	}

	b.out <- utils.Divider()
}

func (b blackjack) sumHand(hand entities.Hand) int {
	total := 0
	nAces := 0
	for _, card := range hand.Cards {
		if card.Hidden {
			continue
		}
		if card.Rank == entities.Ace {
			nAces++
		}
		total += card.GetValue()
	}

	for total > 21 && nAces > 0 {
		nAces--
		total -= 10
	}

	return total
}
