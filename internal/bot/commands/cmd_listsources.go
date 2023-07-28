package bot_commands

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
	news_bot "github.com/smonkeymonkey/news-bot/internal/bot"
	"github.com/smonkeymonkey/news-bot/internal/bot/helpers"
	"github.com/smonkeymonkey/news-bot/internal/model"
)

type SourceListner interface {
	Sources(cttx context.Context) ([]model.Source, error)
}

func ViewCmdListSources(lister SourceListner) news_bot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		sources, err := lister.Sources(ctx)
		if err != nil {
			return err
		}
		formetted := lo.Map(sources, func(source model.Source, _ int) string {
			return formatSource(source)
		})
		msg := fmt.Sprintf("Sources: \\(%d\\):\n\n%s", len(sources), strings.Join(formetted, "\n\n"))

		reply := tgbotapi.NewMessage(update.Message.Chat.ID, msg)
		reply.ParseMode = tgbotapi.ModeMarkdownV2

		if _, err := bot.Send(reply); err != nil {
			return err
		}
		return nil
	}
}

func formatSource(source model.Source) string {
	return fmt.Sprintf("*%s*\n ID:`%d`\n URL: %s",
		helpers.EscapeForMarkdown(source.Name),
		source.ID,
		helpers.EscapeForMarkdown(source.FeedURL),
	)
}
