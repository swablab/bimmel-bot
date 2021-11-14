package discord

import (
	"log"
	"swablab-bot/config"
	"swablab-bot/handler"
	"time"

	"github.com/bwmarrin/discordgo"
)

var componentBellId string = "swablab-enty-bell"

type discordBot struct {
	messageHandler handler.MessageHandler
	config         *config.DiscordConfig
	activeServers  map[string]*discordServer
	discordSession *discordgo.Session
}

type discordServer struct {
	ChannelID       string
	CreatedChannel  bool
	CategoryID      string
	CreatedCategory bool
}

func NewRingBot(config *config.DiscordConfig, messageHandler handler.MessageHandler) (*discordBot, error) {
	bot := new(discordBot)
	bot.config = config
	bot.messageHandler = messageHandler
	bot.activeServers = make(map[string]*discordServer)

	discordSession, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return nil, err
	}
	bot.discordSession = discordSession

	bot.discordSession.AddHandler(bot.onMessageCreate)
	bot.discordSession.AddHandler(bot.onGuildCreate)
	bot.discordSession.AddHandler(bot.onInteractionCreate)

	err = bot.discordSession.Open()
	if err != nil {
		return nil, err
	}

	return bot, nil
}

func (bot *discordBot) Close() {
	bot.cleanupServers()
}

func (bot *discordBot) onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	server, exists := bot.activeServers[m.GuildID]
	if exists && m.ChannelID == server.ChannelID {
		go bot.messageHandler.SendMessage(m.Content)
		go func() {
			time.Sleep(2 * time.Second)
			s.ChannelMessageDelete(m.ChannelID, m.ID)
		}()
	}
}

func (bot *discordBot) onGuildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	if event.Guild.Unavailable {
		return
	}
	server, err := bot.prepareGuild(event.Guild)
	if err != nil {
		log.Print("could not prepare discord server", err)
		return
	}

	ringMessage := createRingMessage()
	_, err = s.ChannelMessageSendComplex(server.ChannelID, ringMessage)
	if err != nil {
		log.Print(err)
	}
}

func (bot *discordBot) onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionMessageComponent && i.MessageComponentData().CustomID == componentBellId {
		bot.messageHandler.SendMessage(i.Interaction.Member.User.Username)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
		})
	}
}

func (bot *discordBot) prepareGuild(guild *discordgo.Guild) (*discordServer, error) {
	server := new(discordServer)
	for _, channel := range guild.Channels {
		if channel.Name == config.DiscordConfiguration.ChannelName && channel.Type == discordgo.ChannelTypeGuildText {
			server.ChannelID = channel.ID
			server.CreatedChannel = false

		} else if channel.Name == config.DiscordConfiguration.ChannelCategory && channel.Type == discordgo.ChannelTypeGuildCategory {
			server.CategoryID = channel.ID
			server.CreatedCategory = false
		}
	}

	if server.CategoryID == "" {
		category, err := bot.discordSession.GuildChannelCreate(guild.ID, bot.config.ChannelCategory, discordgo.ChannelTypeGuildCategory)
		if err != nil {
			return nil, err
		}
		server.CategoryID = category.ID
		server.CreatedCategory = true
	}

	if server.ChannelID == "" {
		channel, err := bot.discordSession.GuildChannelCreate(guild.ID, bot.config.ChannelName, discordgo.ChannelTypeGuildText)
		if err == nil {
			return nil, err
		}
		server.ChannelID = channel.ID
		server.CreatedChannel = true

		edit := new(discordgo.ChannelEdit)
		edit.ParentID = server.CategoryID
		bot.discordSession.ChannelEditComplex(channel.ID, edit)
	}
	bot.activeServers[guild.ID] = server
	return server, nil
}

func (bot *discordBot) cleanupServers() {
	for _, server := range bot.activeServers {
		if server.CreatedChannel {
			bot.discordSession.ChannelDelete(server.ChannelID)

		}
		if server.CreatedCategory {
			bot.discordSession.ChannelDelete(server.CategoryID)
		}
	}
}

func createRingMessage() *discordgo.MessageSend {
	ringMessage := new(discordgo.MessageSend)
	ringMessage.Content = "Bitte hier klingeln"
	ringMessage.Components = []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Klingeln",
					Style:    discordgo.SuccessButton,
					Disabled: false,
					CustomID: componentBellId,
				},
			},
		},
	}
	return ringMessage
}
