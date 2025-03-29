package bot

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kirinyoku/twitch-kit/internal/fetcher"
	"github.com/kirinyoku/twitch-kit/internal/formatter"
	"github.com/kirinyoku/twitch-kit/internal/utils"
)

// UserState tracks the current state of a user's interaction with the bot.
type UserState struct {
	AwaitingUsername bool   // Indicates if the bot is waiting for a username input
	PressedButton    string // Stores the button pressed by the user
}

// Fetcher defines an interface for fetching Twitch channel data.
type Fetcher interface {
	FetchFollows(ctx context.Context, username string) ([]fetcher.Follow, error)
	FetchMods(ctx context.Context, username string) ([]fetcher.Mod, error)
	FetchVips(ctx context.Context, username string) ([]fetcher.Vip, error)
	FetchFounders(ctx context.Context, username string) ([]fetcher.Founders, error)
}

// Bot manages Telegram bot operations and state.
type Bot struct {
	api        *tgbotapi.BotAPI    // Telegram Bot API instance
	cmdViewMap map[string]ViewFunc // Maps commands to their view functions
	userState  map[int64]UserState // Tracks user interaction states by chat ID
	fetcher    Fetcher             // Interface for fetching Twitch data
}

// ViewFunc defines a function type for handling bot view commands.
//
// Parameters:
//
//	ctx - Context for controlling the operation
//	bot - Telegram Bot API instance
//	update - Update received from Telegram
//
// Returns:
//
//	An error if the operation fails
type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error

const telegramMessageLimit = 4096 // Maximum length of a Telegram message

// New creates a new Bot instance with the provided API and fetcher.
//
// Parameters:
//
//	api - Telegram Bot API instance
//	fetcher - Implementation of the Fetcher interface
//
// Returns:
//
//	A pointer to a new Bot instance
func New(api *tgbotapi.BotAPI, fetcher Fetcher) *Bot {
	return &Bot{
		api:       api,
		userState: make(map[int64]UserState),
		fetcher:   fetcher,
	}
}

// RegisterCommand associates a command with its view function.
//
// Parameters:
//
//	name - Command name (e.g., "start")
//	view - View function to handle the command
func (b *Bot) RegisterCommand(name string, view ViewFunc) {
	if b.cmdViewMap == nil {
		b.cmdViewMap = make(map[string]ViewFunc)
	}

	b.cmdViewMap[name] = view
}

// Start runs the bot and listens for updates until the context is done.
//
// Parameters:
//
//	ctx - Context for controlling the bot's lifecycle
//
// Returns:
//
//	An error if the bot fails to start or stops unexpectedly
func (b *Bot) Start(ctx context.Context) error {
	const op = "bot.Start"

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			updateCtx, updateCancel := context.WithTimeout(ctx, 5*time.Second)
			b.handleUpdate(updateCtx, update)
			updateCancel()
		case <-ctx.Done():
			return fmt.Errorf("%s: context done", op)
		}
	}
}

// handleUpdate processes incoming Telegram updates.
//
// Parameters:
//
//	ctx - Context for the operation
//	update - Telegram update to process
func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	const op = "bot.handleUpdate"

	defer func() {
		if p := recover(); p != nil {
			log.Fatalf("panic recovered: %v\n%s", p, string(debug.Stack()))
		}
	}()

	if update.CallbackQuery != nil {
		b.handleCallback(ctx, update.CallbackQuery)
		return
	}

	if update.Message == nil {
		return
	}

	if b.isAwaitingUsername(update.Message.From.ID) {
		b.handleUserInput(ctx, update)
		return
	}

	if update.Message.IsCommand() {
		b.handleCommand(ctx, update)
	} else {
		b.sendStartKeyboard(ctx, update)
	}
}

// handleCallback processes inline keyboard button presses.
//
// Parameters:
//
//	ctx - Context for the operation
//	callback - Callback query from Telegram
func (b *Bot) handleCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) {
	const op = "bot.handleCallback"

	cb := tgbotapi.NewCallback(callback.ID, callback.Data)
	if _, err := b.api.Request(cb); err != nil {
		log.Printf("%s: failed to send callback: %v", op, err)
	}

	b.userState[callback.Message.Chat.ID] = UserState{
		AwaitingUsername: true,
		PressedButton:    callback.Data,
	}

	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Enter the channel name:")
	b.api.Send(msg)
}

