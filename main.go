package main

import (
	"fmt"

	"github.com/erdemkosk/golang-x-kafka/pkg/models"
	"github.com/erdemkosk/golang-x-kafka/pkg/repositories"
	"github.com/erdemkosk/golang-x-kafka/pkg/services"
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	tweetService := services.NewSaveTweet(repositories.NewRedis[models.Tweet](rdb))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/tweet/:uid", func(c *fiber.Ctx) error {
		fmt.Fprintf(c, "%s\n", c.Params("uid"))
		return nil
	})

	app.Post("/tweet", func(c *fiber.Ctx) error {
		payload := models.Tweet{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		tweet, err := tweetService.Save(c.Context(), payload)
		if err != nil {
			return c.SendString(err.Error())
		}

		return c.JSON(tweet)

	})

	app.Listen(":3000")
}
