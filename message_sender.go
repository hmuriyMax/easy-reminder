package main

import (
	"errors"
	"fmt"
	"github.com/smsaero/smsaero_golang"
)

const (
	email  = "mshchemilkin@gmail.com"
	apiKey = "wmNFt7rppo_atUjw4HOrYkJz4Vx3W2iR"
	signer = "SMSAERO"
)

type Message struct {
	PhoneNumber string
	Text        string
}

func SendMessage(message Message) (*smsaero_golang.Send, error) {
	client := smsaero_golang.NewSmsAeroClient(email, apiKey)

	balance, err := client.Balance()
	if err != nil {
		return nil, fmt.Errorf("error getting balance: %w", err)
	}

	if balance <= 0 {
		return nil, errors.New("insufficient balance")
	}

	sendResult, err := client.Send(message.PhoneNumber, message.Text, signer)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}
	return &sendResult, nil
}
