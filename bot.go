package main

import (
	"encoding/json"
	"log"
	"math"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/bwmarrin/discordgo"
)

func main() {
	var err error
	token, ok := os.LookupEnv("UNITTEST")
	if !ok {
		token, err = getBotToken()
		if err != nil {
			log.Fatalln("error getting bot token,", err)
		}
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

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Unit Bot is now running")
	log.Println("Press CTRL-C to stop")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
	log.Println("Unit Bot has stopped running")
}

func getBotToken() (string, error) {
	const secretName = "UnitBot"

	//Create a Secrets Manager client
	sess, err := session.NewSession()
	if err != nil {
		return "", err
	}

	svc := secretsmanager.New(sess)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		return "", err
	}

	var secret struct{ UnitBotToken string }

	err = json.Unmarshal([]byte(*result.SecretString), &secret)
	if err != nil {
		return "", err
	}

	return secret.UnitBotToken, nil
}

var prefix = regexp.MustCompile(`(?i)^!conv(?:ert)?\s*`)

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Just in case
	defer func() {
		if err := recover(); err != nil {
			log.Println("Function panicked:", err)
		}
	}()

	loc := prefix.FindStringIndex(m.Content)
	if loc == nil {
		return
	}

	out, err := generateResponse(m.Content[loc[1]:])
	if err != nil {
		out = err.Error()
	}

	_, err = s.ChannelMessageSend(m.ChannelID, out)
	if err != nil {
		log.Println("Unable to send message", err)
	}
}

// TODO: Fix this and change out any uses of %g in formats
func calculatePrecision(givenPrecision int, num float64) int {

	num = math.Abs(num)
	precisionTo := 0

	const epsilon = 1e-9
	if _, f := math.Modf(num); f > epsilon && f < 1-epsilon {
		precisionTo += 2
	}

	if precisionTo < givenPrecision {
		precisionTo = givenPrecision
	}

	if log10 := int(math.Log10(num)); log10 < 0 {
		precisionTo -= log10
	}

	return precisionTo
}
