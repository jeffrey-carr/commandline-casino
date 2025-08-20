package utils

import (
	"fmt"
	"strings"

	"casino/entities"
)

const (
	AnsiReset  = "\x1b[0m"
	AnsiBold   = "\x1b[1m"
	AnsiDim    = "\x1b[2m"
	AnsiRed    = "\x1b[31m"
	AnsiGreen  = "\x1b[32m"
	AnsiYellow = "\x1b[33m"
	AnsiBlue   = "\x1b[34m"
	AnsiCyan   = "\x1b[36m"
)

// Text coloring

func Bold(s string) string   { return AnsiBold + s + AnsiReset }
func Dim(s string) string    { return AnsiDim + s + AnsiReset }
func Blue(s string) string   { return AnsiBlue + s + AnsiReset }
func Cyan(s string) string   { return AnsiCyan + s + AnsiReset }
func Green(s string) string  { return AnsiGreen + s + AnsiReset }
func Red(s string) string    { return AnsiRed + s + AnsiReset }
func Yellow(s string) string { return AnsiYellow + s + AnsiReset }

// Rendering

func Clear(out chan string) {
	out <- "\x1b[2J\x1b[H"
}

func Divider() string { return strings.Repeat("─", 52) }

func Banner(title string) (top, middle, bottom string) {
	titleDecor := "═" + title + "═"
	left := max(0, (26 - len(StripANSI(titleDecor))))

	top = strings.Repeat("═", 52)
	middle = strings.Repeat(" ", left) + Bold(title)
	bottom = top
	return
}

// renderHand draws cards horizontally using ASCII/Unicode.
func RenderHand(h entities.Hand) []string {
	var rows [][]string
	for _, c := range h.Cards {
		rows = append(rows, RenderCard(c))
	}
	if len(rows) == 0 {
		return []string{"(empty hand)"}
	}

	height := len(rows[0])
	out := make([]string, height)
	for r := range height {
		var parts []string
		for _, cardRows := range rows {
			parts = append(parts, cardRows[r])
		}
		out[r] = strings.Join(parts, " ")
	}
	return out
}

var (
	cardTop = "┌─────────┐"
	cardBot = "└─────────┘"
)

func RenderCard(c entities.Card) []string {
	if c.Hidden {
		return []string{
			cardTop,
			"│░░░░░░░░░│",
			"│░░░░░░░░░│",
			"│░░░░░░░░░│",
			"│░░░░░░░░░│",
			cardBot,
		}
	}

	suitColor := ""
	if c.Suit == entities.Heart || c.Suit == entities.Diamond {
		suitColor = AnsiRed
	}

	// Fit rank to two chars on edges (handles "10")
	left := PadRight(string(c.Rank), 2)
	right := PadLeft(string(c.Rank), 2)
	center := c.Suit

	line2 := fmt.Sprintf("│%s%-2s%s       │", suitColor, left, AnsiReset)
	line3 := fmt.Sprintf("│    %s%s%s    │", suitColor, center, AnsiReset)
	line4 := "│         │"
	line5 := fmt.Sprintf("│       %s%2s%s│", suitColor, right, AnsiReset)

	return []string{
		cardTop,
		line2,
		line3,
		line4,
		line5,
		cardBot,
	}
}

// Helpers

func StripANSI(s string) string {
	// Remove sequences like \x1b[...m
	var out strings.Builder
	skip := false
	for i := 0; i < len(s); i++ {
		if !skip && i+1 < len(s) && s[i] == 0x1b && s[i+1] == '[' {
			skip = true
			continue
		}
		if skip {
			if s[i] == 'm' {
				skip = false
			}
			continue
		}
		out.WriteByte(s[i])
	}
	return out.String()
}
