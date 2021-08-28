package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	convert "unit-bot"

	"github.com/chzyer/readline"
)

func main() {
	convert.InitCurrency(func() (string, error) {
		key, ok := os.LookupEnv("CURRENCY_API_KEY")
		if !ok {
			return "", errors.New("CURRENCY_API_KEY not found")
		}
		return key, nil
	})

	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	log.SetOutput(rl.Stderr())
	for {
		line, err := rl.Readline()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		fmt.Println(convert.Process(line))
	}
}
