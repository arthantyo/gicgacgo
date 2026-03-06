package shared

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
)

type PlayerStats struct {
	UserId   string
	Username string
	Wins     int
	Losses   int
	Draws    int
	Level    int
	XP       int
}

type Rank struct {
	Name  string
	Tier  int
	Icon  string
	MinXP int
}

type GuildStats struct {
	GuildId string
	Stats   map[string]*PlayerStats // userId -> stats
}

// XP rewards for different outcomes
const (
	XPPerWin  = 100
	XPPerDraw = 25
	XPPerLoss = 10
)

// Rank definitions with XP thresholds
var Ranks = []Rank{
	{Name: "Bronze", Tier: 1, Icon: "🥉", MinXP: 0},
	{Name: "Bronze", Tier: 2, Icon: "🥉", MinXP: 300},
	{Name: "Bronze", Tier: 3, Icon: "🥉", MinXP: 600},
	{Name: "Silver", Tier: 1, Icon: "🥈", MinXP: 1000},
	{Name: "Silver", Tier: 2, Icon: "🥈", MinXP: 1500},
	{Name: "Silver", Tier: 3, Icon: "🥈", MinXP: 2100},
	{Name: "Gold", Tier: 1, Icon: "🥇", MinXP: 2800},
	{Name: "Gold", Tier: 2, Icon: "🥇", MinXP: 3600},
	{Name: "Gold", Tier: 3, Icon: "🥇", MinXP: 4500},
	{Name: "Platinum", Tier: 1, Icon: "💎", MinXP: 5500},
	{Name: "Platinum", Tier: 2, Icon: "💎", MinXP: 6600},
	{Name: "Platinum", Tier: 3, Icon: "💎", MinXP: 7800},
	{Name: "Diamond", Tier: 1, Icon: "💠", MinXP: 9100},
	{Name: "Diamond", Tier: 2, Icon: "💠", MinXP: 10500},
	{Name: "Diamond", Tier: 3, Icon: "💠", MinXP: 12000},
	{Name: "Master", Tier: 1, Icon: "👑", MinXP: 13600},
	{Name: "Grandmaster", Tier: 1, Icon: "⭐", MinXP: 15300},
}

// Global stats across all guilds
var GlobalStats = make(map[string]*PlayerStats) // userId -> stats

// Guild-specific stats
var GuildLeaderboards = make(map[string]*GuildStats) // guildId -> GuildStats

func (ps *PlayerStats) TotalGames() int {
	return ps.Wins + ps.Losses + ps.Draws
}

func (ps *PlayerStats) WinRate() float64 {
	total := ps.TotalGames()
	if total == 0 {
		return 0
	}
	return float64(ps.Wins) / float64(total) * 100
}

func (ps *PlayerStats) WinRateString() string {
	return fmt.Sprintf("%.1f%%", ps.WinRate())
}

// GetRank returns the current rank based on XP
func (ps *PlayerStats) GetRank() Rank {
	currentRank := Ranks[0]
	for _, rank := range Ranks {
		if ps.XP >= rank.MinXP {
			currentRank = rank
		} else {
			break
		}
	}
	return currentRank
}

// GetRankString returns a formatted rank string
func (ps *PlayerStats) GetRankString() string {
	rank := ps.GetRank()
	return fmt.Sprintf("%s %s %d", rank.Icon, rank.Name, rank.Tier)
}

// GetNextRank returns the next rank and XP needed to reach it
func (ps *PlayerStats) GetNextRank() (Rank, int) {
	for i, rank := range Ranks {
		if ps.XP < rank.MinXP {
			xpNeeded := rank.MinXP - ps.XP
			return rank, xpNeeded
		}
		// If we're at the last rank and have enough XP
		if i == len(Ranks)-1 && ps.XP >= rank.MinXP {
			return rank, 0 // Already at max rank
		}
	}
	return Ranks[len(Ranks)-1], 0
}

// AddXP adds experience points and updates level
func (ps *PlayerStats) AddXP(xp int) {
	ps.XP += xp
	ps.Level = ps.XP / 100 // Simple level calculation: 1 level per 100 XP
}

// CalculateDynamicXP calculates XP based on winner and loser ranks
// Higher ranked players get less XP for beating lower ranked players
// Lower ranked players get more XP for beating higher ranked players
func CalculateDynamicXP(baseXP int, winnerXP int, loserXP int) (int, int) {
	winnerRankIndex := getRankIndex(winnerXP)
	loserRankIndex := getRankIndex(loserXP)
	
	rankDiff := loserRankIndex - winnerRankIndex
	
	// Calculate multiplier based on rank difference
	// If winner is higher ranked (rankDiff < 0), they get less XP
	// If winner is lower ranked (rankDiff > 0), they get more XP
	var winnerMultiplier, loserMultiplier float64
	
	if rankDiff == 0 {
		// Same rank - normal XP
		winnerMultiplier = 1.0
		loserMultiplier = 1.0
	} else if rankDiff < 0 {
		// Winner is higher ranked than loser
		// Winner gets less, loser gets normal
		winnerMultiplier = math.Max(0.3, 1.0 - float64(-rankDiff)*0.1)
		loserMultiplier = 1.0
	} else {
		// Winner is lower ranked than loser (upset!)
		// Winner gets bonus, loser gets less penalty
		winnerMultiplier = math.Min(2.0, 1.0 + float64(rankDiff)*0.15)
		loserMultiplier = 0.8
	}
	
	winnerXPGain := int(float64(baseXP) * winnerMultiplier)
	loserXPGain := int(float64(XPPerLoss) * loserMultiplier)
	
	return winnerXPGain, loserXPGain
}

