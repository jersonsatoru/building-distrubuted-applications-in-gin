// 	Recipes API
// 	This is a simple recipes API. You can find out more about the API at  google, thanks
//	Schemes: http
//	Host: localhost:8080
//	BasePath:
//	Version: 1.0.0
//	Contact: jersonsatoru@yahoo.com.br
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
// swagger:meta
package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

var (
	Recipes = make([]Recipe, 0)
)

func main() {
	r := gin.Default()
	r.POST("/recipes", NewRecipeHandler)
	r.GET("/recipes", ListRecipesHandler)
	r.PUT("/recipes/:id", UpdateRecipeHandler)
	r.DELETE("/recipes/:id", DeleteRecipeHandler)
	r.GET("/recipes/search", SearchRecipesHandler)
	r.Run()
}

// swagger:operation POST /recipes recipe createRecipe
// Create a new recipe
// ---
// produces:
// - apllication/json
// responses:
//   '201':
//     description: Successful operation
func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	err := c.ShouldBindJSON(&recipe)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	id := xid.New().String()
	publishedAt := time.Now()
	recipe.ID = id
	recipe.PublishedAt = publishedAt
	Recipes = append(Recipes, recipe)
	c.Writer.Header().Set("Location", fmt.Sprintf("/recipes/%s", recipe.ID))
	c.JSON(http.StatusCreated, map[string]interface{}{
		"recipe": recipe,
	})
}

// swagger:operation GET /recipes recipe listRecipe
// Returns list of recipes
// ---
// produces:
// - application/json
// responses:
//   '200':
//     description: Successful operation
func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"data": Recipes,
	})
}

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
func DeleteRecipeHandler(c *gin.Context) {
	id := c.Params.ByName("id")
	registryIndex := -1
	for index, recipe := range Recipes {
		if recipe.ID == id {
			registryIndex = index
			Recipes = append(Recipes[:index], Recipes[index+1:]...)
		}
	}
	if registryIndex == -1 {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

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
func UpdateRecipeHandler(c *gin.Context) {
	id := c.Params.ByName("id")
	var recipe Recipe
	err := c.ShouldBindJSON(&recipe)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	recipeIndex := -1
	for index := range Recipes {
		if id == Recipes[index].ID {
			recipeIndex = index
			r := &Recipes[index]
			r.Name = recipe.Name
			r.Ingredients = recipe.Ingredients
			r.Tags = recipe.Tags
			r.Instructions = recipe.Instructions
			break
		}
	}
	if recipeIndex == -1 {
		c.Writer.WriteHeader(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, Recipes[recipeIndex])
}

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
func SearchRecipesHandler(c *gin.Context) {
	searchedTag := c.Query("tag")
	found := []*Recipe{}
	for _, recipe := range Recipes {
		for _, tag := range recipe.Tags {
			if strings.EqualFold(searchedTag, tag) {
				found = append(found, &recipe)
				log.Printf("%v", recipe)
				break
			}
		}
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"data": found,
	})
}

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publisehdAt"`
}
