package monitoring

import (
	"context"
	"easy-reminder/internal/pkg/sender"
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
	availableMessage   = "корзин"
)

type Monitoring struct {
	Name        string
	Url         string
	PhoneNumber string
	CheckDelay  time.Duration
	Senders     []sender.Sender

	prevStatus Status
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
		log.Printf("%s:\t\t %s: %v\n", m.Name, statusUnknown, err)
		m.prevStatus = statusUnknown
		return nil
	}

	if response == nil || response.StatusCode != http.StatusOK {
		log.Printf("%s:\t\t %s\n", m.Name, statusUnknown)
		m.prevStatus = statusUnknown
		return nil
	}

	buf := new(strings.Builder)
	if _, err = io.Copy(buf, response.Body); err != nil {
		log.Printf("%s:\t\t %s: %v\n", m.Name, statusUnknown, err)
		m.prevStatus = statusUnknown
		return nil
	}

	// TODO: получше парсить
	if strings.Contains(buf.String(), unavailableMessage) {
		log.Printf("%s:\t\t %s\n", m.Name, statusInactive)
		m.prevStatus = statusInactive
		return nil
	}

	if !strings.Contains(buf.String(), availableMessage) {
		log.Printf("%s:\t\t %s\n", m.Name, statusUnknown)
		m.prevStatus = statusUnknown
		return nil
	}

	log.Printf("%s:\t\t %s (%s)\n", m.Name, statusActive, m.Url)

	if m.prevStatus == statusActive {
		return nil
	}
	m.prevStatus = statusActive
	for _, sndr := range m.Senders {
		message, sendErr := sndr.Send(&sender.Message{
			ID:   m.PhoneNumber,
			Text: fmt.Sprintf("Тут %s стал доступен в %s! \n Ссыл очка: %s", m.Name, time.Now().Format(time.DateTime), m.Url),
		})
		if sendErr != nil {
			log.Printf("Error sending message: %v\n", sendErr)
			return nil
		}

		jsonMessage, _ := json.MarshalIndent(message, "", "\t")
		log.Printf("Message sent: %v\n", string(jsonMessage))
	}

	return nil
}
