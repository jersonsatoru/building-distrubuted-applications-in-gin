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
	recipeHandler := NewRecipeHandler(context.TODO(), db, redisClient)
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	srv := http.Server{
		Addr:    port,
		Handler: recipeHandler.RecipeRoutes(),
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
