package mp1

func (g *Game) PreBowserCheck(player int) {
	//Special events when player has 0 coins
	if g.Players[player].Coins == 0 {
		if g.Players[player].Stars > 0 {
			g.AwardCoins(player, 10, false)
			g.Players[player].Stars--
		} else {
			g.AwardCoins(player, 20, false)
		}
		g.EndCharacterTurn()
	} else {
		g.ExtraEvent = BowserEvent{player}
	}
}

type BowserEvent struct {
	Player int
}

type BowserResponse int

const (
	CoinsForBowser BowserResponse = iota
	BowserBalloonBurst
	BowsersFaceLift
	BowsersTugoWar
	BashnCash
	BowserRevolution
	BowsersChanceTime
	StarPresent
)

func (b BowserEvent) Responses() []Response {
	return []Response{
		CoinsForBowser,
		BowserBalloonBurst,
		BowsersFaceLift,
		BowsersTugoWar,
		BashnCash,
		BowserRevolution,
		BowsersChanceTime,
		StarPresent,
	}
}

func (b BowserEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

func (b BowserEvent) Handle(r Response, g *Game) {
	choice := r.(BowserResponse)
	switch choice {
	case CoinsForBowser:
		coinsLost := GetBowserMinigameCoinLoss(g.Turn)
		g.AwardCoins(b.Player, -coinsLost, false)
		g.EndCharacterTurn()
	case BowserBalloonBurst:
		g.ExtraEvent = BowserBalloonBurstEvent{}
	case BowsersFaceLift:
		g.ExtraEvent = BowsersFaceLiftEvent{b.Player}
	case BowsersTugoWar:
		g.ExtraEvent = BowsersTugoWarEvent{b.Player}
	case BashnCash:
		g.ExtraEvent = BowsersBashnCash{b.Player, g.Players[b.Player].Coins}
	case BowserRevolution:
		coins := 0
		for i := range g.Players {
			coins += g.Players[i].Coins
		}
		coins /= 4
		for i := range g.Players {
			g.Players[i].Coins = coins
		}
		g.EndCharacterTurn()
	case BowsersChanceTime:
		g.ExtraEvent = BowsersChanceTimeEvent{}
	}
}

type BowserBalloonBurstEvent struct{}

type BowserBalloonBurstResult int

const (
	BBBDraw BowserBalloonBurstResult = iota
	BBBP1Win
	BBBP2Win
	BBBP3Win
	BBBP4Win
)

func GetBowserMinigameCoinLoss(turn uint8) int {
	if turn <= 9 {
		return 10
	} else if turn <= 19 {
		return 20
	} else if turn <= 29 {
		return 30
	}
	return 40
}

var BBBResults = []Response{
	BBBDraw,
	BBBP1Win,
	BBBP2Win,
	BBBP3Win,
	BBBP4Win,
}

func (b BowserBalloonBurstEvent) Responses() []Response {
	return BBBResults
}

func (b BowserBalloonBurstEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

func (b BowserBalloonBurstEvent) Handle(r Response, g *Game) {
	results := r.(BowserBalloonBurstResult)
	coinLoss := -GetBowserMinigameCoinLoss(g.Turn)
	switch results {
	case BBBDraw:
		for p := range g.Players {
			g.AwardCoins(p, -20, true)
		}
	case BBBP1Win:
		g.AwardCoins(1, coinLoss, true)
		g.AwardCoins(2, coinLoss, true)
		g.AwardCoins(3, coinLoss, true)
	case BBBP2Win:
		g.AwardCoins(0, coinLoss, true)
		g.AwardCoins(2, coinLoss, true)
		g.AwardCoins(3, coinLoss, true)
	case BBBP3Win:
		g.AwardCoins(0, coinLoss, true)
		g.AwardCoins(1, coinLoss, true)
		g.AwardCoins(3, coinLoss, true)
	case BBBP4Win:
		g.AwardCoins(0, coinLoss, true)
		g.AwardCoins(1, coinLoss, true)
		g.AwardCoins(2, coinLoss, true)
	}
	g.EndCharacterTurn()
}

type BowsersFaceLiftEvent struct {
	Player int
}

func (b BowsersFaceLiftEvent) Responses() []Response {
	return CPURangeEvent{1, 15}.Responses()
}

func (b BowsersFaceLiftEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

func (b BowsersFaceLiftEvent) Handle(r Response, g *Game) {
	results := r.(int)
	if results == 15 { //All players won
		g.AwardCoins(b.Player, -50, true)
		return
	}

	coinLoss := -GetBowserMinigameCoinLoss(g.Turn)
	for p := range g.Players {
		if results&(1<<p) == 0 {
			g.AwardCoins(p, coinLoss, true)
		}
	}
	g.EndCharacterTurn()
}

type BowsersTugoWarEvent struct {
	Player int
}

type BowsersTugoWarResult int

const (
	BTWDraw BowsersTugoWarResult = iota
	BTW1TWin
	BTW3TWin
)

var BTWResults = []Response{
	BTWDraw,
	BTW1TWin,
	BTW3TWin,
}

func (b BowsersTugoWarEvent) Responses() []Response {
	return BTWResults
}

func (b BowsersTugoWarEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

func (b BowsersTugoWarEvent) Handle(r Response, g *Game) {
	results := r.(BowsersTugoWarResult)
	switch results {
	case BTWDraw:
		for p := range g.Players {
			g.AwardCoins(p, -30, true)
		}
	case BTW1TWin:
		coinLoss := -GetBowserMinigameCoinLoss(g.Turn)
		for p := range g.Players {
			if p != b.Player {
				g.AwardCoins(p, coinLoss, true)
			}
		}
	case BTW3TWin:
		g.AwardCoins(b.Player, -10, true)
	}
	g.EndCharacterTurn()
}

type BowsersBashnCash struct {
	Player int
	Coins  int
}

func (b BowsersBashnCash) Responses() []Response {
	max := b.Coins / 5
	max += b.Coins % 5
	return CPURangeEvent{1, max}.Responses()
}

func (b BowsersBashnCash) ControllingPlayer() int {
	return CPU_PLAYER
}

func (b BowsersBashnCash) Handle(r Response, g *Game) {
	timesHit := r.(int)
	coinsLost := 0
	if b.Coins/5 < timesHit {
		coinsLost += b.Coins - (b.Coins % 5)
		timesHit -= b.Coins / 5
		coinsLost += timesHit
	} else {
		coinsLost += timesHit * 5
	}
	g.AwardCoins(b.Player, -coinsLost, true)
	g.EndCharacterTurn()
}

type BowsersChanceTimeEvent struct{}

type BCTResponse struct {
	Player int
	Coins  int
}

var BCTResponses = []Response{
	BCTResponse{0, 10},
	BCTResponse{0, 20},
	BCTResponse{0, 30},
	BCTResponse{1, 10},
	BCTResponse{1, 20},
	BCTResponse{1, 30},
	BCTResponse{2, 10},
	BCTResponse{2, 20},
	BCTResponse{2, 30},
}

func (b BowsersChanceTimeEvent) Responses() []Response {
	return BCTResponses
}

func (b BowsersChanceTimeEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

func (b BowsersChanceTimeEvent) Handle(r Response, g *Game) {
	res := r.(BCTResponse)
	g.AwardCoins(res.Player, -res.Coins, false)
	g.EndCharacterTurn()
}
