package rgapi

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	SoloQueue = 420
	FlexQueue = 440
)

func GetRankedQueues() []int {
	return []int{SoloQueue, FlexQueue}
}

type ParticipantStats struct {
	PhysicalDamageDealt             int64 `json:"physicalDamageDealt"`
	NeutralMinionsKilledTeamJungle  int   `json:"neutralMinionsKilledTeamJungle"`
	MagicDamageDealt                int64 `json:"magicDamageDealt"`
	TotalPlayerScore                int   `json:"totalPlayerScore"`
	Deaths                          int   `json:"deaths"`
	Win                             bool  `json:"win"`
	NeutralMinionsKilledEnemyJungle int   `json:"neutralMinionsKilledEnemyJungle"`
	AltarsCaptured                  int   `json:"altarsCaptured"`
	LargestCriticalStrike           int   `json:"largestCriticalStrike"`
	TotalDamageDealt                int64 `json:"totalDamageDealt"`
	MagicDamageDealtToChampions     int64 `json:"magicDamageDealtToChampions"`
	VisionWardsBoughtInGame         int   `json:"visionWardsBoughtInGame"`
	DamageDealtToObjectives         int64 `json:"damageDealtToObjectives"`
	LargestKillingSpree             int   `json:"largestKillingSpree"`
	Item1                           int   `json:"item1"`
	QuadraKills                     int   `json:"quadraKills"`
	TeamObjective                   int   `json:"teamObjective"`
	TotalTimeCrowdControlDealt      int   `json:"totalTimeCrowdControlDealt"`
	LongestTimeSpentLiving          int   `json:"longestTimeSpentLiving"`
	WardsKilled                     int   `json:"wardsKilled"`
	FirstTowerAssist                bool  `json:"firstTowerAssist"`
	FirstTowerKill                  bool  `json:"firstTowerKill"`
	Item2                           int   `json:"item2"`
	Item3                           int   `json:"item3"`
	Item0                           int   `json:"item0"`
	FirstBloodAssist                bool  `json:"firstBloodAssist"`
	VisionScore                     int64 `json:"visionScore"`
	WardsPlaced                     int   `json:"wardsPlaced"`
	Item4                           int   `json:"item4"`
	Item5                           int   `json:"item5"`
	Item6                           int   `json:"item6"`
	TurretKills                     int   `json:"turretKills"`
	TripleKills                     int   `json:"tripleKills"`
	DamageSelfMitigated             int64 `json:"damageSelfMitigated"`
	ChampLevel                      int   `json:"champLevel"`
	NodeNeutralizeAssist            int   `json:"nodeNeutralizeAssist"`
	FirstInhibitorKill              bool  `json:"firstInhibitorKill"`
	GoldEarned                      int   `json:"goldEarned"`
	MagicalDamageTaken              int64 `json:"magicalDamageTaken"`
	Kills                           int   `json:"kills"`
	DoubleKills                     int   `json:"doubleKills"`
	NodeCaptureAssist               int   `json:"nodeCaptureAssist"`
	TrueDamageTaken                 int64 `json:"trueDamageTaken"`
	NodeNeutralize                  int   `json:"nodeNeutralize"`
	FirstInhibitorAssist            bool  `json:"firstInhibitorAssist"`
	Assists                         int   `json:"assists"`
	UnrealKills                     int   `json:"unrealKills"`
	NeutralMinionsKilled            int   `json:"neutralMinionsKilled"`
	ObjectivePlayerScore            int   `json:"objectivePlayerScore"`
	CombatPlayerScore               int   `json:"combatPlayerScore"`
	DamageDealtToTurrets            int64 `json:"damageDealtToTurrets"`
	AltarsNeutralized               int   `json:"altarsNeutralized"`
	PhysicalDamageDealtToChampions  int64 `json:"physicalDamageDealtToChampions"`
	GoldSpent                       int   `json:"goldSpent"`
	TrueDamageDealt                 int64 `json:"trueDamageDealt"`
	TrueDamageDealtToChampions      int64 `json:"trueDamageDealtToChampions"`
	ParticipantId                   int   `json:"participantId"`
	PentaKills                      int   `json:"pentaKills"`
	TotalHeal                       int64 `json:"totalHeal"`
	TotalMinionsKilled              int   `json:"totalMinionsKilled"`
	FirstBloodKill                  bool  `json:"firstBloodKill"`
	NodeCapture                     int   `json:"nodeCapture"`
	LargestMultiKill                int   `json:"largestMultiKill"`
	SightWardsBoughtInGame          int   `json:"sightWardsBoughtInGame"`
	TotalDamageDealtToChampions     int64 `json:"totalDamageDealtToChampions"`
	TotalUnitsHealed                int   `json:"totalUnitsHealed"`
	InhibitorKills                  int   `json:"inhibitorKills"`
	TotalScoreRank                  int   `json:"totalScoreRank"`
	TotalDamageTaken                int64 `json:"totalDamageTaken"`
	KillingSprees                   int   `json:"killingSprees"`
	TimeCCingOthers                 int64 `json:"timeCCingOthers"`
	PhysicalDamageTaken             int64 `json:"physicalDamageTaken"`
}

type Player struct {
	CurrentPlatformId string `json:"currentPlatformId"`
	SummonerName      string `json:"summonerName"`
	MatchHistoryUri   string `json:"matchHistoryUri"`
	PlatformId        string `json:"platformId"`
	CurrentAccountId  int64  `json:"currentAccountId"`
	ProfileIcon       int    `json:"profileIcon"`
	SummonerId        int64  `json:"summonerId"`
	AccountId         int64  `json:"accountId"`
}

