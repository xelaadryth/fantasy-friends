package fantasy

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/TrevorSStone/goriot"
)

//RankedPrefix is a prefix used for all ranked queues for Game.SubType
const RankedPrefix = "RANKED_"

//Used to identify the teams
const (
	Order = "ORDER"
	Chaos = "CHAOS"
)

//PlayerScore tracks the points scored for each different kind of point event
type PlayerScore struct {
	Kills       float32
	Deaths      float32
	Assists     float32
	CS          float32
	TenKA       float32
	TripleKills float32
	QuadraKills float32
	Pentakills  float32
	Score       float32
}

//TeamScore has scores for the team and each team member
type TeamScore struct {
	Top     PlayerScore
	Jungle  PlayerScore
	Mid     PlayerScore
	Bottom  PlayerScore
	Support PlayerScore
	Score   float32
}

//MatchScore contains the subtotal and total fantasy points for a match
type MatchScore struct {
	Order  TeamScore
	Chaos  TeamScore
	Winner string
}

//CreatePlayerScore uses basic information to construct a PlayerScore object
func CreatePlayerScore(kills int, deaths int, assists int, cs int, tripleKills int, quadraKills int, pentaKills int) PlayerScore {
	score := PlayerScore{
		Kills:       float32(kills) * PointValues[KillsString],
		Deaths:      float32(deaths) * PointValues[DeathsString],
		Assists:     float32(assists) * PointValues[AssistsString],
		CS:          float32(cs) * PointValues[CSString],
		TripleKills: float32(tripleKills) * PointValues[TripleKillsString],
		QuadraKills: float32(quadraKills) * PointValues[QuadraKillsString],
		Pentakills:  float32(pentaKills) * PointValues[PentakillsString],
	}

	if kills+assists >= 10 {
		score.TenKA = PointValues[TenKAString]
	}

	score.Score = score.Kills + score.Deaths + score.Assists + score.CS + score.TenKA + score.TripleKills +
		score.QuadraKills + score.Pentakills

	return score
}

//PlayerScoreFromGame calculates a score for a given game
func PlayerScoreFromGame(stats goriot.GameStat) PlayerScore {
	return CreatePlayerScore(stats.ChampionsKilled, stats.NumDeaths, stats.Assists, stats.MinionsKilled,
		stats.TripleKills, stats.QuadraKills, stats.PentaKills)
}

//CalculateScores returns a score map for each summonerID passed in
func CalculateScores(region string, summonerIDs []int64) (MatchScore, error) {
	playerScores := make([]PlayerScore, len(summonerIDs), len(summonerIDs))
	for i := 0; i < len(summonerIDs); i++ {
		log.Println(i+1, "- summonerID", summonerIDs[i])
		games, err := goriot.RecentGameBySummoner(region, summonerIDs[i])
		time.Sleep(2 * time.Second)
		if err != nil {
			return MatchScore{}, err
		}
		if len(games) == 0 {
			return MatchScore{}, errors.New(fmt.Sprint("No games in history for summonerID ", summonerIDs[i]))
		}

		for j := len(games) - 1; j >= 0; j-- {
			if !games[j].Invalid && strings.HasPrefix(games[j].SubType, RankedPrefix) {
				playerScores[i] = PlayerScoreFromGame(games[j].Statistics)
			}
		}
	}

	//TODO: Split this function into helper methods
	orderScore := TeamScore{Top: playerScores[0], Jungle: playerScores[1], Mid: playerScores[2], Bottom: playerScores[3], Support: playerScores[4], Score: playerScores[0].Score + playerScores[1].Score + playerScores[2].Score + playerScores[3].Score + playerScores[4].Score}
	chaosScore := TeamScore{Top: playerScores[5], Jungle: playerScores[6], Mid: playerScores[7], Bottom: playerScores[8], Support: playerScores[9], Score: playerScores[5].Score + playerScores[6].Score + playerScores[7].Score + playerScores[8].Score + playerScores[9].Score}
	matchScore := MatchScore{Order: orderScore, Chaos: chaosScore}

	if orderScore.Score > chaosScore.Score {
		matchScore.Winner = Order
	} else if chaosScore.Score > orderScore.Score {
		matchScore.Winner = Chaos
	}

	return matchScore, nil
}
