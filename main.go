package main

import (
	"context"
	"time"
)

type Status string

const (
	phoneNumber = "+79213918575"

	parkingURL      = "https://tickets.stoyanie.ru/shop/parkovka-p1-4/"
	parkingFastRack = "https://tickets.stoyanie.ru/shop/parkovka-fast-track-2/"
	transferURL     = "https://tickets.stoyanie.ru/shop/transfer-tuda-obratno-otpravlenie-iz-moskvy-v-800-16-marta-2024/"

	checkDelay          = time.Second * 5
	stopMessageDuration = time.Minute * 10
)

func main() {
	var (
		parking = Monitoring{
			Name:                "Parking",
			Url:                 parkingURL,
			PhoneNumber:         phoneNumber,
			CheckDelay:          checkDelay,
			StopMessageDuration: stopMessageDuration,
		}

		parkingFT = Monitoring{
			Name:                "Parking Fast Track",
			Url:                 parkingFastRack,
			PhoneNumber:         phoneNumber,
			CheckDelay:          checkDelay,
			StopMessageDuration: stopMessageDuration,
		}

		transfer = Monitoring{
			Name:                "Transfer",
			Url:                 transferURL,
			PhoneNumber:         phoneNumber,
			CheckDelay:          checkDelay,
			StopMessageDuration: stopMessageDuration,
		}
	)

	ctx := context.Background()

	parking.RunAsync(ctx)
	parkingFT.RunAsync(ctx)
	transfer.RunAsync(ctx)

	select {
	case <-ctx.Done():
		return
	}
}
