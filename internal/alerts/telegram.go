package alerts

import (
	"fmt"
	"strconv"

	telebot "gopkg.in/telebot.v3"
)

type TelegramBot struct {
	bot       *telebot.Bot
	chatIDInt int64
}

func NewBot(token, chatID string) *TelegramBot {
	b, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Verbose: false,
	})
	if err != nil {
		fmt.Printf("Failed to create Telegram bot: %v\n", err)
		return nil
	}

	chatIDInt, _ := strconv.ParseInt(chatID, 10, 64)

	return &TelegramBot{bot: b, chatIDInt: chatIDInt}
}

func (t *TelegramBot) SendMessage(text string) error {
	if t == nil || t.bot == nil {
		return fmt.Errorf("bot not initialized")
	}

	_, err := t.bot.Send(telebot.ChatID(t.chatIDInt), text)
	return err
}