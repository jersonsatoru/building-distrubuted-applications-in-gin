package handlers

import (
	"crypto/sha256"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jersonsatoru/building-distributed-applications-in-gin/internal/models"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// swagger:operation POST /refresh_token token refreshToken
// Refresh JWT Token
// ---
// produces:
// - application/json
// responses:
//   '201':
//     description: Successful operation
func (handler AuthHandler) SignInHandler(c *gin.Context) {
	var input models.User
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	singleResult := handler.collection.FindOne(handler.ctx, bson.M{
		"username": input.Username,
	})
	var user models.User
	err = singleResult.Decode(&user)
	if err != nil {
		switch {
		case errors.Is(mongo.ErrNoDocuments, err):
			c.JSON(http.StatusUnauthorized, "invalid username and password")
		default:
			c.JSON(http.StatusInternalServerError, err)
		}
		return
	}
	sha := sha256.New()
	password := sha.Sum([]byte(input.Password))

	if user.Password != string(password) {
		c.JSON(http.StatusUnauthorized, "invalid username and password")
		return
	}

	sessionToken := xid.New().String()
	session := sessions.Default(c)
	session.Set("username", user.Username)
	session.Set("token", sessionToken)
	session.Save()

	expirationTime := time.Now().Add(time.Minute * 10)
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	jwtOutput := &JWTOutput{
		Token:   tokenString,
		Expires: expirationTime,
	}
	c.JSON(http.StatusOK, jwtOutput)
}
