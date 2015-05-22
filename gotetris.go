/*
Package main contains a console-based implementation of Tetris.
See the README for more details.

I don't have any tests or that much in the way of documentation. It's
just a simple video game ;)
*/

package main

import (
	"fmt"
	//// 	"math/rand"
	"strings"
	"time"

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
const boardEndX = boardStartX + boardWidth
const boardEndY = boardStartY + boardHeight
const instructionsStartX = boardEndX + defaultMarginWidth
const instructionsStartY = boardStartY

// Speeds
const animationSpeed = 10 * time.Millisecond
const slowestSpeed = 700 * time.Millisecond
const fastestSpeed = 60 * time.Millisecond

// Text in the UI
const title = "Tetris Written in Go"

var instructions = []string{
	"Goal: Fill in 5 lines!",
	"",
	"\u2190      Left",
	"\u2192      Right",
	"\u2191      Rotate",
	"\u2193      Drop faster",
	"s      Start",
	"p      Pause",
	"esc    Exit",
	"",
	"Level: %v",
	"Lines: %v",
	"",
	"GAME OVER!"
}

// Game play
const numSquares = 4
const numTypes = 7
const defaultLevel = 1
const maxLevel = 10
const rowsPerLevel = 5

type Game struct {
	curLevel     int
	curX         int
	curY         int
	curPiece     int
	skyline      int
	gameStarted  bool
	gamePaused   bool
	gameOver     bool
	numLines     int
	board        [][]int // [y][x]
	xToErase     []int
	yToErase     []int
	dx           []int
	dy           []int
	dxPrime      []int
	dyPrime      []int
	dxBank       [][]int
	dyBank       [][]int
	fallingTimer *time.Timer
}

// NewGame returns a fully-initialized game.
func NewGame() *Game {
	g := new(Game)
	g.resetGame()
	return g
}

// Reset the game in order to play again.
func (g *Game) resetGame() {
	g.curLevel = 1
	g.curX = 1
	g.curY = 1
	g.skyline = boardHeight - 1
	g.gameStarted = false
	g.gamePaused = false
	g.gameOver = false
	g.numLines = 0

	g.board = make([][]int, boardHeight)
	for y := 0; y < boardHeight; y++ {
		g.board[y] = make([]int, boardWidth)
		for x := 0; x < boardWidth; x++ {
			g.board[y][x] = 0
		}
	}

	g.xToErase = []int{0, 0, 0, 0}
	g.yToErase = []int{0, 0, 0, 0}
	g.dx = []int{0, 0, 0, 0}
	g.dy = []int{0, 0, 0, 0}
	g.dxPrime = []int{0, 0, 0, 0}
	g.dyPrime = []int{0, 0, 0, 0}

	g.dxBank = [][]int{
		{},
		{0, 1, -1, 0},
		{0, 1, -1, -1},
		{0, 1, -1, 1},
		{0, -1, 1, 0},
		{0, 1, -1, 0},
		{0, 1, -1, -2},
		{0, 1, 1, 0},
	}

	g.dyBank = [][]int{
		{},
		{0, 0, 0, 1},
		{0, 0, 0, 1},
		{0, 0, 0, 1},
		{0, 0, 1, 1},
		{0, 0, 1, 1},
		{0, 0, 0, 0},
		{0, 0, 1, 1},
	}

	g.fallingTimer = time.NewTimer(time.Duration(1000000 * time.Second))
	g.fallingTimer.Stop()
}

// Function run initializes termbox, draws everything, and starts handling events.
func (g *Game) Run() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	g.resetGame()
	g.drawBoard()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				switch ev.Key {
				case termbox.KeyEsc:
					return
				case termbox.Key('s'):
					g.start()
				case termbox.Key('p'):
					g.pause()
				case termbox.KeyArrowLeft:
					g.moveLeft()
				case termbox.KeyArrowRight:
					g.moveRight()
				case termbox.KeyArrowUp:
					g.rotate()
				case termbox.KeyArrowDown:
					g.moveDown()
				}
			}
		case <-g.fallingTimer.C:
			g.play()
		default:
			g.drawBoard()
			time.Sleep(animationSpeed)
		}
	}
}

// Function speed calculates the speed based on the curLevel.
func (g *Game) speed() time.Duration {
	return slowestSpeed - fastestSpeed*time.Duration(g.curLevel)
}

