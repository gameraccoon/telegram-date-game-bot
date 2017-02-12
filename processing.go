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
	world := player.World()
	
	if *message == "/disconnect" {
		player.SetWorld(nil)
		opponent.SetWorld(nil)
		
		resultMessages.messageToSender = "disconnected"
		resultMessages.messageToOpponent = "player left"
	} else if player == world.CurrentPlayer() {
		if *message == "/act" {
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
				if world.IsTrue(action.Requirements) {
					world.Execute(action.Action)
					
					if player.Gender() == game.Female {
						resultMessages.messageToSender = action.TextInstigatorF
						resultMessages.messageToOpponent = action.TextReceiverF
					} else if player.Gender() == game.Male {
						if action.TextInstigatorM != nil {
							resultMessages.messageToSender = *action.TextInstigatorM
						} else {
							resultMessages.messageToSender = action.TextInstigatorF
						}
						
						if action.TextReceiverM != nil {
							resultMessages.messageToOpponent = *action.TextReceiverM
						} else {
							resultMessages.messageToOpponent = action.TextReceiverF
						}
					}
					
					world.ChangeTurn()
					
					var buffer bytes.Buffer
					
					buffer.WriteString(resultMessages.messageToOpponent + "\n\n")
					for _, id := range action.Reactions {
						reaction := staticData.Actions[id]
						if reaction != nil {
							if world.IsTrue(reaction.Requirements) {
								buffer.WriteString("/act_" + id + " - " + reaction.Name + "\n")
							}
						} else {
							panic("unknown acton " + id)
						}
					}
					buffer.WriteString("/act - full list of available actions")

					resultMessages.messageToOpponent = buffer.String()
				} else {
					resultMessages.messageToSender = "You can't use this action"
				}
			}
		}
	} else {
		resultMessages.messageToSender = "It's not your turn"
	}
	
	return
}
