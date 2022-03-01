package database

import (
	"fmt"

	"github.com/go-redis/redis"
)

func GetRedisConnection(host, port string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "",
		DB:       0,
	})
	status := client.Ping()
	if status.Val() == "PONG" {
		return client, nil
	}
	return nil, fmt.Errorf("redis ping failed")
}
