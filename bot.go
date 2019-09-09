package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/bwmarrin/discordgo"
)

var conv = regexp.MustCompile(`(?i)^!conv ((?:[+-])?\d+(?:\.(\d+))?)\s*([[:alpha:]]+)\s+to\s+([[:alpha:]]+)`)

func main() {
	token, err := getBotToken()
	if err != nil {
		log.Fatalln("error getting bot token,", err)
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
	const (
		secretName = "UnitBot"
	)

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
		return "", nil
	}

	return secret.UnitBotToken, nil
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	match := conv.FindStringSubmatch(m.Content)
	if match == nil {
		return
	}

	ch := m.ChannelID

	precisionFrom := len(match[2])
	// ignore error because the regex ensures it will always parse
	num, _ := strconv.ParseFloat(match[1], 64)

	unitFrom, err := parseUnit(match[3])
	if err != nil {
		sendMessage(s, ch, err.Error())
		return
	}

	unitTo, err := parseUnit(match[4])
	if err != nil {
		sendMessage(s, ch, err.Error())
		return
	}

	converted, err := convert(num, unitFrom, unitTo)
	if err != nil {
		sendMessage(s, ch, err.Error())
		return
	}

	// precisionTo := calculatePrecision(precisionFrom, converted)

	// matchedMessage := fmt.Sprintf("Matched: %#v", match)
	// sendMessage(s, ch, matchedMessage)

	send := fmt.Sprintf("%.*f %s = %.6g %s",
		precisionFrom, num, unitFrom.name(), converted, unitTo.name())
	sendMessage(s, ch, send)
}

func sendMessage(s *discordgo.Session, channelID, message string) {
	_, err := s.ChannelMessageSend(channelID, message)
	if err != nil {
		log.Println("Unable to send message", err)
	}
}

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
