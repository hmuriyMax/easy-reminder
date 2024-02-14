package tg

import (
	"context"
	"easy-reminder/internal/pkg/sender"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
)

type Sender struct {
	bot   *tgbotapi.BotAPI
	chats map[tgbotapi.Chat]struct{}
}

func NewSender(ctx context.Context, token string) *Sender {
	bot, _ := tgbotapi.NewBotAPI(token)
	sndr := &Sender{
		bot:   bot,
		chats: make(map[tgbotapi.Chat]struct{}),
	}
	go func() {
		err := sndr.listenAndSaveChats(ctx)
		if err != nil {
			log.Fatalln(err)
		}
	}()
	return sndr
}

func (s *Sender) Send(message *sender.Message) (sent *sender.Result, err error) {
	for chat, _ := range s.chats {
		var res tgbotapi.Message
		res, err = s.bot.Send(tgbotapi.NewMessage(chat.ID, message.Text))
		if err != nil {
			err = fmt.Errorf("failed to send message: %w", err)
			continue
		}
		sent = &sender.Result{
			Message: res.Text,
			Status:  sender.SendStatusSuccess,
		}
	}
	return
}

func (s *Sender) listenAndSaveChats(ctx context.Context) error {
	chn, err := s.bot.GetUpdatesChan(tgbotapi.NewUpdate(0))
	if err != nil {
		return fmt.Errorf("error getting updates: %w", err)
	}

	for {
		select {
		case event := <-chn:
			if event.Message == nil || event.Message.Chat == nil {
				continue
			}
			chat := event.Message.Chat
			s.chats[*chat] = struct{}{}
			_, _ = s.bot.Send(tgbotapi.NewMessage(chat.ID, "This chat was registered"))

		case <-ctx.Done():
			return nil
		}
	}
}