func (g *Game) drawBoard() {
	termbox.Clear(backgroundColor, backgroundColor)
	tbprint(titleStartX, titleStartY, instructionsColor, backgroundColor, title)
	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {
			cellColor := pieceColors[g.board[y][x]]
			termbox.SetCell(boardStartX+x, boardStartY+y, ' ', cellColor, cellColor)
		}
	}
	for y, instruction := range instructions {
		if strings.HasPrefix(instruction, "Level:") {
			instruction = fmt.Sprintf(instruction, g.curLevel)
		} else if strings.HasPrefix(instruction, "Lines:") {
			instruction = fmt.Sprintf(instruction, g.numLines)
		} else if strings.HasPrefix(instruction, "GAME OVER") && !g.gameOver {
			instruction = ""
		}
		tbprint(instructionsStartX, instructionsStartY+y, instructionsColor, backgroundColor, instruction)
	}
	termbox.Flush()
}

func (g *Game) play() {
	//// 	if g.moveDown() {
	//// 		g.fallingTimer.Reset(g.speed())
	//// 	} else {
	//// 		g.fillMatrix()
	//// 		g.removeLines()
	//// 		if g.skyline > 0 && g.getPiece() {
	//// 			g.fallingTimer.Reset(g.speed())
	//// 		} else {
	//// 			g.gameOver = true
	//// 		}
	//// 	}
}

//// func (g *Game) fillMatrix() {
//// 	for k := 0; k < numSquares; k++ {
//// 		x := g.curX + g.dx[k]
//// 		y := g.curY + g.dy[k]
//// 		if 0 <= y && y < g.boardHeight && 0 <= x && x < g.boardWidth {
//// 			g.board[y][x] = g.curPiece
//// 			if y < g.skyline {
//// 				g.skyline = y
//// 			}
//// 		}
//// 	}
//// }
////
//// func (g *Game) removeLines() {
//// 	for y := 0; y < g.boardHeight; y++ {
//// 		gapFound := false
//// 		for x := 0; x < g.boardWidth; x++ {
//// 			if g.board[y][x] == 0 {
//// 				gapFound = true
//// 				break
//// 			}
//// 		}
//// 		if !gapFound {
//// 			for k := y; k >= g.skyline; k-- {
//// 				for x = 0; x < g.boardWidth; x++ {
//// 					g.board[k][x] = g.board[k - 1][x]
//// 					ImageElement img = query("#s-$k-$x")
//// 					img.src = pieceColors[g.board[k][x]].src
//// 				}
//// 			}
//// 			for x := 0; x < g.boardWidth; x++ {
//// 				g.board[0][x] = 0
//// 				ImageElement img = query("#s-0-$x")
//// 				img.src = pieceColors[0].src
//// 			}
//// 			g.numLines++
//// 			g.skyline++
//// 			InputElement g.numLinesField = query("#num-lines")
//// 			g.numLinesField.value = g.numLines.toString()
//// 			if g.numLines % rowsPerLevel == 0 && g.curLevel < maxLevel {
//// 				g.curLevel++
//// 			}
//// 			SelectElement levelSelect = query("#level-select")
//// 			levelSelect.selectedIndex = g.curLevel - 1
//// 		}
//// 	}
//// }

func (g *Game) pieceFits(x, y int) bool {
	//// 	for k := 0; k < numSquares; k++ {
	//// 		theX := x + g.dxPrime[k]
	//// 		theY := y + g.dyPrime[k]
	//// 		if theX < 0 || theX >= g.boardWidth || theY >= g.boardHeight {
	//// 			return false
	//// 		}
	//// 		if theY > -1 && g.board[theY][theX] > 0 {
	//// 			return false
	//// 		}
	//// 	}
	return true
}

func (g *Game) erasePiece() {
//// 	for k := 0; k < numSquares; k++ {
//// 		x := g.curX + g.dx[k]
//// 		y := g.curY + g.dy[k]
//// 		if 0 <= y && y < g.boardHeight && 0 <= x && x < g.boardWidth {
//// 			g.xToErase[k] = x
//// 			g.yToErase[k] = y
//// 			g.board[y][x] = 0
//// 		}
//// 	}
}

func (g *Game) drawPiece() {
	//// 	for k := 0; k < numSquares; k++ {
	//// 		x = g.curX + g.dx[k]
	//// 		y = g.curY + g.dy[k]
	//// 		if 0 <= y && y < g.boardHeight && 0 <= x && x < g.boardWidth && g.board[y][x] != -g.curPiece {
	//// 			ImageElement img = query("#s-$y-$x")
	//// 			img.src = pieceColors[g.curPiece].src
	//// 			g.board[y][x] = -g.curPiece
	//// 		}
	//// 		x := g.xToErase[k]
	//// 		y := g.yToErase[k]
	//// 		if g.board[y][x] == 0 {
	//// 			ImageElement img = query("#s-$y-$x")
	//// 			img.src = pieceColors[0].src
	//// 		}
	//// 	}
}

