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

// swagger:operation PUT /recipes recipe updateRecipe
// Update a recipe
// ---
// parameters:
// - in: path
//   name: id
//   type: string
//   description: ID of the recipe
//   required: true
// produces:
// - application/json
// responses:
//   '200':
//     description: Successful operation
func (app *RecipeHandler) UpdateRecipeHandler(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	var input models.Recipe
	err = c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	var recipe models.Recipe
	singleResult := app.collection.FindOne(app.ctx, bson.M{"_id": id})
	if err := singleResult.Err(); err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			c.JSON(http.StatusNotFound, "recipe not found")
		default:
			c.JSON(http.StatusInternalServerError, err)
		}
		return
	}
	err = singleResult.Decode(&recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	_, err = app.collection.UpdateOne(app.ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"name":         input.Name,
		"tags":         input.Tags,
		"instructions": input.Instructions,
		"ingredients":  input.Ingredients,
	}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"recipe": recipe})
}
