package mp1

type Player struct {
	Char         string
	Stars        int
	Coins        int
	CurrentSpace ChainSpace
	SkipTurn     bool
}
