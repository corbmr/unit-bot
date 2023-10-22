//go:build linux

package main

import (
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	convert "unit-bot"

	"github.com/bwmarrin/discordgo"
)

var applicationId, commandGuildId string

func init() {
	applicationId = os.Getenv("UNIT_BOT_APPLICATION_ID")
	commandGuildId = os.Getenv("UNIT_BOT_COMMAND_GUILD_ID")
}

func main() {
	convert.SetCurrencyApiKey(os.Getenv("CURRENCY_API_KEY"))

	discordToken, ok := os.LookupEnv("UNIT_BOT_TOKEN")
	if ok {
		stopDiscord := startDiscord(discordToken)
		defer stopDiscord()
	} else {
		slog.Warn("Discord token not found, skipping")
	}

	// Wait here until CTRL-C or other term signal is received.
	slog.Info("Unit Bot is now running")
	slog.Info("Press CTRL-C to stop")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGTERM)
	<-sc
	slog.Info("Stopping Unit Bot")
}

const convertPrefix = "!conv "

func startDiscord(discordToken string) func() {
	discordClient, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		slog.Error("error creating Discord session:", err)
		panic("error creating Discord session")
	}

	discordClient.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages
	discordClient.LogLevel = discordgo.LogWarning
	discordClient.AddHandler(processMessage)

	if len(commandGuildId) > 0 {
		slog.Info("creating commands in guild", "Guild", commandGuildId)
	} else {
		slog.Info("creating global commands")
	}

	createCommand(discordClient, &discordgo.ApplicationCommand{
		Name:        "convert",
		Description: "converts your units",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "from-value",
				Description: "value and unit to convert from",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "to-unit",
				Description: "unit to convert to",
				Required:    true,
			},
		},
	}, handleConvertInteraction)

	discordClient.AddHandler(func(discord *discordgo.Session, i *discordgo.InteractionCreate) {
		cmd := i.ApplicationCommandData().Name
		if handler, ok := commandHandlerMap[cmd]; ok {
			handler(discord, i)
		} else {
			slog.Warn("unknown command:", cmd)
		}
	})

	slog.Info("connecting to Discord...")
	if err := discordClient.Open(); err != nil {
		slog.Error("error connecting to Discord:", err)
		panic("error connecting to Discord")
	}
	slog.Info("successfully connected to Discord")
	return func() {
		if err := cleanupCommands(discordClient); err != nil {
			slog.Error("error cleaning up commands", err)
		} else {
			slog.Info("successfully cleaned up commands")
		}

		if err := discordClient.Close(); err != nil {
			slog.Info("error disconnecting from Discord:", err)
		} else {
			slog.Info("successfully disconnected from Discord")
		}
	}
}

var commandHandlerMap = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}

func createCommand(
	discord *discordgo.Session,
	command *discordgo.ApplicationCommand,
	handler func(s *discordgo.Session, i *discordgo.InteractionCreate),
) {
	cmd, err := discord.ApplicationCommandCreate(applicationId, commandGuildId, command)
	if err != nil {
		slog.Warn("error creating command:", command.Name, err)
		return
	}

	slog.Info("application command created", "Name", cmd.Name, "Id", cmd.ID)
	commandHandlerMap[command.Name] = handler
}

func cleanupCommands(discord *discordgo.Session) error {
	commands, err := discord.ApplicationCommands(applicationId, commandGuildId)
	if err != nil {
		return err
	}

	var errorList []error
	for _, command := range commands {
		err = discord.ApplicationCommandDelete(applicationId, commandGuildId, command.ID)
		errorList = append(errorList, err)
	}
	return errors.Join(errorList...)
}

func handleConvertInteraction(discord *discordgo.Session, i *discordgo.InteractionCreate) {
	var fromValue, toUnit string
	for _, o := range i.ApplicationCommandData().Options {
		switch o.Name {
		case "from-value":
			fromValue = o.StringValue()
		case "to-unit":
			toUnit = o.StringValue()
		default:
			slog.Warn("unexpected command option:", o.Name)
		}
	}

	convertResult := convert.Convert(fromValue, toUnit)

	discord.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: convertResult,
		},
	})
}

func processMessage(discord *discordgo.Session, m *discordgo.MessageCreate) {
	// Just in case
	defer func() {
		if err := recover(); err != nil {
			slog.Error("Function panicked:", err)
		}
	}()

	if m.Author.Bot {
		return
	}

	message := m.Content
	if !strings.HasPrefix(message, convertPrefix) {
		return
	}

	reply := convert.Process(strings.TrimPrefix(message, convertPrefix))

	_, err := discord.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: reply,
		Reference: &discordgo.MessageReference{
			MessageID: m.ID,
			ChannelID: m.ChannelID,
		},
		AllowedMentions: &discordgo.MessageAllowedMentions{},
	})
	if err != nil {
		slog.Info("Unable to send message", err)
	}
}
