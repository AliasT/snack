package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"
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

	// templates
	r.LoadHTMLGlob("templates/*")

	// assets
	r.Static("/public", "./public")

	// routes
	r.GET("/", func(c *gin.Context) {
		key := GetStoreKey()
		stringifiedPersons := client.SMembers(key).Val()
		var persons []Person
		for _, person := range stringifiedPersons {
			var p Person
			err := msgpack.Unmarshal([]byte(person), &p)
			if err == nil {
				persons = append(persons, p)
			}
		}

		c.HTML(200, "index.html", gin.H{
			"persons": persons,
		})
	})

	r.GET("/members/list", func(c *gin.Context) {
		// TODO:
	})

	r.POST("/toggle", func(c *gin.Context) {
		var person Person
		key := GetStoreKey()
		if c.ShouldBind(&person) == nil {
			stringifiedPerson, _ := msgpack.Marshal(&person)
			client.SAdd(key, stringifiedPerson)
		}

		c.JSON(200, gin.H{
			"message": key,
			"pong":    client.SMembers(key).Val(),
		})
	})

	r.Run(":8080")
}

// GetStoreKey
func GetStoreKey() string {
	// 获取当前的日期
	now := time.Now()

	// 使用当前日期作为redis key
	key := fmt.Sprintf("%d-%02d-%02d", now.Year(), now.Month(), now.Day())

	return key
}
