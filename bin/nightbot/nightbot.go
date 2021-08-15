package main

import (
	"encoding/json"
	"errors"
	convert "unit-bot"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

const secretName = "UnitBotSecret"

func init() {
	convert.InitCurrency(func() (string, error) {
		session := session.Must(session.NewSession())
		secrets := secretsmanager.New(session)

		out, err := secrets.GetSecretValue(&secretsmanager.GetSecretValueInput{SecretId: aws.String(secretName)})
		if err != nil {
			return "", err
		}

		var s struct{ CurrencyAPIKey string }
		err = json.Unmarshal([]byte(*out.SecretString), &s)
		if err != nil {
			return "", err
		}

		return s.CurrencyAPIKey, nil
	})
}

func main() {
	lambda.Start(handler)
}

func handler(req events.APIGatewayV2HTTPRequest) (string, error) {
	query, ok := req.QueryStringParameters["q"]
	if !ok {
		return "", errors.New("Request missing query parameter")
	}
	return convert.Process(query), nil
}
