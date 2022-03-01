package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jersonsatoru/building-distributed-applications-in-gin/internal/models"
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
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if user.Username != "admin" || user.Password != "admin" {
		c.JSON(http.StatusUnauthorized, "invalid username and password")
		return
	}

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
