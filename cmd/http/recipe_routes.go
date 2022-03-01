package main

import "github.com/gin-gonic/gin"

func (app *RecipeHandler) RecipeRoutes() *gin.Engine {
	r := gin.Default()
	r.POST("/recipes", app.NewRecipeHandler)
	r.GET("/recipes", app.ListRecipesHandler)
	r.PUT("/recipes/:id", app.UpdateRecipeHandler)
	r.DELETE("/recipes/:id", app.DeleteRecipeHandler)
	r.GET("/recipes/search", app.SearchRecipesHandler)
	return r
}