// getRankIndex returns the index of the rank for a given XP
func getRankIndex(xp int) int {
	for i := len(Ranks) - 1; i >= 0; i-- {
		if xp >= Ranks[i].MinXP {
			return i
		}
	}
	return 0
}

// RecordGameResult updates both global and guild-specific stats
func RecordGameResult(winnerId, loserId, guildId, winnerName, loserName string, isDraw bool) {
	// Update global stats
	if GlobalStats[winnerId] == nil {
		GlobalStats[winnerId] = &PlayerStats{UserId: winnerId, Username: winnerName}
	}
	if GlobalStats[loserId] == nil {
		GlobalStats[loserId] = &PlayerStats{UserId: loserId, Username: loserName}
	}

	if isDraw {
		GlobalStats[winnerId].Draws++
		GlobalStats[loserId].Draws++
		GlobalStats[winnerId].AddXP(XPPerDraw)
		GlobalStats[loserId].AddXP(XPPerDraw)
	} else {
		GlobalStats[winnerId].Wins++
		GlobalStats[loserId].Losses++
		
		// Calculate dynamic XP based on rank difference
		winnerXPGain, loserXPGain := CalculateDynamicXP(XPPerWin, GlobalStats[winnerId].XP, GlobalStats[loserId].XP)
		GlobalStats[winnerId].AddXP(winnerXPGain)
		GlobalStats[loserId].AddXP(loserXPGain)
	}

	// Update guild stats
	if GuildLeaderboards[guildId] == nil {
		GuildLeaderboards[guildId] = &GuildStats{
			GuildId: guildId,
			Stats:   make(map[string]*PlayerStats),
		}
	}

	guildStats := GuildLeaderboards[guildId].Stats
	if guildStats[winnerId] == nil {
		guildStats[winnerId] = &PlayerStats{UserId: winnerId, Username: winnerName}
	}
	if guildStats[loserId] == nil {
		guildStats[loserId] = &PlayerStats{UserId: loserId, Username: loserName}
	}

	if isDraw {
		guildStats[winnerId].Draws++
		guildStats[loserId].Draws++
		guildStats[winnerId].AddXP(XPPerDraw)
		guildStats[loserId].AddXP(XPPerDraw)
	} else {
		guildStats[winnerId].Wins++
		guildStats[loserId].Losses++
		
		// Calculate dynamic XP for guild stats too
		winnerXPGain, loserXPGain := CalculateDynamicXP(XPPerWin, guildStats[winnerId].XP, guildStats[loserId].XP)
		guildStats[winnerId].AddXP(winnerXPGain)
		guildStats[loserId].AddXP(loserXPGain)
	}
	
	// Save stats after each game
	SaveStats()
}

// GetGlobalLeaderboard returns sorted global stats
func GetGlobalLeaderboard(limit int) []*PlayerStats {
	stats := make([]*PlayerStats, 0, len(GlobalStats))
	for _, ps := range GlobalStats {
		stats = append(stats, ps)
	}

	sort.Slice(stats, func(i, j int) bool {
		// Sort by XP/rank first, then by wins
		if stats[i].XP != stats[j].XP {
			return stats[i].XP > stats[j].XP
		}
		return stats[i].Wins > stats[j].Wins
	})

	if limit > 0 && len(stats) > limit {
		return stats[:limit]
	}
	return stats
}

// GetGuildLeaderboard returns sorted guild-specific stats
func GetGuildLeaderboard(guildId string, limit int) []*PlayerStats {
	guildStats, exists := GuildLeaderboards[guildId]
	if !exists || guildStats.Stats == nil {
		return []*PlayerStats{}
	}

	stats := make([]*PlayerStats, 0, len(guildStats.Stats))
	for _, ps := range guildStats.Stats {
		stats = append(stats, ps)
	}

	sort.Slice(stats, func(i, j int) bool {
		// Sort by XP/rank first, then by wins
		if stats[i].XP != stats[j].XP {
			return stats[i].XP > stats[j].XP
		}
		return stats[i].Wins > stats[j].Wins
	})

	if limit > 0 && len(stats) > limit {
		return stats[:limit]
	}
	return stats
}

// SaveData is the structure for JSON persistence
type SaveData struct {
	GlobalStats       map[string]*PlayerStats            `json:"global_stats"`
	GuildLeaderboards map[string]map[string]*PlayerStats `json:"guild_leaderboards"`
}

// SaveStats saves all stats to a JSON file
func SaveStats() error {
	// Convert GuildLeaderboards to a serializable format
	guildData := make(map[string]map[string]*PlayerStats)
	for guildId, guildStats := range GuildLeaderboards {
		guildData[guildId] = guildStats.Stats
	}

	saveData := SaveData{
		GlobalStats:       GlobalStats,
		GuildLeaderboards: guildData,
	}

	data, err := json.MarshalIndent(saveData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal stats: %w", err)
	}

	err = os.WriteFile("stats.json", data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write stats file: %w", err)
	}

	return nil
}

// LoadStats loads all stats from a JSON file
func LoadStats() error {
	data, err := os.ReadFile("stats.json")
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist yet, that's okay
			return nil
		}
		return fmt.Errorf("failed to read stats file: %w", err)
	}

	var saveData SaveData
	err = json.Unmarshal(data, &saveData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal stats: %w", err)
	}

	// Load global stats
	GlobalStats = saveData.GlobalStats
	if GlobalStats == nil {
		GlobalStats = make(map[string]*PlayerStats)
	}

	// Load guild stats
	GuildLeaderboards = make(map[string]*GuildStats)
	for guildId, stats := range saveData.GuildLeaderboards {
		GuildLeaderboards[guildId] = &GuildStats{
			GuildId: guildId,
			Stats:   stats,
		}
	}

	return nil
}
