// Cameling Down

package main

import (
	"fmt"
)

// Tiles are probably either -1 or 0 or 1
type tile int

// We have five different types of camels, 'Z' is the blank camel
type camel rune

// A stack has zero to five different camels, this can be extended at some point but I don't want to do that right now
type stack [5]camel

// A space can have either a tile or a stack, or both very temporarily I guess, depending on how resolution works
type space struct {
	tile
	stack
}

// A board is some number of spaces, I think it's 18 maybe, but that might be 20 if we include the start and end
type board []space

// At the beginning of the game, no spaces have tiles or camels on them
func initializeBoard(boardSize int) board {
	var aBoard board
	var blankCamel camel = 'Z'
	var blankStack stack
	var blankSpace space

	// Now we need a blank stack, maybe
	for c := 0; c < 5; c++ {
		blankStack[c] = blankCamel
	}

	// Finally, a blank space
	blankSpace.tile = 0
	blankSpace.stack = blankStack

	// Now this works to fill up a new board.
	for i := 0; i < boardSize; i++ {
		aBoard = append(aBoard, blankSpace)
	}

	return aBoard
}

func (b board) setTile(tileValue tile, spaceNumber int) {
	b[spaceNumber].tile = tileValue	
}

func main() {
	theBoard := initializeBoard(18)
	
	
	// Everything below is to scratch code for validating functionality...is this "testing?" You be the judge...
	fmt.Println(theBoard[5])

	theBoard.setTile(1, 5)
	
	fmt.Println(theBoard[5])
	
}
