package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jersonsatoru/building-distributed-applications-in-gin/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

// swagger:operation GET /recipes recipe listRecipe
// Returns list of recipes
// ---
// produces:
// - application/json
// responses:
//   '200':
//     description: Successful operation
func (app *RecipeHandler) ListRecipesHandler(c *gin.Context) {
	cur, err := app.collection.Find(app.ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	defer cur.Close(app.ctx)
	recipes := make([]models.Recipe, 0)
	for cur.Next(app.ctx) {
		var recipe models.Recipe
		if err := cur.Decode(&recipe); err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		recipes = append(recipes, recipe)
	}
	c.JSON(http.StatusOK, gin.H{
		"data": recipes,
	})
}
