package buttons

import (
	"fmt"
	shared "gicgacgo/shared"
	"time"

	"github.com/bwmarrin/discordgo"
)

func AcceptDuel(s *discordgo.Session, i *discordgo.InteractionCreate, inviter string, invitee string) {
	if i.Member.User.ID != invitee {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "respectfully, you are not authorized to respond to this invitation",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	if shared.Players[invitee] != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "sorry, ure already in a game",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	gameId := fmt.Sprintf("%s_%s", inviter, invitee)

	shared.Players[inviter] = &shared.Player{GameId: gameId, Id: inviter}
	shared.Players[invitee] = &shared.Player{GameId: gameId, Id: invitee}

	// Get player usernames
	inviterUser, _ := s.User(inviter)
	inviteeUser, _ := s.User(invitee)
	inviterName := inviterUser.Username
	inviteeName := inviteeUser.Username

	shared.Games[gameId] = &shared.Game{
		StartedTimestamp: time.Now(),
		Players:          []shared.Player{*shared.Players[inviter], *shared.Players[invitee]},
		Turn:             shared.RandomizeTurn(),
		PlayerX:          *shared.Players[inviter],
		PlayerY:          *shared.Players[invitee],
		PlayerXName:      inviterName,
		PlayerYName:      inviteeName,
		Game: shared.Board{
			{"", "", ""},
			{"", "", ""},
			{"", "", ""},
		},
		ChannelId: i.ChannelID,
		GuildId:   i.GuildID,
	}

	shared.DisableAllButtons(s, i)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("duel has started between <@%s> and <@%s>! Use `/place <x> <y>` to play.", inviter, invitee),
		},
	})

	shared.StartGame(s, i, gameId)
}
