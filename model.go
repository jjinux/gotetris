package main

import (
	"math/rand"
	"time"
)

// Speeds
const slowestSpeed = 700 * time.Millisecond
const fastestSpeed = 60 * time.Millisecond

// Game play
const numSquares = 4
const numTypes = 7
const defaultLevel = 1
const maxLevel = 10
const rowsPerLevel = 5

// Pieces
var dxBank = [][]int{
	{},
	{0, 1, -1, 0},
	{0, 1, -1, -1},
	{0, 1, -1, 1},
	{0, -1, 1, 0},
	{0, 1, -1, 0},
	{0, 1, -1, -2},
	{0, 1, 1, 0},
}

var dyBank = [][]int{
	{},
	{0, 0, 0, 1},
	{0, 0, 0, 1},
	{0, 0, 0, 1},
	{0, 0, 1, 1},
	{0, 0, 1, 1},
	{0, 0, 0, 0},
	{0, 0, 1, 1},
}

type gameState int

const (
	gameIntro gameState = iota
	gameStarted
	gamePaused
	gameOver
)

// Struct Game contains all the game state.
type Game struct {
	board        [][]int // [y][x]
	state        gameState
	level        int
	numLines     int
	piece        int
	x            int
	y            int
	dx           []int
	dy           []int
	dxPrime      []int
	dyPrime      []int
	skyline      int
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
	g.board = make([][]int, boardHeight)
	for y := 0; y < boardHeight; y++ {
		g.board[y] = make([]int, boardWidth)
		for x := 0; x < boardWidth; x++ {
			g.board[y][x] = 0
		}
	}

	g.state = gameIntro
	g.level = 1
	g.numLines = 0
	g.x = 1
	g.y = 1
	g.dx = []int{0, 0, 0, 0}
	g.dy = []int{0, 0, 0, 0}
	g.dxPrime = []int{0, 0, 0, 0}
	g.dyPrime = []int{0, 0, 0, 0}
	g.skyline = boardHeight - 1

	g.fallingTimer = time.NewTimer(time.Duration(1000000 * time.Second))
	g.fallingTimer.Stop()
}

// Set the timer to make the pieces fall again.
func (g *Game) resetFallingTimer() {
	g.fallingTimer.Reset(g.speed())
}

// Function speed calculates the speed based on the level.
func (g *Game) speed() time.Duration {
	return slowestSpeed - fastestSpeed*time.Duration(g.level)
}

// This gets called everytime g.fallingTimer goes off.
func (g *Game) play() {
	if g.moveDown() {
		g.resetFallingTimer()
	} else {
		g.lockPiece()
	}
}

func (g *Game) lockPiece() {
	g.fillMatrix()
	g.removeLines()
	if g.skyline > 0 && g.getPiece() {
		g.resetFallingTimer()
	} else {
		g.state = gameOver
	}
}

// This gets called as part of the piece falling.
func (g *Game) fillMatrix() {
	for k := 0; k < numSquares; k++ {
		x := g.x + g.dx[k]
		y := g.y + g.dy[k]
		if 0 <= y && y < boardHeight && 0 <= x && x < boardWidth {
			g.board[y][x] = g.piece
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
			if g.numLines%rowsPerLevel == 0 && g.level < maxLevel {
				g.level++
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
		x := g.x + g.dx[k]
		y := g.y + g.dy[k]
		if 0 <= y && y < boardHeight && 0 <= x && x < boardWidth {
			g.board[y][x] = 0
		}
	}
}

// Place the piece in the board.
func (g *Game) placePiece() {
	for k := 0; k < numSquares; k++ {
		x := g.x + g.dx[k]
		y := g.y + g.dy[k]
		if 0 <= y && y < boardHeight && 0 <= x && x < boardWidth && g.board[y][x] != -g.piece {
			g.board[y][x] = -g.piece
		}
	}
}

// The user pressed the 's' key to start the game.
func (g *Game) start() {
	switch g.state {
	case gameStarted:
		return
	case gamePaused:
		g.resume()
		return
	case gameOver:
		g.resetGame()
		fallthrough
	default:
		g.state = gameStarted
		g.getPiece()
		g.placePiece()
		g.resetFallingTimer()
	}
}

// The user pressed the 'p' key to pause the game.
func (g *Game) pause() {
	switch g.state {
	case gameStarted:
		g.state = gamePaused
		g.fallingTimer.Stop()
	case gamePaused:
		g.resume()
	}
}

// The user pressed the left arrow.
func (g *Game) moveLeft() {
	if g.state != gameStarted {
		return
	}
	for k := 0; k < numSquares; k++ {
		g.dxPrime[k] = g.dx[k]
		g.dyPrime[k] = g.dy[k]
	}
	if g.pieceFits(g.x-1, g.y) {
		g.erasePiece()
		g.x--
		g.placePiece()
	}
}

// The user pressed the right arrow.
func (g *Game) moveRight() {
	if g.state != gameStarted {
		return
	}
	for k := 0; k < numSquares; k++ {
		g.dxPrime[k] = g.dx[k]
		g.dyPrime[k] = g.dy[k]
	}
	if g.pieceFits(g.x+1, g.y) {
		g.erasePiece()
		g.x++
		g.placePiece()
	}
}

// The user pressed the up arrow in order to rotate the piece.
func (g *Game) rotate() {
	if g.state != gameStarted {
		return
	}
	for k := 0; k < numSquares; k++ {
		g.dxPrime[k] = g.dy[k]
		g.dyPrime[k] = -g.dx[k]
	}
	if g.pieceFits(g.x, g.y) {
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
	if g.state != gameStarted {
		return false
	}
	for k := 0; k < numSquares; k++ {
		g.dxPrime[k] = g.dx[k]
		g.dyPrime[k] = g.dy[k]
	}
	if !g.pieceFits(g.x, g.y+1) {
		return false
	}
	g.erasePiece()
	g.y++
	g.placePiece()
	return true
}

// The user pressed the space bar to make the piece fall.
func (g *Game) fall() {
	if g.state != gameStarted {
		return
	}
	for k := 0; k < numSquares; k++ {
		g.dxPrime[k] = g.dx[k]
		g.dyPrime[k] = g.dy[k]
	}
	if !g.pieceFits(g.x, g.y+1) {
		return
	}
	g.fallingTimer.Stop()
	g.erasePiece()
	for g.pieceFits(g.x, g.y+1) {
		g.y++
	}
	g.placePiece()
	g.resetFallingTimer()
	g.lockPiece()
}

// Get a random piece and try to place it.
func (g *Game) getPiece() bool {
	g.piece = 1 + rand.Int()%numTypes
	g.x = boardWidth / 2
	g.y = 0
	for k := 0; k < numSquares; k++ {
		g.dx[k] = dxBank[g.piece][k]
		g.dy[k] = dyBank[g.piece][k]
	}
	for k := 0; k < numSquares; k++ {
		g.dxPrime[k] = g.dx[k]
		g.dyPrime[k] = g.dy[k]
	}
	if !g.pieceFits(g.x, g.y) {
		return false
	}
	g.placePiece()
	return true
}

// Resume after pausing.
func (g *Game) resume() {
	g.state = gameStarted
	g.play()
}
