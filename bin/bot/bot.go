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

func main() {
	convert.CurrencyInit = func() (string, error) {
		return os.Getenv("CURRENCY_API_KEY"), nil
	}

	token := os.Getenv("UNIT_BOT_TOKEN")
	if len(token) == 0 {
		log.Fatalln("Discord token not found")
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln("error creating discord session,", err)
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages
	dg.LogLevel = discordgo.LogWarning
	dg.AddHandler(onMessageCreate)

	err = dg.Open()
	if err != nil {
		log.Fatalln("error opening connection,", err)
	}
	defer func() {
		dg.Close()
		log.Println("Unit Bot has stopped running")
	}()

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Unit Bot is now running")
	log.Println("Press CTRL-C to stop")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGTERM)
	<-sc
}

const (
	convertPrefix = "!conv "
)

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Just in case
	defer func() {
		if err := recover(); err != nil {
			log.Println("Function panicked:", err)
		}
	}()

	if strings.HasPrefix(m.Content, convertPrefix) {
		res := convert.Process(strings.TrimPrefix(m.Content, convertPrefix))
		err := reply(s, m.Message, res)
		if err != nil {
			log.Println("Unable to send message", err)
		}
	}
}

func reply(s *discordgo.Session, m *discordgo.Message, content string) error {
	_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: content,
		Reference: &discordgo.MessageReference{
			MessageID: m.ID,
			ChannelID: m.ChannelID,
		},
		AllowedMentions: &discordgo.MessageAllowedMentions{},
	})
	return err
}
