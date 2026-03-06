package command

import (
	"fmt"
	shared "gicgacgo/shared"
	"time"

	"github.com/bwmarrin/discordgo"
)

func Train(s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.Member.User

	if shared.Players[user.ID] != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "you're already in a game! finish it before starting a training session",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	gameId := fmt.Sprintf("%s_ai", user.ID)

	shared.Players[user.ID] = &shared.Player{GameId: gameId, Id: user.ID}

	// Create AI player
	aiPlayerId := "ai_bot"
	shared.Players[aiPlayerId] = &shared.Player{GameId: gameId, Id: aiPlayerId}

	turn := shared.RandomizeTurn()

	shared.Games[gameId] = &shared.Game{
		StartedTimestamp: time.Now(),
		Players:          []shared.Player{*shared.Players[user.ID], *shared.Players[aiPlayerId]},
		Turn:             turn,
		PlayerX:          *shared.Players[user.ID],
		PlayerY:          *shared.Players[aiPlayerId],
		PlayerXName:      user.Username,
		PlayerYName:      "AI Bot",
		Game: shared.Board{
			{"", "", ""},
			{"", "", ""},
			{"", "", ""},
		},
		ChannelId:      i.ChannelID,
		GuildId:        i.GuildID,
		IsSinglePlayer: true,
		AI:             &shared.AI{},
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("🤖 training mode started! you are playing against an robotic opponent.\nuse `/place <row> <col>` to make your move!"),
		},
	})

	shared.StartGame(s, i, gameId)

	// If AI goes first, make its move
	if turn == "O" {
		shared.MakeAIMove(s, i, gameId)
	}
}
