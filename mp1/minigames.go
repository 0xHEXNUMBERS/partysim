package mp1

import (
	"fmt"
	"strconv"
)

//MinigameFFAReward handles Free-For-All minigame rewards. One player
//will gain coins from this event. It contains a CoinMinigame if extra
//coins can be gained by any player.
type MinigameFFAReward struct {
	IsCoinMinigame bool
	Coin           CoinMinigameFFAReward
}

func (m MinigameFFAReward) Question(g *Game) string {
	return "Which character won the minigame?"
}

func (m MinigameFFAReward) Type() EventType {
	return PLAYER_EVT_TYPE
}

//DrawableFFAReward handles drawable Free-For-All minigame rewards. One
//player may gain coins from this event. It contains a CoinMinigame if
//extra coins can be gained by any player.
type DrawableFFAReward struct {
	MinigameFFAReward
}

var DrawableFFAPlayers = []Response{0, 1, 2, 3, 4}

var MinigameFFAPlayers = DrawableFFAPlayers[:4]

//Responses returns a slice of ints from [0, 4]
func (d DrawableFFAReward) Responses() []Response {
	return DrawableFFAPlayers
}

//Responses returns a slice of ints from [0, 3]
func (m MinigameFFAReward) Responses() []Response {
	return MinigameFFAPlayers
}

func (m MinigameFFAReward) ControllingPlayer() int {
	return CPU_PLAYER
}

//MinigameFFAReward gives out 10 coins to player r. If r == 4, then no one
//gains coins. If m.IsCoinMinigame is true, then the game's next event is
//set to the containing coin minigame. Otherwise, the game's turn ends.
func (m MinigameFFAReward) Handle(r Response, g *Game) {
	player := r.(int)
	if player != 4 {
		g.AwardCoins(player, 10, true)
	}
	if m.IsCoinMinigame {
		g.NextEvent = m.Coin
	} else {
		g.EndGameTurn()
	}
}

//CoinMinigameFFAReward distributes coins gained from coin minigames.
type CoinMinigameFFAReward struct {
	Range
	Player int
}

func (c CoinMinigameFFAReward) Question(g *Game) string {
	return fmt.Sprintf("How many extra coins did %s gain?",
		g.Players[c.Player].Char)
}

func (c CoinMinigameFFAReward) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle gives r coins to c.Player. The games next event is either set to
//the same event with c.Player incremented, or ends the game's turn.
func (c CoinMinigameFFAReward) Handle(r Response, g *Game) {
	coins := r.(int)
	g.AwardCoins(c.Player, coins, true)
	if c.Player == 3 {
		g.EndGameTurn()
	} else {
		c.Player++
		g.NextEvent = c
	}
}

//MinigameFFAMultiWinReward handles Free-For-All minigame rewards. 0 or
//more players will gain coins from this event, and 0 or more players will
//lose coins from this event.
type MinigameFFAMultiWinReward struct {
	CoinsToWin      int
	CoinsToLose     int
	CoinsIfNoWinner int
}

func (m MinigameFFAMultiWinReward) Question(g *Game) string {
	return "Which players won this minigame?"
}

func (m MinigameFFAMultiWinReward) Type() EventType {
	return MULTIWIN_PLAYER_EVT_TYPE
}

//Responses returns a slice of ints from [0, 15].
func (m MinigameFFAMultiWinReward) Responses() []Response {
	return NewRange(0, 15)
}

func (m MinigameFFAMultiWinReward) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle gives out gives out coins to players based on the value of r. The
//bits 0-3 in r determine if player x won the minigame. For example: if r
//is 9, player 0 and player 2 both win coins, while player 1 and player 3
//lose coins.
func (m MinigameFFAMultiWinReward) Handle(r Response, g *Game) {
	wins := r.(int)
	defer g.EndGameTurn()
	if wins == 0 {
		for p := 0; p < 4; p++ {
			g.AwardCoins(p, m.CoinsIfNoWinner, true)
		}
		return
	}
	for p := 0; p < 4; p++ {
		if wins&(1<<p) > 0 {
			g.AwardCoins(p, m.CoinsToWin, true)
		} else {
			g.AwardCoins(p, m.CoinsToLose, true)
		}
	}
}

//MinigameFFA1Loser handles Free-For-All minigame rewards. Either 1 player
//gives the other 3 5 coins each, or all players win 10 coins.
type MinigameFFA1Loser struct{}

func (m MinigameFFA1Loser) Question(g *Game) string {
	return "Which player lost this minigame?"
}

var FFA1LoserPlayers = DrawableFFAPlayers

func (m MinigameFFA1Loser) Type() EventType {
	return PLAYER_EVT_TYPE
}

//Responses returns a slice of ints from [0, 4].
func (m MinigameFFA1Loser) Responses() []Response {
	return FFA1LoserPlayers
}

func (m MinigameFFA1Loser) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle gives coins to player based on r. If r == 4, then all players win
//10 coins. Otherwise, player r gives 5 coins to every other player.
func (m MinigameFFA1Loser) Handle(r Response, g *Game) {
	player := r.(int)
	defer g.EndGameTurn()
	if player == 4 {
		for i := 0; i < 4; i++ {
			g.AwardCoins(i, 10, true)
		}
		return
	}
	for i := 0; i < 4; i++ {
		if i == player {
			g.AwardCoins(i, -15, true)
		} else {
			g.AwardCoins(i, 5, true)
		}
	}
}

type MinigameFFACoopResponse int

const (
	MinigameFFACoopWin MinigameFFACoopResponse = iota
	MinigameFFACoopLoss
)

func (m MinigameFFACoopResponse) String() string {
	switch m {
	case MinigameFFACoopWin:
		return "Everyone won the minigame"
	case MinigameFFACoopLoss:
		return "Everyone lost the minigame"
	}
	return ""
}

//MinigameFFACoop handles Free-For-All cooperative minigame rewards.
//Players either win 10 coins each or lose 5 coins each.
type MinigameFFACoop struct{}

func (m MinigameFFACoop) Question(g *Game) string {
	return "Did everyone win this minigame?"
}

