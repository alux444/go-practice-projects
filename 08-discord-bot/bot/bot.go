package bot

import (
	"fmt"
	"ping-bot/config"

	"github.com/bwmarrin/discordgo"
)

var BotId string
var goBot *discordgo.Session

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	user, err := goBot.User("@me")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	BotId = user.ID
	goBot.AddHandler(messageHandler)

	err = goBot.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running.")
}

func messageHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == BotId {
		return
	}

	if message.Content == "ping" {
		session.ChannelMessageSend(message.ChannelID, "pong")
	}
}
