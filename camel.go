// Cameling Down

package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Tiles are probably either -1 or 0 or 1
type tile int

// We have five different types of camels in the typical game, but trying to be general now (should have branched probably, will learn git soon enough).
//No blank camel required in this version, I think?
type camel rune

// A stack has zero to five different camels in a typical game, trying to make this flexible.
type stack []camel

// The rollList has some number of camels that have yet to be rolled for this leg of the race.
type rollList []camel

// A space can have either a tile or a stack, or both very temporarily I guess, depending on how resolution works
type space struct {
	tile
	stack
}

// A board is some number of spaces, I think it's 16 maybe, but that's probably 18 if we include the start and end.
type board []space

// A game has a board and a rollList.
type game struct {
	gameBoard    board
	gameRollList rollList
}

var whiteCamel camel = 'W'
var yellowCamel camel = 'Y'
var orangeCamel camel = 'O'
var greenCamel camel = 'G'
var blueCamel camel = 'B'
var initialStack stack = []camel{whiteCamel, yellowCamel, orangeCamel, greenCamel, blueCamel}
var blankStack stack

// initializeBoard sets the board so no spaces have tiles or camels on them, but sets the camels at the start
func initializeBoard(boardSize int) board {
	var aBoard board
	var blankSpace space

	// Initializing a blank space
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

// findSpace takes a particular camel and returns the space number on the board that that camel is currently on
func (g game) findSpace(c camel) (s int) {
	for s := 0; s < len(g.gameBoard); s++ {
		for i := 0; i < 5; i++ {
			if g.gameBoard[s].stack[i] == c {
				return s
			}
		}
	}
	return 0
}

// moveStack moves a stack of camels a certain number of spaces on a board
func (g game) moveStack(s int, moveValue int) {

	fartherPosition := s + moveValue

	// First, check if we're going to be on a desert or oasis tile. If so, act as if we are moving from the desert/oasis tile to the appropriate one next to it
	if g.gameBoard[fartherPosition].tile != 0 {
		g.gameBoard.moveStack(fartherPosition, g.gameBoard[fartherPosition].tile)
	}

	// Sign check, relevant for desert tiles
	nearerPosition := s
	if moveValue < 0 {
		fartherPosition = s
		nearerPosition = s + moveValue
	}

	// We build up the destination space's stack. The one starting nearer the start will wind up on top of the one starting farther. Then we put it on the right space.
	var newStack stack
	for i := 0; i < len(g.gameBoard[nearerPosition].stack); i++ {
		newStack = append(newStack, g.gameBoard[nearerPosition].stack[i])
	}

	for j := 0; j < len(g.gameBoard[fartherPosition].stack); j++ {
		newStack = append(newStack, g.gameBoard[fartherPosition].stack)
	}

	g.gameBoard[s+moveValue].stack = newStack
	g.gameBoard[s].stack = blankStack

}

// rollDice will roll one die for each camel and then move the camel stacks appropriately. Maybe we should have a camel struct that keeps track of positions.
func (g game) rollDice() (leftR rollList) {

	startR := g.gameRollList

	rSeed := rand.NewSource(time.Now().UnixNano())
	rGen := rand.New(rSeed)

	// Need to select a camel at random.
	camelNum := rGen.Intn(len(startR))
	selectedCamel := startR.getCamel(camelNum)

	fmt.Println(fmt.Sprintf("%c", selectedCamel))

	// Then need to roll that camel's die.
	rollValue := selectedCamel.rollMe()
	fmt.Println(rollValue)

	// Then need to resolve the camel positions...
	// I suppose the first thing to do is figure out which space the camel is on, and then move that  space's stack by that many spaces.
	g.moveStack(g.findSpace(selectedCamel), rollValue)

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
	theRollList := initializeRollList()

	var theGame game
	theGame.gameBoard = theBoard
	theGame.gameRollList = theRollList

	i := len(theGame.gameRollList)
	for i > 0 {
		theGame.gameRollList = theGame.rollDice()
		i = len(theGame.gameRollList)
	}

	// Everything below is to scratch code for validating functionality...is this "testing?" You be the judge...

	fmt.Println(theBoard[5])
	fmt.Println(theBoard.getTileState(5))

	theBoard.setTile(1, 5)

	fmt.Println(theBoard[5])
	fmt.Println(theBoard.getTileState(5))
}
