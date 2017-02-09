package game

import (
	"github.com/mattn/anko/vm"
)

type Game struct {
	playerM *Player
	playerW *Player
	env *vm.Env
}

func (thisGame *Game) Init() {
	thisGame.env = vm.NewEnv()
}

func (thisGame *Game) Execute(command string) {
	_, err := thisGame.env.Execute(command)
	if err != nil {
		panic(err)
	}
}

func (thisGame *Game) SetPlayerM(player *Player) {
	thisGame.playerM = player
}

func (thisGame *Game) SetPlayerW(player *Player) {
	thisGame.playerW = player
}

func (thisGame *Game) PlayerM() *Player {
	return thisGame.playerM
}

func (thisGame *Game) PlayerW() *Player {
	return thisGame.playerW
}
