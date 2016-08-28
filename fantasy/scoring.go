package fantasy

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/TrevorSStone/goriot"
)

//RankedPrefix is a prefix used for all ranked queues for Game.SubType
const RankedPrefix = "RANKED_"

//Used to identify the teams
const (
	Order = "Order"
	Chaos = "Chaos"
)

//PlayerScore tracks the points scored for each different kind of point event
type PlayerScore struct {
	SummonerName string
	Kills        float32
	Deaths       float32
	Assists      float32
	CS           float32
	TenKA        float32
	TripleKills  float32
	QuadraKills  float32
	Pentakills   float32
	Score        float32
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
func CreatePlayerScore(summoner goriot.Summoner, kills int, deaths int, assists int, cs int, tripleKills int, quadraKills int, pentaKills int) PlayerScore {
	score := PlayerScore{
		SummonerName: summoner.Name,
		Kills:        float32(kills) * PointValues[KillsString],
		Deaths:       float32(deaths) * PointValues[DeathsString],
		Assists:      float32(assists) * PointValues[AssistsString],
		CS:           float32(cs) * PointValues[CSString],
		TripleKills:  float32(tripleKills) * PointValues[TripleKillsString],
		QuadraKills:  float32(quadraKills) * PointValues[QuadraKillsString],
		Pentakills:   float32(pentaKills) * PointValues[PentakillsString],
	}

	//Bonus points for getting at least 10 kills/assists in the same game
	if kills+assists >= 10 {
		score.TenKA = PointValues[TenKAString]
	}

	score.Score = score.Kills + score.Deaths + score.Assists + score.CS + score.TenKA + score.TripleKills +
		score.QuadraKills + score.Pentakills

	return score
}

//PlayerScoreFromGame calculates a score for a given game
func PlayerScoreFromGame(summoner goriot.Summoner, stats goriot.GameStat) PlayerScore {
	return CreatePlayerScore(summoner, stats.ChampionsKilled, stats.NumDeaths, stats.Assists, stats.MinionsKilled,
		stats.TripleKills, stats.QuadraKills, stats.PentaKills)
}

//CalculateScores returns a match score for the match played by the summoners passed in
func CalculateScores(region string, summoners []goriot.Summoner) (MatchScore, error) {
	//TODO: Split this function into helper methods

	//For each summoner, get the best game they have in their recent history
	//TODO: Save these games to DB
	playerScores := make([]PlayerScore, len(summoners), len(summoners))
	for i := 0; i < len(summoners); i++ {
		log.Println(i+1, "- summonerID", summoners[i].ID)
		games, err := goriot.RecentGameBySummoner(region, summoners[i].ID)
		if err != nil {
			return MatchScore{}, err
		}

		//Find the best ranked game in a list of recent games
		foundFlag := false
		for j := 0; j <= len(games); j++ {
			if !games[j].Invalid && strings.HasPrefix(games[j].SubType, RankedPrefix) {
				playerScores[i] = PlayerScoreFromGame(summoners[i], games[j].Statistics)
				foundFlag = true
				break
			}
		}
		if !foundFlag {
			return MatchScore{}, errors.New(fmt.Sprint("No ranked games found for ", summoners[i]))
		}
	}

	//Create teams and match score objects
	orderScore := TeamScore{Top: playerScores[0], Jungle: playerScores[1], Mid: playerScores[2], Bottom: playerScores[3], Support: playerScores[4], Score: playerScores[0].Score + playerScores[1].Score + playerScores[2].Score + playerScores[3].Score + playerScores[4].Score}
	chaosScore := TeamScore{Top: playerScores[5], Jungle: playerScores[6], Mid: playerScores[7], Bottom: playerScores[8], Support: playerScores[9], Score: playerScores[5].Score + playerScores[6].Score + playerScores[7].Score + playerScores[8].Score + playerScores[9].Score}
	matchScore := MatchScore{Order: orderScore, Chaos: chaosScore}

	//Pick the winner
	//TODO: Do this by user-selected team name instead of Order/Chaos
	if orderScore.Score > chaosScore.Score {
		matchScore.Winner = Order
	} else if chaosScore.Score > orderScore.Score {
		matchScore.Winner = Chaos
	}

	return matchScore, nil
}
