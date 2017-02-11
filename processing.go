package main

import (
	"bytes"
	"strings"
	"github.com/gameraccoon/telegram-date-game-bot/game"
)

func collectActions(player *game.Player, staticData *game.StaticData) (actions map[string]*game.Action) {
	actions = make(map[string]*game.Action)
	
	world := player.World()
	
	for id, action := range staticData.Actions {
		if world.IsTrue(action.Requirements) {
			actions[id] = action
		}
	}
	
	return
}

// resultMessages[0] is a message for message sender and resultMessages[1] is a message for his opponent
func processCommand(player *game.Player, opponent *game.Player, staticData *game.StaticData, message *string) (resultMessages *messages) {
	resultMessages = &messages{}
	
	if *message == "/disconnect" {
		player.SetWorld(nil)
		opponent.SetWorld(nil)
		
		resultMessages.messageToSender = "disconnected"
		resultMessages.messageToOpponent = "player left"
	} else if *message == "/act" {
		actions := collectActions(player, staticData)
		var buffer bytes.Buffer
		for id, action := range actions {
			buffer.WriteString("/act_" + id + " " + action.Name + "\n")
		}
		resultMessages.messageToSender = buffer.String()
	} else if strings.HasPrefix(*message, "/act_") {
		actionId := (*message)[5:len(*message)]
		action := staticData.Actions[actionId]
		if action != nil {
			world := player.World()
			if world.IsTrue(action.Requirements) {
				world.Execute(action.Action)

				if player.Gender() == game.Female {
					resultMessages.messageToSender = action.TextInstigatorF
					resultMessages.messageToOpponent = action.TextReceiverF
				} else if player.Gender() == game.Male {
					resultMessages.messageToSender = action.TextInstigatorM
					resultMessages.messageToOpponent = action.TextReceiverM
				}
			} else {
				resultMessages.messageToSender = "You can't use this action"
			}
		}
	} else {
		resultMessages.messageToOpponent = "> " + *message
	}
	
	return
}
