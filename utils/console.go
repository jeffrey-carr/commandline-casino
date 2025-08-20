package utils

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/term"
)

type Console struct {
	In        <-chan string // messages to show to the user (from other components)
	Out       chan string   // complete user-entered lines (to other components)
	Prompt    string        // prompt to display (e.g., "> ")
	outWriter io.Writer

	mu        sync.Mutex
	origState *term.State
	cancel    context.CancelFunc
}

// NewConsole creates a Console. Provide a channel you will write messages into (In).
// The Console exposes Out for you to read user-entered lines.
func NewConsole(in <-chan string, opts ...func(*Console)) *Console {
	c := &Console{
		In:        in,
		Out:       make(chan string, 64),
		Prompt:    "> ",
		outWriter: os.Stdout,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// WithWriter lets you redirect output (e.g., to a buffer in tests).
func WithWriter(w io.Writer) func(*Console) {
	return func(c *Console) { c.outWriter = w }
}

// Run starts the interactive loop. It returns when the context is canceled
// or stdin is closed. It restores the terminal state on exit.
func (c *Console) Run(ctx context.Context) (err error) {
	ctx, c.cancel = context.WithCancel(ctx)
	defer close(c.Out)

	// Handle SIGINT/SIGTERM to cleanly restore the terminal.
	sigCh := make(chan os.Signal, 2)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	// Put terminal into raw mode for responsive input/editing.
	if term.IsTerminal(int(os.Stdin.Fd())) {
		state, e := term.MakeRaw(int(os.Stdin.Fd()))
		if e != nil {
			return e
		}
		c.mu.Lock()
		c.origState = state
		c.mu.Unlock()
		defer c.restore()
	}

	// Reader for stdin runes (handles UTF-8 properly).
	r := bufio.NewReader(os.Stdin)
	runeCh := make(chan rune, 16)
	readErr := make(chan error, 1)

	// Blocking rune reader goroutine
	go func() {
		defer close(runeCh)
		for {
			ch, _, e := r.ReadRune()
			if e != nil {
				readErr <- e
				return
			}
			runeCh <- ch
		}
	}()

	// Draw initial prompt.
	c.printPrompt("")

	var line []rune

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case s := <-sigCh:
			// On Ctrl-C (SIGINT), just move to a new line and keep running.
			if s == os.Interrupt {
				line = line[:0]
				c.newline()
				c.printPrompt("")
				continue
			}
			return nil

		case msg, ok := <-c.In:
			if !ok {
				// Inbound channel closed; keep running for user input.
				c.In = nil
				continue
			}
			// Print inbound message on its own line, then redraw current input.
			c.repaintWithMessage(msg, string(line))

		case ch, ok := <-runeCh:
			if !ok {
				// stdin reader ended; exit
				return nil
			}
			switch ch {
			case '\r', '\n': // Enter
				c.newline()
				text := string(line)
				line = line[:0]
				c.Out <- text
				c.printPrompt("")
			case 3: // Ctrl-C (ETX)
				// Clear the line and show a fresh prompt (common shell behavior).
				line = line[:0]
				c.newline()
				c.printPrompt("")
			case 4: // Ctrl-D (EOT)
				// If line is empty, treat as EOF and exit. If not empty, submit the line.
				if len(line) == 0 {
					c.newline()
					return nil
				}
				c.newline()
				text := string(line)
				line = line[:0]
				c.Out <- text
				c.printPrompt("")
			case 127, 8: // Backspace / Ctrl-H
				if len(line) > 0 {
					line = line[:len(line)-1]
					c.redraw(string(line))
				}
			default:
				line = append(line, ch)
				_, _ = fmt.Fprint(c.outWriter, string(ch))
			}

		case e := <-readErr:
			// stdin read error or EOF
			_ = e
			return nil
		}
	}
}

// Close stops the loop (if running) and restores the terminal.
func (c *Console) Close() {
	if c.cancel != nil {
		c.cancel()
	}
	c.restore()
}

// --- helpers ---

func (c *Console) restore() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.origState != nil {
		_ = term.Restore(int(os.Stdin.Fd()), c.origState)
		c.origState = nil
	}
}

func (c *Console) printPrompt(current string) {
	_, _ = fmt.Fprint(c.outWriter, c.Prompt)
	_, _ = fmt.Fprint(c.outWriter, current)
}

func (c *Console) newline() {
	_, _ = fmt.Fprint(c.outWriter, "\r\n")
}

func (c *Console) redraw(current string) {
	// Move to line start, print prompt+buffer, clear to end.
	_, _ = fmt.Fprintf(c.outWriter, "\r%s%s\x1b[K", c.Prompt, current)
}

func (c *Console) repaintWithMessage(msg, current string) {
	// Go to new line, show message, then redraw prompt+current buffer.
	_, _ = fmt.Fprintf(c.outWriter, "\r\x1b[K%s\r\n", msg)
	c.printPrompt(current)
}
