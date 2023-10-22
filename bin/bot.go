//go:build linux

package main

import (
	"log"
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
		log.Println("Discord token not found, skipping")
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Unit Bot is now running")
	log.Println("Press CTRL-C to stop")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGTERM)
	<-sc
	log.Println("Stopping Unit Bot")
}

const convertPrefix = "!conv "

func startDiscord(discordToken string) func() {
	discordClient, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Panicln("error creating Discord session:", err)
	}

	discordClient.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages
	discordClient.LogLevel = discordgo.LogWarning
	discordClient.AddHandler(processMessage)

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
			log.Println("unknown command:", cmd)
		}
	})

	log.Println("connecting to Discord...")
	if err := discordClient.Open(); err != nil {
		log.Panicln("error connecting to Discord:", err)
	}
	log.Println("successfully connected to Discord")
	return func() {
		if err := discordClient.Close(); err != nil {
			log.Println("error disconnecting from Discord:", err)
		} else {
			log.Println("successfully disconnected from Discord")
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
		log.Println("error creating command:", command.Name, err)
		return
	}

	log.Printf("application command created: %v id: %v\n", cmd.Name, cmd.ID)
	commandHandlerMap[command.Name] = handler
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
			log.Println("unexpected command option:", o.Name)
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
			log.Println("Function panicked:", err)
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
		log.Println("Unable to send message", err)
	}
}
