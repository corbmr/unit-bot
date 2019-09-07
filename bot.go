package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/bwmarrin/discordgo"
)

const region = "us-west-2"

var conv = regexp.MustCompile(`^(\d+(?:\.(\d+))?)\s*([[:alpha:]]+)\s+to\s+([[:alpha:]]+)`)

func main() {
	token, err := getBotToken()
	if err != nil {
		log.Fatalln("error getting bot token,", err)
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln("error creating discord session,", err)
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		log.Fatalln("error opening connection,", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func getBotToken() (string, error) {
	const (
		secretName = "UnitBot"
	)

	//Create a Secrets Manager client
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if !strings.HasPrefix(m.Content, "!conv ") {
		return
	}

	match := conv.FindStringSubmatch(m.Content[6:])
	if match == nil {
		const usage = "Usage: !conv {from}{units} to {units}"
		sendMessage(s, m.ChannelID, usage)
		return
	}

	precision := len(match[2])
	// ignore error because the regex ensures it will always parse
	num, _ := strconv.ParseFloat(match[1], 64)
	unitFrom, unitTo := match[3], match[4]
	converted := convert(num, unitFrom, unitTo)

	matchedMessage := fmt.Sprintf("Matched: %#v", match)
	sendMessage(s, m.ChannelID, matchedMessage)

	send := fmt.Sprintf("%.*f %s = %*.f %s", precision, num, unitFrom, precision, converted, unitTo)
	sendMessage(s, m.ChannelID, send)
}

func sendMessage(s *discordgo.Session, channelID, message string) {
	_, err := s.ChannelMessageSend(channelID, message)
	if err != nil {
		log.Println("Unable to send message", message)
	}
}
