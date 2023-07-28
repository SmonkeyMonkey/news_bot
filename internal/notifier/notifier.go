package notifier

import (
	"context"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/smonkeymonkey/news-bot/internal/model"
)

type Notifier struct {
	articles     ArticleProvider
	bot          *tgbotapi.BotAPI
	sendInterval time.Duration
	lookupTime   time.Duration
	channelId    int64
}
type ArticleProvider interface {
	AllNotPosted(ctx context.Context, since time.Time, limit uint64) ([]model.Article, error)
	MarkAsPosted(ctx context.Context, article model.Article) error
}

func New(articleP ArticleProvider,
	bot *tgbotapi.BotAPI,
	sendInterval time.Duration,
	lookup time.Duration,
	channelId int64) *Notifier {
	return &Notifier{
		articles:     articleP,
		bot:          bot,
		sendInterval: sendInterval,
		lookupTime:   lookup,
		channelId:    channelId,
	}
}

func (n *Notifier) Start(ctx context.Context) error {
	ticker := time.NewTicker(n.sendInterval)
	defer ticker.Stop()

	if err := n.SelectAndSendArticle(ctx); err != nil {
		return err
	}

	for {
		select {
		case <-ticker.C:
			if err := n.SelectAndSendArticle(ctx); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (n *Notifier) SelectAndSendArticle(ctx context.Context) error {
	topArticle, err := n.articles.AllNotPosted(ctx, time.Now().Add(-n.lookupTime), 1)
	if err != nil {
		return err
	}
	if len(topArticle) == 0 {
		return nil
	}

	article := topArticle[0]
	if err := n.sendArticle(article); err != nil {
		return err
	}

	return n.articles.MarkAsPosted(ctx, article)
}

func (n *Notifier) sendArticle(article model.Article) error {

	msg := tgbotapi.NewMessage(n.channelId, fmt.Sprintf("*%s*\n%s", n.EscapeForMarkdown(article.Title), n.EscapeForMarkdown(article.Link)))

	msg.ParseMode = tgbotapi.ModeMarkdownV2

	_, err := n.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func (n *Notifier) EscapeForMarkdown(src string) string {
	replacer := strings.NewReplacer(
		"-",
		"\\-",
		"_",
		"\\_",
		"*",
		"\\*",
		"[",
		"\\[",
		"]",
		"\\]",
		"(",
		"\\(",
		")",
		"\\)",
		"~",
		"\\~",
		"`",
		"\\`",
		">",
		"\\>",
		"#",
		"\\#",
		"+",
		"\\+",
		"=",
		"\\=",
		"|",
		"\\|",
		"{",
		"\\{",
		"}",
		"\\}",
		".",
		"\\.",
		"!",
		"\\!",
	)

	return replacer.Replace(src)
}
