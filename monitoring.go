package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Status string

const (
	statusActive   Status = "AVAILABLE"
	statusInactive Status = "OUT OF STOCK"
	statusUnknown  Status = "UNKNOWN"

	unavailableMessage = "Нет в наличии"
)

type Monitoring struct {
	Name                string
	Url                 string
	PhoneNumber         string
	CheckDelay          time.Duration
	StopMessageDuration time.Duration

	messageSendSleepStop time.Time
}

func (m *Monitoring) RunAsync(ctx context.Context) {
	go func() {
		err := m.Run(ctx)
		if err != nil {
			log.Printf("Error running: %v\n", err)
		}
	}()
}

func (m *Monitoring) Run(ctx context.Context) error {
	m.messageSendSleepStop = time.Now()
	for {
		select {
		case <-time.After(m.CheckDelay):
			err := m.checker()
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (m *Monitoring) checker() error {
	response, err := http.Get(m.Url)
	if err != nil {
		log.Printf("%s: %s: %v\n", m.Name, statusUnknown, err)
		return nil
	}

	if response == nil || response.StatusCode != http.StatusOK {
		log.Printf("%s: %s\n", m.Name, statusUnknown)
		return nil
	}

	buf := new(strings.Builder)
	if _, err = io.Copy(buf, response.Body); err != nil {
		log.Printf("%s: %s: %v\n", m.Name, statusUnknown, err)
		return nil
	}

	// TODO: получше парсить
	if strings.Contains(buf.String(), unavailableMessage) {
		log.Printf("%s: %s\n", m.Name, statusInactive)
		return nil
	}

	log.Printf("%s: %s (%s)\n", m.Name, statusActive, m.Url)

	if time.Now().Before(m.messageSendSleepStop) {
		return nil
	}
	message, sendErr := SendMessage(Message{
		PhoneNumber: m.PhoneNumber,
		Text:        fmt.Sprintf("Йо, тут \"%s\" доступен в %s! \n Ссыл очка: %s", m.Name, time.Now().Format(time.DateTime), m.Url),
	})
	if sendErr != nil {
		log.Printf("Error sending message: %v\n", sendErr)
		return nil
	}

	jsonMessage, _ := json.MarshalIndent(message, "", "\t")
	log.Printf("Message sent: %v\n", string(jsonMessage))
	m.messageSendSleepStop = time.Now().Add(m.StopMessageDuration)

	return nil
}
