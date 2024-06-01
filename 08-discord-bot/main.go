package main

import (
	"fmt"
	"ping-bot/bot"
	"ping-bot/config"
)

func main() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bot.Start()

	//creates a channel that can send and receive values of type struct
	<-make(chan struct{})
}
