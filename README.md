
news-bot is a telegram bot for getting news via rss

## Installation

Clone repository
```bash
https://github.com/smonkeymonkey/news-bot
```

## Usage
news-bot can publish news in your Telegram channel via rss, for this you need to register a new bot in @BotFather and type your token into config.hcl file into telegram_bot_token

#### Available env variable
- NB_TELEGRAM_BOT_TOKEN — token for Telegram Bot API 
- NB_TELEGRAM_CHANNEL_ID — ID of the channel to post to, can be obtained via @JsonDumpBot
- NB_DATABASE_DSN — PostgreSQL connection string
- NB_FETCH_INTERVAL — the interval of checking for new articles, default 10m
- NB_NOTIFICATION_INTERVAL — the interval of delivering new articles to Telegram channel, default 1m
- NB_FILTER_KEYWORDS — comma separated list of words to skip articles containing these words


### Run with Docker

```bash
cd news-bot
docker-compose up
```
This repository uses goose for manage migrations
First need install goose:
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```
For macOS users goose is available as a Homebrew Formulae:
```bash
brew install goose
``` 
after installation you can up migrations
```bash
make up
```

### example usage
##### add new source: 
- /addsource {"name": "reddit", "url": "https://www.reddit.com/r/golang/.rss"} 
##### list sources: 
- /listsources 
##### delete source: 
- /deletesource 1
