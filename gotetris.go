/*
Package main contains a console-based implementation of Tetris.
See the README for more details.

I don't have any tests or that much in the way of documentation. It's
just a simple video game ;)
*/

package main

import (
	"time"

	"github.com/nsf/termbox-go"
)

const backgroundColor = termbox.ColorBlack
const instructionsColor = termbox.ColorWhite
const defaultMarginWidth = 2
const defaultMarginHeight = 1
const boardStartX = defaultMarginWidth
const boardStartY = defaultMarginHeight
const boardWidth = 10
const boardHeight = 16
const boardEndX = boardStartX + boardWidth
const boardEndY = boardStartY + boardHeight
const instructionsStartX = boardEndX + defaultMarginWidth
const instructionsStartY = defaultMarginHeight

var instructions = []string{
	"Use arrow keys",
	"Press down to make the piece fall",
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func draw() {
	termbox.Clear(backgroundColor, backgroundColor)
	for y := boardStartY; y < boardEndY; y++ {
		for x := boardStartX; x < boardEndX; x++ {
			termbox.SetCell(x, y, ' ', termbox.ColorGreen, termbox.ColorGreen)
		}
	}
	for i, instruction := range instructions {
		tbprint(instructionsStartX, instructionsStartY+i, instructionsColor, backgroundColor, instruction)
	}
	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	draw()

loop:
	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
				break loop
			}
		default:
			draw()
			time.Sleep(10 * time.Millisecond)
		}
	}
}
