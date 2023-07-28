package bot_commands

import (
	"context"
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	news_bot "github.com/smonkeymonkey/news-bot/internal/bot"
	"github.com/smonkeymonkey/news-bot/internal/bot/helpers"
	"github.com/smonkeymonkey/news-bot/internal/model"
)

type SourceStorage interface {
	Add(ctx context.Context, source model.Source) (int64, error)
}

func VievCmdAddSource(storage SourceStorage) news_bot.ViewFunc {
	type addSourceArgs struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		args, err := helpers.ParseJson[addSourceArgs](update.Message.CommandArguments())
		if err != nil {
			return errors.New("incorrect input data")
		}

		source := model.Source{
			Name:    args.Name,
			FeedURL: args.URL,
		}
		sourceId, err := storage.Add(ctx, source)
		if err != nil {
			return err
		}

		msg := fmt.Sprintf("Source with id `%d` succesfully added", sourceId)
		reply := tgbotapi.NewMessage(update.Message.Chat.ID, msg)

		reply.ParseMode = tgbotapi.ModeMarkdownV2

		if _, err := bot.Send(reply); err != nil {
			return err
		}
		return nil
	}
}
