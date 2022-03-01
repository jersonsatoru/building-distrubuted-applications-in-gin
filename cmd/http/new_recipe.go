package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jersonsatoru/building-distributed-applications-in-gin/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// swagger:operation POST /recipes recipe createRecipe
// Create a new recipe
// ---
// produces:
// - apllication/json
// responses:
//   '201':
//     description: Successful operation
func (app *RecipeHandler) NewRecipeHandler(c *gin.Context) {
	var recipe models.Recipe
	err := c.ShouldBindJSON(&recipe)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	_, err = app.collection.InsertOne(app.ctx, recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	app.redisClient.Del("recipes")
	c.Writer.Header().Set("Location", fmt.Sprintf("/recipes/%s", recipe.ID))
	c.JSON(http.StatusCreated, map[string]interface{}{
		"recipe": recipe,
	})
}