func (m MinigameFFACoop) Type() EventType {
	return ENUM_EVT_TYPE
}

func (m MinigameFFACoop) ControllingPlayer() int {
	return CPU_PLAYER
}

func (m MinigameFFACoop) Responses() []Response {
	return []Response{MinigameFFACoopWin, MinigameFFACoopLoss}
}

//Handle gives coins to players based on r. If r == true, then each player
//gains 10 coins. Otherwise, each player loses 5 coins.
func (m MinigameFFACoop) Handle(r Response, g *Game) {
	won := r.(MinigameFFACoopResponse)
	coins := -5
	if won == MinigameFFACoopWin {
		coins = 10
	}
	for i := 0; i < 4; i++ {
		g.AwardCoins(i, coins, true)
	}
	g.EndGameTurn()
}

//MinigameGrabBag handles the grab bag FFA minigame. Players steal coins
//from each other in the duration of this minigame.
type MinigameGrabBag struct {
	Range
	Player int
	Acc    int
}

func (m MinigameGrabBag) Question(g *Game) string {
	return fmt.Sprintf("How many coins did %s gain/lose?",
		g.Players[m.Player].Char)
}

func (m MinigameGrabBag) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle gives player m.Player r coins. An accumulator is kept to save one
//round of event execution. If m.Player is 2, then Player 3 receives m.Acc
//coins. Otherwise, the next event is set to this event with m.Player
//incremented.
func (m MinigameGrabBag) Handle(r Response, g *Game) {
	coins := r.(int)
	m.Acc += coins
	g.AwardCoins(m.Player, coins, true)
	m.Player++
	if m.Player >= 3 {
		g.AwardCoins(3, m.Acc, true)
		g.EndGameTurn()
	} else {
		g.NextEvent = m
	}
}

//MinigameFFAGame is a enumeration of the available FFA minigames.
type MinigameFFAGame int

const (
	MinigameFFABurriedTreasure MinigameFFAGame = iota
	MinigameFFATreasureDivers
	MinigameFFAHotBobomb
	MinigameFFAMusicalMushroom
	MinigameFFACrazyCutter
	MinigameFFAFaceLift
	MinigameFFABalloonBurst
	MinigameFFACoinBlockBlitz
	MinigameFFASkateboardScamper
	MinigameFFABoxMountainMayhem
	MinigameFFAPlatformPeril
	MinigameFFAMushroomMixup
	MinigameFFAGrabBag
	MinigameFFABumperBalls
	MinigameFFATipsyTourney
	MinigameFFABombsAway
	MinigameFFAMarioBandstand
	MinigameFFAShyGuySays
	MinigameFFACastAways
	MinigameFFAKeypaWay
	MinigameFFARunningoftheBulb
	MinigameFFAHotRopeJump
	MinigameFFAHammerDrop
	MinigameFFASlotCarDerby
)

func (m MinigameFFAGame) String() string {
	switch m {
	case MinigameFFABurriedTreasure:
		return "Burried Treasure"
	case MinigameFFATreasureDivers:
		return "Treasure Divers"
	case MinigameFFAHotBobomb:
		return "Hot Bobomb"
	case MinigameFFAMusicalMushroom:
		return "Musical Mushroom"
	case MinigameFFACrazyCutter:
		return "Crazy Cutter"
	case MinigameFFAFaceLift:
		return "Face Lift"
	case MinigameFFABalloonBurst:
		return "Balloon Burst"
	case MinigameFFACoinBlockBlitz:
		return "Coin Block Blitz"
	case MinigameFFASkateboardScamper:
		return "Skateboard Scamper"
	case MinigameFFABoxMountainMayhem:
		return "Box Mountain Mayhem"
	case MinigameFFAPlatformPeril:
		return "Platform Peril"
	case MinigameFFAMushroomMixup:
		return "Mushroom Mixup"
	case MinigameFFAGrabBag:
		return "Grab Bag"
	case MinigameFFABumperBalls:
		return "Bumper Balls"
	case MinigameFFATipsyTourney:
		return "Tipsy Tourney"
	case MinigameFFABombsAway:
		return "Bombs Away"
	case MinigameFFAMarioBandstand:
		return "Mario Bandstand"
	case MinigameFFAShyGuySays:
		return "Shy Guy Says"
	case MinigameFFACastAways:
		return "Cast Aways"
	case MinigameFFAKeypaWay:
		return "Key Pa Way"
	case MinigameFFARunningoftheBulb:
		return "Running of the Bulb"
	case MinigameFFAHotRopeJump:
		return "Hot Rope Jump"
	case MinigameFFAHammerDrop:
		return "Hammer Drop"
	case MinigameFFASlotCarDerby:
		return "Slot Car Derby"
	}
	return ""
}

//MinigameFFASelector selects which FFA minigame to play.
type MinigameFFASelector struct{}

func (m MinigameFFASelector) Question(g *Game) string {
	return "Which minigame was selected?"
}

func (m MinigameFFASelector) Type() EventType {
	return ENUM_EVT_TYPE
}

//Responses returns a slice of all MinigameFFAGame enumerations.
func (m MinigameFFASelector) Responses() []Response {
	return []Response{
		MinigameFFABurriedTreasure,
		MinigameFFATreasureDivers,
		MinigameFFAHotBobomb,
		MinigameFFAMusicalMushroom,
		MinigameFFACrazyCutter,
		MinigameFFAFaceLift,
		MinigameFFABalloonBurst,
		MinigameFFACoinBlockBlitz,
		MinigameFFASkateboardScamper,
		MinigameFFABoxMountainMayhem,
		MinigameFFAPlatformPeril,
		MinigameFFAMushroomMixup,
		MinigameFFAGrabBag,
		MinigameFFABumperBalls,
		MinigameFFATipsyTourney,
		MinigameFFABombsAway,
		MinigameFFAMarioBandstand,
		MinigameFFAShyGuySays,
		MinigameFFACastAways,
		MinigameFFAKeypaWay,
		MinigameFFARunningoftheBulb,
		MinigameFFAHotRopeJump,
		MinigameFFAHammerDrop,
		MinigameFFASlotCarDerby,
	}
}

