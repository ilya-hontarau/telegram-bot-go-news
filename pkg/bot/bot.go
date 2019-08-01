package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/illfate/telegram-bot-go-news/pkg/cache"
	"github.com/pkg/errors"
)

type Bot struct {
	*tgbotapi.BotAPI
}

func New(token string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create new bot api")
	}
	return &Bot{
		BotAPI: bot,
	}, nil
}

func (bot *Bot) NextCommand(update tgbotapi.Update, cache *cache.Cache) {
	category := update.Message.CommandArguments()
	if category == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Empty categories arguments.")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("couldn't send message: %s", err)
		}
		return
	}

	userName := update.Message.From.UserName
	link := cache.GetLink(category, userName)
	if link == "" {
		photoConfig := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, "gopher-no.png")
		photoConfig.Caption = "No new posts"
		if _, err := bot.Send(photoConfig); err != nil {
			log.Printf("couldn't send photo: %s", err)
		}
		return
	}

	if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, link)); err != nil {
		log.Printf("couldn't send post link: %s", err)
		return
	}
	cache.AddUserUrl(userName, link)
}

func (bot *Bot) StartCommand(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, `Hi, i'm a go news bot!
Type /next [category] to get new post with this category.`)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("couldn't send message: %s", err)
	}
}
