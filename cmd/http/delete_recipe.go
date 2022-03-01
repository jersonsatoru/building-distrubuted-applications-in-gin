package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jersonsatoru/building-distributed-applications-in-gin/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// swagger:operation DELETE /recipes/{id} recipe deleteRecipe
// Delete a recipe by its ID
//---
// parameters:
// - name: id
//   in: path
//   description: ID of the recipe
//   required: true
//   type: string
// produces:
// - application/json
// responses:
//   '204':
//     description: Successful operation
func (app *RecipeHandler) DeleteRecipeHandler(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	singleRow := app.collection.FindOne(app.ctx, bson.M{"_id": id})
	if err := singleRow.Err(); err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			c.JSON(http.StatusNotFound, "recipe not found")
		default:
			c.JSON(http.StatusInternalServerError, err)
		}
		return
	}
	var recipe models.Recipe
	err = singleRow.Decode(&recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	app.collection.DeleteOne(app.ctx, bson.M{"_id": id})
	c.Writer.WriteHeader(http.StatusNoContent)
}
