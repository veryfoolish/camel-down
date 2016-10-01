// Camel-Down is an engine for investigating a popular children's board game.
package camel

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
type Tile int

// We have five different types of camels in the typical game, but trying to be general now (should have branched probably, will learn git soon enough).
//No blank camel required in this version, I think?
type camel rune

// A Stack has zero to five different camels in a typical game, trying to make this flexible.
type Stack []camel

// The rollList has some number of camels that have yet to be rolled for this leg of the race.
type rollList []camel

// A space can have either a Tile or a Stack, or both very temporarily I guess, depending on how resolution works
type space struct {
	Tile
	Stack
}

// A board is some number of spaces, I think it's 16 maybe, but that's probably 18 if we include the start and end.
type board []space

// A game has a board and a rollList.
type Game struct {
	GameBoard    board
	GameRollList rollList
}

var whiteCamel camel = 'W'
var yellowCamel camel = 'Y'
var orangeCamel camel = 'O'
var greenCamel camel = 'G'
var blueCamel camel = 'B'

// initialStack should probably be a parameter similar to how the size of the game board is, but then our camels won't be colourful
var initialStack Stack = []camel{whiteCamel, yellowCamel, orangeCamel, greenCamel, blueCamel}
var blankStack Stack

// initializeBoard sets the board so no spaces have Tiles or camels on them, but sets the camels at the start
func (g Game) InitializeBoard(boardSize int) board {
	var aBoard board
	var blankSpace space

	// Initializing a blank space
	blankSpace.Stack = blankStack

	// Now this works to fill up a new board.
	for i := 0; i < boardSize; i++ {
		aBoard = append(aBoard, blankSpace)
	}

	// Let's set up our camels off the board
	aBoard[0].Stack = initialStack

	return aBoard
}

// initializeRollList resets the roll list to include all starting camels.
func (g Game) InitializeRollList() (r rollList) {
	for i := 0; i < len(initialStack); i++ {
		r = append(r, initialStack[i])
	}
	return r
}

// setTile is to flip a Tile to oasis or desert. There's no check to ensure it's -1 or 0 or 1. Also Tiles aren't tied to players so we can just do as many as we want.
func (b board) setTile(TileValue Tile, spaceNumber int) {
	b[spaceNumber].Tile = TileValue
}

// getTileState is to get the state of a Tile. Should be -1 or 0 or 1.
func (b board) getTileState(spaceNumber int) Tile {
	return b[spaceNumber].Tile
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
func (g Game) findCamel(c camel) (s int, p int) {
	for s := 0; s < len(g.GameBoard); s++ {
		for p := 0; p < len(g.GameBoard[s].Stack); p++ {
			if g.GameBoard[s].Stack[p] == c {
				return s, p
			}
		}
	}
	return 0, 0
}

// winner takes a camel and says who won the game. This panic hack is probably something to refactor away soon.
func (g Game) winner(c camel) {
	cSpace, _ := g.findCamel(c)
	theWinner := g.GameBoard[cSpace].Stack[0]
	fmt.Print(fmt.Sprintf("\n\n%c won the game! Suck it, everyone else.\n\n", theWinner))
	panic("Game over!")
}

// moveStack moves a Stack of camels a certain number of spaces on a board, for the selected camel and all camels above it
func (g Game) moveStack(c camel, moveValue int) {

	// Find the space the camel is on and position in the Stack for that camel
	cSpace, cPos := g.findCamel(c)

	fartherPosition := cSpace + moveValue

	// See if we won the game - this is a terrible implementation will want to fix this immediately.
	if fartherPosition >= len(g.GameBoard) {
		g.winner(c)
	}

	// Sign check, relevant for desert Tiles
	nearerPosition := cSpace
	if moveValue < 0 {
		fartherPosition = cSpace
		nearerPosition = cSpace + moveValue
	}

	// We build up the destination space's Stack. The one starting nearer the start will wind up on top of the one starting farther. Then we put it on the right space.
	var newStack Stack

	for i := 0; i < len(g.GameBoard[nearerPosition].Stack); i++ {
		if i <= cPos && moveValue > 0 || moveValue < 0 {
			newStack = append(newStack, g.GameBoard[nearerPosition].Stack[i])
		}
	}

	for j := 0; j < len(g.GameBoard[fartherPosition].Stack); j++ {
		if j >= cPos && moveValue < 0 || moveValue > 0 {
			newStack = append(newStack, g.GameBoard[fartherPosition].Stack[j])
		}
	}

	g.GameBoard[cSpace+moveValue].Stack = newStack

	if len(g.GameBoard[cSpace].Stack[cPos:]) > 0 {
		g.GameBoard[cSpace].Stack = g.GameBoard[cSpace].Stack[cPos+1:]
	} else {
		g.GameBoard[cSpace].Stack = blankStack
	}

	// Lastly, check if we're going to be on a desert or oasis Tile. If so, move the Stack to that Tile and act as if we are moving from the desert/oasis Tile to the appropriate one next to it
	if g.GameBoard[fartherPosition].Tile != 0 {
		g.moveStack(c, int(g.GameBoard[fartherPosition].Tile))
	}

}

// rollDice will roll one die for each camel and then move the camel Stacks appropriately. Maybe we should have a camel struct that keeps track of positions.
func (g Game) RollDice() (leftR rollList) {

	g.PrettyPrintGame()

	startR := g.GameRollList

	// Need to select a camel at random.
	camelNum := rand.Intn(len(startR))
	selectedCamel := startR.getCamel(camelNum)

	// Then need to roll that camel's die.
	rollValue := selectedCamel.rollMe()
	fmt.Print(fmt.Sprintf("%c", selectedCamel))
	fmt.Print(" rolled a ")
	fmt.Println(rollValue)

	// Then need to resolve the camel positions...
	// I suppose the first thing to do is figure out which space the camel is on, and then move that  space's Stack by that many spaces.
	g.moveStack(selectedCamel, rollValue)

	// Then need to remove the camel from the list of camels yet to roll. This code still isn't 100% clear to me, so I should reread it a few times when more awake.
	leftR = append(startR[:camelNum], startR[camelNum+1:]...)
	fmt.Print("Camels left to roll = ")
	fmt.Println(len(leftR))

	// Then if there are still camels left to roll, need to do the whole thing over...probably can do all this recursively...
	return leftR

}

// PrettyPrintGame will display the current status of the game board but in a pretty format.
func (g Game) PrettyPrintGame() {
	fmt.Println()
	for i := 0; i < len(g.GameBoard); i++ {
		if len(g.GameBoard[i].Stack) > 0 || g.GameBoard[i].Tile != 0 {
			fmt.Print(i)
			fmt.Print(": ")
			fmt.Print(g.GameBoard[i].Tile)
			fmt.Print(": ")
			for j := 0; j < len(g.GameBoard[i].Stack); j++ {
				fmt.Print(fmt.Sprintf("%c", g.GameBoard[i].Stack[j]))
			}
			fmt.Println("")
		}
	}
	fmt.Println()
}
