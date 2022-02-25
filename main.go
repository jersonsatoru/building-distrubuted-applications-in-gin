package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

var (
	Recipes = make([]Recipe, 0)
)

func main() {
	r := gin.Default()
	r.GET("/recipes", NewRecipeHandler)
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
	c.JSON(http.StatusCreated, recipe)
}

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publisehdAt"`
}
