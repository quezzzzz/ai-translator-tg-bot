package tg_bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"tg_bot/common/postgres"
	"tg_bot/config"
	"tg_bot/internal/translator"
)

type TGBotTranslator struct {
	config     *config.Config
	bot        *tgbotapi.BotAPI
	translator *translator.YdxAITranslator
	postgres   *postgres.PostgresDB
}

func New(config *config.Config) (*TGBotTranslator, error) {
	t := translator.New(config.AITranslator)

	ctx := context.Background()
	p, err := postgres.NewWithConfig(ctx, config.Postgres.User, config.Postgres.Password, config.Postgres)
	if err != nil {
		return nil, err
	}

	bot, err := tgbotapi.NewBotAPI(config.TGBotToken)
	if err != nil {
		return nil, err
	}

	return &TGBotTranslator{
		config:     config,
		bot:        bot,
		translator: t,
		postgres:   p,
	}, nil
}

func (t *TGBotTranslator) StartBot() {
	t.bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.bot.GetUpdatesChan(u)

	ctx := context.Background()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		switch update.Message.Text {
		case "/start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет, я бот-переводчик, отправь мне любой текст и я переведу его для себя")
			t.bot.Send(msg)
		case "/history":
			history, err := t.postgres.GetUserHistory(ctx, update.CallbackQuery.Message.Chat.ID)
			if err != nil {
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Произошла ошибка")
				t.bot.Send(msg)

			} else {
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, history)
				t.bot.Send(msg)
			}
		case "":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я не могу такое перевести")
			t.bot.Send(msg)
		default:
			translatedText := t.translator.TranslateText(update.Message.Text, "en")
			if err := t.postgres.SaveTranslation(ctx, update.Message.Chat.ID, update.Message.Text, translatedText.Text); err != nil {
				log.Println(fmt.Sprintf("failed to save translation:"+
					"chat id : %d"+
					"text : %s"+
					"translated text : %s"+
					"error : %s ", update.Message.Chat.ID, update.Message.Text, translatedText.Text, err.Error()))
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, translatedText.Text)
			t.bot.Send(msg)
		}

	}
}
