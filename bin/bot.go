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

var applicationId, guildId string

func init() {
	applicationId = getenvOrPanic("UNIT_BOT_APPLICATION_ID")
	guildId = getenvOrPanic("UNIT_BOT_COMMAND_GUILD_ID")
}

func getenvOrPanic(env string) string {
	v := os.Getenv(env)
	if len(v) == 0 {
		log.Panicln("environment variable expected:", env)
	}
	return v
}

func main() {
	convert.CurrencyInit = func() (string, error) {
		return os.Getenv("CURRENCY_API_KEY"), nil
	}

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
	discordClient.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.Bot {
			return
		}
		processMessage(m.Content, func(reply string) error {
			_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content: reply,
				Reference: &discordgo.MessageReference{
					MessageID: m.ID,
					ChannelID: m.ChannelID,
				},
				AllowedMentions: &discordgo.MessageAllowedMentions{},
			})
			return err
		})
	})

	createCommand(discordClient, &discordgo.ApplicationCommand{
		Name:        "convert",
		Description: "converts your units",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionNumber,
				Name:        "value",
				Description: "value of unit to convert",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "from-unit",
				Description: "unit to convert from",
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

	createCommand(discordClient, &discordgo.ApplicationCommand{
		Name:        "convert2",
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
	}, handleConvert2Interaction)

	discordClient.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		cmd := i.ApplicationCommandData().Name
		if handler, ok := commandHandlerMap[cmd]; ok {
			handler(s, i)
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
	s *discordgo.Session,
	command *discordgo.ApplicationCommand,
	handler func(s *discordgo.Session, i *discordgo.InteractionCreate),
) {
	cmd, err := s.ApplicationCommandCreate(applicationId, guildId, command)
	if err != nil {
		log.Println("error creating command:", command.Name)
		return
	}

	log.Printf("application command created: %v id: %v\n", cmd.Name, cmd.ID)
	commandHandlerMap[command.Name] = handler
}

func handleConvertInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var (
		value            float64
		fromUnit, toUnit string
	)
	for _, o := range i.ApplicationCommandData().Options {
		switch o.Name {
		case "value":
			value = o.FloatValue()
		case "from-unit":
			fromUnit = o.StringValue()
		case "to-unit":
			toUnit = o.StringValue()
		default:
			log.Println("unexpected command option:", o.Name)
		}
	}

	convertResult := convert.Convert(value, fromUnit, toUnit)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: convertResult,
		},
	})
}

func handleConvert2Interaction(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	convertResult := convert.Convert2(fromValue, toUnit)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: convertResult,
		},
	})
}

func processMessage(message string, reply func(string) error) {
	// Just in case
	defer func() {
		if err := recover(); err != nil {
			log.Println("Function panicked:", err)
		}
	}()

	if strings.HasPrefix(message, convertPrefix) {
		res := convert.Process(strings.TrimPrefix(message, convertPrefix))
		err := reply(res)
		if err != nil {
			log.Println("Unable to send message", err)
		}
	}
}
