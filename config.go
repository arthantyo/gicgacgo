package main

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func setupBot() (*discordgo.Session, error) {
	botToken := os.Getenv("DISCORD_BOT_TOKEN")
	if botToken == "" {
		log.Fatalf("bot token not provided")
	}

	s, err := discordgo.New("Bot " + botToken)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// getGuildID returns the guild ID for command registration
// Returns the dev guild ID if DEV_MODE is set, otherwise returns empty string for global commands
func getGuildID() string {
	if os.Getenv("DEV_MODE") == "true" {
		return os.Getenv("DEV_GUILD_ID")
	}
	return "" // Empty string registers commands globally
}
