package main

import (
	"bufio"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"io/ioutil"
	"log"
	"os"
	"pictrick/mongo"
	"time"
)

const version = "0.0.0"

func main() {
	app := fiber.New(fiber.Config{
		ServerHeader:                 "Pictrick",
		DisableStartupMessage:        true,
		AppName:                      "Pictrick",
	})

	err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	mongo.Init()

	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/:oid", Get)

	app.Use(limiter.New(limiter.Config{
		Max:          1,
		Expiration:   1 * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
	}))

	app.Post("/store", Store)



	artfile, err := os.Open("art.txt")

	if err == nil {
		art, err := ioutil.ReadAll(bufio.NewReader(artfile))
		if err == nil {
			fmt.Println(string(art))
		}
	}

	fmt.Println("Version:", version)
	fmt.Println("Listening on:", Config.ListenOn)

	if Config.TLS {
		panic(app.ListenTLS(Config.ListenOn, Config.TLSPem, Config.TLSKey))
	} else {
		panic(app.Listen(Config.ListenOn))
	}
}