package controller

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

//routeHome to home page
func routeHome(c *gin.Context) {
	session := sessions.Default(c)
	session.Set(sessionNavActive, "home")
	sessionMap := sessionAsMap(&session)
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		sessionSession: *sessionMap,
	})
}

func routeAbout(c *gin.Context) {
	session := sessions.Default(c)
	session.Set(sessionNavActive, "help")
	sessionMap := sessionAsMap(&session)
	c.HTML(http.StatusOK, "about.tmpl", gin.H{
		sessionSession: *sessionMap,
	})
}

func routeLogin(c *gin.Context) {
	session := sessions.Default(c)
	session.Set(sessionNavActive, "user")
	sessionMap := sessionAsMap(&session)
	c.HTML(http.StatusOK, "login.tmpl", gin.H{
		sessionSession: *sessionMap,
	})
}
