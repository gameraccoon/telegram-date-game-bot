package main

import (
	"log"
	"fmt"
	"bytes"
	"strings"
	"text/template"
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
		
		(*players)[message.Chat.ID] = player
	}
	
	return player
}

func formatMessage(message *string, player *game.Player, opponent *game.Player) string {
	data := map[string]interface{}{}
	
	if opponent != nil {
		data["Name"] = opponent.Name()
	}
	
	if player != nil {
		data["YourName"] = player.Name()
	}

	t := template.Must(template.New("").Parse(*message))
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, data); err != nil {
		panic(err)
	}
	return buf.String()
}

func sendMessages(bot *tgbotapi.BotAPI, sender *game.Player, opponent *game.Player, messages *messages) {
	if sender != nil && sender.ChatId() != 0 && messages.messageToSender != "" {
		msg := tgbotapi.NewMessage(sender.ChatId(), formatMessage(&messages.messageToSender, sender, opponent))
		bot.Send(msg)
	}
	if opponent != nil && opponent.ChatId() != 0 && messages.messageToOpponent != "" {
		msg := tgbotapi.NewMessage(opponent.ChatId(), formatMessage(&messages.messageToOpponent, opponent, sender))
		bot.Send(msg)
	}
}

func chooseGender(update *tgbotapi.Update, bot *tgbotapi.BotAPI, player *game.Player) (succeeded bool) {
	
	if (update.Message.Text == "/female") {
		player.SetGender(game.Female)
		succeeded = true
	} else if (update.Message.Text == "/male") {
		player.SetGender(game.Male)
		succeeded = true
	} else {
		messages := &messages{
			messageToSender : "Select your gender: /male /female. It can't be changed.",
		}
		sendMessages(bot, player, nil, messages)
		
		succeeded = false
	}
	return
}

func chooseName(update *tgbotapi.Update, bot *tgbotapi.BotAPI, player *game.Player, staticData *game.StaticData) (succeeded bool) {
	succeeded = false
	var namesList *[]string
	if player.Gender() == game.Female {
		namesList = &staticData.NamesF
	} else if player.Gender() == game.Male {
		namesList = &staticData.NamesM
	} else {
		panic("gender hasn't set")
	}
	
	if strings.HasPrefix(update.Message.Text, "/") {
		receivedName := update.Message.Text[1:len(update.Message.Text)]
		for _, name := range *namesList {
			if name == receivedName {
				player.SetName(receivedName)
				succeeded = true
				break
			}
		}
	}
	
	if !succeeded {
		var buffer bytes.Buffer
		buffer.WriteString("Select your name:\n")
		
		for _, name := range *namesList {
			buffer.WriteString("/" + name + " ")
		}
		
		messages := &messages{
			messageToSender : buffer.String(),
		}
		sendMessages(bot, player, nil, messages)
		
		succeeded = false
	}
	return
}

func genderRelatedText(player *game.Player, femaleText string, maleText string) string {
	if player.Gender() == game.Female {
		return femaleText
	} else {
		return maleText
	}
}

func matchPlayer(bot *tgbotapi.BotAPI, freePlayers *freePlayers, player *game.Player) {
	matchedPlayer := findAMatch(player, freePlayers)
	
	// send messages to players if an opponent is found
	if matchedPlayer != nil {		
		summaryText := fmt.Sprintf("You came to your friend. %s must've been waiting you for a long time.\nThis is %s door\n/act - see the full list of available actions",
								   genderRelatedText(player, "He", "She"),
								   genderRelatedText(player, "his", "her"))
	
		messages := &messages {
			messageToSender : summaryText,
		}

		sendMessages(bot, player, nil, messages)
	} else {
		messages := &messages {
			messageToSender : fmt.Sprintf("You're waiting for your friend's coming. %s said that %s'll come very soon.",
										  genderRelatedText(player, "He", "She"),
										  genderRelatedText(player, "he", "she")),
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
			chooseName(update, bot, player, staticData)
		}
	} else if player.Name() == "" { // if we need to choose name
		isSucceeded := chooseName(update, bot, player, staticData)
		
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
		
		messages := processCommand(player, opponent, staticData, &update.Message.Text)
		
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
