package server

import (
	"back/search"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gofiber/fiber/v2"
)

type Query struct {
	Query string `json:"query"`
}

func Run(es *elasticsearch.Client) {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		c.Set("Access-Control-Allow-Credentials", "true")
		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusOK)
		}
		return c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Server is running!")
	})

	app.Post("/api/search", func(c *fiber.Ctx) error {
		query := &Query{}

		if err := c.BodyParser(query); err != nil {
			return err
		}
		fmt.Println("Post method with query:", query.Query)

		resp := search.Search(query.Query, es)

		return c.JSON(resp)
	})

	log.Fatal(app.Listen(":3000"))
}
