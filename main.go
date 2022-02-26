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

func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"data": Recipes,
	})
}

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
