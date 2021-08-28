package main

import (
	"log"
	"os"
	"os/signal"
	"strings"

	convert "unit-bot"

	"github.com/bwmarrin/discordgo"
)

func main() {
	convert.InitCurrency(func() (string, error) {
		return os.Getenv("CURRENCY_API_KEY"), nil
	})

	token := os.Getenv("UNIT_BOT_TOKEN")
	if len(token) == 0 {
		log.Fatalln("Discord token not found")
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln("error creating discord session,", err)
	}

	dg.AddHandler(onMessageCreate)

	err = dg.Open()
	if err != nil {
		log.Fatalln("error opening connection,", err)
	}
	defer dg.Close()

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Unit Bot is now running")
	log.Println("Press CTRL-C to stop")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc
	log.Println("Unit Bot has stopped running")
}

const cmdPrefix = "!conv "

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Just in case
	defer func() {
		if err := recover(); err != nil {
			log.Println("Function panicked:", err)
		}
	}()

	if !strings.HasPrefix(m.Content, cmdPrefix) {
		return
	}

	res := convert.Process(strings.TrimPrefix(m.Content, cmdPrefix))

	_, err := s.ChannelMessageSend(m.ChannelID, res)
	if err != nil {
		log.Println("Unable to send message", err)
	}
}
