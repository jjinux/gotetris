/*
Package main contains a console-based implementation of Tetris.
See the README for more details.

I'm using MVC to structure the application. However, I haven't yet
needed to make actual types for the view or the controller; I can
get away with just a simple function for each.

I don't have any tests. It's just a simple video game ;)
*/

package main

import (
	"time"

	"github.com/nsf/termbox-go"
)

// Function main initializes termbox, renders the view, and starts
// handling events.
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

	g := NewGame()
	render(g)

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
				case ev.Key == termbox.KeySpace:
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
			render(g)
			time.Sleep(animationSpeed)
		}
	}
}
