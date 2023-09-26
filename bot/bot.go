package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"git.s8k.top/SeraphJACK/beanbot/config"
	"git.s8k.top/SeraphJACK/beanbot/syntax"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

type messageHandleContext struct {
	bot  *tgbotapi.BotAPI
	chat *tgbotapi.Chat
	msg  *tgbotapi.Message
}

func Start() error {
	bot, err := tgbotapi.NewBotAPI(config.Cfg.BotToken)
	if err != nil {
		return err
	}

	// Watch for transactions that need to be committed
	go func() {
		for {
			commitAll()
			time.Sleep(time.Second)
		}
	}()

	// Polling message updates
	id := 0
	for {
		id++
		u := tgbotapi.NewUpdate(id)
		u.Timeout = 60

		updates := bot.GetUpdatesChan(u)

		for update := range updates {
			if update.Message != nil {
				ctx := &messageHandleContext{bot: bot, chat: update.Message.Chat, msg: update.Message}
				go handleMessage(ctx)
			}
			if update.EditedMessage != nil {
				ctx := &messageHandleContext{bot: bot, chat: update.EditedMessage.Chat, msg: update.EditedMessage}
				go handleMessage(ctx)
			}
			if update.CallbackQuery != nil {
				handleCallbackQuery(update.CallbackQuery)
			}
		}
	}
}

func handleMessage(ctx *messageHandleContext) {
	msg := ctx.msg
	// We only process private messages
	if msg.Chat.Type != "private" {
		return
	}

	// User is not authorized, break
	if !authorized(ctx, msg.From) {
		return
	}

	if msg.Command() == "recent" {
		updateRecentKeyboard(ctx)
		return
	}

	raw := msg.Text

	txn, err := syntax.Parse(strings.Split(raw, " "), &config.Cfg.Syntax)
	if err != nil {
		msg := tgbotapi.NewMessage(ctx.chat.ID, "Failed to parse txn syntax: "+err.Error())
		go ctx.bot.Send(msg)
		return
	}

	txnID := uuid.New().String()
	msgCfg := tgbotapi.NewMessage(ctx.chat.ID,
		fmt.Sprintf("The following transaction is about to be committed:```\n%s\n```",
			txn.ToBeanLanguageSyntax()),
	)
	msgCfg.ParseMode = tgbotapi.ModeMarkdown
	msgCfg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Cancel", txnID),
	})

	confirmMsg, err := ctx.bot.Send(msgCfg)
	if err != nil {
		log.Printf("Failed to send txn confirmation message: %v", err)
		return
	}

	aboutToCommit(txnID, &transaction{
		ctx:        ctx,
		raw:        raw,
		txn:        txn,
		confirmMsg: confirmMsg,
		commitTime: time.Now().Add(10 * time.Second),
	})
}

func handleCallbackQuery(query *tgbotapi.CallbackQuery) {
	cancel(query.Data)
}

func updateRecentKeyboard(ctx *messageHandleContext) {
	var btns [][]tgbotapi.KeyboardButton
	for _, v := range recentCmds {
		btns = append(btns, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(v.command),
		))
	}
	msg := tgbotapi.NewMessage(ctx.chat.ID, "Recent Commands")
	if len(btns) > 0 {
		markup := tgbotapi.NewReplyKeyboard(btns...)
		markup.OneTimeKeyboard = true
		msg.ReplyMarkup = markup
	}
	_, err := ctx.bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send recent menu: %v", err)
	}
	return
}
