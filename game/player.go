package game

import (
	//"log"
	//"fmt"
	//"strings"
)

type Player struct {
	chatId int64
	name string // for debug purposes
	game *Game
}

func (thisPlayer *Player) SetChatId(chatId int64) {
	thisPlayer.chatId = chatId
}

func (thisPlayer *Player) ChatId() int64 {
	return thisPlayer.chatId
}

func (thisPlayer *Player) SetGame(game *Game) {
	thisPlayer.game = game
}

func (thisPlayer *Player) Game() *Game {
	return thisPlayer.game
}

func (thisPlayer *Player) SetName(name string) {
	thisPlayer.name = name
}

func (thisPlayer *Player) Name() string {
	return thisPlayer.name
}