func (m MinigameFFASelector) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle sets the next event to the selected Minigame.
func (m MinigameFFASelector) Handle(r Response, g *Game) {
	game := r.(MinigameFFAGame)
	switch game {
	case MinigameFFABurriedTreasure:
		g.NextEvent = MinigameFFAReward{}
	case MinigameFFATreasureDivers:
		g.NextEvent = CoinMinigameFFAReward{Range{0, 50}, 0}
	case MinigameFFAHotBobomb:
		g.NextEvent = MinigameFFA1Loser{}
	case MinigameFFAMusicalMushroom:
		g.NextEvent = MinigameFFAReward{}
	case MinigameFFACrazyCutter:
		g.NextEvent = MinigameFFAMultiWinReward{10, -5, 0}
	case MinigameFFAFaceLift:
		g.NextEvent = MinigameFFAMultiWinReward{10, -5, 0}
	case MinigameFFABalloonBurst:
		g.NextEvent = MinigameFFAReward{}
	case MinigameFFACoinBlockBlitz:
		g.NextEvent = CoinMinigameFFAReward{Range{0, 40}, 0}
	case MinigameFFASkateboardScamper:
		//TODO: Separate coin from coinbag
		g.NextEvent = MinigameFFAReward{true, CoinMinigameFFAReward{Range{0, 10}, 0}}
	case MinigameFFABoxMountainMayhem:
		g.NextEvent = CoinMinigameFFAReward{Range{0, 25}, 0}
	case MinigameFFAPlatformPeril:
		//TODO: Separate coin from coinbag
		g.NextEvent = MinigameFFAReward{true, CoinMinigameFFAReward{Range{0, 10}, 0}}
	case MinigameFFAMushroomMixup:
		g.NextEvent = DrawableFFAReward{}
	case MinigameFFAGrabBag:
		//TODO: I'm not sure how I feel about max of 50.
		//Theoritically, players can steal > 50 coins, but probably not
		//feasible
		g.NextEvent = MinigameGrabBag{Range{-50, 50}, 0, 0}
	case MinigameFFABumperBalls:
		g.NextEvent = DrawableFFAReward{}
	case MinigameFFATipsyTourney:
		g.NextEvent = MinigameFFAReward{}
	case MinigameFFABombsAway:
		g.NextEvent = MinigameFFAReward{}
	case MinigameFFAMarioBandstand:
		//TODO: Find out how many coins are distributed
		g.NextEvent = MinigameFFAMultiWinReward{}
	case MinigameFFAShyGuySays:
		g.NextEvent = MinigameFFAReward{}
	case MinigameFFACastAways:
		//TODO: Should optimize; ask for chests/bags/coins
		g.NextEvent = CoinMinigameFFAReward{Range{0, 80}, 0}
	case MinigameFFAKeypaWay:
		g.NextEvent = MinigameFFACoop{}
	case MinigameFFARunningoftheBulb:
		g.NextEvent = MinigameFFAMultiWinReward{10, 0, -5}
	case MinigameFFAHotRopeJump:
		g.NextEvent = MinigameFFA1Loser{}
	case MinigameFFAHammerDrop:
		g.NextEvent = CoinMinigameFFAReward{Range{0, 20}, 0}
	case MinigameFFASlotCarDerby:
		g.NextEvent = MinigameFFAReward{}
	}
}

//Minigame2V2Reward handles 2v2 minigame rewards. One team will gain coins
//from this event, while the other team will lose coins.
type Minigame2V2Reward struct {
	BlueTeam [2]int
	RedTeam  [2]int
}

func (m Minigame2V2Reward) Question(g *Game) string {
	return fmt.Sprintf(
		"Which team won the minigame (Blue: %s/%s | Red: %s/%s",
		g.Players[m.BlueTeam[0]].Char,
		g.Players[m.BlueTeam[1]].Char,
		g.Players[m.RedTeam[0]].Char,
		g.Players[m.RedTeam[1]].Char,
	)
}

//Drawable2V2Reward handles 2v2 minigame rewards. Zero or One team will
//gain coins from this event, while the other team may lose coins.
type DrawableMinigame2V2Reward Minigame2V2Reward

var Drawable2V2Players = []Response{
	Minigame2V2BlueWin, Minigame2V2RedWin, Minigame2V2Draw,
}

var Minigame2V2Players = Drawable2V2Players[:2]

type Minigame2V2Result int

const (
	Minigame2V2BlueWin Minigame2V2Result = iota
	Minigame2V2RedWin
	Minigame2V2Draw
)

func (m Minigame2V2Result) String() string {
	switch m {
	case Minigame2V2BlueWin:
		return "Blue Team Wins"
	case Minigame2V2RedWin:
		return "Red Team Wins"
	case Minigame2V2Draw:
		return "Draw"
	}
	return ""
}

//Responses returns a slice of ints from [0, 2]
func (d DrawableMinigame2V2Reward) Responses() []Response {
	return Drawable2V2Players
}

func (m Minigame2V2Reward) Type() EventType {
	return ENUM_EVT_TYPE
}

//Responses returns a slice of ints from [0, 1]
func (m Minigame2V2Reward) Responses() []Response {
	return Minigame2V2Players
}

func (m Minigame2V2Reward) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle gives coins out to team members based on r. If r == 0, then the
//blue team takes coins from the red team. If r == 1, then the red team
//takes coins from the blue team. Otherwise, it is considered a draw.
func (m Minigame2V2Reward) Handle(r Response, g *Game) {
	team := r.(Minigame2V2Result)
	if team == Minigame2V2BlueWin {
		g.AwardCoins(m.BlueTeam[0], 10, true)
		g.AwardCoins(m.BlueTeam[1], 10, true)
		g.AwardCoins(m.RedTeam[0], -10, true)
		g.AwardCoins(m.RedTeam[1], -10, true)
	} else if team == Minigame2V2RedWin {
		g.AwardCoins(m.RedTeam[0], 10, true)
		g.AwardCoins(m.RedTeam[1], 10, true)
		g.AwardCoins(m.BlueTeam[0], -10, true)
		g.AwardCoins(m.BlueTeam[1], -10, true)
	}
	g.EndGameTurn()
}

