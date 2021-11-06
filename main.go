package main

import (
	"log"
	"os"
	"os/signal"
	"swablab-bot/config"
	"swablab-bot/handler"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var messageHandler handler.MessageHandler

func main() {
	var err error

	//build MessageHandler
	messageHandler, err = handler.NewMqttMessageHandler(config.MqttConfiguration)
	if err != nil {
		log.Fatal(err)
	}
	defer messageHandler.Close()
	log.Println("Successfully created MessageHandler")

	//build discord api
	discord, err := discordgo.New("Bot " + config.DiscordConfiguration.Token)
	if err != nil {
		log.Fatal(err)
	}
	defer discord.Close()

	discord.AddHandler(messageCreate)

	err = discord.Open()
	if err != nil {
		log.Fatalf("error opening connection %s", err)
		return
	}
	log.Println("Successfully connected to the discord API")

	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	go messageHandler.SendMessage(m.Content)
}
