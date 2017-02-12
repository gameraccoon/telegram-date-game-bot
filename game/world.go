package game

import (
	"github.com/mattn/anko/vm"
)

type World struct {
	playerF *Player
	playerM *Player
	env *vm.Env
	currentPlayer *Player
	opponent *Player
}

func (thisWorld *World) Init() {
	thisWorld.env = vm.NewEnv()
	// create blackboard
	thisWorld.env.Execute("b = {}; playerF = {}; playerM = {}")
}

func (thisWorld *World) IsTrue(command string) bool {
	result, err := thisWorld.env.Execute(command)
	if err != nil {
		panic(err)
	}
	
	return result.Bool()
}

func (thisWorld *World) Execute(command string) {
	_, err := thisWorld.env.Execute(command)
	if err != nil {
		panic(err)
	}
}

func (thisWorld *World) SetPlayerF(player *Player) {
	thisWorld.playerF = player
}

func (thisWorld *World) SetPlayerM(player *Player) {
	thisWorld.playerM = player
}

func (thisWorld *World) PlayerF() *Player {
	return thisWorld.playerF
}

func (thisWorld *World) PlayerM() *Player {
	return thisWorld.playerM
}

func (thisWorld *World) CurrentPlayer() *Player {
	return thisWorld.currentPlayer
}

func (thisWorld *World) Opponent() *Player {
	return thisWorld.opponent
}

func (thisWorld *World) InitTurn(forcedPlayer *Player) {
	if (forcedPlayer == thisWorld.playerF) {
		thisWorld.currentPlayer = thisWorld.playerF
		thisWorld.opponent = thisWorld.playerM
		thisWorld.env.Execute("me = playerF; opponent = playerM")
	} else if (forcedPlayer == thisWorld.playerM) {
		thisWorld.currentPlayer = thisWorld.playerM
		thisWorld.opponent = thisWorld.playerF
		thisWorld.env.Execute("me = playerM; opponent = playerF")
	} else {
		panic("unknown player")
	}
}

func (thisWorld *World) ChangeTurn() {
	if thisWorld.currentPlayer == nil || thisWorld.opponent == nil {
		panic("players haven't initialized")
	} else {
		temp := thisWorld.currentPlayer
		thisWorld.currentPlayer = thisWorld.opponent
		thisWorld.opponent = temp
		thisWorld.env.Execute("temp = me; me = opponent; opponent = temp; temp = nil")
	}
}
