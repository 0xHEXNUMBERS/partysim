package mp1

import "testing"

func TestHiddenBlock(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{0, 21}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
		Config: GameConfig{MaxTurns: 20, EventsDice: true},
	}
	g.MovePlayer(0, 1) //Blue Space
	EventIs(HiddenBlockEvent{0}, g.NextEvent, "", t)

	gHidden := g
	gHidden.NextEvent.Handle(true, &gHidden)
	EventIs(EventDiceBlock{0}, gHidden.NextEvent, "Dice", t)

	gBlue := g
	gBlue.NextEvent.Handle(false, &gBlue)
	EventIs(PickDiceBlock{1, gBlue.Config}, gBlue.NextEvent, "Player", t)
	CoinsIs(13, 0, gBlue, "", t)
}
