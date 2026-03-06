package command

import (
	"fmt"
	shared "gicgacgo/shared"

	"github.com/bwmarrin/discordgo"
)

func Leaderboard(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	// leaderboard type (local or global)
	var leaderboardType string
	if len(options) > 0 {
		leaderboardType = options[0].StringValue()
	} else {
		leaderboardType = "global"
	}

	var stats []*shared.PlayerStats
	var title string

	if leaderboardType == "local" {
		stats = shared.GetGuildLeaderboard(i.GuildID, 10)
		title = "local server leaderboard - top 10"
	} else {
		stats = shared.GetGlobalLeaderboard(10)
		title = "global leaderboard - top 10"
	}

	embed := &discordgo.MessageEmbed{
		Title: title,
		Color: 0x00ff00,
	}

	if len(stats) == 0 {
		embed.Description = "no games have been played yet! use `/duel` to start playing."
	} else {
		for rank, player := range stats {
			var medal string
			switch rank {
			case 0:
				medal = "🥇"
			case 1:
				medal = "🥈"
			case 2:
				medal = "🥉"
			default:
				medal = fmt.Sprintf("**#%d**", rank+1)
			}

			name := fmt.Sprintf("%s %s - %s", medal, player.Username, player.GetRankString())
			value := fmt.Sprintf("**level:** %d | **xp:** %d\n**w/l/d:** %d/%d/%d | **win rate:** %s\n**total games:** %d",
				player.Level, player.XP,
				player.Wins, player.Losses, player.Draws, player.WinRateString(),
				player.TotalGames())

			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   name,
				Value:  value,
				Inline: false,
			})
		}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
