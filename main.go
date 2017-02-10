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

type messages struct {
	messageToSender string
	messageToOpponent string
}

type freePlayers struct {
	female []*game.Player
	male []*game.Player
}

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

func sendMessages(bot *tgbotapi.BotAPI, sender *game.Player, opponent *game.Player, messages *messages) {
	if sender != nil && sender.ChatId() != 0 && messages.messageToSender != "" {
		msg := tgbotapi.NewMessage(sender.ChatId(), messages.messageToSender)
		bot.Send(msg)
	}
	if opponent != nil && opponent.ChatId() != 0 && messages.messageToOpponent != "" {
		msg := tgbotapi.NewMessage(opponent.ChatId(), messages.messageToOpponent)
		bot.Send(msg)
	}
}

func chooseGender(update *tgbotapi.Update, bot *tgbotapi.BotAPI, player *game.Player) (succeeded bool) {
	messages := &messages{}
	if (update.Message.Text == "F") {
		player.SetGender(game.Female)
		messages.messageToSender = "So you're a woman"
		succeeded = true
	} else if (update.Message.Text == "M") {
		player.SetGender(game.Male)
		messages.messageToSender = "So you're a man"
		succeeded = true
	} else {
		messages.messageToSender = "Pleace select your gender (M/F). It can't be changed."
		succeeded = false
	}

	sendMessages(bot, player, nil, messages)
	return
}

func matchPlayer(bot *tgbotapi.BotAPI, freePlayers *freePlayers, player *game.Player) {
	matchedPlayer := findAMatch(player, freePlayers)
	
	// send messages to players if an opponent is found
	if matchedPlayer != nil {
		messages := &messages {
			messageToSender : "match " + matchedPlayer.Name(),
			messageToOpponent : "match " + player.Name(),
		}

		sendMessages(bot, player, matchedPlayer, messages)
	} else {
		messages := &messages {
			messageToSender : "Searching for players",
		}

		sendMessages(bot, player, nil, messages)
	}
}

func processUpdate(update *tgbotapi.Update, bot *tgbotapi.BotAPI, players *map[int64]*game.Player,
				   staticData *game.StaticData, freePlayers *freePlayers) {
	
	player := makeOrFindPlayer(update.Message, players)
	
	if player.Gender() == game.Undefined { // if we need to choose a gender
		isSucceeded := chooseGender(update, bot, player)
		if isSucceeded {
			matchPlayer(bot, freePlayers, player)
		}
	} else if player.World() == nil { // if we need a match
		matchPlayer(bot, freePlayers, player)
	} else {
		var opponent *game.Player
		{
			world := player.World()
			if world.PlayerM() == player {
				opponent = world.PlayerF()
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
	
	freePlayers := &freePlayers{
		female : make([]*game.Player, 0),
		male : make([]*game.Player, 0),
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		
		processUpdate(&update, bot, &players, &staticData, freePlayers)
	}
}
