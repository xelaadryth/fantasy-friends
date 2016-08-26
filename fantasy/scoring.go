package fantasy

//PointValues is a mapping from "eventType" to fantasy point value
var PointValues = map[string]float64{
	"Kills":        2,
	"Deaths":       -0.5,
	"Assists":      1.5,
	"CS":           0.01,
	"10+ K/A":      2,
	"Triple Kills": 2,
	"Quadra Kills": 5,
	"Pentakills":   10}
