package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	Token     string
	BotPrefix string
	config    *configStruct
)

type configStruct struct {
	BotPrefix string `json:"botPrefix"`
}

func ReadConfig() error {
	fmt.Println("Reading from config.json")
	file, err := os.ReadFile("./config.json")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(string(file))
	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading ENV: %v\n", err)
		return err
	}
	Token = os.Getenv("TOKEN")
	BotPrefix = config.BotPrefix
	return nil
}
