package p2p

type GameStatus uint32

func (gs GameStatus) String() string {
	switch gs {
	case GameStatusWaiting:
		return "WAITING"
	case GameStatusDealing:
		return "DEALING"
	case GameStatusPreFlop:
		return "PRE_FLOP"
	case GameStatusFlop:
		return "FLOP"
	case GameStatusTurn:
		return "TURN"
	case GameStatusRiver:
		return "RIVER"
	default:
		return "unknown"
	}
}

const (
	GameStatusWaiting GameStatus = iota
	GameStatusDealing
	GameStatusPreFlop
	GameStatusFlop
	GameStatusTurn
	GameStatusRiver
)

type GameState struct {
	isDeal     bool       // should be atomic accessible
	gameStatus GameStatus // should be atomic accessible
}

func NewGameState() *GameState {
	return &GameState{}
}

func (gs *GameState) loop() {
	for {
		select {}
	}
}
