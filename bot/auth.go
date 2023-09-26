package bot

import (
	"fmt"

	"git.s8k.top/SeraphJACK/beanbot/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func authorized(ctx *messageHandleContext, user *tgbotapi.User) bool {
	for _, v := range config.Cfg.AuthorizedUserIDs {
		if v == user.ID {
			return true
		}
	}

	msg := tgbotapi.NewMessage(ctx.chat.ID,
		fmt.Sprintf("User with ID `%d` is not authorized to access this bot.", user.ID))
	msg.ParseMode = tgbotapi.ModeMarkdown

	go ctx.bot.Send(msg)

	return false
}
