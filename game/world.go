package game

import (
	"github.com/mattn/anko/vm"
)

type World struct {
	playerF *Player
	playerM *Player
	env *vm.Env
}

func (thisWorld *World) Init() {
	thisWorld.env = vm.NewEnv()
	// create blackboard
	thisWorld.env.Execute("b = {}")
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