//CoinMinigame2V2Reward distributes coins gained from 2v2 coin minigames.
type CoinMinigame2V2Reward struct {
	Range
	BlueTeam [2]int
	RedTeam  [2]int
	Team     int
}

func (c CoinMinigame2V2Reward) Question(g *Game) string {
	var team string
	var char0, char1 string
	if c.Team == 0 {
		team = "blue"
		char0 = g.Players[c.BlueTeam[0]].Char
		char0 = g.Players[c.BlueTeam[1]].Char
	} else if c.Team == 1 {
		team = "red"
		char0 = g.Players[c.RedTeam[0]].Char
		char0 = g.Players[c.RedTeam[1]].Char
	}

	return fmt.Sprintf("How many coins did the %s team gain (%s/%s)",
		team, char0, char1)
}

func (c CoinMinigame2V2Reward) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle gives r coins to team c.Team (0 == Blue, 1 == Red).
func (c CoinMinigame2V2Reward) Handle(r Response, g *Game) {
	coins := r.(int)
	if c.Team == 0 {
		g.AwardCoins(c.BlueTeam[0], coins, true)
		g.AwardCoins(c.BlueTeam[1], coins, true)
	} else if c.Team == 1 {
		g.AwardCoins(c.RedTeam[0], coins, true)
		g.AwardCoins(c.RedTeam[1], coins, true)
		g.EndGameTurn()
		return
	}
	c.Team++
	g.NextEvent = c
}

//Minigame2V2Game is a enumeration of the available 2V2 minigames.
type Minigame2V2Game int

const (
	Minigame2V2BobsledRun Minigame2V2Game = iota
	Minigame2V2DesertDash
	Minigame2V2Bombsketball
	Minigame2V2HandcarHavoc
	Minigame2V2DeepSeaDivers
)

func (m Minigame2V2Game) String() string {
	switch m {
	case Minigame2V2BobsledRun:
		return "Bobsled Run"
	case Minigame2V2DesertDash:
		return "Desert Dash"
	case Minigame2V2Bombsketball:
		return "Bombsketball"
	case Minigame2V2HandcarHavoc:
		return "Handcar Havoc"
	case Minigame2V2DeepSeaDivers:
		return "Deep Sea Divers"
	}
	return ""
}

//Minigame2V2Selector selects which 2V2 minigame to play.
type Minigame2V2Selector struct {
	Team1 [2]int
	Team2 [2]int
}

func (m Minigame2V2Selector) Question(g *Game) string {
	return "Which minigame was selected?"
}

func (m Minigame2V2Selector) Type() EventType {
	return ENUM_EVT_TYPE
}

//Responses returns a slice of all Minigame2V2Game enumerations.
func (m Minigame2V2Selector) Responses() []Response {
	return []Response{
		Minigame2V2BobsledRun,
		Minigame2V2DesertDash,
		Minigame2V2Bombsketball,
		Minigame2V2HandcarHavoc,
		Minigame2V2DeepSeaDivers,
	}
}

func (m Minigame2V2Selector) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle sets the next event to the selected Minigame.
func (m Minigame2V2Selector) Handle(r Response, g *Game) {
	game := r.(Minigame2V2Game)
	switch game {
	case Minigame2V2BobsledRun:
		g.NextEvent = Minigame2V2Reward{
			m.Team1, m.Team2,
		}
	case Minigame2V2DesertDash:
		g.NextEvent = Minigame2V2Reward{
			m.Team1, m.Team2,
		}
	case Minigame2V2Bombsketball:
		g.NextEvent = Minigame2V2Reward{
			m.Team1, m.Team2,
		}
	case Minigame2V2HandcarHavoc:
		g.NextEvent = Minigame2V2Reward{
			m.Team1, m.Team2,
		}
	case Minigame2V2DeepSeaDivers:
		g.NextEvent = CoinMinigame2V2Reward{
			Range{0, 50}, m.Team1, m.Team2, 0,
		}
	}
}

//Minigame1V3Reward handles 1v3 minigame rewards. One team will gain coins
//from this event, while the other team will lose coins.
type Minigame1V3Reward struct {
	SingleTeam int
}

type Minigame1V3Result int

const (
	Minigame1V3SingleWin Minigame1V3Result = iota
	Minigame1V3TeamWin
	Minigame1V3Draw
)

func (m Minigame1V3Result) String() string {
	switch m {
	case Minigame1V3SingleWin:
		return "Single Player Wins"
	case Minigame1V3TeamWin:
		return "Team of 3 Wins"
	case Minigame1V3Draw:
		return "Draw"
	}
	return ""
}

//Drawable1V3Reward handles 1v3 minigame rewards. Zero or One team will
//gain coins from this event, while the other team may lose coins.
type Drawable1V3Reward Minigame1V3Reward

var Drawable1V3Players = []Response{
	Minigame1V3SingleWin, Minigame1V3TeamWin, Minigame1V3Draw,
}
var Minigame1V3Players = Drawable1V3Players[:2]

func (m Minigame1V3Reward) Question(g *Game) string {
	return fmt.Sprintf("Which team won the minigame? (Single Player: %s)",
		g.Players[m.SingleTeam].Char)
}

func (m Minigame1V3Reward) Type() EventType {
	return ENUM_EVT_TYPE
}

//Responses returns a slice of ints from [0, 2]
func (d Drawable1V3Reward) Responses() []Response {
	return Drawable1V3Players
}

//Responses returns a slice of ints from [0, 1]
func (m Minigame1V3Reward) Responses() []Response {
	return Minigame1V3Players
}

func (m Minigame1V3Reward) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle gives coins out to team members based on r. If r == 0, then the
//solo player takes coins from the 3player team. If r == 1, then the
//3player takes coins from the solo player. Otherwise, it is considered a
//draw.
func (m Minigame1V3Reward) Handle(r Response, g *Game) {
	team := r.(Minigame1V3Result)
	if team == Minigame1V3SingleWin {
		g.AwardCoins(m.SingleTeam, 15, true)
		for i := 0; i < 4; i++ {
			if i != m.SingleTeam {
				g.AwardCoins(i, -5, true)
			}
		}
	} else if team == Minigame1V3TeamWin {
		g.AwardCoins(m.SingleTeam, -15, true)
		for i := 0; i < 4; i++ {
			if i != m.SingleTeam {
				g.AwardCoins(i, 5, true)
			}
		}
	}
	g.EndGameTurn()
}

