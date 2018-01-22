package fantasy

//PlayersPerMatch is the number of players in a fantasy match
const PlayersPerMatch = 10

//An enumeration of different point events used to calculate score
const (
	KillsString       = "Kills"
	DeathsString      = "Deaths"
	AssistsString     = "Assists"
	CSString          = "CS"
	TenKAString       = "10+ K/A"
	TripleKillsString = "Triple Kills"
	QuadraKillsString = "Quadra Kills"
	PentakillsString  = "Pentakills"
)

//PointValues is a mapping from "eventType" to fantasy point value
var PointValues = map[string]float32{
	KillsString:       2,
	DeathsString:      -0.5,
	AssistsString:     1.5,
	CSString:          0.01,
	TenKAString:       2,
	TripleKillsString: 2,
	QuadraKillsString: 5,
	PentakillsString:  10,
}
