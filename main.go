package main

import (
	"errors"
	"fmt"

	"github.com/erdemkosk/golang-x-kafka/internal"
	"github.com/erdemkosk/golang-x-kafka/pkg/models"
	"github.com/erdemkosk/golang-x-kafka/pkg/repositories"
	"github.com/erdemkosk/golang-x-kafka/pkg/services"
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New()

	app.Use(recover.New())
	app.Use(cors.New())

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	tweetService := services.CreateTweetService(repositories.NewRedis[models.Tweet](rdb))
	followerService := services.CreateFollowerService()
	timelineService := services.CreateTimelineService(rdb)
	kafkaWriter, closeWriter := internal.NewWriter[models.Tweet]("127.0.0.1:9092", "twitter.newTweets")
	defer closeWriter()

	// @Summary Get the home page
	// @Description Retrieve the home page of the application
	// @ID get-home-page
	// @Produce plain
	// @Success 200 {string} string "Hello, World!"
	// @Router / [get]
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// @Summary Get a user's tweet
	// @Description Retrieve a tweet of the specified user
	// @ID get-user-tweet
	// @Produce json
	// @Param uid path string true "Tweet ID"
	// @Success 200 {object} Tweet
	// @Failure 404 {string} string "Tweet not found"
	// @Router /tweet/{uid} [get]
	app.Get("/tweet/:uid", func(c *fiber.Ctx) error {
		uid := c.Params("uid")

		tweet, err := tweetService.Get(c.Context(), uid)

		if errors.Is(err, redis.Nil) {
			return c.Status(fiber.StatusNotFound).SendString("tweet not found")
		} else if err != nil {
			return c.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
		}
		return c.JSON(tweet)
	})

	// @Summary Create a new tweet
	// @Description Create a new tweet by the user
	// @ID create-new-tweet
	// @Accept json
	// @Produce json
	// @Param request body Tweet true "Tweet data"
	// @Success 200 {object} Tweet
	// @Failure 400 {string} string "Invalid request"
	// @Router /tweet [post]
	app.Post("/tweet", func(c *fiber.Ctx) error {
		payload := models.Tweet{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		tweet, err := tweetService.Save(c.Context(), payload)
		if err != nil {
			return c.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
		}

		kafkaWriter.WriteBatch(c.Context(), tweet)

		followers, _ := followerService.Followers(c.Context(), tweet.Author)
		for _, follower := range followers {
			if err := timelineService.Push(c.Context(), follower, tweet.UID); err != nil {
				return c.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
			}
		}

		return c.JSON(tweet)
	})

	// @Summary Get user's timeline
	// @Description Retrieve the timeline of the specified user
	// @ID get-user-timeline
	// @Produce plain
	// @Param user path string true "Username"
	// @Success 200 {string} string "User's timeline"
	// @Failure 400 {string} string "Bad request"
	// @Router /timeline/{user} [get]
	app.Get("/timeline/:user", func(c *fiber.Ctx) error {
		user := c.Params("user")

		tweetIDs, err := timelineService.Latest(c.Context(), user, 10)

		if errors.Is(err, redis.Nil) {
			return c.Status(fiber.ErrBadRequest.Code).SendString("timeline not found")
		} else if err != nil {
			return c.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
		}
		tweets, err := tweetService.MGet(c.Context(), tweetIDs...)
		if err != nil {
			return c.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
		}
		timeline := ""
		for i := len(tweets) - 1; i >= 0; i-- {
			tweet := tweets[i]
			timeline += fmt.Sprintf("%s: %s\n________________\n", tweet.Author, tweet.Tweet)
		}

		return c.SendString(timeline)
	})

	app.Listen(":3000")
}
