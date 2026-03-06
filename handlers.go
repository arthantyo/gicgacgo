package main

import (
	buttons "gicgacgo/buttons"
	command "gicgacgo/commands"
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// params is not just right
// TODO not viable when implementing the leaderboard
var buttonHandler = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, p1 string, p2 string){
	"accept": buttons.AcceptDuel,
	"reject": buttons.RejectDuel,
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"ping":        command.Ping,
	"duel":        command.Duel,
	"leaderboard": command.Leaderboard,
	"place":       command.Place,
	"stats":       command.Stats,
	"ranks":       command.Ranks,
	"train":       command.Train,
}

func handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		handleCommandInteraction(s, i)
	case discordgo.InteractionMessageComponent:
		handleButtonInteraction(s, i)
	}
}

func handleCommandInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
		h(s, i)
	}
}

func handleButtonInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// custom parsing
	// eg: <action>_<userid>_<userid> -> depends on the action (some actions may not have other params)
	// TODO doesnt seem to scale well

	parts := strings.Split(i.MessageComponentData().CustomID, "_")
	action := parts[0]
	inviter := parts[1]
	invitee := parts[2]

	if action == "accept" {
		buttonHandler["accept"](s, i, inviter, invitee)
	} else if action == "decline" {
		buttonHandler["reject"](s, i, inviter, invitee)
	}
}

func handleReady(s *discordgo.Session, r *discordgo.Ready) {
	slog.Info("logged in as",
		slog.String("username", s.State.User.Username),
		slog.String("discriminator", s.State.User.Discriminator))
}
