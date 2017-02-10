package game

import (
	"github.com/mattn/anko/vm"
)

type World struct {
	playerM *Player
	playerW *Player
	env *vm.Env
}

func (thisWorld *World) Init() {
	thisWorld.env = vm.NewEnv()
}

func (thisWorld *World) Execute(command string) {
	_, err := thisWorld.env.Execute(command)
	if err != nil {
		panic(err)
	}
}

func (thisWorld *World) SetPlayerM(player *Player) {
	thisWorld.playerM = player
}

func (thisWorld *World) SetPlayerW(player *Player) {
	thisWorld.playerW = player
}

func (thisWorld *World) PlayerM() *Player {
	return thisWorld.playerM
}

func (thisWorld *World) PlayerW() *Player {
	return thisWorld.playerW
}
