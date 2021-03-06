package board

import (
	"testing"

	"github.com/0xhexnumbers/partysim/mp1"
)

func TestStartMovement(t *testing.T) {
	g := *mp1.InitializeGame(ES, mp1.GameConfig{MaxTurns: 20})
	g.NextEvent.Handle(1, &g)

	SpaceIs(esStartingSpace, 0, g, "", t)
}

func TestWarpCDeterminization(t *testing.T) {
	g := *mp1.InitializeGame(ES, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(3, 3)
	g.NextEvent.Handle(1, &g)

	gGate1 := g
	gGate1.NextEvent.Handle(esEntrance7, &gGate1)
	g1bd := gGate1.Board.Data.(esBoardData)
	IntIs(1, g1bd.Gate, "G1 Gate", t)
	SpaceIs(mp1.NewChainSpace(10, 1), 0, gGate1, "G1", t)

	gGate2or3 := g
	gGate2or3.NextEvent.Handle(esEntrance1, &gGate2or3)
	g23bd := gGate2or3.Board.Data.(esBoardData)
	if !g23bd.Gate2or3 {
		t.Error("Gate2or3 not set")
	}
	SpaceIs(esStartingSpace, 0, gGate2or3, "G2o3", t)

	gGate2or3.Players[1].CurrentSpace = mp1.NewChainSpace(4, 4)
	gGate2or3.NextEvent.Handle(1, &gGate2or3)
	gGate2or3.NextEvent.Handle(esEntrance6, &gGate2or3)
	g23bd = gGate2or3.Board.Data.(esBoardData)
	IntIs(3, g23bd.Gate, "G2o3 Gate", t)
	SpaceIs(mp1.NewChainSpace(9, 1), 1, gGate2or3, "G2o3", t)
}

func TestWarpC(t *testing.T) {
	g := *mp1.InitializeGame(ES, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(3, 3)

	g1 := g
	g1.Board.Data = esBoardData{Gate: 1}
	g1.NextEvent.Handle(1, &g1)
	SpaceIs(mp1.NewChainSpace(10, 1), 0, g1, "G1", t)

	g2 := g
	g2.Board.Data = esBoardData{Gate: 2}
	g2.NextEvent.Handle(1, &g2)
	SpaceIs(esStartingSpace, 0, g2, "G2", t)

	g3 := g
	g3.Board.Data = esBoardData{Gate: 3}
	g3.NextEvent.Handle(1, &g3)
	SpaceIs(esStartingSpace, 0, g3, "G3", t)
}

func TestWarpD(t *testing.T) {
	g := *mp1.InitializeGame(ES, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(4, 4)

	g1 := g
	g1.Board.Data = esBoardData{Gate: 1}
	g1.NextEvent.Handle(1, &g1)
	SpaceIs(esStartingSpace, 0, g1, "G1", t)

	g2 := g
	g2.Board.Data = esBoardData{Gate: 2}
	g2.NextEvent.Handle(1, &g2)
	SpaceIs(mp1.NewChainSpace(10, 1), 0, g2, "G2", t)

	g3 := g
	g3.Board.Data = esBoardData{Gate: 3}
	g3.NextEvent.Handle(1, &g3)
	SpaceIs(mp1.NewChainSpace(9, 1), 0, g3, "G3", t)
}

func TestWarpE(t *testing.T) {
	g := *mp1.InitializeGame(ES, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(5, 4)

	g1 := g
	g1.Board.Data = esBoardData{Gate: 1}
	g1.NextEvent.Handle(1, &g1)
	SpaceIs(mp1.NewChainSpace(9, 1), 0, g1, "G1", t)

	g2 := g
	g2.Board.Data = esBoardData{Gate: 2}
	g2.NextEvent.Handle(1, &g2)
	SpaceIs(mp1.NewChainSpace(7, 1), 0, g2, "G2", t)

	g3 := g
	g3.Board.Data = esBoardData{Gate: 3}
	g3.NextEvent.Handle(1, &g3)
	SpaceIs(esStartingSpace, 0, g3, "G3", t)
}

func TestWarpF(t *testing.T) {
	g := *mp1.InitializeGame(ES, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(6, 1)

	g1 := g
	g1.Board.Data = esBoardData{Gate: 1}
	g1.NextEvent.Handle(1, &g1)
	SpaceIs(esStartingSpace, 0, g1, "G1", t)

	g2 := g
	g2.Board.Data = esBoardData{Gate: 2}
	g2.NextEvent.Handle(1, &g2)
	SpaceIs(mp1.NewChainSpace(9, 1), 0, g2, "G2", t)

	g3 := g
	g3.Board.Data = esBoardData{Gate: 3}
	g3.NextEvent.Handle(1, &g3)
	SpaceIs(mp1.NewChainSpace(10, 1), 0, g3, "G3", t)
}

func TestWarpG(t *testing.T) {
	g := *mp1.InitializeGame(ES, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(7, 4)

	g1 := g
	g1.Board.Data = esBoardData{Gate: 1}
	g1.NextEvent.Handle(1, &g1)                //Move
	g1.NextEvent.Handle(ESBranchGotoWarp, &g1) //Goto Warp G
	g1.NextEvent.Handle(Gate(2), &g1)          //Set to Gate 2
	SpaceIs(esStartingSpace, 0, g1, "G1", t)

	g2 := g
	g2.Board.Data = esBoardData{Gate: 2}
	g2.NextEvent.Handle(1, &g2)                //Move
	g2.NextEvent.Handle(ESBranchGotoWarp, &g2) //Goto Warp G
	g2.NextEvent.Handle(Gate(3), &g2)          //Set to Gate 3
	SpaceIs(esStartingSpace, 0, g2, "G2", t)

	g3 := g
	g3.Board.Data = esBoardData{Gate: 3}
	g3.NextEvent.Handle(1, &g3)                //Move
	g3.NextEvent.Handle(ESBranchGotoWarp, &g3) //Goto Warp G
	SpaceIs(mp1.NewChainSpace(11, 1), 0, g3, "G3", t)
}

func TestWarpH(t *testing.T) {
	g := *mp1.InitializeGame(ES, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(7, 10)

	g1 := g
	g1.Board.Data = esBoardData{Gate: 1}
	g1.NextEvent.Handle(1, &g1)
	SpaceIs(mp1.NewChainSpace(8, 1), 0, g1, "G1", t)

	g2 := g
	g2.Board.Data = esBoardData{Gate: 2}
	g2.NextEvent.Handle(1, &g2)
	SpaceIs(mp1.NewChainSpace(11, 1), 0, g2, "G2", t)

	g3 := g
	g3.Board.Data = esBoardData{Gate: 3}
	g3.NextEvent.Handle(1, &g3)
	SpaceIs(mp1.NewChainSpace(8, 1), 0, g3, "G3", t)
}

func TestWarpI(t *testing.T) {
	g := *mp1.InitializeGame(ES, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(8, 8)

	g1 := g
	g1.Board.Data = esBoardData{Gate: 1}
	g1.NextEvent.Handle(1, &g1)
	SpaceIs(esStartingSpace, 0, g1, "G1", t)

	g2 := g
	g2.Board.Data = esBoardData{Gate: 2}
	g2.NextEvent.Handle(1, &g2)
	SpaceIs(esStartingSpace, 0, g2, "G2", t)

	g3 := g
	g3.Board.Data = esBoardData{Gate: 3}
	g3.NextEvent.Handle(1, &g3)
	SpaceIs(esStartingSpace, 0, g3, "G3", t)
}

func TestWarpJ(t *testing.T) {
	g := *mp1.InitializeGame(ES, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(9, 9)

	g1 := g
	g1.Board.Data = esBoardData{Gate: 1}
	g1.NextEvent.Handle(1, &g1)                //Move
	g1.NextEvent.Handle(ESBranchGotoWarp, &g1) //Goto warp J
	g1.NextEvent.Handle(Gate(2), &g1)          //Set to gate 2
	SpaceIs(esStartingSpace, 0, g1, "G1", t)

	g2 := g
	g2.Board.Data = esBoardData{Gate: 2}
	g2.NextEvent.Handle(1, &g2)                //Move
	g2.NextEvent.Handle(ESBranchGotoWarp, &g2) //Goto warp J
	g2.NextEvent.Handle(Gate(3), &g2)          //Set to gate 3
	SpaceIs(esStartingSpace, 0, g2, "G2", t)

	g3 := g
	g3.Board.Data = esBoardData{Gate: 3}
	g3.NextEvent.Handle(1, &g3)                //Move
	g3.NextEvent.Handle(ESBranchGotoWarp, &g3) //Goto warp J
	SpaceIs(mp1.NewChainSpace(11, 1), 0, g3, "G3", t)
}

func TestWarpK(t *testing.T) {
	g := *mp1.InitializeGame(ES, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(9, 17)

	g1 := g
	g1.Board.Data = esBoardData{Gate: 1}
	g1.NextEvent.Handle(1, &g1)
	SpaceIs(mp1.NewChainSpace(11, 1), 0, g1, "G1", t)

	g2 := g
	g2.Board.Data = esBoardData{Gate: 2}
	g2.NextEvent.Handle(1, &g2)
	SpaceIs(esStartingSpace, 0, g2, "G2", t)

	g3 := g
	g3.Board.Data = esBoardData{Gate: 3}
	g3.NextEvent.Handle(1, &g3)
	SpaceIs(mp1.NewChainSpace(7, 1), 0, g3, "G3", t)
}

func TestStarSpace(t *testing.T) {
	g := *mp1.InitializeGame(ES, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(4, 1)

	gNoCoins := g
	gNoCoins.NextEvent.Handle(1, &gNoCoins)
	CoinsIs(13, 0, gNoCoins, "NoCoins", t)
	SpaceIs(mp1.NewChainSpace(4, 3), 0, gNoCoins, "NoCoins Space", t)

	gCoins := g
	gCoins.Players[0].Coins = 20
	gCoins.NextEvent.Handle(1, &gCoins)                     //Move
	gCoins.NextEvent.Handle(ESVisitBabyBowserPlay, &gCoins) //Pay 20 coins
	gCoins.NextEvent.Handle(ESBattleBabyBowserWin, &gCoins) //Won dice roll

	StarsIs(1, 0, gCoins, "", t)
	CoinsIs(3, 0, gCoins, "", t)
	SpaceIs(mp1.NewChainSpace(4, 3), 0, gCoins, "Space", t)
}

func TestChanceSpace(t *testing.T) {
	g := *mp1.InitializeGame(ES, mp1.GameConfig{MaxTurns: 20})
	var bd esBoardData
	bd.StarTaken[0] = true
	g.Board.Data = bd
	g.Players[0].CurrentSpace = mp1.NewChainSpace(4, 1)

	g.NextEvent.Handle(1, &g)
	EventIs(mp1.ChanceTime{Player: 0}, g.NextEvent, "", t)
}

func TestStarSpaceReset(t *testing.T) {
	g := *mp1.InitializeGame(ES, mp1.GameConfig{MaxTurns: 20})
	var bd esBoardData
	bd.StarTaken = [7]bool{false, true, true, true, true, true, true}
	g.Board.Data = bd
	g.Players[0].Coins = 20
	g.Players[0].CurrentSpace = mp1.NewChainSpace(4, 1)

	g.NextEvent.Handle(1, &g)                     //Move
	g.NextEvent.Handle(ESVisitBabyBowserPlay, &g) //Pay 20 coins
	g.NextEvent.Handle(ESBattleBabyBowserWin, &g) //Won dice roll

	bd = g.Board.Data.(esBoardData)
	if bd.StarTaken != [7]bool{} {
		t.Errorf("Star Data not reset")
	}
}

func TestSendToStart(t *testing.T) {
	g := *mp1.InitializeGame(ES, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(4, 3)
	g.Players[1].CurrentSpace = mp1.NewChainSpace(4, 3)
	g.Players[2].CurrentSpace = mp1.NewChainSpace(4, 3)
	g.Players[3].CurrentSpace = mp1.NewChainSpace(4, 3)

	g.NextEvent.Handle(1, &g)
	for i := range g.Players {
		SpaceIs(esStartingSpace, i, g, "", t)
	}
	CoinsIs(13, 0, g, "", t)
	IntIs(1, g.Players[0].HappeningCount, "HappeningCount", t)
}

func TestVisitBowser(t *testing.T) {
	g := *mp1.InitializeGame(ES, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(12, 0)
	g.Players[0].Coins = 30

	gCoins := g
	gCoins.NextEvent.Handle(1, &gCoins)       //Move
	gCoins.NextEvent.Handle(Gate(1), &gCoins) //Set to gate 1
	CoinsIs(13, 0, gCoins, "CoinsTaken", t)

	gStars := g
	gStars.Players[0].Stars = 2
	gStars.NextEvent.Handle(1, &gStars)       //Move
	gStars.NextEvent.Handle(Gate(1), &gStars) //Set to gate 1
	CoinsIs(33, 0, gStars, "StarTaken", t)
	StarsIs(1, 0, gStars, "StarTaken", t)
}