type ParticipantIdentity struct {
	Player        Player `json:"player"`
	ParticipantId int    `json:"participantId"`
}

type TeamStats struct {
	FirstDragon          bool       `json:"firstDragon"`
	FirstInhibitor       bool       `json:"firstInhibitor"`
	Bans                 []TeamBans `json:"bans"`
	BaronKills           int        `json:"baronKills"`
	FirstRiftHerald      bool       `json:"firstRiftHerald"`
	FirstBaron           bool       `json:"firstBaron"`
	RiftHeraldKills      int        `json:"riftHeraldKills"`
	FirstBlood           bool       `json:"firstBlood"`
	TeamId               int        `json:"teamId"`
	FirstTower           bool       `json:"firstTower"`
	VilemawKills         int        `json:"vilemawKills"`
	InhibitorKills       int        `json:"inhibitorKills"`
	TowerKills           int        `json:"towerKills"`
	DominionVictoryScore int        `json:"dominionVictoryScore"`
	Win                  string     `json:"win"`
	DragonKills          int        `json:"dragonKills"`
}

type TeamBans struct {
	PickTurn   int `json:"pickTurn"`
	ChampionId int `json:"championId"`
}

type Rune struct {
	RuneId int `json:"rune"`
	Rank   int `json:"rank"`
}

type ParticipantTimeline struct {
	Lane                        string             `json:"lane"`
	ParticipantId               int                `json:"participantId"`
	CsDiffPerMinDeltas          map[string]float64 `json:"csDiffPerMinDeltas"`
	GoldPerMinDeltas            map[string]float64 `json:"goldPerMinDeltas"`
	XpDiffPerMinDeltas          map[string]float64 `json:"xpDiffPerMinDeltas"`
	CreepsPerMinDeltas          map[string]float64 `json:"creepsPerMinDeltas"`
	XpPerMinDeltas              map[string]float64 `json:"xpPerMinDeltas"`
	Role                        string             `json:"role"`
	DamageTakenDiffPerMinDeltas map[string]float64 `json:"damageTakenDiffPerMinDeltas"`
	DamageTakenPerMinDeltas     map[string]float64 `json:"damageTakenPerMinDeltas"`
}

type Mastery struct {
	MasteryId int `json:"masteryId"`
	Rank      int `json:"rank"`
}

type Participant struct {
	Stats                     ParticipantStats    `json:"stats"`
	ParticipantId             int                 `json:"participantId"`
	Runes                     []Rune              `json:"runes"`
	Timeline                  ParticipantTimeline `json:"timeline"`
	TeamId                    int                 `json:"teamId"`
	Spell2Id                  int                 `json:"spell2Id"`
	Masteries                 []Mastery           `json:"masteries"`
	HighestAchievedSeasonTier string              `json:"highestAchievedSeasonTier"`
	Spell1Id                  int                 `json:"spell1Id"`
	ChampionId                int                 `json:"championId"`
}

type Match struct {
	SeasonId              int                   `json:"seasonId"`
	QueueId               int                   `json:"queueId"`
	GameId                int64                 `json:"gameId"`
	ParticipantIdentities []ParticipantIdentity `json:"participantIdentities"`
	GameVersion           string                `json:"gameVersion"`
	PlatformId            string                `json:"platformId"`
	GameMode              string                `json:"gameMode"`
	MapId                 int                   `json:"mapId"`
	GameType              string                `json:"gameType"`
	Teams                 []TeamStats           `json:"teams"`
	Participants          []Participant         `json:"participants"`
	GameDuration          int64                 `json:"gameDuration"`
	GameCreation          int64                 `json:"gameCreation"`
}

type MatchReference struct {
	Lane       string `json:"lane"`
	GameId     int64  `json:"gameId"`
	Champion   int    `json:"champion"`
	PlatformId string `json:"platformId"`
	Timestamp  int64  `json:"timestamp"`
	Queue      int    `json:"queue"`
	Role       string `json:"role"`
	Season     int    `json:"season"`
}

type Matchlist struct {
	Matches    []MatchReference `json:"matches"`
	TotalGames int              `json:"totalGames"`
	StartIndex int              `json:"startIndex"`
	EndIndex   int              `json:"endIndex"`
}

//MatchDetails gets specific stats on a match
func MatchDetails(region string, matchId int64) (match *Match, err error) {
	match = &(Match{})
	err = apiGet(region, fmt.Sprintf("/match/v3/matches/%d", matchId), match)
	return
}

//FilterMatchlist pulls your 100 most recent games that fulfill the given constraints; unfortunately lane isn't one
func FilterMatchlist(region string, accountId int64, numGames int, queues []int) (matchlist Matchlist, err error) {
	queuesParams := ""
	if len(queues) > 0 {
		strQueues := make([]string, len(queues))
		for _, queue := range queues {
			strQueues = append(strQueues, strconv.Itoa(queue))
		}
		queuesParams = "&queues=%s" + strings.Join(strQueues, "&queues=")
	}
	err = apiGet(region, fmt.Sprintf("/match/v3/matchlists/by-account/%d?endIndex=%d%s", accountId, numGames, queuesParams), &matchlist)
	return
}

//RecentMatchlist pulls your 20 most recent games
func RecentMatchlist(region string, accountId int64) (matchlist Matchlist, err error) {
	err = apiGet(region, fmt.Sprintf("/match/v3/matchlists/by-account/%d/recent", accountId), &matchlist)
	return
}
