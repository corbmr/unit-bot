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
	"github.com/gempir/go-twitch-irc/v3"
)

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

	twitchToken, ok := os.LookupEnv("TWITCH_TOKEN")
	if ok {
		stopTwitch := startTwitch(twitchToken)
		defer stopTwitch()
	} else {
		log.Println("Twitch token not found, skipping")
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

func startTwitch(twitchToken string) func() {
	twitchClient := twitch.NewClient("UnitBot", "oauth:"+twitchToken)
	twitchClient.SetJoinRateLimiter(twitch.CreateVerifiedRateLimiter())
	twitchClient.OnPrivateMessage(func(message twitch.PrivateMessage) {
		processMessage(message.Message, func(reply string) error {
			twitchClient.Reply(message.Channel, message.ID, reply)
			return nil
		})
	})

	log.Println("connecting to Twitch...")
	if err := connectTwitch(twitchClient); err != nil {
		log.Panicln("error connecting to Twitch,", err)
	}
	log.Println("successfully connected to Twitch")

	joinTwitchChannels(twitchClient)

	return func() {
		if err := twitchClient.Disconnect(); err != nil {
			log.Println("error disconnecting from Twitch:", err)
		} else {
			log.Println("successfully disconnected from Twitch")
		}
	}
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

func connectTwitch(t *twitch.Client) error {
	twitchConnected := make(chan error)
	t.OnConnect(func() {
		twitchConnected <- nil
	})

	go func() {
		twitchConnected <- t.Connect()
	}()

	return <-twitchConnected
}

func joinTwitchChannels(t *twitch.Client) {
	f, err := os.ReadFile("channels.txt")
	if err != nil {
		log.Println("Unable to read channels.txt:", err)
		return
	}

	t.Join(strings.Split(string(f), "\n")...)
}
