package board

import (
	"fmt"

	"github.com/0xhexnumbers/partysim/mp1"
)

type ESBranchResponse int

const (
	ESBranchGotoWarp ESBranchResponse = iota
	ESBranchContinue
)

func (e ESBranchResponse) String() string {
	switch e {
	case ESBranchGotoWarp:
		return "Go to warp"
	case ESBranchContinue:
		return "Continue down main path"
	}
	return ""
}

//ESBranchEvent let's the player decide if they want to take the warp.
type ESBranchEvent struct {
	Player int
	Moves  int
	Warp1  mp1.ChainSpace
	Warp2  mp1.ChainSpace
	Warp3  mp1.ChainSpace
}

func (e ESBranchEvent) Question(g *mp1.Game) string {
	return fmt.Sprintf("Does %s take the warp?",
		g.Players[e.Player].Char)
}

func (e ESBranchEvent) Type() mp1.EventType {
	return mp1.ENUM_EVT_TYPE
}

func (e ESBranchEvent) ControllingPlayer() int {
	return e.Player
}

func (e ESBranchEvent) Responses() []mp1.Response {
	return []mp1.Response{ESBranchGotoWarp, ESBranchContinue}
}

//Handle executes based on r. If r is true, the player's new position is set
//based on the current gate (setting the next event to set the gate if
//the gate is unknown). If r is false, the player continues down their
//current chain.
func (e ESBranchEvent) Handle(r mp1.Response, g *mp1.Game) {
	gotoWarp := r.(ESBranchResponse)
	bd := g.Board.Data.(esBoardData)
	if gotoWarp == ESBranchGotoWarp {
		switch bd.Gate {
		case 0:
			g.NextEvent = ESWarpDest{
				e.Player,
				e.Moves,
				bd.Gate2or3,
				e.Warp1,
				e.Warp2,
				e.Warp3,
			}
			return
		case 1:
			g.Players[e.Player].CurrentSpace = e.Warp1
		case 2:
			g.Players[e.Player].CurrentSpace = e.Warp2
		case 3:
			g.Players[e.Player].CurrentSpace = e.Warp3
		}
		g.MovePlayer(e.Player, e.Moves)
	} else {
		g.MovePlayer(e.Player, e.Moves)
	}
}

type ESVisitBabyBowserResponse int

const (
	ESVisitBabyBowserPlay ESVisitBabyBowserResponse = iota
	ESVisitBabyBowserIgnore
)

func (e ESVisitBabyBowserResponse) String() string {
	switch e {
	case ESVisitBabyBowserPlay:
		return "Pay 20 coins to play minigame"
	case ESVisitBabyBowserIgnore:
		return "Do not play minigame"
	}
	return ""
}

//ESVisitBabyBowser let's the player decide if they want to play baby
//bowser's minigame to win a star.
type ESVisitBabyBowser struct {
	Player int
	Moves  int
	Index  int
}

func (e ESVisitBabyBowser) Question(g *mp1.Game) string {
	return fmt.Sprintf(
		"Does %s pay 20 coins to play Baby Bowser's star minigame?",
		g.Players[e.Player].Char,
	)
}

func (e ESVisitBabyBowser) Type() mp1.EventType {
	return mp1.ENUM_EVT_TYPE
}

func (e ESVisitBabyBowser) ControllingPlayer() int {
	return e.Player
}

func (e ESVisitBabyBowser) Responses() []mp1.Response {
	return []mp1.Response{ESVisitBabyBowserPlay, ESVisitBabyBowserIgnore}
}

//Handle sets the next event to the baby bowser minigame if r is true. If r
//is false, then nothing happens.
func (e ESVisitBabyBowser) Handle(r mp1.Response, g *mp1.Game) {
	battle := r.(ESVisitBabyBowserResponse)
	if battle == ESVisitBabyBowserPlay {
		g.AwardCoins(e.Player, -20, false)
		g.NextEvent = ESBattleBabyBowser{
			e.Player, e.Moves, e.Index,
		}
	} else {
		g.MovePlayer(e.Player, e.Moves)
	}
}

type ESBattleBabyBowserResponse int

const (
	ESBattleBabyBowserWin ESBattleBabyBowserResponse = iota
	ESBattleBabyBowserLose
)

func (e ESBattleBabyBowserResponse) String() string {
	switch e {
	case ESBattleBabyBowserWin:
		return "Baby Bowser loses the minigame"
	case ESBattleBabyBowserLose:
		return "Baby Bowser wins the minigame"
	}
	return ""
}

//ESBattleBabyBowser decides if the player wins the minigame.
type ESBattleBabyBowser struct {
	Player int
	Moves  int
	Index  int
}

func (e ESBattleBabyBowser) Question(g *mp1.Game) string {
	return fmt.Sprintf("Does %s win Baby Bowser's minigame?",
		g.Players[e.Player].Char)
}

func (e ESBattleBabyBowser) Type() mp1.EventType {
	return mp1.ENUM_EVT_TYPE
}

func (e ESBattleBabyBowser) ControllingPlayer() int {
	return mp1.CPU_PLAYER
}

func (e ESBattleBabyBowser) Responses() []mp1.Response {
	return []mp1.Response{ESBattleBabyBowserWin, ESBattleBabyBowserLose}
}

