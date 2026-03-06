package command

import (
	"fmt"
	shared "gicgacgo/shared"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Ranks(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "ranking system",
		Description: "climb the ranks by playing games and earning exp!",
		Color:       0xffd700,
	}

	var currentTier strings.Builder
	var lastTierName string

	for _, rank := range shared.Ranks {
		if rank.Name != lastTierName {
			// Add the previous tier to embed if exists
			if lastTierName != "" {
				embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
					Name:   lastTierName,
					Value:  currentTier.String(),
					Inline: true,
				})
				currentTier.Reset()
			}
			lastTierName = rank.Name
		}

		currentTier.WriteString(fmt.Sprintf("%s **Tier %d** - %d XP\n", rank.Icon, rank.Tier, rank.MinXP))
	}

	// Add the last tier
	if lastTierName != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   lastTierName,
			Value:  currentTier.String(),
			Inline: true,
		})
	}

	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("earn xp: win +%d | draw +%d | loss +%d", shared.XPPerWin, shared.XPPerDraw, shared.XPPerLoss),
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
