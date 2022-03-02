// 	Recipes API
// 	This is a simple recipes API. You can find out more about the API at  google, thanks
//	Schemes: https
//	Host: http://api.recipes.io:44004
//	BasePath:
//	Version: 1.0.0
//	Contact: jersonsatoru@yahoo.com.br
//  SecurityDefinitions:
//    api_key:
//      type: apiKey
//      name: Authorization
//      in: header
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
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
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

	sha := sha256.New()
	password := sha.Sum([]byte("234234"))
	db.Collection("auth").InsertOne(context.TODO(), map[string]string{
		"username": "jersonsatoru",
		"password": string(password),
	})

	redisClient, err := database.GetRedisConnection(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	if err != nil {
		log.Fatal(err)
	}
	recipeHandler := handlers.NewRecipeHandler(context.TODO(), db.Collection("recipes"), redisClient)
	authHandler := handlers.NewAuthHandler(context.TODO(), db.Collection("auth"))

	r := gin.Default()
	host := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	store, err := redis.NewStore(10, "tcp", host, "", []byte("secret"))
	if err != nil {
		log.Fatal(err)
	}
	r.Use(sessions.Sessions("recipes_api", store))
	authRouter := r.Group("/api/recipe/v1")
	{
		authRouter.POST("/signin", authHandler.SignInHandler)
		authRouter.POST("/refresh", authHandler.RefreshTokenHandler)
		authRouter.POST("/signout", authHandler.SignoutHandler)
	}
	recipeRouter := r.Group("/api/recipe/v1")
	{
		recipeRouter.GET("/recipes", recipeHandler.ListRecipesHandler)
	}
	protectedRouter := recipeRouter.Group("/")
	protectedRouter.Use(middlewares.OAuth2Session())
	{
		protectedRouter.POST("/recipes", recipeHandler.NewRecipeHandler)
		protectedRouter.PUT("/recipes/:id", recipeHandler.UpdateRecipeHandler)
		protectedRouter.DELETE("/recipes/:id", recipeHandler.DeleteRecipeHandler)
		protectedRouter.GET("/recipes/search", recipeHandler.SearchRecipesHandler)
	}

	// port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	srv := http.Server{
		Addr:    ":44004",
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
	err = srv.ListenAndServeTLS("certs/localhost.crt", "certs/localhost.key")
	if err != nil {
		log.Fatal(err)
	}
	err = <-shutdownErr
	if err != nil {
		log.Fatal(err)
	}
}
