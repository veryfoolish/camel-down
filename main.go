package main

import (
	"fmt"
	"github.com/ThePlightOfOthers/camel-down/camel"
)

func main() {

	var theGame camel.Game

	theGame.GameBoard = theGame.InitializeBoard(18)
	theGame.GameRollList = theGame.InitializeRollList()

	fmt.Println("\n\nWelcome to Camel Down. Camels are Stacked below one another from right to left. Rightmost is bottommost.\n")
	fmt.Println("The initial number is the space number, the second is the current Tile status (desert or oasis). Then the camels on that space.")
	fmt.Println("If no camels are on a space, and the space doesn't have a desert or oasis Tile, I deem it 'uninteresting' and do not display it.")
	fmt.Println("We're going to start a few legs here, let's see how it goes.")
	for leg := 1; leg <= 3; leg++ {
		fmt.Printf("\n\n\nStarting leg number %d\n\n\n", leg)
		i := len(theGame.GameRollList)
		for i > 0 {
			theGame.GameRollList = theGame.RollDice()
			i = len(theGame.GameRollList)
		}
		if len(theGame.GameBoard[4].Stack) == 0 {
			theGame.GameBoard[4].Tile = -1
		}
		if len(theGame.GameBoard[6].Stack) == 0 {
			theGame.GameBoard[6].Tile = -1
		}
		theGame.GameRollList = theGame.InitializeRollList()
	}

	// Everything below is to scratch code for validating functionality...is this "testing?" You be the judge...

	theGame.PrettyPrintGame()

}
