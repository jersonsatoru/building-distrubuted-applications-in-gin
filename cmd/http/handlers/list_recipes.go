package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jersonsatoru/building-distributed-applications-in-gin/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

// swagger:operation GET /recipes recipe listRecipe
// Returns list of recipes
// ---
// produces:
// - application/json
// responses:
//   '200':
//     description: Successful operation
func (app *RecipeHandler) ListRecipesHandler(c *gin.Context) {
	val, err := app.redisClient.Get("recipes").Result()
	if err == redis.Nil {
		cur, err := app.collection.Find(app.ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		defer cur.Close(app.ctx)
		recipes := make([]models.Recipe, 0)
		for cur.Next(app.ctx) {
			var recipe models.Recipe
			if err := cur.Decode(&recipe); err != nil {
				c.JSON(http.StatusInternalServerError, err)
				return
			}
			recipes = append(recipes, recipe)
		}

		data, err := json.Marshal(recipes)
		if err != nil {
			c.JSON(http.StatusOK, err)
		}
		status := app.redisClient.Set("recipes", data, time.Second*30)
		log.Println(status.Val())
		c.JSON(http.StatusOK, gin.H{
			"data": recipes,
		})
	} else if err != nil {
		c.JSON(http.StatusOK, err)
	} else {
		recipes := make([]models.Recipe, 0)
		err := json.Unmarshal([]byte(val), &recipes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": recipes,
		})
	}
}