//Handle gives the player a star and sets the baby bowser's StarTaken flag
//to true if r is true. If r is false, a star is taken from the plaeyr.
func (e ESBattleBabyBowser) Handle(r mp1.Response, g *mp1.Game) {
	star := r.(ESBattleBabyBowserResponse)
	bd := g.Board.Data.(esBoardData)
	if star == ESBattleBabyBowserWin {
		g.Players[e.Player].Stars++
		bd.StarTaken[e.Index] = true
		if esAllStarsCollected(bd) {
			bd.StarTaken = [7]bool{
				false, false, false, false, false, false, false,
			}
		}
		g.Board.Data = bd
	} else if g.Players[e.Player].Stars > 0 {
		g.Players[e.Player].Stars--
	}
	g.MovePlayer(e.Player, e.Moves)
}

//ESWarpCDest decides which Warp C destination the player goes to.
type ESWarpCDest struct {
	Player int
	Moves  int
}

func (e ESWarpCDest) Question(g *mp1.Game) string {
	return fmt.Sprintf("Where did %s warp to?", g.Players[e.Player].Char)
}

func (e ESWarpCDest) Type() mp1.EventType {
	return mp1.CHAINSPACE_EVT_TYPE
}

//Resopnses returns a slice of the 2 possible spaces the player can warp
//to.
func (e ESWarpCDest) Responses() []mp1.Response {
	return []mp1.Response{esEntrance1, esEntrance7}
}

func (e ESWarpCDest) ControllingPlayer() int {
	return mp1.CPU_PLAYER
}

//Handle moves the player to the ChainSpace r and sets various flags if
//needed.
func (e ESWarpCDest) Handle(r mp1.Response, g *mp1.Game) {
	dest := r.(mp1.ChainSpace)
	g.Players[e.Player].CurrentSpace = dest

	bd := g.Board.Data.(esBoardData)
	if dest == esEntrance7 {
		bd.Gate = 1
	} else {
		bd.Gate2or3 = true
	}
	g.Board.Data = bd

	g.MovePlayer(e.Player, e.Moves)
}

//ESWarpDest decides which gate the board is playing under currently.
type ESWarpDest struct {
	Player   int
	Moves    int
	Gate2or3 bool
	Island1  mp1.ChainSpace
	Island2  mp1.ChainSpace
	Island3  mp1.ChainSpace
}

func (e ESWarpDest) Question(g *mp1.Game) string {
	return fmt.Sprintf("Where did %s warp to?", g.Players[e.Player].Char)
}

func (e ESWarpDest) Type() mp1.EventType {
	return mp1.CHAINSPACE_EVT_TYPE
}

//Responses returns the list of possible ChainSpaces that the player can
//warp to.
func (e ESWarpDest) Responses() []mp1.Response {
	ret := []mp1.Response{
		e.Island1,
		e.Island2,
		e.Island3,
	}
	if e.Gate2or3 {
		ret = ret[1:]
	}
	return ret
}

func (e ESWarpDest) ControllingPlayer() int {
	return mp1.CPU_PLAYER
}

//Handle moves the player to the ChainSpace in r and set's the current
//gate the board is under in r.
func (e ESWarpDest) Handle(r mp1.Response, g *mp1.Game) {
	dest := r.(mp1.ChainSpace)
	bd := g.Board.Data.(esBoardData)

	//Set Gate
	switch dest {
	case e.Island1:
		bd.Gate = 1
	case e.Island2:
		bd.Gate = 2
	case e.Island3:
		bd.Gate = 3
	}
	g.Board.Data = bd
	g.Players[e.Player].CurrentSpace = dest
	g.MovePlayer(e.Player, e.Moves)
}

//ESChangeGates decides which Gate the board will change to.
type ESChangeGates struct {
	Player  int
	Moves   int
	Current int
}

type Gate int

func (g Gate) String() string {
	return fmt.Sprintf("Gate %d", int(g))
}

var esChangeGatesResponses = [3][]mp1.Response{
	{Gate(2), Gate(3)},
	{Gate(1), Gate(3)},
	{Gate(1), Gate(2)},
}

func (e ESChangeGates) Question(g *mp1.Game) string {
	return "What warp gate number did Bowser set the course to?"
}

func (e ESChangeGates) Type() mp1.EventType {
	return mp1.ENUM_EVT_TYPE
}

//Responses returns the gates that can be switched to.
func (e ESChangeGates) Responses() []mp1.Response {
	return esChangeGatesResponses[e.Current-1]
}

func (e ESChangeGates) ControllingPlayer() int {
	return mp1.CPU_PLAYER
}

//Handle switches the current gate configuration to r, moves the player to
//the starting space, and moves the player their remaining spaces.
func (e ESChangeGates) Handle(r mp1.Response, g *mp1.Game) {
	gate := r.(Gate)
	bd := g.Board.Data.(esBoardData)
	bd.Gate = int(gate)
	bd.Gate2or3 = (gate != 1)
	g.Board.Data = bd
	g.Players[e.Player].CurrentSpace = esEntrance1
	g.MovePlayer(e.Player, e.Moves)
}
