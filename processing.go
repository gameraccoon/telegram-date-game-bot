package main

import "github.com/gameraccoon/telegram-date-game-bot/game"

// resultMessages[0] is a message for message sender and resultMessages[1] is a message for his opponent
func processCommand(player *game.Player, opponent *game.Player, message *string) (resultMessages *messages) {
	resultMessages = &messages{}
	
	if *message == "/disconnect" {
		player.SetWorld(nil)
		opponent.SetWorld(nil)
		
		resultMessages.messageToSender = "disconnected"
		resultMessages.messageToOpponent = "player left"
	} else {
		resultMessages.messageToOpponent = "> " + *message
	}
	
	return
}
