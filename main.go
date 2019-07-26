package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gocolly/colly"
	"log"
	"os"
	"strings"
	"time"
)

const (
	botToken = "BOT_TOKEN"
)

func userHasLink(userName, userUrl string) bool {
	for _, url := range userUrls[userName] {
		if url == userUrl {
			return true
		}
	}
	return false
}

var userUrls map[string][]string

type postsCache map[string][]Post

var cache postsCache


type Post struct {
	Link      string
	AddedTime time.Time
}

type scraper struct {
	tagsToPost postsCache
	*colly.Collector
}

func ScrapPosts() {
	c := colly.NewCollector()
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	c.OnXML("/rss/channel/item", func(e *colly.XMLElement) {
		link := e.ChildText("/link")
		idx := strings.Index(link, "?")
		if idx != -1 {
			link = link[:idx]
		}
		for _, category := range e.ChildTexts("//category") {
			cache[category] = append(cache[category], Post{
				Link:      link,
				AddedTime: time.Now(),
			})
		}
	})
	c.Visit("https://habr.com/ru/rss/hubs/all/")
	c.Wait()
}


func DeleteOldCache() {
	for category, posts := range cache {
		for _, post := range posts {
			if time.Now().Sub(post.AddedTime) > time.Duration(24*time.Hour) {
				delete(cache, category)
			}
		}
	}
}

func GetLink(categories []string, userName string) string {
	for _, category := range categories {
		for _, post := range cache[category] {
			if !userHasLink(userName, post.Link) {
				return post.Link
			}
		}
	}
	return ""
}

func logCache(){
	for k, v := range cache {
		log.Print(k, " - ", v)
	}
}

func main() {
	token := os.Getenv(botToken)
	if token == "" {
		log.Printf(`no "%s"env var `, botToken)
		return
	}
	userUrls = make(map[string][]string)
	cache = postsCache{}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Printf("couldn't start bot: %s", err)
		return
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return
	}

	ScrapPosts()
	logCache()

	ticker := time.NewTicker(10 * time.Second)
	go func() {
		select {
		case <-ticker.C:
			ScrapPosts()
			DeleteOldCache()
		}
	}()

	for update := range updates {
		switch update.Message.Command() {
		case "next":
			categoriesArgs := update.Message.CommandArguments()
			if categoriesArgs == "" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Empty categories arguments.")
				bot.Send(msg)
				break
			}
			categories := strings.Split(categoriesArgs, " ")
			userName := update.Message.From.UserName
			link := GetLink(categories, userName)
			if len(link) == 0 {
				photoConfig := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, "gopher-no.png")
				photoConfig.Caption = "No new posts"
				bot.Send(photoConfig)
				break
			}
			_, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, link))
			if err != nil {
				log.Print("couldn't send post link")
				break
			}
			userUrls[userName] = append(userUrls[userName], link)
		case "start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s\n%s",
				"Hi, i'm a go news bot!", "Type /next [tag] to get new post with this tag."))
			bot.Send(msg)
		}

	}
}
