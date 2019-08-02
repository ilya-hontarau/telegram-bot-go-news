package main

import (
	"log"
	"os"
	"time"

	"github.com/alecthomas/kingpin"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/illfate/telegram-bot-go-news/pkg/bot"
	"github.com/illfate/telegram-bot-go-news/pkg/cache"
)

const botToken = "BOT_TOKEN"

var (
	lifeTime   = kingpin.Flag("lifetime", "Life time of posts").Default("24h").Duration()
	updateTime = kingpin.Flag("updtime", "Posts update time").Default("1h").Duration()
	scrapeURL  = kingpin.Flag("url", "Scrape rss channel").Default("https://habr.com/ru/rss/hubs/all/").String()
)

func main() {
	kingpin.Parse()
	token := os.Getenv(botToken)
	if token == "" {
		log.Fatalf(`no %q env var`, botToken)
	}
	tgBot, err := bot.New(token)
	if err != nil {
		log.Fatalf("couldn't start bot: %s", err)
	}
	u := tgbotapi.UpdateConfig{
		Timeout: 60,
	}
	updates, err := tgBot.GetUpdatesChan(u)
	if err != nil {
		log.Printf("couldn't get update chan: %s", err)
		return
	}

	c := cache.New()
	c.ScrapePosts(*scrapeURL)

	ticker := time.NewTicker(*updateTime)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			c.UpdatePosts(*lifeTime, *scrapeURL)
		}
	}()

	for update := range updates {
		switch update.Message.Command() {
		case "next":
			tgBot.NextCommand(update, c)
		case "start":
			tgBot.StartCommand(update)
		}
	}
}
