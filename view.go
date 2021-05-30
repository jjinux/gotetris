package main

import (
	"fmt"
	"math"
	"strings"
	"github.com/nsf/termbox-go"
)

// Colors
const (
  backgroundColor = termbox.ColorBlue
  boardColor = termbox.ColorBlack
  instructionsColor = termbox.ColorYellow
)
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
const (
  defaultMarginWidth = 2
  defaultMarginHeight = 1
  titleStartX = defaultMarginWidth
  titleStartY = defaultMarginHeight
  titleHeight = 1
  titleEndY = titleStartY + titleHeight
  boardStartX = defaultMarginWidth
  boardStartY = titleEndY + defaultMarginHeight
  boardWidth = 10
  boardHeight = 16
  cellWidth = 2
  boardEndX = boardStartX + boardWidth*cellWidth
  boardEndY = boardStartY + boardHeight
  instructionsStartX = boardEndX + defaultMarginWidth
  instructionsStartY = boardStartY
)
// Text in the UI
const title = "TETRIS WRITTEN IN GO"

var instructions = []string{
	"Goal: Fill in 5 lines!",
	"",
	"left   Left",
	"right  Right",
	"up     Rotate",
	"down   Down",
	"space  Fall",
	"s      Start",
	"p      Pause",
	"esc,q  Exit",
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
			for i := 0; i < cellWidth; i++ {
				termbox.SetCell(boardStartX+cellWidth*x+i, boardStartY+y, ' ', cellColor, cellColor)
			}
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
