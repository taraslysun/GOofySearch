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
