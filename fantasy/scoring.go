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
	BlueTeam = "blue"
	RedTeam  = "red"
)

//PlayerScore tracks the points scored for each different kind of point event
type PlayerScore struct {
	SummonerName string
	Kills        float32
	Deaths       float32
	Assists      float32
	CS           float32
	CSString     string //Round to 2 decimal places
	TenKA        float32
	TripleKills  float32
	QuadraKills  float32
	Pentakills   float32
	Score        float32
	ScoreString  string //Round to 2 decimal places
}

//TeamScore has scores for the team and each team member
type TeamScore struct {
	Top         PlayerScore
	Jungle      PlayerScore
	Mid         PlayerScore
	Bottom      PlayerScore
	Support     PlayerScore
	Score       float32
	ScoreString string //Round to 2 decimal places
}

//MatchScore contains the subtotal and total fantasy points for a match
type MatchScore struct {
	BlueTeam    TeamScore
	RedTeam     TeamScore
	WinningSide string
	WinningTeam string
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
	score.CSString = fmt.Sprintf("%.2f", score.CS)

	//Bonus points for getting at least 10 kills/assists in the same game
	if kills+assists >= 10 {
		score.TenKA = PointValues[TenKAString]
	}

	score.Score = score.Kills + score.Deaths + score.Assists + score.CS + score.TenKA + score.TripleKills +
		score.QuadraKills + score.Pentakills
	score.ScoreString = fmt.Sprintf("%.2f", score.Score)

	return score
}

//PlayerScoreFromGame calculates a score for a given game
func PlayerScoreFromGame(summoner goriot.Summoner, stats goriot.GameStat) PlayerScore {
	return CreatePlayerScore(summoner, stats.ChampionsKilled, stats.NumDeaths, stats.Assists, stats.MinionsKilled,
		stats.TripleKills, stats.QuadraKills, stats.PentaKills)
}

//PlayerScoreBestRecent gets the highest score from recent games
func PlayerScoreBestRecent(region string, summoner goriot.Summoner) (PlayerScore, error) {
	log.Println("Fetching recent games for player", summoner.Name)
	games, err := goriot.RecentGameBySummoner(region, summoner.ID)
	if err != nil {
		return PlayerScore{}, err
	}

	//Find the best ranked game in a list of recent games
	playerScores := make([]PlayerScore, len(games), len(games))
	maxIndex := -1
	//TODO: See if there's a reasonable definition for a "minimum float32" but probably varies depending on architecture
	var maxScore float32 = -999999
	for i := 0; i < len(games); i++ {
		if !games[i].Invalid && strings.HasPrefix(games[i].SubType, RankedPrefix) {
			playerScores[i] = PlayerScoreFromGame(summoner, games[i].Statistics)
			if playerScores[i].Score > maxScore {
				maxScore = playerScores[i].Score
				maxIndex = i
			}
		}
	}
	if maxIndex >= 0 {
		return playerScores[maxIndex], nil
	}

	return PlayerScore{}, errors.New(fmt.Sprint("No ranked games found for ", summoner))
}

//CalculateMatchScore returns a match score for the match played by the summoners passed in
func CalculateMatchScore(region string, summoners []goriot.Summoner) (MatchScore, error) {
	if len(summoners) != PlayersPerMatch {
		return MatchScore{}, errors.New(fmt.Sprint(
			"Calculating match score requires ", PlayersPerMatch, " players, only given ", len(summoners)))
	}

	//For each summoner, get the best game they have in their recent history
	//TODO: Save these games to DB
	playerScores := make([]PlayerScore, len(summoners), len(summoners))
	for i := 0; i < len(summoners); i++ {
		playerScore, err := PlayerScoreBestRecent(region, summoners[i])
		if err != nil {
			return MatchScore{}, err
		}
		playerScores[i] = playerScore
	}

	//Create teams and match score objects
	blueTeamScore := TeamScore{
		Top:     playerScores[0],
		Jungle:  playerScores[1],
		Mid:     playerScores[2],
		Bottom:  playerScores[3],
		Support: playerScores[4],
		Score:   playerScores[0].Score + playerScores[1].Score + playerScores[2].Score + playerScores[3].Score + playerScores[4].Score,
	}
	blueTeamScore.ScoreString = fmt.Sprintf("%.2f", blueTeamScore.Score)
	redTeamScore := TeamScore{
		Top:     playerScores[5],
		Jungle:  playerScores[6],
		Mid:     playerScores[7],
		Bottom:  playerScores[8],
		Support: playerScores[9],
		Score:   playerScores[5].Score + playerScores[6].Score + playerScores[7].Score + playerScores[8].Score + playerScores[9].Score,
	}
	redTeamScore.ScoreString = fmt.Sprintf("%.2f", redTeamScore.Score)
	matchScore := MatchScore{
		BlueTeam: blueTeamScore,
		RedTeam:  redTeamScore,
	}

	//Pick the winner
	//TODO: Do this by user-selected team name instead of Blue/Red
	if blueTeamScore.Score > redTeamScore.Score {
		matchScore.WinningSide = BlueTeam
		matchScore.WinningTeam = "Blue"
	} else if redTeamScore.Score > blueTeamScore.Score {
		matchScore.WinningSide = RedTeam
		matchScore.WinningTeam = "Red"
	}

	return matchScore, nil
}
