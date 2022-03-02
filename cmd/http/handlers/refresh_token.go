package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// swagger:operation POST /refresh token refreshToken
// Refresh JWT Token
// --
// produces:
// - application/jons
// responses:
//   '200':
//     description: Successful operation
func (h *AuthHandler) RefreshTokenHandler(c *gin.Context) {
	bearer := c.GetHeader("Authorization")
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(bearer, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if token == nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, "invalid token")
		return
	}
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		c.JSON(http.StatusUnauthorized, "token expired")
		return
	}
	expirationTime := time.Now().Add(15 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(os.Getenv("JWT_SECRET"))
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
