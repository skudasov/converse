package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/f4hrenh9it/parley/log"
	"go.uber.org/ratelimit"
)

var B *Bot

type Bot struct {
	Token                  string
	Api                    *tgbotapi.BotAPI
	UpdatesChan            tgbotapi.UpdatesChannel
	DeliveryRateLimit      ratelimit.Limiter
	DefaultConversationSla int
}

func StartNewBot(token string, rl int, defaultConversationSla int) {
	log.L.Debugf("using token: %s", token)
	b, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.L.Fatal(err)
	}
	b.Debug = false

	log.L.Debugf("Authorized on account %s", b.Self.UserName)

	ucfg := tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60

	updChan, err := b.GetUpdatesChan(ucfg)
	if err != nil {
		log.L.Fatalf("failed to get bot updates channel: %s", err)
	}

	B = &Bot{
		Token:             token,
		Api:               b,
		UpdatesChan:       updChan,
		DeliveryRateLimit: ratelimit.New(rl),
		DefaultConversationSla: defaultConversationSla,
	}

	for update := range B.UpdatesChan {
		if update.CallbackQuery != nil {
			err := B.HandleKbCallback(update)
			if err != nil {
				log.L.Error(err)
			}
		}
		if update.Message == nil {
			continue
		}
		if update.Message.Chat.ID == CS.SupportChat {
			if err := B.PublicRoute(update); err != nil {
				log.L.Error(err)
			}
		} else {
			if err := B.PrivateRoute(update); err != nil {
				log.L.Error(err)
			}
		}
		log.L.Infof("[%s] %s", update.Message.From.UserName, update.Message.Text)
	}
}
