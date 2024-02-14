package main

import (
	"context"
	"easy-reminder/internal/pkg/monitoring"
	"easy-reminder/internal/pkg/sender"
	"easy-reminder/internal/pkg/sender/tg"
	"time"
)

const (
	phoneNumber = ""

	email    = ""
	apiKey   = ""
	tgBotKey = ""

	parkingURL      = "https://tickets.stoyanie.ru/shop/parkovka-p1-4/"
	parkingFastRack = "https://tickets.stoyanie.ru/shop/parkovka-fast-track-2/"
	transferURL     = "https://tickets.stoyanie.ru/shop/transfer-tuda-obratno-otpravlenie-iz-moskvy-v-800-16-marta-2024/"

	checkDelay = time.Second * 5
)

func main() {
	ctx := context.Background()

	var (
		//smsSender = sms.NewSender(smsaero_golang.NewSmsAeroClient(email, apiKey))
		tgSender = tg.NewSender(ctx, tgBotKey)
		parking  = monitoring.Monitoring{
			Name:        "Parking",
			Url:         parkingURL,
			PhoneNumber: phoneNumber,
			CheckDelay:  checkDelay,
			Senders:     []sender.Sender{tgSender},
		}

		parkingFT = monitoring.Monitoring{
			Name:        "Parking FT",
			Url:         parkingFastRack,
			PhoneNumber: phoneNumber,
			CheckDelay:  checkDelay,
			Senders:     []sender.Sender{tgSender},
		}

		//transfer = monitoring.Monitoring{
		//	Name:                "Transfer",
		//	Url:                 transferURL,
		//	PhoneNumber:         phoneNumber,
		//	CheckDelay:          checkDelay,
		//	StopMessageDuration: stopMessageDuration,
		//}
	)

	parking.RunAsync(ctx)
	parkingFT.RunAsync(ctx)
	//transfer.RunAsync(ctx)

	select {
	case <-ctx.Done():
		return
	}
}