// handleCommand executes the appropriate view function for a command.
//
// Parameters:
//
//	ctx - Context for the operation
//	update - Telegram update containing the command
func (b *Bot) handleCommand(ctx context.Context, update tgbotapi.Update) {
	const op = "bot.handleCommand"

	cmd := update.Message.Command()
	cmdView, ok := b.cmdViewMap[cmd]
	if !ok {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command.")
		b.api.Send(msg)
		b.sendStartKeyboard(ctx, update)
		return
	}

	if err := cmdView(ctx, b.api, update); err != nil {
		log.Printf("%s: %v", op, err)
	}
}

// handleUserInput processes username input from the user.
//
// Parameters:
//
//	ctx - Context for the operation
//	update - Telegram update containing the user input
func (b *Bot) handleUserInput(ctx context.Context, update tgbotapi.Update) {
	const op = "bot.handleUserInput"

	username := update.Message.Text
	state := b.userState[update.Message.From.ID]

	delete(b.userState, update.Message.Chat.ID)

	response, err := b.processRequest(ctx, username, state.PressedButton)
	if err != nil {
		b.sendError(update.Message.Chat.ID, "Failed to fetch data", err)
		b.sendStartKeyboard(ctx, update)
		return
	}

	messages := utils.SplitMessage(response, telegramMessageLimit)
	for _, part := range messages {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, part)
		msg.ParseMode = tgbotapi.ModeHTML
		msg.DisableWebPagePreview = true
		b.api.Send(msg)
	}

	b.sendStartKeyboard(ctx, update)
}

// isAwaitingUsername checks if the bot is waiting for a username from a user.
//
// Parameters:
//
//	chatID - Telegram chat ID of the user
//
// Returns:
//
//	True if awaiting username, false otherwise
func (b *Bot) isAwaitingUsername(chatID int64) bool {
	state, ok := b.userState[chatID]
	return ok && state.AwaitingUsername
}

// processRequest fetches and formats data based on the user's selection.
//
// Parameters:
//
//	ctx - Context for the operation
//	username - Twitch username to fetch data for
//	button - Selected option (e.g., "follows", "moders")
//
// Returns:
//
//	Formatted response string and an error if any
func (b *Bot) processRequest(ctx context.Context, username, button string) (string, error) {
	switch button {
	case "follows":
		follows, err := b.fetcher.FetchFollows(ctx, username)
		if err != nil {
			return "", err
		}
		return formatter.FormatFollows(username, follows), nil

	case "moders":
		mods, err := b.fetcher.FetchMods(ctx, username)
		if err != nil {
			return "", err
		}
		return formatter.FormatMods(username, mods), nil

	case "vips":
		vips, err := b.fetcher.FetchVips(ctx, username)
		if err != nil {
			return "", err
		}
		return formatter.FormatVips(username, vips), nil

	case "founders":
		founders, err := b.fetcher.FetchFounders(ctx, username)
		if err != nil {
			return "", err
		}
		return formatter.FormatFounders(username, founders), nil
	}

	return "", fmt.Errorf("unknown button: %s", button)
}

// sendError sends an error message to the user.
//
// Parameters:
//
//	chatID - Telegram chat ID of the user
//	message - Error message text
//	err - Error details to include
func (b *Bot) sendError(chatID int64, message string, err error) {
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s: %v.", message, err))
	if _, err := b.api.Send(msg); err != nil {
		log.Printf("bot.handleUpdate: failed to send error message: %v", err)
	}
}

// sendStartKeyboard sends the initial command selection keyboard.
//
// Parameters:
//
//	ctx - Context for the operation
//	update - Telegram update to respond to
func (b *Bot) sendStartKeyboard(ctx context.Context, update tgbotapi.Update) {
	view := ViewCmdStart()
	if err := view(ctx, b.api, update); err != nil {
		log.Printf("bot.handleUpdate: failed to send inline keyboard: %v", err)
	}
}
