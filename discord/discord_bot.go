package discord

import (
	"log"
	"swablab-bot/config"
	"swablab-bot/handler"
	"time"

	"github.com/bwmarrin/discordgo"
)

type discordBot struct {
	messageHandler handler.MessageHandler
	config         *config.DiscordConfig
	activeServers  map[string]*DiscordServer
	discordSession *discordgo.Session
}

type DiscordServer struct {
	channelId       string
	createdChannel  bool
	categoryId      string
	createdCategory bool
}

func NewRingBot(config *config.DiscordConfig, messageHandler handler.MessageHandler) (*discordBot, error) {
	bot := new(discordBot)
	bot.config = config
	bot.messageHandler = messageHandler
	bot.activeServers = make(map[string]*DiscordServer)

	discordSession, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return nil, err
	}
	bot.discordSession = discordSession

	bot.discordSession.AddHandler(bot.onMessageCreate)
	bot.discordSession.AddHandler(bot.OnGuildCreate)

	err = bot.discordSession.Open()
	if err != nil {
		return nil, err
	}

	return bot, nil
}

func (bot *discordBot) onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	go bot.messageHandler.SendMessage(m.Content)
	go func() {
		time.Sleep(2 * time.Second)
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}()

}

func (bot *discordBot) OnGuildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	if event.Guild.Unavailable {
		return
	}

	server := new(DiscordServer)
	for _, channel := range event.Guild.Channels {
		if channel.Name == config.DiscordConfiguration.ChannelName && channel.Type == discordgo.ChannelTypeGuildText {
			server.channelId = channel.ID
			server.createdChannel = false

		} else if channel.Name == config.DiscordConfiguration.ChannelCategory && channel.Type == discordgo.ChannelTypeGuildCategory {
			server.categoryId = channel.ID
			server.createdCategory = false
		}
	}

	if server.categoryId == "" {
		category, err := s.GuildChannelCreate(event.Guild.ID, bot.config.ChannelCategory, discordgo.ChannelTypeGuildCategory)
		if err == nil {
			server.categoryId = category.ID
			server.createdCategory = true
		}
	}

	if server.channelId == "" {
		channel, err := s.GuildChannelCreate(event.Guild.ID, bot.config.ChannelName, discordgo.ChannelTypeGuildText)
		if err == nil {
			server.channelId = channel.ID
			server.createdChannel = true
		}

		edit := new(discordgo.ChannelEdit)
		edit.ParentID = server.categoryId
		s.ChannelEditComplex(channel.ID, edit)
	}

	bot.activeServers[event.Guild.ID] = server

	msg := new(discordgo.MessageSend)
	msg.Content = "Bitte hier klingeln (einfach eine Nachricht schreiben)"

	_, err := s.ChannelMessageSendComplex(server.channelId, msg)
	if err != nil {
		log.Fatal(err)
	}
}

func (bot *discordBot) cleanupChannels() {
	for _, server := range bot.activeServers {
		if server.createdChannel {
			bot.discordSession.ChannelDelete(server.channelId)

		}
		if server.createdCategory {
			bot.discordSession.ChannelDelete(server.categoryId)
		}
	}
}

func (bot *discordBot) Close() {
	bot.cleanupChannels()
}
