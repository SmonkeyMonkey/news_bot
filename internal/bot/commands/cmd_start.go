package bot_commands

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	news_bot "github.com/smonkeymonkey/news-bot/internal/bot"
)

func ViewStart() news_bot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "started")); err != nil {
			return err
		}
		return nil
	}
}
