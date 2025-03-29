package formatter

import (
	"fmt"

	"github.com/kirinyoku/twitch-kit/internal/fetcher"
)

// FormatFollows creates a formatted string of Twitch users that the specified user follows.
// It generates HTML links and includes follow dates.
//
// Parameters:
//
//	username - The Twitch username of the follower
//	follows - Slice of Follow structs containing follow information
//
// Returns:
//
//	A formatted string with HTML links and follow details
func FormatFollows(username string, follows []fetcher.Follow) string {
	response := fmt.Sprintf("<a href=\"https://twitch.tv/%s\">%s</a> is following:\n", username, username)
	for i, follow := range follows {
		twitchURL := fmt.Sprintf("https://twitch.tv/%s", follow.Login)
		formattedDate := follow.FollowedAt.Format("2006-01-02")
		response += fmt.Sprintf("%d. <a href=\"%s\">%s</a> (followed at %s)\n", i+1, twitchURL, follow.DisplayName, formattedDate)
	}
	return response
}

// FormatMods creates a formatted string of Twitch channel moderators for a user.
// It generates HTML links and includes moderation grant dates.
//
// Parameters:
//
//	username - The Twitch username of the channel owner
//	mods - Slice of Mod structs containing moderator information
//
// Returns:
//
//	A formatted string with HTML links and moderator details
func FormatMods(username string, mods []fetcher.Mod) string {
	response := fmt.Sprintf("<a href=\"https://twitch.tv/%s\">%s</a>'s list of channel moders:\n", username, username)
	for i, mod := range mods {
		twitchURL := fmt.Sprintf("https://twitch.tv/%s", mod.Login)
		formattedDate := mod.GrantedAt.Format("2006-01-02")
		response += fmt.Sprintf("%d. <a href=\"%s\">%s</a> (moded at %s)\n", i+1, twitchURL, mod.DisplayName, formattedDate)
	}
	return response
}

// FormatVips creates a formatted string of Twitch channel VIPs for a user.
// It generates HTML links and includes VIP grant dates.
//
// Parameters:
//
//	username - The Twitch username of the channel owner
//	vips - Slice of Vip structs containing VIP information
//
// Returns:
//
//	A formatted string with HTML links and VIP details
func FormatVips(username string, vips []fetcher.Vip) string {
	response := fmt.Sprintf("<a href=\"https://twitch.tv/%s\">%s</a>'s list of channel vips:\n", username, username)
	for i, vip := range vips {
		twitchURL := fmt.Sprintf("https://twitch.tv/%s", vip.Login)
		formattedDate := vip.GrantedAt.Format("2006-01-02")
		response += fmt.Sprintf("%d. <a href=\"%s\">%s</a> (viped at %s)\n", i+1, twitchURL, vip.DisplayName, formattedDate)
	}
	return response
}

// FormatFounders creates a formatted string of Twitch channel founders for a user.
// It generates HTML links and includes founding dates.
//
// Parameters:
//
//	username - The Twitch username of the channel owner
//	founders - Slice of Founders structs containing founder information
//
// Returns:
//
//	A formatted string with HTML links and founder details
func FormatFounders(username string, founders []fetcher.Founders) string {
	response := fmt.Sprintf("<a href=\"https://twitch.tv/%s\">%s</a>'s list of channel founders:\n", username, username)
	for i, founder := range founders {
		twitchURL := fmt.Sprintf("https://twitch.tv/%s", founder.Login)
		formattedDate := founder.FirstMonth.Format("2006-01-02")

		response += fmt.Sprintf("%d. <a href=\"%s\">%s</a> (founded at %s)\n", i+1, twitchURL, founder.DisplayName, formattedDate)
	}
	return response
}
