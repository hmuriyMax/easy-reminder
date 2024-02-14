package sms

import (
	"easy-reminder/internal/pkg/identity/user_providers"
	"easy-reminder/internal/pkg/sender"
	"errors"
	"fmt"
	"github.com/smsaero/smsaero_golang"
)

const signer = "SMSAERO"

type Sender struct {
	client *smsaero_golang.Client
	user_providers.Provider
}

func NewSender(client *smsaero_golang.Client) *Sender {
	return &Sender{client: client}
}

func (s *Sender) Send(message *sender.Message) (*sender.Result, error) {
	balance, err := s.client.Balance()
	if err != nil {
		return nil, fmt.Errorf("error getting balance: %w", err)
	}

	if balance <= 0 {
		return nil, errors.New("insufficient balance")
	}

	sendResult, err := s.client.Send(message.ID, message.Text, signer)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}
	return &sender.Result{
		Message: sendResult.Text,
		Status:  sender.SendStatusSuccess,
	}, nil
}
