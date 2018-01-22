package fantasy

import (
	"errors"
	"fmt"
	"log"

	"github.com/xelaadryth/fantasy-friends/rgapi"
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
	Top         *PlayerScore
	Jungle      *PlayerScore
	Mid         *PlayerScore
	Bottom      *PlayerScore
	Support     *PlayerScore
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
func CreatePlayerScore(gameName string, kills int, deaths int, assists int, cs int, tripleKills int, quadraKills int, pentaKills int) *PlayerScore {
	score := PlayerScore{
		SummonerName: gameName,
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

	return &score
}

//PlayerScoreFromGame calculates a score for a given game
func PlayerScoreFromGame(region string, summoner rgapi.Summoner, matchId int64) (score *PlayerScore, err error) {
	match, err := rgapi.MatchDetails(region, matchId)
	if err != nil {
		return
	}

	//Only grab the stats of the player we care about
	participantId := -1
	err = errors.New(fmt.Sprintf("Summoner %s with account id %d not found in match %d", summoner.Name, summoner.AccountId, matchId))
	for _, participantIdentity := range match.ParticipantIdentities {
		if participantIdentity.Player.AccountId == summoner.AccountId {
			participantId = participantIdentity.ParticipantId
			err = nil
			break
		}
	}
	if err != nil {
		log.Println("Match data:", match)
		return
	}

	stats := rgapi.ParticipantStats{}
	err = errors.New(fmt.Sprintf("Participant with id %d not found in match %d", participantId, matchId))
	for _, participant := range match.Participants {
		if participant.ParticipantId == participantId {
			stats = participant.Stats
			err = nil
			break
		}
	}
	if err != nil {
		return
	}

	score = CreatePlayerScore(summoner.Name, stats.Kills, stats.Deaths, stats.Assists, stats.TotalMinionsKilled,
		stats.TripleKills, stats.QuadraKills, stats.PentaKills)

	return
}

//PlayerScoreBestRecent gets the highest score from recent games
func PlayerScoreBestRecent(region string, summoner rgapi.Summoner) (maxScore *PlayerScore, err error) {
	log.Println("Fetching recent games for player", summoner.Name)
	//TODO: Cache results and limit searches based on most recent timestamps checked
	//Currently grabs 5 (max 100) most recent ranked games, but should instead grab all games since a timestamp
	matchlist, err := rgapi.FilterMatchlist(region, summoner.AccountId, 5, rgapi.GetRankedQueues())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error getting data for summoner %s, try another summoner", summoner.Name))
	}

	//Find the best ranked game in a list of recent games
	//We should never go lower than this score in a game; definition of a minimum float is iffy
	maxScore = &PlayerScore{Score: -8192}
	err = errors.New(fmt.Sprint("No ranked games found for ", summoner))
	for _, match := range matchlist.Matches {
		score, scoreErr := PlayerScoreFromGame(region, summoner, match.GameId)
		if scoreErr != nil {
			return maxScore, scoreErr
		}
		if score.Score > maxScore.Score {
			maxScore = score
			err = nil
		}
	}
	return
}

//CalculateMatchScore returns a match score for the match played by the summoners passed in
func CalculateMatchScore(region string, summoners []rgapi.Summoner) (*MatchScore, error) {
	if len(summoners) != PlayersPerMatch {
		return nil, errors.New(fmt.Sprint(
			"Calculating match score requires ", PlayersPerMatch, " players, was given ", len(summoners)))
	}

	//For each summoner, get the best game they have in their recent history
	//TODO: Save timestamp of most recently checked games in DB or Redis
	playerScores := make([]*PlayerScore, len(summoners))
	for i := 0; i < len(summoners); i++ {
		playerScore, err := PlayerScoreBestRecent(region, summoners[i])
		if err != nil {
			return nil, err
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

	return &matchScore, nil
}
