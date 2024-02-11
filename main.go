package main

import (
	"encoding/json"
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
	url                = "https://tickets.stoyanie.ru/shop/bilet-na-masleniczu-2024/"
)

func main() {
	messageSendSleepStop := time.Now()
	for {
		time.Sleep(time.Second * 5)

		response, err := http.Get(url)
		if err != nil {
			log.Printf("%s: %v\n", statusUnknown, err)
			continue
		}

		if response == nil || response.StatusCode != http.StatusOK {
			log.Println(statusUnknown)
			continue
		}

		buf := new(strings.Builder)
		if _, err = io.Copy(buf, response.Body); err != nil {
			log.Printf("%s: %v\n", statusUnknown, err)
			continue
		}

		if strings.Contains(buf.String(), unavailableMessage) {
			log.Println(statusInactive)
			continue
		}

		log.Println(statusActive)

		if time.Now().Before(messageSendSleepStop) {
			continue
		}

		log.Println(url)
		message, sendErr := SendMessage(Message{
			PhoneNumber: "+79164352929",
			Text:        "Бегом за билетом! \n https://tickets.stoyanie.ru/shop/bilet-na-masleniczu-2024/",
		})
		if sendErr != nil {
			log.Printf("Error sending message: %v\n", sendErr)
			continue
		}

		jsonMessage, _ := json.MarshalIndent(message, "", "\t")
		log.Printf("Message sent: %v\n", string(jsonMessage))
		messageSendSleepStop = time.Now().Add(time.Minute * 10)
	}
}
