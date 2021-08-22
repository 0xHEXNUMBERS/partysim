package mp1

//pbcSeedCheck decides if the player got a toad seed or a bowser seed.
type pbcSeedCheck struct {
	Player int
	Moves  int
}

//Responses returns a slice of bools (true/false).
func (p pbcSeedCheck) Responses() []Response {
	return []Response{true, false}
}

func (p pbcSeedCheck) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle moves the player based on r. If r is true, the player moves to the
//bowser path. If r is false, the player moves to the toad path.
func (p pbcSeedCheck) Handle(r Response, g *Game) {
	bowser := r.(bool)
	if bowser {
		g.Players[p.Player].CurrentSpace = ChainSpace{1, 0}
		data := g.Board.Data.(pbcBoardData)
		data.BowserSeedPlanted = true
		g.Board.Data = data
	} else {
		g.Players[p.Player].CurrentSpace = ChainSpace{0, 0}
	}
	g.MovePlayer(p.Player, p.Moves-1)
}

//pbcPiranhaDecision decides if the player wants to pay 30 coins for a
//piranha.
type pbcPiranhaDecision struct {
	Player  int
	Piranha int
}

//Responses returns a slice of bools (true/false).
func (p pbcPiranhaDecision) Responses() []Response {
	return []Response{true, false}
}

func (p pbcPiranhaDecision) ControllingPlayer() int {
	return p.Player
}

//Handle performs the decision r. If r is true, then the player pays 30
//coins and gains a piranha at their current space. If r is false, nothing
//happens.
func (p pbcPiranhaDecision) Handle(r Response, g *Game) {
	plantPiranha := r.(bool)
	data := g.Board.Data.(pbcBoardData)
	if plantPiranha {
		data.PiranhaPlant[p.Piranha] = p.Player
		data.PiranhaOccupied[p.Piranha] = true
		g.AwardCoins(p.Player, -30, false)
	}
	g.Board.Data = data
	g.EndCharacterTurn()
}
