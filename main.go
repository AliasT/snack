package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"time"
)

type Person struct {
	Name string `form:"name"`
}

func main() {
	// gin server
	r := gin.Default()

	// redis
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// routes

	r.GET("/members/list", func(c *gin.Context) {
		// TODO:
	})

	r.POST("/toggle", func(c *gin.Context) {
		// 获取当前的日期
		now := time.Now()

		// 使用当前日期作为redis key
		key := fmt.Sprintf("%d-%02d-%02d", now.Year(), now.Month(), now.Day())

		var person Person
		if c.ShouldBind(&person) == nil {
			client.SAdd(key, person.Name).Result()
		} else {
			fmt.Println("+error+")
		}

		c.JSON(200, gin.H{
			"message": key,
			"pong":    client.SMembers(key).Val(),
		})
	})

	r.Run(":8080")
}