package main

import (
	"log"
	"fmt"
	"strings"
	"io/ioutil"
	"encoding/json"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gameraccoon/telegram-date-game-bot/game"
)

func getFileStringContent(filePath string) string {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Print(err)
	}

	return strings.TrimSpace(string(fileContent))
}

func getApiToken() string {
	return getFileStringContent("./telegramApiToken.txt")
}

func remove(s []int, i int) []int {
    s[len(s)-1], s[i] = s[i], s[len(s)-1]
    return s[:len(s)-1]
}

func getStaticDataJsonString() string {
	return getFileStringContent("./staticData.json")
}

func main() {
	var apiToken string = getApiToken()

	bot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	
	players := make(map[int64]*game.Player)
	
	var staticData game.StaticData
	
	{
		dec := json.NewDecoder(strings.NewReader(getStaticDataJsonString()))
		err := dec.Decode(&staticData);
		if err != nil {
			log.Fatal(err)
		}
	}
	
	// only one because there are no matching by sex yet
	var freePlayer *game.Player;

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		
		// find already created player
		player, ok := players[update.Message.Chat.ID]
		
		// create if not found
		if !ok {
			player = &game.Player{}
			
			player.SetChatId(update.Message.Chat.ID)
			
			user := update.Message.From
			if user != nil {
				if len(user.UserName) > 0 {
					player.SetName(user.UserName)
				} else {
					player.SetName(user.FirstName)
				}
			}
			
			players[update.Message.Chat.ID] = player
		}
		
		{	
			message := "action: " + staticData.Actions["TestActionId"].Name
			msg := tgbotapi.NewMessage(player.ChatId(), message)
			bot.Send(msg)
		}
		
		// if we need a match
		if player.Game() == nil {
			// if we are first
			if freePlayer == nil {
				// put ourselves into match queue
				freePlayer = player
			} else {
				// if there are another player to match
				if player != freePlayer {
					// make game
					game := &game.Game{}
					game.SetPlayerM(player)
					game.SetPlayerW(freePlayer)
					freePlayer.SetGame(game)
					player.SetGame(game)

					// remove the player from queue
					freePlayer = nil

					// send messages to players
					{
						message := "match " + game.PlayerM().Name()
						msg := tgbotapi.NewMessage(game.PlayerW().ChatId(), message)
						bot.Send(msg)
					}
					{
						message := "match " + game.PlayerW().Name()
						msg := tgbotapi.NewMessage(game.PlayerM().ChatId(), message)
						bot.Send(msg)
					}
				}
			}
		} else {
			game := player.Game()
			gamePlayer := game.PlayerM()
			if gamePlayer == player {
				gamePlayer = game.PlayerW()
			}
			
			if update.Message.Text == "/disconnect" {
				player.SetGame(nil)
				gamePlayer.SetGame(nil)
				
				{
					message := "disconnected"
					msg := tgbotapi.NewMessage(player.ChatId(), message)
					bot.Send(msg)
				}
				{
					message := "player left"
					msg := tgbotapi.NewMessage(gamePlayer.ChatId(), message)
					bot.Send(msg)
				}
				
				continue
			}
			
			if gamePlayer != nil {
				msg := tgbotapi.NewMessage(gamePlayer.ChatId(), "> " + update.Message.Text)
				bot.Send(msg)
			}
		}
	}
}

