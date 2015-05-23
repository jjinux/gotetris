/*
Package main contains a console-based implementation of Tetris.
See the README for more details.

I don't have any tests. It's just a simple video game ;)
*/

package main

import (
	"fmt"
	"math/rand"
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

// Game play
const numSquares = 4
const numTypes = 7
const defaultLevel = 1
const maxLevel = 10
const rowsPerLevel = 5

// Struct Game contains all the game state.
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
				switch {
				case ev.Key == termbox.KeyArrowLeft:
					g.moveLeft()
				case ev.Key == termbox.KeyArrowRight:
					g.moveRight()
				case ev.Key == termbox.KeyArrowUp:
					g.rotate()
				case ev.Key == termbox.KeyArrowDown:
					g.moveDown()
				case ev.Ch == ' ':
					g.fall()
				case ev.Ch == 's':
					g.start()
				case ev.Ch == 'p':
					g.pause()
				case ev.Key == termbox.KeyEsc:
					return
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

// Set the timer to make the pieces fall again.
func (g *Game) resetFallingTimer() {
	g.fallingTimer.Reset(g.speed())
}

// Function speed calculates the speed based on the curLevel.
func (g *Game) speed() time.Duration {
	return slowestSpeed - fastestSpeed*time.Duration(g.curLevel)
}

// This takes care of drawing everything.
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

// This gets called everytime g.fallingTimer goes off.
func (g *Game) play() {
	if g.moveDown() {
		g.resetFallingTimer()
	} else {
		g.fillMatrix()
		g.removeLines()
		if g.skyline > 0 && g.getPiece() {
			g.resetFallingTimer()
		} else {
			g.gameOver = true
		}
	}
}

// This gets called as part of the piece falling.
func (g *Game) fillMatrix() {
	for k := 0; k < numSquares; k++ {
		x := g.curX + g.dx[k]
		y := g.curY + g.dy[k]
		if 0 <= y && y < boardHeight && 0 <= x && x < boardWidth {
			g.board[y][x] = g.curPiece
			if y < g.skyline {
				g.skyline = y
			}
		}
	}
}

// Look for completed lines and remove them.
func (g *Game) removeLines() {
	for y := 0; y < boardHeight; y++ {
		gapFound := false
		for x := 0; x < boardWidth; x++ {
			if g.board[y][x] == 0 {
				gapFound = true
				break
			}
		}
		if !gapFound {
			for k := y; k >= g.skyline; k-- {
				for x := 0; x < boardWidth; x++ {
					g.board[k][x] = g.board[k-1][x]
				}
			}
			for x := 0; x < boardWidth; x++ {
				g.board[0][x] = 0
			}
			g.numLines++
			g.skyline++
			if g.numLines%rowsPerLevel == 0 && g.curLevel < maxLevel {
				g.curLevel++
			}
		}
	}
}

// Return whether or not a piece fits.
func (g *Game) pieceFits(x, y int) bool {
	for k := 0; k < numSquares; k++ {
		theX := x + g.dxPrime[k]
		theY := y + g.dyPrime[k]
		if theX < 0 || theX >= boardWidth || theY >= boardHeight {
			return false
		}
		if theY > -1 && g.board[theY][theX] > 0 {
			return false
		}
	}
	return true
}

// This gets called when a piece moves to a new location.
func (g *Game) erasePiece() {
	for k := 0; k < numSquares; k++ {
		x := g.curX + g.dx[k]
		y := g.curY + g.dy[k]
		if 0 <= y && y < boardHeight && 0 <= x && x < boardWidth {
			g.board[y][x] = 0
		}
	}
}

// Place the piece in the board.
func (g *Game) placePiece() {
	for k := 0; k < numSquares; k++ {
		x := g.curX + g.dx[k]
		y := g.curY + g.dy[k]
		if 0 <= y && y < boardHeight && 0 <= x && x < boardWidth && g.board[y][x] != -g.curPiece {
			g.board[y][x] = -g.curPiece
		}
	}
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
	g.placePiece()
	g.gameStarted = true
	g.gamePaused = false
	g.resetFallingTimer()
}

// The user pressed the 'p' key to pause the game.
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

// The user pressed the left arrow.
func (g *Game) moveLeft() {
	if !g.gameStarted || g.gamePaused || g.gameOver {
		return
	}
	for k := 0; k < numSquares; k++ {
		g.dxPrime[k] = g.dx[k]
		g.dyPrime[k] = g.dy[k]
	}
	if g.pieceFits(g.curX-1, g.curY) {
		g.erasePiece()
		g.curX--
		g.placePiece()
	}
}

// The user pressed the right arrow.
func (g *Game) moveRight() {
	if !g.gameStarted || g.gamePaused || g.gameOver {
		return
	}
	for k := 0; k < numSquares; k++ {
		g.dxPrime[k] = g.dx[k]
		g.dyPrime[k] = g.dy[k]
	}
	if g.pieceFits(g.curX+1, g.curY) {
		g.erasePiece()
		g.curX++
		g.placePiece()
	}
}

// The user pressed the up arrow in order to rotate the piece.
func (g *Game) rotate() {
	if !g.gameStarted || g.gamePaused || g.gameOver {
		return
	}
	for k := 0; k < numSquares; k++ {
		g.dxPrime[k] = g.dy[k]
		g.dyPrime[k] = -g.dx[k]
	}
	if g.pieceFits(g.curX, g.curY) {
		g.erasePiece()
		for k := 0; k < numSquares; k++ {
			g.dx[k] = g.dxPrime[k]
			g.dy[k] = g.dyPrime[k]
		}
		g.placePiece()
	}
}

// Move the piece downward if possible.
func (g *Game) moveDown() bool {
	if !g.gameStarted || g.gamePaused || g.gameOver {
		return false
	}
	for k := 0; k < numSquares; k++ {
		g.dxPrime[k] = g.dx[k]
		g.dyPrime[k] = g.dy[k]
	}
	if !g.pieceFits(g.curX, g.curY+1) {
		return false
	}
	g.erasePiece()
	g.curY++
	g.placePiece()
	return true
}

// The user pressed the space bar to make the piece fall.
func (g *Game) fall() {
	if !g.gameStarted || g.gamePaused || g.gameOver {
		return
	}
	for k := 0; k < numSquares; k++ {
		g.dxPrime[k] = g.dx[k]
		g.dyPrime[k] = g.dy[k]
	}
	if !g.pieceFits(g.curX, g.curY+1) {
		return
	}
	g.fallingTimer.Stop()
	g.erasePiece()
	for g.pieceFits(g.curX, g.curY+1) {
		g.curY++
	}
	g.placePiece()
	g.resetFallingTimer()
}

// Get a random piece and try to place it.
func (g *Game) getPiece() bool {
	g.curPiece = 1 + rand.Int()%numTypes
	g.curX = boardWidth / 2
	g.curY = 0
	for k := 0; k < numSquares; k++ {
		g.dx[k] = g.dxBank[g.curPiece][k]
		g.dy[k] = g.dyBank[g.curPiece][k]
	}
	for k := 0; k < numSquares; k++ {
		g.dxPrime[k] = g.dx[k]
		g.dyPrime[k] = g.dy[k]
	}
	if !g.pieceFits(g.curX, g.curY) {
		return false
	}
	g.placePiece()
	return true
}

// Resume after pausing.
func (g *Game) resume() {
	if g.gameStarted && g.gamePaused && !g.gameOver {
		g.gamePaused = false
		g.play()
	}
}

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
