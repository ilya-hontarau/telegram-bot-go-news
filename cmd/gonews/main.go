package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/illfate/telegram-bot-go-news/pkg/config"

	"github.com/alecthomas/kingpin"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/illfate/telegram-bot-go-news/pkg/bot"
	"github.com/illfate/telegram-bot-go-news/pkg/cache"
)

const botToken = "BOT_TOKEN"

var (
	lifeTime       = kingpin.Flag("lifetime", "Life time of posts").Default("24h").Duration()
	updateTime     = kingpin.Flag("updtime", "Posts update time").Default("1h").Duration()
	scrapeURL      = kingpin.Flag("url", "Scrape rss channel").Default("https://habr.com/ru/rss/hubs/all/").String()
	configFilePath = kingpin.Flag("config", "Path to config file of synonyms").Default("config.yml").String()
)

func main() {
	kingpin.Parse()
	token := os.Getenv(botToken)
	if token == "" {
		log.Fatalf(`No %q env var`, botToken)
	}
	tgBot, err := bot.New(token)
	if err != nil {
		log.Fatalf("Couldn't start bot: %s", err)
	}
	conf, err := config.New(*configFilePath)
	if err != nil {
		log.Fatalf("Couldn't get config: %s", err)
	}
	u := tgbotapi.UpdateConfig{
		Timeout: 60,
	}
	updates, err := tgBot.GetUpdatesChan(u)
	if err != nil {
		log.Printf("Couldn't get update chan: %s", err)
		return
	}

	c := cache.New(*conf)
	c.ScrapePosts(*scrapeURL)

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGUSR1)
	go func() {
		for range s {
			updatedConfig, err := config.New(*configFilePath)
			if err != nil {
				log.Printf("Couldn't create config: %s", err)
				continue
			}
			c.UpdateConfig(*updatedConfig)
		}
	}()

	ticker := time.NewTicker(*updateTime)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			c.UpdatePosts(*lifeTime, *scrapeURL)
		}
	}()

	for update := range updates {
		if update.Message == nil {
			continue
		}
		switch update.Message.Command() {
		case "next":
			tgBot.NextCommand(update, c)
		case "start":
			tgBot.StartCommand(update)
		}
	}
}
