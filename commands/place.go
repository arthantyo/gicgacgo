package command

import (
	"fmt"
	shared "gicgacgo/shared"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func Place(s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.Member.User
	player, exists := shared.Players[user.ID]
	if !exists {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "grrr, you have to be in a game to place a marker!!!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	game, exists := shared.Games[player.GameId]
	if !exists {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "hmm, there was an error finding your game.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	var row, col int
	for _, option := range i.ApplicationCommandData().Options {
		if option.Name == "row" {
			row = int(option.IntValue()) - 1
		}
		if option.Name == "col" {
			col = int(option.IntValue()) - 1
		}
	}

	if (game.Turn == "X" && game.PlayerX.Id != user.ID) || (game.Turn == "O" && game.PlayerY.Id != user.ID) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "bro, it's not your turn!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	if game.Game[row][col] != "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "yo, the cell is already occupied!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	// place points to board
	shared.PlaceMarker(s, i, player.GameId, row, col)
	shared.EditBoardEmbed(s, i, player.GameId)

	_, won := shared.CheckWin(game.Game)

	if won {
		shared.EndGame(s, i, game, user.ID, false)
		return
	}

	draw := shared.CheckDraw(game.Game)

	if draw {
		shared.EndGame(s, i, game, "", true)
		return
	}

	if game.Turn == "X" {
		game.Turn = "O"
	} else {
		game.Turn = "X"
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("moved ur marker to row %s and col %s", strconv.Itoa(row), strconv.Itoa(col)),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	shared.EditMessageBoardEmbed(s, i, player.GameId)
}
