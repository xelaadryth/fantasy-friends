package controller

//PlayMatch runs a match from TODO: Split this into controller and fantasy
import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xelaadryth/fantasy-friends/fantasy"
)

//MatchForm input fields required for running a match
type MatchForm struct {
	Region       string `form:"region" binding:"required"`
	OrderTop     string `form:"orderTop" binding:"required"`
	OrderJungle  string `form:"orderJungle" binding:"required"`
	OrderMid     string `form:"orderMid" binding:"required"`
	OrderBottom  string `form:"orderBottom" binding:"required"`
	OrderSupport string `form:"orderSupport" binding:"required"`
	ChaosTop     string `form:"chaosTop" binding:"required"`
	ChaosJungle  string `form:"chaosJungle" binding:"required"`
	ChaosMid     string `form:"chaosMid" binding:"required"`
	ChaosBottom  string `form:"chaosBottom" binding:"required"`
	ChaosSupport string `form:"chaosSupport" binding:"required"`
}

func playMatch(c *gin.Context) {
	var matchForm MatchForm
	err := c.Bind(&matchForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
			"form":   matchForm,
		})
		return
	}

	matchScore, err := fantasy.PlayMatch(
		matchForm.Region,
		matchForm.OrderTop,
		matchForm.OrderJungle,
		matchForm.OrderMid,
		matchForm.OrderBottom,
		matchForm.OrderSupport,
		matchForm.ChaosTop,
		matchForm.ChaosJungle,
		matchForm.ChaosMid,
		matchForm.ChaosBottom,
		matchForm.ChaosSupport,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	//TODO: Figure out a way to round floats to 2 decimal places
	c.HTML(http.StatusOK, "match.tmpl", gin.H{
		"matchScore": matchScore,
	})
}
