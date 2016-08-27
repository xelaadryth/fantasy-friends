package fantasy

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/TrevorSStone/goriot"
)

//RankedPrefix is a prefix used for all ranked queues for Game.SubType
const RankedPrefix = "RANKED_"

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
	Total       float32
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

	score.Total = score.Kills + score.Deaths + score.Assists + score.CS + score.TenKA + score.TripleKills +
		score.QuadraKills + score.Pentakills

	return score
}

//PlayerScoreFromGame calculates a score for a given game
func PlayerScoreFromGame(stats goriot.GameStat) PlayerScore {
	return CreatePlayerScore(stats.ChampionsKilled, stats.NumDeaths, stats.Assists, stats.MinionsKilled,
		stats.TripleKills, stats.QuadraKills, stats.PentaKills)
}

//CalculateScores returns a score map for each summonerID passed in
func CalculateScores(region string, summonerIDs []int64) ([]PlayerScore, error) {
	scores := make([]PlayerScore, len(summonerIDs), len(summonerIDs))
	for i := 0; i < len(summonerIDs); i++ {
		fmt.Println("Player", i)
		games, err := goriot.RecentGameBySummoner(region, summonerIDs[i])
		time.Sleep(2 * time.Second)
		if err != nil {
			return scores, err
		}
		if len(games) == 0 {
			return scores, errors.New(fmt.Sprint("No games in history for summonerID ", summonerIDs[i]))
		}

		for j := len(games) - 1; j >= 0; j-- {
			if !games[j].Invalid && strings.HasPrefix(games[j].SubType, RankedPrefix) {
				scores[i] = PlayerScoreFromGame(games[j].Statistics)
			}
		}
	}
	return scores, nil
}
