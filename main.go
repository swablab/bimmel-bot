package main

import (
	"log"
	"os"
	"os/signal"
	"swablab-bot/config"
	"swablab-bot/discord"
	"swablab-bot/handler"
	"syscall"
)

var messageHandler handler.MessageHandler

func main() {
	var err error

	//build MessageHandler
	messageHandler, err = handler.NewDummyMessageHandler() //handler.NewMqttMessageHandler(config.MqttConfiguration)
	if err != nil {
		log.Fatal(err)
	}
	defer messageHandler.Close()
	log.Println("Successfully created MessageHandler")

	//build discord api
	ringBot, err := discord.NewRingBot(&config.DiscordConfiguration, messageHandler)
	if err != nil {
		log.Fatal(err)
	}
	defer ringBot.Close()
	log.Println("Successfully connected to the discord API")

	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
