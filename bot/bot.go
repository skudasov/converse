package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/f4hrenh9it/converse/log"
	"go.uber.org/ratelimit"
)

var B *Bot

type Bot struct {
	ResponseChan           chan tgbotapi.Chattable
	DeliveryRateLimit      ratelimit.Limiter
	DefaultConversationSla int
}

func StartSender(b *tgbotapi.BotAPI, msgChan chan tgbotapi.Chattable) {
	for {
		msgCfg := <- msgChan
		log.L.Debugf("sending msg: %s", msgCfg)
		if _, err := b.Send(msgCfg); err != nil {
			log.L.Errorf(TgApiErr, err)
		}
	}
}

func StartReceiver(
	updChan tgbotapi.UpdatesChannel,
	respChan chan tgbotapi.Chattable,
	rl int, defaultConversationSla int) {
	B = &Bot{
		ResponseChan:           respChan,
		DeliveryRateLimit:      ratelimit.New(rl),
		DefaultConversationSla: defaultConversationSla,
	}

	for update := range updChan {
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
