package game

import (
	//"log"
	//"fmt"
	//"strings"
	//"encoding/json"
)

type Game struct {
	playerM *Player
	playerW *Player
	blackboard Blackboard
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
