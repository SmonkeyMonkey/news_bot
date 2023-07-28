package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"


	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"


	"github.com/smonkeymonkey/news-bot/internal/bot"
	bot_commands "github.com/smonkeymonkey/news-bot/internal/bot/commands"
	"github.com/smonkeymonkey/news-bot/internal/config"
	"github.com/smonkeymonkey/news-bot/internal/fetcher"
	"github.com/smonkeymonkey/news-bot/internal/notifier"
	"github.com/smonkeymonkey/news-bot/internal/storage"
)

func main() {
	botApi, err := tgbotapi.NewBotAPI(config.Get().TelegramBotToken)
	if err != nil {
		log.Printf("failed to create bot.Error: %v", err)
		return
	}

	
	db, err := sqlx.Connect("postgres", config.Get().DatabaseDSN)
	if err != nil {
		log.Printf("failed to connect to db.Error: %v", err)
		return
	}
	defer db.Close()

	var (
		articleStorage = storage.NewArticleStorage(db)
		sourceStorage  = storage.NewSourceStorage(db)
		fetcher        = fetcher.New(
			articleStorage,
			sourceStorage,
			config.Get().FetchInterval,
			config.Get().FilterKeywords,
		)
		notifier = notifier.New(
			articleStorage,
			botApi,
			config.Get().FetchInterval,
			2*config.Get().FetchInterval,
			config.Get().TelegramChannelID,
		)
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	newsBot := bot.New(botApi)
	newsBot.RegisterCmdView("start", bot_commands.ViewStart())
	newsBot.RegisterCmdView("addsource", bot_commands.VievCmdAddSource(sourceStorage))
	newsBot.RegisterCmdView("listsources", bot_commands.ViewCmdListSources(sourceStorage))
	newsBot.RegisterCmdView("deletesource", bot_commands.ViewCmdDeleteSource(sourceStorage))

	go func(ctx context.Context) {
		if err := fetcher.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("failed to start fetcher.Error: %v", err)
				return
			}
			log.Println("fetcher stopped")
		}
	}(ctx)

	go func(ctx context.Context) {
		if err := notifier.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("failed to start notifier.Error: %v", err)
				return
			}
			log.Println("notifier stopped")
		}
	}(ctx)

	if err := newsBot.Run(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			log.Printf("failed to run bot.Error msg: %v", err)
			return
		}
		log.Println("bot stopped")
	}
}
