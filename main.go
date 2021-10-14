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

var myHandler handler.MessageHandler

func main() {
	var err error
	myHandler, err = handler.NewMqttHandler(config.MqttConfiguration)
	if err != nil {
		log.Fatal(err)
	}
	defer myHandler.Close()

	dg, err := discordgo.New("Bot " + config.DiscordConfiguration.Token)
	if err != nil {
		log.Fatal(err)
	}
	defer dg.Close()
	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		log.Fatalf("error opening connection %s", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	go myHandler.Message(m.Content)
}
