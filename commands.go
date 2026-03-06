package main

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func registerCommands(s *discordgo.Session) ([]*discordgo.ApplicationCommand, error) {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "gives ping pong",
		}, {
			Name: "leaderboard",
			Options: []*discordgo.ApplicationCommandOption{{
				Name:        "type",
				Type:        discordgo.ApplicationCommandOptionString,
				Choices:     []*discordgo.ApplicationCommandOptionChoice{{Name: "local", Value: "local"}, {Name: "global", Value: "global"}},
				Description: "scope of the placement",
				Required:    true,
			}},
			Description: "shows the guild ranking",
		},
		{
			Name:        "duel",
			Description: "invite someone to play against",
			Options: []*discordgo.ApplicationCommandOption{{
				Name:        "username",
				Type:        discordgo.ApplicationCommandOptionUser,
				Description: "person you want to duel",
				Required:    true,
			}},
		},
		{
			Name:        "place",
			Description: "point in the grid to place ur marker",
			Options: []*discordgo.ApplicationCommandOption{{
				Name:        "row",
				Type:        discordgo.ApplicationCommandOptionInteger,
				Description: "row of the tictactoe grid",
				Choices:     []*discordgo.ApplicationCommandOptionChoice{{Name: "one", Value: 1}, {Name: "two", Value: 2}, {Name: "three", Value: 3}},
				Required:    true,
			}, {
				Name:        "col",
				Type:        discordgo.ApplicationCommandOptionInteger,
				Description: "column of the tictactoe grid",
				Choices:     []*discordgo.ApplicationCommandOptionChoice{{Name: "one", Value: 1}, {Name: "two", Value: 2}, {Name: "three", Value: 3}},
				Required:    true,
			},
			},
		},
		{
			Name:        "stats",
			Description: "view your personal game statistics",
		},
		{
			Name:        "ranks",
			Description: "view all available ranks and XP requirements",
		},
		{
			Name:        "train",
			Description: "practice against an AI opponent",
		},
	}

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, getGuildID(), v)
		if err != nil {
			slog.Error("cannot create command", slog.String("command", v.Name), slog.Any("error", err))
		}
		registeredCommands[i] = cmd
	}

	return registeredCommands, nil
}

func removeCommands(s *discordgo.Session, commands []*discordgo.ApplicationCommand) {
	for _, v := range commands {
		err := s.ApplicationCommandDelete(s.State.User.ID, getGuildID(), v.ID)
		if err != nil {
			slog.Error("cannot delete command", slog.String("command", v.Name), slog.Any("error", err))
		}
	}
}
