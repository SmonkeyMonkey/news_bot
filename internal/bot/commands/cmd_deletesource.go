package bot_commands

import (
	"context"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/smonkeymonkey/news-bot/internal/bot"
)

type SourceDeleter interface {
	Delete(ctx context.Context, sourceID int64) error
}

func ViewCmdDeleteSource(deleter SourceDeleter) bot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		idStr := update.Message.CommandArguments()
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return err
		}

		if err := deleter.Delete(ctx, id); err != nil {
			log.Printf("error deleting: %v", err)
			return nil
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Source succesfully deleted")
		if _, err := bot.Send(msg); err != nil {
			return err
		}

		return nil
	}
}
