package mp1

import (
	"testing"
)

func TestMove(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 23}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}
	g.MovePlayer(0, 4)
	expected := ChainSpace{1, 27}
	got := g.Players[0].CurrentSpace
	if expected != got {
		t.Errorf("Position expected: %#v, got: %#v", expected, got)
	}

}

func TestCanPayThwomp(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 23}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}
	g.MovePlayer(0, 10)
	expected := ytiThwompBranchEvent{0, 6, 1}
	got := g.ExtraEvent
	if expected != got {
		t.Errorf("Event expected: %#v, got: %#v", expected, got)
	}
}

func TestCanNotPayThwomp(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 0, ChainSpace{1, 23}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}

	g.MovePlayer(0, 10)
	expectedEvt := NormalDiceBlock{1}
	gotEvt := g.ExtraEvent
	if expectedEvt != gotEvt {
		t.Errorf("Expected event: %#v, got: %#v", expectedEvt, gotEvt)
	}
}

func TestGainCoins(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 23}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
		Config: GameConfig{
			MaxTurns: 25,
		},
	}

	g.MovePlayer(0, 1)
	expected := 13
	got := g.Players[0].Coins
	if expected != got {
		t.Errorf("Coins expected: %d, got: %d", expected, got)
	}
}

func TestPayThwompAndGainCoins(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 23}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
		Config: GameConfig{
			MaxTurns: 25,
		},
	}

	//Move player to invisible space
	g.MovePlayer(0, 10)
	//Accept payment to thwomp 1
	g.ExtraEvent.Handle(true, &g)
	//Pay thwomp 3 coins, move and land on blue space
	g.ExtraEvent.Handle(3, &g)

	expectedEvt := NormalDiceBlock{1}
	gotEvt := g.ExtraEvent
	if expectedEvt != gotEvt {
		t.Errorf("Expected event: %#v, got: %#v", expectedEvt, gotEvt)
	}

	expectedSquare := ChainSpace{0, 12}
	gotSquare := g.Players[0].CurrentSpace
	if expectedSquare != gotSquare {
		t.Errorf("Space expected: %#v, got: %#v", expectedSquare, gotSquare)
	}

	expectedCoins := 10
	gotCoins := g.Players[0].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Coins expected: %d, got: %d", expectedCoins, gotCoins)
	}
}

func TestIgnoreThwomp(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 23}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
		Config: GameConfig{
			MaxTurns: 25,
		},
	}

	g.MovePlayer(0, 10)
	g.ExtraEvent.Handle(false, &g)

	expectedEvt := NormalDiceBlock{1}
	gotEvt := g.ExtraEvent
	if expectedEvt != gotEvt {
		t.Errorf("Expected event: %#v, got: %#v", expectedEvt, gotEvt)
	}

	expectedSquare := ChainSpace{1, 5}
	gotSquare := g.Players[0].CurrentSpace
	if expectedSquare != gotSquare {
		t.Errorf("Space expected: %#v, got: %#v", expectedSquare, gotSquare)
	}

	expectedCoins := 13
	gotCoins := g.Players[0].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Coins expected: %d, got: %d", expectedCoins, gotCoins)
	}
}

func TestStarSwapViaHappening(t *testing.T) {
	t.SkipNow()
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 23}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}

	g.MovePlayer(0, 3)
	expectedEvt := NormalDiceBlock{1}
	gotEvt := g.ExtraEvent
	if expectedEvt != gotEvt {
		t.Errorf("Expected event: %#v, got: %#v", expectedEvt, gotEvt)
	}

	bd := g.Board.Data.(ytiBoardData)
	if bd.StarPosition == ytiRightIslandStar {
		t.Errorf("Expected star position: %#v, got: %#v", ytiRightIslandStar, bd.StarPosition)
	}
}

func TestCoinsOnStart(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{0, 22}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}

	g.MovePlayer(0, 1)
	expectedCoins := 20
	gotCoins := g.Players[0].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Coins expected: %d, got: %d", expectedCoins, gotCoins)
	}
}

