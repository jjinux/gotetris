/*
Package main contains a console-based implementation of Tetris.
See the README for more details.

I don't have any tests or that much in the way of documentation. It's
just a simple video game ;)
*/

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

const backgroundColor = termbox.ColorBlack
const instructionsColor = termbox.ColorWhite

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
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func draw() {
	termbox.Clear(backgroundColor, backgroundColor)
	tbprint(titleStartX, titleStartY, instructionsColor, backgroundColor, title)
	for y := boardStartY; y < boardEndY; y++ {
		for x := boardStartX; x < boardEndX; x++ {
			termbox.SetCell(x, y, ' ', termbox.ColorGreen, termbox.ColorGreen)
		}
	}
	for i, instruction := range instructions {
		if strings.HasPrefix(instruction, "Level:") {
			instruction = fmt.Sprintf(instruction, 0)
		} else if strings.HasPrefix(instruction, "Lines:") {
			instruction = fmt.Sprintf(instruction, 0)
		}
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

// /**
//  * This is an implementation of Tetris in Dart.
//  *
//  * Warning: This is my first Dart program, and I haven't even read the tutorial
//  * yet!  Nonetheless, it does work, and it's fairly clean.
//  *
//  * This code was inspired by Alexei Kourbatov (http://www.javascripter.net).
//  */

// library game;
// import 'dart:html';
// import 'dart:math';
// import 'dart:async';

// class Game {

//   static final NUM_SQUARES = 4;
//   static final NUM_TYPES = 7;
//   static final NUM_IMAGES = 8;
//   static final DEFAULT_LEVEL = 1;
//   static final MAX_LEVEL = 10;
//   static final ROWS_PER_LEVEL = 5;
//   static final BOARD_HEIGHT = 16;
//   static final BOARD_WIDTH = 10;
//   static final SLOWEST_SPEED = 700;
//   static final FASTEST_SPEED = 60;
//   static final SQUARE_WIDTH = 16;
//   static final SQUARE_HEIGHT = 16;

//   // Keystroke processing
//   static final INITIAL_DELAY = 200;
//   static final REPEAT_DELAY = 20;

//   // These are for Netscape.
//   static final LEFT_NN = ' 52 ';
//   static final RIGHT_NN = ' 54 ';
//   static final UP_NN = ' 56 53 ';
//   static final DOWN_NN = ' 50 ';
//   static final SPACE_NN = ' 32 ';

//   // These are for Internet Explorer.
//   static final LEFT_IE = ' 37 52 100 ';
//   static final RIGHT_IE = ' 39 54 102 ';
//   static final UP_IE = ' 38 56 53 104 101 ';
//   static final DOWN_IE = ' 40 50 98 ';
//   static final SPACE_IE = ' 32 ';

//   num curLevel = DEFAULT_LEVEL;
//   num curX = 1;
//   num curY = 1;
//   num curPiece;
//   num skyline = BOARD_HEIGHT - 1;
//   bool boardDrawn = false;
//   bool gamePaused = false;
//   bool gameStarted = false;
//   bool sayingBye = false;
//   Timer timer;
//   num numLines = 0;
//   num speed = SLOWEST_SPEED - FASTEST_SPEED * DEFAULT_LEVEL;
//   List<ImageElement> squareImages;
//   List<List<num>> board;
//   List<num> xToErase;
//   List<num> yToErase;
//   List<num> dx;
//   List<num> dy;
//   List<num> dxPrime;
//   List<num> dyPrime;
//   List<List<num>> dxBank;
//   List<List<num>> dyBank;
//   Random random = new Random();

//   // Keystroke processing
//   bool isActiveLeft = false;
//   bool isActiveRight = false;
//   bool isActiveUp = false;
//   bool isActiveDown = false;
//   bool isActiveSpace = false;
//   Timer timerLeft;
//   Timer timerRight;
//   Timer timerDown;

//   Game() {
//     squareImages = [];
//     board = [];
//     xToErase = [0, 0, 0, 0];
//     yToErase = [0, 0, 0, 0];
//     dx = [0, 0, 0, 0];
//     dy = [0, 0, 0, 0];
//     dxPrime = [0, 0, 0, 0];
//     dyPrime = [0, 0, 0, 0];
//     dxBank = [[], [0, 1, -1, 0], [0, 1, -1, -1], [0, 1, -1, 1], [0, -1, 1, 0], [0, 1, -1, 0], [0, 1, -1, -2], [0, 1, 1, 0]];
//     dyBank = [[], [0, 0, 0, 1], [0, 0, 0, 1], [0, 0, 0, 1], [0, 0, 1, 1], [0, 0, 1, 1], [0, 0, 0, 0], [0, 0, 1, 1]];

//     for (num i = 0; i < NUM_IMAGES; i++) {
//       ImageElement img = new Element.tag("img");
//       img.src = 'images/s${i}.png';
//       squareImages.add(img);
//     }

//     for (num i = 0; i < BOARD_HEIGHT; i++) {
//       board.add([]);
//       for (num j = 0; j < BOARD_WIDTH; j++) {
//         board[i].add(0);
//       }
//     }

//     document.on.keyDown.add(onKeyDown);
//     document.on.keyUp.add(onKeyUp);

//     SelectElement levelSelect = query("#level-select");
//     levelSelect.on.change.add((Event e) {
//       onLevelSelectChange();
//       levelSelect.blur();
//     });

//     InputElement startButton = query("#start-button");
//     InputElement pauseButton = query("#pause-button");
//     startButton.on.click.add((event) => start());
//     pauseButton.on.click.add((event) => pause());
//   }

//   void run() {
//     drawBoard();
//     resetGame();
//   }

//   void start() {
//     if (sayingBye) {
//       window.history.back();
//       sayingBye = false;
//     }
//     if (gameStarted) {
//       if (!boardDrawn) {
//         return;
//       }
//       if (gamePaused) {
//         resume();
//       }
//       return;
//     }
//     getPiece();
//     drawPiece();
//     gameStarted = true;
//     gamePaused = false;
//     InputElement numLinesField = query("#num-lines");
//     numLinesField.value = numLines.toString();
//     timer = new Timer(speed, (timer) => play());
//   }

//   void drawBoard() {
//     DivElement boardDiv = query("#board-div");
//     PreElement pre = new Element.tag("pre");
//     boardDiv.nodes.add(pre);
//     pre.classes.add("board");
//     for (num i = 0; i < BOARD_HEIGHT; i++) {
//       DivElement div = new Element.tag("div");
//       pre.nodes.add(div);
//       for (num j = 0; j < BOARD_WIDTH; j++) {
//         ImageElement img = new Element.tag("img");
//         div.nodes.add(img);
//         img.id = "s-$i-$j";
//         img.src = "images/s${board[i][j].abs()}.png";
//         img.width = SQUARE_WIDTH;
//         img.height = SQUARE_HEIGHT;
//       }
//       ImageElement rightMargin = new Element.tag("img");
//       div.nodes.add(rightMargin);
//       rightMargin.src = "images/g.png";
//       rightMargin.width = 1;
//       rightMargin.height = SQUARE_HEIGHT;
//     }
//     DivElement trailingDiv = new Element.tag("div");
//     pre.nodes.add(trailingDiv);
//     ImageElement trailingImg = new Element.tag("img");
//     trailingDiv.nodes.add(trailingImg);
//     trailingImg.src = "images/g.png";
//     trailingImg.id = "board-trailing-img";
//     trailingImg.width = BOARD_WIDTH * 16 + 1;
//     trailingImg.height = 1;
//     boardDrawn = true;
//   }

//   void resetGame() {
//     for (num i = 0; i < BOARD_HEIGHT; i++) {
//       for (num j = 0; j < BOARD_WIDTH; j++) {
//         board[i][j] = 0;
//         ImageElement img = query("#s-$i-$j");
//         img.src = 'images/s0.png';
//       }
//     }
//     gameStarted = false;
//     gamePaused = false;
//     numLines = 0;
//     curLevel = 1;
//     skyline = BOARD_HEIGHT - 1;
//     InputElement numLinesField = query("#num-lines");
//     numLinesField.value = numLines.toString();
//     SelectElement levelSelect = query("#level-select");
//     levelSelect.selectedIndex = 0;

//     // I shouldn't have to call this manually, but I do.
//     // See: http://code.google.com/p/dart/issues/detail?id=2325&thanks=2325&ts=1332879888
//     onLevelSelectChange();
//   }

//   void play() {
//     if (moveDown()) {
//       timer = new Timer(speed, (timer) => play());
//     } else {
//       fillMatrix();
//       removeLines();
//       if (skyline > 0 && getPiece()) {
//         timer = new Timer(speed, (timer) => play());
//       } else {
//         isActiveLeft = false;
//         isActiveUp = false;
//         isActiveRight = false;
//         isActiveDown = false;
//         window.alert('Game over!');
//         resetGame();
//       }
//     }
//   }

//   void pause() {
//     if (boardDrawn && gameStarted) {
//       if (gamePaused) {
//         resume();
//         return;
//       }
//       timer.cancel();
//       gamePaused = true;
//     }
//   }

//   void onKeyDown(KeyboardEvent event) {
//     // I'm positive there are more modern ways to do keyboard event handling.
//     String keyNN, keyIE;
//     keyNN = ' ${event.keyCode} ';
//     keyIE = ' ${event.keyCode} ';

//     // Only preventDefault if we can actually handle the keyDown event.  If we
//     // capture all keyDown events, we break things like using ^r to reload the page.

//     if (!gameStarted || !boardDrawn || gamePaused) {
//       return;
//     }

//     if (LEFT_NN.indexOf(keyNN) != -1 || LEFT_IE.indexOf(keyIE) != -1) {
//       if (!isActiveLeft) {
//         isActiveLeft = true;
//         isActiveRight = false;
//         moveLeft();
//         timerLeft = new Timer(INITIAL_DELAY, (timer) => slideLeft());
//       }
//     } else if (RIGHT_NN.indexOf(keyNN) != -1 || RIGHT_IE.indexOf(keyIE) != -1) {
//       if (!isActiveRight) {
//         isActiveRight = true;
//         isActiveLeft = false;
//         moveRight();
//         timerRight = new Timer(INITIAL_DELAY, (timer) => slideRight());
//       }
//     } else if (UP_NN.indexOf(keyNN) != -1 || UP_IE.indexOf(keyIE) != -1) {
//       if (!isActiveUp) {
//         isActiveUp = true;
//         isActiveDown = false;
//         rotate();
//       }
//     } else if (SPACE_NN.indexOf(keyNN) != -1 || SPACE_IE.indexOf(keyIE) != -1) {
//       if (!isActiveSpace) {
//         isActiveSpace = true;
//         isActiveDown = false;
//         fall();
//       }
//     } else if (DOWN_NN.indexOf(keyNN) != -1 || DOWN_IE.indexOf(keyIE) != -1) {
//       if (!isActiveDown) {
//         isActiveDown = true;
//         isActiveUp = false;
//         moveDown();
//         timerDown = new Timer(INITIAL_DELAY, (timer) => slideDown());
//       }
//     } else {
//       return;
//     }
//     event.preventDefault();
//   }

//   void onKeyUp(KeyboardEvent event) {
//     // See comments in onKeyDown.
//     var keyNN, keyIE;
//     keyNN = ' ${event.keyCode} ';
//     keyIE = ' ${event.keyCode} ';
//     if (LEFT_NN.indexOf(keyNN) != -1 || LEFT_IE.indexOf(keyIE) != -1) {
//       isActiveLeft = false;
//       timerLeft.cancel();
//     } else if (RIGHT_NN.indexOf(keyNN) != -1 || RIGHT_IE.indexOf(keyIE) != -1) {
//       isActiveRight = false;
//       timerRight.cancel();
//     } else if (UP_NN.indexOf(keyNN) != -1 || UP_IE.indexOf(keyIE) != -1) {
//       isActiveUp = false;
//     } else if (DOWN_NN.indexOf(keyNN) != -1 || DOWN_IE.indexOf(keyIE) != -1) {
//       isActiveDown = false;
//       timerDown.cancel();
//     } else if (SPACE_NN.indexOf(keyNN) != -1 || SPACE_IE.indexOf(keyIE) != -1) {
//       isActiveSpace = false;
//     } else {
//       return;
//     }
//     event.preventDefault();
//   }

//   void fillMatrix() {
//     num k, x, y;
//     for (k = 0; k < NUM_SQUARES; k++) {
//       x = curX + dx[k];
//       y = curY + dy[k];
//       if (0 <= y && y < BOARD_HEIGHT && 0 <= x && x < BOARD_WIDTH) {
//         board[y][x] = curPiece;
//         if (y < skyline) {
//           skyline = y;
//         }
//       }
//     }
//   }

//   void removeLines() {
//     num i, j, k;
//     bool gapFound;
//     for (i = 0; i < BOARD_HEIGHT; i++) {
//       gapFound = false;
//       for (j = 0; j < BOARD_WIDTH; j++) {
//         if (board[i][j] == 0) {
//           gapFound = true;
//           break;
//         }
//       }
//       if (!gapFound) {
//         for (k = i; k >= skyline; k--) {
//           for (j = 0; j < BOARD_WIDTH; j++) {
//             board[k][j] = board[k - 1][j];
//             ImageElement img = query("#s-$k-$j");
//             img.src = squareImages[board[k][j]].src;
//           }
//         }
//         for (j = 0; j < BOARD_WIDTH; j++) {
//           board[0][j] = 0;
//           ImageElement img = query("#s-0-$j");
//           img.src = squareImages[0].src;
//         }
//         numLines++;
//         skyline++;
//         InputElement numLinesField = query("#num-lines");
//         numLinesField.value = numLines.toString();
//         if (numLines % ROWS_PER_LEVEL == 0 && curLevel < MAX_LEVEL) {
//           curLevel++;
//         }
//         speed = SLOWEST_SPEED - FASTEST_SPEED * curLevel;
//         SelectElement levelSelect = query("#level-select");
//         levelSelect.selectedIndex = curLevel - 1;
//       }
//     }
//   }

//   bool pieceFits(x, y) {
//     num k, theX, theY;
//     for (k = 0; k < NUM_SQUARES; k++) {
//       theX = x + dxPrime[k];
//       theY = y + dyPrime[k];
//       if (theX < 0 || theX >= BOARD_WIDTH || theY >= BOARD_HEIGHT) {
//         return false;
//       }
//       if (theY > -1 && board[theY][theX] > 0) {
//         return false;
//       }
//     }
//     return true;
//   }

//   void erasePiece() {
//     num k, x, y;
//     if (boardDrawn) {
//       for (k = 0; k < NUM_SQUARES; k++) {
//         x = curX + dx[k];
//         y = curY + dy[k];
//         if (0 <= y && y < BOARD_HEIGHT && 0 <= x && x < BOARD_WIDTH) {
//           xToErase[k] = x;
//           yToErase[k] = y;
//           board[y][x] = 0;
//         }
//       }
//     }
//   }

//   void drawPiece() {
//     num k, x, y;
//     if (boardDrawn) {
//       for (k = 0; k < NUM_SQUARES; k++) {
//         x = curX + dx[k];
//         y = curY + dy[k];
//         if (0 <= y && y < BOARD_HEIGHT && 0 <= x && x < BOARD_WIDTH && board[y][x] != -curPiece) {
//           ImageElement img = query("#s-$y-$x");
//           img.src = squareImages[curPiece].src;
//           board[y][x] = -curPiece;
//         }
//         x = xToErase[k];
//         y = yToErase[k];
//         if (board[y][x] == 0) {
//           ImageElement img = query("#s-$y-$x");
//           img.src = squareImages[0].src;
//         }
//       }
//     }
//   }

//   bool moveDown() {
//     num k;
//     for (k = 0; k < NUM_SQUARES; k++) {
//       dxPrime[k] = dx[k];
//       dyPrime[k] = dy[k];
//     }
//     if (pieceFits(curX, curY + 1)) {
//       erasePiece();
//       curY++;
//       drawPiece();
//       return true;
//     }
//     return false;
//   }

//   void moveLeft() {
//     num k;
//     for (k = 0; k < NUM_SQUARES; k++) {
//       dxPrime[k] = dx[k];
//       dyPrime[k] = dy[k];
//     }
//     if (pieceFits(curX - 1, curY)) {
//       erasePiece();
//       curX--;
//       drawPiece();
//     }
//   }

//   void moveRight() {
//     num k;
//     for (k = 0; k < NUM_SQUARES; k++) {
//       dxPrime[k] = dx[k];
//       dyPrime[k] = dy[k];
//     }
//     if (pieceFits(curX + 1, curY)) {
//       erasePiece();
//       curX++;
//       drawPiece();
//     }
//   }

//   bool getPiece() {
//     num k;
//     curPiece = 1 + random.nextInt(NUM_TYPES);
//     curX = 5;
//     curY = 0;
//     for (k = 0; k < NUM_SQUARES; k++) {
//       dx[k] = dxBank[curPiece][k];
//       dy[k] = dyBank[curPiece][k];
//     }
//     for (k = 0; k < NUM_SQUARES; k++) {
//       dxPrime[k] = dx[k];
//       dyPrime[k] = dy[k];
//     }
//     if (pieceFits(curX, curY)) {
//       drawPiece();
//       return true;
//     }
//     return false;
//   }

//   void resume() {
//     if (boardDrawn && gameStarted && gamePaused) {
//       play();
//       gamePaused = false;
//     }
//   }

//   void onLevelSelectChange() {
//     SelectElement levelSelect = query("#level-select");
//     OptionElement selectedOption = levelSelect.options[levelSelect.selectedIndex];
//     curLevel = int.parse(selectedOption.value);
//     speed = SLOWEST_SPEED - FASTEST_SPEED * curLevel;
//   }

//   void rotate() {
//     num k;
//     for (k = 0; k < NUM_SQUARES; k++) {
//       dxPrime[k] = dy[k];
//       dyPrime[k] = -dx[k];
//     }
//     if (pieceFits(curX, curY)) {
//       erasePiece();
//       for (k = 0; k < NUM_SQUARES; k++) {
//         dx[k] = dxPrime[k];
//         dy[k] = dyPrime[k];
//       }
//       drawPiece();
//     }
//   }

//   void fall() {
//     num k;
//     for (k = 0; k < NUM_SQUARES; k++) {
//       dxPrime[k] = dx[k];
//       dyPrime[k] = dy[k];
//     }
//     if (!pieceFits(curX, curY + 1)) {
//       return;
//     }
//     timer.cancel();
//     erasePiece();
//     while (pieceFits(curX, curY + 1)) {
//       curY++;
//     }
//     drawPiece();
//     timer = new Timer(speed, (timer) => play());
//   }

//   void slideLeft() {
//     if (isActiveLeft) {
//       moveLeft();
//       timerLeft = new Timer(REPEAT_DELAY, (timer) => slideLeft());
//     }
//   }

//   void slideRight() {
//     if (isActiveRight) {
//       moveRight();
//       timerRight = new Timer(REPEAT_DELAY, (timer) => slideRight());
//     }
//   }

//   void slideDown() {
//     if (isActiveDown) {
//       moveDown();
//       timerDown = new Timer(REPEAT_DELAY, (timer) => slideDown());
//     }
//   }

//   void write(String message) {
//     ParagraphElement p = new Element.tag('p');
//     p.text = message;
//     document.body.nodes.add(p);
//   }

// }

// void main() {
//   new Game().run();
// }
