package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"fmt"
	"github.com/f4hrenh9it/parley/db"
	"bytes"
)

func startMsg(chatid int64) tgbotapi.MessageConfig {
	text := "Hello! This is *SONM* support bot.\n" +
		"*Start* or *select* conversation using /active.\n" +
		"*Find* your solved conversations using /history.\n" +
		"*Finish* conversation using /current.\n"
	msg := tgbotapi.NewMessage(chatid, text)
	msg.ParseMode = "markdown"
	return msg
}

func startMsgAgent(chatid int64) tgbotapi.MessageConfig {
	text := "Hello! This is *SONM* support bot.\n" +
		"*Start* or *select* conversation using /active.\n" +
		"*Find* your solved conversations using /history.\n" +
		"*Finish* conversation using /current.\n\n" +
		"Join support group: `https://t.me/joinchat/BvQoZBLhkTqbbxurOO621g`"
	msg := tgbotapi.NewMessage(chatid, text)
	msg.ParseMode = "markdown"
	return msg
}


func searchNotAllowed(chatid int64) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatid, "search is allowed only for agents")
}

func nothingToSearch(chatid int64) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatid, "search query is empty")
}

func searchResult(chatid int64, sr *db.SearchResult) tgbotapi.MessageConfig {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("*Conversation*: `%d`\n*Name*: `%s`\n", sr.Conversation, sr.UserName))
	if sr.Caption != "" {
		buf.WriteString(fmt.Sprintf("*Caption*: `%s`\n", sr.Caption))
	}
	if sr.Msg != "" {
		buf.WriteString(fmt.Sprintf("*Msg*: `%s`\n", sr.Msg))
	}
	return ConvLinkKb(chatid, sr.Conversation, buf.String())
}


func myChatId(chatid int64) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatid, fmt.Sprintf("your chat id: %d", chatid))
}

func unrecognizedCommand(chatid int64) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatid, "Command is not recognized, see command list")
}

func noSuchConversation(chatid int64) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatid, "No such conversation")
}

func conversationCreated(chatid int64, convId int) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatid, fmt.Sprintf("*Conversation*: %d\nDescribe your problem, first message *must* be text\nUse /current to end conversation", convId))
	msg.ParseMode = "markdown"
	return msg
}

func conversationJoined(chatid int64, convId int) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatid, fmt.Sprintf("Now participating in conversation: %d", convId))
}

func creationRestricted(chatid int64) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatid, "Sorry, conversation creation for support members is restricted for now")
}

func conversationClosed(chatid int64, convId int) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatid, fmt.Sprintf("Conversation closed: %d", convId))
}

func conversationReopened(chatid int64, convId int) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatid, fmt.Sprintf("Conversation reopened: %d", convId))
}

func conversationCreatedAlert(cid int, chatid int64) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatid, fmt.Sprintf("Conversation: `%d`\n*status*: `opened`", cid))
}

func conversationEndedAlert(cid int, chatid int64, status db.ConvStatus, who string) tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	if who != "" {
		msg = tgbotapi.NewMessage(chatid, fmt.Sprintf("Conversation: `%d`\n*status*: `%s`\n*Who*: `%s`\n", cid, status.String(), who))
	} else {
		msg = tgbotapi.NewMessage(chatid, fmt.Sprintf("Conversation: `%d`\n*status*: `%s`\n", cid, status.String()))
	}
	msg.ParseMode = "markdown"
	return msg
}

func conversationsListEmpty(chatid int64) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatid, "Conversations list is empty")
}

func conversationOccupied(chatid int64, name string) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatid, fmt.Sprintf("Conversation is occupied by: %s", name))
}

func newAnswerAlert(chatid int64, convId int, desc string, name string) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatid, fmt.Sprintf("Coversation: `%d`\n*Desc*: `%s`\n*New answer from %s*", convId, desc, name))
}

func newQuestionAlert(chatid int64, convId int, desc string, name string) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatid, fmt.Sprintf("Coversation: `%d`\n*Desc*: `%s`\n*New question from %s*", convId, desc, name))
}

func notInAnyConversation(chatid int64) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(chatid, "You are not participating in any conversation, please select one using /active or /history")
}

func visitAlert(chatId int64, cid int, name string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatId, fmt.Sprintf("Conversation: `%d`\n`%s` joined", cid, name))
	msg.ParseMode = "markdown"
	return msg
}
