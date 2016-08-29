package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//invalidHandler takes a non-nil error
//TODO: Display a nicely formatted page with link to go back to previous page; may require sessions implemented
func invalidHandler(c *gin.Context, httpStatusCode int, err error) {
	if err == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  "Invalid error.",
		})
	}
	c.JSON(httpStatusCode, gin.H{
		"status": httpStatusCode,
		"error":  err.Error(),
	})
}
