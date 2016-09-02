package controller

import (
	"log"
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
	cookieStore.Options(sessions.Options{
		MaxAge:   60 * 60 * 24 * 30,
		Secure:   !gin.IsDebugging(),
		HttpOnly: true,
	})
	router.Use(sessions.Sessions("fantasy-friends", cookieStore))
}

func sessionAsMap(session *sessions.Session) *map[string]interface{} {
	sessionMap := make(map[string]interface{})
	for _, fieldName := range sessionFields {
		sessionMap[fieldName] = (*session).Get(fieldName)
	}

	return &sessionMap
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
	router.GET("/", routeHome)
	router.GET("/about", routeAbout)
	//User Accounts ======================================================================================================
	router.GET("/login", routeLogin)
	router.POST("/login", processUser)
	router.GET("/logout", logout)
	//Fantasy ============================================================================================================
	router.GET("/team", routeTeam)
	router.POST("/team", saveTeam)
	router.POST("/matchResults", playMatch)
	//====================================================================================================================

	router.Run(":" + port)
}
