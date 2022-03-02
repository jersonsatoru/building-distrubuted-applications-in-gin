package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// swagger:operation /signout token signout
// Log out and clear the session client
// ---
// produces:
// - application/json
// responses:
//   '200':
//     description: Successful operation
func (h *AuthHandler) SignoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Writer.WriteHeader(http.StatusOK)
}
