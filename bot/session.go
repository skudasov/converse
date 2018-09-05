package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/f4hrenh9it/parley/db"
	"github.com/f4hrenh9it/parley/log"
	"fmt"
)

type ConversationsStore struct {
	SupportChat         int64
	Conversations       map[int]*Conversation
	CurrentConversation map[int64]int
}

var CS *ConversationsStore

func NewStore(supChatId int64) {
	CS = &ConversationsStore{
		supChatId,
		make(map[int]*Conversation),
		make(map[int64]int),
	}
}

func (m *ConversationsStore) Debug() {
	for cid, c := range m.Conversations {
		log.L.Debugf("Conversation: %d", cid)
		log.L.Debugf("\tLastSeenUserId: %d", c.LastSeenUserId)
		log.L.Debugf("\tUserId: %d", c.UserId)
		log.L.Debugf("\tAgentId: %d", c.AgentId)
	}
	for k, v := range m.CurrentConversation {
		log.L.Debugf("User: %d --> ConversationId: %d", k, v)
	}
}

func (m *ConversationsStore) AlertConversationCreated(cid int) error {
	msg := conversationCreatedAlert(cid, CS.SupportChat)
	msg.ParseMode = "markdown"
	if _, err := B.Api.Send(msg); err != nil {
		return fmt.Errorf(TgApiErr, err)
	}
	return nil
}

func (m *ConversationsStore) AlertConvStatusChanged(cid int, status db.ConvStatus, who string) error {
	var msgs []tgbotapi.MessageConfig
	lastSeenUserId := CS.Conversations[cid].LastSeenUserId
	msgs = append(msgs, conversationEndedAlert(cid, lastSeenUserId, status, ""))
	msgs = append(msgs, conversationEndedAlert(cid, CS.SupportChat, status, who))

	for _, m := range msgs {
		if _, err := B.Api.Send(m); err != nil {
			return fmt.Errorf(TgApiErr, err)
		}
	}
	return nil
}

func (m *ConversationsStore) CreateConv(chatId int64) (int, error) {
	pid, err := db.NextSequence("conversations_msgs", "id")
	if err != nil {
		return 0, fmt.Errorf(db.DbSequenceErr, err)
	}
	//exit previous conversation
	cid := m.CurrentConversation[chatId]
	if cid != 0 {
		m.Conversations[cid].UserId = 0
	}
	m.Conversations[pid] = &Conversation{UserId: chatId, LastSeenUserId: chatId}
	m.CurrentConversation[chatId] = pid
	return pid, nil
}

func (m *ConversationsStore) Load(convId int) {
	if _, ok := m.Conversations[convId]; !ok {
		creator, desc := db.LoadInfo(convId)
		m.Conversations[convId] = &Conversation{LastSeenUserId: creator, Description: desc}
	}
}

func (m *ConversationsStore) isOccupiedConversation(convId int, agentChatId int64) bool {
	return m.Conversations[convId].AgentId != 0 && m.Conversations[convId].AgentId != agentChatId
}

func (m *ConversationsStore) Visit(convId int, userId int64, isAgent bool) (int64, bool) {
	if isAgent {
		if m.isOccupiedConversation(convId, userId) {
			return m.Conversations[convId].AgentId, false
		}
		currConvId := m.CurrentConversation[userId]
		if _, ok := m.Conversations[currConvId]; ok {
			m.Conversations[currConvId].AgentId = 0
		}
		m.Conversations[convId].AgentId = userId
		m.CurrentConversation[userId] = convId
	} else {
		if db.IsCreator(userId, convId) {
			currConvId := m.CurrentConversation[userId]
			if _, ok := m.Conversations[currConvId]; ok {
				m.Conversations[currConvId].UserId = 0
			}

			m.Conversations[convId].UserId = userId
			m.CurrentConversation[userId] = convId
			// we need this to notify user if his CurrentConversation != convId
			m.Conversations[convId].LastSeenUserId = userId
		} else {
			return 0, true
		}
	}
	return 0, false
}

func (m *ConversationsStore) Close(chatId int64, status db.ConvStatus) error {
	if cc, ok := m.CurrentConversation[chatId]; ok {
		if err := db.EndConversation(cc, chatId, status); err != nil {
			return fmt.Errorf(db.EndConvErr, err)
		}
		resolverName := db.NameByChatId(chatId)
		if err := m.AlertConvStatusChanged(cc, status, resolverName); err != nil {
			return fmt.Errorf(db.AlertEndConvErr, err)
		}
		m.ExitConv(chatId)
	} else {
		msg := notInAnyConversation(chatId)
		if _, err := B.Api.Send(msg); err != nil {
			return fmt.Errorf(TgApiErr, err)
		}
	}
	return nil
}

func (m *ConversationsStore) Reopen(chatId int64, status db.ConvStatus) error {
	if cc, ok := m.CurrentConversation[chatId]; ok {
		if err := db.ReopenConversation(cc); err != nil {
			return fmt.Errorf(AlertConvReopen, err)
		}
		name := db.NameByChatId(chatId)
		if err := m.AlertConvStatusChanged(cc, status, name); err != nil {
			return fmt.Errorf(AlertConvReopen, err)
		}
		delete(CS.CurrentConversation, chatId)
	} else {
		msg := notInAnyConversation(chatId)
		if _, err := B.Api.Send(msg); err != nil {
			return fmt.Errorf(TgApiErr, err)
		}
	}
	return nil
}

func (m *ConversationsStore) ExitConv(chatId int64) {
	cid := CS.CurrentConversation[chatId]
	userId := CS.Conversations[cid].UserId
	agentId := CS.Conversations[cid].AgentId
	delete(CS.CurrentConversation, userId)
	delete(CS.CurrentConversation, agentId)
	delete(CS.Conversations, cid)
}

type Conversation struct {
	UserId         int64
	AgentId        int64
	LastSeenUserId int64
	Description    string
}

func (m *Conversation) Deliver(bot *Bot, u tgbotapi.Update, isAgent bool) error {
	if isAgent {
		fc := tgbotapi.NewForward(m.UserId, m.AgentId, u.Message.MessageID)
		if _, err := B.Api.Send(fc); err != nil {
			return fmt.Errorf(TgApiErr, err)
		}
	} else {
		fc := tgbotapi.NewForward(m.AgentId, m.UserId, u.Message.MessageID)
		if _, err := B.Api.Send(fc); err != nil {
			return fmt.Errorf(TgApiErr, err)
		}
	}
	return nil
}
