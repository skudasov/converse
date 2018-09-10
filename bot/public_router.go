package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/f4hrenh9it/converse/log"
	"github.com/f4hrenh9it/converse/db"
	"fmt"
)

func (m *Bot) PublicRoute(update tgbotapi.Update) error {
	log.L.Debugf("routing public msg from chatid: %d", update.Message.Chat.ID)
	userId := update.Message.From.ID
	chatId := update.Message.Chat.ID

	cmd := update.Message.Command()

	isAgent := db.IsAgent(int64(userId))
	switch cmd {
	case ConversationsListNonResolved:
		var convList []*db.ListConvInfo
		var err error
		convList, err = db.ListActiveConversations(isAgent, userId)
		if err != nil {
			return fmt.Errorf(ListActiveConvsErr, err)
		}
		if err = m.SendConversationList(chatId, convList); err != nil {
			return fmt.Errorf(SendConvList, err)
		}
	}
	return nil
}
