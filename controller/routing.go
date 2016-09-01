package controller

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/contrib/secure"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func addMiddleware(router *gin.Engine) {
	//Redirect to HTTPS when not in debug mode
	router.Use(
		secure.Secure(secure.Options{
			SSLRedirect:          true,
			SSLProxyHeaders:      map[string]string{"X-Forwarded-Proto": "https"},
			STSSeconds:           315360000,
			STSIncludeSubdomains: true,
			FrameDeny:            true,
			ContentTypeNosniff:   true,
			BrowserXssFilter:     true,
			IsDevelopment:        gin.IsDebugging(),
		}))

	//Use cookies which expire in 30 days by default
	cookieStore := sessions.NewCookieStore(
		[]byte(os.Getenv("COOKIE_AUTH_KEY")),
		[]byte(os.Getenv("COOKIE_ENCRYPTION_KEY")),
		[]byte(os.Getenv("COOKIE_OLD_AUTH_KEY")),
		[]byte(os.Getenv("COOKIE_OLD_ENCRYPTION_KEY")),
	)
	router.Use(sessions.Sessions("fantasy-friends", cookieStore))
}

//Route does all the routing for the app
func Route() {
	//Get port number to listen for
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	//Set up routing
	router := gin.Default()
	addMiddleware(router)

	router.LoadHTMLGlob("templates/*.tmpl")
	router.Static("/static", "./static")

	//TODO: Split into routing groups
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"displayName": sessions.Default(c).Get("displayName"),
		})
	})
	router.GET("/about", func(c *gin.Context) {
		c.HTML(http.StatusOK, "about.tmpl", gin.H{
			"displayName": sessions.Default(c).Get("displayName"),
		})
	})

	//User Accounts ==============================================================================================================
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"displayName": sessions.Default(c).Get("displayName"),
		})
	})
	router.POST("/login", processUser)

	//Fantasy ====================================================================================================================
	router.POST("/matchResults", playMatch)
	router.Run(":" + port)
}
