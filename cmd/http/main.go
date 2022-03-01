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
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jersonsatoru/building-distributed-applications-in-gin/cmd/http/handlers"
	"github.com/jersonsatoru/building-distributed-applications-in-gin/cmd/http/middlewares"
	"github.com/jersonsatoru/building-distributed-applications-in-gin/internal/database"
)

func main() {
	db, err := database.GetMongoCollection(os.Getenv("MONGO_URI"))
	if err != nil {
		log.Fatal(err)
	}
	redisClient, err := database.GetRedisConnection(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	if err != nil {
		log.Fatal(err)
	}
	recipeHandler := handlers.NewRecipeHandler(context.TODO(), db, redisClient)
	authHandler := handlers.NewAuthHandler(context.TODO(), db)

	r := gin.Default()
	authRouter := r.Group("/api/recipe/v1")
	{
		authRouter.POST("/signin", authHandler.SignInHandler)
		authRouter.POST("/refresh", authHandler.RefreshTokenHandler)
	}
	recipeRouter := r.Group("/api/recipe/v1")
	{
		recipeRouter.GET("/recipes", recipeHandler.ListRecipesHandler)
	}
	protectedRouter := recipeRouter.Group("/")
	protectedRouter.Use(middlewares.Auth())
	{
		protectedRouter.POST("/recipes", recipeHandler.NewRecipeHandler)
		protectedRouter.PUT("/recipes/:id", recipeHandler.UpdateRecipeHandler)
		protectedRouter.DELETE("/recipes/:id", recipeHandler.DeleteRecipeHandler)
		protectedRouter.GET("/recipes/search", recipeHandler.SearchRecipesHandler)
	}

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	srv := http.Server{
		Addr:    port,
		Handler: r,
	}
	shutdownErr := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownErr <- err
		}
	}()
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	err = <-shutdownErr
	if err != nil {
		log.Fatal(err)
	}
}