type Throwable1V3MinigameResponse int

const (
	Throwable1V3MinigameThrow Throwable1V3MinigameResponse = iota
	Throwable1V3MinigameNoThrow
)

func (t Throwable1V3MinigameResponse) String() string {
	switch t {
	case Throwable1V3MinigameThrow:
		return "Single Player throws the minigame"
	case Throwable1V3MinigameNoThrow:
		return "Single Player plays the minigame as normal"
	}
	return ""
}

//Throwable1V3Minigame is a minigame that the Solo player may choose to
//lose, causing no one to gain coins.
type Throwable1V3Minigame struct {
	Player   int
	Minigame Event
}

func (t Throwable1V3Minigame) Question(g *Game) string {
	return fmt.Sprintf(
		"Did %s intentionally fail the minigame (no team gained coins)?",
		g.Players[t.Player].Char,
	)
}

func (t Throwable1V3Minigame) Type() EventType {
	return ENUM_EVT_TYPE
}

func (t Throwable1V3Minigame) ControllingPlayer() int {
	return t.Player
}

func (t Throwable1V3Minigame) Responses() []Response {
	return []Response{
		Throwable1V3MinigameThrow,
		Throwable1V3MinigameNoThrow,
	}
}

//Handle chooses whether to throw the minigame based on r. If r == true,
//then the game turn ends. If r == false, the next event becomes the
//underlying minigame.
func (t Throwable1V3Minigame) Handle(r Response, g *Game) {
	throw := r.(Throwable1V3MinigameResponse)
	if throw == Throwable1V3MinigameThrow {
		g.EndGameTurn()
	} else {
		g.NextEvent = t.Minigame
	}
}

//Specific 1V3 Minigames

//MinigamePipeMaze holds the implementation for Pipe Maze.
type MinigamePipeMaze struct {
	Player int
}

func (m MinigamePipeMaze) Question(g *Game) string {
	return "Which player won this minigame?"
}

func (m MinigamePipeMaze) Type() EventType {
	return PLAYER_EVT_TYPE
}

//Responses returns a slice of ints from [0,3].
func (m MinigamePipeMaze) Responses() []Response {
	return []Response{0, 1, 2, 3}
}

func (m MinigamePipeMaze) ControllingPlayer() int {
	return m.Player
}

//Handle gives player r 10 coins.
func (m MinigamePipeMaze) Handle(r Response, g *Game) {
	player := r.(int)
	g.AwardCoins(player, 10, true)
	g.EndGameTurn()
}

//MinigameBashnCash holds the implementation for Bash n Cash.
type MinigameBashnCash struct {
	BowsersBashnCash
}

//Handle calculates the number of coins taken from the solo player. The
//folowing events are set to distribute those coins among the other 3
//players.
func (m MinigameBashnCash) Handle(r Response, g *Game) {
	//TODO: code is copied from BowsersBashnCash.Handle()
	//Event interface needs to include *Game to remove this
	//effectively
	timesHit := r.(int)
	if timesHit == 0 {
		g.EndGameTurn()
		return
	}
	coinsLost := 0
	if m.Coins/5 < timesHit {
		coinsLost += m.Coins - (m.Coins % 5)
		timesHit -= m.Coins / 5
		coinsLost += timesHit
	} else {
		coinsLost += timesHit * 5
	}
	g.AwardCoins(m.Player, -coinsLost, true)
	nextEvent := MinigameBashnCashCoinAwards{Range{0, coinsLost}, 0, m.Player}
	if m.Player == 0 {
		nextEvent.CurrentPlayer = 1
	}
	g.NextEvent = nextEvent
}

//MinigameBashnCashCoinAwards distributes a set of coins from a player to
//the other 3 players.
type MinigameBashnCashCoinAwards struct {
	Range
	CurrentPlayer int
	LosingPlayer  int
}

func (m MinigameBashnCashCoinAwards) Question(g *Game) string {
	return fmt.Sprintf("How many coins did %s gain?",
		g.Players[m.CurrentPlayer].Char)
}

func (m MinigameBashnCashCoinAwards) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle gives r coins to m.CurrentPlayer. The next event is set to the
//next player to receive coins, or ends the Game's turn if all players
//have received coins.
func (m MinigameBashnCashCoinAwards) Handle(r Response, g *Game) {
	coins := r.(int)
	g.AwardCoins(m.CurrentPlayer, coins, true)
	m.CurrentPlayer++
	if m.CurrentPlayer == m.LosingPlayer {
		m.CurrentPlayer++
	}
	if m.CurrentPlayer >= 4 {
		g.EndGameTurn()
	} else {
		g.NextEvent = m
	}
}

//MinigameBowlOver holds the implementation for Bowl Over.
type MinigameBowlOver struct {
	Player int
}

//MinigameBowlOverResponse is an ending result of the Bowl Over minigame.
type MinigameBowlOverResponse struct {
	SinglePlayer int
	Pins         int
	CharPins     [3]bool
}

func (m MinigameBowlOverResponse) String() string {
	str := fmt.Sprintf("%d regular pins knocked over, and", m.Pins)
	if m.CharPins == [3]bool{} {
		str += " no character pins knocked down"
	} else {
		str += " players "
		i := 0
		for p := 0; p < 4; p++ {
			if m.SinglePlayer == p {
				continue
			}

			if m.CharPins[i] {
				playerNum := strconv.Itoa(p + 1)
				str += playerNum + " "
			}
			i++
		}
		str += "pins knocked over"
	}
	return str
}

func (m MinigameBowlOver) Question(g *Game) string {
	return fmt.Sprintf("How did %s perform?",
		g.Players[m.Player].Char)
}

func (m MinigameBowlOver) Type() EventType {
	return ENUM_EVT_TYPE
}

