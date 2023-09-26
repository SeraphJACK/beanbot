package bot

import (
	"fmt"
	"sync"
	"time"

	"git.s8k.top/SeraphJACK/beanbot/repo"
	"git.s8k.top/SeraphJACK/beanbot/syntax"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type transaction struct {
	ctx         *messageHandleContext
	txn         *syntax.Transaction
	txmMsg      tgbotapi.Message
	callbackMsg tgbotapi.Message
	commitTime  time.Time
}

var lock sync.Mutex
var aboutToCommitTxn map[string]*transaction

func aboutToCommit(id string, txn *transaction) {
	lock.Lock()
	defer lock.Unlock()

	aboutToCommitTxn[id] = txn
}

func cancel(id string) {
	lock.Lock()
	defer lock.Unlock()

	if v, ok := aboutToCommitTxn[id]; ok {
		delete(aboutToCommitTxn, id)
		// delete transaction confirmation message
		go v.ctx.bot.Send(tgbotapi.NewDeleteMessage(v.ctx.chat.ID, v.txmMsg.MessageID))
		// delete transaction cancel callback message
		go v.ctx.bot.Send(tgbotapi.NewDeleteMessage(v.ctx.chat.ID, v.callbackMsg.MessageID))
	}
}

func commitAll() {
	lock.Lock()
	defer lock.Unlock()

	for k, v := range aboutToCommitTxn {
		// txn not due to commit, skip
		if time.Now().Before(v.commitTime) {
			continue
		}

		err := repo.CommitTransaction(v.txn)
		if err != nil {
			go v.ctx.bot.Send(tgbotapi.NewMessage(v.ctx.chat.ID, fmt.Sprintf("Failed to commit txn: %v", err)))
		}

		delete(aboutToCommitTxn, k)

		if err == nil {
			// delete transaction message
			go v.ctx.bot.Send(tgbotapi.NewDeleteMessage(v.ctx.chat.ID, v.ctx.msg.MessageID))
		}
		// delete transaction confirmation message
		go v.ctx.bot.Send(tgbotapi.NewDeleteMessage(v.ctx.chat.ID, v.txmMsg.MessageID))
		// delete transaction cancel callback message
		go v.ctx.bot.Send(tgbotapi.NewDeleteMessage(v.ctx.chat.ID, v.callbackMsg.MessageID))
	}
}
