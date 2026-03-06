package command

import (
	"fmt"
	shared "gicgacgo/shared"

	"github.com/bwmarrin/discordgo"
)

func Stats(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userId := i.Member.User.ID
	username := i.Member.User.Username

	// Get global stats
	globalStats := shared.GlobalStats[userId]
	
	// Get local stats
	var localStats *shared.PlayerStats
	if guildStats, exists := shared.GuildLeaderboards[i.GuildID]; exists {
		localStats = guildStats.Stats[userId]
	}

	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("📊 stats for %s", username),
		Color: 0x3498db,
	}

	// Global stats field
	if globalStats != nil {
		nextRank, xpNeeded := globalStats.GetNextRank()
		progressInfo := ""
		if xpNeeded > 0 {
			progressInfo = fmt.Sprintf("\n**next rank:** %s %s %d (%d xp needed)", nextRank.Icon, nextRank.Name, nextRank.Tier, xpNeeded)
		} else {
			progressInfo = "\n🏆 **max rank achieved!**"
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: "🌍 global stats",
			Value: fmt.Sprintf("**rank:** %s | **level:** %d\n**xp:** %d%s\n**w/l/d:** %d/%d/%d | **win rate:** %s\n**total games:** %d",
				globalStats.GetRankString(), globalStats.Level, globalStats.XP, progressInfo,
				globalStats.Wins, globalStats.Losses, globalStats.Draws,
				globalStats.WinRateString(), globalStats.TotalGames()),
			Inline: false,
		})
	} else {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "🌍 global Stats",
			Value:  "no games played yet",
			Inline: false,
		})
	}

	// Local stats field
	if localStats != nil {
		nextRank, xpNeeded := localStats.GetNextRank()
		progressInfo := ""
		if xpNeeded > 0 {
			progressInfo = fmt.Sprintf("\n**next rank:** %s %s %d (%d xp needed)", nextRank.Icon, nextRank.Name, nextRank.Tier, xpNeeded)
		} else {
			progressInfo = "\n🏆 **max rank achieved!**"
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: "🏠 server stats",
			Value: fmt.Sprintf("**rank:** %s | **level:** %d\n**xp:** %d%s\n**w/l/d:** %d/%d/%d | **win rate:** %s\n**total games:** %d",
				localStats.GetRankString(), localStats.Level, localStats.XP, progressInfo,
				localStats.Wins, localStats.Losses, localStats.Draws,
				localStats.WinRateString(), localStats.TotalGames()),
			Inline: false,
		})
	} else {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "🏠 server Stats",
			Value:  "no games played on this server yet",
			Inline: false,
		})
	}

	// Add current game status
	if shared.Players[userId] != nil {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "current status",
			Value:  "in an active game",
			Inline: false,
		})
	}

	// Add XP info footer
	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("xp rewards: win +%d | draw +%d | loss +%d", shared.XPPerWin, shared.XPPerDraw, shared.XPPerLoss),
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral, // Only visible to the user
		},
	})
}