//Responses returns all valid Bowl Over minigame endings.
func (m MinigameBowlOver) Responses() []Response {
	return []Response{
		MinigameBowlOverResponse{m.Player, 0, [3]bool{false, false, false}},
		MinigameBowlOverResponse{m.Player, 1, [3]bool{false, false, false}},
		MinigameBowlOverResponse{m.Player, 2, [3]bool{false, false, false}},
		MinigameBowlOverResponse{m.Player, 0, [3]bool{false, false, true}},
		MinigameBowlOverResponse{m.Player, 1, [3]bool{false, false, true}},
		MinigameBowlOverResponse{m.Player, 2, [3]bool{false, false, true}},
		MinigameBowlOverResponse{m.Player, 0, [3]bool{false, true, false}},
		MinigameBowlOverResponse{m.Player, 1, [3]bool{false, true, false}},
		MinigameBowlOverResponse{m.Player, 2, [3]bool{false, true, false}},
		MinigameBowlOverResponse{m.Player, 0, [3]bool{false, true, true}},
		MinigameBowlOverResponse{m.Player, 1, [3]bool{false, true, true}},
		MinigameBowlOverResponse{m.Player, 2, [3]bool{false, true, true}},
		MinigameBowlOverResponse{m.Player, 0, [3]bool{true, false, false}},
		MinigameBowlOverResponse{m.Player, 1, [3]bool{true, false, false}},
		MinigameBowlOverResponse{m.Player, 2, [3]bool{true, false, false}},
		MinigameBowlOverResponse{m.Player, 0, [3]bool{true, false, true}},
		MinigameBowlOverResponse{m.Player, 1, [3]bool{true, false, true}},
		MinigameBowlOverResponse{m.Player, 2, [3]bool{true, false, true}},
		MinigameBowlOverResponse{m.Player, 0, [3]bool{true, true, false}},
		MinigameBowlOverResponse{m.Player, 1, [3]bool{true, true, false}},
		MinigameBowlOverResponse{m.Player, 2, [3]bool{true, true, false}},
		MinigameBowlOverResponse{m.Player, 0, [3]bool{true, true, true}},
		MinigameBowlOverResponse{m.Player, 1, [3]bool{true, true, true}},
		MinigameBowlOverResponse{m.Player, 2, [3]bool{true, true, true}},
	}
}

func (m MinigameBowlOver) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle calculates coin awards given a valid Response r.
func (m MinigameBowlOver) Handle(r Response, g *Game) {
	res := r.(MinigameBowlOverResponse)
	g.AwardCoins(m.Player, res.Pins, true)
	pinIndex := 0
	for i := 0; i < 4; i++ {
		if i == m.Player {
			continue
		}
		if res.CharPins[pinIndex] {
			g.GiveCoins(i, m.Player, 3, true)
		}
		pinIndex++
	}
	g.EndGameTurn()
}

//MinigameCraneGameCoins holds the implementation of the coin portion for
//Crane Game.
type MinigameCraneGameCoins struct {
	Player int
}

func (m MinigameCraneGameCoins) Question(g *Game) string {
	return fmt.Sprintf(
		"How many coins did %s gain from picking up an object?",
		g.Players[m.Player].Char)
}

func (m MinigameCraneGameCoins) Type() EventType {
	return ENUM_EVT_TYPE
}

//Responses returns a slice of ints: {0, 1, 5, 10}.
func (m MinigameCraneGameCoins) Responses() []Response {
	return []Response{0, 1, 5, 10}
}

func (m MinigameCraneGameCoins) ControllingPlayer() int {
	return m.Player
}

//Handle gives the solo player r coins. If r == 0, the solo player decided
//to go after one of the 3 other players, and the next event is set to
//MinigameCraneGamePlayers.
func (m MinigameCraneGameCoins) Handle(r Response, g *Game) {
	coins := r.(int)
	if coins > 0 {
		g.AwardCoins(m.Player, coins, true)
		g.EndGameTurn()
	} else {
		groupPlayers := [3]int{}
		for i := 0; i < 3; i++ {
			if i >= m.Player {
				groupPlayers[i] = i + 1
			} else {
				groupPlayers[i] = i
			}
		}
		g.NextEvent = MinigameCraneGamePlayers{
			m.Player,
			groupPlayers,
		}
	}
}

//MinigameCraneGamePlayers holds the implementation of the player portion
//for Crane Game.
type MinigameCraneGamePlayers struct {
	SoloPlayer int
	Team       [3]int
}

func (m MinigameCraneGamePlayers) Question(g *Game) string {
	return fmt.Sprintf("Which player did %s successfully pick up?",
		g.Players[m.SoloPlayer].Char)
}

func (m MinigameCraneGamePlayers) Type() EventType {
	return PLAYER_EVT_TYPE
}

//Responses returns a slice of ints containing the indexes of the 3 other
//players, and 4 for failing to pick up any object.
func (m MinigameCraneGamePlayers) Responses() []Response {
	return []Response{m.Team[0], m.Team[1], m.Team[2], 4}
}

func (m MinigameCraneGamePlayers) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle gives a third of Player r's coins to m.SoloPlayer. If r == 4, the
//solo player gains no coins.
func (m MinigameCraneGamePlayers) Handle(r Response, g *Game) {
	losingPlayer := r.(int)
	if losingPlayer != 4 {
		coins := g.Players[losingPlayer].Coins / 3
		g.GiveCoins(losingPlayer, m.SoloPlayer, coins, true)
	}
	g.EndGameTurn()
}

//MinigamePaddleBattle holds the implementation for Paddle Battle
type MinigamePaddleBattle struct {
	Range
	Player int
}

func (m MinigamePaddleBattle) Question(g *Game) string {
	//TODO: Rephrase question, or refactor Paddle Battle.
	//Question is not clear
	return fmt.Sprintf("How times did %s effectivily get hit?",
		g.Players[m.Player].Char)
}

func (m MinigamePaddleBattle) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle calculates the amount of coins going from the solo player to/from
//the other 3 players. The solo player effectively loses r*3 coins.
func (m MinigamePaddleBattle) Handle(r Response, g *Game) {
	hits := r.(int)
	//TODO: Find out how coins are distributed when m.Player has 1 or 2 coins
	for i := 0; i < 4; i++ {
		if i == m.Player {
			i++
			continue
		}
		g.GiveCoins(m.Player, i, hits, true)
	}
	g.EndGameTurn()
}

