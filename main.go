package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"casino/games"
	"casino/games/blackjack"
	"casino/games/poker"
	"casino/utils"
)

func main() {
	// out is the channel to write _to_ the user
	out := make(chan string, 32)

	// inPipe is the pipe from the main process to the game process
	inPipe := make(chan string) // we'll close this on quit to stop the game
	console := utils.NewConsole(out)
	defer console.Close()

	// Run console in background, but wait for it explicitly.
	runDone := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		runDone <- console.Run(ctx)
	}()

	saveManager := utils.NewInMemorySaveDataManager()
	dealer := utils.NewDealer()
	b := blackjack.NewBlackjack(dealer, saveManager, inPipe, out, cancel)
	p := poker.NewPoker(dealer, saveManager, inPipe, out, cancel)
	gameMap := map[int]games.Game{
		1: b,
		2: p,
	}

	utils.Clear(out)
	bannerTop, bannerMiddle, bannerBottom := utils.Banner("THE CASINO")
	out <- bannerTop
	out <- bannerMiddle
	out <- bannerBottom
	out <- utils.Dim("Type 'quit' at any time to leave")
	out <- "Select a number from the menu below to play:"
	for id, g := range gameMap {
		out <- fmt.Sprintf("%d. %s", id, g.Name())
	}

	var choice int
	for input := range console.Out {
		i64, err := strconv.ParseInt(strings.TrimSpace(input), 10, 64)
		if err != nil {
			out <- utils.Red(fmt.Sprintf("%q is not a valid number", input))
			continue
		}
		i := int(i64)
		if _, ok := gameMap[i]; ok {
			choice = i
			break
		}
		out <- utils.Yellow(fmt.Sprintf("Unknown option: %s", input))
	}

	game := gameMap[choice]
	go game.Play()

	// Forward subsequent console lines to the game until "quit".
	go func() {
		defer close(inPipe) // let blackjack.Play exit its range
		for line := range console.Out {
			if strings.EqualFold(strings.TrimSpace(line), "quit") {
				console.Close() // cancels Console.Run; closes console.Out
				return
			}
			inPipe <- line
		}
	}()

	// Block until Console.Run returns (e.g., after Close or EOF).
	if err := <-runDone; err != nil {
		fmt.Println(utils.Dim("Thanks for playing!"))
	}
}
