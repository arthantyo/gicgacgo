package shared

import (
	"fmt"
	"sort"
)

type PlayerStats struct {
	UserId   string
	Username string
	Wins     int
	Losses   int
	Draws    int
}

type GuildStats struct {
	GuildId string
	Stats   map[string]*PlayerStats // userId -> stats
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
	} else {
		GlobalStats[winnerId].Wins++
		GlobalStats[loserId].Losses++
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
	} else {
		guildStats[winnerId].Wins++
		guildStats[loserId].Losses++
	}
}

// GetGlobalLeaderboard returns sorted global stats
func GetGlobalLeaderboard(limit int) []*PlayerStats {
	stats := make([]*PlayerStats, 0, len(GlobalStats))
	for _, ps := range GlobalStats {
		stats = append(stats, ps)
	}

	sort.Slice(stats, func(i, j int) bool {
		// Sort by wins first, then by win rate
		if stats[i].Wins != stats[j].Wins {
			return stats[i].Wins > stats[j].Wins
		}
		return stats[i].WinRate() > stats[j].WinRate()
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
		// Sort by wins first, then by win rate
		if stats[i].Wins != stats[j].Wins {
			return stats[i].Wins > stats[j].Wins
		}
		return stats[i].WinRate() > stats[j].WinRate()
	})

	if limit > 0 && len(stats) > limit {
		return stats[:limit]
	}
	return stats
}
