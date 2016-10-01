// Camel-Down is an engine for investigating a popular children's board game.

package main

import (
	"fmt"
	"math/rand"
	"time"
)

// We're going to be rolling some dice.
func init() {
	rand.Seed(time.Now().UnixNano())
}

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
	for i := 0; i < len(initialStack); i++ {
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

	return rand.Intn(3) + 1

}

// findSpace takes a particular camel and returns the space number on the board that that camel is currently on
func (g game) findCamel(c camel) (s int, p int) {
	for s := 0; s < len(g.gameBoard); s++ {
		for p := 0; p < len(g.gameBoard[s].stack); p++ {
			if g.gameBoard[s].stack[p] == c {
				return s, p
			}
		}
	}
	return 0, 0
}

// moveStack moves a stack of camels a certain number of spaces on a board, for the selected camel and all camels above it
func (g game) moveStack(c camel, moveValue int) {

	// Find the space the camel is on and position in the stack for that camel
	cSpace, cPos := g.findCamel(c)

	fartherPosition := cSpace + moveValue

	// Sign check, relevant for desert tiles
	nearerPosition := cSpace
	if moveValue < 0 {
		fartherPosition = cSpace
		nearerPosition = cSpace + moveValue
	}

	// We build up the destination space's stack. The one starting nearer the start will wind up on top of the one starting farther. Then we put it on the right space.
	var newStack stack

	for i := 0; i < len(g.gameBoard[nearerPosition].stack); i++ {
		if i <= cPos && moveValue > 0 || moveValue < 0 {
			newStack = append(newStack, g.gameBoard[nearerPosition].stack[i])
		}
	}

	for j := 0; j < len(g.gameBoard[fartherPosition].stack); j++ {
		if j >= cPos && moveValue < 0 || moveValue > 0 {
			newStack = append(newStack, g.gameBoard[fartherPosition].stack[j])
		}
	}

	g.gameBoard[cSpace+moveValue].stack = newStack

	if len(g.gameBoard[cSpace].stack[cPos:]) > 0 {
		g.gameBoard[cSpace].stack = g.gameBoard[cSpace].stack[cPos+1:]
	} else {
		g.gameBoard[cSpace].stack = blankStack
	}

	// Lastly, check if we're going to be on a desert or oasis tile. If so, move the stack to that tile and act as if we are moving from the desert/oasis tile to the appropriate one next to it
	if g.gameBoard[fartherPosition].tile != 0 {
		g.moveStack(c, int(g.gameBoard[fartherPosition].tile))
	}

}

// rollDice will roll one die for each camel and then move the camel stacks appropriately. Maybe we should have a camel struct that keeps track of positions.
func (g game) rollDice() (leftR rollList) {

	prettyPrintGame(g)

	startR := g.gameRollList

	// Need to select a camel at random.
	camelNum := rand.Intn(len(startR))
	selectedCamel := startR.getCamel(camelNum)

	// Then need to roll that camel's die.
	rollValue := selectedCamel.rollMe()
	fmt.Print(fmt.Sprintf("%c", selectedCamel))
	fmt.Print(" rolled a ")
	fmt.Println(rollValue)

	// Then need to resolve the camel positions...
	// I suppose the first thing to do is figure out which space the camel is on, and then move that  space's stack by that many spaces.
	g.moveStack(selectedCamel, rollValue)

	// Then need to remove the camel from the list of camels yet to roll. This code still isn't 100% clear to me, so I should reread it a few times when more awake.
	leftR = append(startR[:camelNum], startR[camelNum+1:]...)
	fmt.Print("Camels left to roll = ")
	fmt.Println(len(leftR))

	// Then if there are still camels left to roll, need to do the whole thing over...probably can do all this recursively...
	return leftR

}

// prettyPrintGame will display the current status of the game board but in a pretty format.
func prettyPrintGame(g game) {
	fmt.Println()
	for i := 0; i < len(g.gameBoard); i++ {
		if len(g.gameBoard[i].stack) > 0 || g.gameBoard[i].tile != 0 {
			fmt.Print(i)
			fmt.Print(": ")
			fmt.Print(g.gameBoard[i].tile)
			fmt.Print(": ")
			for j := 0; j < len(g.gameBoard[i].stack); j++ {
				fmt.Print(fmt.Sprintf("%c", g.gameBoard[i].stack[j]))
			}
			fmt.Println("")
		}
	}
	fmt.Println()
}

func main() {

	theBoard := initializeBoard(18)
	theRollList := initializeRollList()

	var theGame game
	theGame.gameBoard = theBoard
	theGame.gameRollList = theRollList

	fmt.Println("")
	fmt.Println("")
	fmt.Println("Welcome to Camel Down. Camels are stacked below one another from right to left. Rightmost is bottommost.")
	fmt.Println("")
	fmt.Println("The initial number is the space number, the second is the current tile status (desert or oasis). Then the camels on that space.")
	fmt.Println("If no camels are on a space, and the space doesn't have a desert or oasis tile, I deem it 'uninteresting' and do not display it.")
	fmt.Println("We're going to start a few legs here, let's see how it goes.")
	for leg := 1; leg <= 3; leg++ {
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
		fmt.Print("Starting leg number ")
		fmt.Println(leg)
		fmt.Println("")
		fmt.Println("")
		i := len(theGame.gameRollList)
		for i > 0 {
			theGame.gameRollList = theGame.rollDice()
			i = len(theGame.gameRollList)
		}
		if len(theGame.gameBoard[4].stack) == 0 {
			theGame.gameBoard[4].tile = 1
		}
		if len(theGame.gameBoard[6].stack) == 0 {
			theGame.gameBoard[6].tile = -1
		}
		theGame.gameRollList = initializeRollList()
	}

	// Everything below is to scratch code for validating functionality...is this "testing?" You be the judge...

	prettyPrintGame(theGame)

}
