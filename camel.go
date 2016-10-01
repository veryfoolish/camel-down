// Cameling Down

package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Tiles are probably either -1 or 0 or 1
type tile int

// We have five different types of camels, 'Z' is the blank camel
type camel rune

// A stack has zero to five different camels, this can be extended at some point but I don't want to do that right now
type stack [5]camel

// The rollList has some number of camels that have yet to be rolled for this leg of the race.
type rollList []camel

// A space can have either a tile or a stack, or both very temporarily I guess, depending on how resolution works
type space struct {
	tile
	stack
}

// A board is some number of spaces, I think it's 16 maybe, but that's probably 18 if we include the start and end
type board []space

var whiteCamel camel = 'W'
var yellowCamel camel = 'Y'
var orangeCamel camel = 'O'
var greenCamel camel = 'G'
var blueCamel camel = 'B'
var blankCamel camel = 'Z'
var initialStack stack = [5]camel{whiteCamel, yellowCamel, orangeCamel, greenCamel, blueCamel}
var blankStack stack = [5]camel{blankCamel, blankCamel, blankCamel, blankCamel, blankCamel}

	
// initializeBoard sets the board so no spaces have tiles or camels on them, but sets the camels at the start
func initializeBoard(boardSize int) board {
	var aBoard board
	var blankSpace space

	// Initializing a blank space
	blankSpace.tile = 0
	blankSpace.stack = blankStack

	// Now this works to fill up a new board.
	for i := 0; i < boardSize; i++ {
		aBoard = append(aBoard, blankSpace)
	}

	// Let's set up our camels off the board
	aBoard[0].stack = initialStack

	return aBoard
}

// initializeRollList resets the roll list to include all starting camels.
func initializeRollList() (r rollList) {
	for i := 0; i < 5; i++ {
		r = append(r, initialStack[i])
	}
	return r
}

// setTile is to flip a tile to oasis or desert. There's no check to ensure it's -1 or 0 or 1. Also tiles aren't tied to players so we can just do as many as we want.
func (b board) setTile(tileValue tile, spaceNumber int) {
	b[spaceNumber].tile = tileValue
}

// getTileState is to get the state of a tile. Should be -1 or 0 or 1.
func (b board) getTileState(spaceNumber int) tile {
	return b[spaceNumber].tile
}

// getCamel will return the camel in a specific position in a rollList
func (r rollList) getCamel(position int) camel {
	return r[position]
}

// rollMe will roll the die for a specific camel and return the value rolled. Optionally, we could weight the dice...
func (c camel) rollMe() int {

	rSeed := rand.NewSource(time.Now().UnixNano())
	rGen := rand.New(rSeed)
	
	return rGen.Intn(3) + 1
}

// rollDice will roll one die for each camel and then move the camel stacks appropriately. Maybe we should have a camel struct that keeps track of positions.
func rollDice(startR rollList) (leftR rollList) {

	rSeed := rand.NewSource(time.Now().UnixNano())
	rGen := rand.New(rSeed)

	// Need to select a camel at random.
	camelNum := rGen.Intn(len(startR))
	selectedCamel := startR.getCamel(camelNum)

	fmt.Println(fmt.Sprintf("%c", selectedCamel))
	
	// Then need to roll that camel's die.
	rollValue := selectedCamel.rollMe()
	fmt.Println(rollValue)
	
	// Then need to resolve the camel positions?

	// Then need to remove the camel from the list of camels yet to roll. This code still isn't 100% clear to me, so I should reread it a few times when more awake.
	leftR = append(startR[:camelNum], startR[camelNum+1:]...)
	fmt.Print("leftR = ")
	fmt.Println(leftR)
	fmt.Println(len(leftR))
	
	// Then if there are still camels left to roll, need to do the whole thing over...probably can do all this recursively...
	return leftR

}

func main() {

	theBoard := initializeBoard(18)
	rollList := initializeRollList()



	// Everything below is to scratch code for validating functionality...is this "testing?" You be the judge...

	i := len(rollList)
	fmt.Print("i = ")
	fmt.Println(i)
	
	for i > 0 {	
		rollList = rollDice(rollList)
		i = len(rollList)
		fmt.Print("i = ")
		fmt.Println(i)
	}
	
	fmt.Println(rollList)
	fmt.Println(theBoard[5])
	fmt.Println(theBoard.getTileState(5))

	theBoard.setTile(1, 5)

	fmt.Println(theBoard[5])
	fmt.Println(theBoard.getTileState(5))
}
