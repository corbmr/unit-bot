package main

import (
	"encoding/json"
	"log"
	"math"
	"os"
	"os/signal"
	"syscall"

	"unit-bot/internal/convert"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/bwmarrin/discordgo"
)

func main() {
	secret, err := getSecret()
	if err != nil {
		log.Fatalln("unable to get secrets:", err)
	}
	convert.InitCurrency(secret.CurrencyAPIKey)

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
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
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

	token, tokenOk := os.LookupEnv("UNITTEST")
	apiKey, apiOk := os.LookupEnv("CURRENCY_KEY")

	if !tokenOk || !apiOk {
		secret, err := getBotSecret()
		if err != nil {
			return secrets{}, err
		}

		if !tokenOk {
			token = secret.UnitBotToken
		}

		if !apiOk {
			apiKey = secret.CurrencyAPIKey
		}
	}

	return secrets{token, apiKey}, nil
}

func getBotSecret() (secrets, error) {
	const secretName = "UnitBot"

	//Create a Secrets Manager client
	sess, err := session.NewSession()
	if err != nil {
		return secrets{}, err
	}

	svc := secretsmanager.New(sess)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		return secrets{}, err
	}

	var secret secrets

	err = json.Unmarshal([]byte(*result.SecretString), &secret)
	return secret, err
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Just in case
	defer func() {
		if err := recover(); err != nil {
			log.Println("Function panicked:", err)
		}
	}()

	res := commandExpr([]byte(m.Content))
	if !res.Ok {
		return
	}

	resp, err := res.V.(command).Do()
	if err != nil {
		resp = err.Error()
	}

	_, err = s.ChannelMessageSend(m.ChannelID, resp)
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
