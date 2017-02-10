package game

type Player struct {
	chatId int64
	name string // only for debug purposes
	world *World
}

func (thisPlayer *Player) SetChatId(chatId int64) {
	thisPlayer.chatId = chatId
}

func (thisPlayer *Player) ChatId() int64 {
	return thisPlayer.chatId
}

func (thisPlayer *Player) SetWorld(world *World) {
	thisPlayer.world = world
}

func (thisPlayer *Player) World() *World {
	return thisPlayer.world
}

func (thisPlayer *Player) SetName(name string) {
	thisPlayer.name = name
}

func (thisPlayer *Player) Name() string {
	return thisPlayer.name
}