// The user pressed the 's' key to start the game.
func (g *Game) start() {
	if g.gameOver {
		g.resetGame()
	}
	if g.gameStarted {
		if g.gamePaused {
			g.resume()
		}
		return
	}
	g.getPiece()
	g.drawPiece()
	g.gameStarted = true
	g.gamePaused = false
	g.fallingTimer.Reset(g.speed())
}

func (g *Game) pause() {
	if g.gameStarted {
		if g.gamePaused {
			g.resume()
			return
		}
		g.fallingTimer.Stop()
		g.gamePaused = true
	}
}

func (g *Game) moveLeft() {
	if !g.gameStarted || g.gamePaused {
		return
	}
	for k := 0; k < numSquares; k++ {
		g.dxPrime[k] = g.dx[k]
		g.dyPrime[k] = g.dy[k]
	}
	if g.pieceFits(g.curX - 1, g.curY) {
		g.erasePiece()
		g.curX--
		g.drawPiece()
	}
}

func (g *Game) moveRight() {
	//// 	if !g.gameStarted || g.gamePaused {
	////		return
	////	}
	//// 	for k := 0; k < numSquares; k++ {
	//// 		g.dxPrime[k] = g.dx[k]
	//// 		g.dyPrime[k] = g.dy[k]
	//// 	}
	//// 	if g.pieceFits(g.curX + 1, g.curY) {
	//// 		g.erasePiece()
	//// 		g.curX++
	//// 		g.drawPiece()
	//// 	}
}

func (g *Game) rotate() {
	//// 	if !g.gameStarted || g.gamePaused {
	////		return
	////	}
	//// 	for k := 0; k < numSquares; k++ {
	//// 		g.dxPrime[k] = g.dy[k]
	//// 		g.dyPrime[k] = -g.dx[k]
	//// 	}
	//// 	if g.pieceFits(g.curX, g.curY) {
	//// 		g.erasePiece()
	//// 		for k = 0; k < numSquares; k++ {
	//// 			g.dx[k] = g.dxPrime[k]
	//// 			g.dy[k] = g.dyPrime[k]
	//// 		}
	//// 		g.drawPiece()
	//// 	}
}

func (g *Game) moveDown() {
	//// 	if !g.gameStarted || g.gamePaused {
	////		return
	////	}
	//// 	for k := 0; k < numSquares; k++ {
	//// 		g.dxPrime[k] = g.dx[k]
	//// 		g.dyPrime[k] = g.dy[k]
	//// 	}
	//// 	if g.pieceFits(g.curX, g.curY + 1) {
	//// 		g.erasePiece()
	//// 		g.curY++
	//// 		g.drawPiece()
	//// 		return true
	//// 	}
	//// 	return false
}

func (g *Game) getPiece() bool {
	//// 	g.curPiece = 1 + rand.Int()%numTypes
	//// 	g.curX = 5
	//// 	g.curY = 0
	//// 	for k := 0; k < numSquares; k++ {
	//// 		g.dx[k] = g.dxBank[g.curPiece][k]
	//// 		g.dy[k] = g.dyBank[g.curPiece][k]
	//// 	}
	//// 	for k = 0; k < numSquares; k++ {
	//// 		g.dxPrime[k] = g.dx[k]
	//// 		g.dyPrime[k] = g.dy[k]
	//// 	}
	//// 	if g.pieceFits(g.curX, g.curY) {
	//// 		g.drawPiece()
	//// 		return true
	//// 	}
	return false
}

func (g *Game) resume() {
	//// 	if g.gameStarted && g.gamePaused {
	//// 		g.play()
	//// 		g.gamePaused = false
	//// 	}
}

//// func (g *Game) fall() {
//// 	for k := 0; k < numSquares; k++ {
//// 		g.dxPrime[k] = g.dx[k]
//// 		g.dyPrime[k] = g.dy[k]
//// 	}
//// 	if !g.pieceFits(g.curX, g.curY + 1) {
//// 		return
//// 	}
//// 	g.fallingTimer.Stop()
//// 	g.erasePiece()
//// 	while g.pieceFits(g.curX, g.curY + 1) {
//// 		g.curY++
//// 	}
//// 	g.drawPiece()
//// 	g.fallingTimer.Reset(g.speed())
//// }
////
//// func (g *Game) write(String message) {
//// 	p = new(Element.tag("p"))
//// 	p.text = message
//// 	document.body.nodes.add(p)
//// }

// Function tbprint draws a string.
func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

// Function main runs a new Game.
func main() {
	NewGame().Run()
}
