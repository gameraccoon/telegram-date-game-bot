package main

import "github.com/gameraccoon/telegram-date-game-bot/game"

func findAMatch(player *game.Player, freePlayer **game.Player) (match *game.Player) {
	// if we are first
	if *freePlayer == nil {
		// put ourselves into match queue
		*freePlayer = player
	} else {
		// if there are another player to match
		if player != *freePlayer {
			match = *freePlayer
			// remove the player from queue
			*freePlayer = nil
			
			// make game
			world := &game.World{}
			world.SetPlayerM(player)
			world.SetPlayerW(match)
			match.SetWorld(world)
			player.SetWorld(world)
		}
	}
	
	return
}