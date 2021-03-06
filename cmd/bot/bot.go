package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"unit-bot/internal/convert"

	"github.com/bwmarrin/discordgo"
)

func main() {
	secret, err := getSecret()
	if err != nil {
		log.Fatalln("unable to get secrets:", err)
	}

	convert.InitCurrency(func() (string, error) {
		return secret.CurrencyAPIKey, nil
	})

	if len(secret.UnitBotToken) == 0 {
		log.Fatalln("Discord token not found")
	}

	dg, err := discordgo.New("Bot " + secret.UnitBotToken)
	if err != nil {
		log.Fatalln("error creating discord session,", err)
	}

	dg.AddHandler(onMessageCreate)

	err = dg.Open()
	if err != nil {
		log.Fatalln("error opening connection,", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Unit Bot is now running")
	log.Println("Press CTRL-C to stop")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
	log.Println("Unit Bot has stopped running")
}

type secrets struct {
	UnitBotToken   string
	CurrencyAPIKey string
}

func getSecret() (secrets, error) {
	s, ok := os.LookupEnv("UNIT_BOT_SECRET")
	if !ok {
		return secrets{}, errors.New("UNIT_BOT_SECRET not found")
	}

	var secret secrets
	err := json.Unmarshal([]byte(s), &secret)
	if err != nil {
		return secrets{}, err
	}

	return secret, nil
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Just in case
	defer func() {
		if err := recover(); err != nil {
			log.Println("Function panicked:", err)
		}
	}()

	const prefix = "!conv "
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	res := convert.Process(strings.TrimPrefix(m.Content, prefix))

	_, err := s.ChannelMessageSend(m.ChannelID, res)
	if err != nil {
		log.Println("Unable to send message", err)
	}
}
