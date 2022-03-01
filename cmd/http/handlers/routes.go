package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jersonsatoru/building-distributed-applications-in-gin/cmd/http/middlewares"
)

func (app *RecipeHandler) RecipeRoutes() *gin.Engine {
	r := gin.Default()
	recipeRouter := r.Group("/api/recipe/v1")
	{
		recipeRouter.GET("/recipes", app.ListRecipesHandler)
	}
	authRouter := recipeRouter.Group("/")
	authRouter.Use(middlewares.Auth())
	{
		authRouter.POST("/recipes", app.NewRecipeHandler)
		authRouter.PUT("/recipes/:id", app.UpdateRecipeHandler)
		authRouter.DELETE("/recipes/:id", app.DeleteRecipeHandler)
		authRouter.GET("/recipes/search", app.SearchRecipesHandler)
	}
	return r
}
