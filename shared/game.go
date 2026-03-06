package shared

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Board [3][3]string

type Game struct {
	Game             Board
	StartedTimestamp time.Time
	EndedTimestamp   *time.Time
	PlayerX          Player
	PlayerY          Player
	PlayerXName      string
	PlayerYName      string
	Players          []Player
	Turn             string
	BoardId          string
	BoardMessageId   string
	ChannelId        string
	GuildId          string
}

type Player struct {
	Id     string
	GameId string
}

var Games = make(map[string]*Game)
var Players = make(map[string]*Player)

func EndGame(s *discordgo.Session, i *discordgo.InteractionCreate, game *Game, winnerId string, isDraw bool) {
	// Determine winner and loser IDs
	var loserId string
	var message string
	if isDraw {
		message = "woof. that was an equal game you should duel again to see who's the real deal"
		RecordGameResult(game.PlayerX.Id, game.PlayerY.Id, game.GuildId, game.PlayerXName, game.PlayerYName, true)
	} else {
		if winnerId == game.PlayerX.Id {
			loserId = game.PlayerY.Id
		} else {
			loserId = game.PlayerX.Id
		}
		message = fmt.Sprintf("congratulations <@%s>! you won the game!", winnerId)
		var winnerName, loserName string
		if winnerId == game.PlayerX.Id {
			winnerName = game.PlayerXName
			loserName = game.PlayerYName
		} else {
			winnerName = game.PlayerYName
			loserName = game.PlayerXName
		}
		RecordGameResult(winnerId, loserId, game.GuildId, winnerName, loserName, false)
	}

	delete(Games, game.PlayerX.GameId)
	delete(Players, game.PlayerX.Id)
	delete(Players, game.PlayerY.Id)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}

func CheckWin(board Board) (string, bool) {
	for i := 0; i < 3; i++ {
		if board[i][0] != "" && board[i][0] == board[i][1] && board[i][1] == board[i][2] {
			return board[i][0], true
		}
		if board[0][i] != "" && board[0][i] == board[1][i] && board[1][i] == board[2][i] {
			return board[0][i], true
		}
	}

	if board[0][0] != "" && board[0][0] == board[1][1] && board[1][1] == board[2][2] {
		return board[0][0], true
	}
	if board[0][2] != "" && board[0][2] == board[1][1] && board[1][1] == board[2][0] {
		return board[0][2], true
	}

	return "", false
}

func EditBoardEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, gameId string) {
	game, exists := Games[gameId]

	if !exists {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "404, game not found. prob a bug",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	embed := createBoardEmbed(game.Game)

	s.ChannelMessageEditEmbed(game.ChannelId, game.BoardId, embed)

}

func EditMessageBoardEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, gameId string) {
	game := Games[gameId]

	var userTurn string
	if game.Turn == "X" {
		userTurn = game.PlayerX.Id
	} else {
		userTurn = game.PlayerY.Id
	}

	s.ChannelMessageEdit(i.ChannelID, game.BoardMessageId, fmt.Sprintf("<@%s> it's your turn right now", userTurn))
}

func CheckDraw(board Board) bool {
	for _, row := range board {
		for _, cell := range row {
			if cell == "" {
				return false
			}
		}
	}
	return true
}

func RandomizeTurn() string {
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(2) == 0 {
		return "X"
	}
	return "O"
}

func StartGame(s *discordgo.Session, i *discordgo.InteractionCreate, gameId string) {
	game := Games[gameId]
	board, _ := s.ChannelMessageSendEmbed(i.ChannelID, createBoardEmbed(Games[gameId].Game))

	game.BoardId = board.ID

	var userTurn string
	if game.Turn == "X" {
		userTurn = game.PlayerX.Id
	} else {
		userTurn = game.PlayerY.Id
	}

	boardMessage, _ := s.ChannelMessageSend(i.ChannelID, fmt.Sprintf("<@%s> its your turn right now", userTurn))
	Games[gameId].BoardMessageId = boardMessage.ID
}

func createBoardEmbed(board Board) *discordgo.MessageEmbed {
	var boardStr strings.Builder
	for _, row := range board {
		for _, cell := range row {
			if cell == "" {
				boardStr.WriteString("🟦")
			} else if cell == "X" {
				boardStr.WriteString("🔼")
			} else if cell == "O" {
				boardStr.WriteString("⏺")
			}
		}
		boardStr.WriteString("\n")
	}

	embed := &discordgo.MessageEmbed{
		Title:       "the duel",
		Description: boardStr.String(),
		Color:       0x00ff00,
	}

	return embed
}

func PlaceMarker(s *discordgo.Session, i *discordgo.InteractionCreate, gameId string, row int, col int) {
	Games[gameId].Game[row][col] = Games[gameId].Turn
}