//Minigame1V3Game is a enumeration of the available 1V3 minigames.
type Minigame1V3Game int

const (
	Minigame1V3PipeMaze Minigame1V3Game = iota
	Minigame1V3BashnCash
	Minigame1V3BowlOver
	Minigame1V3CoinBlockBash
	Minigame1V3TightropeTreachery
	Minigame1V3CraneGame
	Minigame1V3PiranhaPursuit
	Minigame1V3TugoWar
	Minigame1V3PaddleBattle
	Minigame1V3CoinShowerFlower
)

func (m Minigame1V3Game) String() string {
	switch m {
	case Minigame1V3PipeMaze:
		return "Pipe Maze"
	case Minigame1V3BashnCash:
		return "Bash n Cash"
	case Minigame1V3BowlOver:
		return "Bowl Over"
	case Minigame1V3CoinBlockBash:
		return "Coin Block Bash"
	case Minigame1V3TightropeTreachery:
		return "Tightrope Treachery"
	case Minigame1V3CraneGame:
		return "Crane Game"
	case Minigame1V3PiranhaPursuit:
		return "Piranha Pursuit"
	case Minigame1V3TugoWar:
		return "Tug o War"
	case Minigame1V3PaddleBattle:
		return "Paddle Battle"
	case Minigame1V3CoinShowerFlower:
		return "Coin Shower Flower"
	}
	return ""
}

//Minigame1V3Selector selects which 1V3 minigame to play.
type Minigame1V3Selector struct {
	Player    int
	SoloCoins int
}

func (m Minigame1V3Selector) Question(g *Game) string {
	return "Which minigame was selected?"
}

func (m Minigame1V3Selector) Type() EventType {
	return ENUM_EVT_TYPE
}

//Responses returns a slice of all Minigame1V3Game enumerations. If the
//solo player has 0 coins, then BashnCash is not selected.
func (m Minigame1V3Selector) Responses() []Response {
	if m.SoloCoins == 0 {
		return []Response{
			Minigame1V3PipeMaze,
			Minigame1V3BowlOver,
			Minigame1V3CoinBlockBash,
			Minigame1V3TightropeTreachery,
			Minigame1V3CraneGame,
			Minigame1V3PiranhaPursuit,
			Minigame1V3TugoWar,
			Minigame1V3PaddleBattle,
			Minigame1V3CoinShowerFlower,
		}
	}
	return []Response{
		Minigame1V3PipeMaze,
		Minigame1V3BashnCash,
		Minigame1V3BowlOver,
		Minigame1V3CoinBlockBash,
		Minigame1V3TightropeTreachery,
		Minigame1V3CraneGame,
		Minigame1V3PiranhaPursuit,
		Minigame1V3TugoWar,
		Minigame1V3PaddleBattle,
		Minigame1V3CoinShowerFlower,
	}
}

func (m Minigame1V3Selector) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle sets the next event to the selected Minigame.
func (m Minigame1V3Selector) Handle(r Response, g *Game) {
	minigame := r.(Minigame1V3Game)
	switch minigame {
	case Minigame1V3PipeMaze:
		g.NextEvent = MinigamePipeMaze{m.Player}
	case Minigame1V3BashnCash:
		coins := g.Players[m.Player].Coins
		g.NextEvent = MinigameBashnCash{NewBowsersBashnCash(m.Player, coins)}
	case Minigame1V3BowlOver:
		g.NextEvent = MinigameBowlOver{m.Player}
	case Minigame1V3CoinBlockBash:
		g.NextEvent = CoinMinigameFFAReward{Range{0, 30}, 0}
	case Minigame1V3TightropeTreachery:
		g.NextEvent = Minigame1V3Reward{m.Player}
	case Minigame1V3CraneGame:
		g.NextEvent = MinigameCraneGameCoins{m.Player}
	case Minigame1V3PiranhaPursuit:
		g.NextEvent = Minigame1V3Reward{m.Player}
	case Minigame1V3TugoWar:
		g.NextEvent = Minigame1V3Reward{m.Player}
	case Minigame1V3PaddleBattle:
		g.NextEvent = MinigamePaddleBattle{Range{-10, 10}, m.Player}
	case Minigame1V3CoinShowerFlower:
		g.NextEvent = Throwable1V3Minigame{
			m.Player,
			CoinMinigameFFAReward{Range{0, 30}, 0},
		}
	}
}

//Minigame1PRewards handles 1P minigame rewards. The player will either
//gain 10 coins, or lose 5 coins.
type Minigame1PRewards struct {
	Player int
}

func (m Minigame1PRewards) Question(g *Game) string {
	return fmt.Sprintf("How many coins did %s win?",
		g.Players[m.Player].Char)
}

func (m Minigame1PRewards) Type() EventType {
	return ENUM_EVT_TYPE
}

//Responses returns a slice of ints, {-5, 10}.
func (m Minigame1PRewards) Responses() []Response {
	return []Response{-5, 10}
}

func (m Minigame1PRewards) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle gives m.Player r coins.
func (m Minigame1PRewards) Handle(r Response, g *Game) {
	coins := r.(int)
	g.AwardCoins(m.Player, coins, true)
	g.EndCharacterTurn()
}

//MinigameMemoryMatch holds the implementation for Memory Match.
type MinigameMemoryMatch struct {
	Minigame1PRewards
}

//Responses returns a slice of ints, {0, 2, 4, 6, 10}.
func (m MinigameMemoryMatch) Responses() []Response {
	return []Response{0, 2, 4, 6, 10}
}

//MinigameSlotMachine hodls the implementation for Slot Machine.
type MinigameSlotMachine struct {
	Minigame1PRewards
}

//Responses returns a slice of ints, {0, 1, 3, 5, 6, 8, 10, 20}.
func (m MinigameSlotMachine) Responses() []Response {
	return []Response{
		0, 1, 3, 5,
		6, 8, 10, 20,
	}
}

