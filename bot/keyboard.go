package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"fmt"
)

func ConvsToKb(convs []string, rowCapacity int) [][]string {
	var res [][]string
	for i := 0; i < len(convs); i += rowCapacity {
		end := i + rowCapacity
		if end > len(convs) {
			end = len(convs)
		}
		res = append(res, convs[i:end])
	}
	return res
}

func ConvsToStr(convs []int) []string {
	var c []string
	for _, conv := range convs {
		c = append(c, strconv.Itoa(conv))
	}
	return c
}

func setupKb() tgbotapi.InlineKeyboardMarkup {
	kb := tgbotapi.NewInlineKeyboardMarkup()
	kb.InlineKeyboard = make([][]tgbotapi.InlineKeyboardButton, 0)
	return kb
}

func ConversationSelectorKb(chatid int64, convs []int, creationAllowed bool) tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	var msgText string
	var kb tgbotapi.InlineKeyboardMarkup
	convsStr := ConvsToStr(convs)
	splittedConvs := ConvsToKb(convsStr, 3)
	kb = setupKb()
	if creationAllowed {
		kb.InlineKeyboard = append(
			kb.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("New conversation", "newConversation"),
			),
		)
	}
	if len(convs) == 0 {
		msgText = "No conversations found"
	} else {
		msgText = "Select conversation"
		for _, convsRow := range splittedConvs {
			row := tgbotapi.NewInlineKeyboardRow()
			for _, conv := range convsRow {
				row = append(row, tgbotapi.NewInlineKeyboardButtonData(conv, conv))
			}
			kb.InlineKeyboard = append(kb.InlineKeyboard, row)
		}
	}
	msg = tgbotapi.NewMessage(chatid, msgText)
	msg.ReplyMarkup = kb
	return msg
}

func ConversationCurrentKb(chatid int64, curConv int) tgbotapi.MessageConfig {
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Resolve", ResolveCurrentConversationCbData),
			tgbotapi.NewInlineKeyboardButtonData("Close", CloseCurrentConversationCbData),
			tgbotapi.NewInlineKeyboardButtonData("Reopen", ReopenCurrentConversationCbData),
		),
	)
	msg := tgbotapi.NewMessage(chatid, fmt.Sprintf("Current conversation: %d", curConv))
	msg.ReplyMarkup = kb
	return msg
}

func ConvLinkKb(chatId int64, convId int, text string) tgbotapi.MessageConfig {
	cid := strconv.Itoa(convId)
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("Link: %s", cid),
				cid),
		),
	)
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ReplyMarkup = kb
	msg.ParseMode = "markdown"
	return msg
}
