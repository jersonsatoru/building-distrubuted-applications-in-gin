package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"os"

	auth0 "github.com/auth0-community/go-auth0"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jersonsatoru/building-distributed-applications-in-gin/cmd/http/handlers"
	jose "gopkg.in/square/go-jose.v2"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		claims := &handlers.Claims{}
		tkn, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		if tkn == nil || !tkn.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Next()
	}
}

func AuthSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionToken := session.Get("token")
		log.Println(sessionToken)
		if sessionToken == nil {
			c.JSON(http.StatusForbidden, "not logged")
			c.Abort()
		}
		c.Next()
	}
}

func OAuth2Session() gin.HandlerFunc {
	return func(c *gin.Context) {
		var auth0Domain = fmt.Sprintf("https://%s/", os.Getenv("AUTH0_DOMAIN"))
		client := auth0.NewJWKClient(auth0.JWKClientOptions{
			URI: auth0Domain + ".well-known/jwks.json",
		}, nil)
		configuration := auth0.NewConfiguration(
			client,
			[]string{os.Getenv("AUTH0_API_IDENTIFIER")},
			auth0Domain,
			jose.RS256)
		validator := auth0.NewValidator(configuration, nil)
		_, err := validator.ValidateRequest(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "invalid token")
			c.Abort()
			return
		}
		c.Next()
	}
}
