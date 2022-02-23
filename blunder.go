package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
)

type AnalysisMove struct {
	Color    chess.Color
	Actual   int
	Applied  *chess.Position
	Best     int
	BestMove *chess.Move
	Position *chess.Position
	Status   chess.Method
}

func analyze(eng *uci.Engine, pos *chess.Position, ifErr AnalysisMove) AnalysisMove {
	cmdPos := uci.CmdPosition{Position: pos}
	cmdGo := uci.CmdGo{MoveTime: time.Second / 1000}
	if err := eng.Run(cmdPos, cmdGo); err != nil {
		return ifErr
	}
	return AnalysisMove{
		Color:    pos.Turn(),
		Position: pos,
		BestMove: eng.SearchResults().BestMove,
	}
}

func Blunders(pgn_s string) []string {
	pgn, err := chess.PGN(strings.NewReader(pgn_s))
	if err != nil {
		fmt.Println("Error! %+v", err)
		return []string{}
	}
	game := chess.NewGame(pgn)
	fmt.Printf("%v\n", game.GetTagPair("Date"))
	eng, err := uci.New("stockfish")
	if err != nil {
		panic(err)
	}
	defer eng.Close()
	// initialize UCI with new game
	if err := eng.Run(uci.CmdUCI, uci.CmdIsReady, uci.CmdUCINewGame); err != nil {
		panic(err)
	}
	blunders := []string{}
	analyzed := []AnalysisMove{}
	errorMove := AnalysisMove{}

	/*
		   Walk through list of all positions.

		   1) Capture a new struct with each one with a subset of info that we care about:
			    a) Who made the move (black or white)
			    b) The position
		        c) What the engine thinks is the best move
		   2) Score the position if the best move was taken
		   3) If we're not on the first position, score the previous actual move
	*/
	for index, actualMove := range game.Positions() {
		// 1)
		aMove := analyze(eng, actualMove, errorMove)

		// 2)
		if aMove != errorMove {
			bestMoveApplied := actualMove.Update(aMove.BestMove)
			bestAnalysis := analyze(eng, bestMoveApplied, errorMove)
			aMove.Best = eng.SearchResults().Info.Score.CP
			aMove.Status = bestMoveApplied.Status()
			aMove.Applied = bestAnalysis.Position
		} else {
			fmt.Printf("  !! Unable to perform analysis from position %s\n", actualMove)
		}
		analyzed = append(analyzed, aMove)

		// 3)
		if index == 0 {
			continue
		}
		fmt.Printf("  Index: %d, Len: %d Status: %s\n", index, len(analyzed), aMove.Position.Status())
		analyzed[index-1].Actual = eng.SearchResults().Info.Score.CP
	}
	return blunders
}
