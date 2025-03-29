package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ViewCmdStart creates a view handler for the bot's start command.
// It presents an inline keyboard with options for the user to select.
//
// Returns:
//
//	A ViewFunc that handles the start command interaction
func ViewCmdStart() ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("follows", "follows"),
				tgbotapi.NewInlineKeyboardButtonData("moders", "moders"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("vips", "vips"),
				tgbotapi.NewInlineKeyboardButtonData("founders", "founders"),
			),
		)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Select the option:")
		msg.ReplyMarkup = inlineKeyboard

		_, err := bot.Send(msg)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}

		return nil
	}
}
