package utils

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func GetBet(in chan string, out chan string, message string, minChips, maxChips int) int {
	var wager int

	out <- message
	for line := range in {
		// Attempt to convert to int
		i64, err := strconv.ParseInt(strings.TrimSpace(line), 10, 64)
		if err != nil {
			out <- Yellow(Bold("Wager must be a whole number"))
			out <- message
			continue
		}

		wager = int(i64)
		if wager < minChips {
			out <- Yellow(Bold("Wager must be at least 1"))
		} else if wager > maxChips {
			out <- Yellow(Bold(fmt.Sprintf("You cannot wager %d, you only have %d", wager, maxChips)))
		} else {
			break
		}

		out <- message
	}

	return wager
}

// GetInput takes in a list of commands and waits for user to choose one. It returns
// the command that was chosen
func GetInput(in, out chan string, commands []string, message string) string {
	out <- message
	for line := range in {
		lowerLine := strings.TrimSpace(strings.ToLower(line))
		if slices.Contains(commands, lowerLine) {
			return lowerLine
		}
	}

	return ""
}