func TestMushroomSpace(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{0, 7}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}

	g.MovePlayer(0, 4)
	expected := MushroomEvent{0}
	got := g.ExtraEvent
	if expected != got {
		t.Errorf("Expected event: %#v, got: %#v", expected, got)
	}

	//Received red mushroom
	gRed := g
	gRed.ExtraEvent.Handle(true, &gRed)
	expectedEvent := NormalDiceBlock{0}
	gotEvent := gRed.ExtraEvent
	if expectedEvent != gotEvent {
		t.Errorf("Expected Red Mushroom Event: %#v, got: %#v",
			expectedEvent,
			gotEvent,
		)
	}

	//Received poison mushroom
	gPoison := g
	gPoison.ExtraEvent.Handle(false, &gPoison)
	expectedEvt := NormalDiceBlock{1}
	gotEvt := gPoison.ExtraEvent
	if expectedEvt != gotEvt {
		t.Errorf("Expected event: %#v, got: %#v", expectedEvt, gotEvt)
	}

	if !gPoison.Players[0].SkipTurn {
		t.Errorf("SkipTurn not set")
	}
}

func TestSkipTurnViaMinigame(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{0, 7}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 7}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 7}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 7}),
		},
		ExtraEvent: NormalDiceBlock{0},
		Config:     GameConfig{MaxTurns: 20},
	}

	//All players recieve poison mushroom
	g.ExtraEvent.Handle(4, &g)
	g.ExtraEvent.Handle(false, &g)
	g.ExtraEvent.Handle(4, &g)
	g.ExtraEvent.Handle(false, &g)
	g.ExtraEvent.Handle(4, &g)
	g.ExtraEvent.Handle(false, &g)
	g.ExtraEvent.Handle(4, &g)
	g.ExtraEvent.Handle(false, &g)

	//Perform FFA Minigame
	g.ExtraEvent.Handle(MinigameFFAMusicalMushroom, &g)
	g.ExtraEvent.Handle(0, &g)

	//Should recieve 2nd minigame as all players were poisoned
	expectedEvent := MinigameFFASelector{}
	gotEvent := g.ExtraEvent
	if expectedEvent != gotEvent {
		t.Errorf("Expected Minigame event: %#v, got: %#v",
			expectedEvent, gotEvent,
		)
	}

	//Perform 2nd minigame
	g.ExtraEvent.Handle(MinigameFFAMusicalMushroom, &g)
	g.ExtraEvent.Handle(0, &g)

	var expectedTurn uint8 = 2
	gotTurn := g.Turn
	if expectedTurn != gotTurn {
		t.Errorf("Expected turn #%d, got turn #%d",
			expectedTurn, gotTurn,
		)
	}

	expectedDiceEvent := NormalDiceBlock{0}
	gotDiceEvent := g.ExtraEvent
	if expectedDiceEvent != gotDiceEvent {
		t.Errorf("Expected Dice event: %#v, got: %#v",
			expectedDiceEvent, gotDiceEvent,
		)
	}
}

func TestSkipTurnViaCharacterTurn(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{0, 7}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 7}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 7}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 7}),
		},
		CurrentPlayer: 2,
		ExtraEvent:    NormalDiceBlock{2},
		Config:        GameConfig{MaxTurns: 20},
	}
	g.Players[0].LastSpaceType = Blue
	g.Players[1].LastSpaceType = Blue

	//Player 2 fails mushroom check
	g.ExtraEvent.Handle(4, &g)
	g.ExtraEvent.Handle(false, &g)

	//Player 3 moves to blue space
	g.ExtraEvent.Handle(1, &g)

	//Handle Minigame
	g.ExtraEvent.Handle(MinigameFFAMusicalMushroom, &g)
	g.ExtraEvent.Handle(0, &g)

	//Players 0 & 1 move to blue space
	g.ExtraEvent.Handle(1, &g)
	g.ExtraEvent.Handle(1, &g)

	expectedEvent := NormalDiceBlock{3}
	gotEvent := g.ExtraEvent
	if expectedEvent != gotEvent {
		t.Errorf("Expected event: %#v, got: %#v",
			expectedEvent, gotEvent,
		)
	}
}

