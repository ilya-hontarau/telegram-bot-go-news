package gonews

import (
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/illfate/telegram-bot-go-news/pkg/bot"
	"github.com/illfate/telegram-bot-go-news/pkg/cache"
)

const (
	botToken            = "BOT_TOKEN"
	lifeTimePostsEnvVar = "LIFE_TIME"
	gopherPicFileEnvVar = "GOPHER_PIC_FILE"
	updateTimeEnvVar    = "UPDATE_TIME"
	scrapeUrlEnvVar     = "SCRAPE_URL"
)

func main() {
	token := os.Getenv(botToken)
	if token == "" {
		log.Fatalf(`no %q env var`, botToken)
	}
	lifeTimePosts := os.Getenv(lifeTimePostsEnvVar)
	if token == "" {
		log.Fatalf(`no %q env var`, lifeTimePostsEnvVar)
	}
	lifeTime, err := time.ParseDuration(lifeTimePosts)
	if err != nil {
		log.Fatalf("couldn't parse time duration: %s", err)
	}
	gopherPicFile := os.Getenv(gopherPicFileEnvVar)
	if gopherPicFile == "" {
		log.Fatalf(`no %q env var`, gopherPicFileEnvVar)
	}
	updateTime := os.Getenv(updateTimeEnvVar)
	if updateTime == "" {
		log.Fatalf(`no %q env var`, updateTimeEnvVar)
	}
	scrapeUrl := os.Getenv(scrapeUrlEnvVar)
	if scrapeUrl == "" {
		log.Fatalf(`no %q env var`, scrapeUrlEnvVar)
	}
	updTime, err := time.ParseDuration(updateTime)
	if err != nil {
		log.Fatalf("couldn't parse update time: %s", err)
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
	c.ScrapePosts(scrapeUrl)

	ticker := time.NewTicker(updTime)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ticker.C:
				c.UpdatePosts(lifeTime, scrapeUrl)
			}
		}
	}()

	for update := range updates {
		switch update.Message.Command() {
		case "next":
			tgBot.NextCommand(update, c, gopherPicFile)
		case "start":
			tgBot.StartCommand(update)
		}
	}
}
