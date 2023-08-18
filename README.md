# gotetris: Tetris Written in Go

This is a console-based version of Tetris written in Go.

![Screen shot](https://raw.githubusercontent.com/jjinux/gotetris/master/screen_shot.png)

## Run directly

	go run github.com/jjinux/gotetris@latest

## Install

	go get github.com/jjinux/gotetris

## Working on the Code

See [How to Write Go Code](https://golang.org/doc/code.html).

Setup your Go workspace:

	mkdir -p gotetris/src/github.com/jjinux
	cd gotetris
	(cd src/github.com/jjinux &&
	  git clone https://github.com/jjinux/gotetris.git)
	echo "See src/github.com/jjinux/gotetris/README." > README

	# Do this each time to work on the code from the top-level
	# gotetris directory.
	export GOPATH=`pwd`
	export PATH=$PATH:$GOPATH/bin

	# Install dependencies.
	go get -u github.com/nsf/termbox-go

Build:

	go install github.com/jjinux/gotetris

Execute:

	gotetris

## Credits

gotetris is built on top of the excellent
[termbox-go](https://github.com/nsf/termbox-go) library.

It was inspired by [Alexei Kourbatov](http://www.javascripter.net) as
well as my port of [Tetris to Dart](http://code.google.com/p/tetris-in-dart/).