func TestStealCoinsViaBoo(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 21}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}
	g.MovePlayer(0, 4) //Land on happening
	g.ExtraEvent.Handle(BooStealAction{0, 1, false}, &g)
	expectedEvent := BooCoinsEvent{
		PayRangeEvent: PayRangeEvent{
			Player: 1,
			Min:    1,
			Max:    10, //Max of 10 coins
		},
		RecvPlayer: 0,
		Moves:      4,
	}
	gotEvent := g.ExtraEvent
	if expectedEvent != gotEvent {
		t.Errorf("Expected movement: %#v, got: %#v",
			expectedEvent,
			gotEvent,
		)
	}

	expectedEvent.Handle(5, &g)
	expectedDaisyCoins := 15
	gotDaisyCoins := g.Players[0].Coins

	if expectedDaisyCoins != gotDaisyCoins {
		t.Errorf("Daisy expected: %d coins, got: %d coins",
			expectedDaisyCoins,
			gotDaisyCoins,
		)
	}

	expectedLuigiCoins := 5
	gotLuigiCoins := g.Players[1].Coins
	if expectedLuigiCoins != gotLuigiCoins {
		t.Errorf("Luigi expected: %d coins, got: %d coins",
			expectedLuigiCoins,
			gotLuigiCoins,
		)
	}
}

func TestStealTooManyCoinsViaBoo(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 21}),
			NewPlayer("Luigi", 0, 4, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}
	g.MovePlayer(0, 4) //Want to land on happening
	g.ExtraEvent.Handle(BooStealAction{0, 1, false}, &g)
	expectedEvent := BooCoinsEvent{
		PayRangeEvent: PayRangeEvent{
			Player: 1,
			Min:    1,
			Max:    4, //Max of 4 coins
		},
		RecvPlayer: 0,
		Moves:      4,
	}
	gotEvent := g.ExtraEvent
	if expectedEvent != gotEvent {
		t.Errorf("Expected event: %#v, got: %#v",
			expectedEvent,
			gotEvent,
		)
	}

	expectedEvent.Handle(4, &g)
	expectedDaisyCoins := 14
	gotDaisyCoins := g.Players[0].Coins

	if expectedDaisyCoins != gotDaisyCoins {
		t.Errorf("Daisy expected: %d coins, got: %d coins",
			expectedDaisyCoins,
			gotDaisyCoins,
		)
	}

	expectedLuigiCoins := 0
	gotLuigiCoins := g.Players[1].Coins
	if expectedLuigiCoins != gotLuigiCoins {
		t.Errorf("Luigi expected: %d coins, got: %d coins",
			expectedLuigiCoins,
			gotLuigiCoins,
		)
	}
}

func TestPassEmptyBooSpace(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 21}),
			NewPlayer("Luigi", 0, 4, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
		Config: GameConfig{
			NoBoo: true,
		},
	}
	g.MovePlayer(0, 4)
	expectedEvt := NormalDiceBlock{1}
	gotEvt := g.ExtraEvent
	if expectedEvt != gotEvt {
		t.Errorf("Expected event: %#v, got: %#v", expectedEvt, gotEvt)
	}

	expectedSpace := ChainSpace{1, 26}
	gotSpace := g.Players[0].CurrentSpace
	if expectedSpace != gotSpace {
		t.Errorf("Expected space: %#v, got: %#v",
			expectedSpace,
			gotSpace,
		)
	}
}

func TestBuyStar(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 20, ChainSpace{0, 18}),
			NewPlayer("Luigi", 0, 4, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
		Config: GameConfig{
			MaxTurns: 25,
		},
	}

	g.MovePlayer(0, 1) //Land on blue space

	expectedSpace := ChainSpace{0, 20}
	expectedCoins := 3
	expectedStars := 1
	gotSpace := g.Players[0].CurrentSpace
	gotCoins := g.Players[0].Coins
	gotStars := g.Players[0].Stars

	if expectedSpace != gotSpace {
		t.Errorf("Expected Space: %#v, got: %#v", expectedSpace, gotSpace)
	}
	if expectedCoins != gotCoins {
		t.Errorf("Expected Coins: %#v, got: %#v", expectedCoins, gotCoins)
	}
	if expectedStars != gotStars {
		t.Errorf("Expected Stars: %#v, got: %#v", expectedStars, gotStars)
	}
}
