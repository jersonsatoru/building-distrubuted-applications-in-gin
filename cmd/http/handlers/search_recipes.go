package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jersonsatoru/building-distributed-applications-in-gin/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

// swagger:operation GET /recipes/search recipe searchRecipe
// Search for recipes filtered by tag
// ---
// parameters:
// - in: query
//   name: tag
//   description: Recipe's tag that you are looking for
//   type: string
//   required: true
// produces:
// - application/json
// responses:
//   '200':
//     description: Successful operation
func (app *RecipeHandler) SearchRecipesHandler(c *gin.Context) {
	searchedTag := c.Query("tag")
	cursor, err := app.collection.Find(app.ctx, bson.M{"tags": bson.M{"$in": []string{searchedTag}}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	recipes := make([]models.Recipe, 0)
	defer cursor.Close(app.ctx)
	for cursor.Next(app.ctx) {
		var recipe models.Recipe
		err := cursor.Decode(&recipe)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		recipes = append(recipes, recipe)
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"data": recipes,
	})
}
