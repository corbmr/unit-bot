package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"unit-bot/internal/convert"
)

func main() {
	convert.InitCurrency(func() (string, error) {
		key, ok := os.LookupEnv("CURRENCY_KEY")
		if !ok {
			return "", errors.New("CURRENCY_KEY not found")
		}
		return key, nil
	})

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println(convert.Process(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
