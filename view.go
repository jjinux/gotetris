package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/nsf/termbox-go"
)

// Colors
const backgroundColor = termbox.ColorBlue
const boardColor = termbox.ColorBlack
const instructionsColor = termbox.ColorYellow

var pieceColors = []termbox.Attribute{
	termbox.ColorBlack,
	termbox.ColorRed,
	termbox.ColorGreen,
	termbox.ColorYellow,
	termbox.ColorBlue,
	termbox.ColorMagenta,
	termbox.ColorCyan,
	termbox.ColorWhite,
}

// Layout
const defaultMarginWidth = 2
const defaultMarginHeight = 1
const titleStartX = defaultMarginWidth
const titleStartY = defaultMarginHeight
const titleHeight = 1
const titleEndY = titleStartY + titleHeight
const boardStartX = defaultMarginWidth
const boardStartY = titleEndY + defaultMarginHeight
const boardWidth = 10
const boardHeight = 16
const boardEndX = boardStartX + boardWidth*2
const boardEndY = boardStartY + boardHeight
const instructionsStartX = boardEndX + defaultMarginWidth
const instructionsStartY = boardStartY

// Text in the UI
const title = "TETRIS WRITTEN IN GO"

var instructions = []string{
	"Goal: Fill in 5 lines!",
	"",
	"\u2190      Left",
	"\u2192      Right",
	"\u2191      Rotate",
	"\u2193      Down",
	"Space  Fall",
	"s      Start",
	"p      Pause",
	"esc    Exit",
	"",
	"Level: %v",
	"Lines: %v",
	"",
	"GAME OVER!",
}

// This takes care of rendering everything.
func render(g *Game) {
	termbox.Clear(backgroundColor, backgroundColor)
	tbprint(titleStartX, titleStartY, instructionsColor, backgroundColor, title)
	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {
			cellValue := g.board[y][x]
			absCellValue := int(math.Abs(float64(cellValue)))
			cellColor := pieceColors[absCellValue]
			termbox.SetCell(boardStartX+2*x, boardStartY+y, ' ', cellColor, cellColor)
			termbox.SetCell(boardStartX+2*x+1, boardStartY+y, ' ', cellColor, cellColor)
		}
	}
	for y, instruction := range instructions {
		if strings.HasPrefix(instruction, "Level:") {
			instruction = fmt.Sprintf(instruction, g.level)
		} else if strings.HasPrefix(instruction, "Lines:") {
			instruction = fmt.Sprintf(instruction, g.numLines)
		} else if strings.HasPrefix(instruction, "GAME OVER") && g.state != gameOver {
			instruction = ""
		}
		tbprint(instructionsStartX, instructionsStartY+y, instructionsColor, backgroundColor, instruction)
	}
	termbox.Flush()
}

// Function tbprint draws a string.
func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}
