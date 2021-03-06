package mp1

import (
	"reflect"
	"testing"
)

func MakeSimpleBoard(s SpaceType) Board {
	return Board{
		Chains: &[]Chain{
			{
				{Type: Start},
				{Type: s},
			},
		},
	}
}

func MakeRepeatedBoard(s SpaceType, times int) Board {
	chain := make(Chain, times+1)
	for i := range chain {
		chain[i] = Space{Type: s}
	}
	chain[0] = Space{Type: Start}
	return Board{Chains: &[]Chain{chain}}
}

func SpaceIs(expected ChainSpace, player int, g Game, flavour string, t *testing.T) {
	got := g.Players[player].CurrentSpace
	if expected != got {
		t.Errorf("Expected %s %d Space: %#v, got: %#v",
			flavour, player, expected, got)
	}
}

func StarsIs(expected, player int, g Game, flavour string, t *testing.T) {
	got := g.Players[player].Stars
	if expected != got {
		t.Errorf("Expected Player %d %s Stars: %d, got: %d",
			player, flavour, expected, got)
	}
}

func CoinsIs(expected, player int, g Game, flavour string, t *testing.T) {
	got := g.Players[player].Coins
	if expected != got {
		t.Errorf("Expected Player %d %s Coins: %d, got: %d",
			player, flavour, expected, got)
	}
}

func MinigameCoinsIs(expected, player int, g Game, flavour string, t *testing.T) {
	got := g.Players[player].MinigameCoins
	if expected != got {
		t.Errorf("Expected Player %d %s Coins: %d, got: %d",
			player, flavour, expected, got)
	}
}

func SpaceTypeIs(expected, got SpaceType, flavour string, t *testing.T) {
	if expected != got {
		t.Errorf("Expected %s Space Type: %d, got: %d",
			flavour, expected, got)
	}
}

func IntIs(expected, got int, flavour string, t *testing.T) {
	if expected != got {
		t.Errorf("Expected %s: %d, got: %d",
			flavour, expected, got)
	}
}

func EventIs(expected, got Event, flavour string, t *testing.T) {
	if expected != got {
		t.Errorf("Expected %s Event: %#v, got: %#v",
			flavour, expected, got)
	}
}

func ResIs(expected []Response, g Game, flavour string, t *testing.T) {
	got := g.NextEvent.Responses()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected %s Res: %#v, got: %#v",
			flavour, expected, got)
	}
}
