package main

import "github.com/gameraccoon/telegram-date-game-bot/game"

func matchPlayers(femalePlayer *game.Player, malePlayer *game.Player) {
	// make a new world for these players
	world := &game.World{}
	world.Init()
	world.SetPlayerF(femalePlayer)
	world.SetPlayerM(malePlayer)
	femalePlayer.SetWorld(world)
	malePlayer.SetWorld(world)
}

func findAMatch(player *game.Player, freePlayers *freePlayers) (opponent *game.Player) {
	if player.Gender() == game.Female {
		if len(freePlayers.male) > 0 {
			opponent = freePlayers.male[0]
			freePlayers.male = freePlayers.male[1:]
			
			matchPlayers(player, opponent)
		} else {
			freePlayers.female = append(freePlayers.female, player)
		}
	} else if player.Gender() == game.Male {
		if len(freePlayers.female) > 0 {
			opponent = freePlayers.female[0]
			freePlayers.female = freePlayers.female[1:]
			
			matchPlayers(player, opponent)
		} else {
			freePlayers.male = append(freePlayers.male, player)
		}
	} else {
		panic("Gender is not set")
	}
	
	return
}