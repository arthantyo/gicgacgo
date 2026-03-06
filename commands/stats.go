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
		Title: fmt.Sprintf("📊 Stats for %s", username),
		Color: 0x3498db,
	}

	// Global stats field
	if globalStats != nil {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: "🌍 Global Stats",
			Value: fmt.Sprintf("**Wins:** %d | **Losses:** %d | **Draws:** %d\n**Win Rate:** %s | **Total Games:** %d",
				globalStats.Wins, globalStats.Losses, globalStats.Draws, 
				globalStats.WinRateString(), globalStats.TotalGames()),
			Inline: false,
		})
	} else {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "🌍 Global Stats",
			Value:  "No games played yet",
			Inline: false,
		})
	}

	// Local stats field
	if localStats != nil {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: "🏠 Server Stats",
			Value: fmt.Sprintf("**Wins:** %d | **Losses:** %d | **Draws:** %d\n**Win Rate:** %s | **Total Games:** %d",
				localStats.Wins, localStats.Losses, localStats.Draws,
				localStats.WinRateString(), localStats.TotalGames()),
			Inline: false,
		})
	} else {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "🏠 Server Stats",
			Value:  "No games played on this server yet",
			Inline: false,
		})
	}

	// Add current game status
	if shared.Players[userId] != nil {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "⚔️ Current Status",
			Value:  "In an active game",
			Inline: false,
		})
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral, // Only visible to the user
		},
	})
}
