package mp1

//HiddenBlockEvent holds the implementation for hidden blocks.
type HiddenBlockEvent struct {
	Player int
}

//Responses returns a slice of bools (true/false).
func (h HiddenBlockEvent) Responses() []Response {
	return []Response{true, false}
}

func (h HiddenBlockEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle sets the hidden block action to be taken depending on r. If r is
//true, then the next event will be a EventDiceBlock. If r is false, then
//then the player will land on the space they're currently on.
func (h HiddenBlockEvent) Handle(r Response, g *Game) {
	isHiddenBlock := r.(bool)
	if isHiddenBlock {
		g.ExtraEvent = EventDiceBlock{h.Player}
	} else {
		g.ActivateSpace(h.Player)
	}
}
