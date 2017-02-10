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

func getStaticDataJsonString() string {
	return getFileStringContent("./staticData.json")
}

func makeOrFindPlayer(message *tgbotapi.Message, players *map[int64]*game.Player) *game.Player {
	// find an already created player
	player, ok := (*players)[message.Chat.ID]
		
	// create if not found
	if !ok {
		player = &game.Player{}

		player.SetChatId(message.Chat.ID)

		user := message.From
		if user != nil {
			if len(user.UserName) > 0 {
				player.SetName(user.UserName)
			} else {
				player.SetName(user.FirstName)
			}
		}

		(*players)[message.Chat.ID] = player
	}
	
	return player
}

type messages struct {
	messageToSender string
	messageToOpponent string
}

func sendMessages(bot *tgbotapi.BotAPI, sender *game.Player, opponent *game.Player, messages *messages) {
	if messages.messageToSender != "" && sender.ChatId() != 0 {
		msg := tgbotapi.NewMessage(sender.ChatId(), messages.messageToSender)
		bot.Send(msg)
	}
	if messages.messageToOpponent != "" && opponent.ChatId() != 0 {
		msg := tgbotapi.NewMessage(opponent.ChatId(), messages.messageToOpponent)
		bot.Send(msg)
	}
}

func processUpdate(update *tgbotapi.Update, bot *tgbotapi.BotAPI, players *map[int64]*game.Player,
				   staticData *game.StaticData, freePlayer **game.Player) {
	
	player := makeOrFindPlayer(update.Message, players)

	// if we need a match
	if player.World() == nil {
		matchedPlayer := findAMatch(player, freePlayer)

		// send messages to players if an opponent is found
		if matchedPlayer != nil {
			messages := &messages {
				messageToSender : "match " + matchedPlayer.Name(),
				messageToOpponent : "match " + player.Name(),
			}
			
			sendMessages(bot, player, matchedPlayer, messages)
		}
	} else {
		
		// find opponent
		var opponent *game.Player
		{
			world := player.World()
			if world.PlayerM() == player {
				opponent = world.PlayerW()
			} else {
				opponent = world.PlayerM()
			}
		}
		
		messages := processCommand(player, opponent, &update.Message.Text)
		
		sendMessages(bot, player, opponent, messages)
	}
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
		
		processUpdate(&update, bot, &players, &staticData, &freePlayer)
	}
}