//MinigameWhackaPlant holds the implementation for Whack a Plant.
type MinigameWhackaPlant struct {
	Minigame1PRewards
}

//Responses returns a slice of Coins from [0, 36].
func (m MinigameWhackaPlant) Responses() []Response {
	return NewRange(0, 36)
}

//MinigameTeeteringTowers holds the implementation for Teetering Towers.
type MinigameTeeteringTowers struct {
	Minigame1PRewards
}

//Responses returns a slice of ints, {-5, 10, 11, 15, 16}.
func (m MinigameTeeteringTowers) Responses() []Response {
	return []Response{ //Mix of coin and coinbag
		-5, 10, 11, 15, 16,
	}
}

//Minigame1PGame is a enumeration of the available 1P minigames.
type Minigame1PGame int

const (
	Minigame1PMemoryMatch Minigame1PGame = iota
	Minigame1PSlotMachine
	Minigame1PShellGame
	Minigame1PGhostGuess
	Minigame1PPedalPower
	Minigame1PWhackaPlant
	Minigame1PGroundPound
	Minigame1PTeeteringTowers
	Minigame1PKnockBlockTower
	Minigame1PLimboDance
)

func (m Minigame1PGame) String() string {
	switch m {
	case Minigame1PMemoryMatch:
		return "Memory Match"
	case Minigame1PSlotMachine:
		return "Slot Machine"
	case Minigame1PShellGame:
		return "Shell Game"
	case Minigame1PGhostGuess:
		return "Ghost Guess"
	case Minigame1PPedalPower:
		return "Pedal Power"
	case Minigame1PWhackaPlant:
		return "Whack a Plant"
	case Minigame1PGroundPound:
		return "Ground Pound"
	case Minigame1PTeeteringTowers:
		return "Teetering Towers"
	case Minigame1PKnockBlockTower:
		return "Knock Block Tower"
	case Minigame1PLimboDance:
		return "Limbo Dance"
	}
	return ""
}

//Minigame1PSelector selects which 1P minigame to play.
type Minigame1PSelector struct {
	Player int
}

func (m Minigame1PSelector) Question(g *Game) string {
	return "Which minigame was selected?"
}

func (m Minigame1PSelector) Type() EventType {
	return ENUM_EVT_TYPE
}

//Responses returns a slice of all Minigame1PGame enumerations.
func (m Minigame1PSelector) Responses() []Response {
	return []Response{
		Minigame1PMemoryMatch,
		Minigame1PSlotMachine,
		Minigame1PShellGame,
		Minigame1PGhostGuess,
		Minigame1PPedalPower,
		Minigame1PWhackaPlant,
		Minigame1PGroundPound,
		Minigame1PTeeteringTowers,
		Minigame1PKnockBlockTower,
		Minigame1PLimboDance,
	}
}

func (m Minigame1PSelector) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle sets the next event to the selected Minigame.
func (m Minigame1PSelector) Handle(r Response, g *Game) {
	game := r.(Minigame1PGame)
	baseGame := Minigame1PRewards{m.Player}
	switch game {
	case Minigame1PMemoryMatch:
		g.NextEvent = MinigameMemoryMatch{baseGame}
	case Minigame1PSlotMachine:
		g.NextEvent = MinigameSlotMachine{baseGame}
	case Minigame1PShellGame:
		g.NextEvent = baseGame
	case Minigame1PGhostGuess:
		g.NextEvent = baseGame
	case Minigame1PPedalPower:
		g.NextEvent = baseGame
	case Minigame1PWhackaPlant:
		g.NextEvent = MinigameWhackaPlant{baseGame}
	case Minigame1PGroundPound:
		g.NextEvent = baseGame
	case Minigame1PTeeteringTowers:
		g.NextEvent = MinigameTeeteringTowers{baseGame}
	case Minigame1PKnockBlockTower:
		g.NextEvent = baseGame
	case Minigame1PLimboDance:
		g.NextEvent = baseGame
	}
}

//MinigameTeam is an enumeration of the available teams players can be on.
type MinigameTeam int

const (
	BlueTeam MinigameTeam = iota
	RedTeam
	GreenTeam
)

func (m MinigameTeam) String() string {
	switch m {
	case BlueTeam:
		return "Blue Team"
	case RedTeam:
		return "Red Team"
	case GreenTeam:
		return "Green Team (undecided)"
	}
	return ""
}

//SpaceToTeam is a mapping from SpaceType to MinigameTeam.
func SpaceToTeam(s SpaceType) MinigameTeam {
	switch s {
	case Blue, Mushroom, MinigameSpace, Chance:
		return BlueTeam
	case Red, Bowser:
		return RedTeam
	default:
		return GreenTeam
	}
}

//GetMinigame sets up the next minigame type (FFA/2v2/1v3). If any players
//are on the *green* team, the next event is set to determine that
//player's team selection.
func (g *Game) GetMinigame() {
	var blueTeam []int
	var redTeam []int
	for i, p := range g.Players {
		if SpaceToTeam(p.LastSpaceType) == BlueTeam {
			blueTeam = append(blueTeam, i)
		} else if SpaceToTeam(p.LastSpaceType) == RedTeam {
			redTeam = append(redTeam, i)
		}
	}

	var minigame Event
	switch len(blueTeam) {
	case 0, 4:
		minigame = MinigameFFASelector{}
	case 1:
		minigame = Minigame1V3Selector{blueTeam[0], g.Players[blueTeam[0]].Coins}
	case 2:
		minigame = Minigame2V2Selector{
			[2]int{blueTeam[0], blueTeam[1]},
			[2]int{redTeam[0], redTeam[1]},
		}
	case 3:
		minigame = Minigame1V3Selector{redTeam[0], g.Players[redTeam[0]].Coins}
	}
	g.NextEvent = minigame
}

//FindGreenPlayer looks through the 4 players to find one that is on the
//*green* team.
func (g *Game) FindGreenPlayer() {
	for i, p := range g.Players {
		if SpaceToTeam(p.LastSpaceType) == GreenTeam {
			g.NextEvent = DeterminePlayerTeamEvent{
				Player: i,
			}
			return
		}
	}
	g.GetMinigame()
}
